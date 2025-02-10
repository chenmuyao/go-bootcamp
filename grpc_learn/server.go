package grpc_learn

import "context"

type Server struct {
	UnimplementedUserServiceServer
}

// GetByID implements UserServiceServer.
func (s *Server) GetByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error) {
	return &GetByIDResponse{
		User: &User{
			Id:   123,
			Name: "my name",
		},
	}, nil
}

var _ UserServiceServer = &Server{}
