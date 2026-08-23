package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"log/slog"
	"time"
)

type HandoffService struct {
	repo *repository.Repository
	log  *slog.Logger
	now  func() time.Time
}

func (s *HandoffService) Offer(ctx context.Context, h domain.Handoff) error {
	if h.ID == "" || h.WorkOrderID == "" || h.ContractorID == "" {
		return domain.ErrValidation
	}
	h.State = "offered"
	h.CreatedAt = s.now()
	return s.repo.CreateHandoff(ctx, h)
}
func (s *HandoffService) Accept(ctx context.Context, id, notes string) error {
	return s.repo.TransitionHandoff(ctx, id, "offered", "accepted", notes, s.now().Format(time.RFC3339Nano))
}
func (s *HandoffService) Complete(ctx context.Context, id, notes string) error {
	return s.repo.TransitionHandoff(ctx, id, "accepted", "completed", notes, s.now().Format(time.RFC3339Nano))
}
func (s *HandoffService) Window(ctx context.Context, w domain.MaintenanceWindow) error {
	if !w.EndsAt.After(w.StartsAt) {
		return fmt.Errorf("%w: window", domain.ErrValidation)
	}
	overlap, err := s.repo.HasWindowOverlap(ctx, w)
	if err != nil {
		return err
	}
	if overlap {
		return domain.ErrConflict
	}
	w.State = "planned"
	return s.repo.CreateWindow(ctx, w)
}
