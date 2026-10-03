package sampleapp

import (
	"context"
	"database/sql"
	"errors"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

func Transfer(ctx context.Context, db *sql.DB, from, to int, amount int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromBal, toBal int64
	if err := tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = $1`, from).Scan(&fromBal); err != nil {
		return err
	}
	if fromBal < amount {
		return ErrInsufficientFunds
	}
	if err := tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = $1`, to).Scan(&toBal); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE accounts SET balance = $1 WHERE id = $2`, fromBal-amount, from); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE accounts SET balance = $1 WHERE id = $2`, toBal+amount, to); err != nil {
		return err
	}
	return tx.Commit()
}

func Withdraw(ctx context.Context, db *sql.DB, account int, amount int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var bal int64
	if err := tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, account).Scan(&bal); err != nil {
		return err
	}
	if bal < amount {
		return ErrInsufficientFunds
	}
	if _, err := tx.ExecContext(ctx, `UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, account); err != nil {
		return err
	}
	return tx.Commit()
}

func TotalDeposits(ctx context.Context, db *sql.DB, branch string) (int64, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var total int64
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(SUM(balance), 0) FROM accounts WHERE branch = $1`, branch).Scan(&total); err != nil {
		return 0, err
	}
	return total, tx.Commit()
}
