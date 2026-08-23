package consistency

import (
	"sort"
	"time"
)

type Versioned struct {
	Version int
	At      time.Time
	Actor   string
}

func Latest(history []Versioned) (Versioned, bool) {
	if len(history) == 0 {
		return Versioned{}, false
	}
	copy := append([]Versioned(nil), history...)
	sort.Slice(copy, func(i, j int) bool {
		if copy[i].Version == copy[j].Version {
			return copy[i].At.Before(copy[j].At)
		}
		return copy[i].Version < copy[j].Version
	})
	return copy[len(copy)-1], true
}
func Monotonic(history []Versioned) bool {
	latest := 0
	for _, entry := range history {
		if entry.Version <= latest {
			return false
		}
		latest = entry.Version
	}
	return true
}
