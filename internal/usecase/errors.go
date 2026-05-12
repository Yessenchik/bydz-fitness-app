package usecase

import "errors"

var (
	ErrInvalidMembershipName = errors.New("invalid membership name")
	ErrInvalidDuration       = errors.New("invalid membership duration")
	ErrInvalidPrice          = errors.New("invalid price")

	ErrMembershipPlanNotFound = errors.New("membership plan not found")
	ErrSubscriptionNotFound   = errors.New("subscription not found")
	ErrAccessDenied           = errors.New("access denied")

	ErrInvalidUserID = errors.New("invalid user id")
	ErrInvalidPlanID = errors.New("invalid membership plan id")
)
