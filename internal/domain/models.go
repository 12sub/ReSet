package domain

import (
	"time"
	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	StatusActive   SubscriptionStatus = "active"
	StatusCanceled SubscriptionStatus = "canceled"
	StatusPending  SubscriptionStatus = "pending"
)

type Subscription struct {
	ID            uuid.UUID          `json:"id"`
	CustomerEmail string             `json:"customer_email"`
	PlanCode      string             `json:"plan_code"`      // Paystack plan code
	PaystackSubID string             `json:"paystack_sub_id"` // Paystack subscription ID
	PaystackEmailToken string             `json:"-"`               // Paystack email token
	Status        SubscriptionStatus `json:"status"`
	Amount        int                `json:"amount"`         // in kobo
	CreatedAt     time.Time          `json:"created_at"`
	CanceledAt    *time.Time         `json:"canceled_at,omitempty"`
}

type CancelRequest struct {
	SubscriptionID string `json:"subscription_id"`
	Reason         string `json:"reason,omitempty"`
}

type TrackedSubscription struct {
	ID          uuid.UUID          `json:"id"`
	UserEmail   string             `json:"user_email"`
	Merchant    string             `json:"merchant"`
	Category    string             `json:"category"`
	Amount      float64            `json:"amount"`
	RawText     string             `json:"raw_text"`
	Status      SubscriptionStatus `json:"status"`
	Provider    string             `json:"provider"`
	ProviderRef string             `json:"provider_ref"`
	EmailToken  string             `json:"-"`
	Confidence  float64            `json:"confidence"`
	MatchType   string             `json:"match_type"`
	CreatedAt   time.Time          `json:"created_at"`
	CanceledAt  *time.Time         `json:"canceled_at,omitempty"`
}