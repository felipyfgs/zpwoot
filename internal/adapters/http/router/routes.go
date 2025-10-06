package router

import (
	"net/http"

	"zpwoot/internal/adapters/http/handlers"
	"zpwoot/internal/adapters/http/middleware"
	"zpwoot/internal/container"

	_ "zpwoot/docs/swagger"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(c *container.Container) http.Handler {
	r := chi.NewRouter()

	middleware.SetupMiddleware(r)

	h := handlers.NewHandlers(
		c.GetDatabase(),
		c.GetLogger(),
		c.GetConfig(),
		c.GetSessionUseCases(),
		c.GetMessageUseCases(),
		c.GetWhatsAppClient(),
	)

	setupPublicRoutes(r, h)
	setupAPIRoutes(r, c, h)

	return r
}

func setupPublicRoutes(r *chi.Mux, h *handlers.Handlers) {
	r.Get("/", h.Health.Info)
	r.Get("/health", h.Health.Health)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}

func setupAPIRoutes(r *chi.Mux, c *container.Container, h *handlers.Handlers) {
	r.Group(func(r chi.Router) {
		middleware.SetupAuthMiddleware(r, c.GetConfig())

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
	})

	setupMessageRoutes(r, h)
}

func setupMessageRoutes(r chi.Router, h *handlers.Handlers) {
	r.Route("/sessions/{sessionId}/send/message", func(r chi.Router) {
		r.Post("/text", h.Message.SendTextMessage)
		r.Post("/image", h.Message.SendImage)
		r.Post("/audio", h.Message.SendAudio)
		r.Post("/video", h.Message.SendVideo)
		r.Post("/document", h.Message.SendDocument)
		r.Post("/sticker", h.Message.SendSticker)
		r.Post("/location", h.Message.SendLocation)
		r.Post("/contact", h.Message.SendContact)
		r.Post("/contacts", h.Message.SendContactsArray)
		r.Post("/reaction", h.Message.SendReaction)
		r.Post("/template", h.Message.SendTemplate)
		r.Post("/buttons", h.Message.SendButtons)
		r.Post("/list", h.Message.SendList)
		r.Post("/poll", h.Message.SendPoll)
		r.Post("/viewonce", h.Message.SendViewOnce)
	})
}
