---
name: check-tx
description: Check whether database transactions in the user's code can produce serializability anomalies (write skew, read skew, lost updates) by running their read/write shapes concurrently against a throwaway Postgres in Docker and verifying the history with elle-go. Use when the user runs /check-tx, asks whether transactional code is safe under concurrency or a given isolation level, or just wrote code that reads and then writes rows inside a transaction.
---

# check-tx

You turn the user's transactional code into **shapes**, the user confirms them, and a deterministic checker (elle-go) decides. You never decide the verdict yourself: you draft the input and explain the output.

## 1. Pick the target

- If the user passed an argument (a function name, file, or path), use it.
- Otherwise, look at transactional code in the current change (`git diff` and `git diff --staged`): functions that open a transaction, or ORM blocks such as `transaction do`, `@Transactional`, `with transaction.atomic()`, `BEGIN ... COMMIT`.
- If there are several candidates, or none, ask the user which one.
- Then find the **other transactions that touch the same rows** (grep for the same tables or models). Anomalies happen *between* transactions, so a shape set with a single transaction type is usually incomplete. Include each concurrent writer and reader of those rows.

## 2. Extract shapes

A shape is one transaction type: an ordered list of operations on **row keys**.

- Reading a row (`SELECT`, ORM `find` or `get`) → `{"type": "r", "key": "<table>:<id>"}`
- Writing a row (`UPDATE`, `INSERT`, `DELETE`, ORM `save`, `update`) → `{"type": "append", "key": "<table>:<id>"}`
- Keep the order the code performs them in.
- Use a **small set of concrete ids** so transactions collide. For example, `transfer(from, to)` becomes two shapes, `account:1 → account:2` and `account:2 → account:1`. Contention is what exposes anomalies.

Stop and tell the user the check would be unsound, rather than guessing, when the code:
- reads by **predicate or range** (`WHERE balance > 0`, `COUNT(*)`, `SELECT ... WHERE date BETWEEN`). Phantoms are not modeled.
- relies on **explicit locks** (`SELECT ... FOR UPDATE`, advisory locks). The analog does not take locks yet, so it would report anomalies the locks prevent.
- picks which rows to touch **based on values it read**, in a way that fixed ids can't represent.

## 3. Confirm with the user

Show only: the shapes table (transaction → ordered ops), the isolation level with a few words on where it came from (Postgres defaults to `read-committed`), and at most two short bullets on assumptions that could make the check wrong. Then ask: "Run it?" The user can correct the shapes or the level first. This review is the safety gate against a wrong extraction.

## 4. Run

Write the shapes to a temporary JSON file in this format:

```json
[
  {"mops": [{"type": "r", "key": "doctor:alice"}, {"type": "r", "key": "doctor:bob"}, {"type": "append", "key": "doctor:alice"}]},
  {"mops": [{"type": "r", "key": "doctor:alice"}, {"type": "r", "key": "doctor:bob"}, {"type": "append", "key": "doctor:bob"}]}
]
```

Then run:

```
${CLAUDE_PLUGIN_ROOT}/skills/check-tx/scripts/check.sh <shapes.json> <read-committed|repeatable-read|serializable> [txns]
```

The script starts a disposable container (Postgres plus elle-go, pulled automatically on first use), runs the workload inside it, and removes the container. It never touches the user's databases. The user needs only Docker.

Exit codes: `0` = no anomaly found, `1` = anomaly found, `2` = setup or input error (report it; Docker not running and the image failing to pull are the common causes).

## 5. Report

Keep the report short and in plain words. Say what breaks for the user's data, not the theory. Don't use anomaly codes (G2, G-single) unless the user asks; use the plain name. Use exactly this format, with no extra sections:

**Problem found (exit 1):**

```
❌ <plain name> at <level>
<One sentence: which calls run at the same time, and what goes wrong for the data.>
Fix: <one concrete change, e.g. "use serializable for this transaction" or "add FOR UPDATE to the doctors SELECT">
```

Then ask one question: "Re-run with the fix to confirm?"

Plain names: lost or overwritten update / write skew (G2), inconsistent read (G-single), transactions seeing each other's writes (G1c), conflicting writes (G0).

Example:

```
❌ Write skew at read-committed
Alice and Bob can go off call at the same moment; both see two doctors on call, so the shift ends up with nobody.
Fix: run GoOffCall at serializable (and retry on serialization errors).
```

**No problem found (exit 0):**

```
✅ No problem found in <N> runs at <level>
This is not a guarantee, only that the race didn't show up in these runs.
```

Add one line only if more than a third of transactions aborted: "Note: <X>% of runs were rejected by Postgres; your code needs to retry on serialization errors."

**Error (exit 2):** one line saying what failed and how to fix it (e.g. "Docker isn't running; start it and try again.").
