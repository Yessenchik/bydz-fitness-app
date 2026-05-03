package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
)

const (
	profileTTL      = 15 * time.Minute
	prefixProfile   = "user:profile:"
	prefixBlacklist = "token:blacklist:"
	prefixRefresh   = "token:refresh:"
)

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) domain.CacheRepository {
	return &cacheRepository{client: client}
}

// ─────────────────────────────────────────
//  User Profile Cache
// ─────────────────────────────────────────

func (r *cacheRepository) SetUserProfile(ctx context.Context, user *domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("cache set profile marshal: %w", err)
	}
	key := prefixProfile + user.ID.String()
	return r.client.Set(ctx, key, data, profileTTL).Err()
}

func (r *cacheRepository) GetUserProfile(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	key := prefixProfile + id.String()
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss — не ошибка
		}
		return nil, fmt.Errorf("cache get profile: %w", err)
	}
	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("cache get profile unmarshal: %w", err)
	}
	return &user, nil
}

func (r *cacheRepository) InvalidateUserProfile(ctx context.Context, id uuid.UUID) error {
	return r.client.Del(ctx, prefixProfile+id.String()).Err()
}

// ─────────────────────────────────────────
//  JWT Blacklist
// ─────────────────────────────────────────

func (r *cacheRepository) BlacklistToken(ctx context.Context, tokenID string, ttlSeconds int64) error {
	key := prefixBlacklist + tokenID
	return r.client.Set(ctx, key, "1", time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *cacheRepository) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := prefixBlacklist + tokenID
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("check blacklist: %w", err)
	}
	return result > 0, nil
}

// ─────────────────────────────────────────
//  Refresh Token Store
// ─────────────────────────────────────────

func (r *cacheRepository) SetRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string, ttlSeconds int64) error {
	key := prefixRefresh + userID.String()
	return r.client.Set(ctx, key, tokenID, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *cacheRepository) GetRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	key := prefixRefresh + userID.String()
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrTokenInvalid
		}
		return "", fmt.Errorf("get refresh token: %w", err)
	}
	return val, nil
}

func (r *cacheRepository) DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error {
	return r.client.Del(ctx, prefixRefresh+userID.String()).Err()
}
