package service

import (
	"context"
	"fmt"
	"time"

	"github.com/12sub/Reset/internal/domain"
	"github.com/12sub/Reset/internal/paystack"
	"github.com/12sub/Reset/internal/python"
	"github.com/12sub/Reset/internal/repository"
	"github.com/google/uuid"
)

type TrackedService struct {
	repo     *repository.TrackedRepo
	python   *python.Client
	paystack *paystack.Client
}

func NewTrackedService(repo *repository.TrackedRepo, py *python.Client, pc *paystack.Client) *TrackedService {
	return &TrackedService{repo: repo, python: py, paystack: pc}
}

func (s *TrackedService) Detect(ctx context.Context, email, text string) (*domain.TrackedSubscription, error) {
	result, err := s.python.Classify(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("AI classification failed: %w", err)
	}
	if !result.IsSubscription {
		return nil, fmt.Errorf("not a subscription: %s", text)
	}

	sub := &domain.TrackedSubscription{
		ID:         uuid.New(),
		UserEmail:  email,
		Merchant:   result.Merchant,
		Category:   result.Category,
		Amount:     result.Amount,
		RawText:    result.RawText,
		Status:     domain.StatusActive,
		Provider:   "direct",
		Confidence: result.Confidence,
		MatchType:  result.MatchType,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *TrackedService) ListByEmail(ctx context.Context, email string) ([]domain.TrackedSubscription, error) {
	return s.repo.ListByEmail(ctx, email)
}

func (s *TrackedService) Cancel(ctx context.Context, id uuid.UUID) (*domain.TrackedSubscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, fmt.Errorf("subscription not found")
	}
	if sub.Status == domain.StatusCanceled {
		return nil, fmt.Errorf("already canceled")
	}

	if sub.Provider == "paystack" && sub.ProviderRef != "" && sub.EmailToken != "" {
		if err := s.paystack.DisableSubscription(sub.ProviderRef, sub.EmailToken); err != nil {
			return nil, fmt.Errorf("paystack cancellation failed: %w", err)
		}
	}

	if err := s.repo.UpdateStatus(ctx, id, domain.StatusCanceled); err != nil {
		return nil, err
	}

	sub.Status = domain.StatusCanceled
	now := time.Now()
	sub.CanceledAt = &now
	return sub, nil
}