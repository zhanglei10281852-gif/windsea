package metrics

import (
	"sync"
	"sync/atomic"
)

type Counters struct {
	requests      atomic.Int64
	failures      atomic.Int64
	workerRetries atomic.Int64
	mu            sync.Mutex
	byRoute       map[string]int64
}

func New() *Counters { return &Counters{byRoute: map[string]int64{}} }
func (c *Counters) Request(route string) {
	c.requests.Add(1)
	c.mu.Lock()
	c.byRoute[route]++
	c.mu.Unlock()
}
func (c *Counters) Failure() { c.failures.Add(1) }
func (c *Counters) Retry()   { c.workerRetries.Add(1) }
func (c *Counters) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	copy := map[string]int64{"requests": c.requests.Load(), "failures": c.failures.Load(), "worker_retries": c.workerRetries.Load()}
	for route, count := range c.byRoute {
		copy["route:"+route] = count
	}
	return copy
}
