package ioc

import (
	intrv1 "github.com/chenmuyao/go-bootcamp/api/proto/gen/intr/v1"
	"github.com/chenmuyao/go-bootcamp/interactive/config"
	"github.com/chenmuyao/go-bootcamp/interactive/service"
	"github.com/chenmuyao/go-bootcamp/internal/client"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitIntrClient(intrSvc service.InteractiveService) intrv1.InteractiveServiceClient {
	var opts []grpc.DialOption
	if !config.Cfg.GRPC.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.NewClient(config.Cfg.GRPC.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewLocalInteractiveAdapter(intrSvc)
	res := client.NewInteractiveClient(remote, local)
	viper.OnConfigChange(func(in fsnotify.Event) {
		th := viper.GetInt32("grpc.intr.threshold")
		res.UpdateThreshold(th)
	})
	return res
}
