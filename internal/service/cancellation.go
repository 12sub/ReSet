package service

import (
	"context"
	"fmt"

	"github.com/12sub/Reset/internal/domain"
	"github.com/12sub/Reset/internal/paystack"
	"github.com/12sub/Reset/internal/repository"
	"github.com/google/uuid"
)

type CancellationService struct {
	repo     *repository.SubscriptionRepo
	paystack *paystack.Client
}

func NewCancellationService(repo *repository.SubscriptionRepo, pc *paystack.Client) *CancellationService {
	return &CancellationService{repo: repo, paystack: pc}
}

func (s *CancellationService) Cancel(ctx context.Context, subID uuid.UUID, reason string) error {
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
	if sub.PaystackEmailToken == "" {
		return fmt.Errorf("email token missing — cannot cancel")
	}

	// Use the REAL email token from the database
	if err := s.paystack.DisableSubscription(sub.PaystackSubID, sub.PaystackEmailToken); err != nil {
		return fmt.Errorf("paystack disable failed: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, subID, domain.StatusCanceled); err != nil {
		return err
	}

	return nil
}