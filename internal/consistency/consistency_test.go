package consistency_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/consistency"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"testing"
	"time"
)

func TestInventoryInvariant(t *testing.T) {
	part := domain.Part{OnHand: 10, Reserved: 4}
	if err := consistency.PartAvailable(part, 6); err != nil {
		t.Fatal(err)
	}
	if err := consistency.PartAvailable(part, 7); err == nil {
		t.Fatal("capacity ignored")
	}
	if err := consistency.PartAvailable(domain.Part{OnHand: 1, Reserved: 2}, 1); err == nil {
		t.Fatal("invalid inventory accepted")
	}
}
func TestWindowAndOwnership(t *testing.T) {
	start := time.Now()
	if err := consistency.CampaignWindow(start, start.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := consistency.CampaignWindow(start, start); err == nil {
		t.Fatal("zero window accepted")
	}
	order := domain.WorkOrder{ID: "wo", FarmID: "f", Title: "repair"}
	if err := consistency.WorkOrderOwnership(order, "f"); err != nil {
		t.Fatal(err)
	}
	if err := consistency.WorkOrderOwnership(order, "other"); err == nil {
		t.Fatal("wrong farm accepted")
	}
}
func TestVersionHistory(t *testing.T) {
	history := []consistency.Versioned{{Version: 2, At: time.Now(), Actor: "b"}, {Version: 1, At: time.Now().Add(-time.Hour), Actor: "a"}}
	latest, ok := consistency.Latest(history)
	if !ok || latest.Version != 2 {
		t.Fatalf("latest=%+v", latest)
	}
	if !consistency.Monotonic([]consistency.Versioned{{Version: 1}, {Version: 2}}) {
		t.Fatal("monotonic failed")
	}
	if consistency.Monotonic([]consistency.Versioned{{Version: 2}, {Version: 1}}) {
		t.Fatal("nonmonotonic accepted")
	}
}
