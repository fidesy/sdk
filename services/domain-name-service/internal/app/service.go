package app

import (
	"context"
	desc "github.com/fidesy/sdk/services/domain-name-service/pkg/domain-name-service"
)

type (
	Implementation struct {
		desc.UnimplementedDomainNameServiceServer

		domainNameService DomainNameService
	}

	DomainNameService interface {
		GetAddress(ctx context.Context, serviceName string) (string, error)
		UpdateAddress(ctx context.Context, serviceName string, address string) error
	}
)

func New(domainNameService DomainNameService) *Implementation {
	return &Implementation{
		domainNameService: domainNameService,
	}
}
