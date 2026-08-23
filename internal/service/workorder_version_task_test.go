package service_test

import (
	"context"
	"testing"

	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
)

func TestWorkOrderAssignmentRejectsStaleVersion(t *testing.T) {
	db, sv := testkit.Open(t)
	farm := testkit.SeedFarm(t, db)
	testkit.SeedUser(t, db, farm, "operator-a", "operator")
	testkit.SeedUser(t, db, farm, "operator-b", "operator")
	if err := sv.WorkOrders.Create(context.Background(), domain.WorkOrder{ID: "wo-version", FarmID: farm.ID, Title: "inspect", Priority: 3}); err != nil {
		t.Fatal(err)
	}
	if err := sv.WorkOrders.Assign(context.Background(), "wo-version", "operator-a", 1); err != nil {
		t.Fatal(err)
	}
	if err := sv.WorkOrders.Assign(context.Background(), "wo-version", "operator-b", 1); err == nil {
		t.Fatal("stale assignment unexpectedly succeeded")
	}
	order, err := repository.New(db.DB).GetWorkOrder(context.Background(), "wo-version")
	if err != nil {
		t.Fatal(err)
	}
	if order.AssigneeID != "operator-a" {
		t.Fatalf("assignee overwritten: %s", order.AssigneeID)
	}
}
