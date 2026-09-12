# elle-go — Anomaly Reference (local / unpublished)

Status: DRAFT · A revisitable catalogue of the isolation anomalies elle-go aims
to detect. Companion to `DESIGN.md` (the *why/how*) — this doc is the *what*.

The framing throughout: nodes are transactions, edges are ordering constraints
(`ww` / `wr` / `rw`). Most anomalies are **a cycle of a particular shape** in
that graph. A few are **integrity violations** caught without the graph.

Notation for examples: version order of a key is recovered from a list read,
e.g. a read returning `[1 2 3]` means appends happened in the order `1 → 2 → 3`.
`Tn` = transaction n.

---

## 1. The core idea (recap)

A history is serializable **iff** its dependency graph is acyclic. Every anomaly
below is either:

- a **cycle** — ordering constraints contradict each other (no serial order), or
- an **integrity violation** — the history is internally impossible regardless
  of ordering (e.g. reading a value nobody committed).

The *shape* of a cycle (which edge types it contains) determines which isolation
level it violates. That mapping is the anomaly hierarchy.

---

## 2. The dependency-graph anomalies (Adya's G-hierarchy)

Ordered weakest → strongest. Each level *includes* the ones above it: prohibiting
G2 means you also prohibit G0, G1, G-single.

| Code | Name | Cycle shape | Violates (first level that forbids it) |
|------|------|-------------|----------------------------------------|
| **G0** | Dirty Write | cycle of **only `ww`** edges | Read Uncommitted |
| **G1a** | Aborted Read | *(integrity, not a cycle)* | Read Committed |
| **G1b** | Intermediate Read | *(integrity, not a cycle)* | Read Committed |
| **G1c** | Circular Information Flow | cycle of `ww` + `wr` (**no `rw`**) | Read Committed |
| **G-single** | Single Anti-dependency | cycle with **exactly one `rw`** | Snapshot Isolation / Repeatable Read |
| **G2-item** | Item Anti-dependency | cycle with **≥1 `rw`** (over items) | Serializable |
| **G2** | Anti-dependency (general) | cycle with **≥1 `rw`** (incl. predicates) | Serializable |

### G0 — Dirty Write
Two transactions' writes to the same keys are interleaved so their version
orders disagree across keys.

```
On key x:  T1 → T2   (T1's append precedes T2's)
On key y:  T2 → T1   (T2's append precedes T1's)
```
Cycle `T1 → T2 → T1` made entirely of `ww` edges. Means writes were not applied
atomically as a unit. Even Read Uncommitted forbids this.

### G1a — Aborted Read
A transaction reads a value written by a transaction that **later aborted**
(`Type == Fail`). The value never officially existed.

```
T1: append 5 to x        then FAILS (aborts)
T2: read x → [.. 5 ..]   saw a value from a doomed transaction
```
Not a cycle — an integrity violation. We detect it by checking that every value
observed in a read was written by a committed (`Ok`) transaction.

### G1b — Intermediate Read
A transaction reads a value that was **not the final value** a writer left for
that key — i.e. an intermediate state within another transaction.

```
T1: append 1 to x ; append 2 to x     (x's final state from T1 ends in 2)
T2: read x → [1]                        saw the mid-transaction state, missing 2
```
Only the *last* write a transaction makes to a key should be visible to others.
Reading an earlier one is G1b. Also an integrity check, not a cycle.

### G1c — Circular Information Flow
A cycle using only `ww` and `wr` edges (no anti-dependencies). Information flows
in a loop: each txn in the cycle observed or overwrote the next.

```
T1 → T2   (T2 read what T1 wrote — wr)
T2 → T1   (T1 read what T2 wrote — wr)
```
Two transactions each read the other's write → impossible to order. Forbidden
at Read Committed and above.

### G-single — Single Anti-dependency ("read skew")
A cycle with **exactly one** `rw` edge; the rest are `ww`/`wr`. This is the
classic **read skew**: you read one key before an update and another key after.

```
Initial: x=[1], y=[1]
T1: read x → [1]        (before T2)      ... read y → [1 2]  (after T2)
T2: append 2 to x ; append 2 to y ; commits
```
T1 saw x's *old* value but y's *new* value — a snapshot that never existed as a
single point in time. Allowed by Read Committed; forbidden by Snapshot
Isolation / Repeatable Read.

### G2-item — Item Anti-dependency ("write skew")
A cycle with **one or more `rw`** edges over individual items (keys). The
textbook **write skew**:

```
Initial: x=[], y=[]
T1: read y → []  ; append 1 to x       (decides based on y, writes x)
T2: read x → []  ; append 1 to y       (decides based on x, writes y)
```
- On y: T1 read old → T2 appended later ⇒ `T1 → T2` (rw)
- On x: T2 read old → T1 appended later ⇒ `T2 → T1` (rw)

Cycle `T1 → T2 → T1` with two `rw` edges. Each read a stale value the other then
invalidated. Allowed by Snapshot Isolation (!); forbidden only by Serializable.
This is the anomaly Postgres `SERIALIZABLE` rejects but `REPEATABLE READ` (SI)
permits — the motivating case for the whole project.

### G2 — General Anti-dependency
Same as G2-item but the anti-dependency may come from a **predicate** read
(e.g. "read all rows where …") rather than a single item. Prohibiting G1 + G2 is
the definition of **full Serializability (PL-3)**. Predicate reads are a
non-goal for the MVP (see `DESIGN.md §9`), so elle-go targets G2-item first.

---

## 3. Classic named anomalies → G-hierarchy

The ANSI/Berenson ("A Critique of ANSI SQL Isolation Levels") phenomena, mapped
onto the dependency-graph codes above.

| Classic name | Description | Corresponds to |
|--------------|-------------|----------------|
| **Dirty Write** (P0) | write over another txn's uncommitted write | G0 |
| **Dirty Read** (P1) | read another txn's uncommitted / aborted write | G1a (aborted flavour) |
| **Non-repeatable / Fuzzy Read** (P2) | re-reading a key yields a changed value | G-single |
| **Phantom** (P3) | a predicate re-read gains/loses rows | G2 (predicate) |
| **Lost Update** (P4) | two txns read-modify-write, one update vanishes | G-single / G2-item |
| **Read Skew** (A5A) | read keys across an update boundary | G-single |
| **Write Skew** (A5B) | disjoint writes on correlated reads | G2-item |

---

## 4. Isolation-level lattice

Which anomalies each level *permits* (✓ = can occur, ✗ = prohibited):

| Anomaly \ Level | Read Uncommitted | Read Committed | Snapshot / Rep. Read | Serializable |
|-----------------|:----------------:|:--------------:|:--------------------:|:------------:|
| G0 (dirty write) | ✗ | ✗ | ✗ | ✗ |
| G1a/G1b/G1c | ✓ | ✗ | ✗ | ✗ |
| G-single (read skew) | ✓ | ✓ | ✗ | ✗ |
| G2-item (write skew) | ✓ | ✓ | ✓ | ✗ |
| G2 (predicate) | ✓ | ✓ | ✓ | ✗ |

Reading down a column tells you what a database at that level still lets slip
through. Serializable is the only column that's all ✗.

---

## 5. Integrity anomalies (no cycle needed)

Elle also checks that the history is even *coherent* before graph analysis.
These are list-append specific and catch broken data, not just bad ordering:

- **Internal inconsistency** — within one transaction, a read doesn't reflect
  that same transaction's own earlier appends to the key. A txn must see its own
  writes.
- **Aborted read (G1a)** — reading an element from a `Fail` transaction.
- **Intermediate read (G1b)** — reading a non-final element from a writer.
- **Duplicate elements** — the same value appears twice in a list read. Breaks
  the "appends are unique" assumption that makes version recovery unambiguous
  (`DESIGN.md §4`).
- **Incompatible order** — two reads of the same key return orders that can't
  both be prefixes of one true order, e.g. `[1 2]` and `[2 1]`.
- **Cyclic version order** — the recovered per-key version order itself contains
  a cycle (a append precedes b and b precedes a on the same key).

---

## 6. What elle-go targets (MVP)

In scope for the first working checker (`ROADMAP.md`):
- **G0, G1c, G-single, G2-item** via cycle detection.
- **G1a, G1b, internal, duplicate/incompatible/cyclic** via integrity checks.

Deferred (see `DESIGN.md §9`):
- **G2 predicate** anomalies (needs predicate reads).
- **process / realtime** edge variants (`G0-process`, `G-single-realtime`, …).

---

## 7. References

- Adya, "Weak Consistency" (PhD thesis, MIT, 1999) — the G-hierarchy.
- Berenson et al., "A Critique of ANSI SQL Isolation Levels" (1995) — P/A phenomena.
- Kingsbury & Alvaro, "Elle: Inferring Isolation Anomalies…" (VLDB 2020).
- Upstream Elle `consistency_model` namespace — the anomaly ↔ level lattice in code.
