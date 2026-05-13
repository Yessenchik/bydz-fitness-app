package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	deliverygrpc "membership-service/internal/delivery/grpc"
	"membership-service/internal/repository/postgres"
	"membership-service/internal/usecase"
	membershippb "membership-service/proto"
)

func TestGRPCHandler_CreateMembershipPlan(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	membershipRepo := postgres.NewMembershipPostgresRepository(db)
	membershipUC := usecase.NewMembershipUsecase(membershipRepo)

	handler := deliverygrpc.NewHandler(
		membershipUC,
		nil,
		nil,
		nil,
	)

	response, err := handler.CreateMembershipPlan(
		ctx,
		&membershippb.CreateMembershipPlanRequest{
			Name:           "Test Handler Plan",
			DurationMonths: 1,
			Price:          15000,
			Description:    "Created from gRPC handler test",
		},
	)

	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Plan)

	assert.Equal(t, "Test Handler Plan", response.Plan.Name)
	assert.Equal(t, int32(1), response.Plan.DurationMonths)
	assert.Equal(t, float64(15000), response.Plan.Price)
}
