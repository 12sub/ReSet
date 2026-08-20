package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/12sub/Reset/internal/domain"
	"github.com/12sub/Reset/internal/service"
	"github.com/12sub/Reset/internal/templates"
)

type SubscriptionHandler struct {
	subService *service.SubscriptionService
	cancelSvc  *service.CancellationService
}

func NewSubscriptionHandler(sub *service.SubscriptionService, cancel *service.CancellationService) *SubscriptionHandler {
	return &SubscriptionHandler{subService: sub, cancelSvc: cancel}
}

func (h *SubscriptionHandler) Index(w http.ResponseWriter, r *http.Request) {
	templates.T.ExecuteTemplate(w, "layout", nil)
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderError(w, "Invalid form data")
		return
	}

	email := r.FormValue("email")
	planCode := r.FormValue("plan_code")
	reference := r.FormValue("reference")
	amountStr := r.FormValue("amount")

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		renderError(w, "Invalid amount")
		return
	}

	sub, err := h.subService.Create(r.Context(), email, planCode, reference, amount)
	if err != nil {
		renderError(w, err.Error())
		return
	}

	data := map[string]interface{}{
		"ID":            sub.ID.String(),
		"CustomerEmail": sub.CustomerEmail,
		"PlanCode":      sub.PlanCode,
		"Amount":        sub.Amount / 100,
		"CreatedAt":     sub.CreatedAt.Format("Jan 02, 2006 15:04"),
	}

	templates.T.ExecuteTemplate(w, "active", data)
}

func (h *SubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderError(w, "Invalid form data")
		return
	}

	subIDStr := r.FormValue("subscription_id")
	if subIDStr == "" {
		var req domain.CancelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			renderError(w, "Invalid request")
			return
		}
		subIDStr = req.SubscriptionID
	}

	id, err := uuid.Parse(subIDStr)
	if err != nil {
		renderError(w, "Invalid subscription ID")
		return
	}

	sub, err := h.subService.Get(r.Context(), id)
	if err != nil || sub == nil {
		renderError(w, "Subscription not found")
		return
	}

	if err := h.cancelSvc.Cancel(r.Context(), id, "User requested via ReSet"); err != nil {
		renderError(w, err.Error())
		return
	}

	data := map[string]interface{}{
		"ID":            sub.ID.String(),
		"CustomerEmail": sub.CustomerEmail,
		"PlanCode":      sub.PlanCode,
		"CanceledAt":    "Just now",
	}

	templates.T.ExecuteTemplate(w, "canceled", data)
}

func renderError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusOK)
	templates.T.ExecuteTemplate(w, "error", map[string]string{"Message": msg})
}