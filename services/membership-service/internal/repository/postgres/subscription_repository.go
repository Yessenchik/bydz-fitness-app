package postgres

import (
	"context"
	"database/sql"

	"github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/domain"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, plan_id, status, starts_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		s.UserID,
		s.PlanID,
		s.Status,
		s.StartsAt,
		s.ExpiresAt,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*domain.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, starts_at, expires_at, cancelled_at, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var s domain.Subscription

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.PlanID,
		&s.Status,
		&s.StartsAt,
		&s.ExpiresAt,
		&s.CancelledAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *SubscriptionRepository) GetActiveByUserID(ctx context.Context, userID string) (*domain.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, starts_at, expires_at, cancelled_at, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		  AND status = 'active'
		  AND expires_at > NOW()
		ORDER BY expires_at DESC
		LIMIT 1
	`

	var s domain.Subscription

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&s.ID,
		&s.UserID,
		&s.PlanID,
		&s.Status,
		&s.StartsAt,
		&s.ExpiresAt,
		&s.CancelledAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *SubscriptionRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, starts_at, expires_at, cancelled_at, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*domain.Subscription

	for rows.Next() {
		var s domain.Subscription

		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.PlanID,
			&s.Status,
			&s.StartsAt,
			&s.ExpiresAt,
			&s.CancelledAt,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, &s)
	}

	return subscriptions, rows.Err()
}

func (r *SubscriptionRepository) Cancel(ctx context.Context, id string) (*domain.Subscription, error) {
	query := `
		UPDATE subscriptions
		SET status = 'cancelled',
		    cancelled_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, plan_id, status, starts_at, expires_at, cancelled_at, created_at, updated_at
	`

	var s domain.Subscription

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.PlanID,
		&s.Status,
		&s.StartsAt,
		&s.ExpiresAt,
		&s.CancelledAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
