package service_test

import (
	"context"
	"github.com/zhanglei10281852-gif/windsea/internal/service"
	"github.com/zhanglei10281852-gif/windsea/internal/testkit"
	"testing"
	"time"
)

func TestHealthReportsReadyDatabase(t *testing.T) {
	database, _ := testkit.Open(t)
	health := service.CheckHealth(context.Background(), database.DB, time.Now())
	if health.Status != "ready" || health.Database != "up" {
		t.Fatalf("health=%+v", health)
	}
}
