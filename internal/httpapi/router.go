package httpapi

import (
	"encoding/json"
	"github.com/zhanglei10281852-gif/windsea/internal/service"
	"log/slog"
	"net/http"
)

func NewRouter(services *service.Registry, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	h := &handler{services: services, logger: logger}
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("GET /readyz", h.ready)
	mux.HandleFunc("POST /v1/campaigns", h.createCampaign)
	mux.HandleFunc("POST /v1/campaigns/{id}/review", h.reviewCampaign)
	mux.HandleFunc("POST /v1/campaigns/{id}/publish", h.publishCampaign)
	mux.HandleFunc("POST /v1/campaigns/{id}/close", h.closeCampaign)
	mux.HandleFunc("POST /v1/work-orders", h.createWorkOrder)
	mux.HandleFunc("GET /v1/work-orders", h.listWorkOrders)
	mux.HandleFunc("POST /v1/work-orders/{id}/assign", h.assignWorkOrder)
	mux.HandleFunc("POST /v1/work-orders/{id}/start", h.startWorkOrder)
	mux.HandleFunc("POST /v1/work-orders/{id}/complete", h.completeWorkOrder)
	mux.HandleFunc("POST /v1/reservations", h.holdReservation)
	mux.HandleFunc("POST /v1/reservations/{id}/consume", h.consumeReservation)
	mux.HandleFunc("POST /v1/reservations/{id}/release", h.releaseReservation)
	mux.HandleFunc("POST /v1/alerts", h.ingestAlert)
	mux.HandleFunc("POST /v1/alerts/{id}/ack", h.ackAlert)
	mux.HandleFunc("POST /v1/alerts/{id}/resolve", h.resolveAlert)
	mux.HandleFunc("GET /v1/alerts", h.listAlerts)
	mux.HandleFunc("POST /v1/telemetry/batches", h.receiveTelemetry)
	mux.HandleFunc("POST /v1/telemetry/batches/{id}/complete", h.completeTelemetry)
	mux.HandleFunc("POST /v1/handoffs", h.offerHandoff)
	mux.HandleFunc("POST /v1/handoffs/{id}/accept", h.acceptHandoff)
	mux.HandleFunc("POST /v1/handoffs/{id}/complete", h.completeHandoff)
	return withRecovery(withRequestID(mux, logger), logger)
}

type handler struct {
	services *service.Registry
	logger   *slog.Logger
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func decode(r *http.Request, value any) error { return json.NewDecoder(r.Body).Decode(value) }
