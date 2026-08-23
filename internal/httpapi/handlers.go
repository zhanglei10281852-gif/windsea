package httpapi

import (
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"github.com/zhanglei10281852-gif/windsea/internal/observability"
	"net/http"
	"strconv"
	"time"
)

func (h *handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
func (h *handler) ready(w http.ResponseWriter, r *http.Request) {
	if h.services.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready"})
		return
	}
	if err := h.services.DB.PingContext(r.Context()); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

type campaignInput struct {
	ID, FarmID, Name string
	StartAt, EndAt   time.Time
}

func (h *handler) createCampaign(w http.ResponseWriter, r *http.Request) {
	var input campaignInput
	if err := decode(r, &input); err != nil {
		fail(w, err)
		return
	}
	campaign := domain.Campaign{ID: input.ID, FarmID: input.FarmID, Name: input.Name, StartAt: input.StartAt, EndAt: input.EndAt}
	if err := h.services.Campaigns.Create(r.Context(), campaign); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, campaign)
}
func (h *handler) reviewCampaign(w http.ResponseWriter, r *http.Request) {
	h.transitionCampaign(w, r, string(domain.CampaignReview))
}
func (h *handler) publishCampaign(w http.ResponseWriter, r *http.Request) {
	h.transitionCampaign(w, r, string(domain.CampaignPublished))
}
func (h *handler) closeCampaign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	version, _ := strconv.Atoi(r.URL.Query().Get("version"))
	campaign, err := h.services.Campaigns.Close(r.Context(), id, version)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, campaign)
}
func (h *handler) transitionCampaign(w http.ResponseWriter, r *http.Request, target string) {
	id := r.PathValue("id")
	version, _ := strconv.Atoi(r.URL.Query().Get("version"))
	var campaign domain.Campaign
	var err error
	if target == string(domain.CampaignReview) {
		campaign, err = h.services.Campaigns.SubmitReview(r.Context(), id, version)
	} else {
		campaign, err = h.services.Campaigns.Publish(r.Context(), id, version)
	}
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, campaign)
}

type workOrderInput struct {
	ID, FarmID, CampaignID, TurbineID, Title string
	Priority                                 int
}

func (h *handler) createWorkOrder(w http.ResponseWriter, r *http.Request) {
	var input workOrderInput
	if err := decode(r, &input); err != nil {
		fail(w, err)
		return
	}
	order := domain.WorkOrder{ID: input.ID, FarmID: input.FarmID, CampaignID: input.CampaignID, TurbineID: input.TurbineID, Title: input.Title, Priority: input.Priority}
	if err := h.services.WorkOrders.Create(r.Context(), order); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, order)
}
func (h *handler) listWorkOrders(w http.ResponseWriter, r *http.Request) {
	limit, offset := 50, 0
	if value, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && value > 0 {
		limit = value
	}
	if value, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && value >= 0 {
		offset = value
	}
	page, err := h.services.WorkOrders.List(r.Context(), r.URL.Query().Get("farm_id"), r.URL.Query().Get("state"), limit, offset)
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}
func (h *handler) assignWorkOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	assignee := r.URL.Query().Get("assignee_id")
	version, _ := strconv.Atoi(r.URL.Query().Get("version"))
	if err := h.services.WorkOrders.Assign(r.Context(), id, assignee, version); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "assigned"})
}
func (h *handler) startWorkOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	version, _ := strconv.Atoi(r.URL.Query().Get("version"))
	if err := h.services.WorkOrders.Start(r.Context(), id, version); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "in_progress"})
}
func (h *handler) completeWorkOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	version, _ := strconv.Atoi(r.URL.Query().Get("version"))
	if err := h.services.WorkOrders.Complete(r.Context(), id, version); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

type reservationInput struct {
	ID, PartID, WorkOrderID, RequestedBy string
	Quantity                             int
}

func (h *handler) holdReservation(w http.ResponseWriter, r *http.Request) {
	var input reservationInput
	if err := decode(r, &input); err != nil {
		fail(w, err)
		return
	}
	if err := h.services.Inventory.Hold(r.Context(), input.ID, input.PartID, input.WorkOrderID, input.RequestedBy, input.Quantity); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": input.ID, "state": "held"})
}
func (h *handler) consumeReservation(w http.ResponseWriter, r *http.Request) {
	if err := h.services.Inventory.Consume(r.Context(), r.PathValue("id")); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": "consumed"})
}
func (h *handler) releaseReservation(w http.ResponseWriter, r *http.Request) {
	if err := h.services.Inventory.Release(r.Context(), r.PathValue("id")); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": "released"})
}

type alertInput struct {
	ID, FarmID, TurbineID, Code, Severity, Message string
	OccurredAt                                     time.Time
}

func (h *handler) ingestAlert(w http.ResponseWriter, r *http.Request) {
	var input alertInput
	if err := decode(r, &input); err != nil {
		fail(w, err)
		return
	}
	alert := domain.Alert{ID: input.ID, FarmID: input.FarmID, TurbineID: input.TurbineID, Code: input.Code, Severity: input.Severity, Message: input.Message, OccurredAt: &input.OccurredAt}
	if err := h.services.Alerts.Ingest(r.Context(), alert); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, alert)
}
func (h *handler) ackAlert(w http.ResponseWriter, r *http.Request) {
	if err := h.services.Alerts.Acknowledge(r.Context(), r.PathValue("id")); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": "acknowledged"})
}
func (h *handler) resolveAlert(w http.ResponseWriter, r *http.Request) {
	if err := h.services.Alerts.Resolve(r.Context(), r.PathValue("id")); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": "resolved"})
}
func (h *handler) listAlerts(w http.ResponseWriter, r *http.Request) {
	items, err := h.services.Alerts.Open(r.Context(), r.URL.Query().Get("farm_id"))
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type telemetryInput struct {
	ID, FarmID, Source string
	Samples            int
}

func (h *handler) receiveTelemetry(w http.ResponseWriter, r *http.Request) {
	var input telemetryInput
	if err := decode(r, &input); err != nil {
		fail(w, err)
		return
	}
	batch := domain.TelemetryBatch{ID: input.ID, FarmID: input.FarmID, Source: input.Source, Samples: input.Samples}
	if err := h.services.Telemetry.Receive(r.Context(), batch); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, batch)
}
func (h *handler) completeTelemetry(w http.ResponseWriter, r *http.Request) {
	accepted, _ := strconv.Atoi(r.URL.Query().Get("accepted"))
	rejected, _ := strconv.Atoi(r.URL.Query().Get("rejected"))
	if err := h.services.Telemetry.Complete(r.Context(), r.PathValue("id"), accepted, rejected); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": "completed"})
}

type handoffInput struct{ ID, WorkOrderID, ContractorID, Notes string }

func (h *handler) offerHandoff(w http.ResponseWriter, r *http.Request) {
	var input handoffInput
	if err := decode(r, &input); err != nil {
		fail(w, err)
		return
	}
	if err := h.services.Handoffs.Offer(r.Context(), domain.Handoff{ID: input.ID, WorkOrderID: input.WorkOrderID, ContractorID: input.ContractorID, Notes: input.Notes}); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": input.ID, "state": "offered"})
}
func (h *handler) acceptHandoff(w http.ResponseWriter, r *http.Request) {
	if err := h.services.Handoffs.Accept(r.Context(), r.PathValue("id"), r.URL.Query().Get("notes")); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": "accepted"})
}
func (h *handler) completeHandoff(w http.ResponseWriter, r *http.Request) {
	if err := h.services.Handoffs.Complete(r.Context(), r.PathValue("id"), r.URL.Query().Get("notes")); err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"state": "completed"})
}

var _ = observability.RequestID
