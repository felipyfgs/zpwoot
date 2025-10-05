package container

import (
	"context"
	"fmt"

	"zpwoot/internal/adapters/database"
	"zpwoot/internal/adapters/database/repository"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/adapters/waclient"
	"zpwoot/internal/config"
	"zpwoot/internal/core/application/usecase/message"
	"zpwoot/internal/core/application/usecase/session"
	domainSession "zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"
)

// Container is a simple dependency injection container
type Container struct {
	config   *config.Config
	logger   *logger.Logger
	database *database.Database
	migrator *database.Migrator

	// Domain services
	sessionService *domainSession.Service

	// External adapters
	whatsappClient output.WhatsAppClient

	// Use cases
	sessionUseCases input.SessionUseCases
	messageUseCases input.MessageUseCases
}

// NewContainer creates a new dependency injection container
func NewContainer(cfg *config.Config) *Container {
	return &Container{
		config: cfg,
	}
}

// Initialize sets up all dependencies
func (c *Container) Initialize() error {
	// Initialize logger
	logger.Init(c.config.LogLevel)
	c.logger = logger.NewFromAppConfig(c.config)
	c.logger.Info().Msg("Initializing zpwoot container")

	// Initialize database
	c.logger.Info().Msg("Connecting to database")
	db, err := database.New(c.config, c.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	c.database = db
	c.logger.Info().Msg("Database connection established")

	// Initialize migrator and run migrations
	c.logger.Info().Msg("Running database migrations")
	c.migrator = database.NewMigrator(db, c.logger)
	if err := c.migrator.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	c.logger.Info().Msg("Database migrations completed")

	// Initialize domain services
	c.logger.Info().Msg("Initializing domain services")
	sessionRepo := repository.NewSessionRepository(c.database.DB)
	c.sessionService = domainSession.NewService(sessionRepo)

	// Initialize WhatsApp client
	c.logger.Info().Msg("Initializing WhatsApp client")
	c.initializeWhatsAppClient()

	// Initialize use cases
	c.logger.Info().Msg("Initializing use cases")
	c.sessionUseCases = session.NewUseCases(c.sessionService, c.whatsappClient)
	c.messageUseCases = message.NewMessageUseCases(c.sessionService, c.whatsappClient)

	c.logger.Info().Msg("Container initialization completed successfully")
	return nil
}

// initializeWhatsAppClient sets up the WhatsApp client
func (c *Container) initializeWhatsAppClient() {
	sessionRepository := repository.NewSessionRepository(c.database.DB)
	sessionRepo := repository.NewSessionRepo(sessionRepository)
	waContainer := waclient.NewWAStoreContainer(
		c.database.DB,
		c.logger,
		c.config.Database.URL,
	)
	waClient := waclient.NewWAClient(waContainer, c.logger, sessionRepo)
	c.whatsappClient = waclient.NewWAClientAdapter(waClient)
}

// Start initializes the container
func (c *Container) Start(ctx context.Context) error {
	return c.Initialize()
}

// Stop shuts down the container
func (c *Container) Stop(ctx context.Context) error {
	if c.database != nil {
		return c.database.Close()
	}
	return nil
}

// Getters
func (c *Container) GetConfig() *config.Config {
	return c.config
}

func (c *Container) GetLogger() *logger.Logger {
	return c.logger
}

func (c *Container) GetDatabase() *database.Database {
	return c.database
}

func (c *Container) GetMigrator() *database.Migrator {
	return c.migrator
}

func (c *Container) GetSessionService() *domainSession.Service {
	return c.sessionService
}

func (c *Container) GetWhatsAppClient() output.WhatsAppClient {
	return c.whatsappClient
}

func (c *Container) GetSessionUseCases() input.SessionUseCases {
	return c.sessionUseCases
}

func (c *Container) GetMessageUseCases() input.MessageUseCases {
	return c.messageUseCases
}


