package router

import (
	"github.com/go-chi/chi/v5"

	"zpwoot/internal/api/handlers"
	"zpwoot/platform/logger"
)

func setupWebhookRoutes(r chi.Router, sessionService *usecases.SessionService, appLogger *logger.Logger) {
	webhookHandler := handler.NewWebhookHandler(sessionService, appLogger)

	r.Route("/{sessionId}/webhook", func(r chi.Router) {

		r.Post("/set", webhookHandler.SetConfig)
		r.Get("/find", webhookHandler.FindConfig)

		r.Post("/test", webhookHandler.TestWebhook)
	})
}
