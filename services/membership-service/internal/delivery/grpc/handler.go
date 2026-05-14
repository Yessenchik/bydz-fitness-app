package grpc

import (
	"context"

	membershipv1 "github.com/Yessenchik/bydz-fitness-app/services/membership-service/gen/proto/membership/v1"
	"github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	membershipv1.UnimplementedMembershipServiceServer

	uc *usecase.MembershipUsecase
}

func NewHandler(uc *usecase.MembershipUsecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) CreatePlan(ctx context.Context, req *membershipv1.CreatePlanRequest) (*membershipv1.Plan, error) {
	plan, err := h.uc.CreatePlan(ctx, req.Name, req.DurationDays, req.PriceKzt)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return planToProto(plan), nil
}

func (h *Handler) GetPlan(ctx context.Context, req *membershipv1.GetPlanRequest) (*membershipv1.Plan, error) {
	plan, err := h.uc.GetPlan(ctx, req.PlanId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return planToProto(plan), nil
}

func (h *Handler) ListPlans(ctx context.Context, req *membershipv1.ListPlansRequest) (*membershipv1.ListPlansResponse, error) {
	plans, err := h.uc.ListPlans(ctx, req.OnlyActive)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	items := make([]*membershipv1.Plan, 0, len(plans))
	for _, p := range plans {
		items = append(items, planToProto(p))
	}

	return &membershipv1.ListPlansResponse{
		Plans: items,
	}, nil
}

func (h *Handler) CreateSubscription(ctx context.Context, req *membershipv1.CreateSubscriptionRequest) (*membershipv1.Subscription, error) {
	sub, err := h.uc.CreateSubscription(ctx, req.UserId, req.PlanId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return subscriptionToProto(sub), nil
}

func (h *Handler) GetSubscription(ctx context.Context, req *membershipv1.GetSubscriptionRequest) (*membershipv1.Subscription, error) {
	sub, err := h.uc.GetSubscription(ctx, req.SubscriptionId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return subscriptionToProto(sub), nil
}

func (h *Handler) GetActiveSubscription(ctx context.Context, req *membershipv1.GetActiveSubscriptionRequest) (*membershipv1.Subscription, error) {
	sub, err := h.uc.GetActiveSubscription(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return subscriptionToProto(sub), nil
}

func (h *Handler) ListUserSubscriptions(ctx context.Context, req *membershipv1.ListUserSubscriptionsRequest) (*membershipv1.ListUserSubscriptionsResponse, error) {
	subs, err := h.uc.ListUserSubscriptions(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	items := make([]*membershipv1.Subscription, 0, len(subs))
	for _, s := range subs {
		items = append(items, subscriptionToProto(s))
	}

	return &membershipv1.ListUserSubscriptionsResponse{
		Subscriptions: items,
	}, nil
}

func (h *Handler) CancelSubscription(ctx context.Context, req *membershipv1.CancelSubscriptionRequest) (*membershipv1.Subscription, error) {
	sub, err := h.uc.CancelSubscription(ctx, req.SubscriptionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return subscriptionToProto(sub), nil
}

func (h *Handler) ValidateMembership(ctx context.Context, req *membershipv1.ValidateMembershipRequest) (*membershipv1.ValidateMembershipResponse, error) {
	ok, reason, sub := h.uc.ValidateMembership(ctx, req.UserId)

	return &membershipv1.ValidateMembershipResponse{
		IsValid:      ok,
		Reason:       reason,
		Subscription: subscriptionToProto(sub),
	}, nil
}

func (h *Handler) UpdatePlan(ctx context.Context, req *membershipv1.UpdatePlanRequest) (*membershipv1.Plan, error) {
	return nil, status.Error(codes.Unimplemented, "UpdatePlan not implemented yet")
}

func (h *Handler) DeletePlan(ctx context.Context, req *membershipv1.DeletePlanRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "DeletePlan not implemented yet")
}

func (h *Handler) ExtendSubscription(ctx context.Context, req *membershipv1.ExtendSubscriptionRequest) (*membershipv1.Subscription, error) {
	return nil, status.Error(codes.Unimplemented, "ExtendSubscription not implemented yet")
}

func (h *Handler) ExpireSubscription(ctx context.Context, req *membershipv1.ExpireSubscriptionRequest) (*membershipv1.Subscription, error) {
	return nil, status.Error(codes.Unimplemented, "ExpireSubscription not implemented yet")
}
