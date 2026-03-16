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

func TestGetToDoByID(t *testing.T) {
	repo := repository.NewToDoRepository()
	s := NewToDoService(repo)

	todo1 := &model.ToDo{Title: "Title", Description: "Description"}
	repo.Create(todo1)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{name: "found", id: todo1.ID, wantErr: false},
		{name: "not found", id: 67, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetToDoByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToDoByID(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
			if !tt.wantErr && got.ID != tt.id {
				t.Errorf("GetToDoByID(%d) wrong todo = %+v", tt.id, got)
			}
			if tt.wantErr && !errors.Is(err, ErrTaskNotFound) {
				t.Errorf("GetToDoByID(%d) wrong error: %v", tt.id, err)
			}
		})
	}
}

func TestUpdateToDo(t *testing.T) {
	repo := repository.NewToDoRepository()
	s := NewToDoService(repo)

	originalTodo := &model.ToDo{Title: "Title", Description: "Description"}
	repo.Create(originalTodo)

	tests := []struct {
		name        string
		id          int
		title       string
		description string
		completed   bool
		wantErr     bool
	}{
		{
			name:        "successful update",
			id:          originalTodo.ID,
			title:       "Updated Title",
			description: "Updated Desc",
			completed:   true,
			wantErr:     false,
		},
		{
			name:        "validation error",
			id:          originalTodo.ID,
			title:       "",
			description: "Desc",
			completed:   false,
			wantErr:     true,
		},
		{
			name:        "not found",
			id:          52,
			title:       "New Title",
			description: "New Desc",
			completed:   false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.UpdateToDo(tt.id, tt.title, tt.description, tt.completed)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateToDo(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}

			if !tt.wantErr {
				if got.Title != tt.title || got.Description != tt.description || got.Completed != tt.completed {
					t.Errorf("UpdateToDo(%d) wrong result: %+v", tt.id, got)
				}
			}
		})
	}
}

func TestDeleteToDo(t *testing.T) {
	repo := repository.NewToDoRepository()
	s := NewToDoService(repo)

	todo := &model.ToDo{Title: "Task", Description: "Description"}
	repo.Create(todo)

	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{
			name:    "successful delete",
			id:      todo.ID,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      123,
			wantErr: ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteToDo(tt.id)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				_, err := repo.GetByID(tt.id)
				if err != nil {
					t.Fatalf("repo error: %v", err)
				}
				todo, _ := repo.GetByID(tt.id)
				if todo != nil {
					t.Fatalf("todo was not deleted")
				}
			}
		})
	}
}
