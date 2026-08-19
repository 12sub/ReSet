package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/12sub/Reset/internal/domain"
)

type SubscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, customer_email, plan_code, paystack_sub_id, status, amount, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		sub.ID, sub.CustomerEmail, sub.PlanCode, sub.PaystackSubID,
		sub.Status, sub.Amount, sub.CreatedAt,
	)
	return err
}

func (r *SubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var sub domain.Subscription
	query := `SELECT id, customer_email, plan_code, paystack_sub_id, status, amount, created_at, canceled_at 
	          FROM subscriptions WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&sub.ID, &sub.CustomerEmail, &sub.PlanCode, &sub.PaystackSubID,
		&sub.Status, &sub.Amount, &sub.CreatedAt, &sub.CanceledAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &sub, err
}

func (r *SubscriptionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.SubscriptionStatus) error {
	query := `UPDATE subscriptions SET status = $1, canceled_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *SubscriptionRepo) ListActive(ctx context.Context) ([]domain.Subscription, error) {
	query := `SELECT id, customer_email, plan_code, paystack_sub_id, status, amount, created_at 
	          FROM subscriptions WHERE status = 'active'`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		if err := rows.Scan(&sub.ID, &sub.CustomerEmail, &sub.PlanCode, &sub.PaystackSubID,
			&sub.Status, &sub.Amount, &sub.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, rows.Err()
}
func (r *SubscriptionRepo) UpdateStatusByPaystackID(ctx context.Context, paystackSubID string, status domain.SubscriptionStatus) error {
	query := `UPDATE subscriptions SET status = $1, canceled_at = NOW() WHERE paystack_sub_id = $2`
	_, err := r.db.ExecContext(ctx, query, status, paystackSubID)
	return err
}