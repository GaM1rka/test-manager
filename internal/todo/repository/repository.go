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

func (r *ToDoRepository) Create(todo *model.ToDo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo.ID = r.nextID
	r.nextID++
	r.todos[todo.ID] = todo
	return nil
}
