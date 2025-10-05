package container

import (
	"context"
	"fmt"
	"net/http"

	"zpwoot/internal/adapters/database"
	"zpwoot/internal/adapters/database/repository"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/config"
	domainSession "zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"
)

type Container struct {
	config   *config.Config
	logger   *logger.Logger
	database *database.Database
	migrator *database.Migrator

	sessionService *domainSession.Service

	whatsappClient output.WhatsAppClient

	sessionUseCases input.SessionUseCases
	messageUseCases input.MessageUseCases
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{
		config: cfg,
	}
}

func (c *Container) Initialize() error {

	logger.Init(c.config.LogLevel)
	c.logger = logger.NewFromAppConfig(c.config)

	c.logger.Info().Msg("Initializing zpwoot container")

	c.logger.Info().Msg("Connecting to database")
	db, err := database.New(c.config, c.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	c.database = db
	c.logger.Info().Msg("Database connection established")

	c.logger.Info().Msg("Initializing database migrator")
	c.migrator = database.NewMigrator(db, c.logger)

	c.logger.Info().Msg("Running database migrations")
	if err := c.migrator.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	c.logger.Info().Msg("Database migrations completed")

	c.logger.Info().Msg("Initializing domain services")
	sessionRepo := repository.NewSessionRepository(c.database.DB)
	c.sessionService = domainSession.NewService(sessionRepo)

	c.logger.Info().Msg("Initializing use cases")
	c.initializeUseCases()

	c.logger.Info().Msg("Container initialization completed successfully")
	return nil
}

func (c *Container) GetDatabase() *database.Database {
	return c.database
}

func (c *Container) GetLogger() *logger.Logger {
	return c.logger
}

func (c *Container) GetMigrator() *database.Migrator {
	return c.migrator
}

func (c *Container) GetConfig() *config.Config {
	return c.config
}

func (c *Container) Start(ctx context.Context) error {
	return c.Initialize()
}

func (c *Container) Stop(ctx context.Context) error {
	if c.database != nil {
		return c.database.Close()
	}
	return nil
}

func (c *Container) Handler() http.Handler {
	mux := http.NewServeMux()

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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"zpwoot API is running","version":"1.0.0"}`))
	})

	return mux
}

func (c *Container) initializeUseCases() {

	c.sessionUseCases = c.createSessionUseCases()
	c.messageUseCases = c.createMessageUseCases()
}

func (c *Container) createSessionUseCases() input.SessionUseCases {

	return nil
}

func (c *Container) createMessageUseCases() input.MessageUseCases {

	return nil
}

func (c *Container) GetSessionService() *domainSession.Service {
	return c.sessionService
}

func (c *Container) GetSessionUseCases() input.SessionUseCases {
	return c.sessionUseCases
}

func (c *Container) GetMessageUseCases() input.MessageUseCases {
	return c.messageUseCases
}

func (c *Container) SetSessionUseCases(useCases input.SessionUseCases) {
	c.sessionUseCases = useCases
}

func (c *Container) SetMessageUseCases(useCases input.MessageUseCases) {
	c.messageUseCases = useCases
}

func (c *Container) SetWhatsAppClient(client output.WhatsAppClient) {
	c.whatsappClient = client
}
