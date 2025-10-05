package main

import (
	"log"
	"net/http"

	"github.com/zpwoot/internal/infra/http/router"
	"github.com/zpwoot/platform/config"
	"github.com/zpwoot/platform/container"
	"github.com/zpwoot/platform/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Initialize dependencies container
	container := container.NewContainer(cfg)

	// Setup HTTP router
	r := router.NewRouter(container)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
