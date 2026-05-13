package grpc

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"membership-service/internal/domain/entity"
	"membership-service/internal/infrastructure/cache"
	membershippb "membership-service/proto"
)

func mapMembershipPlanToProto(plan *entity.MembershipPlan) *membershippb.MembershipPlan {
	if plan == nil {
		return nil
	}

	return &membershippb.MembershipPlan{
		Id:             plan.ID.String(),
		Name:           plan.Name,
		DurationMonths: int32(plan.DurationMonths),
		Price:          plan.Price,
		Description:    plan.Description,
		IsDeleted:      plan.IsDeleted,
		CreatedAt:      timestamppb.New(plan.CreatedAt),
		UpdatedAt:      timestamppb.New(plan.UpdatedAt),
	}
}

func mapSubscriptionToProto(subscription *entity.Subscription) *membershippb.Subscription {
	if subscription == nil {
		return nil
	}

	return &membershippb.Subscription{
		Id:               subscription.ID.String(),
		UserId:           subscription.UserID.String(),
		MembershipPlanId: subscription.MembershipPlanID.String(),
		StartDate:        timestamppb.New(subscription.StartDate),
		EndDate:          timestamppb.New(subscription.EndDate),
		IsActive:         subscription.IsActive,
		Status:           mapSubscriptionStatusToProto(subscription.Status),
		CreatedAt:        timestamppb.New(subscription.CreatedAt),
		UpdatedAt:        timestamppb.New(subscription.UpdatedAt),
	}
}

func mapAttendanceToProto(attendance *entity.Attendance) *membershippb.Attendance {
	if attendance == nil {
		return nil
	}

	return &membershippb.Attendance{
		Id:          attendance.ID.String(),
		UserId:      attendance.UserID.String(),
		CheckInTime: timestamppb.New(attendance.CheckInTime),
		Date:        formatDate(attendance.Date),
		CreatedAt:   timestamppb.New(attendance.CreatedAt),
	}
}

func mapAttendanceStatsToProto(stats *cache.AttendanceStats) *membershippb.MonthlyAttendanceStatsResponse {
	if stats == nil {
		return nil
	}

	return &membershippb.MonthlyAttendanceStatsResponse{
		UserId:          stats.UserID,
		VisitsThisMonth: stats.VisitsThisMonth,
		VisitsLastMonth: stats.VisitsLastMonth,
		TotalVisits:     stats.TotalVisits,
	}
}

func mapSubscriptionStatusToProto(
	status entity.SubscriptionStatus,
) membershippb.SubscriptionStatus {
	switch status {
	case entity.SubscriptionStatusActive:
		return membershippb.SubscriptionStatus_SUBSCRIPTION_STATUS_ACTIVE
	case entity.SubscriptionStatusExpired:
		return membershippb.SubscriptionStatus_SUBSCRIPTION_STATUS_EXPIRED
	case entity.SubscriptionStatusCancelled:
		return membershippb.SubscriptionStatus_SUBSCRIPTION_STATUS_CANCELLED
	default:
		return membershippb.SubscriptionStatus_SUBSCRIPTION_STATUS_UNSPECIFIED
	}
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}
