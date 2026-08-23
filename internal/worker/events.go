package worker

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sync"
)

type Event struct{ ID, Kind, Payload string }
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Event
	closed      bool
}

func NewEventBus() *EventBus { return &EventBus{subscribers: map[string][]chan Event{}} }
func (b *EventBus) Subscribe(kind string) (<-chan Event, func()) {
	b.mu.Lock()
	defer b.mu.Unlock()
	channel := make(chan Event, 8)
	if b.closed {
		close(channel)
		return channel, func() {}
	}
	b.subscribers[kind] = append(b.subscribers[kind], channel)
	cancel := func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		for i, item := range b.subscribers[kind] {
			if item == channel {
				b.subscribers[kind] = append(b.subscribers[kind][:i], b.subscribers[kind][i+1:]...)
				close(channel)
				break
			}
		}
	}
	return channel, cancel
}
func (b *EventBus) Publish(ctx context.Context, event Event) error {
	if err := ctx.Err(); err != nil {
		return errors.Join(domain.ErrCancelled, err)
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, channel := range b.subscribers[event.Kind] {
		select {
		case channel <- event:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
func (b *EventBus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for _, channels := range b.subscribers {
		for _, channel := range channels {
			close(channel)
		}
	}
	b.subscribers = map[string][]chan Event{}
}
