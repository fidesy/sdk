package app

import (
	"context"
	desc "github.com/fidesy/sdk/services/domain-name-service/pkg/domain-name-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UpdateAddress(ctx context.Context, req *desc.UpdateAddressRequest) (*desc.UpdateAddressResponse, error) {
	err := i.domainNameService.UpdateAddress(ctx, req.GetServiceName(), req.GetAddress())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "domainNameService.UpdateAddress: %v", err)
	}

	return &desc.UpdateAddressResponse{}, nil
}
