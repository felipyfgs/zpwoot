package router

import (
	"encoding/json"
	"net/http"

	"zpwoot/internal/adapters/database/repository"
	"zpwoot/internal/adapters/http/handlers"
	httpMiddleware "zpwoot/internal/adapters/http/middleware"
	"zpwoot/internal/adapters/waclient"
	"zpwoot/internal/container"

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
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-API-Key", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(httpMiddleware.CORSMiddleware())
	r.Use(httpMiddleware.SecurityHeadersMiddleware())
	r.Use(httpMiddleware.JSONMiddleware())
}

type Handlers struct {
	Session *handlers.SessionHandler
	Message *handlers.MessageHandler
}

func initializeHandlers(container *container.Container) *Handlers {

	sessionHandler := createSessionHandler(container)
	messageHandler := createMessageHandler(container)

	return &Handlers{
		Session: sessionHandler,
		Message: messageHandler,
	}
}

func createSessionHandler(container *container.Container) *handlers.SessionHandler {

	sessionRepository := repository.NewSessionRepository(container.GetDatabase().DB)
	sessionRepo := repository.NewSessionRepo(sessionRepository)
	waContainer := waclient.NewWAStoreContainer(
		container.GetDatabase().DB,
		container.GetLogger(),
		container.GetConfig().Database.URL,
	)
	waClient := waclient.NewWAClient(waContainer, container.GetLogger(), sessionRepo)

	sessionManager := waclient.NewSessionManagerAdapter(waClient)

	return handlers.NewSessionHandler(
		container.GetSessionUseCases(),
		sessionManager,
		container.GetLogger(),
	)
}

func createMessageHandler(container *container.Container) *handlers.MessageHandler {

	sessionRepository := repository.NewSessionRepository(container.GetDatabase().DB)
	sessionRepo := repository.NewSessionRepo(sessionRepository)
	waContainer := waclient.NewWAStoreContainer(
		container.GetDatabase().DB,
		container.GetLogger(),
		container.GetConfig().Database.URL,
	)
	waClient := waclient.NewWAClient(waContainer, container.GetLogger(), sessionRepo)

	messageSender := waclient.NewMessageSender(waClient)
	messageService := waclient.NewMessageServiceWrapper(messageSender)

	return handlers.NewMessageHandler(
		messageService,
		container.GetLogger(),
	)
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
