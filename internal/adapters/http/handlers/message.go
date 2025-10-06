package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"

	"github.com/go-chi/chi/v5"
)

type MessageHandler struct {
	messageService input.MessageService
	logger         output.Logger
}

func NewMessageHandler(messageService input.MessageService, logger output.Logger) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		logger:         logger,
	}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

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
		err = h.messageService.SendTextMessage(r.Context(), sessionID, req.To, req.Text)

	case "media":
		if req.Media == nil {
			h.writeError(w, http.StatusBadRequest, "validation_error", "media is required for media messages")
			return
		}
		err = h.messageService.SendMediaMessage(r.Context(), sessionID, req.To, req.Media.ToInterfacesMediaData())

	case "location":
		if req.Location == nil {
			h.writeError(w, http.StatusBadRequest, "validation_error", "location is required for location messages")
			return
		}
		err = h.messageService.SendLocationMessage(r.Context(), sessionID, req.To,
			req.Location.Latitude, req.Location.Longitude, req.Location.Name)

	case "contact":
		if req.Contact == nil {
			h.writeError(w, http.StatusBadRequest, "validation_error", "contact is required for contact messages")
			return
		}
		contactInfo := &input.ContactInfo{
			Name:  req.Contact.Name,
			Phone: req.Contact.Phone,
		}
		err = h.messageService.SendContactMessage(r.Context(), sessionID, req.To, contactInfo)

	default:
		h.writeError(w, http.StatusBadRequest, "validation_error", "unsupported message type")
		return
	}

	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to send message")
		if waErr, ok := err.(*output.WhatsAppError); ok {
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

	response := &dto.SendMessageResponse{
		MessageID: messageID,
		Status:    "sent",
		SentAt:    time.Now(),
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

	chatInfo, err := h.messageService.GetChatInfo(r.Context(), sessionID, chatJID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("chat_jid", chatJID).
			Msg("Failed to get chat info")
		if waErr, ok := err.(*output.WhatsAppError); ok {
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

	contacts, err := h.messageService.GetContacts(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to get contacts")
		if waErr, ok := err.(*output.WhatsAppError); ok {
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

	chats, err := h.messageService.GetChats(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to get chats")
		if waErr, ok := err.(*output.WhatsAppError); ok {
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

// SendTextMessage handles sending text messages
// @Summary		Send Text Message
// @Description	Send a text message to a WhatsApp contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.SendTextMessageRequest	true	"Text message request"
// @Success		200			{object}	dto.SendMessageResponse		"Message sent successfully"
// @Failure		400			{object}	dto.ErrorResponse			"Bad request"
// @Failure		404			{object}	dto.ErrorResponse			"Session not found"
// @Failure		412			{object}	dto.ErrorResponse			"Session not connected"
// @Failure		500			{object}	dto.ErrorResponse			"Internal server error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/send/message/text [post]
func (h *MessageHandler) SendTextMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendTextMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Text == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "text is required")
		return
	}

	err := h.messageService.SendTextMessage(r.Context(), sessionID, req.To, req.Text)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.To).
			Msg("Failed to send text message")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: "msg_" + generateID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

// SendAudioMessage handles sending audio messages
// @Summary		Send Audio Message
// @Description	Send an audio message to a WhatsApp contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string						true	"Session ID"
// @Param			request		body		dto.SendAudioMessageRequest	true	"Audio message request"
// @Success		200			{object}	dto.SendMessageResponse			"Message sent successfully"
// @Failure		400			{object}	dto.ErrorResponse				"Bad request"
// @Failure		404			{object}	dto.ErrorResponse				"Session not found"
// @Failure		412			{object}	dto.ErrorResponse				"Session not connected"
// @Failure		500			{object}	dto.ErrorResponse				"Internal server error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/send/message/audio [post]
func (h *MessageHandler) SendAudioMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendAudioMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Audio == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", "audio is required")
		return
	}

	err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.To, req.Audio.ToMediaData())
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.To).
			Msg("Failed to send audio message")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: "msg_" + generateID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

// SendImageMessage handles sending image messages
// @Summary		Send Image Message
// @Description	Send an image message to a WhatsApp contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string						true	"Session ID"
// @Param			request		body		dto.SendImageMessageRequest	true	"Image message request"
// @Success		200			{object}	dto.SendMessageResponse			"Message sent successfully"
// @Failure		400			{object}	dto.ErrorResponse				"Bad request"
// @Failure		404			{object}	dto.ErrorResponse				"Session not found"
// @Failure		412			{object}	dto.ErrorResponse				"Session not connected"
// @Failure		500			{object}	dto.ErrorResponse				"Internal server error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/send/message/image [post]
func (h *MessageHandler) SendImageMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendImageMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Image == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", "image is required")
		return
	}

	err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.To, req.Image.ToMediaData())
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.To).
			Msg("Failed to send image message")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: "msg_" + generateID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

// handleMessageError handles common message sending errors
func (h *MessageHandler) handleMessageError(w http.ResponseWriter, err error) {
	if waErr, ok := err.(*output.WhatsAppError); ok {
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
}

func generateID() string {
	return "123456789"
}

// SendTextMessage handles sending text messages
// @Summary		Send Text Message
// @Description	Send a text message to a WhatsApp contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string					true	"Session ID"
// @Param			request		body		dto.SendTextMessageRequest	true	"Text message request"
// @Success		200			{object}	dto.SendMessageResponse		"Message sent successfully"
// @Failure		400			{object}	dto.ErrorResponse			"Bad request"
// @Failure		404			{object}	dto.ErrorResponse			"Session not found"
// @Failure		412			{object}	dto.ErrorResponse			"Session not connected"
// @Failure		500			{object}	dto.ErrorResponse			"Internal server error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/send/message/text [post]
func (h *MessageHandler) SendTextMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendTextMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Text == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "text is required")
		return
	}

	messageID := "msg_" + generateID()
	err := h.messageService.SendTextMessage(r.Context(), sessionID, req.To, req.Text)
	if err != nil {
		h.handleSendError(w, err, sessionID, "text")
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: messageID,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

// SendAudioMessage handles sending audio messages
// @Summary		Send Audio Message
// @Description	Send an audio message to a WhatsApp contact
// @Tags			Messages
// @Accept			json
// @Produce		json
// @Param			sessionId	path		string						true	"Session ID"
// @Param			request		body		dto.SendAudioMessageRequest	true	"Audio message request"
// @Success		200			{object}	dto.SendMessageResponse			"Message sent successfully"
// @Failure		400			{object}	dto.ErrorResponse				"Bad request"
// @Failure		404			{object}	dto.ErrorResponse				"Session not found"
// @Failure		412			{object}	dto.ErrorResponse				"Session not connected"
// @Failure		500			{object}	dto.ErrorResponse				"Internal server error"
// @Security		ApiKeyAuth
// @Router			/sessions/{sessionId}/send/message/audio [post]
func (h *MessageHandler) SendAudioMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendAudioMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Audio == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", "audio is required")
		return
	}

	messageID := "msg_" + generateID()
	err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.To, req.Audio.ToOutputMediaData())
	if err != nil {
		h.handleSendError(w, err, sessionID, "audio")
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: messageID,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}
