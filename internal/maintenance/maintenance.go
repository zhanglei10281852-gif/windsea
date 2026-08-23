package maintenance

import (
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sort"
	"time"
)

type ChecklistItem struct {
	Code, Label string
	Required    bool
	Completed   bool
	CompletedAt *time.Time
}
type Checklist struct {
	ID, WorkOrderID string
	Items           []ChecklistItem
}

func (c Checklist) Validate() error {
	if c.ID == "" || c.WorkOrderID == "" || len(c.Items) == 0 {
		return domain.ErrValidation
	}
	seen := map[string]bool{}
	for _, item := range c.Items {
		if item.Code == "" || seen[item.Code] {
			return fmt.Errorf("%w: duplicate checklist item", domain.ErrValidation)
		}
		seen[item.Code] = true
	}
	return nil
}
func (c Checklist) Complete() bool {
	for _, item := range c.Items {
		if item.Required && !item.Completed {
			return false
		}
	}
	return true
}
func SortItems(items []ChecklistItem) []ChecklistItem {
	copy := append([]ChecklistItem(nil), items...)
	sort.Slice(copy, func(i, j int) bool { return copy[i].Code < copy[j].Code })
	return copy
}
func Mark(items []ChecklistItem, code string, at time.Time) ([]ChecklistItem, error) {
	result := append([]ChecklistItem(nil), items...)
	for i := range result {
		if result[i].Code == code {
			result[i].Completed = true
			result[i].CompletedAt = &at
			return result, nil
		}
	}
	return result, domain.ErrNotFound
}
