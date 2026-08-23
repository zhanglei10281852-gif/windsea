package service

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/zhanglei10281852-gif/windsea/internal/db"
	"github.com/zhanglei10281852-gif/windsea/internal/repository"
)

type Registry struct {
	DB         *sql.DB
	Repo       *repository.Repository
	Log        *slog.Logger
	Now        func() time.Time
	Campaigns  *CampaignService
	WorkOrders *WorkOrderService
	Inventory  *InventoryService
	Alerts     *AlertService
	Telemetry  *TelemetryService
	Handoffs   *HandoffService
	Queries    *QueryService
	Auth       *AuthService
}

func NewRegistry(database *db.DB, logger *slog.Logger, now func() time.Time) *Registry {
	actual := database.DB
	repo := repository.New(actual)
	r := &Registry{DB: actual, Repo: repo, Log: logger, Now: now}
	r.Campaigns = &CampaignService{repo: repo, log: logger, now: now}
	r.WorkOrders = &WorkOrderService{repo: repo, log: logger, now: now}
	r.Inventory = &InventoryService{repo: repo, log: logger, now: now}
	r.Alerts = &AlertService{repo: repo, log: logger, now: now}
	r.Telemetry = &TelemetryService{repo: repo, log: logger, now: now}
	r.Handoffs = &HandoffService{repo: repo, log: logger, now: now}
	r.Queries = &QueryService{repo: repo, log: logger}
	r.Auth = &AuthService{repo: repo, log: logger, now: now}
	return r
}
