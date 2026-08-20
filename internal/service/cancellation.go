package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/12sub/Reset/internal/domain"
	"github.com/12sub/Reset/internal/paystack"
	"github.com/12sub/Reset/internal/repository"
)

type CancellationService struct {
	repo     *repository.SubscriptionRepo
	paystack *paystack.Client
}

func NewCancellationService(repo *repository.SubscriptionRepo, pc *paystack.Client) *CancellationService {
	return &CancellationService{repo: repo, paystack: pc}
}

func (s *CancellationService) Cancel(ctx context.Context, subID uuid.UUID, reason string) error {
	// 1. Get local subscription
	sub, err := s.repo.GetByID(ctx, subID)
	if err != nil {
		return err
	}
	if sub == nil {
		return fmt.Errorf("subscription not found")
	}
	if sub.Status == domain.StatusCanceled {
		return fmt.Errorf("already canceled")
	}

	// 2. Disable on Paystack (need to fetch email_token - store it in DB in production)
	// For hackathon: Paystack requires the subscription code + email_token
	// You should store email_token when creating. For now, we'll assume you have it.
	// If you didn't store it, you'd need to fetch from Paystack's list endpoint.
	
	// NOTE: In your Create flow, also store the email_token from Paystack response
	// Here we'll use a placeholder - fix this in your actual implementation
	emailToken := "fetch-this-from-paystack-or-store-on-create"
	
	if err := s.paystack.DisableSubscription(sub.PaystackSubID, emailToken); err != nil {
		return fmt.Errorf("paystack disable failed: %w", err)
	}

	// 3. Update local DB
	if err := s.repo.UpdateStatus(ctx, subID, domain.StatusCanceled); err != nil {
		return err
	}

	// TODO: Log cancellation reason to analytics table
	return nil
}