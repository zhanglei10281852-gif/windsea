package auth_test

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/windsea/internal/auth"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"testing"
	"time"
)

func TestPrincipalContext(t *testing.T) {
	ctx := auth.WithPrincipal(context.Background(), auth.Principal{UserID: "u", FarmID: "f", Role: auth.RoleSupervisor})
	principal, ok := auth.PrincipalFrom(ctx)
	if !ok || principal.Role != auth.RoleSupervisor {
		t.Fatalf("principal=%+v", principal)
	}
	if err := auth.RequireRole(ctx, auth.RoleSupervisor); err != nil {
		t.Fatal(err)
	}
	if err := auth.RequireRole(ctx, auth.RoleViewer); err == nil {
		t.Fatal("expected forbidden")
	}
}
func TestSessionExpiry(t *testing.T) {
	now := time.Now()
	if err := auth.ValidateSession(now, now.Add(time.Hour), nil); err != nil {
		t.Fatal(err)
	}
	if err := auth.ValidateSession(now, now.Add(-time.Minute), nil); err == nil {
		t.Fatal("expired session accepted")
	}
	revoked := now.Add(-time.Minute)
	if err := auth.ValidateSession(now, now.Add(time.Hour), &revoked); err == nil {
		t.Fatal("revoked session accepted")
	}
}
func TestAuthorization(t *testing.T) {
	if err := auth.Authorize(auth.RoleSupervisor, "publish_campaign"); err != nil {
		t.Fatal(err)
	}
	if err := auth.Authorize(auth.RoleViewer, "publish_campaign"); err == nil {
		t.Fatal("viewer published")
	}
	if !errors.Is(domain.ErrForbidden, domain.ErrForbidden) {
		t.Fatal("sentinel")
	}
}
