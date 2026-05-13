package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type AttendanceStats struct {
	UserID          string `json:"user_id"`
	VisitsThisMonth int64  `json:"visits_this_month"`
	VisitsLastMonth int64  `json:"visits_last_month"`
	TotalVisits     int64  `json:"total_visits"`
}

type AttendanceCache interface {
	SetUserStats(ctx context.Context, userID string, stats *AttendanceStats) error
	GetUserStats(ctx context.Context, userID string) (*AttendanceStats, error)
	DeleteUserStats(ctx context.Context, userID string) error
}

type RedisAttendanceCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisAttendanceCache(client *redis.Client) *RedisAttendanceCache {
	return &RedisAttendanceCache{
		client: client,
		ttl:    10 * time.Minute,
	}
}

func (c *RedisAttendanceCache) SetUserStats(
	ctx context.Context,
	userID string,
	stats *AttendanceStats,
) error {
	key := "attendance:stats:" + userID

	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *RedisAttendanceCache) GetUserStats(
	ctx context.Context,
	userID string,
) (*AttendanceStats, error) {
	key := "attendance:stats:" + userID

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var stats AttendanceStats

	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (c *RedisAttendanceCache) DeleteUserStats(ctx context.Context, userID string) error {
	key := "attendance:stats:" + userID

	return c.client.Del(ctx, key).Err()
}
