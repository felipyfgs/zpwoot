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

	r.Get("/{sessionName}/info", sessionHandler.GetSessionInfo)
	r.Delete("/{sessionName}/delete", sessionHandler.DeleteSession)

	r.Post("/{sessionName}/connect", sessionHandler.ConnectSession)
	r.Post("/{sessionName}/logout", sessionHandler.LogoutSession)
	r.Get("/{sessionName}/qr", sessionHandler.GetQRCode)
	r.Post("/{sessionName}/pair", sessionHandler.PairPhone)

	r.Post("/{sessionName}/proxy/set", sessionHandler.SetProxy)
	r.Get("/{sessionName}/proxy/find", sessionHandler.GetProxy)

	r.Get("/{sessionName}/stats", sessionHandler.GetSessionStats)
}
