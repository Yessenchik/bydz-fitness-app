package grpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/Yessenchik/bydz-fitness-app/user-auth-service/gen/userauth/v1"
	deliveryGRPC "github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/delivery/grpc"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/usecase"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/mocks"
)

const bufSize = 1024 * 1024

// newTestServer поднимает gRPC-сервер на in-memory буфере (bufconn).
// Это позволяет тестировать полный стек (handler → маршалинг → proto) без сетевого сокета.
func newTestServer(
	t *testing.T,
	authUC usecase.AuthUsecase,
	userUC usecase.UserUsecase,
) pb.UserAuthServiceClient {
	t.Helper()

	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	pb.RegisterUserAuthServiceServer(srv, deliveryGRPC.NewHandler(authUC, userUC, zaptest.NewLogger(t)))

	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Logf("bufconn server stopped: %v", err)
		}
	}()

	t.Cleanup(func() { srv.Stop() })

	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return pb.NewUserAuthServiceClient(conn)
}

func TestHandler_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authUC := mocks.NewMockAuthUsecase(ctrl)
	userUC := mocks.NewMockUserUsecase(ctrl)

	expectedUser := &domain.User{
		Email:     "john@example.com",
		FirstName: "John",
	}

	authUC.EXPECT().
		Register(gomock.Any(), usecase.RegisterInput{
			Email:     "john@example.com",
			Password:  "StrongPass1",
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "+7-999-000-00-00",
		}).
		Return(expectedUser, nil)

	client := newTestServer(t, authUC, userUC)

	resp, err := client.Register(context.Background(), &pb.RegisterRequest{
		Email:     "john@example.com",
		Password:  "StrongPass1",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "+7-999-000-00-00",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Contains(t, resp.Message, "john@example.com")
}

func TestHandler_Register_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authUC := mocks.NewMockAuthUsecase(ctrl)
	userUC := mocks.NewMockUserUsecase(ctrl)

	authUC.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, domain.ErrUserAlreadyExists)

	client := newTestServer(t, authUC, userUC)

	_, err := client.Register(context.Background(), &pb.RegisterRequest{
		Email: "existing@example.com", Password: "StrongPass1",
	})

	require.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestHandler_GetProfile_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := newTestServer(
		t,
		mocks.NewMockAuthUsecase(ctrl),
		mocks.NewMockUserUsecase(ctrl),
	)

	_, err := client.GetProfile(context.Background(), &pb.GetProfileRequest{
		UserId: "not-a-valid-uuid",
	})

	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}
