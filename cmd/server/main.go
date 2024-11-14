package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/egigiffari/go-grpc-todo/pb"
	"github.com/egigiffari/go-grpc-todo/pkg/service"
	"github.com/egigiffari/go-grpc-todo/pkg/store"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()

	service := service.New(store.NewInMemory())
	server := grpc.NewServer()
	pb.RegisterTodoServiceServer(server, service)

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("cannot start server, err: %s", err.Error())
	}

	fmt.Printf("server running on %s\n", addr)
	if err := server.Serve(ln); err != nil {
		log.Fatalf("cannot start server, err: %s", err.Error())
	}
}
