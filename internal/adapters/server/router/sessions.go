package router

import (
	"github.com/go-chi/chi/v5"

	"zpwoot/internal/adapters/server/handler"
	"zpwoot/internal/core/session"
	"zpwoot/internal/services"
	"zpwoot/platform/logger"
)

func setupSessionRoutes(r chi.Router, sessionService *services.SessionService, sessionResolver session.SessionResolver, appLogger *logger.Logger) {
	sessionHandler := handler.NewSessionHandler(sessionService, sessionResolver, appLogger)

	r.Post("/create", sessionHandler.CreateSession)
	r.Get("/list", sessionHandler.ListSessions)

	r.Get("/{sessionId}/info", sessionHandler.GetSessionInfo)
	r.Delete("/{sessionId}/delete", sessionHandler.DeleteSession)

	r.Post("/{sessionId}/connect", sessionHandler.ConnectSession)
	r.Post("/{sessionId}/logout", sessionHandler.LogoutSession)
	r.Get("/{sessionId}/qr", sessionHandler.GetQRCode)
	r.Post("/{sessionId}/pair", sessionHandler.PairPhone)

	r.Post("/{sessionId}/proxy/set", sessionHandler.SetProxy)
	r.Get("/{sessionId}/proxy/find", sessionHandler.GetProxy)

	r.Get("/{sessionId}/stats", sessionHandler.GetSessionStats)
}
