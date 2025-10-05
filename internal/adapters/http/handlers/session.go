package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/adapters/waclient"
	"zpwoot/internal/application/dto"

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

	// Auto-connect and generate QR code if requested
	if req.GenerateQRCode {
		err = h.waClient.ConnectSession(r.Context(), sessionID)
		if err != nil {
			h.logger.Warn().
				Err(err).
				Str("session_id", sessionID).
				Msg("Failed to auto-connect session after creation")
			// Don't fail the request, just log the warning
		} else {
			// Wait a bit for QR code to be generated
			time.Sleep(500 * time.Millisecond)

			// Refresh client to get updated QR code
			client, err = h.waClient.GetSession(r.Context(), sessionID)
			if err != nil {
				h.logger.Warn().
					Err(err).
					Str("session_id", sessionID).
					Msg("Failed to refresh session after QR generation")
			}
		}
	}

	// Convert to response format using DTO
	response := dto.FromWAClient(client)

	h.writeSuccessResponse(w, http.StatusCreated, response)
}

// GetSession retrieves a session by ID
//
//	@Summary		Get WhatsApp Session
//	@Description	Retrieves detailed information about a specific WhatsApp session (without QR code)
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"	Format(uuid)
//	@Success		200			{object}	dto.SessionListInfo		"Session information (without QR code)"
//	@Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
//	@Failure		404			{object}	dto.ErrorResponse		"Session not found"
//	@Failure		500			{object}	dto.ErrorResponse		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/sessions/{sessionId}/info [get]
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	client, err := h.waClient.GetSession(r.Context(), sessionID)
	if err != nil {
		if err == waclient.ErrSessionNotFound {
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "session not found")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to get session")
		}
		return
	}

	// Convert to SessionResponse then to SessionListInfo (removes QR code)
	sessionResp := dto.FromWAClient(client)
	sessionInfo := sessionResp.ToListInfo()

	h.writeSuccessResponse(w, http.StatusOK, sessionInfo)
}

// ListSessions retrieves all sessions
//
//	@Summary		List WhatsApp Sessions
//	@Description	Retrieves a list of all WhatsApp sessions with their current status (without QR codes)
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.SessionListResponse	"List of sessions (without QR codes)"
//	@Failure		500	{object}	dto.ErrorResponse		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/sessions/list [get]
func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	clients, err := h.waClient.ListSessions(r.Context())
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to list sessions")
		h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to list sessions")
		return
	}

	// Convert to SessionListResponse (excludes QR codes)
	response := dto.FromWAClientList(clients)

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// ConnectSession connects a session to WhatsApp
//
//	@Summary		Connect WhatsApp Session
//	@Description	Initiates connection for a WhatsApp session. Returns session info with QR code if generated.
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"	Format(uuid)
//	@Success		200			{object}	dto.SessionResponse		"Session info with QR code (if generated)"
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
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "session not found")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to connect session")
		}
		return
	}

	// Wait a bit for QR code to be generated
	time.Sleep(500 * time.Millisecond)

	// Get updated session with QR code
	client, err := h.waClient.GetSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to get session after connect")
		// Return simple response if we can't get the session
		response := &dto.SessionActionResponse{
			SessionID: sessionID,
			Action:    "connect",
			Status:    "success",
			Message:   "Session connection initiated",
		}
		h.writeSuccessResponse(w, http.StatusOK, response)
		return
	}

	// Return full session response with QR code
	response := dto.FromWAClient(client)
	h.writeSuccessResponse(w, http.StatusOK, response)
}

// DisconnectSession disconnects a session temporarily (keeps credentials)
//
//	@Summary		Disconnect WhatsApp Session
//	@Description	Disconnects an active WhatsApp session temporarily. Credentials are kept for reconnection without QR code.
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
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "session not found")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to disconnect session")
		}
		return
	}

	response := &dto.SessionActionResponse{
		SessionID: sessionID,
		Action:    "disconnect",
		Status:    "success",
		Message:   "Session disconnected (credentials kept for reconnection)",
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// LogoutSession logs out a session permanently (unlinks device)
//
//	@Summary		Logout WhatsApp Session
//	@Description	Logs out a WhatsApp session permanently. Unlinks device from WhatsApp. Requires QR scan to reconnect.
//	@Tags			Sessions
//	@Accept			json
//	@Produce		json
//	@Param			sessionId	path		string					true	"Session ID"	Format(uuid)
//	@Success		200			{object}	dto.SessionResponse		"Session logged out successfully"
//	@Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
//	@Failure		404			{object}	dto.ErrorResponse		"Session not found"
//	@Failure		500			{object}	dto.ErrorResponse		"Logout error"
//	@Security		ApiKeyAuth
//	@Router			/sessions/{sessionId}/logout [post]
func (h *SessionHandler) LogoutSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	err := h.waClient.LogoutSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to logout session")
		if err == waclient.ErrSessionNotFound {
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "session not found")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to logout session")
		}
		return
	}

	response := &dto.SessionActionResponse{
		SessionID: sessionID,
		Action:    "logout",
		Status:    "success",
		Message:   "Session logged out (device unlinked from WhatsApp)",
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
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
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "session not found")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to delete session")
		}
		return
	}

	response := &dto.SessionActionResponse{
		SessionID: sessionID,
		Action:    "delete",
		Status:    "success",
		Message:   "Session deleted successfully",
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
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
			errorCode := dto.ErrorCodeInternalError
			if waErr.Code == "SESSION_NOT_FOUND" {
				status = http.StatusNotFound
				errorCode = dto.ErrorCodeNotFound
			} else if waErr.Code == "ALREADY_CONNECTED" {
				status = http.StatusConflict
				errorCode = dto.ErrorCodeConflict
			}
			h.writeErrorResponse(w, status, errorCode, waErr.Message)
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to get QR code")
		}
		return
	}

	response := dto.FromQREvent(qrEvent)

	h.writeSuccessResponse(w, http.StatusOK, response)
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
			errorCode := dto.ErrorCodeInternalError
			if waErr.Code == "SESSION_NOT_FOUND" {
				status = http.StatusNotFound
				errorCode = dto.ErrorCodeNotFound
			} else if waErr.Code == "ALREADY_CONNECTED" {
				status = http.StatusConflict
				errorCode = dto.ErrorCodeConflict
			}
			h.writeErrorResponse(w, status, errorCode, waErr.Message)
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to refresh QR code")
		}
		return
	}

	response := dto.FromQREvent(qrEvent)

	h.writeSuccessResponse(w, http.StatusOK, response)
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
