package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"zpwoot/internal/application/dto"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/adapters/waclient"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// SessionHandler handles WhatsApp session operations
type SessionHandler struct {
	waClient *waclient.WAClient
	logger   *logger.Logger
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(waClient *waclient.WAClient, logger *logger.Logger) *SessionHandler {
	return &SessionHandler{
		waClient: waClient,
		logger:   logger,
	}
}



// CreateSession creates a new WhatsApp session
//
//	@Summary		Create WhatsApp Session
//	@Description	Creates a new WhatsApp session with the specified configuration
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateSessionRequest	true	"Session configuration"
//	@Success		201		{object}	dto.SessionResponse			"Session created successfully"
//	@Failure		400		{object}	dto.ErrorResponse			"Invalid request body or validation error"
//	@Failure		409		{object}	dto.ErrorResponse			"Session already exists"
//	@Failure		500		{object}	dto.ErrorResponse			"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/sessions/create [post]
func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeBadRequest, "invalid JSON body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeValidation, "name is required")
		return
	}

	// Generate session ID
	sessionID := uuid.New().String()

	// Create session config using DTO conversion
	config := req.ToSessionConfig()
	config.SessionID = sessionID

	// Create session
	client, err := h.waClient.CreateSession(r.Context(), config)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", config.SessionID).
			Str("name", config.Name).
			Msg("Failed to create session")
		if waErr, ok := err.(*waclient.WAError); ok {
			h.writeErrorResponse(w, http.StatusConflict, waErr.Code, waErr.Message)
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to create session")
		}
		return
	}

	// Convert to response format using DTO
	response := dto.FromWAClient(client)

	h.writeSuccessResponse(w, http.StatusCreated, response)
}

// GetSession retrieves a session by ID
//
//	@Summary		Get WhatsApp Session
//	@Description	Retrieves detailed information about a specific WhatsApp session
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"	Format(uuid)
//	@Success		200			{object}	dto.SessionResponse		"Session information"
//	@Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
//	@Failure		404			{object}	dto.ErrorResponse		"Session not found"
//	@Failure		500			{object}	dto.ErrorResponse		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/sessions/{sessionId}/info [get]
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	client, err := h.waClient.GetSession(r.Context(), sessionID)
	if err != nil {
		if err == waclient.ErrSessionNotFound {
			h.writeError(w, http.StatusNotFound, "session_not_found", "session not found")
		} else {
			h.writeError(w, http.StatusInternalServerError, "internal_error", "failed to get session")
		}
		return
	}

	sessionInfo := h.clientToSessionInfo(client)

	response := &waclient.SessionResponse{
		Success: true,
		Session: sessionInfo,
	}

	h.writeJSON(w, http.StatusOK, response)
}

// ListSessions retrieves all sessions
//
//	@Summary		List WhatsApp Sessions
//	@Description	Retrieves a list of all WhatsApp sessions with their current status
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		dto.SessionResponse	"List of sessions"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/sessions/list [get]
func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	clients, err := h.waClient.ListSessions(r.Context())
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to list sessions")
		h.writeError(w, http.StatusInternalServerError, "internal_error", "failed to list sessions")
		return
	}

	sessions := make([]*waclient.SessionInfo, len(clients))
	for i, client := range clients {
		sessions[i] = h.clientToSessionInfo(client)
	}

	response := &waclient.SessionResponse{
		Success:  true,
		Sessions: sessions,
	}

	h.writeJSON(w, http.StatusOK, response)
}

// ConnectSession connects a session to WhatsApp
//
//	@Summary		Connect WhatsApp Session
//	@Description	Initiates connection for a WhatsApp session. Use QR code endpoint to get QR for scanning.
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"	Format(uuid)
//	@Success		200			{object}	dto.SessionResponse		"Session connection initiated"
//	@Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
//	@Failure		404			{object}	dto.ErrorResponse		"Session not found"
//	@Failure		500			{object}	dto.ErrorResponse		"Connection error"
//	@Security		ApiKeyAuth
//	@Router			/sessions/{sessionId}/connect [post]
func (h *SessionHandler) ConnectSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	err := h.waClient.ConnectSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to connect session")
		if err == waclient.ErrSessionNotFound {
			h.writeError(w, http.StatusNotFound, "session_not_found", "session not found")
		} else {
			h.writeError(w, http.StatusInternalServerError, "connection_error", "failed to connect session")
		}
		return
	}

	response := &waclient.SessionResponse{
		Success: true,
	}

	h.writeJSON(w, http.StatusOK, response)
}

// DisconnectSession disconnects a session
//
//	@Summary		Disconnect WhatsApp Session
//	@Description	Disconnects an active WhatsApp session
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"	Format(uuid)
//	@Success		200			{object}	dto.SessionResponse		"Session disconnected successfully"
//	@Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
//	@Failure		404			{object}	dto.ErrorResponse		"Session not found"
//	@Failure		500			{object}	dto.ErrorResponse		"Disconnection error"
//	@Security		ApiKeyAuth
//	@Router			/sessions/{sessionId}/disconnect [post]
func (h *SessionHandler) DisconnectSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	err := h.waClient.DisconnectSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to disconnect session")
		if err == waclient.ErrSessionNotFound {
			h.writeError(w, http.StatusNotFound, "session_not_found", "session not found")
		} else {
			h.writeError(w, http.StatusInternalServerError, "disconnection_error", "failed to disconnect session")
		}
		return
	}

	response := &waclient.SessionResponse{
		Success: true,
	}

	h.writeJSON(w, http.StatusOK, response)
}

// DeleteSession deletes a session
//
//	@Summary		Delete WhatsApp Session
//	@Description	Permanently deletes a WhatsApp session and all associated data
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"	Format(uuid)
//	@Success		200			{object}	dto.SessionResponse		"Session deleted successfully"
//	@Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
//	@Failure		404			{object}	dto.ErrorResponse		"Session not found"
//	@Failure		500			{object}	dto.ErrorResponse		"Deletion error"
//	@Security		ApiKeyAuth
//	@Router			/sessions/{sessionId}/delete [delete]
func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	err := h.waClient.DeleteSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to delete session")
		if err == waclient.ErrSessionNotFound {
			h.writeError(w, http.StatusNotFound, "session_not_found", "session not found")
		} else {
			h.writeError(w, http.StatusInternalServerError, "deletion_error", "failed to delete session")
		}
		return
	}

	response := &waclient.SessionResponse{
		Success: true,
	}

	h.writeJSON(w, http.StatusOK, response)
}

// GetQRCode retrieves the QR code for a session
//
//	@Summary		Get QR Code
//	@Description	Retrieves the QR code for WhatsApp session authentication. Scan with WhatsApp mobile app.
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"	Format(uuid)
//	@Success		200			{object}	dto.QRCodeResponse		"QR code data"
//	@Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
//	@Failure		404			{object}	dto.ErrorResponse		"Session not found"
//	@Failure		409			{object}	dto.ErrorResponse		"Session already connected"
//	@Failure		500			{object}	dto.ErrorResponse		"QR code generation error"
//	@Security		ApiKeyAuth
//	@Router			/sessions/{sessionId}/qr [get]
func (h *SessionHandler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	qrEvent, err := h.waClient.GetQRCodeForSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to get QR code")
		if waErr, ok := err.(*waclient.WAError); ok {
			status := http.StatusInternalServerError
			if waErr.Code == "SESSION_NOT_FOUND" {
				status = http.StatusNotFound
			} else if waErr.Code == "ALREADY_CONNECTED" {
				status = http.StatusConflict
			}
			h.writeError(w, status, waErr.Code, waErr.Message)
		} else {
			h.writeError(w, http.StatusInternalServerError, "qr_error", "failed to get QR code")
		}
		return
	}

	response := &waclient.SessionResponse{
		Success: true,
		QRCode:  qrEvent.Base64,
	}

	h.writeJSON(w, http.StatusOK, response)
}

// RefreshQRCode forces a refresh of the QR code
func (h *SessionHandler) RefreshQRCode(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	qrEvent, err := h.waClient.RefreshQRCode(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to refresh QR code")
		if waErr, ok := err.(*waclient.WAError); ok {
			status := http.StatusInternalServerError
			if waErr.Code == "SESSION_NOT_FOUND" {
				status = http.StatusNotFound
			} else if waErr.Code == "ALREADY_CONNECTED" {
				status = http.StatusConflict
			}
			h.writeError(w, status, waErr.Code, waErr.Message)
		} else {
			h.writeError(w, http.StatusInternalServerError, "qr_refresh_error", "failed to refresh QR code")
		}
		return
	}

	response := &waclient.SessionResponse{
		Success: true,
		QRCode:  qrEvent.Base64,
	}

	h.writeJSON(w, http.StatusOK, response)
}

// Helper methods
func (h *SessionHandler) clientToSessionInfo(client *waclient.Client) *waclient.SessionInfo {
	deviceJID := ""
	if client.WAClient.Store.ID != nil {
		deviceJID = client.WAClient.Store.ID.String()
	}

	return &waclient.SessionInfo{
		ID:          client.SessionID,
		Name:        client.Name,
		DeviceJID:   deviceJID,
		Status:      client.Status,
		Connected:   client.Status == waclient.StatusConnected,
		QRCode:      client.QRCode,
		QRExpiresAt: client.QRExpiresAt,
		ConnectedAt: client.ConnectedAt,
		LastSeen:    client.LastSeen,
		CreatedAt:   time.Now(), // Would need to be stored in client
		UpdatedAt:   time.Now(),
	}
}

// Helper methods for HTTP responses

func (h *SessionHandler) writeSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := dto.NewSuccessResponse(data)
	json.NewEncoder(w).Encode(response)
}

func (h *SessionHandler) writeErrorResponse(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := dto.NewErrorResponse(code, message)
	json.NewEncoder(w).Encode(response)
}

func (h *SessionHandler) writeValidationErrorResponse(w http.ResponseWriter, field, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	response := dto.NewValidationErrorResponse(field, message)
	json.NewEncoder(w).Encode(response)
}

// Legacy methods for backward compatibility (deprecated)
func (h *SessionHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	h.writeSuccessResponse(w, status, data)
}

func (h *SessionHandler) writeError(w http.ResponseWriter, status int, code, message string) {
	h.writeErrorResponse(w, status, code, message)
}
