package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
)

func TestWorkOrderLifecycleAndPaging(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	if err := services.WorkOrders.Create(context.Background(), domain.WorkOrder{ID: "wo-1", FarmID: farm.ID, Title: "Replace yaw motor", Priority: 5}); err != nil {
		t.Fatal(err)
	}
	user := testkit.SeedUser(t, database, farm, "tech", "operator")
	if err := services.WorkOrders.Assign(context.Background(), "wo-1", user.ID, 1); err != nil {
		t.Fatal(err)
	}
	if err := services.WorkOrders.Start(context.Background(), "wo-1", 2); err != nil {
		t.Fatal(err)
	}
	if err := services.WorkOrders.Complete(context.Background(), "wo-1", 3); err != nil {
		t.Fatal(err)
	}
	page, err := services.WorkOrders.List(context.Background(), farm.ID, "completed", 10, 0)
	if err != nil || page.Total != 1 || len(page.Items) != 1 {
		t.Fatalf("page=%+v err=%v", page, err)
	}
}
func TestWorkOrderInvalidTransition(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	if err := services.WorkOrders.Create(context.Background(), domain.WorkOrder{ID: "wo-2", FarmID: farm.ID, Title: "Blade inspection", Priority: 3}); err != nil {
		t.Fatal(err)
	}
	if err := services.WorkOrders.Complete(context.Background(), "wo-2", 1); err == nil {
		t.Fatal("queued order completed")
	}
}
