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

func TestPublishCampaignWithAuditStaysInReviewOnAuditFailure(t *testing.T) {
	database, services := testkit.Open(t)
	farm := testkit.SeedFarm(t, database)

	campaign := domain.Campaign{ID: "campaign-pub", FarmID: farm.ID, Name: "Wave", StartAt: time.Now(), EndAt: time.Now().Add(time.Hour)}
	if err := services.Campaigns.Create(context.Background(), campaign); err != nil {
		t.Fatal(err)
	}
	if _, err := services.Campaigns.SubmitReview(context.Background(), campaign.ID, 1); err != nil {
		t.Fatal(err)
	}

	// Pre-seed a colliding audit_event id so the publish transaction's audit insert fails.
	repo := repository.New(database.DB)
	collision := domain.AuditEvent{ID: "audit-collision", FarmID: farm.ID, ActorID: "operator", ObjectType: "campaign", ObjectID: campaign.ID, Action: "publish", Result: "ok", RequestID: "req", Metadata: "{}", CreatedAt: time.Now()}
	if err := repo.AppendAudit(context.Background(), collision); err != nil {
		t.Fatal(err)
	}

	event := domain.AuditEvent{ID: "audit-collision", FarmID: farm.ID, ActorID: "operator", ObjectType: "campaign", ObjectID: campaign.ID, Action: "publish", Result: "ok", RequestID: "req", Metadata: "{}", CreatedAt: time.Now()}
	if err := service.PublishCampaignWithAudit(context.Background(), database.DB, campaign.ID, 2, event, time.Now()); err == nil {
		t.Fatal("expected audit persistence failure, got nil")
	}

	after, err := repo.GetCampaign(context.Background(), campaign.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.State != string(domain.CampaignReview) {
		t.Fatalf("campaign state=%q, want %q (should stay in review when audit fails)", after.State, domain.CampaignReview)
	}
}
