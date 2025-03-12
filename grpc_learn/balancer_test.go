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
	grpc "google.golang.org/grpc"
	_ "google.golang.org/grpc/balancer/weightedroundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

type BalancerTestSuite struct {
	suite.Suite
	cli *clientv3.Client
}

func (s *BalancerTestSuite) SetupSuite() {
	cli, err := clientv3.NewFromURL("localhost:12379")
	require.NoError(s.T(), err)
	s.cli = cli
}

func (s *BalancerTestSuite) TestClientWRR() {
	t := s.T()
	etcdResolver, err := resolver.NewBuilder(s.cli)
	require.NoError(t, err)
	cc, err := grpc.NewClient(
		"etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`
{
  "loadBalancingConfig": [
    {
      "weighted_round_robin": {}
    }
  ]
}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	client := NewUserServiceClient(cc)
	for range 10 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.GetByID(ctx, &GetByIDRequest{
			Id: 123,
		})
		cancel()
		require.NoError(t, err)
		log.Println(resp.User)
	}
}

func (s *BalancerTestSuite) TestClient() {
	t := s.T()
	etcdResolver, err := resolver.NewBuilder(s.cli)
	require.NoError(t, err)
	cc, err := grpc.NewClient(
		"etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`
{
  "loadBalancingConfig": [
    {
      "round_robin": {}
    }
  ]
}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	client := NewUserServiceClient(cc)
	for range 10 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.GetByID(ctx, &GetByIDRequest{
			Id: 123,
		})
		cancel()
		require.NoError(t, err)
		log.Println(resp.User)
	}
}

func (s *BalancerTestSuite) TestServer() {
	go func() {
		s.startServer(":8090")
	}()
	s.startServer(":8091")
}

func (s *BalancerTestSuite) startServer(port string) {
	t := s.T()
	em, err := endpoints.NewManager(s.cli, "service/user")
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	addr := "127.0.0.1" + port
	key := "service/user/" + addr
	l, err := net.Listen("tcp", addr)
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
		_, err1 := s.cli.KeepAlive(kaCtx, leaseResp.ID)
		require.NoError(t, err1)
	}()

	server := grpc.NewServer()
	RegisterUserServiceServer(server, &Server{Name: addr})
	server.Serve(l)

	kaCancel()
	err = em.DeleteEndpoint(ctx, key)
	if err != nil {
		t.Log(err)
	}
	server.GracefulStop()
	// s.cli.Close()
}

// func TestBalancer(t *testing.T) {
// 	suite.Run(t, new(BalancerTestSuite))
// }
