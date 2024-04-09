package app

import (
	"context"
	desc "github.com/fidesy/sdk/services/realtime-configs-service/pkg/realtime-configs-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) SetValue(ctx context.Context, req *desc.SetValueRequest) (*desc.SetValueResponse, error) {
	err := i.realtimeConfigsService.SetValue(ctx, req.GetKey(), req.GetValue(), req.GetServiceName())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "realtimeConfigsService.SetValue: %v", err)
	}

	return &desc.SetValueResponse{}, nil
}
