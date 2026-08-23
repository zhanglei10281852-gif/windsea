package authorization_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/auth"
	"github.com/zhanglei10281852-gif/windsea/internal/authorization"
	"testing"
)

func TestFarmScopedPolicy(t *testing.T) {
	ctx := auth.WithPrincipal(context.Background(), auth.Principal{FarmID: "farm-1", Role: auth.RoleSupervisor})
	policy := authorization.Policy{Action: "publish", Roles: []string{auth.RoleSupervisor}, FarmScoped: true}
	if err := authorization.Check(ctx, policy, "farm-1"); err != nil {
		t.Fatal(err)
	}
	if err := authorization.Check(ctx, policy, "farm-2"); err == nil {
		t.Fatal("cross farm access")
	}
	policies := []authorization.Policy{{Action: "read", Roles: []string{auth.RoleViewer}}, {Action: "publish", Roles: []string{auth.RoleSupervisor}}}
	if err := authorization.Any(ctx, policies, "farm-1"); err != nil {
		t.Fatal(err)
	}
}
