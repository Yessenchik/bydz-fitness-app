package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	membershippb "membership-service/proto"
)

func (h *Handler) RegisterAttendance(
	ctx context.Context,
	req *membershippb.RegisterAttendanceRequest,
) (*membershippb.AttendanceResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	attendance, err := h.attendanceUC.RegisterAttendance(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &membershippb.AttendanceResponse{
		Attendance: mapAttendanceToProto(attendance),
	}, nil
}

func (h *Handler) GetUserAttendanceHistory(
	ctx context.Context,
	req *membershippb.GetUserAttendanceHistoryRequest,
) (*membershippb.GetUserAttendanceHistoryResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	attendances, err := h.attendanceUC.GetUserHistory(
		ctx,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseAttendances := make([]*membershippb.Attendance, 0, len(attendances))

	for i := range attendances {
		attendance := attendances[i]
		responseAttendances = append(responseAttendances, mapAttendanceToProto(&attendance))
	}

	return &membershippb.GetUserAttendanceHistoryResponse{
		Attendances: responseAttendances,
		Limit:       req.Limit,
		Offset:      req.Offset,
	}, nil
}

func (h *Handler) GetMonthlyAttendanceStats(
	ctx context.Context,
	req *membershippb.GetMonthlyAttendanceStatsRequest,
) (*membershippb.MonthlyAttendanceStatsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	stats, err := h.attendanceUC.GetMonthlyStats(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return mapAttendanceStatsToProto(stats), nil
}
