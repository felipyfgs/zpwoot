package router

import (
	"fmt"
	"net/http"

	"zpwoot/internal/adapters/container"
	"zpwoot/internal/adapters/http/handlers"
	httpMiddleware "zpwoot/internal/adapters/http/middleware"
	"zpwoot/internal/adapters/waclient"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// NewRouter creates a new HTTP router
func NewRouter(container *container.Container) http.Handler {
	r := chi.NewRouter()

	// Basic middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Custom middleware
	r.Use(httpMiddleware.CORSMiddleware())
	r.Use(httpMiddleware.JSONMiddleware())

	// Container should already be initialized by main.go

	// Health check (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if container.GetDatabase() != nil {
			if err := container.GetDatabase().Health(); err != nil {
				http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"zpwoot"}`))
	})

	// API info (no auth required)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"zpwoot WhatsApp API","version":"1.0.0","docs":"/api/v1"}`))
	})

	// Initialize WhatsApp client and handlers
	_, sessionHandler, messageHandler := initializeWAClient(container)

	// Protected API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth middleware for API routes
		r.Use(httpMiddleware.AuthMiddleware(container.GetConfig()))

		// Session management routes
		r.Route("/sessions", func(r chi.Router) {
			r.Post("/", sessionHandler.CreateSession)           // POST /api/v1/sessions
			r.Get("/", sessionHandler.ListSessions)             // GET /api/v1/sessions
			r.Get("/{sessionId}", sessionHandler.GetSession)    // GET /api/v1/sessions/{id}
			r.Delete("/{sessionId}", sessionHandler.DeleteSession) // DELETE /api/v1/sessions/{id}

			// Session actions
			r.Post("/{sessionId}/connect", sessionHandler.ConnectSession)       // POST /api/v1/sessions/{id}/connect
			r.Post("/{sessionId}/disconnect", sessionHandler.DisconnectSession) // POST /api/v1/sessions/{id}/disconnect
			r.Get("/{sessionId}/qr", sessionHandler.GetQRCode)                  // GET /api/v1/sessions/{id}/qr
			r.Post("/{sessionId}/qr/refresh", sessionHandler.RefreshQRCode)     // POST /api/v1/sessions/{id}/qr/refresh

			// Message routes
			r.Post("/{sessionId}/messages", messageHandler.SendMessage)         // POST /api/v1/sessions/{id}/messages
			r.Get("/{sessionId}/chats", messageHandler.GetChats)                // GET /api/v1/sessions/{id}/chats
			r.Get("/{sessionId}/contacts", messageHandler.GetContacts)          // GET /api/v1/sessions/{id}/contacts
			r.Get("/{sessionId}/chat-info", messageHandler.GetChatInfo)         // GET /api/v1/sessions/{id}/chat-info?chatJid=xxx
		})

		// API documentation endpoint
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			apiDocs := map[string]interface{}{
				"name":    "zpwoot WhatsApp API",
				"version": "1.0.0",
				"endpoints": map[string]interface{}{
					"sessions": map[string]string{
						"POST /sessions":                    "Create new session",
						"GET /sessions":                     "List all sessions",
						"GET /sessions/{id}":                "Get session details",
						"DELETE /sessions/{id}":             "Delete session",
						"POST /sessions/{id}/connect":       "Connect session",
						"POST /sessions/{id}/disconnect":    "Disconnect session",
						"GET /sessions/{id}/qr":             "Get QR code",
						"POST /sessions/{id}/qr/refresh":    "Refresh QR code",
						"POST /sessions/{id}/messages":      "Send message",
						"GET /sessions/{id}/chats":          "Get chats",
						"GET /sessions/{id}/contacts":       "Get contacts",
						"GET /sessions/{id}/chat-info":      "Get chat info",
					},
				},
				"authentication": "API Key required in X-API-Key header or Authorization: Bearer {key}",
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("%+v", apiDocs)))
		})
	})

	// Migration status endpoint (no auth required for debugging)
	r.Get("/migrations/status", func(w http.ResponseWriter, r *http.Request) {
		migrator := container.GetMigrator()
		if migrator == nil {
			http.Error(w, "Migrator not available", http.StatusServiceUnavailable)
			return
		}

		migrations, err := migrator.GetMigrationStatus()
		if err != nil {
			http.Error(w, "Failed to get migration status: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf(`{"migrations":{"status":"available","count":%d}}`, len(migrations))
		w.Write([]byte(response))
	})

	return r
}

// initializeWAClient initializes the WhatsApp client and handlers
func initializeWAClient(container *container.Container) (*waclient.WAClient, *handlers.SessionHandler, *handlers.MessageHandler) {
	// Create session manager
	sessionManager := waclient.NewDBSessionManager(container.GetDatabase().DB)

	// Create WhatsApp store container
	waContainer := waclient.NewWAStoreContainer(container.GetDatabase().DB, container.GetLogger())

	// Create WhatsApp client
	waClient := waclient.NewWAClient(waContainer, container.GetLogger(), sessionManager)

	// Create message sender
	messageSender := waclient.NewMessageSender(waClient)

	// Create handlers
	sessionHandler := handlers.NewSessionHandler(waClient, container.GetLogger())
	messageHandler := handlers.NewMessageHandler(messageSender, container.GetLogger())

	return waClient, sessionHandler, messageHandler
}
