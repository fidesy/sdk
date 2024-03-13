package grpc

import (
	"context"
	"fmt"
	grpcResolver "github.com/fidesy/sdk/common/grpc/resolver"
	randomCommon "github.com/fidesy/sdk/common/random"
	"google.golang.org/grpc/resolver"
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
	ServerOption func(s *Server) error
	Server       struct {
		port        string
		metricsPort string
	}
	ServiceDescriptor interface {
		GetDescription() *grpc.ServiceDesc
	}
)

func WithPort(port string) ServerOption {
	return func(s *Server) error {
		s.port = port
		return nil
	}
}

func WithMetricsPort(metricsPort string) ServerOption {
	return func(s *Server) error {
		s.metricsPort = metricsPort
		return nil
	}
}

func WithGraylog(graylogHost string) ServerOption {
	return func(s *Server) error {
		logAppName := strings.ReplaceAll(appName, "-", "_")

		w, err := gsyslog.DialLogger("udp", graylogHost, gsyslog.LOG_ERR, "SYSLOG", logAppName)
		if err != nil {
			return fmt.Errorf("gsyslog.DialLogger: %w", err)
		}

		logger.Init(w)

		return nil
	}
}

func WithTracer(jaegerEndpoint string) ServerOption {
	return func(s *Server) error {
		var err error
		_, closer, err = NewTracer(jaegerEndpoint)
		if err != nil {
			return fmt.Errorf("config.NewJaegerTracer: %v", err)
		}

		return nil
	}
}

func WithDomainNameService(ctx context.Context, dnsHost string) ServerOption {
	return func(s *Server) error {
		err := NewDomainNameService(ctx, dnsHost)
		if err != nil {
			return fmt.Errorf("NewDomainNameService: %w", err)
		}

		rb := &grpcResolver.Builder{
			DomainNameService: domainNameServiceClient,
		}
		resolver.Register(rb)

		return nil
	}
}

func NewServer(options ...ServerOption) (*Server, error) {
	if appName == "" {
		panic("APP_NAME env variable is required")
	}

	s := &Server{}

	initMetrics()

	for _, opt := range options {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	s.fillInDefaultValues()

	return s, nil
}

func (s *Server) Run(ctx context.Context, descs ...ServiceDescriptor) error {
	defer s.shutDown()

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

func (s *Server) shutDown() {
	if closer != nil {
		err := closer.Close()
		if err != nil {
			logger.Errorf("closer.Close: %v", err)
		}
	}
}

func (s *Server) fillInDefaultValues() {
	if logger.Get() == nil {
		logger.Init(os.Stdout)
		logger.Info("GRAYLOG_HOST env not found, using local logger")
	}

	if s.port == "" {
		randomPort := randomCommon.RandomPort()
		logger.Info(fmt.Sprintf("Using default grpc server port %s", randomPort))
		s.port = randomPort
	}

	if s.metricsPort == "" {
		randomPort := randomCommon.RandomPort()
		logger.Info(fmt.Sprintf("Using default metrics server port %s", randomPort))
		s.metricsPort = randomPort
	}
}
