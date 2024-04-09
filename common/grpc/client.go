package grpc

import (
	"context"
	realtime_configs_service "github.com/fidesy/sdk/services/realtime-configs-service/pkg/realtime-configs-service"

	domain_name_service "github.com/fidesy/sdk/common/grpc/pkg/domain-name-service"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	domainNameServiceClient      domain_name_service.DomainNameServiceClient
	realtimeConfigsServiceClient realtime_configs_service.RealtimeConfigsServiceClient
)

type ConnWrapper[Client any] func(_ grpc.ClientConnInterface) Client

func NewDomainNameService(ctx context.Context, domainNameServiceHost string) error {
	conn, err := grpc.DialContext(
		ctx,
		domainNameServiceHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	domainNameServiceClient = domain_name_service.NewDomainNameServiceClient(conn)

	return nil
}

func NewRealtimeConfigsService(ctx context.Context, realtimeConfigsServiceHost string) error {
	conn, err := grpc.DialContext(
		ctx,
		realtimeConfigsServiceHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	realtimeConfigsServiceClient = realtime_configs_service.NewRealtimeConfigsServiceClient(conn)

	return nil
}

func NewClient[Client any](ctx context.Context, connWrapper ConnWrapper[Client], serviceName string) (Client, error) {
	conn, err := grpc.DialContext(
		ctx,
		serviceName,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return lo.Empty[Client](), err
	}

	return connWrapper(conn), nil
}
