package stream_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/stream"
	"testing"
	"time"
)

func TestBroadcaster(t *testing.T) {
	b := stream.New[string]()
	id, events, err := b.Listen(1)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Send(context.Background(), "alert"); err != nil {
		t.Fatal(err)
	}
	select {
	case value := <-events:
		if value != "alert" {
			t.Fatal(value)
		}
	case <-time.After(time.Second):
		t.Fatal("event missing")
	}
	b.Remove(id)
	b.Close()
}
