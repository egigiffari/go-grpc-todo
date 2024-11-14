package service_test

import (
	"context"
	"net"
	"testing"

	"github.com/egigiffari/go-grpc-todo/pb"
	"github.com/egigiffari/go-grpc-todo/pkg/service"
	"github.com/egigiffari/go-grpc-todo/pkg/store"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestClientTest(t *testing.T) {
	t.Parallel()

	server, addr := startServer(t)
	client := newClient(t, addr)

	todo := todoRequest()
	expectedId := todo.Id
	req := &pb.CreateTodoRequest{Todo: todo}

	res, err := client.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res.Id, expectedId)

	storedData, err := server.Store.Read(expectedId)
	require.NoError(t, err)
	require.NotNil(t, storedData)

	requireSameTodo(t, storedData, todo)
}

func startServer(t *testing.T) (*service.Service, string) {
	service := service.New(store.NewInMemory())

	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, service)

	ln, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go grpcServer.Serve(ln)

	return service, ln.Addr().String()
}

func newClient(t *testing.T, serverAddr string) pb.TodoServiceClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client, err := grpc.NewClient(serverAddr, opts...)
	require.NoError(t, err)
	return pb.NewTodoServiceClient(client)
}

func requireSameTodo(t *testing.T, expected *pb.Todo, actually *pb.Todo) {
	json1, err := protojson.Marshal(expected)
	require.NoError(t, err)

	json2, err := protojson.Marshal(actually)
	require.NoError(t, err)

	require.Equal(t, string(json1), string(json2))
}
