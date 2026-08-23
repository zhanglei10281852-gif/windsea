package planning_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/planning"
	"testing"
	"time"
)

func TestAssignmentValidation(t *testing.T) {
	now := time.Now()
	assignments := []planning.Assignment{{WorkOrderID: "1", TurbineID: "t", TechnicianID: "tech", Start: now, End: now.Add(time.Hour)}, {WorkOrderID: "2", TurbineID: "t2", TechnicianID: "tech", Start: now.Add(2 * time.Hour), End: now.Add(3 * time.Hour)}}
	if err := planning.Validate(assignments); err != nil {
		t.Fatal(err)
	}
	assignments[1].Start = now.Add(30 * time.Minute)
	if err := planning.Validate(assignments); err == nil {
		t.Fatal("overlap accepted")
	}
}
func TestNextSlot(t *testing.T) {
	now := time.Now()
	existing := []planning.Assignment{{Start: now, End: now.Add(time.Hour)}}
	if got := planning.NextSlot(existing, 30*time.Minute, now); !got.Equal(now.Add(time.Hour)) {
		t.Fatalf("slot=%v", got)
	}
}
