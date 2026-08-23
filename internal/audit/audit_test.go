package audit_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/audit"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"testing"
	"time"
)

func TestEventBuilderRoundTrip(t *testing.T) {
	builder := audit.EventBuilder{Now: func() time.Time { return time.Unix(1, 0) }}
	event, err := builder.Build("e", "f", "u", "work_order", "wo", "completed", "ok", "req", map[string]any{"priority": 3})
	if err != nil {
		t.Fatal(err)
	}
	metadata, err := audit.DecodeMetadata(event)
	if err != nil || metadata["priority"] != float64(3) {
		t.Fatalf("metadata=%v err=%v", metadata, err)
	}
}
func TestEventBuilderRequiresIdentity(t *testing.T) {
	builder := audit.EventBuilder{Now: time.Now}
	if _, err := builder.Build("", "f", "u", "x", "id", "a", "ok", "r", nil); err == nil {
		t.Fatal("missing event id accepted")
	}
	_ = domain.ErrValidation
}
