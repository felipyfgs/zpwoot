package router

import (
	"github.com/go-chi/chi/v5"

	"zpwoot/internal/adapters/server/handler"
	"zpwoot/internal/core/session"
	"zpwoot/internal/services"
	"zpwoot/platform/logger"
)

func setupWebhookRoutes(r chi.Router, sessionService *services.SessionService, sessionResolver session.SessionResolver, appLogger *logger.Logger) {
	webhookHandler := handler.NewWebhookHandler(sessionService, sessionResolver, appLogger)

	r.Route("/{sessionName}/webhook", func(r chi.Router) {

		r.Post("/set", webhookHandler.SetConfig)
		r.Get("/find", webhookHandler.FindConfig)

		r.Post("/test", webhookHandler.TestWebhook)
	})
}
