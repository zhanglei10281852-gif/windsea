package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
	"time"
)

func TestReservationNeverOversellsPart(t *testing.T) {
	db, sv := testkit.Open(t)
	farm := testkit.SeedFarm(t, db)
	u := testkit.SeedUser(t, db, farm, "operator-race", "operator")
	for _, id := range []string{"wo-race-1", "wo-race-2"} {
		if err := sv.WorkOrders.Create(context.Background(), domain.WorkOrder{ID: id, FarmID: farm.ID, Title: id, Priority: 2}); err != nil {
			t.Fatal(err)
		}
	}
	part := domain.Part{ID: "part-last", FarmID: farm.ID, SKU: "GBX", Description: "gearbox", OnHand: 1, CreatedAt: time.Now()}
	if err := sv.Inventory.CreatePart(context.Background(), part); err != nil {
		t.Fatal(err)
	}
	if err := sv.Inventory.Hold(context.Background(), "res-race-1", part.ID, "wo-race-1", u.ID, 1); err != nil {
		t.Fatal(err)
	}
	if err := sv.Inventory.Hold(context.Background(), "res-race-2", part.ID, "wo-race-2", u.ID, 1); err == nil {
		p, _ := repository.New(db.DB).GetPart(context.Background(), part.ID)
		t.Fatalf("second hold succeeded reserved=%d", p.Reserved)
	}
}
