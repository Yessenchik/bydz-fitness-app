package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"membership-service/internal/domain/entity"
	"membership-service/internal/domain/repository"
	"membership-service/internal/infrastructure/cache"
	"membership-service/internal/infrastructure/events"
)

type AttendanceUsecase struct {
	attendanceRepo repository.AttendanceRepository
	subscriptionUC *SubscriptionUsecase
	cache          cache.AttendanceCache
	producer       events.Producer
}

func NewAttendanceUsecase(
	attendanceRepo repository.AttendanceRepository,
	subscriptionUC *SubscriptionUsecase,
	cache cache.AttendanceCache,
	producer events.Producer,
) *AttendanceUsecase {
	return &AttendanceUsecase{
		attendanceRepo: attendanceRepo,
		subscriptionUC: subscriptionUC,
		cache:          cache,
		producer:       producer,
	}
}

func (u *AttendanceUsecase) RegisterAttendance(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.Attendance, error) {
	allowed, reason, err := u.subscriptionUC.CheckAccess(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, ErrAccessDenied
	}

	attendance := entity.NewAttendance(userID)

	if err := u.attendanceRepo.Create(ctx, attendance); err != nil {
		return nil, err
	}

	_ = u.cache.DeleteUserStats(ctx, userID.String())

	_ = u.producer.Publish(ctx, "attendance.created", map[string]any{
		"attendance_id": attendance.ID.String(),
		"user_id":       userID.String(),
		"check_in_time": attendance.CheckInTime,
		"reason":        reason,
	})

	return attendance, nil
}

func (u *AttendanceUsecase) GetUserHistory(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	offset int,
) ([]entity.Attendance, error) {
	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	return u.attendanceRepo.GetByUserID(ctx, userID, limit, offset)
}

func (u *AttendanceUsecase) GetMonthlyStats(
	ctx context.Context,
	userID uuid.UUID,
) (*cache.AttendanceStats, error) {
	cached, err := u.cache.GetUserStats(ctx, userID.String())
	if err == nil && cached != nil {
		return cached, nil
	}

	now := time.Now().UTC()

	thisMonthStart := time.Date(
		now.Year(),
		now.Month(),
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	nextMonthStart := thisMonthStart.AddDate(0, 1, 0)

	lastMonthStart := thisMonthStart.AddDate(0, -1, 0)
	lastMonthEnd := thisMonthStart

	visitsThisMonth, err := u.attendanceRepo.CountByUserAndPeriod(
		ctx,
		userID,
		thisMonthStart,
		nextMonthStart,
	)
	if err != nil {
		return nil, err
	}

	visitsLastMonth, err := u.attendanceRepo.CountByUserAndPeriod(
		ctx,
		userID,
		lastMonthStart,
		lastMonthEnd,
	)
	if err != nil {
		return nil, err
	}

	totalVisits, err := u.attendanceRepo.CountTotalByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats := &cache.AttendanceStats{
		UserID:          userID.String(),
		VisitsThisMonth: visitsThisMonth,
		VisitsLastMonth: visitsLastMonth,
		TotalVisits:     totalVisits,
	}

	_ = u.cache.SetUserStats(ctx, userID.String(), stats)

	return stats, nil
}
