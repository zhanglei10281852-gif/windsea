package scheduler_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/scheduler"
	"testing"
	"time"
)

func TestSchedulerOrdersDueTasks(t *testing.T) {
	s := scheduler.New()
	ran := []string{}
	now := time.Now()
	s.Add(scheduler.Task{ID: "late", At: now.Add(time.Hour), Run: func(context.Context) error { ran = append(ran, "late"); return nil }})
	s.Add(scheduler.Task{ID: "early", At: now.Add(-time.Minute), Run: func(context.Context) error { ran = append(ran, "early"); return nil }})
	if errs := s.RunDue(context.Background(), now); len(errs) != 0 {
		t.Fatal(errs)
	}
	if len(ran) != 1 || ran[0] != "early" {
		t.Fatalf("ran=%v", ran)
	}
	if len(s.Due(now)) != 0 {
		t.Fatal("task repeated")
	}
}
