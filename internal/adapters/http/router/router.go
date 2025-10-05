package router

import (
	"fmt"
	"net/http"

	"zpwoot/internal/adapters/container"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// NewRouter creates a new HTTP router
func NewRouter(container *container.Container) http.Handler {
	r := chi.NewRouter()

	// Middleware
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

	// Container should already be initialized by main.go

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if container.GetDatabase() != nil {
			if err := container.GetDatabase().Health(); err != nil {
				http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
				return
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"zpwoot"}`))
	})

	// API info
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"zpwoot API is running","version":"1.0.0"}`))
	})

	// Migration status endpoint
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		// Simple JSON response
		response := `{"migrations":["status":"available","count":` + 
			fmt.Sprintf("%d", len(migrations)) + `]}`
		w.Write([]byte(response))
	})

	return r
}
