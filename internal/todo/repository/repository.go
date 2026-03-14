package repository

import (
	"fmt"
	"sync"
	"test-manager/internal/todo/model"
)

type ToDoRepository struct {
	todos  map[int]*model.ToDo
	mu     sync.RWMutex
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

func (r *ToDoRepository) GetAll() ([]*model.ToDo, error) {
	r.mu.RLock()
	defer r.mu.Unlock()

	todos := make([]*model.ToDo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *ToDoRepository) GetByID(id int) (*model.ToDo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, nil
	}
	return todo, nil
}

func (r *ToDoRepository) Update(todo *model.ToDo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[todo.ID]; !exists {
		return fmt.Errorf("todo %d not found", todo.ID)
	}

	r.todos[todo.ID] = todo
	return nil
}

func (r *ToDoRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[id]; !exists {
		return fmt.Errorf("todo %d not found", id)
	}

	delete(r.todos, id)
	return nil
}
