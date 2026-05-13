package usecase

import (
	"context"
	"errors"

	"github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/domain"
	jwtManager "github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/infrastructure/jwt"
	natsPublisher "github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/infrastructure/nats"
	redisRepo "github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/infrastructure/redis"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Create(ctx context.Context, u *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

type AuthUsecase struct {
	users UserRepo
	jwt   *jwtManager.JWTManager
	redis *redisRepo.Client
	nats  *natsPublisher.Publisher
}

func NewAuthUsecase(
	users UserRepo,
	jwt *jwtManager.JWTManager,
	redis *redisRepo.Client,
	nats *natsPublisher.Publisher,
) *AuthUsecase {
	return &AuthUsecase{
		users: users,
		jwt:   jwt,
		redis: redis,
		nats:  nats,
	}
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

func (u *AuthUsecase) Register(ctx context.Context, email, password, firstName, lastName, phone string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
		FirstName:    firstName,
		LastName:     lastName,
		Phone:        phone,
		Role:         domain.RoleClient,
		IsVerified:   false,
	}

	if err := u.users.Create(ctx, user); err != nil {
		return nil, err
	}

	_ = u.nats.Publish("user.registered", map[string]string{
		"user_id":    user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
	})

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (*TokenPair, *domain.User, error) {
	user, err := u.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	tokenPair, err := u.generateAndSaveTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return tokenPair, user, nil
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	userID, err := u.jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	savedToken, err := u.redis.GetRefreshToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	if savedToken != refreshToken {
		return nil, errors.New("refresh token revoked or replaced")
	}

	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return u.generateAndSaveTokens(ctx, user)
}

func (u *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	userID, err := u.jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	return u.redis.DeleteRefreshToken(ctx, userID)
}

func (u *AuthUsecase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return u.users.GetByID(ctx, userID)
}

func (u *AuthUsecase) generateAndSaveTokens(ctx context.Context, user *domain.User) (*TokenPair, error) {
	accessToken, err := u.jwt.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	if err := u.redis.SaveRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	}, nil
}
