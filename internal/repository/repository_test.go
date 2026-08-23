package repository_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
	"time"
)

func TestFarmAndTurbinePersistence(t *testing.T) {
	database, _ := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	repo := repository.New(database.DB)
	turbine := domain.Turbine{ID: "t1", FarmID: farm.ID, Name: "Turbine 1", Model: "WT-8", RatedKW: 8000, Status: "operational", Version: 1, CreatedAt: time.Now().UTC()}
	if err := repo.CreateTurbine(context.Background(), turbine); err != nil {
		t.Fatal(err)
	}
	items, err := repo.ListTurbines(context.Background(), farm.ID)
	if err != nil || len(items) != 1 {
		t.Fatalf("items=%+v err=%v", items, err)
	}
}
func TestCampaignInspectionPersistence(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	user := testkit.SeedUser(t, database, farm, "inspector", "operator")
	campaign := domain.Campaign{ID: "c1", FarmID: farm.ID, Name: "Spring survey", StartAt: time.Now(), EndAt: time.Now().Add(time.Hour)}
	if err := services.Campaigns.Create(context.Background(), campaign); err != nil {
		t.Fatal(err)
	}
	repo := repository.New(database.DB)
	inspection := domain.Inspection{ID: "i1", CampaignID: "c1", TurbineID: "t1", InspectorID: user.ID, Status: "pending", CreatedAt: time.Now()}
	testkit.MustExec(t, database.DB, "INSERT INTO turbines(id,farm_id,name,model,rated_kw,status,created_at) VALUES('t1',?,'T1','M',100,'operational',CURRENT_TIMESTAMP)", farm.ID)
	if err := repo.AddInspection(context.Background(), inspection); err != nil {
		t.Fatal(err)
	}
	if count, err := repo.CountIncompleteInspections(context.Background(), campaign.ID); err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}
	if err := services.Campaigns.CompleteInspection(context.Background(), inspection.ID, "all good"); err != nil {
		t.Fatal(err)
	}
	if count, err := repo.CountIncompleteInspections(context.Background(), campaign.ID); err != nil || count != 0 {
		t.Fatalf("count=%d err=%v", count, err)
	}
}
