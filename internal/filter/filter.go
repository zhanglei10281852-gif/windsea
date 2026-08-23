package filter

import (
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sort"
	"strings"
)

type WorkOrderFilter struct {
	State                    string
	MinPriority, MaxPriority int
	Search                   string
}

func ApplyWorkOrders(items []domain.WorkOrder, filter WorkOrderFilter) []domain.WorkOrder {
	result := make([]domain.WorkOrder, 0, len(items))
	for _, item := range items {
		if filter.State != "" && item.State != filter.State {
			continue
		}
		if filter.MinPriority > 0 && item.Priority < filter.MinPriority {
			continue
		}
		if filter.MaxPriority > 0 && item.Priority > filter.MaxPriority {
			continue
		}
		if filter.Search != "" && !strings.Contains(strings.ToLower(item.Title), strings.ToLower(filter.Search)) {
			continue
		}
		result = append(result, item)
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Priority == result[j].Priority {
			return result[i].CreatedAt.Before(result[j].CreatedAt)
		}
		return result[i].Priority > result[j].Priority
	})
	return result
}
func Page[T any](items []T, limit, offset int) []T {
	if limit < 1 {
		limit = 1
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []T{}
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]T(nil), items[offset:end]...)
}
