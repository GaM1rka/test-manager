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
		h.logger.WithError(err).Warn("Error while parsing json in CreateTask method")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":      "CreateTask",
		"title":       req.Title,
		"description": req.Description,
	}).Info("Creating new task")

	todo, err := h.service.CreateToDo(req.Title, req.Description)
	if err != nil {
		if errors.Is(err, validator.ErrInvalidTitle) {
			h.logger.WithError(err).Warn("Validation failed")
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			h.logger.WithError(err).Warn("Error while creating task")
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	h.logger.WithFields(logrus.Fields{
		"method": "GetTasks",
	}).Info("Getting all tasks")

	todos, err := h.service.GetToDos()
	if err != nil {
		h.logger.WithError(err).Warn("Error while getting tasks")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	h.logger.WithFields(logrus.Fields{
		"method": "GetTaskByID",
		"id":     id,
	}).Info("Getting task by ID")

	todo, err := h.service.GetToDoByID(id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			h.logger.WithError(err).Warn("Task not found")
			http.Error(w, "task not found", http.StatusNotFound)
		} else {
			h.logger.WithError(err).Error("Failed to get task")
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) UpdateTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	h.logger.WithFields(logrus.Fields{
		"method": "UpdateTaskByID",
		"id":     id,
	}).Info("Updating task")

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid JSON")
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	todo, err := h.service.UpdateToDo(id, req.Title, req.Description, req.Completed)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			h.logger.WithError(err).Warn("Task not found for update")
			http.Error(w, "task not found", http.StatusNotFound)
		} else if errors.Is(err, validator.ErrInvalidTitle) {
			h.logger.WithError(err).Warn("Validation failed")
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			h.logger.WithError(err).Error("Failed to update task")
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) DeleteTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	h.logger.WithFields(logrus.Fields{
		"method": "DeleteTaskByID",
		"id":     id,
	}).Info("Deleting task")

	err := h.service.DeleteToDo(id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			h.logger.WithError(err).Warn("Task not found for delete")
			http.Error(w, "task not found", http.StatusNotFound)
		} else {
			h.logger.WithError(err).Error("Failed to delete task")
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
