package config

import (
	"time"
)

type Config struct {
	ServerPort      int
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{}
}
