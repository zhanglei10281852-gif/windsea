package recovery_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/recovery"
	"testing"
	"time"
)

func TestCheckpointSequence(t *testing.T) {
	now := time.Now()
	first, err := recovery.New("c1", "f", "received", 1, now)
	if err != nil {
		t.Fatal(err)
	}
	second, _ := recovery.New("c2", "f", "completed", 2, now.Add(time.Minute))
	if err := recovery.Apply(context.Background(), first, second); err != nil {
		t.Fatal(err)
	}
	third, _ := recovery.New("c3", "f", "bad", 4, now.Add(2*time.Minute))
	if err := recovery.Apply(context.Background(), second, third); err == nil {
		t.Fatal("sequence gap accepted")
	}
}
