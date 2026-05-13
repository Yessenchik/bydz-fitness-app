package usecase

import (
	"context"

	"github.com/google/uuid"

	"membership-service/internal/domain/entity"
	"membership-service/internal/domain/repository"
)

type MembershipUsecase struct {
	repo repository.MembershipRepository
}

func NewMembershipUsecase(repo repository.MembershipRepository) *MembershipUsecase {
	return &MembershipUsecase{
		repo: repo,
	}
}

func (u *MembershipUsecase) CreatePlan(
	ctx context.Context,
	name string,
	durationMonths int,
	price float64,
	description string,
) (*entity.MembershipPlan, error) {
	plan := entity.NewMembershipPlan(
		name,
		durationMonths,
		price,
		description,
	)

	if plan.Name == "" {
		return nil, ErrInvalidMembershipName
	}

	if !plan.IsValidDuration() {
		return nil, ErrInvalidDuration
	}

	if plan.Price < 0 {
		return nil, ErrInvalidPrice
	}

	if err := u.repo.Create(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

func (u *MembershipUsecase) GetPlan(
	ctx context.Context,
	id uuid.UUID,
) (*entity.MembershipPlan, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *MembershipUsecase) ListPlans(
	ctx context.Context,
	limit int,
	offset int,
) ([]entity.MembershipPlan, error) {
	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	return u.repo.List(ctx, limit, offset)
}

func (u *MembershipUsecase) UpdatePlan(
	ctx context.Context,
	id uuid.UUID,
	name string,
	durationMonths int,
	price float64,
	description string,
) (*entity.MembershipPlan, error) {
	plan, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	plan.Update(
		name,
		durationMonths,
		price,
		description,
	)

	if plan.Name == "" {
		return nil, ErrInvalidMembershipName
	}

	if !plan.IsValidDuration() {
		return nil, ErrInvalidDuration
	}

	if plan.Price < 0 {
		return nil, ErrInvalidPrice
	}

	if err := u.repo.Update(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

func (u *MembershipUsecase) DeletePlan(
	ctx context.Context,
	id uuid.UUID,
) error {
	return u.repo.Delete(ctx, id)
}
