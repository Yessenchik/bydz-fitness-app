package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(
	addr string,
) *Client {

	rdb := redis.NewClient(
		&redis.Options{
			Addr: addr,
		},
	)

	return &Client{
		rdb: rdb,
	}
}

func (c *Client) SaveRefreshToken(
	ctx context.Context,
	userID string,
	token string,
) error {

	return c.rdb.Set(
		ctx,
		"refresh:"+userID,
		token,
		7*24*time.Hour,
	).Err()
}

func (c *Client) GetRefreshToken(
	ctx context.Context,
	userID string,
) (string, error) {

	return c.rdb.Get(
		ctx,
		"refresh:"+userID,
	).Result()
}

func (c *Client) DeleteRefreshToken(
	ctx context.Context,
	userID string,
) error {

	return c.rdb.Del(
		ctx,
		"refresh:"+userID,
	).Err()
}
