package usecase

import (
	"context"

	"github.com/google/uuid"

	"membership-service/internal/domain/entity"
	"membership-service/internal/domain/repository"
	"membership-service/internal/infrastructure/cache"
	"membership-service/internal/infrastructure/email"
	"membership-service/internal/infrastructure/events"
)

type SubscriptionUsecase struct {
	subRepo  repository.SubscriptionRepository
	planRepo repository.MembershipRepository
	cache    cache.SubscriptionCache
	producer events.Producer
	email    email.Service
}

func NewSubscriptionUsecase(
	subRepo repository.SubscriptionRepository,
	planRepo repository.MembershipRepository,
	cache cache.SubscriptionCache,
	producer events.Producer,
	email email.Service,
) *SubscriptionUsecase {
	return &SubscriptionUsecase{
		subRepo:  subRepo,
		planRepo: planRepo,
		cache:    cache,
		producer: producer,
		email:    email,
	}
}

func (u *SubscriptionUsecase) CreateSubscription(
	ctx context.Context,
	userID uuid.UUID,
	planID uuid.UUID,
	userEmail string,
) (*entity.Subscription, error) {
	plan, err := u.planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	subscription := entity.NewSubscription(
		userID,
		plan.ID,
		plan.DurationMonths,
	)

	if err := u.subRepo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	_ = u.cache.SetStatus(
		ctx,
		userID.String(),
		string(subscription.Status),
		subscription.EndDate,
	)

	_ = u.producer.Publish(ctx, "subscription.created", map[string]any{
		"subscription_id": subscription.ID.String(),
		"user_id":         userID.String(),
		"plan_id":         plan.ID.String(),
		"end_date":        subscription.EndDate,
	})

	if userEmail != "" {
		_ = u.email.SendPurchaseConfirmation(
			ctx,
			userEmail,
			plan.Name,
			subscription.EndDate,
		)
	}

	return subscription, nil
}

func (u *SubscriptionUsecase) GetUserSubscription(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.Subscription, error) {
	return u.subRepo.GetByUserID(ctx, userID)
}

func (u *SubscriptionUsecase) ExtendSubscription(
	ctx context.Context,
	userID uuid.UUID,
	planID uuid.UUID,
) (*entity.Subscription, error) {
	plan, err := u.planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	subscription, err := u.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	subscription.MembershipPlanID = plan.ID
	subscription.Extend(plan.DurationMonths)

	if err := u.subRepo.Update(ctx, subscription); err != nil {
		return nil, err
	}

	_ = u.cache.SetStatus(
		ctx,
		userID.String(),
		string(subscription.Status),
		subscription.EndDate,
	)

	_ = u.producer.Publish(ctx, "subscription.extended", map[string]any{
		"subscription_id": subscription.ID.String(),
		"user_id":         userID.String(),
		"plan_id":         plan.ID.String(),
		"new_end_date":    subscription.EndDate,
	})

	return subscription, nil
}

func (u *SubscriptionUsecase) CancelSubscription(
	ctx context.Context,
	userID uuid.UUID,
) error {
	if err := u.subRepo.Cancel(ctx, userID); err != nil {
		return err
	}

	_ = u.cache.DeleteStatus(ctx, userID.String())

	_ = u.producer.Publish(ctx, "subscription.cancelled", map[string]any{
		"user_id": userID.String(),
	})

	return nil
}

func (u *SubscriptionUsecase) CheckAccess(
	ctx context.Context,
	userID uuid.UUID,
) (bool, string, error) {
	cachedStatus, err := u.cache.GetStatus(ctx, userID.String())
	if err == nil && cachedStatus == string(entity.SubscriptionStatusActive) {
		return true, "access granted from cache", nil
	}

	subscription, err := u.subRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		_ = u.producer.Publish(ctx, "access.denied", map[string]any{
			"user_id": userID.String(),
			"reason":  "active subscription not found",
		})

		return false, "active subscription not found", nil
	}

	if !subscription.CanAccess(nowUTC()) {
		_ = u.cache.DeleteStatus(ctx, userID.String())

		_ = u.producer.Publish(ctx, "access.denied", map[string]any{
			"user_id": userID.String(),
			"reason":  "subscription expired or cancelled",
		})

		return false, "subscription expired or cancelled", nil
	}

	_ = u.cache.SetStatus(
		ctx,
		userID.String(),
		string(subscription.Status),
		subscription.EndDate,
	)

	return true, "access granted", nil
}

func (u *SubscriptionUsecase) ExpireOldSubscriptions(
	ctx context.Context,
) error {
	if err := u.subRepo.ExpireOldSubscriptions(ctx); err != nil {
		return err
	}

	_ = u.producer.Publish(ctx, "subscription.expired", map[string]any{
		"message": "old subscriptions expired",
	})

	return nil
}
