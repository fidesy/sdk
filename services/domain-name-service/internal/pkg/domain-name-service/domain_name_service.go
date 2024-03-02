package domain_name_service

import (
	"context"
	"fmt"
	"time"
)

type Storage interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, bytes []byte, expiration time.Duration) error
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) GetAddress(ctx context.Context, serviceName string) (string, error) {
	bytes, err := s.storage.Get(ctx, serviceName)
	if err != nil {
		return "", fmt.Errorf("storage.Get: %w", err)
	}

	return string(bytes), nil
}

func (s *Service) UpdateAddress(ctx context.Context, serviceName string, address string) error {
	err := s.storage.Set(ctx, serviceName, []byte(address), 0)
	if err != nil {
		return fmt.Errorf("storage.Set: %w", err)
	}

	return nil
}
