package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	membershippb "membership-service/proto"
)

func (h *Handler) GetGlobalStats(
	ctx context.Context,
	req *membershippb.GetGlobalStatsRequest,
) (*membershippb.GlobalStatsResponse, error) {
	stats, err := h.statsUC.GetGlobalStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &membershippb.GlobalStatsResponse{
		ActiveSubscriptionsCount: stats.ActiveSubscriptionsCount,
		VisitsThisMonth:          stats.VisitsThisMonth,
		TotalVisits:              stats.TotalVisits,
	}, nil
}
