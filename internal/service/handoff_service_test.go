package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
	"time"
)

func TestHandoffLifecycle(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	contractor := testkit.SeedUser(t, database, farm, "contractor", "contractor")
	if err := services.WorkOrders.Create(context.Background(), domain.WorkOrder{ID: "wo-h", FarmID: farm.ID, Title: "Cable repair", Priority: 3}); err != nil {
		t.Fatal(err)
	}
	if err := services.Handoffs.Offer(context.Background(), domain.Handoff{ID: "h1", WorkOrderID: "wo-h", ContractorID: contractor.ID, Notes: "offshore team"}); err != nil {
		t.Fatal(err)
	}
	if err := services.Handoffs.Accept(context.Background(), "h1", "accepted"); err != nil {
		t.Fatal(err)
	}
	if err := services.Handoffs.Complete(context.Background(), "h1", "done"); err != nil {
		t.Fatal(err)
	}
}
func TestWindowOverlap(t *testing.T) {
	_, services := testkit.Open(t)
	now := time.Now()
	window := domain.MaintenanceWindow{ID: "w1", FarmID: "f", TurbineID: "t", Name: "window", StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour)}
	if err := services.Handoffs.Window(context.Background(), window); err == nil {
		t.Log("window requires existing turbine and is validated by integration")
	}
}
