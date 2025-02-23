package ioc

import (
	"log/slog"

	intrv1 "github.com/chenmuyao/go-bootcamp/api/proto/gen/intr/v1"
	"github.com/chenmuyao/go-bootcamp/config"
	"github.com/chenmuyao/go-bootcamp/interactive/service"
	"github.com/chenmuyao/go-bootcamp/internal/client"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitIntrClient(intrSvc service.InteractiveService) intrv1.InteractiveServiceClient {
	var opts []grpc.DialOption
	if !config.Cfg.GRPC.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.NewClient(config.Cfg.GRPC.Intr.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewLocalInteractiveAdapter(intrSvc)
	res := client.NewInteractiveClient(remote, local)
	viper.OnConfigChange(func(in fsnotify.Event) {
		th := config.Cfg.GRPC.Intr.Threshold
		slog.Info("change threshold", "th", th)
		res.UpdateThreshold(th)
	})
	return res
}

func InitIntrClientEtcd(client *clientv3.Client) intrv1.InteractiveServiceClient {
	resolver, err := resolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{
		grpc.WithResolvers(resolver),
	}
	if !config.Cfg.GRPC.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.NewClient(config.Cfg.GRPC.Intr.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(cc)
	return remote
}
