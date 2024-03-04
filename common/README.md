### gRPC SDK

Create gRPC server 
```
grpcServer, err := grpc.NewServer()
if err != nil {
    log.Fatalf("grpc.NewServer: %v", err)
}
```

The default is port 8080 for gRPC and port 8081 for metrics. You can update it using options:

```
grpcServer, err := grpc.NewServer(
    grpc.WithPort("11000"),
    grpc.WithMetricsPort("11001"),
)
...
```

Jaeger tracer option:
```
grpcServer, err := grpc.NewServer(
    grpc.WithJaeger("jaeger:5555"),
)
...
```

Graylog option:
```
grpcServer, err := grpc.NewServer(
    grpc.WithGraylog("graylog:9000"),
)
...
```