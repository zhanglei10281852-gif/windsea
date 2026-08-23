package repository

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

func (r *Repository) CreateTelemetryBatch(ctx context.Context, b domain.TelemetryBatch) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO telemetry_batches(id,farm_id,source,state,samples,accepted,rejected,received_at) VALUES(?,?,?,?,?,?,?,?)", b.ID, b.FarmID, b.Source, b.State, b.Samples, b.Accepted, b.Rejected, text(*b.ReceivedAt))
	if err != nil {
		return fmt.Errorf("create telemetry batch: %w", err)
	}
	return nil
}
func (r *Repository) CompleteTelemetryBatch(ctx context.Context, id string, accepted, rejected int, at string) error {
	result, err := r.DB.ExecContext(ctx, "UPDATE telemetry_batches SET state='completed',accepted=?,rejected=?,completed_at=? WHERE id=? AND state='received'", accepted, rejected, at, id)
	if err != nil {
		return fmt.Errorf("complete telemetry batch: %w", err)
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r *Repository) GetTelemetryBatch(ctx context.Context, id string) (domain.TelemetryBatch, error) {
	var b domain.TelemetryBatch
	var received, completed string
	err := r.DB.QueryRowContext(ctx, "SELECT id,farm_id,source,state,samples,accepted,rejected,received_at,completed_at FROM telemetry_batches WHERE id=?", id).Scan(&b.ID, &b.FarmID, &b.Source, &b.State, &b.Samples, &b.Accepted, &b.Rejected, &received, &completed)
	if err != nil {
		if isNoRows(err) {
			return b, domain.ErrNotFound
		}
		return b, err
	}
	t, _ := parse(received)
	b.ReceivedAt = &t
	b.CompletedAt, _ = optionalString(completed)
	return b, nil
}
