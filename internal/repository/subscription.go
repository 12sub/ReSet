package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/12sub/Reset/internal/cache"
	"github.com/12sub/Reset/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionRepo struct {
	db    *sql.DB
	cache *cache.RedisCache
}

func NewSubscriptionRepo(db *sql.DB, c *cache.RedisCache) *SubscriptionRepo {
	return &SubscriptionRepo{db: db, cache: c}
}

func (r *SubscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, customer_email, plan_code, paystack_sub_id, paystack_email_token, status, amount, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		sub.ID, sub.CustomerEmail, sub.PlanCode, sub.PaystackSubID,
		sub.PaystackEmailToken, sub.Status, sub.Amount, sub.CreatedAt,
	)
	return err
}

func (r *SubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	cacheKey := "sub:" + id.String()
	var sub domain.Subscription

	found, err := r.cache.Get(ctx, cacheKey, &sub)
	if err != nil {
		// Log cache error but continue to DB
	}
	if found {
		return &sub, nil
	}

	query := `SELECT id, customer_email, plan_code, paystack_sub_id, paystack_email_token, status, amount, created_at, canceled_at 
	          FROM subscriptions WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	err = row.Scan(&sub.ID, &sub.CustomerEmail, &sub.PlanCode, &sub.PaystackSubID,
		&sub.PaystackEmailToken, &sub.Status, &sub.Amount, &sub.CreatedAt, &sub.CanceledAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	r.cache.Set(ctx, cacheKey, &sub, 10*time.Minute)
	return &sub, nil
}

func (r *SubscriptionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.SubscriptionStatus) error {
	query := `UPDATE subscriptions SET status = $1, canceled_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err == nil {
		r.cache.Delete(ctx, "sub:"+id.String())
	}
	return err
}

func (r *SubscriptionRepo) UpdateStatusByPaystackID(ctx context.Context, paystackSubID string, status domain.SubscriptionStatus) error {
	query := `UPDATE subscriptions SET status = $1, canceled_at = NOW() WHERE paystack_sub_id = $2`
	_, err := r.db.ExecContext(ctx, query, status, paystackSubID)
	if err == nil {
		r.cache.Delete(ctx, "sub:"+paystackSubID)
	}
	return err
}