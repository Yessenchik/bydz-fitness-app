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

func TestRegisterAttendance_WhenAccessAllowed_ShouldCreateAttendance(t *testing.T) {
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

	subscriptionUC := usecase.NewSubscriptionUsecase(
		subRepo,
		&fakeMembershipRepo{},
		&fakeSubscriptionCache{},
		&fakeProducer{},
		&fakeEmailService{},
	)

	attendanceRepo := &fakeAttendanceRepo{}
	attendanceCache := &fakeAttendanceCache{}
	producer := &fakeProducer{}

	attendanceUC := usecase.NewAttendanceUsecase(
		attendanceRepo,
		subscriptionUC,
		attendanceCache,
		producer,
	)

	attendance, err := attendanceUC.RegisterAttendance(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, attendance)

	assert.Equal(t, userID, attendance.UserID)
	assert.Len(t, attendanceRepo.attendances, 1)
	assert.Contains(t, producer.subjects, "attendance.created")
}

func TestRegisterAttendance_WhenAccessDenied_ShouldReturnError(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()

	subscriptionUC := usecase.NewSubscriptionUsecase(
		&fakeSubscriptionRepo{},
		&fakeMembershipRepo{},
		&fakeSubscriptionCache{},
		&fakeProducer{},
		&fakeEmailService{},
	)

	attendanceUC := usecase.NewAttendanceUsecase(
		&fakeAttendanceRepo{},
		subscriptionUC,
		&fakeAttendanceCache{},
		&fakeProducer{},
	)

	attendance, err := attendanceUC.RegisterAttendance(ctx, userID)

	require.Error(t, err)
	assert.Nil(t, attendance)
	assert.Equal(t, usecase.ErrAccessDenied, err)
}

func TestGetMonthlyStats_ShouldReturnStats(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()

	attendanceRepo := &fakeAttendanceRepo{
		attendances: []entity.Attendance{
			{
				ID:          uuid.New(),
				UserID:      userID,
				CheckInTime: time.Now().UTC(),
				Date:        time.Now().UTC(),
				CreatedAt:   time.Now().UTC(),
			},
			{
				ID:          uuid.New(),
				UserID:      userID,
				CheckInTime: time.Now().UTC(),
				Date:        time.Now().UTC(),
				CreatedAt:   time.Now().UTC(),
			},
		},
	}

	subscriptionUC := usecase.NewSubscriptionUsecase(
		&fakeSubscriptionRepo{},
		&fakeMembershipRepo{},
		&fakeSubscriptionCache{},
		&fakeProducer{},
		&fakeEmailService{},
	)

	attendanceUC := usecase.NewAttendanceUsecase(
		attendanceRepo,
		subscriptionUC,
		&fakeAttendanceCache{},
		&fakeProducer{},
	)

	stats, err := attendanceUC.GetMonthlyStats(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, userID.String(), stats.UserID)
	assert.Equal(t, int64(2), stats.VisitsThisMonth)
	assert.Equal(t, int64(0), stats.VisitsLastMonth)
	assert.Equal(t, int64(2), stats.TotalVisits)
}
