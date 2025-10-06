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

func (h *MessageHandler) SendText(w http.ResponseWriter, r *http.Request) {
	h.SendTextMessage(w, r)
}

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

func (h *MessageHandler) SendVideo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendVideoMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Video == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", "video is required")
		return
	}

	messageID := "msg_" + generateID()
	err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.To, req.Video.ToOutputMediaData())
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: messageID,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) SendDocument(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendDocumentMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Document == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", "document is required")
		return
	}

	messageID := "msg_" + generateID()
	err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.To, req.Document.ToOutputMediaData())
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: messageID,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) SendSticker(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendStickerMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Sticker == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sticker is required")
		return
	}

	messageID := "msg_" + generateID()
	err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.To, req.Sticker.ToOutputMediaData())
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: messageID,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) SendAudio(w http.ResponseWriter, r *http.Request) {
	h.SendAudioMessage(w, r)
}

func (h *MessageHandler) SendImage(w http.ResponseWriter, r *http.Request) {
	h.SendImageMessage(w, r)
}

func (h *MessageHandler) SendLocation(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendLocationMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	messageID := "msg_" + generateID()
	err := h.messageService.SendLocationMessage(r.Context(), sessionID, req.To, req.Latitude, req.Longitude, req.Name)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: messageID,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) SendContact(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendContactMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Contact == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", "contact is required")
		return
	}

	contactInfo := &input.ContactInfo{
		Name:  req.Contact.Name,
		Phone: req.Contact.Phone,
		VCard: req.Contact.VCard,
	}

	messageID := "msg_" + generateID()
	err := h.messageService.SendContactMessage(r.Context(), sessionID, req.To, contactInfo)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: messageID,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) SendContactsArray(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendContactsArrayMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if len(req.Contacts) == 0 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "contacts array is required")
		return
	}

	firstContact := req.Contacts[0]
	contactInfo := &input.ContactInfo{
		Name:  firstContact.Name,
		Phone: firstContact.Phone,
		VCard: firstContact.VCard,
	}

	messageID := "msg_" + generateID()
	err := h.messageService.SendContactMessage(r.Context(), sessionID, req.To, contactInfo)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: messageID,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) SendReaction(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendReactionMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.MessageID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "messageId is required")
		return
	}

	if req.Reaction == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "reaction is required")
		return
	}

	err := h.messageService.SendReactionMessage(r.Context(), sessionID, req.To, req.MessageID, req.Reaction)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.To).
			Str("message_id", req.MessageID).
			Msg("Failed to send reaction")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: req.MessageID,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) SendPoll(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendPollMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "name is required")
		return
	}

	if len(req.Options) < 2 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "at least 2 options are required")
		return
	}

	if len(req.Options) > 12 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "maximum 12 options allowed")
		return
	}

	if req.SelectableOptionsCount < 1 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "selectableOptionsCount must be at least 1")
		return
	}

	if req.SelectableOptionsCount > len(req.Options) {
		h.writeError(w, http.StatusBadRequest, "validation_error", "selectableOptionsCount cannot exceed number of options")
		return
	}

	err := h.messageService.SendPollMessage(r.Context(), sessionID, req.To, req.Name, req.Options, req.SelectableOptionsCount)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.To).
			Msg("Failed to send poll")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: "poll_" + generateID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) SendButtons(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendButtonsMessageRequest
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

	if len(req.Buttons) == 0 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "buttons are required")
		return
	}

	if len(req.Buttons) > 3 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "maximum 3 buttons allowed")
		return
	}

	var buttons []input.ButtonInfo
	for _, btn := range req.Buttons {
		buttons = append(buttons, input.ButtonInfo{
			ID:   btn.ID,
			Text: btn.Text,
		})
	}

	err := h.messageService.SendButtonsMessage(r.Context(), sessionID, req.To, req.Text, buttons)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.To).
			Msg("Failed to send buttons message")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: "btn_" + generateID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *MessageHandler) SendList(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendListMessageRequest
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

	if req.Title == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "title is required")
		return
	}

	if len(req.Sections) == 0 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sections are required")
		return
	}

	h.writeError(w, http.StatusNotImplemented, "not_implemented", "list messages not yet implemented")
}

func (h *MessageHandler) SendTemplate(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendTemplateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Template == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", "template is required")
		return
	}

	h.writeError(w, http.StatusNotImplemented, "not_implemented", "template messages not yet implemented")
}

func (h *MessageHandler) SendViewOnce(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SendViewOnceMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	if req.Media == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", "media is required")
		return
	}

	h.writeError(w, http.StatusNotImplemented, "not_implemented", "view once messages not yet implemented")
}
