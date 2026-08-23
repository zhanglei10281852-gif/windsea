package diagnostics

import (
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"strings"
	"time"
)

type Incident struct {
	ID, FarmID, Severity, Summary string
	OpenedAt, ClosedAt            *time.Time
	Tags                          []string
}

func (i Incident) Validate() error {
	if i.ID == "" || i.FarmID == "" || i.Summary == "" {
		return domain.ErrValidation
	}
	if i.Severity != "critical" && i.Severity != "warning" && i.Severity != "info" {
		return fmt.Errorf("%w: severity", domain.ErrValidation)
	}
	return nil
}
func NormalizeTags(tags []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}
func Open(i Incident, at time.Time) (Incident, error) {
	if err := i.Validate(); err != nil {
		return i, err
	}
	i.OpenedAt = &at
	i.Tags = NormalizeTags(i.Tags)
	return i, nil
}
func Close(i Incident, at time.Time) (Incident, error) {
	if i.OpenedAt == nil || at.Before(*i.OpenedAt) {
		return i, domain.ErrInvalidState
	}
	i.ClosedAt = &at
	return i, nil
}
