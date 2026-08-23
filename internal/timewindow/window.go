package timewindow

import (
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sort"
	"time"
)

type Window struct {
	ID, TurbineID string
	Start, End    time.Time
}

func (w Window) Valid() error {
	if w.ID == "" || w.TurbineID == "" || !w.End.After(w.Start) {
		return domain.ErrValidation
	}
	return nil
}
func Overlaps(a, b Window) bool {
	return a.TurbineID == b.TurbineID && a.Start.Before(b.End) && b.Start.Before(a.End)
}
func Check(windows []Window) error {
	copy := append([]Window(nil), windows...)
	sort.Slice(copy, func(i, j int) bool { return copy[i].Start.Before(copy[j].Start) })
	for _, window := range copy {
		if err := window.Valid(); err != nil {
			return err
		}
	}
	for i := 1; i < len(copy); i++ {
		if Overlaps(copy[i-1], copy[i]) {
			return fmt.Errorf("%w: overlap", domain.ErrConflict)
		}
	}
	return nil
}
func Remaining(window Window, now time.Time) time.Duration {
	if now.Before(window.Start) {
		return window.End.Sub(window.Start)
	}
	if now.After(window.End) {
		return 0
	}
	return window.End.Sub(now)
}
