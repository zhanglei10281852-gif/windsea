package db

import (
	"context"
	"database/sql"
	"fmt"
)

func Checkpoint(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, "PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
		return fmt.Errorf("checkpoint: %w", err)
	}
	return nil
}
func ForeignKeysEnabled(ctx context.Context, db *sql.DB) (bool, error) {
	var value int
	if err := db.QueryRowContext(ctx, "PRAGMA foreign_keys").Scan(&value); err != nil {
		return false, err
	}
	return value == 1, nil
}
func TableCount(ctx context.Context, db *sql.DB) (int, error) {
	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
