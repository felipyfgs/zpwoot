package handlers

import (
	"encoding/json"
	"net/http"

	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"

	"github.com/go-chi/chi/v5"
)


type WebhookHandler struct {
	webhookUseCases input.WebhookUseCases
	logger          *logger.Logger
}


func NewWebhookHandler(webhookUseCases input.WebhookUseCases, logger *logger.Logger) *WebhookHandler {
	return &WebhookHandler{
		webhookUseCases: webhookUseCases,
		logger:          logger,
	}
}


// @Summary		Configure Webhook
// @Description	Configure or update webhook for a session
// @Tags			Webhooks
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string						true	"Session ID"
// @Param			request		body		dto.CreateWebhookRequest	true	"Webhook configuration"
// @Success		200			{object}	dto.WebhookResponse			"Webhook configured successfully"
// @Failure		400			{object}	dto.ErrorResponse			"Invalid request"
// @Failure		500			{object}	dto.ErrorResponse			"Internal server error"
// @Router			/sessions/{sessionId}/webhooks [post]
// @Security		ApiKeyAuth
func (h *WebhookHandler) SetWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	var req dto.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode webhook request")
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeBadRequest, "Invalid JSON body")
		return
	}


	response, err := h.webhookUseCases.Upsert(r.Context(), sessionID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to set webhook")
		h.writeError(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, err.Error())
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("webhook_url", req.URL).
		Msg("Webhook configured successfully")

	h.writeJSON(w, http.StatusOK, response)
}


// @Summary		Get Webhook Configuration
// @Description	Get webhook configuration for a session
// @Tags			Webhooks
// @Produce		json
// @Param			sessionId	path		string				true	"Session ID"
// @Success		200			{object}	dto.WebhookResponse	"Webhook configuration"
// @Failure		404			{object}	dto.ErrorResponse	"Webhook not found"
// @Failure		500			{object}	dto.ErrorResponse	"Internal server error"
// @Router			/sessions/{sessionId}/webhooks [get]
// @Security		ApiKeyAuth
func (h *WebhookHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	response, err := h.webhookUseCases.Get(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to get webhook")
		if err.Error() == "webhook not found" {
			h.writeError(w, http.StatusNotFound, dto.ErrorCodeNotFound, "Webhook not found for this session")
			return
		}
		h.writeError(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}


// @Summary		Delete Webhook Configuration
// @Description	Delete webhook configuration for a session
// @Tags			Webhooks
// @Produce		json
// @Param			sessionId	path		string			true	"Session ID"
// @Success		200			{object}	dto.APIResponse	"Webhook deleted successfully"
// @Failure		404			{object}	dto.ErrorResponse	"Webhook not found"
// @Failure		500			{object}	dto.ErrorResponse	"Internal server error"
// @Router			/sessions/{sessionId}/webhooks [delete]
// @Security		ApiKeyAuth
func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	err := h.webhookUseCases.Delete(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to delete webhook")
		if err.Error() == "webhook not found" {
			h.writeError(w, http.StatusNotFound, dto.ErrorCodeNotFound, "Webhook not found for this session")
			return
		}
		h.writeError(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, err.Error())
		return
	}

	h.logger.Info().Str("session_id", sessionID).Msg("Webhook deleted successfully")

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Webhook deleted successfully",
	})
}


// @Summary		List Available Events
// @Description	List all available webhook event types
// @Tags			Webhooks
// @Produce		json
// @Success		200	{object}	dto.ListEventsResponse	"Available events"
// @Failure		500	{object}	dto.ErrorResponse		"Internal server error"
// @Router			/webhooks/events [get]
// @Security		ApiKeyAuth
func (h *WebhookHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	response, err := h.webhookUseCases.ListEvents(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list webhook events")
		h.writeError(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, response)
}


func (h *WebhookHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode JSON response")
	}
}


func (h *WebhookHandler) writeError(w http.ResponseWriter, statusCode int, errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := dto.ErrorResponse{
		Error:   errorCode,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode error response")
	}
}
