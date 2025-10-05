package handlers

import (
	"encoding/json"
	"net/http"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"

	"github.com/go-chi/chi/v5"
)

type SessionHandler struct {
	useCases       input.SessionUseCases
	sessionManager input.SessionManager
	logger         output.Logger
}

func NewSessionHandler(useCases input.SessionUseCases, sessionManager input.SessionManager, logger output.Logger) *SessionHandler {
	return &SessionHandler{
		useCases:       useCases,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// @Summary		Create WhatsApp Session
// @Description	Creates a new WhatsApp session with the specified configuration
// @Tags			Sessions
// @Accept			json
// @Produce		json
// @Param			request	body		CreateSessionRequest	true	"Session configuration"
// @Success		201		{object}	SessionResponse			"Session created successfully"
// @Failure		400		{object}	ErrorResponse			"Invalid request body or validation error"
// @Failure		409		{object}	ErrorResponse			"Session already exists"
// @Failure		500		{object}	ErrorResponse			"Internal server error"
// @Security		ApiKeyAuth
// @Router			/sessions/create [post]
func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeBadRequest, "invalid JSON body")
		return
	}

	response, err := h.useCases.CreateSession(r.Context(), &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("name", req.Name).
			Msg("Failed to create session")

		if err == dto.ErrSessionAlreadyExists {
			h.writeErrorResponse(w, http.StatusConflict, dto.ErrorCodeConflict, "session already exists")
			return
		}

		h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to create session")
		return
	}

	h.writeSuccessResponse(w, http.StatusCreated, response)
}

// @Summary		Get WhatsApp Session
// @Description	Retrieves detailed information about a specific WhatsApp session (without QR code)
// @Tags			Sessions
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"	Format(uuid)
// @Success		200			{object}	dto.SessionListInfo		"Session information (without QR code)"
// @Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
// @Failure		404			{object}	dto.ErrorResponse		"Session not found"
// @Failure		500			{object}	dto.ErrorResponse		"Internal server error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/info [get]
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	response, err := h.useCases.GetSession(r.Context(), sessionID)
	if err != nil {
		if err == dto.ErrSessionNotFound {
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "session not found")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to get session")
		}
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// @Summary		List WhatsApp Sessions
// @Description	Retrieves a list of all WhatsApp sessions with their current status (without QR codes)
// @Tags			Sessions
// @Accept			json
// @Produce		json
// @Success		200	{object}	dto.APIResponse	"List of sessions (without QR codes)"
// @Failure		500	{object}	dto.ErrorResponse		"Internal server error"
// @Security		ApiKeyAuth
// @Router			/sessions/list [get]
func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {

	pagination := &dto.PaginationRequest{
		Limit:  100,
		Offset: 0,
	}

	response, err := h.useCases.ListSessions(r.Context(), pagination)
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to list sessions")
		h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to list sessions")
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// @Summary		Connect WhatsApp Session
// @Description	Initiates connection for a WhatsApp session. Returns session info with QR code if generated.
// @Tags			Sessions
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"	Format(uuid)
// @Success		200			{object}	dto.SessionResponse		"Session info with QR code (if generated)"
// @Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
// @Failure		404			{object}	dto.ErrorResponse		"Session not found"
// @Failure		500			{object}	dto.ErrorResponse		"Connection error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/connect [post]
func (h *SessionHandler) ConnectSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	response, err := h.useCases.ConnectSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to connect session")
		if err == dto.ErrSessionNotFound {
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "session not found")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to connect session")
		}
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// @Summary		Disconnect WhatsApp Session
// @Description	Disconnects an active WhatsApp session temporarily. Credentials are kept for reconnection without QR code.
// @Tags			Sessions
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"	Format(uuid)
// @Success		200			{object}	dto.SessionResponse		"Session disconnected successfully"
// @Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
// @Failure		404			{object}	dto.ErrorResponse		"Session not found"
// @Failure		500			{object}	dto.ErrorResponse		"Disconnection error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/disconnect [post]
func (h *SessionHandler) DisconnectSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	statusResponse, err := h.useCases.DisconnectSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to disconnect session")
		if err == dto.ErrSessionNotFound {
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "session not found")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to disconnect session")
		}
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, statusResponse)
}

// @Summary		Logout WhatsApp Session
// @Description	Logs out a WhatsApp session permanently. Unlinks device from WhatsApp. Requires QR scan to reconnect.
// @Tags			Sessions
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"	Format(uuid)
// @Success		200			{object}	dto.SessionResponse		"Session logged out successfully"
// @Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
// @Failure		404			{object}	dto.ErrorResponse		"Session not found"
// @Failure		500			{object}	dto.ErrorResponse		"Logout error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/logout [post]
func (h *SessionHandler) LogoutSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	err := h.useCases.LogoutSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to logout session")
		if err == dto.ErrSessionNotFound {
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

// @Summary		Delete WhatsApp Session
// @Description	Permanently deletes a WhatsApp session and all associated data
// @Tags			Sessions
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"	Format(uuid)
// @Success		200			{object}	dto.SessionResponse		"Session deleted successfully"
// @Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
// @Failure		404			{object}	dto.ErrorResponse		"Session not found"
// @Failure		500			{object}	dto.ErrorResponse		"Deletion error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/delete [delete]
func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	err := h.useCases.DeleteSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to delete session")
		if err == dto.ErrSessionNotFound {
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

// @Summary		Get QR Code
// @Description	Retrieves the QR code for WhatsApp session authentication. Scan with WhatsApp mobile app.
// @Tags			Sessions
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"	Format(uuid)
// @Success		200			{object}	dto.QRCodeResponse		"QR code data"
// @Failure		400			{object}	dto.ErrorResponse		"Invalid session ID"
// @Failure		404			{object}	dto.ErrorResponse		"Session not found"
// @Failure		409			{object}	dto.ErrorResponse		"Session already connected"
// @Failure		500			{object}	dto.ErrorResponse		"QR code generation error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/qr [get]
func (h *SessionHandler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	response, err := h.useCases.GetQRCode(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to get QR code")

		if err == dto.ErrSessionNotFound {
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "session not found")
		} else if err.Error() == "session is already connected" {
			h.writeErrorResponse(w, http.StatusConflict, dto.ErrorCodeConflict, "session is already connected")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to get QR code")
		}
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

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

func (h *SessionHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	h.writeSuccessResponse(w, status, data)
}

func (h *SessionHandler) writeError(w http.ResponseWriter, status int, code, message string) {
	h.writeErrorResponse(w, status, code, message)
}
