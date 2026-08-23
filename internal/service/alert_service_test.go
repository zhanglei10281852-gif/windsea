package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
	"time"
)

func TestAlertLifecycle(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	testkit.MustExec(t, database.DB, "INSERT INTO turbines(id,farm_id,name,model,rated_kw,status,created_at) VALUES('ta',?,'T','M',100,'operational',CURRENT_TIMESTAMP)", farm.ID)
	occurred := time.Now()
	if err := services.Alerts.Ingest(context.Background(), domain.Alert{ID: "a1", FarmID: farm.ID, TurbineID: "ta", Code: "temperature_high", Severity: "critical", Message: "bearing temperature", OccurredAt: &occurred}); err != nil {
		t.Fatal(err)
	}
	if err := services.Alerts.Resolve(context.Background(), "a1"); err == nil {
		t.Fatal("resolved unacknowledged alert")
	}
	if err := services.Alerts.Acknowledge(context.Background(), "a1"); err != nil {
		t.Fatal(err)
	}
	if err := services.Alerts.Resolve(context.Background(), "a1"); err != nil {
		t.Fatal(err)
	}
	items, err := services.Alerts.Open(context.Background(), farm.ID)
	if err != nil || len(items) != 0 {
		t.Fatalf("items=%+v err=%v", items, err)
	}
}
func TestAlertDeduplication(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	testkit.MustExec(t, database.DB, "INSERT INTO turbines(id,farm_id,name,model,rated_kw,status,created_at) VALUES('tb',?,'T','M',100,'operational',CURRENT_TIMESTAMP)", farm.ID)
	when := time.Now()
	alert := domain.Alert{ID: "a2", FarmID: farm.ID, TurbineID: "tb", Code: "vibration", Severity: "warning", Message: "vibration", OccurredAt: &when}
	if err := services.Alerts.Ingest(context.Background(), alert); err != nil {
		t.Fatal(err)
	}
	alert.ID = "a3"
	if err := services.Alerts.Ingest(context.Background(), alert); err == nil {
		t.Fatal("duplicate alert accepted")
	}
}
