package stream

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sync"
)

type Broadcaster[T any] struct {
	mu        sync.RWMutex
	listeners map[int]chan T
	next      int
	closed    bool
}

func New[T any]() *Broadcaster[T] { return &Broadcaster[T]{listeners: map[int]chan T{}} }
func (b *Broadcaster[T]) Listen(buffer int) (int, <-chan T, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return 0, nil, domain.ErrCancelled
	}
	id := b.next
	b.next++
	b.listeners[id] = make(chan T, buffer)
	return id, b.listeners[id], nil
}
func (b *Broadcaster[T]) Remove(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if channel, ok := b.listeners[id]; ok {
		delete(b.listeners, id)
		close(channel)
	}
}
func (b *Broadcaster[T]) Send(ctx context.Context, value T) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, channel := range b.listeners {
		select {
		case channel <- value:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
func (b *Broadcaster[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for id, channel := range b.listeners {
		close(channel)
		delete(b.listeners, id)
	}
}
