package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	membershippb "membership-service/proto"
)

func (h *Handler) CreateMembershipPlan(
	ctx context.Context,
	req *membershippb.CreateMembershipPlanRequest,
) (*membershippb.MembershipPlanResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "membership plan name is required")
	}

	plan, err := h.membershipUC.CreatePlan(
		ctx,
		req.Name,
		int(req.DurationMonths),
		req.Price,
		req.Description,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &membershippb.MembershipPlanResponse{
		Plan: mapMembershipPlanToProto(plan),
	}, nil
}

func (h *Handler) GetMembershipPlan(
	ctx context.Context,
	req *membershippb.GetMembershipPlanRequest,
) (*membershippb.MembershipPlanResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid membership plan id")
	}

	plan, err := h.membershipUC.GetPlan(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &membershippb.MembershipPlanResponse{
		Plan: mapMembershipPlanToProto(plan),
	}, nil
}

func (h *Handler) ListMembershipPlans(
	ctx context.Context,
	req *membershippb.ListMembershipPlansRequest,
) (*membershippb.ListMembershipPlansResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)

	plans, err := h.membershipUC.ListPlans(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responsePlans := make([]*membershippb.MembershipPlan, 0, len(plans))

	for i := range plans {
		plan := plans[i]
		responsePlans = append(responsePlans, mapMembershipPlanToProto(&plan))
	}

	return &membershippb.ListMembershipPlansResponse{
		Plans:  responsePlans,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

func (h *Handler) UpdateMembershipPlan(
	ctx context.Context,
	req *membershippb.UpdateMembershipPlanRequest,
) (*membershippb.MembershipPlanResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid membership plan id")
	}

	plan, err := h.membershipUC.UpdatePlan(
		ctx,
		id,
		req.Name,
		int(req.DurationMonths),
		req.Price,
		req.Description,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &membershippb.MembershipPlanResponse{
		Plan: mapMembershipPlanToProto(plan),
	}, nil
}

func (h *Handler) DeleteMembershipPlan(
	ctx context.Context,
	req *membershippb.DeleteMembershipPlanRequest,
) (*membershippb.DeleteMembershipPlanResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid membership plan id")
	}

	if err := h.membershipUC.DeletePlan(ctx, id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &membershippb.DeleteMembershipPlanResponse{
		Success: true,
	}, nil
}
