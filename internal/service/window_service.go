package service

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"time"
)

func ValidateMaintenanceWindow(window domain.MaintenanceWindow, now time.Time) error {
	if window.ID == "" || window.FarmID == "" || window.TurbineID == "" {
		return domain.ErrValidation
	}
	if window.EndsAt.Before(window.StartsAt) || window.EndsAt.Equal(window.StartsAt) {
		return domain.ErrValidation
	}
	if window.EndsAt.Before(now) {
		return domain.ErrValidation
	}
	return nil
}
func ScheduleWindow(ctx context.Context, service *HandoffService, window domain.MaintenanceWindow) error {
	if err := ValidateMaintenanceWindow(window, service.now()); err != nil {
		return err
	}
	return service.Window(ctx, window)
}
