package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"test-manager/internal/todo/repository"
	"test-manager/internal/todo/service"
	"testing"

	"github.com/sirupsen/logrus"
)

func setup() *Handler {
	repo := repository.NewToDoRepository()
	s := service.NewToDoService(repo)
	logger := logrus.New()
	return NewHandler(s, logger)
}

func TestCreateTask(t *testing.T) {
	h := setup()

	body := map[string]string{
		"title":       "task1",
		"description": "desc",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.TaskHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 got %d", w.Code)
	}
}

func TestCreateTaskInvalidJSON(t *testing.T) {
	h := setup()

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer([]byte("bad json")))
	w := httptest.NewRecorder()

	h.TaskHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestGetTasks(t *testing.T) {
	h := setup()

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	h.TaskHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestGetTaskByIDNotFound(t *testing.T) {
	h := setup()

	req := httptest.NewRequest(http.MethodGet, "/todos/67", nil)
	w := httptest.NewRecorder()

	h.TaskHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", w.Code)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	h := setup()

	body := map[string]any{
		"title":       "new",
		"description": "desc",
		"completed":   true,
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/todos/67", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	h.TaskHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", w.Code)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	h := setup()

	req := httptest.NewRequest(http.MethodDelete, "/todos/67", nil)
	w := httptest.NewRecorder()

	h.TaskHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", w.Code)
	}
}
