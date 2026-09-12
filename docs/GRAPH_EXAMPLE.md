# elle-go — Worked Graph Examples (local / unpublished)

Status: DRAFT · Concrete, hand-traced examples of turning a history into a
dependency graph. Companion to `DESIGN.md` and `ANOMALIES.md`. The goal here is
to *see* what a node and an edge actually are, on real data.

Reminder:
- **Node** = one committed transaction.
- **Edge** = "must come before," inferred from a same-key op pair (`ww`/`wr`/`rw`).
- A read returning `[1 2 3]` reveals the append order `1 → 2 → 3` for that key.
- Each appended value is unique, so a value pins down exactly which txn wrote it.

---

## Example A — a serializable history (graph is a DAG)

### The history

Three transactions over keys `x` and `y`. We annotate each append with the value
it wrote so we can trace who wrote what.

```
T1:  append x 1              (writes x=1)
T2:  read x [1]           ;  append x 2      (sees x=1, then writes x=2)
T3:  read x [1 2]         ;  read y []       (sees both x appends, y empty)
```

### Step 1 — recover version order per key

From the reads:
- key `x`: the read `[1 2]` tells us the append order is **1 → 2**.
  - value `1` was written by **T1**, value `2` by **T2**.
- key `y`: only an empty read; no writes; no order to recover.

### Step 2 — derive edges

Go op-pair by op-pair on each key:

| Reason | Key | Pair | Edge |
|--------|-----|------|------|
| T1 wrote `1`, T2 wrote the next version `2` | x | write→write | **T1 →ww→ T2** |
| T2 read `[1]` — it witnessed T1's write of `1` | x | write→read | **T1 →wr→ T2** |
| T3 read `[1 2]` — witnessed T2's write of `2` | x | write→read | **T2 →wr→ T3** |
| T3 read `[1]`-prefix state then T2 wrote `2`? No — T3 saw `2`, so no rw here | x | — | (none) |

(T1→T2 has both a ww and a wr; that's fine — two reasons, same direction.)

### Step 3 — the graph

```
        ww, wr              wr
  (T1) ─────────► (T2) ─────────► (T3)
```

All edges point "forward." No path leads back. **Acyclic → it's a DAG.**

### Step 4 — read off a serial order

A DAG can be topologically sorted. One valid order:

```
T1  →  T2  →  T3
```

Every edge is satisfied (all point the same way as the order). A serial
execution `T1; T2; T3` explains the whole history → **serializable ✅**.

---

## Example B — write skew (graph has a cycle)

The motivating anomaly. Two transactions, two keys, both start empty.

### The history

```
T1:  read y []           ;  append x 1      (looks at y, sees nothing, writes x)
T2:  read x []           ;  append y 1      (looks at x, sees nothing, writes y)
```

Both committed. Final state: `x=[1]`, `y=[1]`.

### Step 1 — recover version order per key

- key `x`: value `1` written by **T1**. T2 read `x` as `[]` (empty).
- key `y`: value `1` written by **T2**. T1 read `y` as `[]` (empty).

### Step 2 — derive edges (the key move)

The interesting edges are **anti-dependencies (`rw`)** — a reader who saw the
*old* (empty) state must come before the writer who filled it:

| Reason | Key | Pair | Edge |
|--------|-----|------|------|
| T2 read `x`=`[]` (before `1` existed); T1 later appended `1` | x | read→write | **T2 →rw→ T1** |
| T1 read `y`=`[]` (before `1` existed); T2 later appended `1` | y | read→write | **T1 →rw→ T2** |

### Step 3 — the graph

```
            rw  (on key y: T1 read old y, T2 wrote it)
        ┌───────────────────────────┐
        ▼                           │
      (T1)                        (T2)
        │                           ▲
        └───────────────────────────┘
            rw  (on key x: T2 read old x, T1 wrote it)
```

Flattened, the edges are:

```
  T1 ─rw(y)─► T2
  T2 ─rw(x)─► T1
```

### Step 4 — try to find a serial order

```
  T1 → T2   (from key y)   "T1 must come before T2"
  T2 → T1   (from key x)   "T2 must come before T1"
```

Contradiction: T1 before T2 **and** T2 before T1. Follow the arrows and you loop
forever: `T1 → T2 → T1 → T2 → …`. **Cycle → not a DAG → no serial order exists.**

### Step 5 — classify

The cycle `T1 → T2 → T1` is made of **two `rw` edges** → **G2-item (write
skew)** → violates **Serializable** (but Snapshot Isolation would have allowed
it). See `ANOMALIES.md §2`.

---

## The contrast, side by side

```
  Example A (serializable)          Example B (write skew)

  T1 ──► T2 ──► T3                   T1 ──► T2
                                      ▲      │
  no way back  →  DAG                 └──────┘   loop  →  cycle
  serial order: T1;T2;T3             no serial order exists
```

Same machinery both times: build nodes (transactions), draw `ww`/`wr`/`rw`
edges from same-key op pairs, then ask **"is there a cycle?"** That single
question — cycle or not — is the entire verdict.

---

## How this maps to the code (once built)

- Nodes ← `history.Op` values with `Type == Ok`.
- Version order per key ← recovered by `listappend` from `Read` mops (Phase 2).
- Edges ← emitted by `listappend` dependency inference (Phase 3).
- "Is there a cycle?" ← `graph` package SCC / cycle detection (Phase 4).
- Cycle shape → anomaly name ← `consistency` classification (Phase 5).

These docs describe the target; the packages don't all exist yet (`ROADMAP.md`).
