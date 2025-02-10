package grpc_learn

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestOneOf(t *testing.T) {
	u := &User{}
	email, ok := u.Contacts.(*User_Email)
	if ok {
		t.Log("Is email", email)
		return
	}
}

func TestServer(t *testing.T) {
	gs := grpc.NewServer()
	us := &Server{}
	RegisterUserServiceServer(gs, us)

	l, err := net.Listen("tcp", ":8090")
	assert.NoError(t, err)
	if err = gs.Serve(l); err != nil {
		t.Log("exited", err)
	}
}

func TestClient(t *testing.T) {
	cc, err := grpc.NewClient(
		"localhost:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	client := NewUserServiceClient(cc)
	resp, err := client.GetByID(context.Background(), &GetByIDRequest{Id: 123})
	assert.NoError(t, err)
	t.Log(resp)
}
