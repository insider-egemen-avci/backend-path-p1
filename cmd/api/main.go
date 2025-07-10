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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/secure"

	"insider-egemen-avci/backend-path-p1/internal/config"
	"insider-egemen-avci/backend-path-p1/internal/logging"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	logging.Init(cfg.LogLevel)

	slog.Info("Server starting up...", "port", cfg.Port)

	r := chi.NewRouter()

	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true, // Prevents clickjacking
		ContentTypeNosniff:    true, // Prevents MIME type sniffing
		BrowserXssFilter:      true, // Adds XSS protection
		ContentSecurityPolicy: "default-src 'self'",
	})
	r.Use(secureMiddleware.Handler)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	/* api handler
		r.Route("/api/v1", func(r chi.Router) {
			r.Post("/auth/register", apiHandler.Register)
	    	r.Post("/auth/login", apiHandler.Login)
			r.Post("/auth/refresh", apiHandler.Refresh)
			and other routes...
		})
	*/

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Server is ready to handle requests", "address", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	case sig := <-quit:
		slog.Info("Shutdown signal received", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Server shutdown failed", "error", err)
		}
	}
	slog.Info("Server shutdown complete")
}
