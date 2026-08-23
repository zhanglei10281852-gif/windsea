package service

import (
	"context"
	"database/sql"
	"time"
)

type Health struct {
	Status    string
	Database  string
	CheckedAt time.Time
}

func CheckHealth(ctx context.Context, db *sql.DB, now time.Time) Health {
	health := Health{Status: "ready", Database: "up", CheckedAt: now}
	if err := db.PingContext(ctx); err != nil {
		health.Status = "not_ready"
		health.Database = err.Error()
	}
	return health
}
