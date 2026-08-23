package consistency

import (
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"time"
)

func PartAvailable(part domain.Part, quantity int) error {
	if quantity <= 0 {
		return domain.ErrValidation
	}
	if part.OnHand < 0 || part.Reserved < 0 || part.Reserved > part.OnHand {
		return fmt.Errorf("%w: inventory invariant", domain.ErrConflict)
	}
	if part.OnHand-part.Reserved < quantity {
		return domain.ErrCapacity
	}
	return nil
}
func CampaignWindow(start, end time.Time) error {
	if !end.After(start) {
		return fmt.Errorf("%w: campaign window", domain.ErrValidation)
	}
	return nil
}
func WorkOrderOwnership(order domain.WorkOrder, farmID string) error {
	if order.FarmID != farmID {
		return domain.ErrForbidden
	}
	if order.ID == "" || order.Title == "" {
		return domain.ErrValidation
	}
	return nil
}
