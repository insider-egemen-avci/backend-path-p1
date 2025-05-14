package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            int
	DSN             string
	LogLevel        string
	ShutdownTimeout time.Duration
}

func LoadConfig() (*Config, error) {
	config := Config{
		Port:            8080,
		DSN:             "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable",
		LogLevel:        "info",
		ShutdownTimeout: 10 * time.Second,
	}

	var err error

	portStr := os.Getenv("PORT")

	if portStr != "" {
		port, convErr := strconv.Atoi(portStr)
		if convErr != nil {
			return nil, fmt.Errorf("failed to convert PORT %s to int: %w", portStr, convErr)
		}
		config.Port = port
	}

	dsn := os.Getenv("DSN")

	if dsn != "" {
		config.DSN = dsn
	} else {
		return nil, fmt.Errorf("DSN is not set")
	}

	level := os.Getenv("LOG_LEVEL")

	if level != "" {
		config.LogLevel = level
	} else {
		return nil, fmt.Errorf("LOG_LEVEL is not set")
	}

	shutdownTimeoutStr := os.Getenv("SHUTDOWN_TIMEOUT")

	if shutdownTimeoutStr != "" {
		shutdownTimeout, convErr := time.ParseDuration(shutdownTimeoutStr)
		if convErr != nil {
			return nil, fmt.Errorf("failed to convert SHUTDOWN_TIMEOUT %s to time.Duration: %w", shutdownTimeoutStr, convErr)
		}

		config.ShutdownTimeout = shutdownTimeout
	}

	return &config, err
}
