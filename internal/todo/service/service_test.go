package service

import (
	"errors"
	"test-manager/internal/todo/model"
	"test-manager/internal/todo/repository"
	validator "test-manager/pkg"
	"testing"
)

func TestCreateToDo(t *testing.T) {
	repo := repository.NewToDoRepository()
	s := NewToDoService(repo)

	_, err := s.CreateToDo("Valid Title", "Description")
	if err != nil {
		t.Errorf("CreateToDo() err = %v, expected no error", err)
	}

	_, err = s.CreateToDo("", "Description")
	if !errors.Is(err, validator.ErrInvalidTitle) {
		t.Errorf("CreateToDo() expected validation error, got %v", err)
	}

	_, err = s.CreateToDo("Second Valid Title", "Description")
	if err != nil {
		t.Errorf("CreateToDo() second call failed = %v", err)
	}
}

func TestGetToDos(t *testing.T) {
	repo := repository.NewToDoRepository()
	s := NewToDoService(repo)

	repo.Create(&model.ToDo{Title: "Task 2", Description: "Desc2"})
	repo.Create(&model.ToDo{Title: "Task 1", Description: "Desc1"})

	todos, err := s.GetToDos()
	if err != nil {
		t.Errorf("GetToDos() error = %v, expected = nil", err)
	}

	if len(todos) != 2 {
		t.Errorf("GetToDos() wrong length = %d, expected = %d", len(todos), 2)
	}

	if todos[0].ID >= todos[1].ID {
		t.Errorf("GetToDos() wrong sort order: %v", todos)
	}
}
