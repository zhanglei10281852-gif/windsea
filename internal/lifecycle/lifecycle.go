package lifecycle

import (
	"context"
	"errors"
	"sync"
)

type Manager struct {
	mu      sync.Mutex
	closed  bool
	closers []func() error
}

func New() *Manager { return &Manager{closers: []func() error{}} }
func (m *Manager) Register(closer func() error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return errors.New("lifecycle closed")
	}
	m.closers = append(m.closers, closer)
	return nil
}
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	m.closed = true
	closers := append([]func() error(nil), m.closers...)
	m.mu.Unlock()
	var combined error
	for i := len(closers) - 1; i >= 0; i-- {
		done := make(chan error, 1)
		go func(fn func() error) { done <- fn() }(closers[i])
		select {
		case err := <-done:
			combined = errors.Join(combined, err)
		case <-ctx.Done():
			return errors.Join(combined, ctx.Err())
		}
	}
	return combined
}
