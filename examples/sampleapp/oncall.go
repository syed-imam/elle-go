package sampleapp

import (
	"context"
	"database/sql"
	"errors"
)

var ErrLastDoctor = errors.New("at least one doctor must stay on call")

func GoOffCall(ctx context.Context, db *sql.DB, shiftID int, doctor string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var primary, backup string
	if err := tx.QueryRowContext(ctx, `SELECT primary_doctor, backup_doctor FROM shifts WHERE id = $1`, shiftID).Scan(&primary, &backup); err != nil {
		return err
	}
	onCall := 0
	for _, name := range []string{primary, backup} {
		var on bool
		if err := tx.QueryRowContext(ctx, `SELECT on_call FROM doctors WHERE name = $1`, name).Scan(&on); err != nil {
			return err
		}
		if on {
			onCall++
		}
	}
	if onCall < 2 {
		return ErrLastDoctor
	}
	if _, err := tx.ExecContext(ctx, `UPDATE doctors SET on_call = false WHERE name = $1`, doctor); err != nil {
		return err
	}
	return tx.Commit()
}
