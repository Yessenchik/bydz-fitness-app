package grpc

import (
	membershipv1 "github.com/Yessenchik/bydz-fitness-app/services/membership-service/gen/proto/membership/v1"
	"github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/domain"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func planToProto(p *domain.Plan) *membershipv1.Plan {
	if p == nil {
		return nil
	}

	return &membershipv1.Plan{
		PlanId:       p.ID,
		Name:         p.Name,
		DurationDays: p.DurationDays,
		PriceKzt:     p.PriceKZT,
		IsActive:     p.IsActive,
		CreatedAt:    timestamppb.New(p.CreatedAt),
		UpdatedAt:    timestamppb.New(p.UpdatedAt),
	}
}

func subscriptionToProto(s *domain.Subscription) *membershipv1.Subscription {
	if s == nil {
		return nil
	}

	return &membershipv1.Subscription{
		SubscriptionId: s.ID,
		UserId:         s.UserID,
		PlanId:         s.PlanID,
		Status:         statusToProto(s.Status),
		StartsAt:       timestamppb.New(s.StartsAt),
		ExpiresAt:      timestamppb.New(s.ExpiresAt),
		CreatedAt:      timestamppb.New(s.CreatedAt),
		UpdatedAt:      timestamppb.New(s.UpdatedAt),
	}
}

func statusToProto(status domain.SubscriptionStatus) membershipv1.SubscriptionStatus {
	switch status {
	case domain.StatusPending:
		return membershipv1.SubscriptionStatus_SUBSCRIPTION_STATUS_PENDING
	case domain.StatusActive:
		return membershipv1.SubscriptionStatus_SUBSCRIPTION_STATUS_ACTIVE
	case domain.StatusCancelled:
		return membershipv1.SubscriptionStatus_SUBSCRIPTION_STATUS_CANCELLED
	case domain.StatusExpired:
		return membershipv1.SubscriptionStatus_SUBSCRIPTION_STATUS_EXPIRED
	default:
		return membershipv1.SubscriptionStatus_SUBSCRIPTION_STATUS_UNSPECIFIED
	}
}
