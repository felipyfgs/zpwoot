package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"

	"github.com/go-chi/chi/v5"
)

const (
	errCodeSessionNotFound = "SESSION_NOT_FOUND"
	errCodeNotConnected    = "NOT_CONNECTED"
	errCodeInvalidJID      = "INVALID_JID"
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

// @Summary      Send a message (generic)
// @Description  Send a message of any type (text, media, location, contact) using a generic endpoint
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        sessionId   path      string                    true  "Session ID"
// @Param        message     body      dto.SendMessageRequest    true  "Message data"
// @Success      200         {object}  dto.SendMessageResponse   "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse         "Invalid request"
// @Failure      404         {object}  dto.ErrorResponse         "Session not found"
// @Failure      412         {object}  dto.ErrorResponse         "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse         "Internal server error"
// @Router       /sessions/{sessionId}/messages [post]
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

	err := h.sendMessageByType(r, sessionID, &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to send message")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: "msg_" + generateID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, response)
}

func (h *MessageHandler) sendMessageByType(r *http.Request, sessionID string, req *dto.SendMessageRequest) error {
	switch req.Type {
	case "text":
		if req.Text == "" {
			return &output.WhatsAppError{Code: "VALIDATION_ERROR", Message: "text is required for text messages"}
		}
		return h.messageService.SendTextMessage(r.Context(), sessionID, req.To, req.Text)

	case "media":
		if req.Media == nil {
			return &output.WhatsAppError{Code: "VALIDATION_ERROR", Message: "media is required for media messages"}
		}
		return h.messageService.SendMediaMessage(r.Context(), sessionID, req.To, req.Media.ToInterfacesMediaData())

	case "location":
		if req.Location == nil {
			return &output.WhatsAppError{Code: "VALIDATION_ERROR", Message: "location is required for location messages"}
		}
		return h.messageService.SendLocationMessage(r.Context(), sessionID, req.To,
			req.Location.Latitude, req.Location.Longitude, req.Location.Name)

	case "contact":
		if req.Contact == nil {
			return &output.WhatsAppError{Code: "VALIDATION_ERROR", Message: "contact is required for contact messages"}
		}
		contactInfo := &input.ContactInfo{
			Name:  req.Contact.Name,
			Phone: req.Contact.Phone,
		}
		return h.messageService.SendContactMessage(r.Context(), sessionID, req.To, contactInfo)

	default:
		return &output.WhatsAppError{Code: "VALIDATION_ERROR", Message: "unsupported message type"}
	}
}

func (h *MessageHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

func (h *MessageHandler) writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	}); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode error response")
	}
}

// @Summary      Send text message
// @Description  Send a simple text message to a WhatsApp contact
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        sessionId   path      string                        true  "Session ID"
// @Param        message     body      dto.SendTextMessageRequest    true  "Text message data"
// @Success      200         {object}  dto.SendMessageResponse       "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse             "Invalid request"
// @Failure      404         {object}  dto.ErrorResponse             "Session not found"
// @Failure      412         {object}  dto.ErrorResponse             "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse             "Internal server error"
// @Router       /sessions/{sessionId}/send/message/text [post]
func (h *MessageHandler) SendText(w http.ResponseWriter, r *http.Request) {
	h.SendTextMessage(w, r)
}

func (h *MessageHandler) sendMediaMessageGeneric(w http.ResponseWriter, r *http.Request, mediaType, fieldName string) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req struct {
		Media    *dto.MediaData `json:"media,omitempty"`
		Audio    *dto.MediaData `json:"audio,omitempty"`
		Image    *dto.MediaData `json:"image,omitempty"`
		Video    *dto.MediaData `json:"video,omitempty"`
		Document *dto.MediaData `json:"document,omitempty"`
		Sticker  *dto.MediaData `json:"sticker,omitempty"`
		To       string         `json:"to"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.To == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "to is required")
		return
	}

	var media *dto.MediaData
	switch fieldName {
	case "audio":
		media = req.Audio
	case "image":
		media = req.Image
	case "video":
		media = req.Video
	case "document":
		media = req.Document
	case "sticker":
		media = req.Sticker
	default:
		media = req.Media
	}

	if media == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", fieldName+" is required")
		return
	}

	err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.To, media.ToMediaData())
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.To).
			Str("media_type", mediaType).
			Msg("Failed to send media message")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: "msg_" + generateID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, response)
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

	h.writeJSON(w, response)
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
	h.sendMediaMessageGeneric(w, r, "audio", "audio")
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
	h.sendMediaMessageGeneric(w, r, "image", "image")
}

func (h *MessageHandler) handleMessageError(w http.ResponseWriter, err error) {
	var waErr *output.WhatsAppError
	if errors.As(err, &waErr) {
		var status int
		switch waErr.Code {
		case errCodeSessionNotFound:
			status = http.StatusNotFound
		case errCodeNotConnected:
			status = http.StatusPreconditionFailed
		case errCodeInvalidJID:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		h.writeError(w, status, waErr.Code, waErr.Message)
	} else {
		h.writeError(w, http.StatusInternalServerError, "send_error", "failed to send message")
	}
}

func generateID() string {
	return "123456789"
}

// SendVideo godoc
// @Summary      Send video message
// @Description  Send a video message to a WhatsApp contact with optional caption
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        sessionId   path      string                        true  "Session ID"
// @Param        message     body      dto.SendVideoMessageRequest   true  "Video message data"
// @Success      200         {object}  dto.SendMessageResponse       "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse             "Invalid request"
// @Failure      404         {object}  dto.ErrorResponse             "Session not found"
// @Failure      412         {object}  dto.ErrorResponse             "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse             "Internal server error"
// @Router       /sessions/{sessionId}/send/message/video [post]
func (h *MessageHandler) SendVideo(w http.ResponseWriter, r *http.Request) {
	h.sendMediaMessageGeneric(w, r, "video", "video")
}

// SendDocument godoc
// @Summary      Send document message
// @Description  Send a document/file message to a WhatsApp contact
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        sessionId   path      string                           true  "Session ID"
// @Param        message     body      dto.SendDocumentMessageRequest   true  "Document message data"
// @Success      200         {object}  dto.SendMessageResponse          "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse                "Invalid request"
// @Failure      404         {object}  dto.ErrorResponse                "Session not found"
// @Failure      412         {object}  dto.ErrorResponse                "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse                "Internal server error"
// @Router       /sessions/{sessionId}/send/message/document [post]
func (h *MessageHandler) SendDocument(w http.ResponseWriter, r *http.Request) {
	h.sendMediaMessageGeneric(w, r, "document", "document")
}

// SendSticker godoc
// @Summary      Send sticker message
// @Description  Send a sticker message to a WhatsApp contact
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        sessionId   path      string                          true  "Session ID"
// @Param        message     body      dto.SendStickerMessageRequest   true  "Sticker message data"
// @Success      200         {object}  dto.SendMessageResponse         "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse               "Invalid request"
// @Failure      404         {object}  dto.ErrorResponse               "Session not found"
// @Failure      412         {object}  dto.ErrorResponse               "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse               "Internal server error"
// @Router       /sessions/{sessionId}/send/message/sticker [post]
func (h *MessageHandler) SendSticker(w http.ResponseWriter, r *http.Request) {
	h.sendMediaMessageGeneric(w, r, "sticker", "sticker")
}

// SendAudio godoc
// @Summary      Send audio message
// @Description  Send an audio/voice message to a WhatsApp contact
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        sessionId   path      string                        true  "Session ID"
// @Param        message     body      dto.SendAudioMessageRequest   true  "Audio message data"
// @Success      200         {object}  dto.SendMessageResponse       "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse             "Invalid request"
// @Failure      404         {object}  dto.ErrorResponse             "Session not found"
// @Failure      412         {object}  dto.ErrorResponse             "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse             "Internal server error"
// @Router       /sessions/{sessionId}/send/message/audio [post]
func (h *MessageHandler) SendAudio(w http.ResponseWriter, r *http.Request) {
	h.SendAudioMessage(w, r)
}

// SendImage godoc
// @Summary      Send image message
// @Description  Send an image message to a WhatsApp contact with optional caption
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        sessionId   path      string                        true  "Session ID"
// @Param        message     body      dto.SendImageMessageRequest   true  "Image message data"
// @Success      200         {object}  dto.SendMessageResponse       "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse             "Invalid request"
// @Failure      404         {object}  dto.ErrorResponse             "Session not found"
// @Failure      412         {object}  dto.ErrorResponse             "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse             "Internal server error"
// @Router       /sessions/{sessionId}/send/message/image [post]
func (h *MessageHandler) SendImage(w http.ResponseWriter, r *http.Request) {
	h.SendImageMessage(w, r)
}

// SendLocation godoc
// @Summary      Send location message
// @Description  Send a location message with GPS coordinates to a WhatsApp contact
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        sessionId   path      string                           true  "Session ID"
// @Param        message     body      dto.SendLocationMessageRequest   true  "Location message data"
// @Success      200         {object}  dto.SendMessageResponse          "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse                "Invalid request"
// @Failure      404         {object}  dto.ErrorResponse                "Session not found"
// @Failure      412         {object}  dto.ErrorResponse                "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse                "Internal server error"
// @Router       /sessions/{sessionId}/send/message/location [post]
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

	h.writeJSON(w, response)
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

	h.writeJSON(w, response)
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

	h.writeJSON(w, response)
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

	h.writeJSON(w, response)
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

	h.writeJSON(w, response)
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

	buttons := make([]input.ButtonInfo, 0, len(req.Buttons))
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

	h.writeJSON(w, response)
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

	sections := make([]input.ListSectionInfo, 0, len(req.Sections))
	for _, section := range req.Sections {
		rows := make([]input.ListRowInfo, 0, len(section.Rows))
		for _, row := range section.Rows {
			rows = append(rows, input.ListRowInfo{
				ID:          row.ID,
				Title:       row.Title,
				Description: row.Description,
			})
		}
		sections = append(sections, input.ListSectionInfo{
			Title: section.Title,
			Rows:  rows,
		})
	}

	err := h.messageService.SendListMessage(r.Context(), sessionID, req.To, req.Text, req.Title, sections)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.To).
			Msg("Failed to send list message")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: "list_" + generateID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, response)
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

	template := input.TemplateInfo{
		Content: req.Template.Name,
		Footer:  req.Template.Language,
	}

	err := h.messageService.SendTemplateMessage(r.Context(), sessionID, req.To, template)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.To).
			Msg("Failed to send template message")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: "tmpl_" + generateID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, response)
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

	media := req.Media.ToInterfacesMediaData()

	err := h.messageService.SendViewOnceMessage(r.Context(), sessionID, req.To, media)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.To).
			Msg("Failed to send view once message")
		h.handleMessageError(w, err)
		return
	}

	response := &dto.SendMessageResponse{
		MessageID: "vo_" + generateID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}

	h.writeJSON(w, response)
}
