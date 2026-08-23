package reporting

import (
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"sort"
	"time"
)

type AlertRow struct {
	ID, Code, Severity, State string
	OccurredAt                time.Time
}

func AlertRows(alerts []domain.Alert) []AlertRow {
	rows := make([]AlertRow, 0, len(alerts))
	for _, alert := range alerts {
		when := time.Time{}
		if alert.OccurredAt != nil {
			when = *alert.OccurredAt
		}
		rows = append(rows, AlertRow{ID: alert.ID, Code: alert.Code, Severity: alert.Severity, State: alert.State, OccurredAt: when})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].OccurredAt.After(rows[j].OccurredAt) })
	return rows
}
func Critical(alerts []domain.Alert) []domain.Alert {
	result := []domain.Alert{}
	for _, alert := range alerts {
		if alert.Severity == "critical" && alert.State != "resolved" {
			result = append(result, alert)
		}
	}
	return result
}
