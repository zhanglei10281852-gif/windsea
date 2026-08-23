package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
	"time"
)

func TestCampaignCloseRejectsInProgressInspection(t *testing.T) {
	db, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, db)
	user := testkit.SeedUser(t, db, farm, "inspector-progress", "inspector")
	turbine := domain.Turbine{ID: "turb-progress", FarmID: farm.ID, Name: "Progress", Model: "M", RatedKW: 8, Status: "ready", Version: 1, CreatedAt: time.Now()}
	if err := repository.New(db.DB).CreateTurbine(context.Background(), turbine); err != nil {
		t.Fatal(err)
	}
	campaign := domain.Campaign{ID: "campaign-progress", FarmID: farm.ID, Name: "Progress close", StartAt: time.Now(), EndAt: time.Now().Add(time.Hour)}
	if err := services.Campaigns.Create(context.Background(), campaign); err != nil {
		t.Fatal(err)
	}
	if _, err := services.Campaigns.SubmitReview(context.Background(), campaign.ID, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := services.Campaigns.Publish(context.Background(), campaign.ID, 2); err != nil {
		t.Fatal(err)
	}
	if err := repository.New(db.DB).AddInspection(context.Background(), domain.Inspection{ID: "inspection-progress", CampaignID: campaign.ID, TurbineID: turbine.ID, InspectorID: user.ID, Status: "in_progress", CreatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	got, err := services.Campaigns.Close(context.Background(), campaign.ID, 3)
	if err == nil {
		t.Fatalf("close unexpectedly succeeded with in-progress inspection: %s", got.State)
	}
	if got.State != string(domain.CampaignPublished) {
		t.Fatalf("campaign state=%s", got.State)
	}
}
