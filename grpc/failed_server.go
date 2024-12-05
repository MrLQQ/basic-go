package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type FailedServer struct {
	UnimplementedUserServiceServer
	Name string
}

func (f *FailedServer) GetById(ctx context.Context, in *GetByIDRequest) (*GetByIDResponse, error) {
	log.Println("进来了 failOver")
	return nil, status.Error(codes.Unavailable, "假设服务被熔断了")
}
