package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
)

func TestTelemetryBatchCompletes(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	if err := services.Telemetry.Receive(context.Background(), domain.TelemetryBatch{ID: "batch-1", FarmID: farm.ID, Source: "scada", Samples: 4}); err != nil {
		t.Fatal(err)
	}
	if err := services.Telemetry.Complete(context.Background(), "batch-1", 3, 1); err != nil {
		t.Fatal(err)
	}
	batch, err := services.Telemetry.Get(context.Background(), "batch-1")
	if err != nil || batch.State != "completed" || batch.Accepted != 3 || batch.Rejected != 1 {
		t.Fatalf("batch=%+v err=%v", batch, err)
	}
}
func TestTelemetryRejectsInvalidCounts(t *testing.T) {
	_, services := testkit.Open(t)
	if err := services.Telemetry.Receive(context.Background(), domain.TelemetryBatch{ID: "bad", FarmID: "f", Source: "scada", Samples: 0}); err == nil {
		t.Fatal("zero samples accepted")
	}
}
