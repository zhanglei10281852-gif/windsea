package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/service"
	"log/slog"
	"sync"
	"time"
)

type Supervisor struct {
	services *service.Registry
	logger   *slog.Logger
	interval time.Duration
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewSupervisor(services *service.Registry, logger *slog.Logger, interval time.Duration) *Supervisor {
	return &Supervisor{services: services, logger: logger, interval: interval}
}
func (s *Supervisor) Start(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	s.cancel = cancel
	s.wg.Add(2)
	go s.alertLoop(ctx)
	go s.dispatchLoop(ctx)
}
func (s *Supervisor) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}
func (s *Supervisor) alertLoop(ctx context.Context) {
	defer s.wg.Done()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.logger.Debug("alert worker tick")
		}
	}
}
func (s *Supervisor) dispatchLoop(ctx context.Context) {
	defer s.wg.Done()
	ticker := time.NewTicker(s.interval * 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.logger.Debug("dispatch worker tick")
		}
	}
}
