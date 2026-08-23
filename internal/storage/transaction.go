package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

type TxStore struct{ DB *sql.DB }

func (s TxStore) Run(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: commit", err)
	}
	return nil
}
func (s TxStore) EnsureFarm(ctx context.Context, farmID string) error {
	if farmID == "" {
		return domain.ErrValidation
	}
	var exists int
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM farms WHERE id=?", farmID).Scan(&exists); err != nil {
		return err
	}
	if exists == 0 {
		return domain.ErrNotFound
	}
	return nil
}
