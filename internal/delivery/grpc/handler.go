package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/Yessenchik/bydz-fitness-app/user-auth-service/gen/userauth/v1"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/usecase"
)

// Handler реализует pb.UserAuthServiceServer.
// Единственная ответственность — трансляция proto ↔ domain + вызов usecase.
type Handler struct {
	pb.UnimplementedUserAuthServiceServer
	authUC usecase.AuthUsecase
	userUC usecase.UserUsecase
	logger *zap.Logger
}

func NewHandler(
	authUC usecase.AuthUsecase,
	userUC usecase.UserUsecase,
	logger *zap.Logger,
) *Handler {
	return &Handler{authUC: authUC, userUC: userUC, logger: logger}
}

// ─────────────────────────────────────────
//  Auth RPCs
// ─────────────────────────────────────────

func (h *Handler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := h.authUC.Register(ctx, usecase.RegisterInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	})
	if err != nil {
		return nil, domainErrToStatus(err)
	}
	return &pb.RegisterResponse{
		UserId:  user.ID.String(),
		Message: "Письмо с подтверждением отправлено на " + user.Email,
	}, nil
}

func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	tokenPair, user, err := h.authUC.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, domainErrToStatus(err)
	}
	return &pb.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         domainUserToProto(user),
	}, nil
}

func (h *Handler) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	if err := h.authUC.Logout(ctx, req.AccessToken); err != nil {
		return nil, domainErrToStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	tokenPair, err := h.authUC.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, domainErrToStatus(err)
	}
	return &pb.RefreshTokenResponse{
		AccessToken: tokenPair.AccessToken,
		ExpiresIn:   tokenPair.ExpiresIn,
	}, nil
}

func (h *Handler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*emptypb.Empty, error) {
	if err := h.authUC.VerifyEmail(ctx, req.VerificationToken); err != nil {
		return nil, domainErrToStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) RequestPasswordReset(ctx context.Context, req *pb.RequestPasswordResetRequest) (*emptypb.Empty, error) {
	if err := h.authUC.RequestPasswordReset(ctx, req.Email); err != nil {
		return nil, domainErrToStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) ConfirmPasswordReset(ctx context.Context, req *pb.ConfirmPasswordResetRequest) (*emptypb.Empty, error) {
	if err := h.authUC.ConfirmPasswordReset(ctx, req.ResetToken, req.NewPassword); err != nil {
		return nil, domainErrToStatus(err)
	}
	return &emptypb.Empty{}, nil
}

// ─────────────────────────────────────────
//  User RPCs
// ─────────────────────────────────────────

func (h *Handler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.UserProfile, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}
	user, err := h.userUC.GetProfile(ctx, id)
	if err != nil {
		return nil, domainErrToStatus(err)
	}
	return domainUserToProto(user), nil
}

func (h *Handler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserProfile, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}
	user, err := h.userUC.UpdateProfile(ctx, usecase.UpdateProfileInput{
		UserID:    id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	})
	if err != nil {
		return nil, domainErrToStatus(err)
	}
	return domainUserToProto(user), nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}
	if err := h.userUC.DeleteUser(ctx, id); err != nil {
		return nil, domainErrToStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}
	if err := h.userUC.ChangePassword(ctx, id, req.OldPassword, req.NewPassword); err != nil {
		return nil, domainErrToStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.UserProfile, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}
	role := protoRoleToDomain(req.Role)
	user, err := h.userUC.AssignRole(ctx, id, role)
	if err != nil {
		return nil, domainErrToStatus(err)
	}
	return domainUserToProto(user), nil
}

func (h *Handler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	filter := domain.ListUsersFilter{
		Page:     page,
		PageSize: pageSize,
	}
	if req.FilterRole != pb.Role_ROLE_UNSPECIFIED {
		r := protoRoleToDomain(req.FilterRole)
		filter.Role = &r
	}

	users, total, err := h.userUC.ListUsers(ctx, filter)
	if err != nil {
		return nil, domainErrToStatus(err)
	}

	protoUsers := make([]*pb.UserProfile, len(users))
	for i, u := range users {
		protoUsers[i] = domainUserToProto(u)
	}

	totalPages := int32(total) / req.PageSize
	if int32(total)%req.PageSize != 0 {
		totalPages++
	}

	return &pb.ListUsersResponse{
		Users:      protoUsers,
		TotalCount: int32(total),
		Page:       req.Page,
		TotalPages: totalPages,
	}, nil
}

// ─────────────────────────────────────────
//  Mapping helpers
// ─────────────────────────────────────────

func domainUserToProto(u *domain.User) *pb.UserProfile {
	return &pb.UserProfile{
		UserId:     u.ID.String(),
		Email:      u.Email,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Phone:      u.Phone,
		Role:       domainRoleToProto(u.Role),
		IsVerified: u.IsVerified,
		CreatedAt:  timestamppb.New(u.CreatedAt),
		UpdatedAt:  timestamppb.New(u.UpdatedAt),
	}
}

func domainRoleToProto(r domain.Role) pb.Role {
	switch r {
	case domain.RoleClient:
		return pb.Role_ROLE_CLIENT
	case domain.RoleTrainer:
		return pb.Role_ROLE_TRAINER
	case domain.RoleAdmin:
		return pb.Role_ROLE_ADMIN
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}

func protoRoleToDomain(r pb.Role) domain.Role {
	switch r {
	case pb.Role_ROLE_CLIENT:
		return domain.RoleClient
	case pb.Role_ROLE_TRAINER:
		return domain.RoleTrainer
	case pb.Role_ROLE_ADMIN:
		return domain.RoleAdmin
	default:
		return domain.RoleClient
	}
}

// domainErrToStatus — единая точка трансляции domain-ошибок в gRPC-статусы
func domainErrToStatus(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrTokenInvalid),
		errors.Is(err, domain.ErrTokenExpired):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrEmailNotVerified):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, domain.ErrWeakPassword):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
