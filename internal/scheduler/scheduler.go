package scheduler

import (
	"context"
	"sort"
	"sync"
	"time"
)

type Task struct {
	ID  string
	At  time.Time
	Run func(context.Context) error
}
type Scheduler struct {
	mu    sync.Mutex
	tasks []Task
}

func New() *Scheduler { return &Scheduler{tasks: []Task{}} }
func (s *Scheduler) Add(task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, task)
	sort.Slice(s.tasks, func(i, j int) bool { return s.tasks[i].At.Before(s.tasks[j].At) })
}
func (s *Scheduler) Due(now time.Time) []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	due := []Task{}
	pending := []Task{}
	for _, task := range s.tasks {
		if !task.At.After(now) {
			due = append(due, task)
		} else {
			pending = append(pending, task)
		}
	}
	s.tasks = pending
	return due
}
func (s *Scheduler) RunDue(ctx context.Context, now time.Time) []error {
	errors := []error{}
	for _, task := range s.Due(now) {
		if err := task.Run(ctx); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
