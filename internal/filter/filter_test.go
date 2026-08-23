package filter_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/filter"
	"testing"
	"time"
)

func TestWorkOrderFilter(t *testing.T) {
	now := time.Now()
	items := []domain.WorkOrder{{ID: "1", Title: "Blade repair", State: "queued", Priority: 2, CreatedAt: now}, {ID: "2", Title: "Gearbox repair", State: "queued", Priority: 5, CreatedAt: now.Add(time.Minute)}, {ID: "3", Title: "Blade inspection", State: "completed", Priority: 5, CreatedAt: now}}
	result := filter.ApplyWorkOrders(items, filter.WorkOrderFilter{State: "queued", Search: "repair"})
	if len(result) != 2 || result[0].ID != "2" {
		t.Fatalf("result=%+v", result)
	}
}
func TestPageCopiesSlice(t *testing.T) {
	items := []int{1, 2, 3, 4}
	page := filter.Page(items, 2, 1)
	page[0] = 99
	if items[1] == 99 {
		t.Fatal("page shares backing array")
	}
	if len(filter.Page(items, 10, 10)) != 0 {
		t.Fatal("out of range page")
	}
}
