package main

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	dockeradapter "svitlopus/internal/adapters/docker_adapter"
	"svitlopus/internal/api"
	"svitlopus/internal/config"
	"svitlopus/internal/database"
	logger "svitlopus/internal/logger"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
	}
	logFile, err := logger.Setup(cfg.Log.LogPath)
	if err != nil {
		slog.Error("failed to setup logger", slog.String("error", err.Error()))
	}
	defer logFile.Close()

	db, err := database.New(&cfg.Database)
	if err != nil {
		slog.Error("failed to initialize database", slog.String("error", err.Error()))
	}
	defer db.Close()

	dockeru := dockeradapter.NewDockerUtil(&cfg.Docker)
	if err := dockeru.RunDockerPipeline(); err != nil {
		slog.Error("failed to run docker container", slog.String("error", err.Error()))
	}

	server := api.NewServer(cfg, db)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	slog.Info("server is running")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("http server error", slog.String("error", err.Error()))
	}

	// Wait for the graceful shutdown to complete
	<-done
	slog.Info("Graceful shutdown complete.")
}
