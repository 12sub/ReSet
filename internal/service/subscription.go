package service

import (
	"context"
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

func (s *SubscriptionService) Create(ctx context.Context, email, planCode, authorization string, amount int) (*domain.Subscription, error) {
	// 1. Create on Paystack first
	psSub, err := s.paystack.CreateSubscription(email, planCode, authorization)
	if err != nil {
		return nil, err
	}

	// 2. Save locally
	sub := &domain.Subscription{
		ID:            uuid.New(),
		CustomerEmail: email,
		PlanCode:      planCode,
		PaystackSubID: psSub.SubscriptionCode,
		Status:        domain.StatusActive,
		Amount:        amount,
		CreatedAt:     time.Now(),
	}
	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err // In production: attempt to disable Paystack sub
	}
	return sub, nil
}

func (s *SubscriptionService) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}
