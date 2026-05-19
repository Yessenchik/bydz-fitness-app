package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/pkg/metrics"
)

type userUsecase struct {
	userRepo  domain.UserRepository
	cache     domain.CacheRepository
	publisher domain.MessagePublisher
	logger    *zap.Logger
}

func NewUserUsecase(
	userRepo domain.UserRepository,
	cache domain.CacheRepository,
	publisher domain.MessagePublisher,
	logger *zap.Logger,
) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		cache:     cache,
		publisher: publisher,
		logger:    logger,
	}
}

// ─────────────────────────────────────────
//  GetProfile — с кэшированием
// ─────────────────────────────────────────

func (uc *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserUsecase.GetProfile")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", userID.String()))

	// Cache-aside: сначала Redis
	if cached, err := uc.cache.GetUserProfile(ctx, userID); err == nil && cached != nil {
		metrics.CacheHits.Inc()
		return cached, nil
	}
	metrics.CacheMisses.Inc()

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Пишем в кэш (ошибка кэша не критична)
	if err := uc.cache.SetUserProfile(ctx, user); err != nil {
		uc.logger.Warn("failed to cache user profile",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	return user, nil
}

// ─────────────────────────────────────────
//  UpdateProfile
// ─────────────────────────────────────────

func (uc *userUsecase) UpdateProfile(ctx context.Context, input UpdateProfileInput) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserUsecase.UpdateProfile")
	defer span.End()

	user, err := uc.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Phone = input.Phone
	user.UpdatedAt = time.Now().UTC()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("update profile: %w", err)
	}

	_ = uc.cache.InvalidateUserProfile(ctx, user.ID)

	uc.logger.Info("profile updated", zap.String("user_id", user.ID.String()))
	return user, nil
}

// ─────────────────────────────────────────
//  DeleteUser
// ─────────────────────────────────────────

func (uc *userUsecase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "UserUsecase.DeleteUser")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", userID.String()))

	if err := uc.userRepo.SoftDelete(ctx, userID); err != nil {
		span.RecordError(err)
		return fmt.Errorf("delete user: %w", err)
	}

	_ = uc.cache.InvalidateUserProfile(ctx, userID)

	// Событие для Membership / Training сервисов
	if err := uc.publisher.Publish(ctx, domain.EventUserDeleted, domain.UserDeletedPayload{
		UserID:     userID,
		OccurredAt: time.Now().UTC(),
	}); err != nil {
		uc.logger.Warn("failed to publish user.deleted event",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	uc.logger.Info("user deleted", zap.String("user_id", userID.String()))
	return nil
}

// ─────────────────────────────────────────
//  ChangePassword
// ─────────────────────────────────────────

func (uc *userUsecase) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	ctx, span := tracer.Start(ctx, "UserUsecase.ChangePassword")
	defer span.End()

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Проверяем старый пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return domain.ErrInvalidCredentials
	}

	if err := validatePassword(newPassword); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return fmt.Errorf("change password: hash: %w", err)
	}

	user.PasswordHash = string(hash)
	user.UpdatedAt = time.Now().UTC()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		span.RecordError(err)
		return fmt.Errorf("change password: update: %w", err)
	}

	_ = uc.cache.InvalidateUserProfile(ctx, user.ID)
	uc.logger.Info("password changed", zap.String("user_id", userID.String()))
	return nil
}

// ─────────────────────────────────────────
//  AssignRole
// ─────────────────────────────────────────

func (uc *userUsecase) AssignRole(ctx context.Context, userID uuid.UUID, role domain.Role) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserUsecase.AssignRole")
	defer span.End()

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	oldRole := user.Role
	user.Role = role
	user.UpdatedAt = time.Now().UTC()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("assign role: %w", err)
	}

	_ = uc.cache.InvalidateUserProfile(ctx, user.ID)

	if err := uc.publisher.Publish(ctx, domain.EventUserRoleChanged, domain.UserRoleChangedPayload{
		UserID:     userID,
		OldRole:    oldRole,
		NewRole:    role,
		OccurredAt: time.Now().UTC(),
	}); err != nil {
		uc.logger.Warn("failed to publish user.role_changed event",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	uc.logger.Info("role assigned",
		zap.String("user_id", userID.String()),
		zap.String("old_role", string(oldRole)),
		zap.String("new_role", string(role)),
	)
	return user, nil
}

// ─────────────────────────────────────────
//  ListUsers
// ─────────────────────────────────────────

func (uc *userUsecase) ListUsers(ctx context.Context, filter domain.ListUsersFilter) ([]*domain.User, int, error) {
	ctx, span := tracer.Start(ctx, "UserUsecase.ListUsers")
	defer span.End()

	users, total, err := uc.userRepo.List(ctx, filter)
	if err != nil {
		span.RecordError(err)
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	return users, total, nil
}
