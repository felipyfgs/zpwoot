package router

import (
	"encoding/json"
	"net/http"

	"zpwoot/internal/container"
	"zpwoot/internal/adapters/database/repository"
	"zpwoot/internal/adapters/http/handlers"
	httpMiddleware "zpwoot/internal/adapters/http/middleware"
	"zpwoot/internal/adapters/waclient"
	sessionUseCase "zpwoot/internal/core/application/usecase/session"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "zpwoot/docs/swagger"
)


func NewRouter(container *container.Container) http.Handler {
	r := chi.NewRouter()


	setupMiddleware(r)


	h := initializeHandlers(container)


	setupPublicRoutes(r, container)


	setupAPIRoutes(r, container, h)

	return r
}


func setupMiddleware(r *chi.Mux) {

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)


	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))


	r.Use(httpMiddleware.CORSMiddleware())
	r.Use(httpMiddleware.JSONMiddleware())
}


type Handlers struct {
	Session *handlers.SessionHandler
	Message *handlers.MessageHandler
}


func initializeHandlers(container *container.Container) *Handlers {

	sessionRepo := repository.NewSessionRepository(container.GetDatabase().DB)


	sessionRepoAdapter := repository.NewSessionRepositoryAdapter(sessionRepo)


	waContainer := waclient.NewWAStoreContainer(container.GetDatabase().DB, container.GetLogger())


	waClient := waclient.NewWAClient(waContainer, container.GetLogger(), sessionRepoAdapter)


	whatsappAdapter := waclient.NewWhatsAppAdapter(waClient)


	sessionUseCases := sessionUseCase.NewUseCases(
		container.GetSessionService(),
		whatsappAdapter,
		container.GetNotificationService(),
	)


	messageSender := waclient.NewMessageSender(waClient)


	return &Handlers{
		Session: handlers.NewSessionHandler(sessionUseCases, waClient, container.GetLogger()),
		Message: handlers.NewMessageHandler(messageSender, container.GetLogger()),
	}
}


func setupPublicRoutes(r *chi.Mux, container *container.Container) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"message": "zpwoot WhatsApp API",
			"version": "1.0.0",
			"docs":    "/swagger/index.html",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})


	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		statusCode := http.StatusOK


		if container.GetDatabase() != nil {
			if err := container.GetDatabase().Health(); err != nil {
				status = "unhealthy"
				statusCode = http.StatusServiceUnavailable
			}
		}

		response := map[string]string{
			"status":  status,
			"service": "zpwoot",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	})


	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}


func setupAPIRoutes(r *chi.Mux, container *container.Container, h *Handlers) {

	r.Group(func(r chi.Router) {

		r.Use(httpMiddleware.AuthMiddleware(container.GetConfig()))


		setupSessionRoutes(r, h)
	})
}


func setupSessionRoutes(r chi.Router, h *Handlers) {
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
