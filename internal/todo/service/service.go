package service

import (
	"errors"
	"fmt"
	"sort"
	"test-manager/internal/todo/model"
	"test-manager/internal/todo/repository"
	validator "test-manager/pkg"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type ToDoService struct {
	repo *repository.ToDoRepository
}

func NewToDoService(repo *repository.ToDoRepository) *ToDoService {
	return &ToDoService{
		repo: repo,
	}
}

func (s *ToDoService) CreateToDo(title, description string) (*model.ToDo, error) {
	if err := validator.ValidateTodo(title); err != nil {
		return nil, err
	}

	todo := &model.ToDo{
		Title:       title,
		Description: description,
		Completed:   false,
	}

	if err := s.repo.Create(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *ToDoService) GetToDos() ([]*model.ToDo, error) {
	todos, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	return todos, nil
}

func (s *ToDoService) GetToDoByID(id int) (*model.ToDo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo %d: %w", id, err)
	}
	if todo == nil {
		return nil, ErrTaskNotFound
	}
	return todo, nil
}

func (s *ToDoService) UpdateToDo(id int, title, description string, completed bool) (*model.ToDo, error) {
	if err := validator.ValidateTodo(title); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo %d: %w", id, err)
	}
	if todo == nil {
		return nil, ErrTaskNotFound
	}

	todo.Title = title
	todo.Description = description
	todo.Completed = completed

	if err := s.repo.Update(todo); err != nil {
		return nil, fmt.Errorf("failed to update todo %d: %w", id, err)
	}

	return todo, nil
}

func (s *ToDoService) DeleteToDo(id int) error {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if todo == nil {
		return ErrTaskNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
