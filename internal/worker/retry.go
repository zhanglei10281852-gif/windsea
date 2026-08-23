package worker

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"time"
)

type RetryPolicy struct {
	Attempts  int
	BaseDelay time.Duration
}

func (p RetryPolicy) Run(ctx context.Context, fn func(context.Context) error) error {
	if p.Attempts < 1 {
		p.Attempts = 1
	}
	var last error
	for attempt := 0; attempt < p.Attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return errors.Join(domain.ErrCancelled, err)
		}
		if err := fn(ctx); err == nil {
			return nil
		} else {
			last = err
		}
		delay := p.BaseDelay * time.Duration(attempt+1)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return errors.Join(domain.ErrCancelled, ctx.Err())
		case <-timer.C:
		}
	}
	return last
}
func Backoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 8 {
		attempt = 8
	}
	return time.Duration(1<<attempt) * 100 * time.Millisecond
}
