package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/egigiffari/go-grpc-todo/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	addr := flag.String("address", "", "the server address")
	flag.Parse()
	fmt.Printf("dial server %s\n", *addr)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(*addr, opts...)
	if err != nil {
		log.Printf("failed to dial with server err: %s", err.Error())
	}

	service := pb.NewTodoServiceClient(conn)

	todo := &pb.Todo{
		Id:   uuid.NewString(),
		Name: "morning yoga",
		Done: false,
	}
	req := &pb.CreateTodoRequest{Todo: todo}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := service.Create(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("todo already exists")
		} else {
			log.Printf("connot to create cause, err: %s", err.Error())
		}

		return
	}

	log.Printf("success create todo with id = %s", res.Id)

}
