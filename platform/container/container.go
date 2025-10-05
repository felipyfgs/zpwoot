package container

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"zpwoot/internal/domain/group"
	"zpwoot/internal/domain/message"
	"zpwoot/internal/domain/session"

	"zpwoot/internal/api"
	"zpwoot/internal/infrastructure/repository"
	"zpwoot/internal/infrastructure/whatsapp"

	"zpwoot/platform/config"
	"zpwoot/platform/database"
	"zpwoot/platform/logger"
)

type Container struct {
	config   *config.Config
	logger   *logger.Logger
	database *database.Database

	sessionCore   *session.Service
	messagingCore *message.Service

	// Services are now directly the domain services
	// sessionService   *session.Service (same as sessionCore)
	// messagingService *message.Service (same as messagingCore)
	// groupService     *group.Service

	sessionRepo     session.Repository
	messageRepo     message.Repository
	whatsappGateway session.WhatsAppGateway
}

type Config struct {
	AppConfig *config.Config
	Logger    *logger.Logger
	Database  *database.Database
}

func New(cfg *Config) (*Container, error) {
	container := &Container{
		config:   cfg.AppConfig,
		logger:   cfg.Logger,
		database: cfg.Database,
	}

	if err := container.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize container: %w", err)
	}

	return container, nil
}

func (c *Container) createWhatsAppContainer() (*sqlstore.Container, error) {
	c.logger.Debug("Creating WhatsApp sqlstore container...")

	ctx := context.Background()

	waLogger := waclient.NewWhatsmeowLogger(c.logger)

	container, err := sqlstore.New(ctx, "postgres", c.config.Database.URL, waLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlstore container: %w", err)
	}

	err = container.Upgrade(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade sqlstore: %w", err)
	}

	return container, nil
}

func (c *Container) initialize() error {
	c.logger.Debug("Initializing container...")

	c.sessionRepo = repository.NewSessionRepository(c.database.DB)
	c.messageRepo = repository.NewMessageRepository(c.database.DB, c.logger)

	waContainer, err := c.createWhatsAppContainer()
	if err != nil {
		return fmt.Errorf("failed to create WhatsApp container: %w", err)
	}

	gateway := waclient.NewGateway(waContainer, c.logger)
	gateway.SetDatabase(c.database.DB)
	c.whatsappGateway = gateway

	qrGenerator := waclient.NewQRGenerator(c.logger)

	c.sessionCore = session.NewService(
		c.sessionRepo,
		c.whatsappGateway,
		qrGenerator,
	)

	c.messagingCore = message.NewService(
		c.messageRepo,
	)

	// Services are now directly the domain services
	// No need for separate use case layer in this simplified architecture

	sessionEventHandler := session.NewSessionEventHandler(c.sessionCore)
	c.whatsappGateway.SetEventHandler(sessionEventHandler)

	c.logger.Debug("Container initialized successfully")
	return nil
}

func (c *Container) Start(ctx context.Context) error {
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

func (c *Container) GetSessionService() *session.Service {
	return c.sessionCore
}

func (c *Container) GetMessageService() *message.Service {
	return c.messagingCore
}

func (c *Container) GetSessionCore() *session.Service {
	return c.sessionCore
}

func (c *Container) GetWhatsAppGateway() session.WhatsAppGateway {
	return c.whatsappGateway
}

func (c *Container) Stop(ctx context.Context) error {

	if stopper, ok := c.whatsappGateway.(interface{ Stop(context.Context) error }); ok {
		stopper.Stop(ctx)
	}

	c.database.Close()

	return nil
}

func (c *Container) Server() *server.Server {
	return server.New(&server.Config{
		Config:         c.config,
		Logger:         c.logger,
		SessionService: c.sessionService,
		MessageService: c.messagingService,
		GroupService:   c.groupService,
	})
}

func (c *Container) Handler() http.Handler {
	return c.Server().Handler()
}

// Group service is not implemented yet in this simplified version
func (c *Container) GetGroupService() *group.Service {
	// Return a basic group service for now
	return group.NewService(nil, nil, nil)
}
