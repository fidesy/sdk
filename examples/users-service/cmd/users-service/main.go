package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fidesy/sdk/common/grpc"
	"github.com/fidesy/sdk/common/logger"
	"github.com/fidesy/sdk/examples/users-service/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer cancel()

	grpcServer, err := grpc.NewServer(
		grpc.WithPort(os.Getenv("GRPC_PORT")),
		grpc.WithMetricsPort(os.Getenv("METRICS_PORT")),
		grpc.WithDomainNameService(ctx, "localhost:10000"),
	)
	if err != nil {
		log.Fatalf("grpc.NewServer: %v", err)
	}

	impl := app.New()

	err = grpcServer.Run(ctx, impl)
	if err != nil {
		logger.Fatalf("grpcServer.Run: %v", err)
	}
}
