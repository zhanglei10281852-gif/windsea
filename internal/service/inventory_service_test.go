package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
)

func TestReservationConsumesCapacityAtomically(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	user := testkit.SeedUser(t, database, farm, "operator", "operator")
	if err := services.WorkOrders.Create(context.Background(), domain.WorkOrder{ID: "wo", FarmID: farm.ID, Title: "Gearbox", Priority: 4}); err != nil {
		t.Fatal(err)
	}
	if err := services.Inventory.CreatePart(context.Background(), domain.Part{ID: "p", FarmID: farm.ID, SKU: "GEAR", Description: "gearbox", OnHand: 2}); err != nil {
		t.Fatal(err)
	}
	if err := services.Inventory.Hold(context.Background(), "r1", "p", "wo", user.ID, 2); err != nil {
		t.Fatal(err)
	}
	if err := services.Inventory.Hold(context.Background(), "r2", "p", "wo", user.ID, 1); err == nil {
		t.Fatal("oversold part")
	}
	if err := services.Inventory.Consume(context.Background(), "r1"); err != nil {
		t.Fatal(err)
	}
}
func TestReservationReleaseRestoresAvailability(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	user := testkit.SeedUser(t, database, farm, "operator2", "operator")
	if err := services.WorkOrders.Create(context.Background(), domain.WorkOrder{ID: "wo2", FarmID: farm.ID, Title: "Pitch motor", Priority: 2}); err != nil {
		t.Fatal(err)
	}
	if err := services.WorkOrders.Create(context.Background(), domain.WorkOrder{ID: "wo3", FarmID: farm.ID, Title: "Pitch motor spare", Priority: 2}); err != nil {
		t.Fatal(err)
	}
	if err := services.Inventory.CreatePart(context.Background(), domain.Part{ID: "p2", FarmID: farm.ID, SKU: "PITCH", Description: "motor", OnHand: 1}); err != nil {
		t.Fatal(err)
	}
	if err := services.Inventory.Hold(context.Background(), "r3", "p2", "wo2", user.ID, 1); err != nil {
		t.Fatal(err)
	}
	if err := services.Inventory.Release(context.Background(), "r3"); err != nil {
		t.Fatal(err)
	}
	if err := services.Inventory.Hold(context.Background(), "r4", "p2", "wo3", user.ID, 1); err != nil {
		t.Fatal(err)
	}
}
