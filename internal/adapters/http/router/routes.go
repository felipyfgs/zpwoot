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
		setupMessageRoutes(r, h)
		setupGroupRoutes(r, h)
	})
}

func setupSessionRoutes(r chi.Router, h *handlers.Handlers) {
	r.Route("/sessions", func(r chi.Router) {
		// CRUD de sessões
		r.Post("/create", h.Session.Create)
		r.Get("/list", h.Session.List)
		r.Get("/{sessionId}/info", h.Session.Get)
		r.Delete("/{sessionId}/delete", h.Session.Delete)

		// Controle de conexão
		r.Post("/{sessionId}/connect", h.Session.Connect)
		r.Post("/{sessionId}/disconnect", h.Session.Disconnect)
		r.Post("/{sessionId}/logout", h.Session.Logout)
		r.Get("/{sessionId}/qr", h.Session.QRCode)
		r.Post("/{sessionId}/pair", h.Session.PairPhone)
	})
}

func setupMessageRoutes(r chi.Router, h *handlers.Handlers) {
	// Envio de mensagens
	r.Route("/sessions/{sessionId}/send/message", func(r chi.Router) {
		r.Post("/text", h.Message.SendText)
		r.Post("/image", h.Message.SendImage)
		r.Post("/audio", h.Message.SendAudio)
		r.Post("/video", h.Message.SendVideo)
		r.Post("/document", h.Message.SendDocument)
		r.Post("/sticker", h.Message.SendSticker)
		r.Post("/location", h.Message.SendLocation)
		r.Post("/contact", h.Message.SendContact)
		r.Post("/contacts", h.Message.SendContactsArray)
		r.Post("/reaction", h.Message.SendReaction)
		r.Post("/poll", h.Message.SendPoll)
		r.Post("/buttons", h.Message.SendButtons)
		r.Post("/list", h.Message.SendList)
		r.Post("/template", h.Message.SendTemplate)
	})

	// Operações de mensagens (não envio)
	r.Route("/sessions/{sessionId}/message", func(r chi.Router) {
		r.Post("/delete", h.Message.DeleteMessage)
		r.Post("/edit", h.Message.EditMessage)
		r.Post("/markread", h.Message.MarkRead)
		r.Post("/historysync", h.Message.RequestHistorySync)
	})
}

func setupGroupRoutes(r chi.Router, h *handlers.Handlers) {
	r.Route("/sessions/{sessionId}/groups", func(r chi.Router) {
		// Informações
		r.Get("/", h.Group.ListGroups)
		r.Get("/info", h.Group.GetGroupInfo)
		r.Post("/invite-info", h.Group.GetGroupInviteInfo)

		// Convites
		r.Get("/invite-link", h.Group.GetGroupInviteLink)
		r.Post("/join", h.Group.JoinGroup)

		// Gerenciamento básico
		r.Post("/create", h.Group.CreateGroup)
		r.Post("/leave", h.Group.LeaveGroup)
		r.Post("/participants", h.Group.UpdateGroupParticipants)

		// Configurações do grupo
		r.Post("/name", h.Group.SetGroupName)
		r.Post("/topic", h.Group.SetGroupTopic)

		// Configurações avançadas
		r.Route("/settings", func(r chi.Router) {
			r.Post("/locked", h.Group.SetGroupLocked)
			r.Post("/announce", h.Group.SetGroupAnnounce)
			r.Post("/disappearing", h.Group.SetDisappearingTimer)
		})

		// Mídia
		r.Post("/photo", h.Group.SetGroupPhoto)
		r.Delete("/photo", h.Group.RemoveGroupPhoto)
	})
}
