package repository

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

func (r *Repository) CreateCampaign(ctx context.Context, campaign domain.Campaign) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO campaigns(id,farm_id,name,state,start_at,end_at,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)", campaign.ID, campaign.FarmID, campaign.Name, campaign.State, text(campaign.StartAt), text(campaign.EndAt), campaign.Version, text(campaign.CreatedAt), text(campaign.UpdatedAt))
	if err != nil {
		return fmt.Errorf("create campaign: %w", err)
	}
	return nil
}
func (r *Repository) GetCampaign(ctx context.Context, id string) (domain.Campaign, error) {
	var c domain.Campaign
	var start, end, created, updated string
	err := r.DB.QueryRowContext(ctx, "SELECT id,farm_id,name,state,start_at,end_at,version,created_at,updated_at FROM campaigns WHERE id=?", id).Scan(&c.ID, &c.FarmID, &c.Name, &c.State, &start, &end, &c.Version, &created, &updated)
	if err != nil {
		if isNoRows(err) {
			return c, domain.ErrNotFound
		}
		return c, fmt.Errorf("get campaign: %w", err)
	}
	c.StartAt, _ = parse(start)
	c.EndAt, _ = parse(end)
	c.CreatedAt, _ = parse(created)
	c.UpdatedAt, _ = parse(updated)
	return c, nil
}
func (r *Repository) TransitionCampaign(ctx context.Context, id string, from, to string, version int, at string) (domain.Campaign, error) {
	result, err := r.DB.ExecContext(ctx, "UPDATE campaigns SET state=?,version=version+1,updated_at=? WHERE id=? AND state=? AND version=?", to, at, id, from, version)
	if err != nil {
		return domain.Campaign{}, fmt.Errorf("transition campaign: %w", err)
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		return domain.Campaign{}, domain.ErrConflict
	}
	return r.GetCampaign(ctx, id)
}
func (r *Repository) AddInspection(ctx context.Context, inspection domain.Inspection) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO inspections(id,campaign_id,turbine_id,inspector_id,status,notes,completed_at,created_at) VALUES(?,?,?,?,?,?,?,?)", inspection.ID, inspection.CampaignID, inspection.TurbineID, inspection.InspectorID, inspection.Status, inspection.Notes, nil, text(inspection.CreatedAt))
	if err != nil {
		return fmt.Errorf("add inspection: %w", err)
	}
	return nil
}
func (r *Repository) CountIncompleteInspections(ctx context.Context, campaignID string) (int, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM inspections WHERE campaign_id=? AND status<>?", campaignID, "completed").Scan(&count)
	return count, err
}
func (r *Repository) CompleteInspection(ctx context.Context, id, notes, at string) error {
	result, err := r.DB.ExecContext(ctx, "UPDATE inspections SET status='completed',notes=?,completed_at=? WHERE id=? AND status<>'completed'", notes, at, id)
	if err != nil {
		return fmt.Errorf("complete inspection: %w", err)
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		return domain.ErrConflict
	}
	return nil
}
