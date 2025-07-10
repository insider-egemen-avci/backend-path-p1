package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	Port            int
	DSN             string
	LogLevel        string
	ShutdownTimeout time.Duration
	JWTSecret       string
}

func LoadConfig() (*AppConfig, error) {
	cfg := AppConfig{
		Port:            8080,
		LogLevel:        "info",
		ShutdownTimeout: 10 * time.Second,
	}

	if portStr := os.Getenv("PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT value: %w", err)
		}
		cfg.Port = port
	}

	if dsn := os.Getenv("DSN"); dsn != "" {
		cfg.DSN = dsn
	} else {
		return nil, fmt.Errorf("environment variable DSN must be set")
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.LogLevel = level
	}

	if timeoutStr := os.Getenv("SHUTDOWN_TIMEOUT"); timeoutStr != "" {
		shutdownTimeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SHUTDOWN_TIMEOUT value: %w", err)
		}
		cfg.ShutdownTimeout = shutdownTimeout
	}

	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWTSecret = secret
	} else {
		return nil, fmt.Errorf("environment variable JWT_SECRET must be set")
	}

	return &cfg, nil
}
