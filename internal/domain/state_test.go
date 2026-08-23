package domain_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"testing"
)

func TestCampaignLifecycle(t *testing.T) {
	cases := []struct {
		from, to string
		ok       bool
	}{{"draft", "review", true}, {"review", "published", true}, {"published", "closed", true}, {"draft", "published", false}, {"closed", "draft", false}}
	for _, tc := range cases {
		if got := domain.CampaignState(tc.from).CanMove(domain.CampaignState(tc.to)); got != tc.ok {
			t.Errorf("%s to %s = %v", tc.from, tc.to, got)
		}
	}
}
func TestWorkOrderLifecycle(t *testing.T) {
	for _, tc := range []struct {
		from, to string
		ok       bool
	}{{"queued", "assigned", true}, {"assigned", "in_progress", true}, {"in_progress", "blocked", true}, {"blocked", "in_progress", true}, {"in_progress", "completed", true}, {"queued", "completed", false}, {"completed", "blocked", false}} {
		if got := domain.WorkOrderState(tc.from).CanMove(domain.WorkOrderState(tc.to)); got != tc.ok {
			t.Errorf("%s to %s = %v", tc.from, tc.to, got)
		}
	}
}
func TestReservationLifecycle(t *testing.T) {
	if !domain.ReservationHeld.CanMove(domain.ReservationConsumed) {
		t.Fatal("held should consume")
	}
	if !domain.ReservationHeld.CanMove(domain.ReservationReleased) {
		t.Fatal("held should release")
	}
	if domain.ReservationConsumed.CanMove(domain.ReservationHeld) {
		t.Fatal("consumed cannot reopen")
	}
}
func TestAlertLifecycle(t *testing.T) {
	if !domain.AlertOpen.CanMove(domain.AlertAcknowledged) {
		t.Fatal("open should acknowledge")
	}
	if !domain.AlertAcknowledged.CanMove(domain.AlertResolved) {
		t.Fatal("ack should resolve")
	}
	if domain.AlertOpen.CanMove(domain.AlertResolved) {
		t.Fatal("open cannot resolve directly")
	}
}
func TestValidateTransition(t *testing.T) {
	if err := domain.ValidateTransition("campaign", "draft", "review"); err != nil {
		t.Fatal(err)
	}
	if err := domain.ValidateTransition("campaign", "draft", "closed"); err == nil {
		t.Fatal("expected invalid transition")
	}
	if err := domain.ValidateTransition("unknown", "a", "b"); err == nil {
		t.Fatal("expected unknown entity error")
	}
}
func TestApplyTransitions(t *testing.T) {
	state, err := domain.ApplyTransitions("draft", []domain.Transition{{Entity: "campaign", From: "draft", To: "review"}, {Entity: "campaign", From: "review", To: "published"}})
	if err != nil || state != "published" {
		t.Fatalf("state=%s err=%v", state, err)
	}
	if _, err := domain.ApplyTransitions("draft", []domain.Transition{{Entity: "campaign", From: "published", To: "closed"}}); err == nil {
		t.Fatal("expected conflict")
	}
}
