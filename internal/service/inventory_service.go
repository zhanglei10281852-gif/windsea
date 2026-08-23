package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"github.com/zhanglei10281852-gif/windsea/internal/validation"
	"log/slog"
	"time"
)

type InventoryService struct {
	repo *repository.Repository
	log  *slog.Logger
	now  func() time.Time
}

func (s *InventoryService) CreatePart(ctx context.Context, part domain.Part) error {
	if err := validation.Required(part.ID, part.FarmID, part.SKU); err != nil {
		return err
	}
	if err := validation.Positive(part.OnHand); err != nil {
		return err
	}
	part.Version = 1
	part.CreatedAt = s.now()
	return s.repo.CreatePart(ctx, part)
}
func (s *InventoryService) Hold(ctx context.Context, id, partID, workOrderID, userID string, quantity int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validation.Required(id, partID, workOrderID, userID); err != nil {
		return err
	}
	if err := validation.Positive(quantity); err != nil {
		return err
	}
	at := s.now()
	return s.repo.HoldReservation(ctx, domain.Reservation{ID: id, PartID: partID, WorkOrderID: workOrderID, RequestedBy: userID, Quantity: quantity, Version: 1, State: string(domain.ReservationHeld), CreatedAt: at, UpdatedAt: at})
}
func (s *InventoryService) Consume(ctx context.Context, id string) error {
	return s.repo.MoveReservation(ctx, id, string(domain.ReservationHeld), string(domain.ReservationConsumed), s.now().Format(time.RFC3339Nano))
}
func (s *InventoryService) Release(ctx context.Context, id string) error {
	return s.repo.MoveReservation(ctx, id, string(domain.ReservationHeld), string(domain.ReservationReleased), s.now().Format(time.RFC3339Nano))
}
func (s *InventoryService) EnsureAvailable(ctx context.Context, partID string, quantity int) error {
	part, err := s.repo.GetPart(ctx, partID)
	if err != nil {
		return err
	}
	if part.OnHand-part.Reserved < quantity {
		return fmt.Errorf("%w: part %s", domain.ErrCapacity, partID)
	}
	return nil
}
