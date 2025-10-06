package handlers

import (
	"encoding/json"
	"net/http"

	"zpwoot/internal/adapters/database"
	"zpwoot/internal/core/application/dto"
)

type HealthHandler struct {
	database *database.Database
}

func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		database: db,
	}
}

type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"zpwoot"`
	Version string `json:"version,omitempty" example:"1.0.0"`
} //@name HealthResponse

type SystemInfoResponse struct {
	Message string `json:"message" example:"zpwoot WhatsApp API is running"`
	Version string `json:"version" example:"1.0.0"`
	Service string `json:"service" example:"zpwoot"`
} //@name SystemInfoResponse

// @Summary		Health Check
// @Description	Check if the service and database are healthy
// @Tags			Health
// @Produce		json
// @Success		200	{object}	HealthResponse	"Service is healthy"
// @Failure		503	{object}	ErrorResponse	"Service is unhealthy"
// @Router			/health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{
		Status:  "ok",
		Service: "zpwoot",
		Version: "1.0.0",
	})
}

// @Summary		Service Information
// @Description	Get basic information about the service
// @Tags			Health
// @Produce		json
// @Success		200	{object}	SystemInfoResponse	"Service information"
// @Router			/ [get]
func (h *HealthHandler) Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(SystemInfoResponse{
		Message: "zpwoot WhatsApp API is running",
		Version: "1.0.0",
		Service: "zpwoot",
	})
}
