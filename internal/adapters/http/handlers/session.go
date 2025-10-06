package handlers

import (
	"encoding/json"
	"errors"
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
// @Router			/sessions [post]
func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeBadRequest, "Invalid JSON body")
		return
	}

	if req.Name == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeValidation, "Session name is required")
		return
	}

	if len(req.Name) < 3 {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeValidation, "Session name must be at least 3 characters")
		return
	}

	if len(req.Name) > 50 {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeValidation, "Session name must be less than 50 characters")
		return
	}

	response, err := h.useCases.CreateSession(r.Context(), &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("name", req.Name).
			Msg("Failed to create session")

		if errors.Is(err, dto.ErrSessionAlreadyExists) {
			h.writeErrorResponse(w, http.StatusConflict, dto.ErrorCodeConflict, "A session with this name already exists")
			return
		}

		h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "Failed to create session")
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
// @Router			/sessions/{sessionId} [get]
func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	response, err := h.useCases.GetSession(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, dto.ErrSessionNotFound) {
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
// @Router			/sessions [get]
func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {

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
// @Description	Connects a WhatsApp session. If already connected, returns current status with appropriate message.
// @Tags			Sessions
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string						true	"Session ID"	Format(uuid)
// @Success		200			{object}	dto.SessionStatusResponse	"Session connected successfully or already connected"
// @Failure		400			{object}	dto.ErrorResponse			"Invalid session ID"
// @Failure		404			{object}	dto.ErrorResponse			"Session not found"
// @Failure		500			{object}	dto.ErrorResponse			"Connection error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/connect [post]
func (h *SessionHandler) Connect(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, dto.ErrSessionNotFound) {
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "session not found")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "failed to connect session")
		}
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// @Summary		Disconnect WhatsApp Session
// @Description	Disconnects an active WhatsApp session temporarily. If already disconnected, returns current status with appropriate message.
// @Tags			Sessions
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string						true	"Session ID"	Format(uuid)
// @Success		200			{object}	dto.SessionStatusResponse	"Session disconnected successfully or already disconnected"
// @Failure		400			{object}	dto.ErrorResponse			"Invalid session ID"
// @Failure		404			{object}	dto.ErrorResponse			"Session not found"
// @Failure		500			{object}	dto.ErrorResponse			"Disconnection error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/disconnect [post]
func (h *SessionHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, dto.ErrSessionNotFound) {
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
func (h *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, dto.ErrSessionNotFound) {
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "Session not found")
		} else if err.Error() == "session is already logged out" {
			h.writeErrorResponse(w, http.StatusConflict, dto.ErrorCodeConflict, "Session is already logged out")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "Failed to logout session")
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
// @Router			/sessions/{sessionId} [delete]
func (h *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, dto.ErrSessionNotFound) {
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
func (h *SessionHandler) QRCode(w http.ResponseWriter, r *http.Request) {
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

		if errors.Is(err, dto.ErrSessionNotFound) {
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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode success response")
	}
}

func (h *SessionHandler) writeErrorResponse(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := dto.NewErrorResponse(code, message)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode error response")
	}
}

func (h *SessionHandler) writeError(w http.ResponseWriter, status int, code, message string) {
	h.writeErrorResponse(w, status, code, message)
}

// @Summary      Pair phone
// @Description  Pair a phone number without QR code
// @Tags         Sessions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                    true  "Session ID"
// @Param        request     body      dto.PairPhoneRequest      true  "Phone number"
// @Success      200  {object}  dto.PairPhoneResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/pair [post]
func (h *SessionHandler) PairPhone(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	var req dto.PairPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeBadRequest, "Invalid JSON body")
		return
	}

	if req.Phone == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, dto.ErrorCodeValidation, "phone is required")
		return
	}

	response, err := h.useCases.PairPhone(r.Context(), sessionID, req.Phone)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("phone", req.Phone).
			Msg("Failed to pair phone")

		if err.Error() == "session not found" {
			h.writeErrorResponse(w, http.StatusNotFound, dto.ErrorCodeNotFound, "Session not found")
			return
		}

		if err.Error() == "session is already paired" {
			h.writeErrorResponse(w, http.StatusConflict, "ALREADY_PAIRED", "Session is already paired")
			return
		}

		h.writeErrorResponse(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "Failed to pair phone")
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("phone", req.Phone).
		Str("linking_code", response.LinkingCode).
		Msg("Phone paired successfully")

	h.writeSuccessResponse(w, http.StatusOK, response)
}
