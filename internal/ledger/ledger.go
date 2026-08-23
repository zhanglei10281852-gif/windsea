package ledger

import (
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sort"
	"time"
)

type Entry struct {
	ID, Account string
	Delta       int
	At          time.Time
	Reference   string
}

func Balance(entries []Entry, account string) int {
	sum := 0
	for _, entry := range entries {
		if entry.Account == account {
			sum += entry.Delta
		}
	}
	return sum
}
func Validate(entries []Entry) error {
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.ID == "" || entry.Account == "" || seen[entry.ID] {
			return domain.ErrValidation
		}
		seen[entry.ID] = true
	}
	return nil
}
func Ordered(entries []Entry) []Entry {
	copy := append([]Entry(nil), entries...)
	sort.SliceStable(copy, func(i, j int) bool { return copy[i].At.Before(copy[j].At) })
	return copy
}
func RequireNonNegative(entries []Entry, account string) error {
	balance := 0
	for _, entry := range Ordered(entries) {
		if entry.Account != account {
			continue
		}
		balance += entry.Delta
		if balance < 0 {
			return fmt.Errorf("%w: account %s", domain.ErrConflict, account)
		}
	}
	return nil
}
