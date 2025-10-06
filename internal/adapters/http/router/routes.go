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
		r.Post("/sessions", h.Session.Create)
		r.Get("/sessions", h.Session.List)
		r.Get("/sessions/{sessionId}", h.Session.Get)
		r.Delete("/sessions/{sessionId}", h.Session.Delete)
		r.Post("/sessions/{sessionId}/connect", h.Session.Connect)
		r.Post("/sessions/{sessionId}/disconnect", h.Session.Disconnect)
		r.Post("/sessions/{sessionId}/logout", h.Session.Logout)
		r.Get("/sessions/{sessionId}/qr", h.Session.QRCode)
		r.Post("/sessions/{sessionId}/pair", h.Session.PairPhone)

		// --- MENSAGENS (prefixo send) ---
		r.Post("/sessions/{sessionId}/messages/send/text", h.Message.SendText)
		r.Post("/sessions/{sessionId}/messages/send/image", h.Message.SendImage)
		r.Post("/sessions/{sessionId}/messages/send/audio", h.Message.SendAudio)
		r.Post("/sessions/{sessionId}/messages/send/video", h.Message.SendVideo)
		r.Post("/sessions/{sessionId}/messages/send/document", h.Message.SendDocument)
		r.Post("/sessions/{sessionId}/messages/send/sticker", h.Message.SendSticker)
		r.Post("/sessions/{sessionId}/messages/send/location", h.Message.SendLocation)
		r.Post("/sessions/{sessionId}/messages/send/contact", h.Message.SendContact)
		r.Post("/sessions/{sessionId}/messages/send/contacts", h.Message.SendContactsArray)
		r.Post("/sessions/{sessionId}/messages/send/reaction", h.Message.SendReaction)
		r.Post("/sessions/{sessionId}/messages/send/poll", h.Message.SendPoll)
		r.Post("/sessions/{sessionId}/messages/send/buttons", h.Message.SendButtons)
		r.Post("/sessions/{sessionId}/messages/send/list", h.Message.SendList)
		r.Post("/sessions/{sessionId}/messages/send/template", h.Message.SendTemplate)
		r.Post("/sessions/{sessionId}/messages/delete", h.Message.DeleteMessage)
		r.Post("/sessions/{sessionId}/messages/edit", h.Message.EditMessage)
		r.Post("/sessions/{sessionId}/messages/markread", h.Message.MarkRead)
		r.Post("/sessions/{sessionId}/messages/historysync", h.Message.RequestHistorySync)

		// --- PRESENÇA ---
		r.Post("/sessions/{sessionId}/presence/send", h.Contact.SendPresence)
		r.Post("/sessions/{sessionId}/presence/chat", h.Contact.ChatPresence)

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

		// --- COMUNIDADES ---
		r.Get("/sessions/{sessionId}/communities", h.Community.ListCommunities)
		r.Get("/sessions/{sessionId}/communities/info", h.Community.GetCommunityInfo)
		r.Post("/sessions/{sessionId}/communities", h.Community.CreateCommunity)
		r.Get("/sessions/{sessionId}/communities/{communityJid}/groups", h.Community.GetSubGroups)
		r.Get("/sessions/{sessionId}/communities/{communityJid}/participants", h.Community.GetParticipants)
		r.Post("/sessions/{sessionId}/communities/{communityJid}/link", h.Community.LinkGroup)
		r.Post("/sessions/{sessionId}/communities/{communityJid}/unlink", h.Community.UnlinkGroup)

		// --- NEWSLETTERS ---
		r.Get("/sessions/{sessionId}/newsletters", h.Newsletter.ListNewsletters)
		r.Get("/sessions/{sessionId}/newsletters/info", h.Newsletter.GetNewsletterInfo)
		r.Post("/sessions/{sessionId}/newsletters/info-invite", h.Newsletter.GetNewsletterInfoWithInvite)
		r.Post("/sessions/{sessionId}/newsletters", h.Newsletter.CreateNewsletter)
		r.Post("/sessions/{sessionId}/newsletters/follow", h.Newsletter.FollowNewsletter)
		r.Post("/sessions/{sessionId}/newsletters/{newsletterJid}/unfollow", h.Newsletter.UnfollowNewsletter)
		r.Get("/sessions/{sessionId}/newsletters/{newsletterJid}/messages", h.Newsletter.GetMessages)
		r.Post("/sessions/{sessionId}/newsletters/{newsletterJid}/mark-viewed", h.Newsletter.MarkViewed)
		r.Post("/sessions/{sessionId}/newsletters/{newsletterJid}/react", h.Newsletter.SendReaction)
		r.Post("/sessions/{sessionId}/newsletters/{newsletterJid}/mute", h.Newsletter.ToggleMute)
		r.Post("/sessions/{sessionId}/newsletters/{newsletterJid}/send", h.Newsletter.SendMessage)

		// --- WEBHOOKS ---
		r.Post("/sessions/{sessionId}/webhooks", h.Webhook.SetWebhook)
		r.Get("/sessions/{sessionId}/webhooks", h.Webhook.GetWebhook)
		r.Delete("/sessions/{sessionId}/webhooks", h.Webhook.DeleteWebhook)
		r.Get("/webhooks/events", h.Webhook.ListEvents)
	})

	return r
}
