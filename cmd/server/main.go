package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sfdeloach/churchsite/internal/config"
	"github.com/sfdeloach/churchsite/internal/database"
	"github.com/sfdeloach/churchsite/internal/handlers"
	"github.com/sfdeloach/churchsite/internal/services"
)

func main() {
	// Structured JSON logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Handle subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			runMigrate()
			return
		default:
			slog.Error("unknown command", "command", os.Args[1])
			os.Exit(1)
		}
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if cfg.IsDevelopment() {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)
	}

	// Connect to databases
	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize services
	siteSettingsSvc := services.NewSiteSettingsService(db.Postgres)
	eventSvc := services.NewEventService(db.Postgres)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)
	homeHandler := handlers.NewHomeHandler(siteSettingsSvc, eventSvc)

	// Build router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	// Static files
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Health checks
	r.Get("/health", healthHandler.Liveness)
	r.Get("/health/ready", healthHandler.Readiness)

	// Pages
	r.Get("/", homeHandler.Index)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "addr", addr, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}

func runMigrate() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		slog.Error("usage: sachapel migrate [up|down]")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "up":
		if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
			slog.Error("migration failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migrations applied successfully")
	case "down":
		if err := database.RollbackMigration(cfg.DatabaseURL); err != nil {
			slog.Error("rollback failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migration rolled back successfully")
	default:
		slog.Error("unknown migrate command, use 'up' or 'down'", "command", os.Args[2])
		os.Exit(1)
	}
}
