package worker_test

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/windsea/internal/worker"
	"testing"
	"time"
)

func TestRetryPolicyRetriesAndSucceeds(t *testing.T) {
	attempts := 0
	policy := worker.RetryPolicy{Attempts: 3, BaseDelay: time.Millisecond}
	if err := policy.Run(context.Background(), func(context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary")
		}
		return nil
	}); err != nil || attempts != 3 {
		t.Fatalf("err=%v attempts=%d", err, attempts)
	}
}
func TestRetryPolicyHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	policy := worker.RetryPolicy{Attempts: 3, BaseDelay: time.Millisecond}
	if err := policy.Run(ctx, func(context.Context) error { return nil }); err == nil {
		t.Fatal("cancelled job ran")
	}
}
func TestQueueProcessesJobs(t *testing.T) {
	queue := worker.NewQueue(2, 2)
	defer queue.Close()
	done := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		if err := queue.Submit(context.Background(), worker.Job{ID: string(rune('a' + i)), Run: func(context.Context) error { done <- struct{}{}; return nil }}); err != nil {
			t.Fatal(err)
		}
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("job not processed")
	}
}
func TestLeaseManagerOwnership(t *testing.T) {
	manager := worker.NewLeaseManager()
	lease, err := manager.Acquire(context.Background(), "job", "worker-a", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := manager.Acquire(context.Background(), "job", "worker-b", time.Second); err == nil {
		t.Fatal("second owner acquired lease")
	}
	if !manager.Release(lease.ID, "worker-a") {
		t.Fatal("release failed")
	}
	if _, err := manager.Acquire(context.Background(), "job", "worker-b", time.Second); err != nil {
		t.Fatal(err)
	}
}
