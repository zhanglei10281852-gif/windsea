package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct{ *sql.DB }

func Open(ctx context.Context, dsn string) (*DB, error) {
	connection, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	connection.SetMaxOpenConns(8)
	connection.SetMaxIdleConns(8)
	if err := connection.PingContext(ctx); err != nil {
		_ = connection.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	database := &DB{DB: connection}
	if err := Migrate(ctx, database); err != nil {
		_ = connection.Close()
		return nil, err
	}
	return database, nil
}

func (d *DB) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
