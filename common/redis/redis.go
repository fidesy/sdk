package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redisClient *redis.Client
}

func Connect(ctx context.Context, options *redis.Options) (*Client, error) {
	client := redis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis.Ping: %w", err)
	}

	return &Client{redisClient: client}, nil
}

func (c *Client) RedisClient() *redis.Client {
	return c.redisClient
}

func (c *Client) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = c.redisClient.Set(ctx, key, bytes, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string, dst interface{}) (bool, error) {
	result := c.redisClient.Get(ctx, key)
	if err := result.Err(); err != nil {
		// not found
		if result.Err() == redis.Nil {
			return false, nil
		}

		return false, fmt.Errorf("redis.Get: %w", err)
	}

	bytes, err := result.Bytes()
	if err != nil {
		return false, fmt.Errorf("result.Bytes: %w", err)
	}

	err = json.Unmarshal(bytes, &dst)
	if err != nil {
		return false, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return true, nil
}
