package app

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	desc "github.com/fidesy/sdk/examples/users-service/pkg/users-service"
)

func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.User, error) {
	if req.GetUsername() == "test" {
		return nil, status.Error(codes.Internal, "FUCK YOU")
	}

	return &desc.User{}, nil
}
