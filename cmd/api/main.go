package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"insider-egemen-avci/backend-path-p1/internal/config"
	"insider-egemen-avci/backend-path-p1/internal/logging"
	"insider-egemen-avci/backend-path-p1/internal/models"
	"insider-egemen-avci/backend-path-p1/internal/processing"
)

func ProcessTransaction(batch []models.Transaction) {
	slog.Info("Processing transaction batch", "size", len(batch))

	startTime := time.Now()
	numberOfJobs := len(batch)
	numberOfWorkers := 10

	pool := processing.NewPool(numberOfWorkers, numberOfJobs)
	pool.Start()

	var wg sync.WaitGroup
	wg.Add(numberOfJobs)

	go func() {
		for result := range pool.Results() {
			slog.Info("Transaction processed", "transaction_id", result.ID)
			wg.Done()
		}
	}()

	for _, transaction := range batch {
		pool.AddJob(transaction)
	}

	pool.CloseJobs()

	slog.Info("Waiting for all transactions to be processed")
	wg.Wait()

	elapsedTime := time.Since(startTime)
	slog.Info("Transaction batch processed in", "time", elapsedTime)
}

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
