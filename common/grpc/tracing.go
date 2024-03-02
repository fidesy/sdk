package grpc

import (
	"context"
	"fmt"
	"github.com/fidesy/sdk/common/logger"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var tracer opentracing.Tracer

func NewTracer(jaegerEndpoint string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			CollectorEndpoint:   jaegerEndpoint,
		},
		ServiceName: appName,
	}

	var (
		closer io.Closer
		err    error
	)

	tracer, closer, err = cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer, nil
}

func tracingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		traceIDStr := ""

		// trace id is present in metadata
		if len(md.Get("x-trace-id")) > 0 {
			spanID, _ := jaeger.SpanIDFromString(md.Get("x-span-id")[0])
			traceID, _ := jaeger.TraceIDFromString(md.Get("x-trace-id")[0])
			traceIDStr = traceID.String()

			parentSpanCtx := jaeger.NewSpanContext(traceID, spanID, jaeger.SpanID(0), false, nil)
			span := tracer.StartSpan(
				info.FullMethod,
				ext.RPCServerOption(parentSpanCtx),
			)
			defer span.Finish()

			jaegerSpan, _ := span.(*jaeger.Span)
			ctx = metadata.AppendToOutgoingContext(ctx, "x-span-id", jaegerSpan.SpanContext().SpanID().String())
			ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceID.String())
		} else {
			span := tracer.StartSpan(info.FullMethod)
			defer span.Finish()

			jaegerSpan, _ := span.(*jaeger.Span)
			ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", fmt.Sprint(jaegerSpan.SpanContext().TraceID()))
			ctx = metadata.AppendToOutgoingContext(ctx, "x-span-id", fmt.Sprint(jaegerSpan.SpanContext().TraceID()))
			traceIDStr = jaegerSpan.SpanContext().TraceID().String()
		}

		// Call the gRPC handler.
		response, handlerErr := handler(ctx, req)

		if handlerErr == nil {
			return response, nil
		}

		st, ok := status.FromError(handlerErr)
		if !ok {
			return response, handlerErr
		}
		// log error
		logger.Errorf(
			handlerErr,
			zap.String("trace_id", traceIDStr),
			zap.String("status_code", st.Code().String()),
		)

		return response, handlerErr
	}
}

func GetTracer() opentracing.Tracer {
	return tracer
}
