package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"time"
)

func (r *Repository) CreateAlert(ctx context.Context, a domain.Alert) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO alerts(id,farm_id,turbine_id,code,severity,state,message,dedup_key,occurred_at,created_at) VALUES(?,?,?,?,?,?,?,?,?,?)", a.ID, a.FarmID, a.TurbineID, a.Code, a.Severity, a.State, a.Message, a.DedupKey, text(*a.OccurredAt), text(a.CreatedAt))
	if err != nil {
		return fmt.Errorf("create alert: %w", err)
	}
	return nil
}
func (r *Repository) GetAlert(ctx context.Context, id string) (domain.Alert, error) {
	var a domain.Alert
	var occurred, created string
	var ack, resolved sql.NullString
	err := r.DB.QueryRowContext(ctx, "SELECT id,farm_id,turbine_id,code,severity,state,message,dedup_key,occurred_at,acknowledged_at,resolved_at,created_at FROM alerts WHERE id=?", id).Scan(&a.ID, &a.FarmID, &a.TurbineID, &a.Code, &a.Severity, &a.State, &a.Message, &a.DedupKey, &occurred, &ack, &resolved, &created)
	if err != nil {
		if isNoRows(err) {
			return a, domain.ErrNotFound
		}
		return a, err
	}
	t, _ := parse(occurred)
	a.OccurredAt = &t
	a.AcknowledgedAt, _ = optional(ack)
	a.ResolvedAt, _ = optional(resolved)
	a.CreatedAt, _ = parse(created)
	return a, nil
}
func optionalString(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := parse(value)
	return &parsed, err
}
func (r *Repository) TransitionAlert(ctx context.Context, id, from, to, at string) error {
	column := "acknowledged_at"
	if to == string(domain.AlertResolved) {
		column = "resolved_at"
	}
	query := fmt.Sprintf("UPDATE alerts SET state=?,%s=? WHERE id=? AND state=?", column)
	result, err := r.DB.ExecContext(ctx, query, to, at, id, from)
	if err != nil {
		return fmt.Errorf("transition alert: %w", err)
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r *Repository) ListOpenAlerts(ctx context.Context, farmID string) ([]domain.Alert, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id,farm_id,turbine_id,code,severity,state,message,dedup_key,occurred_at,created_at FROM alerts WHERE farm_id=? AND state<>? ORDER BY occurred_at", farmID, string(domain.AlertResolved))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.Alert{}
	for rows.Next() {
		var a domain.Alert
		var occurred, created string
		if err := rows.Scan(&a.ID, &a.FarmID, &a.TurbineID, &a.Code, &a.Severity, &a.State, &a.Message, &a.DedupKey, &occurred, &created); err != nil {
			return nil, err
		}
		t, _ := parse(occurred)
		a.OccurredAt = &t
		a.CreatedAt, _ = parse(created)
		items = append(items, a)
	}
	return items, rows.Err()
}
