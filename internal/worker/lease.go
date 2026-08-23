package worker

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Lease struct {
	ID        string
	ExpiresAt time.Time
	Owner     string
}
type LeaseManager struct {
	mu     sync.Mutex
	leases map[string]Lease
}

func NewLeaseManager() *LeaseManager { return &LeaseManager{leases: map[string]Lease{}} }
func (m *LeaseManager) Acquire(ctx context.Context, id, owner string, ttl time.Duration) (Lease, error) {
	if err := ctx.Err(); err != nil {
		return Lease{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	if existing, ok := m.leases[id]; ok && existing.ExpiresAt.After(now) {
		return Lease{}, fmt.Errorf("lease held by %s", existing.Owner)
	}
	lease := Lease{ID: id, Owner: owner, ExpiresAt: now.Add(ttl)}
	m.leases[id] = lease
	return lease, nil
}
func (m *LeaseManager) Renew(ctx context.Context, id, owner string, ttl time.Duration) (Lease, error) {
	if err := ctx.Err(); err != nil {
		return Lease{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	lease, ok := m.leases[id]
	if !ok || lease.Owner != owner || !lease.ExpiresAt.After(time.Now()) {
		return Lease{}, fmt.Errorf("lease unavailable")
	}
	lease.ExpiresAt = time.Now().Add(ttl)
	m.leases[id] = lease
	return lease, nil
}
func (m *LeaseManager) Release(id, owner string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	lease, ok := m.leases[id]
	if !ok || lease.Owner != owner {
		return false
	}
	delete(m.leases, id)
	return true
}
