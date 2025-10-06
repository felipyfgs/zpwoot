package handlers

import (
	"zpwoot/internal/adapters/database"
	"zpwoot/internal/adapters/database/repository"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/adapters/waclient"
	"zpwoot/internal/config"
	"zpwoot/internal/core/ports/input"
)

type Handlers struct {
	Session *SessionHandler
	Message *MessageHandler
	Health  *HealthHandler
}

func NewHandlers(
	db *database.Database,
	logger *logger.Logger,
	cfg *config.Config,
	sessionUseCases input.SessionUseCases,
	messageUseCases input.MessageUseCases,
) *Handlers {
	return &Handlers{
		Session: createSessionHandler(db, logger, cfg, sessionUseCases),
		Message: createMessageHandler(db, logger, cfg),
		Health:  NewHealthHandler(db),
	}
}

// createSessionHandler creates a session handler with all dependencies
func createSessionHandler(
	db *database.Database,
	logger *logger.Logger,
	cfg *config.Config,
	sessionUseCases input.SessionUseCases,
) *SessionHandler {
	sessionRepository := repository.NewSessionRepository(db.DB)
	sessionRepo := repository.NewSessionRepo(sessionRepository)
	waContainer := waclient.NewWAStoreContainer(
		db.DB,
		logger,
		cfg.Database.URL,
	)
	waClient := waclient.NewWAClient(waContainer, logger, sessionRepo)
	sessionManager := waclient.NewSessionManagerAdapter(waClient)

	return NewSessionHandler(
		sessionUseCases,
		sessionManager,
		logger,
	)
}

// createMessageHandler creates a message handler with all dependencies
func createMessageHandler(
	db *database.Database,
	logger *logger.Logger,
	cfg *config.Config,
) *MessageHandler {
	sessionRepository := repository.NewSessionRepository(db.DB)
	sessionRepo := repository.NewSessionRepo(sessionRepository)
	waContainer := waclient.NewWAStoreContainer(
		db.DB,
		logger,
		cfg.Database.URL,
	)
	waClient := waclient.NewWAClient(waContainer, logger, sessionRepo)

	messageSender := waclient.NewMessageSender(waClient)
	messageService := waclient.NewMessageServiceWrapper(messageSender)

	return NewMessageHandler(
		messageService,
		logger,
	)
}
