package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"time"
)

type Journal struct{ DB *sql.DB }

func (j Journal) Append(ctx context.Context, id, kind, payload string, at time.Time) error {
	if id == "" || kind == "" {
		return domain.ErrValidation
	}
	if _, err := j.DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS worker_journal(id TEXT PRIMARY KEY,kind TEXT NOT NULL,payload TEXT NOT NULL,created_at TEXT NOT NULL)"); err != nil {
		return err
	}
	_, err := j.DB.ExecContext(ctx, "INSERT INTO worker_journal(id,kind,payload,created_at) VALUES(?,?,?,?)", id, kind, payload, at.UTC().Format(time.RFC3339Nano))
	return err
}
func (j Journal) Pending(ctx context.Context, kind string) ([]string, error) {
	rows, err := j.DB.QueryContext(ctx, "SELECT payload FROM worker_journal WHERE kind=? ORDER BY created_at", kind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return nil, fmt.Errorf("journal row: %w", err)
		}
		result = append(result, payload)
	}
	return result, rows.Err()
}
