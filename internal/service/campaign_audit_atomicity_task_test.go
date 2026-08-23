package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"github.com/zhanglei10281852-gif/windsea/internal/service"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
)

func TestCampaignPublishKeepsReviewWhenAuditFails(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)
	campaign := domain.Campaign{ID: "audit-campaign", FarmID: farm.ID, Name: "Audit wave", StartAt: time.Now(), EndAt: time.Now().Add(time.Hour)}
	if err := services.Campaigns.Create(context.Background(), campaign); err != nil {
		t.Fatal(err)
	}
	if _, err := services.Campaigns.SubmitReview(context.Background(), campaign.ID, 1); err != nil {
		t.Fatal(err)
	}
	event := domain.AuditEvent{ID: "bad-audit", FarmID: "missing-farm", ObjectType: "campaign", ObjectID: campaign.ID, Action: "publish", Result: "ok", RequestID: "req", Metadata: "{}"}
	if err := service.PublishCampaignWithAudit(context.Background(), database.DB, campaign.ID, 2, event, time.Now()); err == nil {
		t.Fatal("expected audit failure")
	}
	stored, err := repository.New(database.DB).GetCampaign(context.Background(), campaign.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.State != string(domain.CampaignReview) {
		t.Fatalf("campaign state=%s", stored.State)
	}
	valid := event
	valid.ID = "good-audit"
	valid.FarmID = farm.ID
	if err := service.PublishCampaignWithAudit(context.Background(), database.DB, campaign.ID, 2, valid, time.Now()); err != nil {
		t.Fatal(err)
	}
	stored, err = repository.New(database.DB).GetCampaign(context.Background(), campaign.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.State != string(domain.CampaignPublished) {
		t.Fatalf("successful state=%s", stored.State)
	}
}
