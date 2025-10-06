package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"zpwoot/internal/adapters/http/handlers"
	"zpwoot/internal/adapters/http/middleware"
	"zpwoot/internal/container"
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
		c.GetWebhookUseCases(),
		c.GetWhatsAppClient(),
	)

	// === PÚBLICAS ===
	r.Get("/", h.Health.Info)
	r.Get("/health", h.Health.Health)
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	// === PRIVADAS ===
	r.Group(func(r chi.Router) {
		middleware.SetupAuthMiddleware(r, c.GetConfig())

		// --- SESSÕES ---
		r.Post("/sessions/create", h.Session.Create)
		r.Get("/sessions/list", h.Session.List)
		r.Get("/sessions/{sessionId}/info", h.Session.Get)
		r.Delete("/sessions/{sessionId}/delete", h.Session.Delete)
		r.Post("/sessions/{sessionId}/connect", h.Session.Connect)
		r.Post("/sessions/{sessionId}/disconnect", h.Session.Disconnect)
		r.Post("/sessions/{sessionId}/logout", h.Session.Logout)
		r.Get("/sessions/{sessionId}/qr", h.Session.QRCode)
		r.Post("/sessions/{sessionId}/pair", h.Session.PairPhone)

		// --- MENSAGENS ---
		r.Post("/sessions/{sessionId}/send/message/text", h.Message.SendText)
		r.Post("/sessions/{sessionId}/send/message/image", h.Message.SendImage)
		r.Post("/sessions/{sessionId}/send/message/audio", h.Message.SendAudio)
		r.Post("/sessions/{sessionId}/send/message/video", h.Message.SendVideo)
		r.Post("/sessions/{sessionId}/send/message/document", h.Message.SendDocument)
		r.Post("/sessions/{sessionId}/send/message/sticker", h.Message.SendSticker)
		r.Post("/sessions/{sessionId}/send/message/location", h.Message.SendLocation)
		r.Post("/sessions/{sessionId}/send/message/contact", h.Message.SendContact)
		r.Post("/sessions/{sessionId}/send/message/contacts", h.Message.SendContactsArray)
		r.Post("/sessions/{sessionId}/send/message/reaction", h.Message.SendReaction)
		r.Post("/sessions/{sessionId}/send/message/poll", h.Message.SendPoll)
		r.Post("/sessions/{sessionId}/send/message/buttons", h.Message.SendButtons)
		r.Post("/sessions/{sessionId}/send/message/list", h.Message.SendList)
		r.Post("/sessions/{sessionId}/send/message/template", h.Message.SendTemplate)

		// --- PRESENÇA ---
		r.Post("/sessions/{sessionId}/send/presence", h.Contact.SendPresence)
		r.Post("/sessions/{sessionId}/presence", h.Contact.ChatPresence)

		// --- GERENCIAMENTO DE MENSAGENS ---
		r.Post("/sessions/{sessionId}/message/delete", h.Message.DeleteMessage)
		r.Post("/sessions/{sessionId}/message/edit", h.Message.EditMessage)
		r.Post("/sessions/{sessionId}/message/markread", h.Message.MarkRead)
		r.Post("/sessions/{sessionId}/message/historysync", h.Message.RequestHistorySync)

		// --- CONTATOS ---
		r.Get("/sessions/{sessionId}/contacts", h.Contact.GetContacts)
		r.Post("/sessions/{sessionId}/contacts/check", h.Contact.CheckUser)
		r.Post("/sessions/{sessionId}/contacts/user", h.Contact.GetUser)
		r.Post("/sessions/{sessionId}/contacts/avatar", h.Contact.GetAvatar)

		// --- GRUPOS ---
		r.Get("/sessions/{sessionId}/groups", h.Group.ListGroups)
		r.Get("/sessions/{sessionId}/groups/info", h.Group.GetGroupInfo)
		r.Post("/sessions/{sessionId}/groups/invite-info", h.Group.GetGroupInviteInfo)
		r.Get("/sessions/{sessionId}/groups/invite-link", h.Group.GetGroupInviteLink)
		r.Post("/sessions/{sessionId}/groups/join", h.Group.JoinGroup)
		r.Post("/sessions/{sessionId}/groups/create", h.Group.CreateGroup)
		r.Post("/sessions/{sessionId}/groups/leave", h.Group.LeaveGroup)
		r.Post("/sessions/{sessionId}/groups/participants", h.Group.UpdateGroupParticipants)
		r.Post("/sessions/{sessionId}/groups/name", h.Group.SetGroupName)
		r.Post("/sessions/{sessionId}/groups/topic", h.Group.SetGroupTopic)
		r.Post("/sessions/{sessionId}/groups/settings/locked", h.Group.SetGroupLocked)
		r.Post("/sessions/{sessionId}/groups/settings/announce", h.Group.SetGroupAnnounce)
		r.Post("/sessions/{sessionId}/groups/settings/disappearing", h.Group.SetDisappearingTimer)
		r.Post("/sessions/{sessionId}/groups/photo", h.Group.SetGroupPhoto)
		r.Delete("/sessions/{sessionId}/groups/photo", h.Group.RemoveGroupPhoto)

		// --- WEBHOOKS ---
		r.Post("/sessions/{sessionId}/webhook/create", h.Webhook.SetWebhook)
		r.Get("/sessions/{sessionId}/webhook/info", h.Webhook.GetWebhook)
		r.Delete("/sessions/{sessionId}/webhook/delete", h.Webhook.DeleteWebhook)
		r.Get("/webhook/events", h.Webhook.ListEvents)
	})

	return r
}
