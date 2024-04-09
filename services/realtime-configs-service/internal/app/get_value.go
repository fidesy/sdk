package app

import (
	"context"
	desc "github.com/fidesy/sdk/services/realtime-configs-service/pkg/realtime-configs-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetValue(ctx context.Context, req *desc.GetValueRequest) (*desc.GetValueResponse, error) {
	value, err := i.realtimeConfigsService.GetValue(ctx, req.GetKey(), req.GetServiceName())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "realtimeConfigsService.GetValue: %v", err)
	}

	return &desc.GetValueResponse{
		Value: value,
	}, nil
}
