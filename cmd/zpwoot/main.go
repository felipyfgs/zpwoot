// Package main zpwoot WhatsApp API Gateway
//
//	@title			zpwoot WhatsApp API
//	@version		1.0.0
//	@description	A comprehensive WhatsApp Business API built with Go, following Clean Architecture principles.
//	@description	Provides endpoints for session management, messaging, contacts, groups, media handling, and integrations.
//
//	@contact.name	zpwoot Support
//	@contact.url	https://github.com/your-org/zpwoot
//	@contact.email	support@zpwoot.com
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@BasePath	/
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description API key for authentication. Use the value from .env file (API_KEY variable). Send as 'Authorization: your-api-key' (without Bearer prefix)
//
//	@tag.name			Sessions
//	@tag.description	WhatsApp session management operations
//
//	@tag.name			Messages
//	@tag.description	Message sending and retrieval operations
//
//	@tag.name			Contacts
//	@tag.description	Contact management operations
//
//	@tag.name			Health
//	@tag.description	Health check and system status
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

	// Initialize logger with full configuration
	logger.InitWithConfig(cfg)

	// Use the centralized logger following the zerolog pattern
	logger.WithComponent("main").Info().Msg("Starting zpwoot application")

	// Initialize dependencies container
	container := container.NewContainer(cfg)

	// Initialize container (this will run migrations automatically)
	if err := container.Initialize(); err != nil {
		logger.WithComponent("main").Fatal().
			Err(err).
			Msg("Failed to initialize container")
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
		logger.WithComponent("server").Info().
			Str("address", cfg.GetServerAddress()).
			Msg("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithComponent("server").Fatal().
				Err(err).
				Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.WithComponent("main").Info().Msg("Shutting down server")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.WithComponent("main").Fatal().
			Err(err).
			Msg("Server forced to shutdown")
	}

	// Stop container
	if err := container.Stop(ctx); err != nil {
		logger.WithComponent("main").Error().
			Err(err).
			Msg("Error stopping container")
	}

	logger.WithComponent("main").Info().Msg("Server exited")
}
