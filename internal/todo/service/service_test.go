package service

import (
	"errors"
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
