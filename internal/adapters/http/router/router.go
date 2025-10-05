package router

import (
	"encoding/json"
	"net/http"

	"zpwoot/internal/adapters/container"
	"zpwoot/internal/adapters/http/handlers"
	httpMiddleware "zpwoot/internal/adapters/http/middleware"
	"zpwoot/internal/adapters/waclient"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "zpwoot/docs/swagger" // Import swagger docs
)

// NewRouter creates and configures the HTTP router with all routes and middleware
func NewRouter(container *container.Container) http.Handler {
	r := chi.NewRouter()

	// Setup global middleware (must be before routes)
	setupMiddleware(r)

	// Initialize handlers
	h := initializeHandlers(container)

	// Setup public routes (no auth)
	setupPublicRoutes(r, container)

	// Setup protected routes (with auth) - using a group to apply middleware
	setupAPIRoutes(r, container, h)

	return r
}

// setupMiddleware configures all middleware for the router
func setupMiddleware(r *chi.Mux) {
	// Basic middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Custom middleware
	r.Use(httpMiddleware.CORSMiddleware())
	r.Use(httpMiddleware.JSONMiddleware())
}

// Handlers holds all HTTP handlers
type Handlers struct {
	Session *handlers.SessionHandler
	Message *handlers.MessageHandler
}

// initializeHandlers initializes all handlers and their dependencies
func initializeHandlers(container *container.Container) *Handlers {
	// Create session manager
	sessionManager := waclient.NewDBSessionManager(container.GetDatabase().DB)

	// Create WhatsApp store container
	waContainer := waclient.NewWAStoreContainer(container.GetDatabase().DB, container.GetLogger())

	// Create WhatsApp client
	waClient := waclient.NewWAClient(waContainer, container.GetLogger(), sessionManager)

	// Create message sender
	messageSender := waclient.NewMessageSender(waClient)

	// Create and return handlers
	return &Handlers{
		Session: handlers.NewSessionHandler(waClient, container.GetLogger()),
		Message: handlers.NewMessageHandler(messageSender, container.GetLogger()),
	}
}

// setupPublicRoutes configures routes that don't require authentication
func setupPublicRoutes(r *chi.Mux, container *container.Container) {
	// Root endpoint - API information
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

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		statusCode := http.StatusOK

		// Check database health
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

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}

// setupAPIRoutes configures protected API routes that require authentication
func setupAPIRoutes(r *chi.Mux, container *container.Container, h *Handlers) {
	// Create a group for protected routes with auth middleware
	r.Group(func(r chi.Router) {
		// Apply authentication middleware to this group
		r.Use(httpMiddleware.AuthMiddleware(container.GetConfig()))

		// Session routes
		setupSessionRoutes(r, h)
	})
}

// setupSessionRoutes configures all session-related routes
func setupSessionRoutes(r chi.Router, h *Handlers) {
	r.Route("/sessions", func(r chi.Router) {
		// Session management
		r.Post("/create", h.Session.CreateSession)
		r.Get("/list", h.Session.ListSessions)
		r.Get("/{sessionId}/info", h.Session.GetSession)
		r.Delete("/{sessionId}/delete", h.Session.DeleteSession)

		// Session actions
		r.Post("/{sessionId}/connect", h.Session.ConnectSession)
		r.Post("/{sessionId}/disconnect", h.Session.DisconnectSession)
		r.Post("/{sessionId}/logout", h.Session.LogoutSession)
		r.Get("/{sessionId}/qr", h.Session.GetQRCode)
		r.Post("/{sessionId}/qr/refresh", h.Session.RefreshQRCode)

		// Message operations
		r.Post("/{sessionId}/messages", h.Message.SendMessage)
		r.Get("/{sessionId}/chats", h.Message.GetChats)
		r.Get("/{sessionId}/contacts", h.Message.GetContacts)
		r.Get("/{sessionId}/chat-info", h.Message.GetChatInfo)
	})
}
