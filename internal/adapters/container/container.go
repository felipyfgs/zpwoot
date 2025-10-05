package container

import (
	"context"
	"fmt"
	"net/http"

	"zpwoot/internal/adapters/config"
	"zpwoot/internal/adapters/database"
	"zpwoot/internal/adapters/database/repository"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/application/interfaces"
	"zpwoot/internal/application/usecase/message"
	"zpwoot/internal/application/usecase/session"
	domainSession "zpwoot/internal/domain/session"
)

// Container holds all application dependencies
type Container struct {
	// Infrastructure
	config   *config.Config
	logger   *logger.Logger
	database *database.Database
	migrator *database.Migrator

	// Domain Services
	sessionService *domainSession.Service

	// External Adapters
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService

	// Use Cases
	sessionUseCases *SessionUseCases
	messageUseCases *MessageUseCases
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

	c.logger.Info().Msg("Initializing zpwoot container")

	// Initialize database
	c.logger.Info().Msg("Connecting to database")
	db, err := database.New(c.config, c.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	c.database = db
	c.logger.Info().Msg("Database connection established")

	// Initialize migrator
	c.logger.Info().Msg("Initializing database migrator")
	c.migrator = database.NewMigrator(db, c.logger)

	// Run migrations automatically
	c.logger.Info().Msg("Running database migrations")
	if err := c.migrator.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	c.logger.Info().Msg("Database migrations completed")

	// Initialize domain services
	c.logger.Info().Msg("Initializing domain services")
	sessionRepo := repository.NewSessionRepository(c.database.DB)
	c.sessionService = domainSession.NewService(sessionRepo)

	// Initialize use cases
	c.logger.Info().Msg("Initializing use cases")
	c.initializeUseCases()

	c.logger.Info().Msg("Container initialization completed successfully")
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

// initializeUseCases initializes all use cases
func (c *Container) initializeUseCases() {
	// Initialize session use cases
	c.sessionUseCases = &SessionUseCases{
		Create:     session.NewCreateUseCase(c.sessionService, c.whatsappClient, c.notificationSvc),
		Get:        session.NewGetUseCase(c.sessionService, c.whatsappClient),
		List:       session.NewListUseCase(c.sessionService, c.whatsappClient),
		Connect:    session.NewConnectUseCase(c.sessionService, c.whatsappClient, c.notificationSvc),
		Disconnect: session.NewDisconnectUseCase(c.sessionService, c.whatsappClient, c.notificationSvc),
		Delete:     session.NewDeleteUseCase(c.sessionService, c.whatsappClient, c.notificationSvc),
		QR:         session.NewQRUseCase(c.sessionService, c.whatsappClient, c.notificationSvc),
	}

	// Initialize message use cases
	c.messageUseCases = &MessageUseCases{
		Send:    message.NewSendUseCase(c.sessionService, c.whatsappClient, c.notificationSvc),
		Receive: message.NewReceiveUseCase(c.sessionService, c.notificationSvc),
	}
}

// GetSessionService returns the session service
func (c *Container) GetSessionService() *domainSession.Service {
	return c.sessionService
}

// GetSessionUseCases returns all session-related use cases
func (c *Container) GetSessionUseCases() *SessionUseCases {
	return c.sessionUseCases
}

// GetMessageUseCases returns all message-related use cases
func (c *Container) GetMessageUseCases() *MessageUseCases {
	return c.messageUseCases
}

// SetWhatsAppClient sets the WhatsApp client
func (c *Container) SetWhatsAppClient(client interfaces.WhatsAppClient) {
	c.whatsappClient = client
}

// SetNotificationService sets the notification service
func (c *Container) SetNotificationService(svc interfaces.NotificationService) {
	c.notificationSvc = svc
}

// SessionUseCases groups session use cases
type SessionUseCases struct {
	Create     *session.CreateUseCase
	Get        *session.GetUseCase
	List       *session.ListUseCase
	Connect    *session.ConnectUseCase
	Disconnect *session.DisconnectUseCase
	Delete     *session.DeleteUseCase
	QR         *session.QRUseCase
}

// MessageUseCases groups message use cases
type MessageUseCases struct {
	Send    *message.SendUseCase
	Receive *message.ReceiveUseCase
}
