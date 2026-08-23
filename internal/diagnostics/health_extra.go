package diagnostics

import (
	"context"
	"fmt"
	"time"
)

type HealthWindow struct {
	Start, End time.Time
	Checks     []Result
}

func BuildWindow(ctx context.Context, registry *Registry, start, end time.Time) HealthWindow {
	if end.Before(start) {
		end = start
	}
	results := Run(ctx, registry.Snapshot())
	return HealthWindow{Start: start, End: end, Checks: results}
}
func (w HealthWindow) Healthy() bool { return Healthy(w.Checks) }
func (w HealthWindow) FailedNames() []string {
	names := []string{}
	for _, check := range w.Checks {
		if !check.OK {
			names = append(names, check.Name)
		}
	}
	return names
}
func (w HealthWindow) Duration() time.Duration { return w.End.Sub(w.Start) }
func (w HealthWindow) String() string {
	state := "healthy"
	if !w.Healthy() {
		state = "degraded"
	}
	return fmt.Sprintf("%s %s", state, Summarize(w.Checks))
}
func RegisterDefaults(registry *Registry, ping func(context.Context) error) error {
	defaults := []Probe{{Name: "database", Run: ping}, {Name: "scheduler", Run: func(context.Context) error { return nil }}, {Name: "worker", Run: func(context.Context) error { return nil }}}
	for _, probe := range defaults {
		if err := registry.Register(probe); err != nil {
			return err
		}
	}
	return nil
}
func CheckAll(ctx context.Context, registry *Registry) error {
	results := Run(ctx, registry.Snapshot())
	if !Healthy(results) {
		return fmt.Errorf("health checks failed: %s", Summarize(results))
	}
	return nil
}
func MustCheck(ctx context.Context, registry *Registry) error {
	if err := CheckAll(ctx, registry); err != nil {
		return err
	}
	return nil
}
