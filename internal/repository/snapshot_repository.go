package repository

import (
	"context"
	"encoding/json"
	"fmt"
)

type Snapshot struct {
	FarmID     string `json:"farm_id"`
	WorkOrders int    `json:"work_orders"`
	Alerts     int    `json:"alerts"`
	CapturedAt string `json:"captured_at"`
}

func (r *Repository) SaveSnapshot(ctx context.Context, snapshot Snapshot) error {
	raw, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	if _, err := r.DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS farm_snapshots(farm_id TEXT PRIMARY KEY,payload TEXT NOT NULL)"); err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, "INSERT INTO farm_snapshots(farm_id,payload) VALUES(?,?) ON CONFLICT(farm_id) DO UPDATE SET payload=excluded.payload", snapshot.FarmID, string(raw))
	return err
}
func (r *Repository) LoadSnapshot(ctx context.Context, farmID string) (Snapshot, error) {
	var raw string
	err := r.DB.QueryRowContext(ctx, "SELECT payload FROM farm_snapshots WHERE farm_id=?", farmID).Scan(&raw)
	if err != nil {
		return Snapshot{}, fmt.Errorf("load snapshot: %w", err)
	}
	var snapshot Snapshot
	if err := json.Unmarshal([]byte(raw), &snapshot); err != nil {
		return Snapshot{}, err
	}
	return snapshot, nil
}
