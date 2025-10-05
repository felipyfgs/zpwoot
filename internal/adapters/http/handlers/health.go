package handlers

import (
	"encoding/json"
	"net/http"

	"zpwoot/internal/adapters/database"
	"zpwoot/internal/core/application/dto"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	database *database.Database
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		database: db,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version,omitempty"`
}

// InfoResponse represents the service info response
type InfoResponse struct {
	Message string `json:"message"`
	Version string `json:"version"`
	Service string `json:"service"`
}

// Health handles the health check endpoint
// @Summary		Health Check
// @Description	Check if the service and database are healthy
// @Tags			Health
// @Produce		json
// @Success		200	{object}	HealthResponse	"Service is healthy"
// @Failure		503	{object}	ErrorResponse	"Service is unhealthy"
// @Router			/health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check database health
	if h.database != nil {
		if err := h.database.Health(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Error:   "database_unhealthy",
				Message: "Database connection is unhealthy",
			})
			return
		}
	}

	// Service is healthy
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{
		Status:  "ok",
		Service: "zpwoot",
		Version: "1.0.0",
	})
}

// Info handles the service info endpoint
// @Summary		Service Information
// @Description	Get basic information about the service
// @Tags			Health
// @Produce		json
// @Success		200	{object}	InfoResponse	"Service information"
// @Router			/ [get]
func (h *HealthHandler) Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(InfoResponse{
		Message: "zpwoot WhatsApp API is running",
		Version: "1.0.0",
		Service: "zpwoot",
	})
}
