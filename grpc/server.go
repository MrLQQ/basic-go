package grpc

import (
	"context"
)

type Server struct {
	UnimplementedUserServiceServer
	Name string
}

func (s *Server) GetById(ctx context.Context, request *GetByIDRequest) (*GetByIDResponse, error) {
	return &GetByIDResponse{
		User: &User{
			Id:   123,
			Name: "daming",
		},
	}, nil
}
