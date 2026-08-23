package storage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type LockTable struct {
	DB *sql.DB
	mu sync.Mutex
}

func (l *LockTable) Acquire(ctx context.Context, key string, ttl time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, err := l.DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS windsea_locks(key TEXT PRIMARY KEY,expires_at TEXT NOT NULL)"); err != nil {
		return err
	}
	var expires string
	err := l.DB.QueryRowContext(ctx, "SELECT expires_at FROM windsea_locks WHERE key=?", key).Scan(&expires)
	if err == nil {
		when, _ := time.Parse(time.RFC3339Nano, expires)
		if when.After(time.Now()) {
			return fmt.Errorf("lock held")
		}
	}
	_, err = l.DB.ExecContext(ctx, "INSERT INTO windsea_locks(key,expires_at) VALUES(?,?) ON CONFLICT(key) DO UPDATE SET expires_at=excluded.expires_at", key, time.Now().Add(ttl).UTC().Format(time.RFC3339Nano))
	return err
}
func (l *LockTable) Release(ctx context.Context, key string) error {
	_, err := l.DB.ExecContext(ctx, "DELETE FROM windsea_locks WHERE key=?", key)
	return err
}
