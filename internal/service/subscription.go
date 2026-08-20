package service

import (
	"context"
	"fmt"
	"time"

	"github.com/12sub/Reset/internal/domain"
	"github.com/12sub/Reset/internal/paystack"
	"github.com/12sub/Reset/internal/repository"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo     *repository.SubscriptionRepo
	paystack *paystack.Client
}

func NewSubscriptionService(repo *repository.SubscriptionRepo, pc *paystack.Client) *SubscriptionService {
	return &SubscriptionService{repo: repo, paystack: pc}
}

func (s *SubscriptionService) Create(ctx context.Context, email, planCode, reference string, amount int) (*domain.Subscription, error) {
	// 1. Verify transaction with Paystack
	verify, err := s.paystack.VerifyTransaction(ctx, reference)
	if err != nil {
		return nil, fmt.Errorf("transaction verification failed: %w", err)
	}
	if verify.Status != "success" {
		return nil, fmt.Errorf("transaction not successful: %s", verify.Status)
	}

	// 2. Create subscription using the authorization code from verification
	authCode := verify.Authorization.AuthorizationCode
	psSub, err := s.paystack.CreateSubscription(email, planCode, authCode)
	if err != nil {
		return nil, fmt.Errorf("subscription creation failed: %w", err)
	}

	// 3. Save locally WITH the email_token (needed for cancellation)
	sub := &domain.Subscription{
		ID:                 uuid.New(),
		CustomerEmail:      email,
		PlanCode:           planCode,
		PaystackSubID:      psSub.SubscriptionCode,
		PaystackEmailToken: psSub.EmailToken,
		Status:             domain.StatusActive,
		Amount:             amount,
		CreatedAt:          time.Now(),
	}
	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return sub, nil
}

func (s *SubscriptionService) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}