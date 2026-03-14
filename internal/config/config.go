package config

import (
	"os"
	"time"
)

type Config struct {
	ServerPort      string
	ShutdownTimeout time.Duration
}

func Load() Config {
	portStr := os.Getenv("SERVER_PORT")
	timeoutStr := os.Getenv("SHUTDOWN_TIMEOUT")

	port := portStr
	if port == "" {
		port = ":8080"
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		timeout = time.Duration(30 * time.Second)
	}

	return Config{
		ServerPort:      port,
		ShutdownTimeout: timeout,
	}
}
