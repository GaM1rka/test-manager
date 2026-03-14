package service

import (
	"test-manager/internal/todo/model"
	"test-manager/internal/todo/repository"
	validator "test-manager/pkg"

	"github.com/sirupsen/logrus"
)

type ToDoService struct {
	repo   *repository.ToDoRepository
	logger *logrus.Logger
}

func NewToDoService(repo *repository.ToDoRepository, l *logrus.Logger) *ToDoService {
	return &ToDoService{
		repo:   repo,
		logger: l,
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
