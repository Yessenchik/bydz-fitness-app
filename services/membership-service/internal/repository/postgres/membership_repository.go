package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"membership-service/internal/domain/entity"
)

type MembershipPostgresRepository struct {
	db *sql.DB
}

func NewMembershipPostgresRepository(db *sql.DB) *MembershipPostgresRepository {
	return &MembershipPostgresRepository{
		db: db,
	}
}

func (r *MembershipPostgresRepository) Create(
	ctx context.Context,
	plan *entity.MembershipPlan,
) error {
	query := `
		INSERT INTO membership_plans (
			id,
			name,
			duration_months,
			price,
			description,
			is_deleted,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		plan.ID,
		plan.Name,
		plan.DurationMonths,
		plan.Price,
		plan.Description,
		plan.IsDeleted,
		plan.CreatedAt,
		plan.UpdatedAt,
	)

	return err
}

func (r *MembershipPostgresRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.MembershipPlan, error) {
	query := `
		SELECT
			id,
			name,
			duration_months,
			price,
			COALESCE(description, ''),
			is_deleted,
			created_at,
			updated_at
		FROM membership_plans
		WHERE id = $1
		  AND is_deleted = false
	`

	var plan entity.MembershipPlan

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID,
		&plan.Name,
		&plan.DurationMonths,
		&plan.Price,
		&plan.Description,
		&plan.IsDeleted,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (r *MembershipPostgresRepository) List(
	ctx context.Context,
	limit int,
	offset int,
) ([]entity.MembershipPlan, error) {
	query := `
		SELECT
			id,
			name,
			duration_months,
			price,
			COALESCE(description, ''),
			is_deleted,
			created_at,
			updated_at
		FROM membership_plans
		WHERE is_deleted = false
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := make([]entity.MembershipPlan, 0)

	for rows.Next() {
		var plan entity.MembershipPlan

		err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.DurationMonths,
			&plan.Price,
			&plan.Description,
			&plan.IsDeleted,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}

func (r *MembershipPostgresRepository) Update(
	ctx context.Context,
	plan *entity.MembershipPlan,
) error {
	query := `
		UPDATE membership_plans
		SET
			name = $1,
			duration_months = $2,
			price = $3,
			description = $4,
			updated_at = $5
		WHERE id = $6
		  AND is_deleted = false
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		plan.Name,
		plan.DurationMonths,
		plan.Price,
		plan.Description,
		plan.UpdatedAt,
		plan.ID,
	)

	return err
}

func (r *MembershipPostgresRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := `
		UPDATE membership_plans
		SET
			is_deleted = true,
			updated_at = NOW()
		WHERE id = $1
		  AND is_deleted = false
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
