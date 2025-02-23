package grpc_learn

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EtcdTestSuite struct {
	suite.Suite
	cli *clientv3.Client
}

func (s *EtcdTestSuite) SetupSuite() {
	cli, err := clientv3.NewFromURL("localhost:12379")
	require.NoError(s.T(), err)
	s.cli = cli
}

func (s *EtcdTestSuite) TestClient() {
	t := s.T()
	etcdResolver, err := resolver.NewBuilder(s.cli)
	require.NoError(t, err)
	cc, err := grpc.NewClient(
		"etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	client := NewUserServiceClient(cc)
	resp, err := client.GetByID(context.Background(), &GetByIDRequest{
		Id: 123,
	})
	require.NoError(t, err)
	log.Println(resp.User)
}

func (s *EtcdTestSuite) TestServer() {
	t := s.T()
	em, err := endpoints.NewManager(s.cli, "service/user")
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	addr := "127.0.0.1:8090"
	key := "service/user/" + addr
	l, err := net.Listen("tcp", ":8090")
	require.NoError(t, err)

	// Use lease
	var ttl int64 = 5
	leaseResp, err := s.cli.Grant(ctx, ttl)
	require.NoError(t, err)

	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, clientv3.WithLease(leaseResp.ID))
	require.NoError(t, err)

	// keep alive
	kaCtx, kaCancel := context.WithCancel(context.Background())
	go func() {
		ch, err1 := s.cli.KeepAlive(kaCtx, leaseResp.ID)
		require.NoError(t, err1)
		for kaResp := range ch {
			t.Log(kaResp.String())
		}
	}()

	go func() {
		// simulate registration change
		ticker := time.NewTicker(time.Second)
		for now := range ticker.C {
			ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
			err1 := em.Update(ctx1, []*endpoints.UpdateWithOpts{
				{
					Update: endpoints.Update{
						Op:  endpoints.Add,
						Key: key,
						Endpoint: endpoints.Endpoint{
							Addr:     addr,
							Metadata: now.String(),
						},
					},
					Opts: []clientv3.OpOption{
						clientv3.WithLease(leaseResp.ID),
					},
				},
			})
			// err1 := em.AddEndpoint(ctx1, key, endpoints.Endpoint{
			// 	Addr:     addr,
			// 	Metadata: now.String(),
			// }, clientv3.WithLease(leaseResp.ID))
			cancel1()
			if err1 != nil {
				t.Log(err1)
			}
		}
	}()

	server := grpc.NewServer()
	RegisterUserServiceServer(server, &Server{})
	server.Serve(l)

	kaCancel()
	err = em.DeleteEndpoint(ctx, key)
	if err != nil {
		t.Log(err)
	}
	server.GracefulStop()
	s.cli.Close()
}

// func TestEtcd(t *testing.T) {
// 	suite.Run(t, new(EtcdTestSuite))
// }
