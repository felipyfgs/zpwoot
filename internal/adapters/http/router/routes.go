package router

import (
	"net/http"

	"zpwoot/internal/adapters/http/handlers"
	"zpwoot/internal/adapters/http/middleware"
	"zpwoot/internal/container"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "zpwoot/docs/swagger"
)

func NewRouter(container *container.Container) http.Handler {
	r := chi.NewRouter()

	middleware.SetupMiddleware(r)

	// Initialize handlers
	h := handlers.NewHandlers(
		container.GetDatabase(),
		container.GetLogger(),
		container.GetConfig(),
		container.GetSessionUseCases(),
		container.GetMessageUseCases(),
	)

	setupPublicRoutes(r, h)
	setupAPIRoutes(r, container, h)

	return r
}

func setupPublicRoutes(r *chi.Mux, h *handlers.Handlers) {

	r.Get("/", h.Health.Info)

	r.Get("/health", h.Health.Health)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}

func setupAPIRoutes(r *chi.Mux, container *container.Container, h *handlers.Handlers) {
	r.Group(func(r chi.Router) {

		middleware.SetupAuthMiddleware(r, container.GetConfig())

		setupSessionRoutes(r, h)
	})
}

func setupSessionRoutes(r chi.Router, h *handlers.Handlers) {
	r.Route("/sessions", func(r chi.Router) {

		r.Post("/create", h.Session.CreateSession)
		r.Get("/list", h.Session.ListSessions)
		r.Get("/{sessionId}/info", h.Session.GetSession)
		r.Delete("/{sessionId}/delete", h.Session.DeleteSession)

		r.Post("/{sessionId}/connect", h.Session.ConnectSession)
		r.Post("/{sessionId}/disconnect", h.Session.DisconnectSession)
		r.Post("/{sessionId}/logout", h.Session.LogoutSession)
		r.Get("/{sessionId}/qr", h.Session.GetQRCode)

		r.Post("/{sessionId}/messages", h.Message.SendMessage)
		r.Get("/{sessionId}/chats", h.Message.GetChats)
		r.Get("/{sessionId}/contacts", h.Message.GetContacts)
		r.Get("/{sessionId}/chat-info", h.Message.GetChatInfo)
	})
}
