package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"membership-service/internal/domain/entity"
)

type SubscriptionPostgresRepository struct {
	db *sql.DB
}

func NewSubscriptionPostgresRepository(db *sql.DB) *SubscriptionPostgresRepository {
	return &SubscriptionPostgresRepository{
		db: db,
	}
}

func (r *SubscriptionPostgresRepository) Create(
	ctx context.Context,
	subscription *entity.Subscription,
) error {
	query := `
		INSERT INTO subscriptions (
			id,
			user_id,
			membership_plan_id,
			start_date,
			end_date,
			is_active,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		subscription.ID,
		subscription.UserID,
		subscription.MembershipPlanID,
		subscription.StartDate,
		subscription.EndDate,
		subscription.IsActive,
		string(subscription.Status),
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)

	return err
}

func (r *SubscriptionPostgresRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.Subscription, error) {
	query := `
		SELECT
			id,
			user_id,
			membership_plan_id,
			start_date,
			end_date,
			is_active,
			status,
			created_at,
			updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var subscription entity.Subscription
	var status string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.MembershipPlanID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.IsActive,
		&status,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	subscription.Status = entity.SubscriptionStatus(status)

	return &subscription, nil
}

func (r *SubscriptionPostgresRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.Subscription, error) {
	query := `
		SELECT
			id,
			user_id,
			membership_plan_id,
			start_date,
			end_date,
			is_active,
			status,
			created_at,
			updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var subscription entity.Subscription
	var status string

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.MembershipPlanID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.IsActive,
		&status,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	subscription.Status = entity.SubscriptionStatus(status)

	return &subscription, nil
}

func (r *SubscriptionPostgresRepository) GetActiveByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.Subscription, error) {
	query := `
		SELECT
			id,
			user_id,
			membership_plan_id,
			start_date,
			end_date,
			is_active,
			status,
			created_at,
			updated_at
		FROM subscriptions
		WHERE user_id = $1
		  AND is_active = true
		  AND status = 'ACTIVE'
		  AND end_date > NOW()
		ORDER BY end_date DESC
		LIMIT 1
	`

	var subscription entity.Subscription
	var status string

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.MembershipPlanID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.IsActive,
		&status,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	subscription.Status = entity.SubscriptionStatus(status)

	return &subscription, nil
}

func (r *SubscriptionPostgresRepository) Update(
	ctx context.Context,
	subscription *entity.Subscription,
) error {
	query := `
		UPDATE subscriptions
		SET
			membership_plan_id = $1,
			start_date = $2,
			end_date = $3,
			is_active = $4,
			status = $5,
			updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		subscription.MembershipPlanID,
		subscription.StartDate,
		subscription.EndDate,
		subscription.IsActive,
		string(subscription.Status),
		subscription.UpdatedAt,
		subscription.ID,
	)

	return err
}

func (r *SubscriptionPostgresRepository) Cancel(
	ctx context.Context,
	userID uuid.UUID,
) error {
	query := `
		UPDATE subscriptions
		SET
			is_active = false,
			status = 'CANCELLED',
			updated_at = NOW()
		WHERE user_id = $1
		  AND is_active = true
		  AND status = 'ACTIVE'
	`

	_, err := r.db.ExecContext(ctx, query, userID)

	return err
}

func (r *SubscriptionPostgresRepository) ExpireOldSubscriptions(
	ctx context.Context,
) error {
	query := `
		UPDATE subscriptions
		SET
			is_active = false,
			status = 'EXPIRED',
			updated_at = NOW()
		WHERE end_date <= NOW()
		  AND is_active = true
		  AND status = 'ACTIVE'
	`

	_, err := r.db.ExecContext(ctx, query)

	return err
}

func (r *SubscriptionPostgresRepository) CountActive(
	ctx context.Context,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM subscriptions
		WHERE is_active = true
		  AND status = 'ACTIVE'
		  AND end_date > NOW()
	`

	var count int64

	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
