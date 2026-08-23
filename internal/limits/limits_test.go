package limits_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/limits"
	"testing"
	"time"
)

func TestWindowCounter(t *testing.T) {
	now := time.Now()
	current := now
	counter := limits.New(2, time.Minute, func() time.Time { return current })
	if err := counter.Allow(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := counter.Allow(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := counter.Allow(context.Background()); err == nil {
		t.Fatal("limit ignored")
	}
	current = current.Add(time.Minute)
	if err := counter.Allow(context.Background()); err != nil {
		t.Fatal(err)
	}
	if counter.Remaining() != 1 {
		t.Fatal("remaining")
	}
}
