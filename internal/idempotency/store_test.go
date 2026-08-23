package idempotency_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/idempotency"
	"testing"
)

func TestIdempotencyReturnsOriginalValue(t *testing.T) {
	store := idempotency.New()
	calls := 0
	run := func(context.Context) (any, error) { calls++; return map[string]string{"status": "ok"}, nil }
	first, err := store.Execute(context.Background(), "key", "scope", run)
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Execute(context.Background(), "key", "scope", run)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("calls=%d", calls)
	}
	if first.(map[string]string)["status"] != "ok" || second.(map[string]any)["status"] != "ok" {
		t.Fatal("values differ")
	}
}
