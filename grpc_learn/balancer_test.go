package grpc_learn

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	_ "github.com/chenmuyao/go-bootcamp/pkg/grpcx/balancer/wrr"
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

func (s *BalancerTestSuite) TestClientCustomWRR() {
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
      "custom_weighted_round_robin": {}
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

func (s *BalancerTestSuite) TestFailoverClient() {
	t := s.T()
	etcdResolver, err := resolver.NewBuilder(s.cli)
	require.NoError(t, err)
	cc, err := grpc.NewClient(
		"etcd:///service/user",
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`
{
    "loadBalancingConfig": [{"round_robin": {}}],
    "methodConfig": [
        {
            "name": [{"service": "UserService"}],
            "retryPolicy": {
                "maxAttempts": 4,
                "initialBackoff": "0.01s",
                "maxBackoff": "0.1s",
                "backoffMultiplier": 2.0,
                "retryableStatusCodes": ["INVALID_ARGUMENT", "UNAVAILABLE"]
            }
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
		s.startServer(":8090", 10, &Server{Name: ":8090"})
	}()
	go func() {
		s.startServer(":8091", 20, &Server{Name: ":8091"})
	}()
	s.startServer(":8092", 30, &FailedServer{Name: ":8092"})
}

func (s *BalancerTestSuite) startServerV1(port string) {
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

func (s *BalancerTestSuite) startServer(port string, weight int, srv UserServiceServer) {
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
		Metadata: map[string]any{
			"weight": weight,
		},
	}, clientv3.WithLease(leaseResp.ID))
	require.NoError(t, err)

	// keep alive
	kaCtx, kaCancel := context.WithCancel(context.Background())
	go func() {
		_, err1 := s.cli.KeepAlive(kaCtx, leaseResp.ID)
		require.NoError(t, err1)
	}()

	server := grpc.NewServer()
	RegisterUserServiceServer(server, srv)
	server.Serve(l)

	kaCancel()
	err = em.DeleteEndpoint(ctx, key)
	if err != nil {
		t.Log(err)
	}
	server.GracefulStop()
	// s.cli.Close()
}

func TestBalancer(t *testing.T) {
	suite.Run(t, new(BalancerTestSuite))
}
