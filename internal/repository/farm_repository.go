package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zhanglei10281852-gif/windsea/internal/domain"
)

func (r *Repository) CreateFarm(ctx context.Context, farm domain.Farm) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO farms(id,name,timezone,turbine_count,created_at) VALUES(?,?,?,?,?)", farm.ID, farm.Name, farm.Timezone, farm.TurbineCount, text(farm.CreatedAt))
	if err != nil {
		return fmt.Errorf("create farm: %w", err)
	}
	return nil
}
func (r *Repository) GetFarm(ctx context.Context, id string) (domain.Farm, error) {
	var farm domain.Farm
	var created string
	err := r.DB.QueryRowContext(ctx, "SELECT id,name,timezone,turbine_count,created_at FROM farms WHERE id=?", id).Scan(&farm.ID, &farm.Name, &farm.Timezone, &farm.TurbineCount, &created)
	if err != nil {
		if isNoRows(err) {
			return farm, domain.ErrNotFound
		}
		return farm, fmt.Errorf("get farm: %w", err)
	}
	farm.CreatedAt, _ = parse(created)
	return farm, nil
}
func (r *Repository) CreateUser(ctx context.Context, user domain.User) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO users(id,farm_id,email,name,role,active,created_at) VALUES(?,?,?,?,?,?,?)", user.ID, user.FarmID, user.Email, user.Name, user.Role, user.Active, text(user.CreatedAt))
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}
func (r *Repository) GetUser(ctx context.Context, id string) (domain.User, error) {
	var u domain.User
	var active int
	var created string
	err := r.DB.QueryRowContext(ctx, "SELECT id,farm_id,email,name,role,active,created_at FROM users WHERE id=?", id).Scan(&u.ID, &u.FarmID, &u.Email, &u.Name, &u.Role, &active, &created)
	if err != nil {
		if isNoRows(err) {
			return u, domain.ErrNotFound
		}
		return u, fmt.Errorf("get user: %w", err)
	}
	u.Active = active == 1
	u.CreatedAt, _ = parse(created)
	return u, nil
}
func (r *Repository) CreateTurbine(ctx context.Context, turbine domain.Turbine) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO turbines(id,farm_id,name,model,rated_kw,status,version,created_at) VALUES(?,?,?,?,?,?,?,?)", turbine.ID, turbine.FarmID, turbine.Name, turbine.Model, turbine.RatedKW, turbine.Status, turbine.Version, text(turbine.CreatedAt))
	if err != nil {
		return fmt.Errorf("create turbine: %w", err)
	}
	return nil
}
func (r *Repository) ListTurbines(ctx context.Context, farmID string) ([]domain.Turbine, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id,farm_id,name,model,rated_kw,status,version,created_at FROM turbines WHERE farm_id=? ORDER BY name", farmID)
	if err != nil {
		return nil, fmt.Errorf("list turbines: %w", err)
	}
	defer rows.Close()
	result := []domain.Turbine{}
	for rows.Next() {
		var t domain.Turbine
		var created string
		if err := rows.Scan(&t.ID, &t.FarmID, &t.Name, &t.Model, &t.RatedKW, &t.Status, &t.Version, &created); err != nil {
			return nil, fmt.Errorf("scan turbine: %w", err)
		}
		t.CreatedAt, _ = parse(created)
		result = append(result, t)
	}
	return result, rows.Err()
}
func isNoRows(err error) bool { return err != nil && err.Error() == "sql: no rows in result set" }

var _ = time.UTC

func (r *Repository) CreateSession(ctx context.Context, session domain.Session) error {
	if session.ExpiresAt == nil {
		return domain.ErrValidation
	}
	_, err := r.DB.ExecContext(ctx, "INSERT INTO sessions(id,user_id,token_hash,expires_at) VALUES(?,?,?,?)", session.ID, session.UserID, session.TokenHash, text(*session.ExpiresAt))
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) FindSession(ctx context.Context, tokenHash string) (domain.Session, error) {
	var s domain.Session
	var expires string
	var revoked string
	err := r.DB.QueryRowContext(ctx, "SELECT id,user_id,token_hash,expires_at,revoked_at FROM sessions WHERE token_hash=?", tokenHash).Scan(&s.ID, &s.UserID, &s.TokenHash, &expires, &revoked)
	if err != nil {
		if isNoRows(err) {
			return s, domain.ErrNotFound
		}
		return s, err
	}
	parsed, _ := parse(expires)
	s.ExpiresAt = &parsed
	s.RevokedAt, _ = optionalString(revoked)
	return s, nil
}
func (r *Repository) RevokeSession(ctx context.Context, tokenHash, at string) error {
	result, err := r.DB.ExecContext(ctx, "UPDATE sessions SET revoked_at=? WHERE token_hash=? AND revoked_at IS NULL", at, tokenHash)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}
