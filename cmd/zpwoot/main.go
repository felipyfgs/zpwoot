package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zpwoot/internal/infra/http/router"
	"zpwoot/platform/config"
	"zpwoot/platform/container"
	"zpwoot/platform/database"
	"zpwoot/platform/logger"
)

func main() {
	// Parse command line flags
	var (
		migrateUp     = flag.Bool("migrate-up", false, "Run database migrations up")
		migrateDown   = flag.Bool("migrate-down", false, "Run database migrations down (rollback last)")
		migrateStatus = flag.Bool("migrate-status", false, "Show migration status")
	)
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.LogLevel)
	log := logger.NewFromAppConfig(cfg)

	log.Info("Starting zpwoot application...")

	// Handle migration commands
	if *migrateUp || *migrateDown || *migrateStatus {
		handleMigrationCommands(cfg, log, *migrateUp, *migrateDown, *migrateStatus)
		return
	}

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

func handleMigrationCommands(cfg *config.Config, log *logger.Logger, migrateUp, migrateDown, migrateStatus bool) {
	// Initialize database connection
	db, err := database.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize migrator
	migrator := database.NewMigrator(db, log)

	switch {
	case migrateUp:
		log.Info("Running migrations up...")
		if err := migrator.RunMigrations(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Info("Migrations completed successfully")

	case migrateDown:
		log.Info("Rolling back last migration...")
		if err := migrator.Rollback(); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		log.Info("Rollback completed successfully")

	case migrateStatus:
		log.Info("Checking migration status...")
		migrations, err := migrator.GetMigrationStatus()
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

		fmt.Println("\nMigration Status:")
		fmt.Println("=================")
		for _, migration := range migrations {
			status := "PENDING"
			if migration.AppliedAt != nil {
				status = "APPLIED"
			}
			fmt.Printf("Version %d: %s [%s]\n", migration.Version, migration.Name, status)
		}
		fmt.Println()
	}
}
