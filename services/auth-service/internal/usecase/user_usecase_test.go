package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/usecase"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/mocks"
)

func TestGetProfile_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	cache := mocks.NewMockCacheRepository(ctrl)
	publisher := mocks.NewMockMessagePublisher(ctrl)
	log := zaptest.NewLogger(t)

	uc := usecase.NewUserUsecase(repo, cache, publisher, log)

	userID := uuid.New()
	cached := &domain.User{ID: userID, Email: "john@example.com", FirstName: "John"}

	// Кэш возвращает данные — репозиторий НЕ должен вызываться
	cache.EXPECT().
		GetUserProfile(gomock.Any(), userID).
		Return(cached, nil)

	repo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(0)

	user, err := uc.GetProfile(context.Background(), userID)

	require.NoError(t, err)
	assert.Equal(t, cached.Email, user.Email)
}

func TestGetProfile_CacheMiss_FallbackToRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	cache := mocks.NewMockCacheRepository(ctrl)
	publisher := mocks.NewMockMessagePublisher(ctrl)
	log := zaptest.NewLogger(t)

	uc := usecase.NewUserUsecase(repo, cache, publisher, log)

	userID := uuid.New()
	dbUser := &domain.User{
		ID: userID, Email: "john@example.com",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	// Cache miss
	cache.EXPECT().
		GetUserProfile(gomock.Any(), userID).
		Return(nil, nil)

	// Fallback to repo
	repo.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(dbUser, nil)

	// Запись в кэш после получения из репозитория
	cache.EXPECT().
		SetUserProfile(gomock.Any(), dbUser).
		Return(nil)

	user, err := uc.GetProfile(context.Background(), userID)

	require.NoError(t, err)
	assert.Equal(t, dbUser.Email, user.Email)
}

func TestAssignRole_PublishesEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	cache := mocks.NewMockCacheRepository(ctrl)
	publisher := mocks.NewMockMessagePublisher(ctrl)
	log := zaptest.NewLogger(t)

	uc := usecase.NewUserUsecase(repo, cache, publisher, log)

	userID := uuid.New()
	user := &domain.User{
		ID:        userID,
		Email:     "john@example.com",
		Role:      domain.RoleClient,
		UpdatedAt: time.Now(),
	}

	repo.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(user, nil)

	repo.EXPECT().
		Update(gomock.Any(), gomock.AssignableToTypeOf(&domain.User{})).
		DoAndReturn(func(_ context.Context, u *domain.User) error {
			assert.Equal(t, domain.RoleTrainer, u.Role)
			return nil
		})

	cache.EXPECT().
		InvalidateUserProfile(gomock.Any(), userID).
		Return(nil)

	// Главная проверка: событие NATS публикуется с правильным subject
	publisher.EXPECT().
		Publish(gomock.Any(), domain.EventUserRoleChanged, gomock.Any()).
		DoAndReturn(func(_ context.Context, subject string, payload any) error {
			p, ok := payload.(domain.UserRoleChangedPayload)
			require.True(t, ok)
			assert.Equal(t, domain.RoleClient, p.OldRole)
			assert.Equal(t, domain.RoleTrainer, p.NewRole)
			return nil
		})

	updated, err := uc.AssignRole(context.Background(), userID, domain.RoleTrainer)

	require.NoError(t, err)
	assert.Equal(t, domain.RoleTrainer, updated.Role)
}

func TestDeleteUser_SoftDeleteAndEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	cache := mocks.NewMockCacheRepository(ctrl)
	publisher := mocks.NewMockMessagePublisher(ctrl)
	log := zaptest.NewLogger(t)

	uc := usecase.NewUserUsecase(repo, cache, publisher, log)

	userID := uuid.New()

	repo.EXPECT().
		SoftDelete(gomock.Any(), userID).
		Return(nil)

	cache.EXPECT().
		InvalidateUserProfile(gomock.Any(), userID).
		Return(nil)

	publisher.EXPECT().
		Publish(gomock.Any(), domain.EventUserDeleted, gomock.Any()).
		Return(nil)

	err := uc.DeleteUser(context.Background(), userID)
	require.NoError(t, err)
}
