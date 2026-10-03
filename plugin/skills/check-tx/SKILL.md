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

Show the shapes as a short table (transaction type → ordered ops) plus the isolation level you'll test. Determine the level from the code or config; Postgres defaults to `read-committed`. Ask the user to confirm or correct before running. This review is the safety gate against a wrong extraction.

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

The script starts a disposable Postgres container, runs the workload, and removes the container. It never touches the user's databases.

Exit codes: `0` = no anomaly found, `1` = anomaly found, `2` = setup or input error (report it; elle-go not on PATH and Docker not running are the common causes).

## 5. Explain

- **Anomaly (exit 1).** Name it in plain words (G2 = write skew, G-single = read skew, G1c = circular information flow, G0 = write cycle). Then explain the cycle in terms of the user's code: which two code paths interleave, and what goes wrong for their data. Propose a concrete fix:
  - raise the isolation level to the one the verdict says it requires, or
  - lock the rows that were read (`SELECT ... FOR UPDATE`), or
  - add a constraint that turns the invariant into a write conflict.

  Offer to re-run at the fixed isolation level to confirm it comes back clean.
- **Clean (exit 0).** Say exactly "no anomaly found in N transactions at <level>". This is evidence, not proof of safety. Never call the code "safe".
- Mention how many transactions committed. A high abort rate under `repeatable-read` or `serializable` means the app needs retry logic for serialization failures.
