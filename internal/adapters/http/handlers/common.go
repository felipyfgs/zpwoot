package handlers

import (
	"zpwoot/internal/adapters/database"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/adapters/waclient"
	"zpwoot/internal/config"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"
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
	waClient output.WhatsAppClient,
) *Handlers {
	return &Handlers{
		Session: createSessionHandler(logger, sessionUseCases, waClient),
		Message: createMessageHandler(logger, waClient),
		Health:  NewHealthHandler(db),
	}
}

func createSessionHandler(
	logger *logger.Logger,
	sessionUseCases input.SessionUseCases,
	waClient output.WhatsAppClient,
) *SessionHandler {

	waClientAdapter, ok := waClient.(*waclient.WAClientAdapter)
	if !ok {
		panic("waClient is not a WAClientAdapter")
	}
	sessionManager := waclient.NewSessionManagerAdapter(waClientAdapter.GetWAClient())

	return NewSessionHandler(
		sessionUseCases,
		sessionManager,
		logger,
	)
}

func createMessageHandler(
	logger *logger.Logger,
	waClient output.WhatsAppClient,
) *MessageHandler {

	waClientAdapter, ok := waClient.(*waclient.WAClientAdapter)
	if !ok {
		panic("waClient is not a WAClientAdapter")
	}

	messageSender := waclient.NewMessageSender(waClientAdapter.GetWAClient())
	messageService := waclient.NewMessageServiceWrapper(messageSender)

	return NewMessageHandler(
		messageService,
		logger,
	)
}
