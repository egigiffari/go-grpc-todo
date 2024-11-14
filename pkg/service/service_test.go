package service_test

import (
	"context"
	"testing"

	"github.com/egigiffari/go-grpc-todo/pb"
	"github.com/egigiffari/go-grpc-todo/pkg/service"
	"github.com/egigiffari/go-grpc-todo/pkg/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServiceCreateTodo(t *testing.T) {
	t.Parallel()

	todoNoId := todoRequest()
	todoNoId.Id = ""

	todoInvalidId := todoRequest()
	todoInvalidId.Id = "invalid-id"

	todoInvalidName := todoRequest()
	todoInvalidName.Name = "as"

	todoDuplicate := todoRequest()
	duplicateStore := store.NewInMemory()
	err := duplicateStore.Save(todoDuplicate)
	require.NoError(t, err)

	tt := []struct {
		name  string
		todo  *pb.Todo
		store store.Store
		want  codes.Code
	}{
		{
			name:  "success",
			todo:  todoRequest(),
			store: store.NewInMemory(),
			want:  codes.OK,
		},
		{
			name:  "success_no_id",
			todo:  todoNoId,
			store: store.NewInMemory(),
			want:  codes.OK,
		},
		{
			name:  "failed_invalid_id",
			todo:  todoInvalidId,
			store: store.NewInMemory(),
			want:  codes.InvalidArgument,
		},
		{
			name:  "failed_invalid_name",
			todo:  todoInvalidName,
			store: store.NewInMemory(),
			want:  codes.InvalidArgument,
		},
		{
			name:  "failed_alredy_exists",
			todo:  todoDuplicate,
			store: duplicateStore,
			want:  codes.AlreadyExists,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := &pb.CreateTodoRequest{Todo: tc.todo}

			service := service.New(tc.store)
			res, err := service.Create(context.Background(), req)
			if tc.want == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)
				if len(tc.todo.Id) > 0 {
					require.Equal(t, res.Id, tc.todo.Id)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, st.Code(), tc.want)
			}
		})
	}
}

func todoRequest() *pb.Todo {
	todo := &pb.Todo{
		Id:   uuid.NewString(),
		Name: "morning yoga",
		Done: false,
	}

	return todo
}
