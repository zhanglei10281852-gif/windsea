package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sync"
	"time"
)

type Heartbeat struct {
	WorkerID string
	At       time.Time
	Healthy  bool
	Details  map[string]string
}
type HeartbeatRegistry struct {
	mu    sync.RWMutex
	items map[string]Heartbeat
}

func NewHeartbeatRegistry() *HeartbeatRegistry {
	return &HeartbeatRegistry{items: map[string]Heartbeat{}}
}
func (r *HeartbeatRegistry) Record(ctx context.Context, h Heartbeat) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if h.WorkerID == "" || h.At.IsZero() {
		return domain.ErrValidation
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[h.WorkerID] = Heartbeat{WorkerID: h.WorkerID, At: h.At, Healthy: h.Healthy, Details: clone(h.Details)}
	return nil
}
func (r *HeartbeatRegistry) Healthy(now time.Time, ttl time.Duration) []Heartbeat {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := []Heartbeat{}
	for _, heartbeat := range r.items {
		if heartbeat.Healthy && now.Sub(heartbeat.At) <= ttl {
			heartbeat.Details = clone(heartbeat.Details)
			result = append(result, heartbeat)
		}
	}
	return result
}
func clone(input map[string]string) map[string]string {
	copy := map[string]string{}
	for key, value := range input {
		copy[key] = value
	}
	return copy
}
