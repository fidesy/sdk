package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/fidesy/sdk/services/domain-name-service/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrNotFound = errors.New("record with provided serviceName not found")
)

type Service struct {
	db *redis.Client
}

func New(ctx context.Context) (*Service, error) {
	c := &Service{}

	cli := redis.NewClient(&redis.Options{
		Addr:     config.Get(config.RedisHost).(string),
		Password: config.Get(config.RedisPassword).(string),
		DB:       0,
	})

	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis.Ping: %w", err)
	}

	c.db = cli

	return c, nil
}

func (s *Service) Set(ctx context.Context, key string, bytes []byte, expiration time.Duration) error {
	err := s.db.Set(ctx, key, bytes, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, key string) ([]byte, error) {
	result := s.db.Get(ctx, key)
	if err := result.Err(); err != nil {
		// not found
		if result.Err() == redis.Nil {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("redis.Get: %w", err)
	}

	bytes, err := result.Bytes()
	if err != nil {
		return nil, fmt.Errorf("result.Bytes: %w", err)
	}

	return bytes, nil
}

func (s *Service) Size(ctx context.Context) (int, error) {
	size, err := s.db.DBSize(ctx).Result()
	if err != nil {
		return 0, fmt.Errorf("redis.DBSize: %w", err)
	}

	return int(size), nil
}
