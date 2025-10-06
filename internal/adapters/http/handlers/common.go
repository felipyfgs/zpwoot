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
	Session    *SessionHandler
	Message    *MessageHandler
	Group      *GroupHandler
	Contact    *ContactHandler
	Community  *CommunityHandler
	Newsletter *NewsletterHandler
	Health     *HealthHandler
	Webhook    *WebhookHandler
}

func NewHandlers(
	db *database.Database,
	logger *logger.Logger,
	cfg *config.Config,
	sessionUseCases input.SessionUseCases,
	messageUseCases input.MessageUseCases,
	webhookUseCases input.WebhookUseCases,
	waClient output.WhatsAppClient,
) *Handlers {
	return &Handlers{
		Session:    createSessionHandler(logger, sessionUseCases, waClient),
		Message:    createMessageHandler(logger, waClient),
		Group:      createGroupHandler(logger, waClient),
		Contact:    createContactHandler(logger, waClient),
		Community:  createCommunityHandler(logger, waClient),
		Newsletter: createNewsletterHandler(logger, waClient),
		Health:     NewHealthHandler(db),
		Webhook:    NewWebhookHandler(webhookUseCases, logger),
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
	sessionManager := waclient.NewManager(waClientAdapter.GetWAClient())

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

	messageSender := waclient.NewSender(waClientAdapter.GetWAClient())
	messageService := waclient.NewMessageService(messageSender)

	return NewMessageHandler(
		messageService,
		logger,
	)
}

func createGroupHandler(
	logger *logger.Logger,
	waClient output.WhatsAppClient,
) *GroupHandler {
	waClientAdapter, ok := waClient.(*waclient.WAClientAdapter)
	if !ok {
		panic("waClient is not a WAClientAdapter")
	}

	groupService := waclient.NewGroupService(waClientAdapter.GetWAClient())

	return NewGroupHandler(
		groupService,
		logger,
	)
}

func createContactHandler(
	logger *logger.Logger,
	waClient output.WhatsAppClient,
) *ContactHandler {
	waClientAdapter, ok := waClient.(*waclient.WAClientAdapter)
	if !ok {
		panic("waClient is not a WAClientAdapter")
	}

	contactService := waclient.NewContactService(waClientAdapter.GetWAClient())

	return NewContactHandler(
		contactService,
		logger,
	)
}

func createCommunityHandler(
	logger *logger.Logger,
	waClient output.WhatsAppClient,
) *CommunityHandler {
	waClientAdapter, ok := waClient.(*waclient.WAClientAdapter)
	if !ok {
		panic("waClient is not a WAClientAdapter")
	}

	communityService := waclient.NewCommunityService(waClientAdapter.GetWAClient())

	return NewCommunityHandler(
		communityService,
		logger,
	)
}

func createNewsletterHandler(
	logger *logger.Logger,
	waClient output.WhatsAppClient,
) *NewsletterHandler {

	waClientAdapter, ok := waClient.(*waclient.WAClientAdapter)
	if !ok {
		panic("waClient is not a WAClientAdapter")
	}

	newsletterService := waclient.NewNewsletterService(waClientAdapter.GetWAClient())

	return NewNewsletterHandler(
		newsletterService,
		logger,
	)
}
