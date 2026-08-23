package lifecycle_test

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/windsea/internal/lifecycle"
	"testing"
	"time"
)

func TestShutdownRunsReverseOrder(t *testing.T) {
	manager := lifecycle.New()
	order := []string{}
	if err := manager.Register(func() error { order = append(order, "first"); return nil }); err != nil {
		t.Fatal(err)
	}
	if err := manager.Register(func() error { order = append(order, "second"); return nil }); err != nil {
		t.Fatal(err)
	}
	if err := manager.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if order[0] != "second" || order[1] != "first" {
		t.Fatalf("order=%v", order)
	}
	if err := manager.Register(func() error { return nil }); err == nil {
		t.Fatal("registered after shutdown")
	}
}
func TestShutdownHonorsDeadline(t *testing.T) {
	manager := lifecycle.New()
	if err := manager.Register(func() error { time.Sleep(30 * time.Millisecond); return errors.New("late") }); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if err := manager.Shutdown(ctx); err == nil {
		t.Fatal("deadline ignored")
	}
}
