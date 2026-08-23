package repository

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

func (r *Repository) CreateHandoff(ctx context.Context, h domain.Handoff) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO handoffs(id,work_order_id,contractor_id,state,notes,created_at) VALUES(?,?,?,?,?,?)", h.ID, h.WorkOrderID, h.ContractorID, h.State, h.Notes, text(h.CreatedAt))
	if err != nil {
		return fmt.Errorf("create handoff: %w", err)
	}
	return nil
}
func (r *Repository) TransitionHandoff(ctx context.Context, id, from, to, notes, at string) error {
	column := "accepted_at"
	if to == "completed" {
		column = "completed_at"
	}
	query := fmt.Sprintf("UPDATE handoffs SET state=?,notes=?,%s=? WHERE id=? AND state=?", column)
	result, err := r.DB.ExecContext(ctx, query, to, notes, at, id, from)
	if err != nil {
		return fmt.Errorf("transition handoff: %w", err)
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r *Repository) CreateWindow(ctx context.Context, w domain.MaintenanceWindow) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO maintenance_windows(id,farm_id,turbine_id,name,starts_at,ends_at,state) VALUES(?,?,?,?,?,?,?)", w.ID, w.FarmID, w.TurbineID, w.Name, text(w.StartsAt), text(w.EndsAt), w.State)
	if err != nil {
		return fmt.Errorf("create maintenance window: %w", err)
	}
	return nil
}
func (r *Repository) HasWindowOverlap(ctx context.Context, w domain.MaintenanceWindow) (bool, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM maintenance_windows WHERE turbine_id=? AND state<>'cancelled' AND starts_at < ? AND ends_at > ?", w.TurbineID, text(w.EndsAt), text(w.StartsAt)).Scan(&count)
	return count > 0, err
}
