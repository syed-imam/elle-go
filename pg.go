package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"elle-go/anomaly"
	"elle-go/checker"
	"elle-go/generator"
	"elle-go/history"
	"elle-go/pgrun"
)

func runPG(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pg", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dsn := fs.String("dsn", os.Getenv("ELLE_PG"), "Postgres DSN of a disposable database (default $ELLE_PG)")
	isoName := fs.String("isolation", "serializable", "read-committed, repeatable-read, or serializable")
	txns := fs.Int("txns", 300, "number of transactions")
	workers := fs.Int("workers", 8, "concurrent clients")
	keys := fs.String("keys", "x,y,z", "comma-separated keys")
	mops := fs.Int("mops", 4, "max operations per transaction")
	seed := fs.Int64("seed", 1, "workload seed")
	shapesPath := fs.String("shapes", "", "JSON file of transaction shapes to run instead of random transactions")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	iso, err := pgrun.ParseIsolation(*isoName)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	cfg := generator.Config{Keys: strings.Split(*keys, ","), Txns: *txns, MaxMops: *mops}
	ops := generator.Invocations(cfg, *seed)
	if *shapesPath != "" {
		shapes, err := loadShapes(*shapesPath)
		if err != nil {
			fmt.Fprintln(stderr, "shapes:", err)
			return 2
		}
		ops = generator.FromShapes(shapes, *txns, *seed)
	}
	if *dsn == "" {
		fmt.Fprintln(stderr, "no database: pass -dsn or set ELLE_PG")
		return 2
	}
	db, err := sql.Open("pgx", *dsn)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	defer db.Close()
	ctx := context.Background()
	if err := pgrun.Setup(ctx, db); err != nil {
		fmt.Fprintln(stderr, "setup:", err)
		return 2
	}
	h := pgrun.Run(ctx, db, iso, ops, *workers)
	v, err := checker.Check(h)
	if err != nil {
		fmt.Fprintln(stderr, "invalid history:", err)
		return 2
	}
	fmt.Fprintf(stdout, "%s: %d/%d committed\n", *isoName, committed(h), len(h))
	fmt.Fprintln(stdout, v)
	if v.Anomaly != anomaly.None {
		return 1
	}
	return 0
}

func committed(h history.History) int {
	n := 0
	for _, op := range h {
		if op.Type == history.Ok {
			n++
		}
	}
	return n
}

func loadShapes(path string) ([]history.Op, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var shapes []history.Op
	if err := json.Unmarshal(b, &shapes); err != nil {
		return nil, err
	}
	if len(shapes) == 0 {
		return nil, fmt.Errorf("%s has no shapes", path)
	}
	return shapes, nil
}
