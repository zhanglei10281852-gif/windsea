package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"github.com/zhanglei10281852-gif/windsea/internal/validation"
	"log/slog"
	"strings"
	"time"
)

type WorkOrderService struct {
	repo *repository.Repository
	log  *slog.Logger
	now  func() time.Time
}

func (s *WorkOrderService) Create(ctx context.Context, order domain.WorkOrder) error {
	if err := validation.Required(order.ID, order.FarmID, order.Title); err != nil {
		return err
	}
	if err := validation.Range(order.Priority, 1, 5); err != nil {
		return err
	}
	order.State = string(domain.WorkQueued)
	order.Version = 1
	order.CreatedAt = s.now()
	order.UpdatedAt = order.CreatedAt
	return s.repo.CreateWorkOrder(ctx, order)
}
func (s *WorkOrderService) Assign(ctx context.Context, id, assignee string, version int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	assignee = strings.TrimSpace(assignee)
	if err := validation.Required(id, assignee); err != nil {
		return err
	}
	return s.repo.AssignWorkOrder(ctx, id, assignee, version, s.now().Format(time.RFC3339Nano))
}
func (s *WorkOrderService) Start(ctx context.Context, id string, version int) error {
	order, err := s.repo.GetWorkOrder(ctx, id)
	if err != nil {
		return err
	}
	if !domain.WorkOrderState(order.State).CanMove(domain.WorkInProgress) {
		return domain.ErrInvalidState
	}
	return s.repo.TransitionWorkOrder(ctx, id, order.State, string(domain.WorkInProgress), version, s.now().Format(time.RFC3339Nano))
}
func (s *WorkOrderService) Block(ctx context.Context, id string, version int) error {
	order, err := s.repo.GetWorkOrder(ctx, id)
	if err != nil {
		return err
	}
	if !domain.WorkOrderState(order.State).CanMove(domain.WorkBlocked) {
		return domain.ErrInvalidState
	}
	return s.repo.TransitionWorkOrder(ctx, id, order.State, string(domain.WorkBlocked), version, s.now().Format(time.RFC3339Nano))
}
func (s *WorkOrderService) Complete(ctx context.Context, id string, version int) error {
	order, err := s.repo.GetWorkOrder(ctx, id)
	if err != nil {
		return err
	}
	if !domain.WorkOrderState(order.State).CanMove(domain.WorkCompleted) {
		return fmt.Errorf("%w: state %s", domain.ErrInvalidState, order.State)
	}
	return s.repo.TransitionWorkOrder(ctx, id, order.State, string(domain.WorkCompleted), version, s.now().Format(time.RFC3339Nano))
}
func (s *WorkOrderService) List(ctx context.Context, farmID, state string, limit, offset int) (domain.Page[domain.WorkOrder], error) {
	return s.repo.ListWorkOrders(ctx, farmID, state, limit, offset)
}
