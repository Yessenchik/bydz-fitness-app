package domain

import "time"

type SubscriptionStatus string

const (
	StatusPending   SubscriptionStatus = "pending"
	StatusActive    SubscriptionStatus = "active"
	StatusCancelled SubscriptionStatus = "cancelled"
	StatusExpired   SubscriptionStatus = "expired"
)

type Subscription struct {
	ID          string
	UserID      string
	PlanID      string
	Status      SubscriptionStatus
	StartsAt    time.Time
	ExpiresAt   time.Time
	CancelledAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
