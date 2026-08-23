package diagnostics

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

type Probe struct {
	Name string
	Run  func(context.Context) error
}
type Result struct {
	Name     string
	OK       bool
	Error    string
	Duration time.Duration
}

func Run(ctx context.Context, probes []Probe) []Result {
	results := make([]Result, 0, len(probes))
	for _, probe := range probes {
		started := time.Now()
		err := probe.Run(ctx)
		result := Result{Name: probe.Name, OK: err == nil, Duration: time.Since(started)}
		if err != nil {
			result.Error = err.Error()
		}
		results = append(results, result)
	}
	return results
}
func Healthy(results []Result) bool {
	for _, result := range results {
		if !result.OK {
			return false
		}
	}
	return true
}
func RuntimeInfo() map[string]string {
	return map[string]string{"go": runtime.Version(), "os": runtime.GOOS, "arch": runtime.GOARCH}
}
func Summarize(results []Result) string {
	copy := append([]Result(nil), results...)
	sort.Slice(copy, func(i, j int) bool { return copy[i].Name < copy[j].Name })
	parts := make([]string, 0, len(copy))
	for _, result := range copy {
		state := "ok"
		if !result.OK {
			state = "failed"
		}
		parts = append(parts, fmt.Sprintf("%s=%s", result.Name, state))
	}
	return strings.Join(parts, ",")
}

type Registry struct {
	mu     sync.RWMutex
	probes map[string]Probe
}

func NewRegistry() *Registry { return &Registry{probes: map[string]Probe{}} }
func (r *Registry) Register(probe Probe) error {
	if probe.Name == "" || probe.Run == nil {
		return domain.ErrValidation
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.probes[probe.Name]; exists {
		return domain.ErrConflict
	}
	r.probes[probe.Name] = probe
	return nil
}
func (r *Registry) Snapshot() []Probe {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Probe, 0, len(r.probes))
	for _, probe := range r.probes {
		result = append(result, probe)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}
