package app

import (
	"context"
	"errors"
	"github.com/fidesy/sdk/services/domain-name-service/internal/pkg/redis"
	desc "github.com/fidesy/sdk/services/domain-name-service/pkg/domain-name-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetAddress(ctx context.Context, req *desc.GetAddressRequest) (*desc.GetAddressResponse, error) {
	address, err := i.domainNameService.GetAddress(ctx, req.GetServiceName())
	if err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "domainNameService.GetAddress: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "domainNameService.GetAddress: %v", err)
	}

	return &desc.GetAddressResponse{
		Address: address,
	}, nil
}
