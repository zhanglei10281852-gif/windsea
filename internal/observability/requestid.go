package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type requestIDKey struct{}

func NewRequestID() string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return "request-unknown"
	}
	return hex.EncodeToString(bytes)
}
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}
func RequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return "request-unknown"
}
