package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/auth"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"log/slog"
	"time"
)

type AuthService struct {
	repo *repository.Repository
	log  *slog.Logger
	now  func() time.Time
}

func (s *AuthService) Issue(ctx context.Context, user domain.User, ttl time.Duration) (string, time.Time, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", time.Time{}, err
	}
	token := hex.EncodeToString(tokenBytes)
	expires := s.now().Add(ttl)
	if err := s.repo.CreateSession(ctx, domain.Session{ID: token[:16], UserID: user.ID, TokenHash: auth.HashToken(token), ExpiresAt: &expires}); err != nil {
		return "", time.Time{}, fmt.Errorf("issue session: %w", err)
	}
	return token, expires, nil
}
func (s *AuthService) Authenticate(ctx context.Context, token string) (auth.Principal, error) {
	session, err := s.repo.FindSession(ctx, auth.HashToken(token))
	if err != nil {
		return auth.Principal{}, domain.ErrForbidden
	}
	if session.ExpiresAt == nil {
		return auth.Principal{}, domain.ErrForbidden
	}
	if err := auth.ValidateSession(s.now(), *session.ExpiresAt, session.RevokedAt); err != nil {
		return auth.Principal{}, err
	}
	user, err := s.repo.GetUser(ctx, session.UserID)
	if err != nil || !user.Active {
		return auth.Principal{}, domain.ErrForbidden
	}
	return auth.Principal{UserID: user.ID, FarmID: user.FarmID, Role: user.Role}, nil
}
func (s *AuthService) Revoke(ctx context.Context, token string) error {
	return s.repo.RevokeSession(ctx, auth.HashToken(token), s.now().Format(time.RFC3339Nano))
}
