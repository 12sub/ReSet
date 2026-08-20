package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/12sub/Reset/internal/domain"
	"github.com/12sub/Reset/internal/service"
)

// BulkCancelRequest represents a batch cancellation request
type BulkCancelRequest struct {
	IDs []string `json:"ids"`
}

// BulkCancelResponse represents the result of a batch cancellation
type BulkCancelResponse struct {
	Success      int      `json:"success"`
	Failed       int      `json:"failed"`
	TotalSavings float64  `json:"total_savings"`
	Results      []CancelResult `json:"results"`
}

type CancelResult struct {
	ID       string `json:"id"`
	Merchant string `json:"merchant"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
	Amount   float64 `json:"amount"`
}

// BulkCancelHandler handles batch cancellation of tracked subscriptions
type BulkCancelHandler struct {
	trackedService *service.TrackedService
}

func NewBulkCancelHandler(svc *service.TrackedService) *BulkCancelHandler {
	return &BulkCancelHandler{trackedService: svc}
}

// HandleBulkCancel processes multiple subscription cancellations
func (h *BulkCancelHandler) HandleBulkCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BulkCancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Try form data fallback
		if err := r.ParseForm(); err == nil {
			idsStr := r.FormValue("ids")
			req.IDs = strings.Split(idsStr, ",")
		} else {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
	}

	resp := BulkCancelResponse{
		Results: make([]CancelResult, 0, len(req.IDs)),
	}

	for _, idStr := range req.IDs {
		if idStr == "" {
			continue
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			resp.Results = append(resp.Results, CancelResult{
				ID:      idStr,
				Success: false,
				Error:   "Invalid UUID",
			})
			resp.Failed++
			continue
		}

		sub, err := h.trackedService.Cancel(r.Context(), id)
		if err != nil {
			resp.Results = append(resp.Results, CancelResult{
				ID:      idStr,
				Success: false,
				Error:   err.Error(),
			})
			resp.Failed++
			continue
		}

		resp.Results = append(resp.Results, CancelResult{
			ID:       idStr,
			Merchant: sub.Merchant,
			Success:  true,
			Amount:   sub.Amount,
		})
		resp.Success++
		resp.TotalSavings += sub.Amount
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleListActive returns active subscriptions for a user as JSON (for test dashboard)
func (h *BulkCancelHandler) HandleListActive(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		email = "demo@reset.ng"
	}

	subs, err := h.trackedService.ListByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter to only active
	var active []domain.TrackedSubscription
	for _, sub := range subs {
		if sub.Status == domain.StatusActive {
			active = append(active, sub)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(active)
}