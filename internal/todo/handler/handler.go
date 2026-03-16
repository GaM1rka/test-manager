package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"test-manager/internal/todo/service"
	validator "test-manager/pkg"

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
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	switch r.Method {

	case http.MethodPost:
		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}
		h.CreateTask(w, r)

	case http.MethodGet:
		if len(parts) == 1 {
			h.GetTasks(w, r)
			return
		}
		if len(parts) != 2 {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		h.GetTaskByID(w, r, id)

	case http.MethodPut:
		if len(parts) != 2 {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		h.UpdateTaskByID(w, r, id)

	case http.MethodDelete:
		if len(parts) != 2 {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		h.DeleteTaskByID(w, r, id)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Error while parsing json in CreateTask method")
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	todo, err := h.service.CreateToDo(req.Title, req.Description)
	if err != nil {
		if errors.Is(err, validator.ErrInvalidTitle) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetToDos()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	todo, err := h.service.GetToDoByID(id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) UpdateTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	todo, err := h.service.UpdateToDo(id, req.Title, req.Description, req.Completed)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
		} else if errors.Is(err, validator.ErrInvalidTitle) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) DeleteTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	err := h.service.DeleteToDo(id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
