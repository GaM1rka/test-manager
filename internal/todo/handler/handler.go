package handler

import (
	"net/http"
	"test-manager/internal/todo/service"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	service *service.ToDoService
	logger  *logrus.Logger
}

func NewHandler(s *service.ToDoService, l *logrus.Logger) *Handler {
	return &Handler{
		service: s,
		logger:  l,
	}
}

func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
	case http.MethodGet:
	case http.MethodPut:
	case http.MethodDelete:
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {}
