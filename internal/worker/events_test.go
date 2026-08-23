package worker_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/worker"
	"testing"
	"time"
)

func TestEventBusPublishAndCancel(t *testing.T) {
	bus := worker.NewEventBus()
	events, cancel := bus.Subscribe("alert")
	defer cancel()
	if err := bus.Publish(context.Background(), worker.Event{ID: "a", Kind: "alert", Payload: "open"}); err != nil {
		t.Fatal(err)
	}
	select {
	case event := <-events:
		if event.ID != "a" {
			t.Fatal(event)
		}
	case <-time.After(time.Second):
		t.Fatal("event missing")
	}
	bus.Close()
}
