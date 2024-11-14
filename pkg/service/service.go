package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/egigiffari/go-grpc-todo/pb"
	"github.com/egigiffari/go-grpc-todo/pkg/store"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Store store.Store
	pb.UnimplementedTodoServiceServer
}

func New(store store.Store) *Service {
	return &Service{Store: store}
}

func (s *Service) Create(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	todo := req.GetTodo()
	fmt.Printf("todo with id: %s\n", todo.Id)

	if err := s.validate(todo); err != nil {
		return nil, err
	}

	if len(todo.Id) <= 0 {
		todo.Id = uuid.NewString()
	}

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("request timeout")
		return nil, status.Error(codes.DeadlineExceeded, "request timeout")
	}

	if ctx.Err() == context.Canceled {
		fmt.Println("request canceled")
		return nil, status.Error(codes.Canceled, "request canceled")
	}

	if err := s.Store.Save(todo); err != nil {
		code := codes.Internal
		if errors.Is(err, store.ErrAlreadyExists) {
			code = codes.AlreadyExists
		}

		return nil, status.Error(code, err.Error())
	}

	fmt.Printf("success to save todo with id: %s\n", todo.Id)
	res := &pb.CreateTodoResponse{Id: todo.Id}
	return res, nil
}

func (s *Service) validate(todo *pb.Todo) error {
	if len(todo.Id) > 0 {
		_, err := uuid.Parse(todo.Id)
		if err != nil {
			return status.Error(codes.InvalidArgument, "Invalid Todo id")
		}
	}

	if len(todo.Name) < 3 {
		return status.Error(codes.InvalidArgument, "Required Todo name")
	}

	return nil
}
