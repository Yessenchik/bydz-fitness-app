package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type SubscriptionCache interface {
	SetStatus(ctx context.Context, userID string, status string, endDate time.Time) error
	GetStatus(ctx context.Context, userID string) (string, error)
	DeleteStatus(ctx context.Context, userID string) error
}

type RedisSubscriptionCache struct {
	client *redis.Client
}

func NewRedisSubscriptionCache(client *redis.Client) *RedisSubscriptionCache {
	return &RedisSubscriptionCache{
		client: client,
	}
}

func (c *RedisSubscriptionCache) SetStatus(
	ctx context.Context,
	userID string,
	status string,
	endDate time.Time,
) error {
	key := "subscription:status:" + userID

	ttl := time.Until(endDate)
	if ttl <= 0 {
		ttl = time.Minute
	}

	return c.client.Set(ctx, key, status, ttl).Err()
}

func (c *RedisSubscriptionCache) GetStatus(ctx context.Context, userID string) (string, error) {
	key := "subscription:status:" + userID

	return c.client.Get(ctx, key).Result()
}

func (c *RedisSubscriptionCache) DeleteStatus(ctx context.Context, userID string) error {
	key := "subscription:status:" + userID

	return c.client.Del(ctx, key).Err()
}
