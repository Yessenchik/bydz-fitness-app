package usecase

import (
	"context"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
	"github.com/google/uuid"
)

// AuthUsecase — контракт для auth-операций (Login, Register, токены)
type AuthUsecase interface {
	Register(ctx context.Context, input RegisterInput) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.Token, *domain.User, error)
	Logout(ctx context.Context, accessToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error)
	VerifyEmail(ctx context.Context, token string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ConfirmPasswordReset(ctx context.Context, token, newPassword string) error
}

// UserUsecase — контракт для CRUD и ролевых операций
type UserUsecase interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateProfile(ctx context.Context, input UpdateProfileInput) (*domain.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	AssignRole(ctx context.Context, userID uuid.UUID, role domain.Role) (*domain.User, error)
	ListUsers(ctx context.Context, filter domain.ListUsersFilter) ([]*domain.User, int, error)
}

// RegisterInput — входные данные для регистрации
type RegisterInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
}

// UpdateProfileInput — входные данные для обновления профиля
type UpdateProfileInput struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	Phone     string
}
