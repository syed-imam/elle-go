# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

- Module: `elle-go` (see `go.mod`)
- Go version: 1.27

> This is a new project. The sections below are intentionally sparse and should
> be expanded as the codebase grows. Do not invent architecture that does not
> yet exist in the code.

## Goals

Build a **serializability checker in Go**, inspired by Jepsen's
[Elle](https://github.com/jepsen-io/elle). Given a recorded history of
transactions, infer dependencies between them, build a dependency graph, and
detect cycles that prove isolation anomalies (up to serializability
violations). Motivation: production transactional DBs (e.g. Postgres) rarely
run at SERIALIZABLE isolation because it's too slow, but we still want to
*detect* serializability failures after the fact from observed histories.

MVP target: the **list-append** workload (transactions of appends + reads over
named lists), because a single list read recovers the full version order of a
key — making ww/wr/rw dependency inference tractable.

## Working mode (IMPORTANT)

This is a **learning project**. The owner is an experienced engineer (Java,
Elixir, TS/JS, Ruby, PHP) getting hands-on with Go. When making code changes:

- Keep each step **one logical, committable unit** (a small feature + its
  test) so each step is learnable.
- **No comments in code.** The owner reads code directly; do not add code
  comments (including Go doc comments). Explain concepts in chat instead.
- **Never run `git commit` or `git push`.** The owner makes ALL commits
  himself. Just work in small logical steps and say when a step is ready to
  commit — don't stage or commit anything.
- **Explain every step** and the Go-specific concepts involved; don't just ship
  large diffs. Draw analogies to Java/Elixir/TS where useful.
- Move at the pace of understanding, not completion.

## Reference

Upstream Elle (Clojure) is cloned at `../elle-upstream` (sibling of this repo,
not committed). Key namespaces: `core` (dependency graphs + cycle detection),
`list_append` (flagship workload), `graph` (graph utils), `consistency_model`
(anomaly → isolation-level lattice), `rels` (edge types: ww/wr/rw + predicate/
process/realtime variants), `txn`, `rw_register`.

## Commands

Standard Go tooling applies:

- Build: `go build ./...`
- Check a JSON history: `go run . history.json` (or pipe to stdin); exit 0 clean / 1 anomaly / 2 bad input
- Run a workload on Postgres: `go run . pg -isolation read-committed|repeatable-read|serializable` (DSN via `-dsn` or `$ELLE_PG`; disposable DB only — it creates and truncates `elle_lists`)
- Test (all): `go test ./...`
- Test (single package): `go test ./path/to/pkg`
- Test (single test): `go test ./path/to/pkg -run '^TestName$' -v`
- Elle parity (slow, JVM per history): `ELLE_PARITY=1 go test ./differential -timeout 45m`
- Live Postgres tests: `ELLE_PG=<dsn> go test ./pgrun` (skipped without it)
- Vet: `go vet ./...`
- Format: `gofmt -w .` (or `go fmt ./...`)

## Architecture

The checker is a pipeline; data flows through four packages, wrapped by a
façade:

    history → listappend → graph → anomaly      (checker.Check glues them)

- `history` — the data model. A `History` is an ordered slice of `Op`
  (transactions); each `Op` has `Mops` (micro-ops: `Append` or `Read` on a
  named key). An op's index is its transaction identity throughout the pipeline.
  Types encode to JSON as names (`"ok"`, `"append"`, `"r"`) with lowercase keys.
- `listappend` — the workload logic. `VersionOrder` recovers each key's write
  order from list reads; `Dependencies` infers the ww/wr/rw `Edge`s between
  transactions from that order. This is where list-append's "a read recovers the
  full version order of a key" property is exploited. `Validate` rejects reads
  of values no committed op appended.
- `graph` — a generic directed multigraph with **typed edges** (`Rel`, a bitset
  of ww/wr/rw). Nodes are transaction indices; the rel lives on the edge, never
  the node. Provides the cycle/anomaly primitives: `FindCycle` (returns a
  witness cycle path, not just a bool), `SCCs` (Kosaraju), and the typed-search
  building blocks `Filter`, `Reachable`, `Path`, `EdgesWith`.
- `anomaly` — classification. `Check` runs **typed search** in severity order →
  a `Verdict{Anomaly, Level, Cycle}`. `Requires` maps the anomaly to the formal
  isolation level that prevents it; `Verdict.String` renders the English verdict.

Key design decisions (rationale in `DESIGN.md`):

- **Typed search, not label-a-witness.** Each anomaly is found by directed
  search for its own edge-shape, so the verdict is order-independent. G0 = a
  cycle in the ww-only subgraph; G1c / G-single = an offending wr / rw edge
  closed by a ww/wr return `Path`; G2 = the catch-all cycle.
- **Edge types on edges, not nodes.** ww/wr/rw names a *relationship between two
  transactions*, so it can only live on an edge. Mirrors Elle's BitRels.
- **Formal isolation levels only.** `Level` is the Adya/ANSI hierarchy;
  translating to a vendor's (often mislabeled) level names is kept separate.

Around the pipeline:

- `checker` — the public façade: `Check(history) (anomaly.Verdict, error)`
  validates, builds the graph, and classifies.
- `generator` — synthetic workloads: `Serial` (valid histories), one injector
  per cycle class (`WriteCycle` G0, `CircularInformationFlow` G1c, `ReadSkew`
  G-single, `WriteSkew` G2), and `Invocations` (unfilled txns for a real DB).
- `differential` — parity harness against upstream Elle (`RunElle`,
  `TypeAgrees`); opt-in via `ELLE_PARITY=1`.
- `pgrun` — runs list-append txns on Postgres (`Exec` one txn, `Run` a worker
  pool) at a chosen isolation level; keys are rows in `elle_lists` (`int[]`).
- `main` — CLI: check a JSON history, or `pg` to run a workload and verdict it.
