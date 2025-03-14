package grpc_learn

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FailedServer struct {
	UnimplementedUserServiceServer
	Name string
}

// GetByID implements UserServiceServer.
func (s *FailedServer) GetByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error) {
	log.Println("Faileover")
	return nil, status.Errorf(codes.Unavailable, "service cut")
}

var _ UserServiceServer = &Server{}
