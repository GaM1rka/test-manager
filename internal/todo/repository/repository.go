package repository

import (
	"sync"
	"test-manager/internal/todo/model"
)

type ToDoRepository struct {
	todos  map[int]*model.ToDo
	mu     sync.Mutex
	nextID int
}

func NewToDoRepository() *ToDoRepository {
	return &ToDoRepository{
		todos:  make(map[int]*model.ToDo),
		nextID: 1,
	}
}
