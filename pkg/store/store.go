package store

import (
	"errors"
	"fmt"
	"sync"

	"github.com/egigiffari/go-grpc-todo/pb"
	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")

type Store interface {
	Save(todo *pb.Todo) error
	Read(id string) (*pb.Todo, error)
}

type InMemoryStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Todo
}

func NewInMemory() Store {
	return &InMemoryStore{
		data: make(map[string]*pb.Todo),
	}
}

func (im *InMemoryStore) Save(todo *pb.Todo) error {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	if _, ok := im.data[todo.Id]; ok {
		return ErrAlreadyExists
	}

	data := &pb.Todo{}
	if err := copier.Copy(data, todo); err != nil {
		return fmt.Errorf("failed to copy data cause: %s", err.Error())
	}

	im.data[data.Id] = data
	return nil
}

func (im *InMemoryStore) Read(id string) (*pb.Todo, error) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	data, ok := im.data[id]
	if !ok {
		return nil, errors.New("not found data")
	}

	todo := &pb.Todo{}
	if err := copier.Copy(todo, data); err != nil {
		return nil, errors.New("failed to copy data")
	}

	return data, nil
}
