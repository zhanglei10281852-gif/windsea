package repository

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

func (r *Repository) CreateWorkOrder(ctx context.Context, order domain.WorkOrder) error {
	var campaignID, turbineID any
	if order.CampaignID != "" {
		campaignID = order.CampaignID
	}
	if order.TurbineID != "" {
		turbineID = order.TurbineID
	}
	_, err := r.DB.ExecContext(ctx, "INSERT INTO work_orders(id,farm_id,campaign_id,turbine_id,assignee_id,title,state,priority,version,due_at,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)", order.ID, order.FarmID, campaignID, turbineID, nil, order.Title, order.State, order.Priority, order.Version, nil, text(order.CreatedAt), text(order.UpdatedAt))
	if err != nil {
		return fmt.Errorf("create work order: %w", err)
	}
	return nil
}
func (r *Repository) GetWorkOrder(ctx context.Context, id string) (domain.WorkOrder, error) {
	var o domain.WorkOrder
	var campaign, turbine, assignee, due, created, updated interface{}
	err := r.DB.QueryRowContext(ctx, "SELECT id,farm_id,campaign_id,turbine_id,assignee_id,title,state,priority,version,due_at,created_at,updated_at FROM work_orders WHERE id=?", id).Scan(&o.ID, &o.FarmID, &campaign, &turbine, &assignee, &o.Title, &o.State, &o.Priority, &o.Version, &due, &created, &updated)
	if err != nil {
		if isNoRows(err) {
			return o, domain.ErrNotFound
		}
		return o, fmt.Errorf("get work order: %w", err)
	}
	if value, ok := campaign.(string); ok {
		o.CampaignID = value
	}
	if value, ok := turbine.(string); ok {
		o.TurbineID = value
	}
	if value, ok := assignee.(string); ok {
		o.AssigneeID = value
	}
	if value, ok := due.(string); ok {
		parsed, _ := parse(value)
		o.DueAt = &parsed
	}
	if value, ok := created.(string); ok {
		o.CreatedAt, _ = parse(value)
	}
	if value, ok := updated.(string); ok {
		o.UpdatedAt, _ = parse(value)
	}
	return o, nil
}
func (r *Repository) AssignWorkOrder(ctx context.Context, id, assignee string, version int, at string) error {
	result, err := r.DB.ExecContext(ctx, "UPDATE work_orders SET assignee_id=?,state='assigned',version=version+1,updated_at=? WHERE id=?", assignee, at, id)
	if err != nil {
		return fmt.Errorf("assign work order: %w", err)
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r *Repository) TransitionWorkOrder(ctx context.Context, id, from, to string, version int, at string) error {
	result, err := r.DB.ExecContext(ctx, "UPDATE work_orders SET state=?,version=version+1,updated_at=? WHERE id=? AND state=? AND version=?", to, at, id, from, version)
	if err != nil {
		return fmt.Errorf("transition work order: %w", err)
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r *Repository) ListWorkOrders(ctx context.Context, farmID, state string, limit, offset int) (domain.Page[domain.WorkOrder], error) {
	query := "SELECT id,farm_id,COALESCE(campaign_id,''),COALESCE(turbine_id,''),COALESCE(assignee_id,''),title,state,priority,version,created_at,updated_at FROM work_orders WHERE farm_id=?"
	args := []any{farmID}
	if state != "" {
		query += " AND state=?"
		args = append(args, state)
	}
	query += " ORDER BY priority DESC,created_at LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return domain.Page[domain.WorkOrder]{}, fmt.Errorf("list work orders: %w", err)
	}
	defer rows.Close()
	page := domain.Page[domain.WorkOrder]{Items: []domain.WorkOrder{}, Limit: limit, Offset: offset}
	for rows.Next() {
		var o domain.WorkOrder
		var created, updated string
		if err := rows.Scan(&o.ID, &o.FarmID, &o.CampaignID, &o.TurbineID, &o.AssigneeID, &o.Title, &o.State, &o.Priority, &o.Version, &created, &updated); err != nil {
			return page, err
		}
		o.CreatedAt, _ = parse(created)
		o.UpdatedAt, _ = parse(updated)
		page.Items = append(page.Items, o)
	}
	if err := rows.Err(); err != nil {
		return page, err
	}
	countQuery := "SELECT COUNT(*) FROM work_orders WHERE farm_id=?"
	countArgs := []any{farmID}
	if state != "" {
		countQuery += " AND state=?"
		countArgs = append(countArgs, state)
	}
	if err := r.DB.QueryRowContext(ctx, countQuery, countArgs...).Scan(&page.Total); err != nil {
		return page, err
	}
	return page, nil
}
