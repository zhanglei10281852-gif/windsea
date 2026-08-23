package db

import (
	"context"
	"database/sql"
	"fmt"
)

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL);`,
	`PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS farms (id TEXT PRIMARY KEY, name TEXT NOT NULL, timezone TEXT NOT NULL, turbine_count INTEGER NOT NULL DEFAULT 0, created_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY, farm_id TEXT NOT NULL REFERENCES farms(id), email TEXT NOT NULL UNIQUE, name TEXT NOT NULL, role TEXT NOT NULL, active INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS sessions (id TEXT PRIMARY KEY, user_id TEXT NOT NULL REFERENCES users(id), token_hash TEXT NOT NULL UNIQUE, expires_at TEXT NOT NULL, revoked_at TEXT);
CREATE TABLE IF NOT EXISTS turbines (id TEXT PRIMARY KEY, farm_id TEXT NOT NULL REFERENCES farms(id), name TEXT NOT NULL, model TEXT NOT NULL, rated_kw INTEGER NOT NULL, status TEXT NOT NULL, version INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, UNIQUE(farm_id, name));
CREATE TABLE IF NOT EXISTS campaigns (id TEXT PRIMARY KEY, farm_id TEXT NOT NULL REFERENCES farms(id), name TEXT NOT NULL, state TEXT NOT NULL, start_at TEXT NOT NULL, end_at TEXT NOT NULL, version INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL, UNIQUE(farm_id, name));
CREATE TABLE IF NOT EXISTS inspections (id TEXT PRIMARY KEY, campaign_id TEXT NOT NULL REFERENCES campaigns(id), turbine_id TEXT NOT NULL REFERENCES turbines(id), inspector_id TEXT NOT NULL REFERENCES users(id), status TEXT NOT NULL, notes TEXT NOT NULL DEFAULT '', completed_at TEXT, created_at TEXT NOT NULL, UNIQUE(campaign_id, turbine_id));
CREATE TABLE IF NOT EXISTS work_orders (id TEXT PRIMARY KEY, farm_id TEXT NOT NULL REFERENCES farms(id), campaign_id TEXT REFERENCES campaigns(id), turbine_id TEXT REFERENCES turbines(id), assignee_id TEXT REFERENCES users(id), title TEXT NOT NULL, state TEXT NOT NULL, priority INTEGER NOT NULL, version INTEGER NOT NULL DEFAULT 1, due_at TEXT, created_at TEXT NOT NULL, updated_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS parts (id TEXT PRIMARY KEY, farm_id TEXT NOT NULL REFERENCES farms(id), sku TEXT NOT NULL, description TEXT NOT NULL, on_hand INTEGER NOT NULL, reserved INTEGER NOT NULL DEFAULT 0, version INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, UNIQUE(farm_id, sku));
CREATE TABLE IF NOT EXISTS reservations (id TEXT PRIMARY KEY, part_id TEXT NOT NULL REFERENCES parts(id), work_order_id TEXT NOT NULL REFERENCES work_orders(id), requested_by TEXT NOT NULL REFERENCES users(id), state TEXT NOT NULL, quantity INTEGER NOT NULL, version INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL, UNIQUE(part_id, work_order_id));
CREATE TABLE IF NOT EXISTS alerts (id TEXT PRIMARY KEY, farm_id TEXT NOT NULL REFERENCES farms(id), turbine_id TEXT NOT NULL REFERENCES turbines(id), code TEXT NOT NULL, severity TEXT NOT NULL, state TEXT NOT NULL, message TEXT NOT NULL, dedup_key TEXT NOT NULL UNIQUE, occurred_at TEXT NOT NULL, acknowledged_at TEXT, resolved_at TEXT, created_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS telemetry_batches (id TEXT PRIMARY KEY, farm_id TEXT NOT NULL REFERENCES farms(id), source TEXT NOT NULL, state TEXT NOT NULL, samples INTEGER NOT NULL, accepted INTEGER NOT NULL DEFAULT 0, rejected INTEGER NOT NULL DEFAULT 0, received_at TEXT NOT NULL, completed_at TEXT);
CREATE TABLE IF NOT EXISTS handoffs (id TEXT PRIMARY KEY, work_order_id TEXT NOT NULL REFERENCES work_orders(id), contractor_id TEXT NOT NULL REFERENCES users(id), state TEXT NOT NULL, notes TEXT NOT NULL DEFAULT '', accepted_at TEXT, completed_at TEXT, created_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS maintenance_windows (id TEXT PRIMARY KEY, farm_id TEXT NOT NULL REFERENCES farms(id), turbine_id TEXT NOT NULL REFERENCES turbines(id), name TEXT NOT NULL, starts_at TEXT NOT NULL, ends_at TEXT NOT NULL, state TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS audit_events (id TEXT PRIMARY KEY, farm_id TEXT NOT NULL REFERENCES farms(id), actor_id TEXT, object_type TEXT NOT NULL, object_id TEXT NOT NULL, action TEXT NOT NULL, result TEXT NOT NULL, request_id TEXT NOT NULL, metadata TEXT NOT NULL, created_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS idempotency_keys (key TEXT PRIMARY KEY, scope TEXT NOT NULL, response TEXT NOT NULL, created_at TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_work_orders_farm_state ON work_orders(farm_id, state, priority);
CREATE INDEX IF NOT EXISTS idx_alerts_farm_state ON alerts(farm_id, state, occurred_at);
CREATE INDEX IF NOT EXISTS idx_audit_object ON audit_events(object_type, object_id, created_at);`,
}

func Migrate(ctx context.Context, d *DB) error {
	if _, err := d.ExecContext(ctx, migrations[0]); err != nil {
		return fmt.Errorf("migration table: %w", err)
	}
	for version, script := range migrations[1:] {
		var applied int
		if err := d.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version+1).Scan(&applied); err != nil {
			return fmt.Errorf("migration check %d: %w", version+1, err)
		}
		if applied > 0 {
			continue
		}
		if err := d.WithTx(ctx, func(tx *sql.Tx) error {
			if _, err := tx.ExecContext(ctx, script); err != nil {
				return fmt.Errorf("migration %d: %w", version+1, err)
			}
			if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations(version, applied_at) VALUES (?, CURRENT_TIMESTAMP)", version+1); err != nil {
				return fmt.Errorf("record migration %d: %w", version+1, err)
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
