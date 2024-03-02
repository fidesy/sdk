package app

import (
	desc "github.com/fidesy/sdk/examples/users-service/pkg/users-service"
	"google.golang.org/grpc"
)

type Implementation struct {
	desc.UnimplementedUserServiceServer
}

func New() *Implementation {
	return &Implementation{}
}

func (i *Implementation) GetDescription() *grpc.ServiceDesc {
	return &desc.UserService_ServiceDesc
}
