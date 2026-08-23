package testkit

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/db"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/observability"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
	"github.com/zhanglei10281852-gif/windsea/internal/service"
	"log/slog"
	"testing"
	"time"
)

func Open(t *testing.T) (*db.DB, *service.Registry) {
	t.Helper()
	database, err := db.Open(context.Background(), fmt.Sprintf("file:%s?mode=memory&cache=shared&_pragma=foreign_keys(1)", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })
	logger := slog.New(slog.NewTextHandler(testWriter{t}, nil))
	fixed := func() time.Time { return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC) }
	return database, service.NewRegistry(database, logger, fixed)
}

type testWriter struct{ t *testing.T }

func (w testWriter) Write(bytes []byte) (int, error) { w.t.Log(string(bytes)); return len(bytes), nil }
func SeedFarm(t *testing.T, database *db.DB) domain.Farm {
	t.Helper()
	farm := domain.Farm{ID: "farm-1", Name: "North Sea Alpha", Timezone: "UTC", TurbineCount: 2, CreatedAt: time.Now().UTC()}
	if err := repository.New(database.DB).CreateFarm(context.Background(), farm); err != nil {
		t.Fatal(err)
	}
	return farm
}
func SeedUser(t *testing.T, database *db.DB, farm domain.Farm, id, role string) domain.User {
	t.Helper()
	user := domain.User{ID: id, FarmID: farm.ID, Email: id + "@windsea.test", Name: id, Role: role, Active: true, CreatedAt: time.Now().UTC()}
	if err := repository.New(database.DB).CreateUser(context.Background(), user); err != nil {
		t.Fatal(err)
	}
	return user
}
func MustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatal(err)
	}
}

var _ = observability.RequestID
