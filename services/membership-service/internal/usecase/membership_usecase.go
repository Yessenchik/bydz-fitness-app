package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/domain"
	natsPublisher "github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/infrastructure/nats"
)

type PlanRepo interface {
	Create(ctx context.Context, p *domain.Plan) error
	GetByID(ctx context.Context, id string) (*domain.Plan, error)
	List(ctx context.Context, onlyActive bool) ([]*domain.Plan, error)
}

type SubscriptionRepo interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetByID(ctx context.Context, id string) (*domain.Subscription, error)
	GetActiveByUserID(ctx context.Context, userID string) (*domain.Subscription, error)
	ListByUserID(ctx context.Context, userID string) ([]*domain.Subscription, error)
	Cancel(ctx context.Context, id string) (*domain.Subscription, error)
}

type MembershipUsecase struct {
	plans PlanRepo
	subs  SubscriptionRepo
	nats  *natsPublisher.Publisher
}

func NewMembershipUsecase(
	plans PlanRepo,
	subs SubscriptionRepo,
	nats *natsPublisher.Publisher,
) *MembershipUsecase {
	return &MembershipUsecase{
		plans: plans,
		subs:  subs,
		nats:  nats,
	}
}

func (u *MembershipUsecase) CreatePlan(ctx context.Context, name string, durationDays int32, priceKZT int64) (*domain.Plan, error) {
	if name == "" {
		return nil, errors.New("plan name is required")
	}
	if durationDays <= 0 {
		return nil, errors.New("duration_days must be positive")
	}
	if priceKZT < 0 {
		return nil, errors.New("price_kzt cannot be negative")
	}

	plan := &domain.Plan{
		Name:         name,
		DurationDays: durationDays,
		PriceKZT:     priceKZT,
		IsActive:     true,
	}

	if err := u.plans.Create(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

func (u *MembershipUsecase) GetPlan(ctx context.Context, planID string) (*domain.Plan, error) {
	return u.plans.GetByID(ctx, planID)
}

func (u *MembershipUsecase) ListPlans(ctx context.Context, onlyActive bool) ([]*domain.Plan, error) {
	return u.plans.List(ctx, onlyActive)
}

func (u *MembershipUsecase) CreateSubscription(ctx context.Context, userID string, planID string) (*domain.Subscription, error) {
	plan, err := u.plans.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	if !plan.IsActive {
		return nil, errors.New("plan is not active")
	}

	_, err = u.subs.GetActiveByUserID(ctx, userID)
	if err == nil {
		return nil, errors.New("user already has active subscription")
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	now := time.Now()

	sub := &domain.Subscription{
		UserID:    userID,
		PlanID:    planID,
		Status:    domain.StatusActive,
		StartsAt:  now,
		ExpiresAt: now.AddDate(0, 0, int(plan.DurationDays)),
	}

	if err := u.subs.Create(ctx, sub); err != nil {
		return nil, err
	}

	_ = u.nats.Publish(ctx, "membership.subscription_created", map[string]any{
		"subscription_id": sub.ID,
		"user_id":         sub.UserID,
		"plan_id":         sub.PlanID,
		"status":          string(sub.Status),
		"expires_at":      sub.ExpiresAt,
	})

	return sub, nil
}

func (u *MembershipUsecase) GetSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	return u.subs.GetByID(ctx, id)
}

func (u *MembershipUsecase) GetActiveSubscription(ctx context.Context, userID string) (*domain.Subscription, error) {
	return u.subs.GetActiveByUserID(ctx, userID)
}

func (u *MembershipUsecase) ListUserSubscriptions(ctx context.Context, userID string) ([]*domain.Subscription, error) {
	return u.subs.ListByUserID(ctx, userID)
}

func (u *MembershipUsecase) CancelSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	sub, err := u.subs.Cancel(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = u.nats.Publish(ctx, "membership.subscription_cancelled", map[string]any{
		"subscription_id": sub.ID,
		"user_id":         sub.UserID,
		"plan_id":         sub.PlanID,
		"status":          string(sub.Status),
	})

	return sub, nil
}

func (u *MembershipUsecase) ValidateMembership(ctx context.Context, userID string) (bool, string, *domain.Subscription) {
	sub, err := u.subs.GetActiveByUserID(ctx, userID)
	if err != nil {
		return false, "no active subscription", nil
	}

	if sub.Status != domain.StatusActive {
		return false, "subscription is not active", sub
	}

	if time.Now().After(sub.ExpiresAt) {
		return false, "subscription expired", sub
	}

	return true, "membership is valid", sub
}
