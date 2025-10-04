package container

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"zpwoot/internal/core/group"
	"zpwoot/internal/core/messaging"
	"zpwoot/internal/core/session"

	"zpwoot/internal/usecases"
	"zpwoot/internal/usecases/shared/validation"

	"zpwoot/internal/adapters/repository"
	"zpwoot/internal/adapters/http"
	"zpwoot/internal/adapters/waclient"

	"zpwoot/platform/config"
	"zpwoot/platform/database"
	"zpwoot/platform/logger"
)

type Container struct {
	config   *config.Config
	logger   *logger.Logger
	database *database.Database

	sessionCore   *session.Service
	messagingCore *messaging.Service

	sessionService   *usecases.SessionService
	messagingService *usecases.MessageService
	groupService     *usecases.GroupService

	sessionRepo     session.Repository
	messageRepo     messaging.Repository
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

	c.messagingCore = messaging.NewService(
		c.messageRepo,
	)

	validator := validation.New()

	c.sessionService = usecases.NewSessionService(
		c.sessionCore,
		c.logger,
		validator,
	)

	c.messagingService = usecases.NewMessageService(
		c.messagingCore,
		c.sessionCore,
		c.logger,
		validator,
	)

	groupCore := group.NewService(nil)

	c.groupService = usecases.NewGroupService(
		groupCore,
		nil,
		nil,
		c.logger,
		validator,
	)

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

func (c *Container) GetSessionService() *usecases.SessionService {
	return c.sessionService
}

func (c *Container) GetMessageService() *usecases.MessageService {
	return c.messagingService
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

type sessionServiceAdapter struct {
	service *usecases.SessionService
}

func (a *sessionServiceAdapter) UpdateDeviceJID(ctx context.Context, id uuid.UUID, deviceJID string) error {

	return nil
}

func (a *sessionServiceAdapter) UpdateQRCode(ctx context.Context, id uuid.UUID, qrCode string, expiresAt time.Time) error {

	return nil
}

func (a *sessionServiceAdapter) ClearQRCode(ctx context.Context, id uuid.UUID) error {

	return nil
}

func (c *Container) GetGroupService() *usecases.GroupService {
	return c.groupService
}

// GetSessionResolver removed - using sessionId directly
