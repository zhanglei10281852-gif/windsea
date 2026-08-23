package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"log/slog"
	"time"
)

type TelemetryService struct {
	repo *repository.Repository
	log  *slog.Logger
	now  func() time.Time
}

func (s *TelemetryService) Receive(ctx context.Context, batch domain.TelemetryBatch) error {
	if batch.Samples <= 0 {
		return fmt.Errorf("%w: samples", domain.ErrValidation)
	}
	now := s.now()
	batch.ReceivedAt = &now
	batch.State = "received"
	return s.repo.CreateTelemetryBatch(ctx, batch)
}
func (s *TelemetryService) Complete(ctx context.Context, id string, accepted, rejected int) error {
	if accepted < 0 || rejected < 0 || accepted+rejected <= 0 {
		return domain.ErrValidation
	}
	return s.repo.CompleteTelemetryBatch(ctx, id, accepted, rejected, s.now().Format(time.RFC3339Nano))
}
func (s *TelemetryService) Get(ctx context.Context, id string) (domain.TelemetryBatch, error) {
	return s.repo.GetTelemetryBatch(ctx, id)
}
