package app

import (
	"context"
	desc "github.com/fidesy/sdk/services/realtime-configs-service/pkg/realtime-configs-service"
)

type (
	Implementation struct {
		desc.UnimplementedRealtimeConfigsServiceServer

		realtimeConfigsService RealtimeConfigsService
	}

	RealtimeConfigsService interface {
		GetValue(ctx context.Context, key, serviceName string) (string, error)
		SetValue(ctx context.Context, key, value, serviceName string) error
	}
)

func New(realtimeConfigsService RealtimeConfigsService) *Implementation {
	return &Implementation{
		realtimeConfigsService: realtimeConfigsService,
	}
}
