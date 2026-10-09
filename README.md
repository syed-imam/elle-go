# elle-go

A serializability checker for transactional databases, written in Go and
inspired by Jepsen's [Elle](https://github.com/jepsen-io/elle).

You give it a recorded history of transactions. It works out the dependencies
between them, builds a dependency graph, and looks for cycles that **prove** an
isolation anomaly, from dirty writes up to serializability violations. When it
finds one, it names the anomaly, gives the transactions involved, and says which
isolation level would have prevented it.

```
$ elle-go pg -isolation repeatable-read
repeatable-read: 157/300 committed
G2: transactions 256 → 254 → 256 (requires serializable or stronger)
```

## Why

Most production databases don't run at `SERIALIZABLE` because it costs
throughput. Weaker levels allow anomalies such as write skew, read skew and
lost updates, which usually go unnoticed until the data is wrong. elle-go
checks after the fact: run your transactions, record what each one read and
wrote, and get either a proof of an anomaly or a clean result.

## How it works

elle-go uses the **list-append** workload. Each key holds a list, transactions
append to it or read it, and every appended value is unique. A single read such
as `[1 2 3]` therefore shows the full order in which writes to that key landed.
From that order elle-go infers three kinds of edges between transactions:

- `ww`: T1's write was overwritten by T2's
- `wr`: T2 read T1's write
- `rw`: T1 read a version that T2's write later replaced

A history is serializable only if this graph has no cycle. The edge types in a
cycle show how bad the anomaly is (Adya's G-hierarchy). elle-go searches for
each cycle shape directly, from most to least severe, so the verdict doesn't
depend on which cycle happens to be found first.

## What it detects

| Anomaly | What happened | Prevented by |
|---|---|---|
| `G0` | dirty write: a cycle of only `ww` edges | read uncommitted |
| `G1c` | circular information flow: a `ww`/`wr` cycle | read committed |
| `incompatible-order` | two reads of a key disagree on the order of its appends | read committed |
| `G-single` | read skew: a cycle with exactly one `rw` edge | snapshot isolation |
| `internal` | a transaction's read contradicts its own earlier reads or appends | snapshot isolation |
| `lost-update` | two transactions read the same value of a key, then both wrote it | snapshot isolation |
| `G2` | write skew and other cycles with `rw` edges | serializable |

Levels are the formal Adya / ANSI ones. Vendors don't always use these names
the same way.

## Install

```sh
go install github.com/syed-imam/elle-go@latest
```

Or use the Docker image, which bundles elle-go with Postgres 16:

```sh
docker pull ghcr.io/syed-imam/elle-go:latest
```

## Usage

### Check a recorded history

```sh
elle-go history.json      # or: elle-go < history.json
```

A history is a JSON array of completed transactions. `type` is `ok`, `fail` or
`info` (outcome unknown). Each micro-op is an `append` of a unique integer or an
`r` (read) with the list it observed:

```json
[
  {"process": 1, "type": "ok", "mops": [{"type": "r", "key": "y", "read": []}, {"type": "append", "key": "x", "app": 1}]},
  {"process": 2, "type": "ok", "mops": [{"type": "r", "key": "x", "read": []}, {"type": "append", "key": "y", "app": 2}]}
]
```

```
G2: transactions 0 → 1 → 0 (requires serializable or stronger)
```

Exit codes: `0` no anomaly, `1` anomaly found, `2` invalid input.

### Run a workload against Postgres

```sh
export ELLE_PG='postgres://postgres:elle@localhost:5432/postgres?sslmode=disable'
elle-go pg -isolation read-committed     # or repeatable-read, serializable
```

> **Use a disposable database.** `pg` creates and truncates a table named
> `elle_lists`.

This runs concurrent random list-append transactions (flags: `-txns`, `-workers`,
`-keys`, `-mops`, `-seed`), records the history and checks it. On Postgres 16:

| Isolation | Result |
|---|---|
| read committed | `G-single` |
| repeatable read | `G2` (write skew) |
| serializable | no anomaly |

### Check your own transaction shapes

`-shapes` runs your own transaction types in place of random ones. Here are two
on-call doctors, each checking that the other is still on call before going off:

```json
[
  {"mops": [{"type": "r", "key": "doctor:alice"}, {"type": "r", "key": "doctor:bob"}, {"type": "append", "key": "doctor:alice"}]},
  {"mops": [{"type": "r", "key": "doctor:alice"}, {"type": "r", "key": "doctor:bob"}, {"type": "append", "key": "doctor:bob"}]}
]
```

```sh
elle-go pg -isolation repeatable-read -shapes examples/write-skew.json
```

```
repeatable-read: 92/300 committed
G2: transactions 296 → 295 → 296 (requires serializable or stronger)
```

### Claude Code plugin

The `check-tx` skill takes the transaction code you just wrote, extracts its
read/write shapes for you to review, and runs them against a throwaway Postgres
in Docker. The only requirement is Docker.

```
/plugin marketplace add syed-imam/elle-go
/plugin install elle-go@elle-go
/elle-go:check-tx
```

## Agreement with Elle

The `differential` package runs the same histories through upstream Elle and
compares the results. This covers hand-written fixtures, random serial
histories, generated histories for each anomaly class, and real Postgres
histories at all three isolation levels. They currently all agree:
- both tools find a violation in the same histories;
- the anomaly elle-go reports is one of the anomalies Elle reports.

```sh
ELLE_PARITY=1 go test ./differential -timeout 45m   # needs lein and ../elle-upstream
```

## Limitations

- **List-append only.** Read/write registers and other data types aren't
  supported.
- **Item-level only.** Predicate reads (`WHERE`, `COUNT(*)`, ranges) and
  phantoms aren't modeled.
- **No locking.** `pg` and `-shapes` don't take explicit locks, so code that
  relies on `SELECT … FOR UPDATE` will be reported as anomalous.
- **One verdict per history:** the most severe anomaly found, not every one.
- **No aborted or intermediate reads.** A read of a value that no committed or
  in-doubt transaction appended is rejected as invalid input, not reported as
  `G1a`. Intermediate reads (`G1b`) aren't detected.
- **No real-time or per-process order.** Strict serializability isn't checked.

## Development

```sh
go test ./...                     # unit tests
ELLE_PG=<dsn> go test ./pgrun     # live Postgres tests
go vet ./...
```

## License

[Apache-2.0](LICENSE). elle-go is an independent implementation. It uses no
code from Elle, which is distributed under the EPL.
