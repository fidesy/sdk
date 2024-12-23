package resolver

import (
	"context"
	"fmt"
	"google.golang.org/grpc/resolver"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	domain_name_service "github.com/fidesy/sdk/common/grpc/pkg/domain-name-service"
)

type Resolver struct {
	target string
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	cc     resolver.ClientConn

	domainNameService domain_name_service.DomainNameServiceClient
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	//TODO implement me
	return
}

func (r *Resolver) Close() {
	//TODO implement me
	return
}

func (r *Resolver) watch() {
	defer r.wg.Done()
	r.lookup(r.target)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			r.lookup(r.target)
		}
	}
}

func (r *Resolver) lookup(target string) {
	// Using different service names for different environments
	env := strings.ToUpper(os.Getenv("ENV"))
	if env == "LOCAL" || env == "STAGING" {
		target = fmt.Sprintf("%s-stage", target)
	}

	addressResp, err := r.domainNameService.GetAddress(context.TODO(), &domain_name_service.GetAddressRequest{
		ServiceName: target,
	})
	if err != nil {
		log.Printf("domainNameServiceClient for target %s: %v", target, err)
		return
	}

	address := resolver.Address{
		Addr:       addressResp.GetAddress(),
		ServerName: target,
	}
	// Обновляем адреса в ClientConn
	err = r.cc.UpdateState(resolver.State{
		Addresses: []resolver.Address{
			address,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
}
