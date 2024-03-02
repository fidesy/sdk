package grpc

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var metricsAppName = strings.ReplaceAll(os.Getenv("APP_NAME"), "-", "_")

var reg *prometheus.Registry

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_grpc_requests", metricsAppName),
			Help: "Count of grpc requests by handlers",
		},
		[]string{"handler"},
	)
	responseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("%s_grpc_response_time", metricsAppName),
			Help:    "Response time of grpc requests",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"handler"},
	)
	statusCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_status_codes", metricsAppName),
			Help: "Handlers not successful status codes",
		},
		[]string{"handler", "status_code"},
	)
)

func initMetrics() {
	reg = prometheus.NewRegistry()

	reg.MustRegister(
		requests,
		responseTime,
		statusCodes,
	)
}

func runMetrics(ctx context.Context, port string) error {
	// Create a new registry
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return err
	}
}

func GetRegistry() *prometheus.Registry {
	return reg
}

func metricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// update handler count metric
		requests.WithLabelValues(info.FullMethod).Inc()

		start := time.Now()
		resp, err := handler(ctx, req)

		responseTime.WithLabelValues(info.FullMethod).Observe(float64(time.Since(start).Milliseconds()))

		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				return resp, err
			}

			statusCodes.WithLabelValues(info.FullMethod, st.Code().String()).Inc()
		}

		return resp, err
	}
}
