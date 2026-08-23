package limits

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sync"
	"time"
)

type WindowCounter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	count  int
	reset  time.Time
	now    func() time.Time
}

func New(limit int, window time.Duration, now func() time.Time) *WindowCounter {
	return &WindowCounter{limit: limit, window: window, now: now, reset: now().Add(window)}
}
func (c *WindowCounter) Allow(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	if !now.Before(c.reset) {
		c.count = 0
		c.reset = now.Add(c.window)
	}
	if c.count >= c.limit {
		return domain.ErrCapacity
	}
	c.count++
	return nil
}
func (c *WindowCounter) Remaining() int { c.mu.Lock(); defer c.mu.Unlock(); return c.limit - c.count }
