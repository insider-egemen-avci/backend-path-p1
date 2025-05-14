package logging

import (
	"log/slog"
	"os"
	"strings"
)

func Init(logLevel string) {
	var level slog.Level

	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		slog.Warn("Invalid log level, defaulting to info")
		level = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	handler := slog.NewTextHandler(os.Stdout, handlerOptions)

	slog.SetDefault(slog.New(handler))

	slog.Info("Logger initialized")
}
