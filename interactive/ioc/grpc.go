package ioc

import (
	grpcIntr "github.com/chenmuyao/go-bootcamp/interactive/grpc"
	"github.com/chenmuyao/go-bootcamp/pkg/grpcx"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func NewGrpcxServer(intrSvc *grpcIntr.InteractiveServiceServer) *grpcx.Server {
	s := grpc.NewServer()
	intrSvc.Register(s)
	cli, err := clientv3.NewFromURL(viper.GetString("grpc.etcdAddr"))
	if err != nil {
		panic(err)
	}
	return grpcx.NewServer(s, cli, viper.GetInt("grpc.port"), viper.GetString("grpc.serviceName"))
}
