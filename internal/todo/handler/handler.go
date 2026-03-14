package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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
		h.CreateTask(w, r)
	case http.MethodGet:
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 3 {
			h.GetTasks(w, r)
		} else {
			id, err := strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
			}
			h.GetTaskByID(w, r, id)
		}
	case http.MethodPut:
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 3 {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
		}
		h.UpdateTaskByID(w, r, id)
	case http.MethodDelete:
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 3 {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
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

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetToDos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request, id int) {}

func (h *Handler) UpdateTaskByID(w http.ResponseWriter, r *http.Request, id int) {}

func (h *Handler) DeleteTaskByID(w http.ResponseWriter, r *http.Request, id int) {}
