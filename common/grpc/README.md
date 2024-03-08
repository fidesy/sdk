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

Use domain-name-service to automatically resolve addresses and connect to other gRPC services by name only.
First of all you have to run this service [sources](https://github.com/fidesy/sdk/tree/master/services/domain-name-service).
```
grpcServer, err := grpc.NewServer(
    grpc.WithDomainNameService(
         context.TODO(), 
         "localhost:10000",
     ),
)
...
```





Connect to another gRPC service
1. You are not using domain-name-service to resolve addresses by service-name

```
authClient, err := grpc.NewClient[auth_service.AuthServiceClient](
        ctx,
	auth_service.NewAuthServiceClient,
	"localhost:7040",
)
...
```

2. You are using domain-name-service

```
authClient, err := grpc.NewClient[auth_service.AuthServiceClient](
        ctx,
	auth_service.NewAuthServiceClient,
	"rpc:///auth-service",
)
...
```