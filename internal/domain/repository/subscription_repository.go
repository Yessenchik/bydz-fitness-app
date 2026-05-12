package repository

import (
	"context"

	"github.com/google/uuid"
	"membership-service/internal/domain/entity"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *entity.Subscription) error

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Subscription, error)

	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Subscription, error)

	GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*entity.Subscription, error)

	Update(ctx context.Context, subscription *entity.Subscription) error

	Cancel(ctx context.Context, userID uuid.UUID) error

	ExpireOldSubscriptions(ctx context.Context) error

	CountActive(ctx context.Context) (int64, error)
}
