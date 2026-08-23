package cache

import (
	"sync"
	"time"
)

type Entry[T any] struct {
	Value     T
	ExpiresAt time.Time
}
type Cache[T any] struct {
	mu    sync.RWMutex
	items map[string]Entry[T]
	now   func() time.Time
}

func New[T any](now func() time.Time) *Cache[T] {
	return &Cache[T]{items: map[string]Entry[T]{}, now: now}
}
func (c *Cache[T]) Set(key string, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Entry[T]{Value: value, ExpiresAt: c.now().Add(ttl)}
}
func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()
	var zero T
	if !ok {
		return zero, false
	}
	if !c.now().Before(entry.ExpiresAt) {
		c.Delete(key)
		return zero, false
	}
	return entry.Value, true
}
func (c *Cache[T]) Delete(key string) { c.mu.Lock(); defer c.mu.Unlock(); delete(c.items, key) }
func (c *Cache[T]) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	for key, entry := range c.items {
		if !now.Before(entry.ExpiresAt) {
			delete(c.items, key)
		}
	}
}
