package service

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"log/slog"
)

type QueryService struct {
	repo *repository.Repository
	log  *slog.Logger
}

func (s *QueryService) WorkOrders(ctx context.Context, farmID, state string, limit, offset int) (domain.Page[domain.WorkOrder], error) {
	return s.repo.ListWorkOrders(ctx, farmID, state, limit, offset)
}
func (s *QueryService) Alerts(ctx context.Context, farmID string) ([]domain.Alert, error) {
	return s.repo.ListOpenAlerts(ctx, farmID)
}
func (s *QueryService) Audit(ctx context.Context, farmID, objectID string) ([]domain.AuditEvent, error) {
	return s.repo.ListAudit(ctx, farmID, objectID)
}
