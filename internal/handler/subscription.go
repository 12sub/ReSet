package handler

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"github.com/12sub/Reset/internal/domain"
	"github.com/12sub/Reset/internal/service"
)

type SubscriptionHandler struct {
	subService *service.SubscriptionService
	cancelSvc  *service.CancellationService
}

func NewSubscriptionHandler(sub *service.SubscriptionService, cancel *service.CancellationService) *SubscriptionHandler {
	return &SubscriptionHandler{subService: sub, cancelSvc: cancel}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email         string `json:"email"`
		PlanCode      string `json:"plan_code"`
		Authorization string `json:"authorization"` // Paystack auth code from frontend
		Amount        int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	sub, err := h.subService.Create(r.Context(), req.Email, req.PlanCode, req.Authorization, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	var req domain.CancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	id, err := uuid.Parse(req.SubscriptionID)
	if err != nil {
		http.Error(w, "invalid uuid", 400)
		return
	}

	if err := h.cancelSvc.Cancel(r.Context(), id, req.Reason); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(204)
}

func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid uuid", 400)
		return
	}

	sub, err := h.subService.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if sub == nil {
		http.Error(w, "not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}