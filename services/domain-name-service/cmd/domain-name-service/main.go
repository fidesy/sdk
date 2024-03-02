package main

import (
	"context"
	"fmt"
	"github.com/fidesy/sdk/services/domain-name-service/internal/app"
	"github.com/fidesy/sdk/services/domain-name-service/internal/config"
	"github.com/fidesy/sdk/services/domain-name-service/internal/pkg/domain-name-service"
	"github.com/fidesy/sdk/services/domain-name-service/internal/pkg/redis"
	desc "github.com/fidesy/sdk/services/domain-name-service/pkg/domain-name-service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	grpcPort  string
	proxyPort string
)

func main() {
	grpcPort = os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		log.Fatalf("GRPC_PORT ENV is required")
	}

	proxyPort = os.Getenv("PROXY_PORT")
	if proxyPort == "" {
		log.Fatalf("PROXY_PORT ENV is required")
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

	storage, err := redis.New(ctx)
	if err != nil {
		log.Fatalf("redis.New: %v", err)
	}

	domainNameService := domain_name_service.New(storage)

	impl := app.New(domainNameService)

	grpcServer := grpc.NewServer()
	grpcServer.RegisterService(&desc.DomainNameService_ServiceDesc, impl)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	errGroup := errgroup.Group{}

	// run http reverse proxy
	errGroup.Go(func() error {
		log.Printf("http proxy is running at %s port", proxyPort)

		err = runReverseProxy(ctx, impl)
		if err != nil {
			return fmt.Errorf("runReverseProxy: %w", err)
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

func runReverseProxy(ctx context.Context, impl *app.Implementation) error {
	router := runtime.NewServeMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
		Debug:          false,
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", proxyPort),
		Handler: corsHandler.Handler(router),
	}

	desc.RegisterDomainNameServiceHandlerServer(ctx, router, impl)

	return server.ListenAndServe()
}
