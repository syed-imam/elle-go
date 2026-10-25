# What `/check-tx` finds

Each case below is real code from [`sampleapp`](sampleapp), checked with the
`check-tx` plugin skill. Claude reads the transaction, extracts its read/write
shapes, you confirm them, and elle-go runs them concurrently against a
throwaway Postgres. The verdict lines are copied verbatim from the runs.

| Code | Bug class | Read committed | Fix | After fix |
|---|---|---|---|---|
| `GoOffCall` | write skew | **G-single** | `FOR UPDATE`, or `SERIALIZABLE` | clean |
| `Transfer` | lost update | **G-single** | `REPEATABLE READ`, or an in-place `UPDATE` | clean |
| `Withdraw` | none | clean | — | — |
| `TotalDeposits` | predicate read | refused | — | — |

---

## 1. Write skew: the last doctor goes off call

[`oncall.go`](sampleapp/oncall.go). A doctor may leave a shift only if both
doctors are still on call.

```go
for _, name := range []string{primary, backup} {
    tx.QueryRowContext(ctx, `SELECT on_call FROM doctors WHERE name = $1`, name).Scan(&on)
    if on { onCall++ }
}
if onCall < 2 { return ErrLastDoctor }
tx.ExecContext(ctx, `UPDATE doctors SET on_call = false WHERE name = $1`, doctor)
```

**Shapes Claude extracts** (you confirm before anything runs):

| Transaction | Ops |
|---|---|
| alice goes off call | r `shift:1` → r `doctor:alice` → r `doctor:bob` → write `doctor:alice` |
| bob goes off call | r `shift:1` → r `doctor:alice` → r `doctor:bob` → write `doctor:bob` |

**Verdict at each level:**

```
read-committed: 300/300 committed
G-single: transactions 220 → 218 → 220 (requires snapshot isolation or stronger)

repeatable-read: 82/300 committed
G2: transactions 297 → 299 → 297 (requires serializable or stronger)

serializable: 50/300 committed
no anomaly found
```

**What it means.** Alice and Bob both see two doctors on call, both go off
call, and the shift ends up with nobody. `REPEATABLE READ` doesn't fix this,
because each transaction writes a different row, so neither one blocks the
other. That's classic write skew (G2).

**Fix.** Lock both rows you read: `SELECT on_call … FOR UPDATE`. Re-run at the
default level:

```
read-committed: 300/300 committed
no anomaly found
```

`SERIALIZABLE` also comes back clean, but only 50 of 300 transactions commit,
so the app would need retry logic for serialization failures.

---

## 2. Lost update: concurrent transfers

[`bank.go`](sampleapp/bank.go). `Transfer` reads both balances, then writes
back values computed in Go.

```go
tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = $1`, from).Scan(&fromBal)
tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = $1`, to).Scan(&toBal)
tx.ExecContext(ctx, `UPDATE accounts SET balance = $1 WHERE id = $2`, fromBal-amount, from)
tx.ExecContext(ctx, `UPDATE accounts SET balance = $1 WHERE id = $2`, toBal+amount, to)
```

**Shapes:** both directions, so the transfers collide.

| Transaction | Ops |
|---|---|
| 1 → 2 | r `account:1` → r `account:2` → write `account:1` → write `account:2` |
| 2 → 1 | r `account:2` → r `account:1` → write `account:2` → write `account:1` |

**Verdict:**

```
read-committed: 82/300 committed
G-single: transactions 144 → 142 → 143 → 144 (requires snapshot isolation or stronger)

repeatable-read: 42/300 committed
no anomaly found
```

**What it means.** Two transfers read the same balance, and the second write
overwrites the first, so money appears or disappears. Only 82 commit even at
read committed: the opposite-order row updates deadlock.

**Fix.** Either run at `REPEATABLE READ` with retries, or let the database do
the arithmetic (`SET balance = balance - $1`) and update rows in id order to
avoid the deadlocks.

---

## 3. Clean: `Withdraw` already locks

[`bank.go`](sampleapp/bank.go). `Withdraw` reads with `FOR UPDATE` before
writing.

| Transaction | Ops |
|---|---|
| withdraw | r `account:1` (FOR UPDATE) → write `account:1` |

```
read-committed: 300/300 committed
no anomaly found
```

The plugin reports this as "no anomaly found in 300 transactions at read
committed". That's evidence, not proof, so it never calls the code "safe".

---

## 4. Refused: `TotalDeposits`

```sql
SELECT COALESCE(SUM(balance), 0) FROM accounts WHERE branch = $1
```

This is a predicate read: which rows it reads depends on a `WHERE` clause, not
a fixed id. elle-go models individual rows only, so it can't see phantoms. The
plugin says so and stops instead of producing a misleading result.

---

## Reproduce

```sh
/plugin marketplace add syed-imam/elle-go
/plugin install elle-go@elle-go
/elle-go:check-tx GoOffCall
```

Requires Docker. Each run takes a few seconds and never touches your own
databases.
