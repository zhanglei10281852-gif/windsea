package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"github.com/zhanglei10281852-gif/windsea/internal/validation"
)

type CampaignService struct {
	repo *repository.Repository
	log  *slog.Logger
	now  func() time.Time
}

func (s *CampaignService) Create(ctx context.Context, campaign domain.Campaign) error {
	if err := validation.Required(campaign.ID, campaign.FarmID, campaign.Name); err != nil {
		return err
	}
	if !campaign.EndAt.After(campaign.StartAt) {
		return fmt.Errorf("%w: campaign window", domain.ErrValidation)
	}
	campaign.State = string(domain.CampaignDraft)
	campaign.Version = 1
	campaign.CreatedAt = s.now()
	campaign.UpdatedAt = campaign.CreatedAt
	return s.repo.CreateCampaign(ctx, campaign)
}
func (s *CampaignService) SubmitReview(ctx context.Context, id string, version int) (domain.Campaign, error) {
	c, err := s.repo.GetCampaign(ctx, id)
	if err != nil {
		return c, err
	}
	if !domain.CampaignState(c.State).CanMove(domain.CampaignReview) {
		return c, domain.ErrInvalidState
	}
	return s.repo.TransitionCampaign(ctx, id, c.State, string(domain.CampaignReview), version, s.now().Format(time.RFC3339Nano))
}
func (s *CampaignService) Publish(ctx context.Context, id string, version int) (domain.Campaign, error) {
	c, err := s.repo.GetCampaign(ctx, id)
	if err != nil {
		return c, err
	}
	if !domain.CampaignState(c.State).CanMove(domain.CampaignPublished) {
		return c, domain.ErrInvalidState
	}
	return s.repo.TransitionCampaign(ctx, id, c.State, string(domain.CampaignPublished), version, s.now().Format(time.RFC3339Nano))
}
func (s *CampaignService) Close(ctx context.Context, id string, version int) (domain.Campaign, error) {
	if err := ctx.Err(); err != nil {
		return domain.Campaign{}, err
	}
	c, err := s.repo.GetCampaign(ctx, id)
	if err != nil {
		return c, err
	}
	if !domain.CampaignState(c.State).CanMove(domain.CampaignClosed) {
		return c, domain.ErrInvalidState
	}
	incomplete, err := s.repo.CountIncompleteInspections(ctx, id)
	if err != nil {
		return c, err
	}
	if incomplete > 0 {
		return c, fmt.Errorf("%w: %d inspections remain", domain.ErrInvalidState, incomplete)
	}
	return s.repo.TransitionCampaign(ctx, id, c.State, string(domain.CampaignClosed), version, s.now().Format(time.RFC3339Nano))
}
func (s *CampaignService) CompleteInspection(ctx context.Context, id, notes string) error {
	return s.repo.CompleteInspection(ctx, id, notes, s.now().Format(time.RFC3339Nano))
}
