package planning

import (
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sort"
	"time"
)

type Assignment struct {
	WorkOrderID, TurbineID, TechnicianID string
	Start, End                           time.Time
	Priority                             int
}

func Validate(assignments []Assignment) error {
	for _, assignment := range assignments {
		if assignment.WorkOrderID == "" || assignment.TurbineID == "" || assignment.TechnicianID == "" || !assignment.End.After(assignment.Start) {
			return domain.ErrValidation
		}
	}
	sort.Slice(assignments, func(i, j int) bool { return assignments[i].Start.Before(assignments[j].Start) })
	for i := 1; i < len(assignments); i++ {
		if assignments[i].TechnicianID == assignments[i-1].TechnicianID && assignments[i].Start.Before(assignments[i-1].End) {
			return fmt.Errorf("%w: technician overlap", domain.ErrConflict)
		}
	}
	return nil
}
func NextSlot(existing []Assignment, duration time.Duration, from time.Time) time.Time {
	candidate := from
	for _, assignment := range existing {
		if candidate.Before(assignment.End) && candidate.Add(duration).After(assignment.Start) {
			candidate = assignment.End
		}
	}
	return candidate
}
