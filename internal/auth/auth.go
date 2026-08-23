package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

type Principal struct{ UserID, FarmID, Role string }
type key struct{}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, key{}, principal)
}
func PrincipalFrom(ctx context.Context) (Principal, bool) {
	value, ok := ctx.Value(key{}).(Principal)
	return value, ok
}
func RequireRole(ctx context.Context, roles ...string) error {
	principal, ok := PrincipalFrom(ctx)
	if !ok {
		return domain.ErrForbidden
	}
	for _, role := range roles {
		if strings.EqualFold(principal.Role, role) {
			return nil
		}
	}
	return domain.ErrForbidden
}
func HashToken(token string) string {
	digest := sha256.Sum256([]byte(token))
	return hex.EncodeToString(digest[:])
}
func ValidateSession(now time.Time, expires time.Time, revoked *time.Time) error {
	if revoked != nil || !now.Before(expires) {
		return fmt.Errorf("%w: session expired or revoked", domain.ErrForbidden)
	}
	return nil
}
