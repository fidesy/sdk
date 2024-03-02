package fooservice

import "fmt"

//go:generate mockgen -source ./fooservice.go -destination ./mocks/mock_fooservice.go

type (
	ServiceOption func(s *Service)

	Service struct {
		externalService ExternalService
	}

	ExternalService interface {
		Call() error
	}
)

func WithExternalService(externalService ExternalService) ServiceOption {
	return func(s *Service) {
		s.externalService = externalService
	}
}

func New(options ...ServiceOption) *Service {
	s := &Service{}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func (s *Service) Call() error {
	err := s.externalService.Call()
	if err != nil {
		return fmt.Errorf("externalService.Call: %w", err)
	}

	return nil
}
