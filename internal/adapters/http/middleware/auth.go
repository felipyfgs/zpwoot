package middleware

import (
	"net/http"
	"strings"

	"zpwoot/internal/adapters/config"
)

// AuthMiddleware provides API key authentication
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health check and root endpoints
			if r.URL.Path == "/health" || r.URL.Path == "/" {
				next.ServeHTTP(w, r)
				return
			}

			// Get API key from header
			// Priority: Authorization header (direct value) > X-API-Key header
			apiKey := r.Header.Get("Authorization")

			// If Authorization header is empty or has Bearer prefix, try X-API-Key
			if apiKey == "" || strings.HasPrefix(apiKey, "Bearer ") {
				// Try X-API-Key header
				xApiKey := r.Header.Get("X-API-Key")
				if xApiKey != "" {
					apiKey = xApiKey
				} else if strings.HasPrefix(apiKey, "Bearer ") {
					// Support Bearer token format as fallback
					apiKey = strings.TrimPrefix(apiKey, "Bearer ")
				}
			}

			// Check API key
			if apiKey == "" || apiKey != cfg.APIKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"unauthorized","message":"invalid or missing API key. Use 'Authorization: YOUR_API_KEY' or 'X-API-Key: YOUR_API_KEY' header"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// JSONMiddleware sets JSON content type for API responses
func JSONMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}
