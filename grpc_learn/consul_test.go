package grpc_learn

import (
	"net"
	"testing"

	capi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	grpc "google.golang.org/grpc"
)

type ConsulTestSuite struct {
	suite.Suite
	cli *capi.Client
}

func (s *ConsulTestSuite) SetupSuite() {
	cli, err := capi.NewClient(capi.DefaultConfig())
	require.NoError(s.T(), err)
	s.cli = cli
}

func (s *ConsulTestSuite) TestClient() {
	// t := s.T()
	// etcdResolver, err := resolver.NewBuilder(s.cli)
	// require.NoError(t, err)
	// cc, err := grpc.NewClient(
	// 	"etcd:///service/user",
	// 	grpc.WithResolvers(etcdResolver),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	// require.NoError(t, err)
	// client := NewUserServiceClient(cc)
	// resp, err := client.GetByID(context.Background(), &GetByIDRequest{
	// 	Id: 123,
	// })
	// require.NoError(t, err)
	// log.Println(resp.User)
}

func (s *ConsulTestSuite) TestServer() {
	t := s.T()
	service := &capi.AgentServiceRegistration{
		Name:    "consul-center",
		ID:      "service/user",
		Address: "127.0.0.1",
		Port:    8090,
	}
	err := s.cli.Agent().ServiceRegister(service)
	require.NoError(t, err)
	l, err := net.Listen("tcp", ":8091")
	require.NoError(t, err)

	server := grpc.NewServer()
	RegisterUserServiceServer(server, &Server{})
	server.Serve(l)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	// addr := "127.0.0.1:8090"
	// key := "service/user/" + addr
	// l, err := net.Listen("tcp", ":8090")
	// require.NoError(t, err)
	//
	// // Use lease
	// var ttl int64 = 5
	// leaseResp, err := s.cli.Grant(ctx, ttl)
	// require.NoError(t, err)
	//
	// err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
	// 	Addr: addr,
	// }, clientv3.WithLease(leaseResp.ID))
	// require.NoError(t, err)
	//
	// // keep alive
	// kaCtx, kaCancel := context.WithCancel(context.Background())
	// go func() {
	// 	ch, err1 := s.cli.KeepAlive(kaCtx, leaseResp.ID)
	// 	require.NoError(t, err1)
	// 	for kaResp := range ch {
	// 		t.Log(kaResp.String())
	// 	}
	// }()
	//
	// go func() {
	// 	// simulate registration change
	// 	ticker := time.NewTicker(time.Second)
	// 	for now := range ticker.C {
	// 		ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
	// 		err1 := em.Update(ctx1, []*endpoints.UpdateWithOpts{
	// 			{
	// 				Update: endpoints.Update{
	// 					Op:  endpoints.Add,
	// 					Key: key,
	// 					Endpoint: endpoints.Endpoint{
	// 						Addr:     addr,
	// 						Metadata: now.String(),
	// 					},
	// 				},
	// 				Opts: []clientv3.OpOption{
	// 					clientv3.WithLease(leaseResp.ID),
	// 				},
	// 			},
	// 		})
	// 		// err1 := em.AddEndpoint(ctx1, key, endpoints.Endpoint{
	// 		// 	Addr:     addr,
	// 		// 	Metadata: now.String(),
	// 		// }, clientv3.WithLease(leaseResp.ID))
	// 		cancel1()
	// 		if err1 != nil {
	// 			t.Log(err1)
	// 		}
	// 	}
	// }()
	//
	// server := grpc.NewServer()
	// RegisterUserServiceServer(server, &Server{})
	// server.Serve(l)
	//
	// kaCancel()
	// err = em.DeleteEndpoint(ctx, key)
	// if err != nil {
	// 	t.Log(err)
	// }
	// server.GracefulStop()
	// s.cli.Close()
}

func TestConsul(t *testing.T) {
	suite.Run(t, new(ConsulTestSuite))
}
