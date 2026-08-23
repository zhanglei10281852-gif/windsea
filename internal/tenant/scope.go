package tenant

import (
	"context"
	"fmt"

	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

type key struct{}

func WithFarm(ctx context.Context, farmID string) context.Context {
	return context.WithValue(ctx, key{}, farmID)
}
func FarmID(ctx context.Context) (string, error) {
	value, ok := ctx.Value(key{}).(string)
	if !ok || value == "" {
		return "", fmt.Errorf("%w: farm scope missing", domain.ErrForbidden)
	}
	return value, nil
}
func Ensure(ctx context.Context, expected string) error {
	actual, err := FarmID(ctx)
	if err != nil {
		return err
	}
	if actual != expected {
		return fmt.Errorf("%w: farm scope mismatch", domain.ErrForbidden)
	}
	return nil
}
