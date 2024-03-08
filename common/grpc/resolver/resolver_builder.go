package resolver

import (
	"context"
	"google.golang.org/grpc/resolver"
	"sync"

	domain_name_service "github.com/fidesy/sdk/common/grpc/pkg/domain-name-service"
)

type Builder struct {
	DomainNameService domain_name_service.DomainNameServiceClient
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &Resolver{
		target:            target.Endpoint(),
		ctx:               ctx,
		cancel:            cancel,
		wg:                sync.WaitGroup{},
		cc:                cc,
		domainNameService: b.DomainNameService,
	}
	r.wg.Add(1)
	// Та самая горутина, которая в фоне будет обновлять адреса
	go r.watch()
	return r, nil
}

func (b *Builder) Scheme() string {
	return "rpc"
}
