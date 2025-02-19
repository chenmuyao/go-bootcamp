package ioc

import (
	grpcIntr "github.com/chenmuyao/go-bootcamp/interactive/grpc"
	"github.com/chenmuyao/go-bootcamp/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewGrpcxServer(intrSvc *grpcIntr.InteractiveServiceServer) *grpcx.Server {
	s := grpc.NewServer()
	intrSvc.Register(s)
	return &grpcx.Server{
		Server: s,
		Addr:   viper.GetString("grpc.server.addr"),
	}
}
