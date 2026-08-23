package maintenance_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/maintenance"
	"testing"
	"time"
)

func TestChecklistCompletion(t *testing.T) {
	now := time.Now()
	items := []maintenance.ChecklistItem{{Code: "isolate", Required: true}, {Code: "inspect", Required: true}}
	check := maintenance.Checklist{ID: "c", WorkOrderID: "wo", Items: items}
	if err := check.Validate(); err != nil {
		t.Fatal(err)
	}
	if check.Complete() {
		t.Fatal("incomplete checklist marked complete")
	}
	updated, err := maintenance.Mark(items, "isolate", now)
	if err != nil {
		t.Fatal(err)
	}
	updated, err = maintenance.Mark(updated, "inspect", now)
	if err != nil {
		t.Fatal(err)
	}
	check.Items = updated
	if !check.Complete() {
		t.Fatal("completed checklist rejected")
	}
	if _, err := maintenance.Mark(updated, "missing", now); err == nil {
		t.Fatal("missing item accepted")
	}
}
