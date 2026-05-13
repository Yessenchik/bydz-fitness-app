package grpc

import (
	"membership-service/internal/usecase"
	membershippb "membership-service/proto"
)

type Handler struct {
	membershippb.UnimplementedMembershipAccessServiceServer

	membershipUC   *usecase.MembershipUsecase
	subscriptionUC *usecase.SubscriptionUsecase
	attendanceUC   *usecase.AttendanceUsecase
	statsUC        *usecase.StatsUsecase
}

func NewHandler(
	membershipUC *usecase.MembershipUsecase,
	subscriptionUC *usecase.SubscriptionUsecase,
	attendanceUC *usecase.AttendanceUsecase,
	statsUC *usecase.StatsUsecase,
) *Handler {
	return &Handler{
		membershipUC:   membershipUC,
		subscriptionUC: subscriptionUC,
		attendanceUC:   attendanceUC,
		statsUC:        statsUC,
	}
}
