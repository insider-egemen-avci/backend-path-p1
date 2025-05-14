package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"insider-egemen-avci/backend-path-p1/internal/config"
	"insider-egemen-avci/backend-path-p1/internal/logging"
)

func main() {
	fmt.Println("Hello, World!")

	config, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	logging.Init(config.LogLevel)

	slog.Info("Server is starting...",
		"port", config.Port,
		"logLevel", config.LogLevel,
		"shutdownTimeout", config.ShutdownTimeout,
	)

	if config.DSN == "" {
		slog.Warn("DSN is not set")
	}

	slog.Info("Application is running on", "port", config.Port)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))

		slog.Info("Request received", "method", r.Method, "path", r.URL.Path)

		time.Sleep(3 * time.Second)

		slog.Info("Request processed", "method", r.Method, "path", r.URL.Path)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("Starting server", "port", config.Port)
		serverErrors <- server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		slog.Error("Server error", "error", err)

		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("Server closed")
		} else {
			slog.Error("Server error", "error", err)
		}
	case sig := <-quit:
		slog.Info("Server is shutting down...", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			slog.Error("Server shutdown error", "error", err)
		}
	}

	slog.Info("Server shutdown complete")
}
