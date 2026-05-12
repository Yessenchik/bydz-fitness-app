package repository

import (
	"context"

	"github.com/google/uuid"
	"membership-service/internal/domain/entity"
)

type MembershipRepository interface {
	Create(ctx context.Context, plan *entity.MembershipPlan) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.MembershipPlan, error)

	List(ctx context.Context, limit int, offset int) ([]entity.MembershipPlan, error)

	Update(ctx context.Context, plan *entity.MembershipPlan) error

	Delete(ctx context.Context, id uuid.UUID) error
}
