//	@title			zpwoot WhatsApp API
//	@version		1.0.0
//	@description	A comprehensive WhatsApp Business API built with Go, following Clean Architecture principles.
//	@description	Provides endpoints for session management, messaging, contacts, groups, media handling, and integrations.

//	@contact.name	zpwoot Support
//	@contact.url	https://github.com/your-org/zpwoot
//	@contact.email	support@zpwoot.com

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@BasePath	/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description API key for authentication. Use the value from .env file (API_KEY variable). Send as 'Authorization: your-api-key' (without Bearer prefix)

//	@tag.name			Sessions
//	@tag.description	WhatsApp session management operations

//	@tag.name			Messages
//	@tag.description	Message sending and retrieval operations

//	@tag.name			Contacts
//	@tag.description	Contact management operations

//	@tag.name			Groups
//	@tag.description	WhatsApp group management operations

//	@tag.name			Presence
//	@tag.description	Presence and status management operations

//	@tag.name			Webhooks
//	@tag.description	Webhook configuration and event management

// @tag.name			Health
// @tag.description	Health check and system status
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zpwoot/internal/adapters/http/router"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/config"
	"zpwoot/internal/container"

	_ "zpwoot/docs/swagger"
)

func main() {
	cfg := config.Load()

	logger.InitWithConfig(cfg)

	logger.WithComponent("main").Info().Msg("Starting zpwoot application")

	container := container.NewContainer(cfg)

	if err := container.Init(); err != nil {
		logger.WithComponent("main").Fatal().
			Err(err).
			Msg("Failed to initialize container")
	}

	r := router.NewRouter(container)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.WithComponent("main").Info().Msg("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithComponent("main").Error().
			Err(err).
			Msg("Server forced to shutdown")
	} else {
		logger.WithComponent("main").Info().Msg("Server shutdown completed")
	}

	if err := container.Stop(ctx); err != nil {
		logger.WithComponent("main").Error().
			Err(err).
			Msg("Error stopping container")
	}

	logger.WithComponent("main").Info().Msg("Server exited")
}
