package pgrun

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib"

	"elle-go/history"
)

const schema = `CREATE TABLE IF NOT EXISTS elle_lists (key text PRIMARY KEY, vals int[] NOT NULL)`

const appendSQL = `INSERT INTO elle_lists (key, vals) VALUES ($1, ARRAY[$2::int])
ON CONFLICT (key) DO UPDATE SET vals = elle_lists.vals || $2::int`

func Setup(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx, `TRUNCATE elle_lists`)
	return err
}

func Exec(ctx context.Context, db *sql.DB, iso sql.IsolationLevel, op history.Op) history.Op {
	done := history.Op{Process: op.Process, Type: history.Fail, Mops: op.Mops}
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: iso})
	if err != nil {
		return done
	}
	defer tx.Rollback()
	mops := make([]history.Mop, len(op.Mops))
	for i, m := range op.Mops {
		mops[i] = m
		if m.Type == history.Append {
			_, err = tx.ExecContext(ctx, appendSQL, m.Key, m.App)
		} else {
			mops[i].Read, err = readList(ctx, tx, m.Key)
		}
		if err != nil {
			return done
		}
	}
	if tx.Commit() != nil {
		return done
	}
	return history.Op{Process: op.Process, Type: history.Ok, Mops: mops}
}

func readList(ctx context.Context, tx *sql.Tx, key string) ([]int, error) {
	var vals pgtype.FlatArray[int]
	err := tx.QueryRowContext(ctx, `SELECT vals FROM elle_lists WHERE key = $1`, key).Scan(&vals)
	if err == sql.ErrNoRows {
		return []int{}, nil
	}
	return []int(vals), err
}
