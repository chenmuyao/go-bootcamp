package grpc_learn

import "context"

type Server struct {
	UnimplementedUserServiceServer
	Name string
}

// GetByID implements UserServiceServer.
func (s *Server) GetByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error) {
	return &GetByIDResponse{
		User: &User{
			Id:   123,
			Name: "name:" + s.Name,
		},
	}, nil
}

var _ UserServiceServer = &Server{}
