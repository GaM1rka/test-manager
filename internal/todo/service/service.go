package service

import (
	"test-manager/internal/todo/repository"

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
