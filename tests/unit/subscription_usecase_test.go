package unit

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"membership-service/internal/domain/entity"
	"membership-service/internal/usecase"
)

func TestCreateSubscription_ShouldCreateActiveSubscription(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	planID := uuid.New()

	planRepo := &fakeMembershipRepo{
		plan: &entity.MembershipPlan{
			ID:             planID,
			Name:           "1 Month Plan",
			DurationMonths: 1,
			Price:          15000,
			Description:    "Basic gym plan",
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		},
	}

	subRepo := &fakeSubscriptionRepo{}
	subCache := &fakeSubscriptionCache{}
	producer := &fakeProducer{}
	emailService := &fakeEmailService{}

	uc := usecase.NewSubscriptionUsecase(
		subRepo,
		planRepo,
		subCache,
		producer,
		emailService,
	)

	subscription, err := uc.CreateSubscription(
		ctx,
		userID,
		planID,
		"test@example.com",
	)

	require.NoError(t, err)
	require.NotNil(t, subscription)

	assert.Equal(t, userID, subscription.UserID)
	assert.Equal(t, planID, subscription.MembershipPlanID)
	assert.Equal(t, entity.SubscriptionStatusActive, subscription.Status)
	assert.True(t, subscription.IsActive)
	assert.True(t, subscription.EndDate.After(time.Now().UTC()))
	assert.Equal(t, "ACTIVE", subCache.status)
	assert.True(t, emailService.sent)
	assert.Contains(t, producer.subjects, "subscription.created")
}

func TestCheckAccess_WhenSubscriptionActive_ShouldAllowAccess(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	planID := uuid.New()

	subRepo := &fakeSubscriptionRepo{
		subscription: &entity.Subscription{
			ID:               uuid.New(),
			UserID:           userID,
			MembershipPlanID: planID,
			StartDate:        time.Now().UTC(),
			EndDate:          time.Now().UTC().AddDate(0, 1, 0),
			IsActive:         true,
			Status:           entity.SubscriptionStatusActive,
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		},
	}

	uc := usecase.NewSubscriptionUsecase(
		subRepo,
		&fakeMembershipRepo{},
		&fakeSubscriptionCache{},
		&fakeProducer{},
		&fakeEmailService{},
	)

	allowed, reason, err := uc.CheckAccess(ctx, userID)

	require.NoError(t, err)
	assert.True(t, allowed)
	assert.Equal(t, "access granted", reason)
}

func TestCheckAccess_WhenNoSubscription_ShouldDenyAccess(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()

	producer := &fakeProducer{}

	uc := usecase.NewSubscriptionUsecase(
		&fakeSubscriptionRepo{},
		&fakeMembershipRepo{},
		&fakeSubscriptionCache{},
		producer,
		&fakeEmailService{},
	)

	allowed, reason, err := uc.CheckAccess(ctx, userID)

	require.NoError(t, err)
	assert.False(t, allowed)
	assert.Equal(t, "active subscription not found", reason)
	assert.Contains(t, producer.subjects, "access.denied")
}

func TestCancelSubscription_ShouldCancelAndClearCache(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()

	subRepo := &fakeSubscriptionRepo{
		subscription: entity.NewSubscription(userID, uuid.New(), 1),
	}

	subCache := &fakeSubscriptionCache{
		status: "ACTIVE",
	}

	producer := &fakeProducer{}

	uc := usecase.NewSubscriptionUsecase(
		subRepo,
		&fakeMembershipRepo{},
		subCache,
		producer,
		&fakeEmailService{},
	)

	err := uc.CancelSubscription(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, "", subCache.status)
	assert.Equal(t, entity.SubscriptionStatusCancelled, subRepo.subscription.Status)
	assert.Contains(t, producer.subjects, "subscription.cancelled")
}
