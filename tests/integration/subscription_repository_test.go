package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"membership-service/internal/domain/entity"
	"membership-service/internal/repository/postgres"
)

func TestSubscriptionRepository_CreateAndGetActive(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	membershipRepo := postgres.NewMembershipPostgresRepository(db)
	subscriptionRepo := postgres.NewSubscriptionPostgresRepository(db)

	plan := entity.NewMembershipPlan(
		"Test 1 Month Plan",
		1,
		15000,
		"Test plan for subscription repository",
	)

	err := membershipRepo.Create(ctx, plan)
	require.NoError(t, err)

	userID := uuid.New()

	subscription := entity.NewSubscription(
		userID,
		plan.ID,
		plan.DurationMonths,
	)

	err = subscriptionRepo.Create(ctx, subscription)
	require.NoError(t, err)

	found, err := subscriptionRepo.GetActiveByUserID(ctx, userID)
	require.NoError(t, err)

	assert.Equal(t, subscription.ID, found.ID)
	assert.Equal(t, userID, found.UserID)
	assert.Equal(t, entity.SubscriptionStatusActive, found.Status)
	assert.True(t, found.EndDate.After(time.Now().UTC()))
}

func TestSubscriptionRepository_Cancel(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	membershipRepo := postgres.NewMembershipPostgresRepository(db)
	subscriptionRepo := postgres.NewSubscriptionPostgresRepository(db)

	plan := entity.NewMembershipPlan(
		"Test 6 Month Plan",
		6,
		70000,
		"Test plan for cancel",
	)

	err := membershipRepo.Create(ctx, plan)
	require.NoError(t, err)

	userID := uuid.New()

	subscription := entity.NewSubscription(
		userID,
		plan.ID,
		plan.DurationMonths,
	)

	err = subscriptionRepo.Create(ctx, subscription)
	require.NoError(t, err)

	err = subscriptionRepo.Cancel(ctx, userID)
	require.NoError(t, err)

	cancelled, err := subscriptionRepo.GetByUserID(ctx, userID)
	require.NoError(t, err)

	assert.False(t, cancelled.IsActive)
	assert.Equal(t, entity.SubscriptionStatusCancelled, cancelled.Status)
}
