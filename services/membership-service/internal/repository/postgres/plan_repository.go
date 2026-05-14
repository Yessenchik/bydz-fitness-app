package postgres

import (
	"context"
	"database/sql"

	"github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/domain"
)

type PlanRepository struct {
	db *sql.DB
}

func NewPlanRepository(db *sql.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) Create(ctx context.Context, p *domain.Plan) error {
	query := `
		INSERT INTO membership_plans (name, duration_days, price_kzt, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		p.Name,
		p.DurationDays,
		p.PriceKZT,
		p.IsActive,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PlanRepository) GetByID(ctx context.Context, id string) (*domain.Plan, error) {
	query := `
		SELECT id, name, duration_days, price_kzt, is_active, created_at, updated_at
		FROM membership_plans
		WHERE id = $1
	`

	var p domain.Plan

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.DurationDays,
		&p.PriceKZT,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PlanRepository) List(ctx context.Context, onlyActive bool) ([]*domain.Plan, error) {
	query := `
		SELECT id, name, duration_days, price_kzt, is_active, created_at, updated_at
		FROM membership_plans
	`

	if onlyActive {
		query += ` WHERE is_active = true`
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*domain.Plan

	for rows.Next() {
		var p domain.Plan

		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.DurationDays,
			&p.PriceKZT,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}

		plans = append(plans, &p)
	}

	return plans, rows.Err()
}
