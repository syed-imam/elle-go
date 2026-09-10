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

- Keep edits **small** (~15 lines at a time) so each step is learnable.
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
- Run: `go run .`
- Test (all): `go test ./...`
- Test (single package): `go test ./path/to/pkg`
- Test (single test): `go test ./path/to/pkg -run '^TestName$' -v`
- Vet: `go vet ./...`
- Format: `gofmt -w .` (or `go fmt ./...`)

## Architecture

_To be documented once source files exist._ Focus this section on the
"big picture" that requires reading multiple files to understand — not a
file-by-file listing.
