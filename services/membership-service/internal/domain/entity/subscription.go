package entity

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusExpired   SubscriptionStatus = "EXPIRED"
	SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
)

type Subscription struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	MembershipPlanID uuid.UUID
	StartDate        time.Time
	EndDate          time.Time
	IsActive         bool
	Status           SubscriptionStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewSubscription(
	userID uuid.UUID,
	membershipPlanID uuid.UUID,
	durationMonths int,
) *Subscription {
	now := time.Now().UTC()

	return &Subscription{
		ID:               uuid.New(),
		UserID:           userID,
		MembershipPlanID: membershipPlanID,
		StartDate:        now,
		EndDate:          now.AddDate(0, durationMonths, 0),
		IsActive:         true,
		Status:           SubscriptionStatusActive,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func (s *Subscription) Extend(durationMonths int) {
	now := time.Now().UTC()

	if s.EndDate.Before(now) {
		s.StartDate = now
		s.EndDate = now.AddDate(0, durationMonths, 0)
	} else {
		s.EndDate = s.EndDate.AddDate(0, durationMonths, 0)
	}

	s.IsActive = true
	s.Status = SubscriptionStatusActive
	s.UpdatedAt = now
}

func (s *Subscription) Cancel() {
	s.IsActive = false
	s.Status = SubscriptionStatusCancelled
	s.UpdatedAt = time.Now().UTC()
}

func (s *Subscription) Expire() {
	s.IsActive = false
	s.Status = SubscriptionStatusExpired
	s.UpdatedAt = time.Now().UTC()
}

func (s Subscription) CanAccess(now time.Time) bool {
	return s.IsActive &&
		s.Status == SubscriptionStatusActive &&
		s.EndDate.After(now)
}

func (s Subscription) IsExpired(now time.Time) bool {
	return s.EndDate.Before(now) || s.Status == SubscriptionStatusExpired
}
