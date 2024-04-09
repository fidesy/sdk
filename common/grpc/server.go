package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/fidesy/sdk/common/grpc/config"
	grpcResolver "github.com/fidesy/sdk/common/grpc/resolver"
	randomCommon "github.com/fidesy/sdk/common/random"
	realtime_configs_service "github.com/fidesy/sdk/services/realtime-configs-service/pkg/realtime-configs-service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/resolver"

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
		proxyPort   string
		proxyRouter *runtime.ServeMux
		swaggerPort string
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

func WithProxyPort(proxyPort string) ServerOption {
	return func(s *Server) error {
		s.proxyPort = proxyPort
		s.proxyRouter = runtime.NewServeMux()
		return nil
	}
}

func WithSwaggerPort(swaggerPort string) ServerOption {
	return func(s *Server) error {
		s.swaggerPort = swaggerPort
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

func WithRealtimeConfigsService(ctx context.Context, dnsHost string) ServerOption {
	return func(s *Server) error {
		client, err := NewClient[realtime_configs_service.RealtimeConfigsServiceClient](
			ctx,
			realtime_configs_service.NewRealtimeConfigsServiceClient,
			dnsHost,
		)
		if err != nil {
			return fmt.Errorf("NewRealtimeConfigsService: %w", err)
		}

		config.Init(client)

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

	errGroup := errgroup.Group{}

	errGroup.Go(func() error {
		if domainNameServiceClient != nil {
			_, err = domainNameServiceClient.UpdateAddress(ctx, &domain_name_service.UpdateAddressRequest{
				ServiceName: appName,
				Address:     fmt.Sprintf("%s:%s", appName, s.port),
			})
			if err != nil {
				return fmt.Errorf("domainNameServiceClient.UpdatePort: %w", err)
			}
		}

		logger.Info(fmt.Sprintf("grpcServer is running at %s port", s.port))
		if err = grpcServer.Serve(lis); err != nil {
			return fmt.Errorf("grpcServer.Serve: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		logger.Info(fmt.Sprintf("metrics are running at %s port", s.metricsPort))
		if err = runMetrics(ctx, s.metricsPort); err != nil {
			return fmt.Errorf("runMetrics: %w", err)
		}

		return nil
	})

	if s.proxyPort != "" {
		errGroup.Go(func() error {
			server := &http.Server{
				Addr:    fmt.Sprintf(":%s", s.proxyPort),
				Handler: s.proxyRouter,
			}

			if err = server.ListenAndServe(); err != nil {
				return fmt.Errorf("httpProxy: server.ListenAndServe: %w", err)
			}
			return nil
		})
	}

	if s.swaggerPort != "" {
		errGroup.Go(func() error {
			fs := http.FileServer(http.Dir("./swaggerui"))
			http.Handle("/docs/", http.StripPrefix("/docs/", fs))

			if err := http.ListenAndServe(fmt.Sprintf(":%s", s.swaggerPort), nil); err != nil {
				return fmt.Errorf("swagger http.ListenAndServe: %w", err)
			}

			return nil
		})
	}

	return errGroup.Wait()
}

func (s *Server) ProxyRouter() *runtime.ServeMux {
	return s.proxyRouter
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
