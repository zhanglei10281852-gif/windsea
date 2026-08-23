package reporting

import (
	"encoding/csv"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"io"
	"sort"
	"strconv"
	"time"
)

type WorkOrderRow struct {
	ID, Title, State string
	Priority         int
	CreatedAt        time.Time
}

func Rows(orders []domain.WorkOrder) []WorkOrderRow {
	rows := make([]WorkOrderRow, 0, len(orders))
	for _, order := range orders {
		rows = append(rows, WorkOrderRow{ID: order.ID, Title: order.Title, State: order.State, Priority: order.Priority, CreatedAt: order.CreatedAt})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].CreatedAt.Before(rows[j].CreatedAt) })
	return rows
}
func WriteCSV(writer io.Writer, rows []WorkOrderRow) error {
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write([]string{"id", "title", "state", "priority", "created_at"}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := csvWriter.Write([]string{row.ID, row.Title, row.State, strconv.Itoa(row.Priority), row.CreatedAt.UTC().Format(time.RFC3339)}); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return csvWriter.Error()
}
func Summary(orders []domain.WorkOrder) map[string]int {
	summary := map[string]int{}
	for _, order := range orders {
		summary[order.State]++
	}
	return summary
}
func EnsureComplete(orders []domain.WorkOrder) error {
	for _, order := range orders {
		if order.ID == "" || order.Title == "" {
			return fmt.Errorf("%w: report row", domain.ErrValidation)
		}
	}
	return nil
}
