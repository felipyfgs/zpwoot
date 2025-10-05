package container

import (
	"context"
	"fmt"
	"net/http"

	"zpwoot/platform/config"
	"zpwoot/platform/database"
	"zpwoot/platform/logger"
)

// Container holds all application dependencies
type Container struct {
	config   *config.Config
	logger   *logger.Logger
	database *database.Database
	migrator *database.Migrator
}

// NewContainer creates a new dependency injection container
func NewContainer(cfg *config.Config) *Container {
	return &Container{
		config: cfg,
	}
}

// Initialize initializes all dependencies
func (c *Container) Initialize() error {
	// Initialize logger
	logger.Init(c.config.LogLevel)
	c.logger = logger.NewFromAppConfig(c.config)

	// Initialize database
	db, err := database.New(c.config, c.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	c.database = db

	// Initialize migrator
	c.migrator = database.NewMigrator(db, c.logger)

	// Run migrations automatically
	if err := c.migrator.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// GetDatabase returns the database instance
func (c *Container) GetDatabase() *database.Database {
	return c.database
}

// GetLogger returns the logger instance
func (c *Container) GetLogger() *logger.Logger {
	return c.logger
}

// GetMigrator returns the migrator instance
func (c *Container) GetMigrator() *database.Migrator {
	return c.migrator
}

// GetConfig returns the config instance
func (c *Container) GetConfig() *config.Config {
	return c.config
}

// Start starts all container components
func (c *Container) Start(ctx context.Context) error {
	return c.Initialize()
}

// Stop stops all container components
func (c *Container) Stop(ctx context.Context) error {
	if c.database != nil {
		return c.database.Close()
	}
	return nil
}

// Handler returns the HTTP handler (placeholder for now)
func (c *Container) Handler() http.Handler {
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if c.database != nil {
			if err := c.database.Health(); err != nil {
				http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
				return
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"zpwoot"}`))
	})

	// Basic info endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"zpwoot API is running","version":"1.0.0"}`))
	})

	return mux
}
