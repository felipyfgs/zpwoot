package handlers

import (
	"encoding/json"
	"net/http"

	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/adapters/waclient"

	"github.com/go-chi/chi/v5"
)

type MessageHandler struct {
	messageSender *waclient.MessageSenderImpl
	logger        *logger.Logger
}

func NewMessageHandler(messageSender *waclient.MessageSenderImpl, logger *logger.Logger) *MessageHandler {
	return &MessageHandler{
		messageSender: messageSender,
		logger:        logger,
	}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req waclient.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	req.SessionID = sessionID

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Type == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "type is required")
		return
	}

	var err error
	messageID := "msg_" + generateID()

	switch req.Type {
	case "text":
		if req.Text == "" {
			h.writeError(w, http.StatusBadRequest, "validation_error", "text is required for text messages")
			return
		}
		err = h.messageSender.SendTextMessage(r.Context(), sessionID, req.To, req.Text)

	case "media":
		if req.Media == nil {
			h.writeError(w, http.StatusBadRequest, "validation_error", "media is required for media messages")
			return
		}
		err = h.messageSender.SendMediaMessage(r.Context(), sessionID, req.To, req.Media)

	case "location":
		if req.Location == nil {
			h.writeError(w, http.StatusBadRequest, "validation_error", "location is required for location messages")
			return
		}
		err = h.messageSender.SendLocationMessage(r.Context(), sessionID, req.To,
			req.Location.Latitude, req.Location.Longitude, req.Location.Name)

	case "contact":
		if req.Contact == nil {
			h.writeError(w, http.StatusBadRequest, "validation_error", "contact is required for contact messages")
			return
		}
		err = h.messageSender.SendContactMessage(r.Context(), sessionID, req.To, req.Contact)

	default:
		h.writeError(w, http.StatusBadRequest, "validation_error", "unsupported message type")
		return
	}

	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to send message")
		if waErr, ok := err.(*waclient.WAError); ok {
			status := http.StatusInternalServerError
			if waErr.Code == "SESSION_NOT_FOUND" {
				status = http.StatusNotFound
			} else if waErr.Code == "NOT_CONNECTED" {
				status = http.StatusPreconditionFailed
			} else if waErr.Code == "INVALID_JID" {
				status = http.StatusBadRequest
			}
			h.writeError(w, status, waErr.Code, waErr.Message)
		} else {
			h.writeError(w, http.StatusInternalServerError, "send_error", "failed to send message")
		}
		return
	}

	response := &waclient.MessageResponse{
		Success:   true,
		MessageID: messageID,
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) GetChatInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	chatJID := r.URL.Query().Get("chatJid")

	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	if chatJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "chatJid is required")
		return
	}

	chatInfo, err := h.messageSender.GetChatInfo(r.Context(), sessionID, chatJID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("chat_jid", chatJID).
			Msg("Failed to get chat info")
		if waErr, ok := err.(*waclient.WAError); ok {
			status := http.StatusInternalServerError
			if waErr.Code == "SESSION_NOT_FOUND" {
				status = http.StatusNotFound
			} else if waErr.Code == "NOT_CONNECTED" {
				status = http.StatusPreconditionFailed
			} else if waErr.Code == "INVALID_JID" {
				status = http.StatusBadRequest
			}
			h.writeError(w, status, waErr.Code, waErr.Message)
		} else {
			h.writeError(w, http.StatusInternalServerError, "chat_info_error", "failed to get chat info")
		}
		return
	}

	h.writeJSON(w, http.StatusOK, chatInfo)
}

func (h *MessageHandler) GetContacts(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	contacts, err := h.messageSender.GetContacts(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to get contacts")
		if waErr, ok := err.(*waclient.WAError); ok {
			status := http.StatusInternalServerError
			if waErr.Code == "SESSION_NOT_FOUND" {
				status = http.StatusNotFound
			} else if waErr.Code == "NOT_CONNECTED" {
				status = http.StatusPreconditionFailed
			}
			h.writeError(w, status, waErr.Code, waErr.Message)
		} else {
			h.writeError(w, http.StatusInternalServerError, "contacts_error", "failed to get contacts")
		}
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"contacts": contacts,
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	chats, err := h.messageSender.GetChats(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to get chats")
		if waErr, ok := err.(*waclient.WAError); ok {
			status := http.StatusInternalServerError
			if waErr.Code == "SESSION_NOT_FOUND" {
				status = http.StatusNotFound
			} else if waErr.Code == "NOT_CONNECTED" {
				status = http.StatusPreconditionFailed
			}
			h.writeError(w, status, waErr.Code, waErr.Message)
		} else {
			h.writeError(w, http.StatusInternalServerError, "chats_error", "failed to get chats")
		}
		return
	}

	response := map[string]interface{}{
		"success": true,
		"chats":   chats,
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *MessageHandler) writeError(w http.ResponseWriter, status int, code, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	})
}

func generateID() string {
	return "123456789"
}
