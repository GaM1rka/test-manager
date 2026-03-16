package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-manager/internal/config"
	"test-manager/internal/todo/handler"
	"test-manager/internal/todo/repository"
	"test-manager/internal/todo/service"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	cfg := config.Load()

	todoRepo := repository.NewToDoRepository()
	todoService := service.NewToDoService(todoRepo)
	h := handler.NewHandler(todoService, logger)

	http.HandleFunc("/todos", h.TaskHandler)

	server := &http.Server{
		Addr: cfg.ServerPort,
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("HTTP server error: %v", err)
		}
		logger.Info("Stopped serving new connections")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, stop := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer stop()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("HTTP shutdown error: %v", err)
	}
	logger.Info("Graceful shutdown complete.")
}
