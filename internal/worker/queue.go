package worker

import (
	"context"
	"sync"
)

type Job struct {
	ID  string
	Run func(context.Context) error
}
type Queue struct {
	jobs chan Job
	wg   sync.WaitGroup
}

func NewQueue(size, workers int) *Queue {
	q := &Queue{jobs: make(chan Job, size)}
	for i := 0; i < workers; i++ {
		q.wg.Add(1)
		go q.loop()
	}
	return q
}
func (q *Queue) Submit(ctx context.Context, job Job) error {
	select {
	case q.jobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (q *Queue) Close() { close(q.jobs); q.wg.Wait() }
func (q *Queue) loop() {
	defer q.wg.Done()
	for job := range q.jobs {
		_ = job.Run(context.Background())
	}
}
