package realtime_configs_service

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Storage interface {
	Get(ctx context.Context, key string, dst interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) GetValue(ctx context.Context, key string, serviceName string) (string, error) {
	key = fmt.Sprintf("%s_%s", serviceName, key)

	var value string
	ok, err := s.storage.Get(ctx, key, &value)
	if err != nil {
		return "", fmt.Errorf("storage.Get: %w", err)
	}

	if !ok {
		return "", errors.New("not found")
	}

	return value, nil
}

func (s *Service) SetValue(ctx context.Context, key, value, serviceName string) error {
	key = fmt.Sprintf("%s_%s", serviceName, key)

	err := s.storage.Set(ctx, key, value, 0)
	if err != nil {
		return fmt.Errorf("storage.Set: %w", err)
	}

	return nil
}
