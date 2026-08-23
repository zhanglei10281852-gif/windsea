package service_test

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/windsea/internal/service"
	"testing"
)

func TestProcessBatchKeepsPartialResults(t *testing.T) {
	result := service.ProcessBatch(context.Background(), []string{"a", "b", "c"}, func(value string) string { return value }, func(_ context.Context, value string) error {
		if value == "b" {
			return errors.New("bad item")
		}
		return nil
	})
	if len(result.Successful) != 2 || len(result.Failed) != 1 {
		t.Fatalf("result=%+v", result)
	}
}
func TestProcessBatchStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := service.ProcessBatch(ctx, []string{"a", "b"}, func(value string) string { return value }, func(_ context.Context, value string) error { return nil })
	if len(result.Failed) != 1 {
		t.Fatalf("result=%+v", result)
	}
}
