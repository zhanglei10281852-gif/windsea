package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"log/slog"
	"time"
)

type AlertService struct {
	repo *repository.Repository
	log  *slog.Logger
	now  func() time.Time
}

func (s *AlertService) Ingest(ctx context.Context, a domain.Alert) error {
	if a.OccurredAt == nil {
		now := s.now()
		a.OccurredAt = &now
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = s.now()
	}
	a.State = string(domain.AlertOpen)
	if a.DedupKey == "" {
		a.DedupKey = a.TurbineID + ":" + a.Code + ":" + a.OccurredAt.Format(time.RFC3339Nano)
	}
	return s.repo.CreateAlert(ctx, a)
}
func (s *AlertService) Acknowledge(ctx context.Context, id string) error {
	alert, err := s.repo.GetAlert(ctx, id)
	if err != nil {
		return err
	}
	if !domain.AlertState(alert.State).CanMove(domain.AlertAcknowledged) {
		return fmt.Errorf("%w: alert state", domain.ErrInvalidState)
	}
	return s.repo.TransitionAlert(ctx, id, alert.State, string(domain.AlertAcknowledged), s.now().Format(time.RFC3339Nano))
}
func (s *AlertService) Resolve(ctx context.Context, id string) error {
	alert, err := s.repo.GetAlert(ctx, id)
	if err != nil {
		return err
	}
	if !domain.AlertState(alert.State).CanMove(domain.AlertResolved) {
		return fmt.Errorf("%w: acknowledge first", domain.ErrInvalidState)
	}
	return s.repo.TransitionAlert(ctx, id, alert.State, string(domain.AlertResolved), s.now().Format(time.RFC3339Nano))
}
func (s *AlertService) Open(ctx context.Context, farmID string) ([]domain.Alert, error) {
	return s.repo.ListOpenAlerts(ctx, farmID)
}
