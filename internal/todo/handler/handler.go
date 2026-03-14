package handler

import (
	"encoding/json"
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

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	todo, err := h.service.CreateToDo(req.Title, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {}
