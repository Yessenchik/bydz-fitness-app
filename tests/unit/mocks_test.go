package unit

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"membership-service/internal/domain/entity"
	"membership-service/internal/infrastructure/cache"
)

type fakeMembershipRepo struct {
	plan *entity.MembershipPlan
}

func (r *fakeMembershipRepo) Create(ctx context.Context, plan *entity.MembershipPlan) error {
	r.plan = plan
	return nil
}

func (r *fakeMembershipRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.MembershipPlan, error) {
	if r.plan == nil {
		return nil, errors.New("plan not found")
	}

	return r.plan, nil
}

func (r *fakeMembershipRepo) List(ctx context.Context, limit int, offset int) ([]entity.MembershipPlan, error) {
	if r.plan == nil {
		return []entity.MembershipPlan{}, nil
	}

	return []entity.MembershipPlan{*r.plan}, nil
}

func (r *fakeMembershipRepo) Update(ctx context.Context, plan *entity.MembershipPlan) error {
	r.plan = plan
	return nil
}

func (r *fakeMembershipRepo) Delete(ctx context.Context, id uuid.UUID) error {
	r.plan = nil
	return nil
}

type fakeSubscriptionRepo struct {
	subscription *entity.Subscription
}

func (r *fakeSubscriptionRepo) Create(ctx context.Context, subscription *entity.Subscription) error {
	r.subscription = subscription
	return nil
}

func (r *fakeSubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	if r.subscription == nil {
		return nil, errors.New("subscription not found")
	}

	return r.subscription, nil
}

func (r *fakeSubscriptionRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Subscription, error) {
	if r.subscription == nil {
		return nil, errors.New("subscription not found")
	}

	return r.subscription, nil
}

func (r *fakeSubscriptionRepo) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*entity.Subscription, error) {
	if r.subscription == nil {
		return nil, errors.New("active subscription not found")
	}

	if r.subscription.UserID != userID {
		return nil, errors.New("active subscription not found")
	}

	if !r.subscription.CanAccess(time.Now().UTC()) {
		return nil, errors.New("active subscription not found")
	}

	return r.subscription, nil
}

func (r *fakeSubscriptionRepo) Update(ctx context.Context, subscription *entity.Subscription) error {
	r.subscription = subscription
	return nil
}

func (r *fakeSubscriptionRepo) Cancel(ctx context.Context, userID uuid.UUID) error {
	if r.subscription != nil && r.subscription.UserID == userID {
		r.subscription.Cancel()
	}

	return nil
}

func (r *fakeSubscriptionRepo) ExpireOldSubscriptions(ctx context.Context) error {
	if r.subscription != nil {
		r.subscription.Expire()
	}

	return nil
}

func (r *fakeSubscriptionRepo) CountActive(ctx context.Context) (int64, error) {
	if r.subscription == nil {
		return 0, nil
	}

	if r.subscription.CanAccess(time.Now().UTC()) {
		return 1, nil
	}

	return 0, nil
}

type fakeAttendanceRepo struct {
	attendances []entity.Attendance
}

func (r *fakeAttendanceRepo) Create(ctx context.Context, attendance *entity.Attendance) error {
	r.attendances = append(r.attendances, *attendance)
	return nil
}

func (r *fakeAttendanceRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Attendance, error) {
	for i := range r.attendances {
		if r.attendances[i].ID == id {
			return &r.attendances[i], nil
		}
	}

	return nil, errors.New("attendance not found")
}

func (r *fakeAttendanceRepo) GetByUserID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]entity.Attendance, error) {
	result := make([]entity.Attendance, 0)

	for _, attendance := range r.attendances {
		if attendance.UserID == userID {
			result = append(result, attendance)
		}
	}

	return result, nil
}

func (r *fakeAttendanceRepo) CountByUserAndPeriod(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) (int64, error) {
	var count int64

	for _, attendance := range r.attendances {
		if attendance.UserID == userID &&
			!attendance.CheckInTime.Before(from) &&
			attendance.CheckInTime.Before(to) {
			count++
		}
	}

	return count, nil
}

func (r *fakeAttendanceRepo) CountAllByPeriod(ctx context.Context, from time.Time, to time.Time) (int64, error) {
	var count int64

	for _, attendance := range r.attendances {
		if !attendance.CheckInTime.Before(from) && attendance.CheckInTime.Before(to) {
			count++
		}
	}

	return count, nil
}

func (r *fakeAttendanceRepo) CountTotalByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64

	for _, attendance := range r.attendances {
		if attendance.UserID == userID {
			count++
		}
	}

	return count, nil
}

func (r *fakeAttendanceRepo) CountTotal(ctx context.Context) (int64, error) {
	return int64(len(r.attendances)), nil
}

type fakeSubscriptionCache struct {
	status string
}

func (c *fakeSubscriptionCache) SetStatus(ctx context.Context, userID string, status string, endDate time.Time) error {
	c.status = status
	return nil
}

func (c *fakeSubscriptionCache) GetStatus(ctx context.Context, userID string) (string, error) {
	if c.status == "" {
		return "", errors.New("cache miss")
	}

	return c.status, nil
}

func (c *fakeSubscriptionCache) DeleteStatus(ctx context.Context, userID string) error {
	c.status = ""
	return nil
}

type fakeAttendanceCache struct {
	stats *cache.AttendanceStats
}

func (c *fakeAttendanceCache) SetUserStats(ctx context.Context, userID string, stats *cache.AttendanceStats) error {
	c.stats = stats
	return nil
}

func (c *fakeAttendanceCache) GetUserStats(ctx context.Context, userID string) (*cache.AttendanceStats, error) {
	if c.stats == nil {
		return nil, errors.New("cache miss")
	}

	return c.stats, nil
}

func (c *fakeAttendanceCache) DeleteUserStats(ctx context.Context, userID string) error {
	c.stats = nil
	return nil
}

type fakeProducer struct {
	subjects []string
}

func (p *fakeProducer) Publish(ctx context.Context, subject string, payload any) error {
	p.subjects = append(p.subjects, subject)
	return nil
}

type fakeEmailService struct {
	sent bool
}

func (s *fakeEmailService) SendPurchaseConfirmation(ctx context.Context, to string, planName string, endDate time.Time) error {
	s.sent = true
	return nil
}

func (s *fakeEmailService) SendExpirationWarning(ctx context.Context, to string, endDate time.Time) error {
	s.sent = true
	return nil
}

func (s *fakeEmailService) SendExpired(ctx context.Context, to string) error {
	s.sent = true
	return nil
}
