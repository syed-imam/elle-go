package pgrun

import (
	"context"
	"database/sql"
	"os"
	"reflect"
	"testing"

	"elle-go/anomaly"
	"elle-go/checker"
	"elle-go/generator"
	"elle-go/history"
)

func openTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("ELLE_PG")
	if dsn == "" {
		t.Skip("set ELLE_PG to a disposable Postgres DSN to run")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	ctx := context.Background()
	if err := Setup(ctx, db); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestExecAppendThenRead(t *testing.T) {
	db := openTestDB(t)
	op := history.Op{Process: 1, Type: history.Invoke, Mops: []history.Mop{
		{Type: history.Read, Key: "x"},
		{Type: history.Append, Key: "x", App: 1},
		{Type: history.Append, Key: "x", App: 2},
		{Type: history.Read, Key: "x"},
	}}
	got := Exec(context.Background(), db, sql.LevelSerializable, op)
	if got.Type != history.Ok {
		t.Fatalf("txn not ok: %+v", got)
	}
	if !reflect.DeepEqual(got.Mops[0].Read, []int{}) || !reflect.DeepEqual(got.Mops[3].Read, []int{1, 2}) {
		t.Fatalf("reads = %v, %v", got.Mops[0].Read, got.Mops[3].Read)
	}
}

func TestRunSerializableIsClean(t *testing.T) {
	db := openTestDB(t)
	cfg := generator.Config{Keys: []string{"x", "y", "z"}, Txns: 300, MaxMops: 4}
	h := Run(context.Background(), db, sql.LevelSerializable, generator.Invocations(cfg, 1), 8)
	if len(h) != cfg.Txns {
		t.Fatalf("got %d completed ops, want %d", len(h), cfg.Txns)
	}
	v, err := checker.Check(h)
	if err != nil {
		t.Fatal(err)
	}
	if v.Anomaly != anomaly.None {
		t.Fatalf("SERIALIZABLE produced an anomaly: %v", v)
	}
}
