package tx

import (
	"context"
	"database/sql"
	"fmt"
)

type Runner interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func Do(ctx context.Context, db interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}, fn func(*sql.Tx) error) error {
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	if err := fn(transaction); err != nil {
		_ = transaction.Rollback()
		return err
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}
