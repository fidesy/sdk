package grpc

import (
	"context"
	domain_name_service "github.com/fidesy/sdk/common/grpc/pkg/domain-name-service"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	domainNameServiceClient domain_name_service.DomainNameServiceClient
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
