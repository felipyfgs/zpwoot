package router

import (
	"github.com/go-chi/chi/v5"

	"zpwoot/internal/api/handlers"
	"zpwoot/platform/logger"
)

func setupChatwootRoutes(r chi.Router, messageService *usecases.MessageService, sessionService *usecases.SessionService, appLogger *logger.Logger) {
	chatwootHandler := handler.NewChatwootHandler(messageService, sessionService, appLogger)

	r.Route("/{sessionId}/chatwoot", func(r chi.Router) {

		r.Post("/set", chatwootHandler.CreateConfig)
		r.Get("/find", chatwootHandler.FindConfig)

		r.Get("/", chatwootHandler.FindConfig)
		r.Put("/", chatwootHandler.UpdateConfig)
		r.Delete("/", chatwootHandler.DeleteConfig)
		r.Post("/test", chatwootHandler.TestConnection)
		r.Post("/auto-create-inbox", chatwootHandler.AutoCreateInbox)
		r.Get("/stats", chatwootHandler.GetStats)
	})
}
