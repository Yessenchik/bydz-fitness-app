package usecase

import (
	"context"
	"time"

	"membership-service/internal/domain/repository"
)

type GlobalStats struct {
	ActiveSubscriptionsCount int64
	VisitsThisMonth          int64
	TotalVisits              int64
}

type StatsUsecase struct {
	subscriptionRepo repository.SubscriptionRepository
	attendanceRepo   repository.AttendanceRepository
}

func NewStatsUsecase(
	subscriptionRepo repository.SubscriptionRepository,
	attendanceRepo repository.AttendanceRepository,
) *StatsUsecase {
	return &StatsUsecase{
		subscriptionRepo: subscriptionRepo,
		attendanceRepo:   attendanceRepo,
	}
}

func (u *StatsUsecase) GetGlobalStats(ctx context.Context) (*GlobalStats, error) {
	activeSubscriptionsCount, err := u.subscriptionRepo.CountActive(ctx)
	if err != nil {
		return nil, err
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

	visitsThisMonth, err := u.attendanceRepo.CountAllByPeriod(
		ctx,
		thisMonthStart,
		nextMonthStart,
	)
	if err != nil {
		return nil, err
	}

	totalVisits, err := u.attendanceRepo.CountTotal(ctx)
	if err != nil {
		return nil, err
	}

	return &GlobalStats{
		ActiveSubscriptionsCount: activeSubscriptionsCount,
		VisitsThisMonth:          visitsThisMonth,
		TotalVisits:              totalVisits,
	}, nil
}
