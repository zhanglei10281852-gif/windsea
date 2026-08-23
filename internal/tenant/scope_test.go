package tenant_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/tenant"
	"testing"
)

func TestFarmScope(t *testing.T) {
	ctx := tenant.WithFarm(context.Background(), "farm-1")
	if id, err := tenant.FarmID(ctx); err != nil || id != "farm-1" {
		t.Fatalf("id=%s err=%v", id, err)
	}
	if err := tenant.Ensure(ctx, "farm-2"); err == nil {
		t.Fatal("scope mismatch accepted")
	}
	if _, err := tenant.FarmID(context.Background()); err == nil {
		t.Fatal("missing scope accepted")
	}
}
