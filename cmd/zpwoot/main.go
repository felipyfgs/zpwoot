package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zpwoot/internal/adapters/config"
	"zpwoot/internal/adapters/container"
	"zpwoot/internal/adapters/http/router"
	"zpwoot/internal/adapters/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.LogLevel)
	log := logger.NewFromAppConfig(cfg)

	log.Info("Starting zpwoot application...")

	// Initialize dependencies container
	container := container.NewContainer(cfg)

	// Initialize container (this will run migrations automatically)
	if err := container.Initialize(); err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// Setup HTTP router
	r := router.NewRouter(container)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("Starting server on %s", cfg.GetServerAddress())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Stop container
	if err := container.Stop(ctx); err != nil {
		log.Errorf("Error stopping container: %v", err)
	}

	log.Info("Server exited")
}
