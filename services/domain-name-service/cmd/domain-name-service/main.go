package main

import (
	"context"
	"fmt"
	"github.com/fidesy/sdk/services/domain-name-service/internal/app"
	"github.com/fidesy/sdk/services/domain-name-service/internal/config"
	"github.com/fidesy/sdk/services/domain-name-service/internal/pkg/domain-name-service"
	"github.com/fidesy/sdk/services/domain-name-service/internal/pkg/metrics"
	"github.com/fidesy/sdk/services/domain-name-service/internal/pkg/redis"
	desc "github.com/fidesy/sdk/services/domain-name-service/pkg/domain-name-service"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	grpcPort    string
	metricsPort string
)

func main() {
	grpcPort = os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		log.Fatalf("GRPC_PORT ENV is required")
	}

	metricsPort = os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		log.Fatalf("METRICS_PORT ENV is required")
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()

	err := config.Init()
	if err != nil {
		log.Fatalf("config.Init: %v", err)
	}

	metrics.Init()

	storage, err := redis.New(ctx)
	if err != nil {
		log.Fatalf("redis.New: %v", err)
	}

	domainNameService := domain_name_service.New(storage)

	impl := app.New(domainNameService)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			metrics.Interceptor(),
		),
	)
	grpcServer.RegisterService(&desc.DomainNameService_ServiceDesc, impl)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	errGroup := errgroup.Group{}

	// run metrics
	errGroup.Go(func() error {
		log.Printf("metrics endpoint is running at %s port", metricsPort)

		err = metrics.Run(ctx, metricsPort)
		if err != nil {
			return fmt.Errorf("metrics.Run: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		log.Printf("grpcServer is running at %s port", grpcPort)
		err = grpcServer.Serve(lis)
		if err != nil {
			return fmt.Errorf("grpcServer.Serve: %v", err)
		}

		return nil
	})

	err = errGroup.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
