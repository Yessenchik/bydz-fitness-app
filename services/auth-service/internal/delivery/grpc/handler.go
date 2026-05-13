package grpc

import (
	"context"

	authv1 "github.com/Yessenchik/bydz-fitness-app/services/auth-service/gen/proto/auth/v1"
	"github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	authv1.UnimplementedAuthServiceServer

	authUC *usecase.AuthUsecase
}

func NewHandler(authUC *usecase.AuthUsecase) *Handler {
	return &Handler{authUC: authUC}
}

func (h *Handler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, err := h.authUC.Register(ctx, req.Email, req.Password, req.FirstName, req.LastName, req.Phone)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.RegisterResponse{
		UserId:  user.ID,
		Message: "user registered",
	}, nil
}

func (h *Handler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	tokenPair, user, err := h.authUC.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authv1.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: &authv1.UserProfile{
			UserId:    user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
		},
	}, nil
}

func (h *Handler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	tokenPair, err := h.authUC.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authv1.RefreshTokenResponse{
		AccessToken: tokenPair.AccessToken,
		ExpiresIn:   tokenPair.ExpiresIn,
	}, nil
}

func (h *Handler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*emptypb.Empty, error) {
	err := h.authUC.Logout(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) GetProfile(ctx context.Context, req *authv1.GetProfileRequest) (*authv1.UserProfile, error) {
	user, err := h.authUC.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &authv1.UserProfile{
		UserId:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
	}, nil
}
