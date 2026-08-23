package metrics_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/metrics"
	"sync"
	"testing"
)

func TestCountersConcurrent(t *testing.T) {
	counters := metrics.New()
	var group sync.WaitGroup
	for i := 0; i < 20; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for j := 0; j < 50; j++ {
				counters.Request("/v1/work-orders")
				counters.Failure()
				counters.Retry()
			}
		}()
	}
	group.Wait()
	snapshot := counters.Snapshot()
	if snapshot["requests"] != 1000 || snapshot["failures"] != 1000 || snapshot["worker_retries"] != 1000 {
		t.Fatalf("snapshot=%v", snapshot)
	}
}
