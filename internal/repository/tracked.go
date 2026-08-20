package repository

import (
	"context"
	"database/sql"

	"github.com/12sub/Reset/internal/domain"
	"github.com/google/uuid"
)

type TrackedRepo struct {
	db *sql.DB
}

func NewTrackedRepo(db *sql.DB) *TrackedRepo {
	return &TrackedRepo{db: db}
}

func (r *TrackedRepo) Create(ctx context.Context, sub *domain.TrackedSubscription) error {
	query := `
		INSERT INTO tracked_subscriptions 
		(id, user_email, merchant, category, amount, raw_text, status, provider, provider_ref, email_token, confidence, match_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.ExecContext(ctx, query,
		sub.ID, sub.UserEmail, sub.Merchant, sub.Category, sub.Amount,
		sub.RawText, sub.Status, sub.Provider, sub.ProviderRef, sub.EmailToken,
		sub.Confidence, sub.MatchType, sub.CreatedAt,
	)
	return err
}

func (r *TrackedRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.TrackedSubscription, error) {
	var sub domain.TrackedSubscription
	query := `
		SELECT id, user_email, merchant, category, amount, raw_text, status, provider, provider_ref, email_token, confidence, match_type, created_at, canceled_at
		FROM tracked_subscriptions WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&sub.ID, &sub.UserEmail, &sub.Merchant, &sub.Category, &sub.Amount,
		&sub.RawText, &sub.Status, &sub.Provider, &sub.ProviderRef, &sub.EmailToken,
		&sub.Confidence, &sub.MatchType, &sub.CreatedAt, &sub.CanceledAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &sub, err
}

func (r *TrackedRepo) ListByEmail(ctx context.Context, email string) ([]domain.TrackedSubscription, error) {
	query := `
		SELECT id, user_email, merchant, category, amount, raw_text, status, provider, provider_ref, email_token, confidence, match_type, created_at, canceled_at
		FROM tracked_subscriptions WHERE user_email = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []domain.TrackedSubscription
	for rows.Next() {
		var sub domain.TrackedSubscription
		err := rows.Scan(
			&sub.ID, &sub.UserEmail, &sub.Merchant, &sub.Category, &sub.Amount,
			&sub.RawText, &sub.Status, &sub.Provider, &sub.ProviderRef, &sub.EmailToken,
			&sub.Confidence, &sub.MatchType, &sub.CreatedAt, &sub.CanceledAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, rows.Err()
}

func (r *TrackedRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.SubscriptionStatus) error {
	query := `UPDATE tracked_subscriptions SET status = $1, canceled_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}