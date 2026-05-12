package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	membershippb "membership-service/proto"
)

func (h *Handler) CreateSubscription(
	ctx context.Context,
	req *membershippb.CreateSubscriptionRequest,
) (*membershippb.SubscriptionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	planID, err := uuid.Parse(req.MembershipPlanId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid membership plan id")
	}

	subscription, err := h.subscriptionUC.CreateSubscription(
		ctx,
		userID,
		planID,
		req.UserEmail,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &membershippb.SubscriptionResponse{
		Subscription: mapSubscriptionToProto(subscription),
	}, nil
}

func (h *Handler) GetUserSubscription(
	ctx context.Context,
	req *membershippb.GetUserSubscriptionRequest,
) (*membershippb.SubscriptionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	subscription, err := h.subscriptionUC.GetUserSubscription(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &membershippb.SubscriptionResponse{
		Subscription: mapSubscriptionToProto(subscription),
	}, nil
}

func (h *Handler) ExtendSubscription(
	ctx context.Context,
	req *membershippb.ExtendSubscriptionRequest,
) (*membershippb.SubscriptionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	planID, err := uuid.Parse(req.MembershipPlanId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid membership plan id")
	}

	subscription, err := h.subscriptionUC.ExtendSubscription(
		ctx,
		userID,
		planID,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &membershippb.SubscriptionResponse{
		Subscription: mapSubscriptionToProto(subscription),
	}, nil
}

func (h *Handler) CancelSubscription(
	ctx context.Context,
	req *membershippb.CancelSubscriptionRequest,
) (*membershippb.CancelSubscriptionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if err := h.subscriptionUC.CancelSubscription(ctx, userID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &membershippb.CancelSubscriptionResponse{
		Success: true,
	}, nil
}

func (h *Handler) CheckAccess(
	ctx context.Context,
	req *membershippb.CheckAccessRequest,
) (*membershippb.CheckAccessResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	allowed, reason, err := h.subscriptionUC.CheckAccess(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &membershippb.CheckAccessResponse{
		Allowed: allowed,
		Reason:  reason,
	}, nil
}
