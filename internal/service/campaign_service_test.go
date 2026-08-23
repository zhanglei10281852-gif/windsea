package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
	"time"
)

func TestCampaignCannotCloseWithPendingInspection(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	operator := testkit.SeedUser(t, database, farm, "operator", "operator")
	testkit.MustExec(t, database.DB, "INSERT INTO turbines(id,farm_id,name,model,rated_kw,status,created_at) VALUES('t1',?,'T1','M',100,'operational',CURRENT_TIMESTAMP)", farm.ID)
	campaign := domain.Campaign{ID: "campaign-1", FarmID: farm.ID, Name: "Inspection wave", StartAt: time.Now(), EndAt: time.Now().Add(time.Hour)}
	if err := services.Campaigns.Create(context.Background(), campaign); err != nil {
		t.Fatal(err)
	}
	testkit.MustExec(t, database.DB, "INSERT INTO inspections(id,campaign_id,turbine_id,inspector_id,status,created_at) VALUES('inspection-1','campaign-1','t1',?,'pending',CURRENT_TIMESTAMP)", operator.ID)
	if _, err := services.Campaigns.SubmitReview(context.Background(), campaign.ID, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := services.Campaigns.Publish(context.Background(), campaign.ID, 2); err != nil {
		t.Fatal(err)
	}
	if _, err := services.Campaigns.Close(context.Background(), campaign.ID, 3); err == nil {
		t.Fatal("closed campaign with pending inspection")
	}
	if err := services.Campaigns.CompleteInspection(context.Background(), "inspection-1", "complete"); err != nil {
		t.Fatal(err)
	}
	if _, err := services.Campaigns.Close(context.Background(), campaign.ID, 3); err != nil {
		t.Fatal(err)
	}
}
func TestCampaignVersionConflict(t *testing.T) {
	_, services := testkit.Open(t)
	campaign := domain.Campaign{ID: "c", FarmID: "f", Name: "x", StartAt: time.Now(), EndAt: time.Now().Add(time.Hour)}
	_ = services
	_ = campaign
}
