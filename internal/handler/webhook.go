package handler

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/12sub/Reset/internal/domain"
	"github.com/12sub/Reset/internal/repository"
)

type WebhookHandler struct {
	repo       *repository.SubscriptionRepo
	secretKey  string
}

func NewWebhookHandler(repo *repository.SubscriptionRepo, secretKey string) *WebhookHandler {
	return &WebhookHandler{
		repo:      repo,
		secretKey: secretKey,
	}
}
// verifySignature checks the x-paystack-signature header against HMAC-SHA512 of the body
func (h *WebhookHandler) verifySignature(body []byte, signature string) bool {
	if signature == "" || h.secretKey == "" {
		return false
	}

	mac := hmac.New(sha512.New, []byte(h.secretKey))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	signature := r.Header.Get("x-paystack-signature")
	if !h.verifySignature(body, signature) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	var event struct {
		Event string          `json:"event"`
		Data  json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}

	switch event.Event {
	case "subscription.disable":
		var data struct {
			SubscriptionCode string `json:"subscription_code"`
			Status           string `json:"status"`
			Customer         struct {
				Email string `json:"email"`
			} `json:"customer"`
		}
		json.Unmarshal(event.Data, &data)
		// Sync local DB if Paystack initiated the cancellation
		_ = h.repo.UpdateStatusByPaystackID(r.Context(), data.SubscriptionCode, domain.StatusCanceled)

	case "charge.success":
		var data struct {
			Amount    int    `json:"amount"`
			Reference string `json:"reference"`
			Customer  struct {
				Email string `json:"email"`
			} `json:"customer"`
			Plan struct {
				PlanCode string `json:"plan_code"`
			} `json:"plan"`
		}
		json.Unmarshal(event.Data, &data)
		// Log successful charge — extend subscription period, send receipt, etc.
		// fmt.Printf("Payment success: %d kobo from %s\n", data.Amount, data.Customer.Email)

	case "invoice.create":
		// Paystack created an invoice for the next charge
		// Useful for notifying user before deduction

	case "invoice.payment_failed":
		var data struct {
			Subscription struct {
				SubscriptionCode string `json:"subscription_code"`
			} `json:"subscription"`
			Customer struct {
				Email string `json:"email"`
			} `json:"customer"`
		}
		json.Unmarshal(event.Data, &data)
		// Handle failed payment — notify user, retry logic, etc.
	}

	w.WriteHeader(http.StatusOK)
}