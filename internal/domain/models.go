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
	Status        SubscriptionStatus `json:"status"`
	Amount        int                `json:"amount"`         // in kobo
	CreatedAt     time.Time          `json:"created_at"`
	CanceledAt    *time.Time         `json:"canceled_at,omitempty"`
}

type CancelRequest struct {
	SubscriptionID string `json:"subscription_id"`
	Reason         string `json:"reason,omitempty"`
}