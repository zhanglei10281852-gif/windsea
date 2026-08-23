package authorization

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/auth"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

type Policy struct {
	Action     string
	Roles      []string
	FarmScoped bool
}

func Check(ctx context.Context, policy Policy, objectFarm string) error {
	if policy.Action == "" {
		return domain.ErrValidation
	}
	if policy.FarmScoped {
		principal, ok := auth.PrincipalFrom(ctx)
		if !ok || principal.FarmID != objectFarm {
			return domain.ErrForbidden
		}
	}
	for _, role := range policy.Roles {
		if principal, ok := auth.PrincipalFrom(ctx); ok && principal.Role == role {
			return nil
		}
	}
	return fmt.Errorf("%w: policy %s", domain.ErrForbidden, policy.Action)
}
func Any(ctx context.Context, policies []Policy, farm string) error {
	for _, policy := range policies {
		if Check(ctx, policy, farm) == nil {
			return nil
		}
	}
	return domain.ErrForbidden
}
