package differential

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/syed-imam/elle-go/anomaly"
	"github.com/syed-imam/elle-go/generator"
	"github.com/syed-imam/elle-go/pgrun"
)

func TestParityPostgres(t *testing.T) {
	if os.Getenv("ELLE_PARITY") == "" {
		t.Skip("set ELLE_PARITY=1 to run the Elle differential test (slow: JVM per history)")
	}
	if !ElleAvailable() {
		t.Skip("lein or ../../elle-upstream not available; skipping Elle differential test")
	}
	dsn := os.Getenv("ELLE_PG")
	if dsn == "" {
		t.Skip("set ELLE_PG to a disposable Postgres DSN to run")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	cfg := generator.Config{Keys: []string{"x", "y", "z"}, Txns: 200, MaxMops: 4}
	for _, name := range []string{"read-committed", "repeatable-read", "serializable"} {
		iso, err := pgrun.ParseIsolation(name)
		if err != nil {
			t.Fatal(err)
		}
		for seed := int64(0); seed < 3; seed++ {
			if err := pgrun.Setup(ctx, db); err != nil {
				t.Fatal(err)
			}
			h := pgrun.Run(ctx, db, iso, generator.Invocations(cfg, seed), 8)
			elle, err := RunElle(h)
			if err != nil {
				t.Fatalf("%s seed %d: RunElle: %v", name, seed, err)
			}
			ours, err := Verdict(h)
			if err != nil {
				t.Fatalf("%s seed %d: Verdict: %v", name, seed, err)
			}
			if (ours.Anomaly != anomaly.None) != elle.Violation() || !TypeAgrees(ours.Anomaly, elle.AnomalyTypes) {
				t.Errorf("%s seed %d: disagree: ours=%v vs elle valid?=%q %v",
					name, seed, ours.Anomaly, elle.ValidField, elle.AnomalyTypes)
				continue
			}
			t.Logf("%s seed %d: agree — ours=%v, elle=%v", name, seed, ours.Anomaly, elle.AnomalyTypes)
		}
	}
}
