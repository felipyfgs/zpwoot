package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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

// CreateSessionRequest represents a session creation request
type CreateSessionRequest struct {
	Name          string                    `json:"name"`
	WebhookURL    string                    `json:"webhookUrl,omitempty"`
	Events        []waclient.EventType      `json:"events,omitempty"`
	ProxyConfig   map[string]string         `json:"proxyConfig,omitempty"`
	AutoReconnect bool                      `json:"autoReconnect"`
}

// CreateSession creates a new WhatsApp session
func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "name is required")
		return
	}

	// Generate session ID
	sessionID := uuid.New().String()

	// Create session config
	config := &waclient.SessionConfig{
		SessionID:     sessionID,
		Name:          req.Name,
		WebhookURL:    req.WebhookURL,
		Events:        req.Events,
		ProxyConfig:   req.ProxyConfig,
		AutoReconnect: req.AutoReconnect,
	}

	// Create session
	client, err := h.waClient.CreateSession(r.Context(), config)
	if err != nil {
		h.logger.Errorf("Failed to create session: %v", err)
		if waErr, ok := err.(*waclient.WAError); ok {
			h.writeError(w, http.StatusConflict, waErr.Code, waErr.Message)
		} else {
			h.writeError(w, http.StatusInternalServerError, "internal_error", "failed to create session")
		}
		return
	}

	// Convert to response format
	sessionInfo := h.clientToSessionInfo(client)

	response := &waclient.SessionResponse{
		Success: true,
		Session: sessionInfo,
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// GetSession retrieves a session by ID
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
func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	clients, err := h.waClient.ListSessions(r.Context())
	if err != nil {
		h.logger.Errorf("Failed to list sessions: %v", err)
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
func (h *SessionHandler) ConnectSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	err := h.waClient.ConnectSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Errorf("Failed to connect session: %v", err)
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
func (h *SessionHandler) DisconnectSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	err := h.waClient.DisconnectSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Errorf("Failed to disconnect session: %v", err)
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
func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	err := h.waClient.DeleteSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Errorf("Failed to delete session: %v", err)
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
func (h *SessionHandler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	qrEvent, err := h.waClient.GetQRCodeForSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Errorf("Failed to get QR code: %v", err)
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
		h.logger.Errorf("Failed to refresh QR code: %v", err)
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

func (h *SessionHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *SessionHandler) writeError(w http.ResponseWriter, status int, code, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	})
}
