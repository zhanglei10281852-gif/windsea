package observability_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/observability"
	"testing"
)

func TestRequestIDContext(t *testing.T) {
	ctx := observability.WithRequestID(context.Background(), "abc")
	if observability.RequestID(ctx) != "abc" {
		t.Fatal("request id lost")
	}
	if observability.RequestID(context.Background()) != "request-unknown" {
		t.Fatal("unexpected default")
	}
}
