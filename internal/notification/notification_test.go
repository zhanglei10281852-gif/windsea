package notification_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/notification"
	"testing"
	"time"
)

func TestDeliverRetries(t *testing.T) {
	sender := &notification.MemorySender{Failures: 2}
	message, err := notification.BuildAlertMessage(domain.Alert{ID: "a", FarmID: "f", Code: "temp", Message: "hot"}, "ops@example.com", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := notification.Deliver(context.Background(), sender, message, 3); err != nil {
		t.Fatal(err)
	}
	if len(sender.Messages) != 1 {
		t.Fatal(sender.Messages)
	}
}
func TestMessageValidation(t *testing.T) {
	if _, err := notification.BuildAlertMessage(domain.Alert{}, "", time.Now()); err == nil {
		t.Fatal("invalid message accepted")
	}
}
