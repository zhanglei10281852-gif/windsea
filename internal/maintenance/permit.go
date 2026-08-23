package maintenance

import (
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"time"
)

type Permit struct {
	ID, WorkOrderID, IssuedBy string
	IssuedAt, ExpiresAt       time.Time
	Revoked                   bool
}

func (p Permit) Valid(now time.Time) error {
	if p.ID == "" || p.WorkOrderID == "" || p.IssuedBy == "" {
		return domain.ErrValidation
	}
	if p.Revoked || !now.Before(p.ExpiresAt) || p.ExpiresAt.Before(p.IssuedAt) {
		return fmt.Errorf("%w: permit", domain.ErrForbidden)
	}
	return nil
}
func Extend(p Permit, now time.Time, extension time.Duration) (Permit, error) {
	if err := p.Valid(now); err != nil {
		return p, err
	}
	if extension <= 0 || extension > 24*time.Hour {
		return p, domain.ErrValidation
	}
	p.ExpiresAt = p.ExpiresAt.Add(extension)
	return p, nil
}
func Revoke(p Permit) Permit { p.Revoked = true; return p }
