package grpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	domain_name_service "github.com/fidesy/sdk/common/grpc/pkg/domain-name-service"
	"github.com/fidesy/sdk/common/logger"
	gsyslog "github.com/hashicorp/go-syslog"
	"google.golang.org/grpc"
)

var appName = os.Getenv("APP_NAME")

type (
	ServerOption func(s *Server)
	Server       struct {
		port           string
		metricsPort    string
		jaegerEndpoint string

		dnsHost string
	}
	ServiceDescriptor interface {
		GetDescription() *grpc.ServiceDesc
	}
)

func WithPort(port string) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithMetricsPort(metricsPort string) ServerOption {
	return func(s *Server) {
		s.metricsPort = metricsPort
	}
}

func WithTracer(jaegerEndpoint string) ServerOption {
	return func(s *Server) {
		s.jaegerEndpoint = jaegerEndpoint
	}
}

func WithDomainNameService(dnsHost string) ServerOption {
	return func(s *Server) {
		s.dnsHost = dnsHost
	}
}

func NewServer(options ...ServerOption) *Server {
	if appName == "" {
		panic("APP_NAME env variable is required")
	}

	s := &Server{}

	initLogger()
	initMetrics()

	for _, opt := range options {
		opt(s)
	}

	s.fillInDefaultValues()

	return s
}

func (s *Server) Run(ctx context.Context, descs ...ServiceDescriptor) error {
	// if jaeger endpoint passed
	// then init jaeger
	if s.jaegerEndpoint != "" {
		_, closer, err := NewTracer(s.jaegerEndpoint)
		if err != nil {
			logger.Fatalf("config.NewJaegerTracer: %v", err)
		}
		defer closer.Close()
	}

	if s.dnsHost != "" {
		err := NewDomainNameService(ctx, s.dnsHost)
		if err != nil {
			return fmt.Errorf("NewDomainNameService: %w", err)
		}
	}

	interceptors := []grpc.UnaryServerInterceptor{
		metricsInterceptor(),
	}
	if tracer != nil {
		interceptors = append(interceptors, tracingInterceptor())
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors...,
		),
	)

	for _, desc := range descs {
		grpcServer.RegisterService(desc.GetDescription(), desc)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	errChan := make(chan error)

	go func() {
		if domainNameServiceClient != nil {
			_, err = domainNameServiceClient.UpdateAddress(ctx, &domain_name_service.UpdateAddressRequest{
				ServiceName: appName,
				Address:     fmt.Sprintf("%s:%s", appName, s.port),
			})
			if err != nil {
				errChan <- fmt.Errorf("domainNameServiceClient.UpdatePort: %w", err)
				return
			}
		}

		logger.Info(fmt.Sprintf("grpcServer is running at %s port", s.port))
		if err = grpcServer.Serve(lis); err != nil {
			errChan <- fmt.Errorf("grpcServer.Serve: %w", err)
		}
	}()

	go func() {
		logger.Info(fmt.Sprintf("metrics are running at %s port", s.metricsPort))
		if err = runMetrics(ctx, s.metricsPort); err != nil {
			errChan <- fmt.Errorf("runMetrics: %w", err)
		}

	}()

	select {
	case <-ctx.Done():
		return nil
	case err = <-errChan:
		return err
	}
}

func (s *Server) fillInDefaultValues() {
	if s.port == "" {
		logger.Info("Using default grpc server port 8080")
		s.port = "8080"
	}

	if s.metricsPort == "" {
		logger.Info("Using default metrics server port 8081")
		s.metricsPort = "8081"
	}
}

func initLogger() {
	graylogHost := os.Getenv("GRAYLOG_HOST")

	logAppName := strings.ReplaceAll(appName, "-", "_")

	w, err := gsyslog.DialLogger("udp", graylogHost, gsyslog.LOG_ERR, "SYSLOG", logAppName)
	// If there is an error or no graylog host
	// then write logs to STDOUT
	if err != nil || graylogHost == "" {
		logger.Init(os.Stdout)
		logger.Info("GRAYLOG_HOST env not found, using local logger")
		return
	}

	logger.Init(w)
}
