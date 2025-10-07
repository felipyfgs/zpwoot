package container

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"zpwoot/internal/adapters/database"
	"zpwoot/internal/adapters/database/repository"
	"zpwoot/internal/adapters/integration/webhook"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/adapters/waclient"
	"zpwoot/internal/config"
	"zpwoot/internal/core/application/usecase/message"
	"zpwoot/internal/core/application/usecase/session"
	webhookUseCase "zpwoot/internal/core/application/usecase/webhook"
	domainSession "zpwoot/internal/core/domain/session"
	domainWebhook "zpwoot/internal/core/domain/webhook"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"
)

type Container struct {
	config   *config.Config
	logger   *logger.Logger
	database *database.Database
	migrator *database.Migrator

	sessionService *domainSession.Service
	webhookService *domainWebhook.Service

	whatsappClient output.WhatsAppClient
	webhookSender  output.WebhookSender

	sessionUseCases input.SessionUseCases
	messageUseCases input.MessageUseCases
	webhookUseCases input.WebhookUseCases
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{
		config: cfg,
	}
}

func (c *Container) Init() error {
	return c.InitWithContext(context.Background())
}

func (c *Container) InitWithContext(ctx context.Context) error {
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

	c.logger.Info().Msg("Running database migrations")

	c.migrator = database.NewMigrator(db, c.logger)
	if err := c.migrator.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	c.logger.Info().Msg("Database migrations completed")

	c.logger.Info().Msg("Initializing domain services")
	sessionRepo := repository.NewSessionRepository(c.database.DB)
	c.sessionService = domainSession.NewService(sessionRepo)
	c.webhookService = domainWebhook.NewService()

	c.logger.Info().Msg("Initializing webhook sender")
	c.initWebhookSender()

	c.logger.Info().Msg("Initializing WhatsApp client")
	c.initWAClient()

	c.logger.Info().Msg("Initializing use cases")
	c.sessionUseCases = session.NewUseCases(c.sessionService, c.whatsappClient, c.logger)
	c.messageUseCases = message.NewUseCases(c.sessionService, c.whatsappClient, c.logger)
	c.webhookUseCases = c.initWebhookUseCases()

	c.logger.Info().Msg("Container initialization completed successfully")

	return nil
}

func (c *Container) initWAClient() {
	sessionRepository := repository.NewSessionRepository(c.database.DB)
	sessionRepo := repository.NewSessionRepo(sessionRepository)
	webhookRepo := repository.NewWebhookRepository(c.database.DB)

	waContainer := waclient.NewWAStoreContainer(
		c.database.DB,
		c.logger,
		c.config.Database.URL,
	)
	waClient := waclient.NewWAClient(waContainer, c.logger, sessionRepo, c.webhookSender, webhookRepo)
	c.whatsappClient = waclient.NewWAClientAdapter(waClient)
}

func (c *Container) Start(ctx context.Context) error {
	return c.InitWithContext(ctx)
}

func (c *Container) Stop(ctx context.Context) error {
	if c.database != nil {
		return c.database.Close()
	}

	return nil
}

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

func (c *Container) GetWebhookUseCases() input.WebhookUseCases {
	return c.webhookUseCases
}

func (c *Container) GetWebhookSender() output.WebhookSender {
	return c.webhookSender
}

func (c *Container) initWebhookSender() {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	c.webhookSender = webhook.NewHTTPWebhookSender(httpClient, c.logger)
}

func (c *Container) initWebhookUseCases() input.WebhookUseCases {
	webhookRepo := repository.NewWebhookRepository(c.database.DB)
	return webhookUseCase.NewWebhookUseCases(webhookRepo, c.webhookService)
}
