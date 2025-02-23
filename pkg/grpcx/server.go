package grpcx

import (
	"context"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/chenmuyao/go-bootcamp/pkg/netx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	cli         *clientv3.Client
	port        int
	KaCancel    func()
	serviceName string
}

func NewServer(
	server *grpc.Server,
	cli *clientv3.Client,
	port int,
	serviceName string,
) *Server {
	return &Server{
		Server:      server,
		cli:         cli,
		port:        port,
		serviceName: serviceName,
	}
}

func (s *Server) Serve() error {
	addr := ":" + strconv.Itoa(s.port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	err = s.register()
	if err != nil {
		return err
	}
	return s.Server.Serve(l)
}

func (s *Server) register() error {
	em, err := endpoints.NewManager(s.cli, "service/"+s.serviceName)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	addr := netx.GetoutboundIP() + ":" + strconv.Itoa(s.port)
	key := "service/" + s.serviceName + "/" + addr
	if err != nil {
		return err
	}

	// Use lease
	var ttl int64 = 5
	leaseResp, err := s.cli.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	// keep alive
	var kaCtx context.Context
	kaCtx, s.KaCancel = context.WithCancel(context.Background())
	ch, err := s.cli.KeepAlive(kaCtx, leaseResp.ID)
	go func() {
		for kaResp := range ch {
			slog.Debug("debug", "ch", kaResp)
		}
	}()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	if s.KaCancel != nil {
		s.KaCancel()
	}
	// if s.cli != nil {
	// 	return s.cli.Close()
	// }
	s.GracefulStop()
	return nil
}
