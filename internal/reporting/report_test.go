package reporting_test

import (
	"bytes"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/reporting"
	"strings"
	"testing"
	"time"
)

func TestReportRowsAndCSV(t *testing.T) {
	now := time.Now()
	orders := []domain.WorkOrder{{ID: "2", Title: "B", State: "queued", Priority: 2, CreatedAt: now.Add(time.Minute)}, {ID: "1", Title: "A", State: "completed", Priority: 5, CreatedAt: now}}
	rows := reporting.Rows(orders)
	if rows[0].ID != "1" {
		t.Fatalf("rows=%+v", rows)
	}
	var buffer bytes.Buffer
	if err := reporting.WriteCSV(&buffer, rows); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buffer.String(), "id,title,state") {
		t.Fatalf("csv=%s", buffer.String())
	}
	if summary := reporting.Summary(orders); summary["completed"] != 1 {
		t.Fatal(summary)
	}
}
func TestCriticalAlerts(t *testing.T) {
	alerts := []domain.Alert{{ID: "a", Severity: "critical", State: "open"}, {ID: "b", Severity: "warning", State: "open"}, {ID: "c", Severity: "critical", State: "resolved"}}
	if len(reporting.Critical(alerts)) != 1 {
		t.Fatal("critical filter")
	}
	rows := reporting.AlertRows(alerts)
	if rows[0].ID != "a" {
		t.Fatal(rows)
	}
}
