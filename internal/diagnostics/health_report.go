package diagnostics

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Report struct {
	GeneratedAt time.Time
	Window      HealthWindow
	Labels      map[string]string
}

func NewReport(window HealthWindow, at time.Time) Report {
	return Report{GeneratedAt: at, Window: window, Labels: map[string]string{}}
}
func (r *Report) Label(key, value string) {
	if r.Labels == nil {
		r.Labels = map[string]string{}
	}
	if strings.TrimSpace(key) != "" {
		r.Labels[key] = value
	}
}
func (r Report) LabelValue(key string) string { return r.Labels[key] }
func (r Report) FailedCount() int             { return len(r.Window.FailedNames()) }
func (r Report) PassedCount() int             { return len(r.Window.Checks) - r.FailedCount() }
func (r Report) Score() float64 {
	if len(r.Window.Checks) == 0 {
		return 1
	}
	return float64(r.PassedCount()) / float64(len(r.Window.Checks))
}
func (r Report) SortedLabels() []string {
	keys := make([]string, 0, len(r.Labels))
	for key := range r.Labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
func (r Report) Summary() string {
	parts := []string{fmt.Sprintf("score=%.2f", r.Score()), fmt.Sprintf("passed=%d", r.PassedCount()), fmt.Sprintf("failed=%d", r.FailedCount())}
	for _, key := range r.SortedLabels() {
		parts = append(parts, key+"="+r.Labels[key])
	}
	return strings.Join(parts, " ")
}
func (r Report) Since(now time.Time) time.Duration { return now.Sub(r.GeneratedAt) }
