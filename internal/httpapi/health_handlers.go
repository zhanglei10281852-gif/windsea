package httpapi

import (
	"github.com/zhanglei10281852-gif/windsea/internal/service"
	"net/http"
)

func (h *handler) healthDetails(w http.ResponseWriter, r *http.Request) {
	health := service.CheckHealth(r.Context(), h.services.DB, h.services.Now())
	if health.Status != "ready" {
		writeJSON(w, http.StatusServiceUnavailable, health)
		return
	}
	writeJSON(w, http.StatusOK, health)
}
