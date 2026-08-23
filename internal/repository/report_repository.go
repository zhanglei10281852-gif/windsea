package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type FleetStatus struct{ Turbines, ActiveAlerts, OpenWorkOrders, CompletedWorkOrders int }

func (r *Repository) FleetStatus(ctx context.Context, farmID string) (FleetStatus, error) {
	var status FleetStatus
	queries := []struct {
		sql  string
		dest *int
	}{{"SELECT COUNT(*) FROM turbines WHERE farm_id=?", &status.Turbines}, {"SELECT COUNT(*) FROM alerts WHERE farm_id=? AND state<>?", &status.ActiveAlerts}, {"SELECT COUNT(*) FROM work_orders WHERE farm_id=? AND state NOT IN ('completed','cancelled')", &status.OpenWorkOrders}, {"SELECT COUNT(*) FROM work_orders WHERE farm_id=? AND state='completed'", &status.CompletedWorkOrders}}
	for _, query := range queries {
		var err error
		if query.sql == queries[1].sql {
			err = r.DB.QueryRowContext(ctx, query.sql, farmID, "resolved").Scan(query.dest)
		} else {
			err = r.DB.QueryRowContext(ctx, query.sql, farmID).Scan(query.dest)
		}
		if err != nil {
			return status, fmt.Errorf("fleet status: %w", err)
		}
	}
	return status, nil
}
func (r *Repository) CountByState(ctx context.Context, table, farmID string) (map[string]int, error) {
	allowed := map[string]bool{"work_orders": true, "alerts": true, "campaigns": true}
	if !allowed[table] {
		return nil, fmt.Errorf("%w: table", sql.ErrNoRows)
	}
	rows, err := r.DB.QueryContext(ctx, "SELECT state,COUNT(*) FROM "+table+" WHERE farm_id=? GROUP BY state", farmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]int{}
	for rows.Next() {
		var state string
		var count int
		if err := rows.Scan(&state, &count); err != nil {
			return nil, err
		}
		result[state] = count
	}
	return result, rows.Err()
}
