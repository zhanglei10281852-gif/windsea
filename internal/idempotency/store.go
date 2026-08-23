package idempotency

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

type Record struct {
	Key, Scope string
	Response   []byte
	CreatedAt  time.Time
}
type Store struct {
	mu      sync.Mutex
	records map[string]Record
}

func New() *Store { return &Store{records: map[string]Record{}} }
func (s *Store) Execute(ctx context.Context, key, scope string, fn func(context.Context) (any, error)) (any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	if record, ok := s.records[key+"|"+scope]; ok {
		s.mu.Unlock()
		var value any
		err := json.Unmarshal(record.Response, &value)
		return value, err
	}
	s.mu.Unlock()
	value, err := fn(ctx)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.records[key+"|"+scope] = Record{Key: key, Scope: scope, Response: raw, CreatedAt: time.Now()}
	s.mu.Unlock()
	return value, nil
}
func (s *Store) Purge(before time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, record := range s.records {
		if record.CreatedAt.Before(before) {
			delete(s.records, key)
		}
	}
}
