package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/application/utils"
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

// buildMessageResponse creates a standardized message response
func (h *MessageHandler) buildMessageResponse(result *output.MessageResult, to, messageType, content string) *dto.SendMessageResponse {
	response := &dto.SendMessageResponse{
		Success:   true,
		ID:        result.MessageID,
		To:        to,
		Type:      messageType,
		Content:   content,
		Timestamp: result.SentAt.Unix(),
		Status:    result.Status,
	}

	if result.ContextInfo != nil {
		response.ContextInfo = &dto.ContextInfo{
			StanzaID:    result.ContextInfo.StanzaID,
			Participant: result.ContextInfo.Participant,
			QuotedID:    result.ContextInfo.QuotedID,
		}
	}

	return response
}

// @Summary      Send text message
// @Description  Send a simple text message to a WhatsApp contact. Supports reply/quote using contextInfo.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                        true  "Session ID"
// @Param        message     body      dto.SendTextMessageRequest    true  "Text message data"
// @Success      200         {object}  dto.SendMessageResponse       "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse             "Invalid request"
// @Failure      401         {object}  dto.ErrorResponse             "Unauthorized"
// @Failure      404         {object}  dto.ErrorResponse             "Session not found"
// @Failure      412         {object}  dto.ErrorResponse             "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse             "Internal server error"
// @Router       /sessions/{sessionId}/send/message/text [post]
func (h *MessageHandler) SendText(w http.ResponseWriter, r *http.Request) {
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	if req.Text == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "text is required")
		return
	}

	// Convert DTO ContextInfo to output ContextInfo
	var contextInfo *output.MessageContextInfo
	if req.ContextInfo != nil {
		contextInfo = &output.MessageContextInfo{
			StanzaID:    req.ContextInfo.StanzaID,
			Participant: req.ContextInfo.Participant,
		}
	}

	result, err := h.messageService.SendTextMessage(r.Context(), sessionID, req.Phone, req.Text, contextInfo)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("phone", req.Phone).
			Msg("Failed to send text message")
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "text", req.Text)
	h.writeJSON(w, response)
}

// @Summary      Send image message
// @Description  Send an image message to a WhatsApp contact with optional caption. Supports Base64, URL, or file path. Supports reply/quote using contextInfo. Set viewOnce to true to send as a view-once message that disappears after being viewed.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                        true  "Session ID"
// @Param        message     body      dto.SendImageMessageRequest   true  "Image message data"
// @Success      200         {object}  dto.SendMessageResponse       "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse             "Invalid request or media processing error"
// @Failure      401         {object}  dto.ErrorResponse             "Unauthorized"
// @Failure      404         {object}  dto.ErrorResponse             "Session not found"
// @Failure      412         {object}  dto.ErrorResponse             "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse             "Internal server error"
// @Router       /sessions/{sessionId}/send/message/image [post]
func (h *MessageHandler) SendImage(w http.ResponseWriter, r *http.Request) {
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	if req.File == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "file is required")
		return
	}

	// Convert DTO ContextInfo to output ContextInfo
	var contextInfo *output.MessageContextInfo
	if req.ContextInfo != nil {
		contextInfo = &output.MessageContextInfo{
			StanzaID:    req.ContextInfo.StanzaID,
			Participant: req.ContextInfo.Participant,
		}
	}

	// Process media using the new media processor
	mediaProcessor := utils.NewMediaProcessor()
	media, err := mediaProcessor.ProcessMedia(req.File, req.MimeType, req.FileName)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("file", req.File).
			Msg("Failed to process media")
		h.writeError(w, http.StatusBadRequest, "media_processing_error", fmt.Sprintf("Failed to process media: %v", err))
		return
	}

	media.Caption = req.Caption
	media.ViewOnce = req.ViewOnce

	result, err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.Phone, media, contextInfo)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "image", req.Caption)
	h.writeJSON(w, response)
}

// @Summary      Send audio message
// @Description  Send an audio message to a WhatsApp contact. Supports Base64, URL, or file path. Supports reply/quote using contextInfo. Set viewOnce to true to send as a view-once message that disappears after being viewed.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                        true  "Session ID"
// @Param        message     body      dto.SendAudioMessageRequest   true  "Audio message data"
// @Success      200         {object}  dto.SendMessageResponse       "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse             "Invalid request or media processing error"
// @Failure      401         {object}  dto.ErrorResponse             "Unauthorized"
// @Failure      404         {object}  dto.ErrorResponse             "Session not found"
// @Failure      412         {object}  dto.ErrorResponse             "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse             "Internal server error"
// @Router       /sessions/{sessionId}/send/message/audio [post]
func (h *MessageHandler) SendAudio(w http.ResponseWriter, r *http.Request) {
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	if req.File == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "file is required")
		return
	}

	// Convert DTO ContextInfo to output ContextInfo
	var contextInfo *output.MessageContextInfo
	if req.ContextInfo != nil {
		contextInfo = &output.MessageContextInfo{
			StanzaID:    req.ContextInfo.StanzaID,
			Participant: req.ContextInfo.Participant,
		}
	}

	// Process media using the new media processor
	mediaProcessor := utils.NewMediaProcessor()
	media, err := mediaProcessor.ProcessMedia(req.File, req.MimeType, req.FileName)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("file", req.File).
			Msg("Failed to process media")
		h.writeError(w, http.StatusBadRequest, "media_processing_error", fmt.Sprintf("Failed to process media: %v", err))
		return
	}

	media.ViewOnce = req.ViewOnce

	result, err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.Phone, media, contextInfo)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "audio", "")
	h.writeJSON(w, response)
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

func (h *MessageHandler) handleMessageError(w http.ResponseWriter, err error) {
	var waErr *output.WhatsAppError
	if errors.As(err, &waErr) {
		switch waErr.Code {
		case "SESSION_NOT_FOUND":
			h.writeError(w, http.StatusNotFound, "session_not_found", "Session not found")
		case "NOT_CONNECTED":
			h.writeError(w, http.StatusPreconditionFailed, "not_connected", "Session not connected")
		case "INVALID_JID":
			h.writeError(w, http.StatusBadRequest, "invalid_jid", "Invalid recipient JID")
		default:
			h.writeError(w, http.StatusInternalServerError, "whatsapp_error", waErr.Message)
		}
		return
	}
	h.writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
}

// @Summary      Send video message
// @Description  Send a video message to a WhatsApp contact with optional caption. Supports Base64, URL, or file path. Supports reply/quote using contextInfo. Set viewOnce to true to send as a view-once message that disappears after being viewed.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                        true  "Session ID"
// @Param        message     body      dto.SendVideoMessageRequest   true  "Video message data"
// @Success      200         {object}  dto.SendMessageResponse       "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse             "Invalid request or media processing error"
// @Failure      401         {object}  dto.ErrorResponse             "Unauthorized"
// @Failure      404         {object}  dto.ErrorResponse             "Session not found"
// @Failure      412         {object}  dto.ErrorResponse             "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse             "Internal server error"
// @Router       /sessions/{sessionId}/send/message/video [post]
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	if req.File == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "file is required")
		return
	}

	var contextInfo *output.MessageContextInfo
	if req.ContextInfo != nil {
		contextInfo = &output.MessageContextInfo{
			StanzaID:    req.ContextInfo.StanzaID,
			Participant: req.ContextInfo.Participant,
		}
	}

	// Process media using the new media processor
	mediaProcessor := utils.NewMediaProcessor()
	media, err := mediaProcessor.ProcessMedia(req.File, req.MimeType, req.FileName)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("file", req.File).
			Msg("Failed to process media")
		h.writeError(w, http.StatusBadRequest, "media_processing_error", fmt.Sprintf("Failed to process media: %v", err))
		return
	}

	media.Caption = req.Caption
	media.ViewOnce = req.ViewOnce

	result, err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.Phone, media, contextInfo)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "video", req.Caption)
	h.writeJSON(w, response)
}

// @Summary      Send document message
// @Description  Send a document message to a WhatsApp contact. Supports Base64, URL, or file path. Supports reply/quote using contextInfo.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                           true  "Session ID"
// @Param        message     body      dto.SendDocumentMessageRequest   true  "Document message data"
// @Success      200         {object}  dto.SendMessageResponse          "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse                "Invalid request or media processing error"
// @Failure      401         {object}  dto.ErrorResponse                "Unauthorized"
// @Failure      404         {object}  dto.ErrorResponse                "Session not found"
// @Failure      412         {object}  dto.ErrorResponse                "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse                "Internal server error"
// @Router       /sessions/{sessionId}/send/message/document [post]
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	if req.File == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "file is required")
		return
	}

	var contextInfo *output.MessageContextInfo
	if req.ContextInfo != nil {
		contextInfo = &output.MessageContextInfo{
			StanzaID:    req.ContextInfo.StanzaID,
			Participant: req.ContextInfo.Participant,
		}
	}

	// Process media using the new media processor
	mediaProcessor := utils.NewMediaProcessor()
	media, err := mediaProcessor.ProcessMedia(req.File, req.MimeType, req.FileName)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("file", req.File).
			Msg("Failed to process media")
		h.writeError(w, http.StatusBadRequest, "media_processing_error", fmt.Sprintf("Failed to process media: %v", err))
		return
	}

	media.Caption = req.Caption

	result, err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.Phone, media, contextInfo)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "document", req.Caption)
	h.writeJSON(w, response)
}

// @Summary      Send location message
// @Description  Send a location message to a WhatsApp contact with latitude, longitude, and optional name/address. Supports reply/quote using contextInfo.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                           true  "Session ID"
// @Param        message     body      dto.SendLocationMessageRequest   true  "Location message data"
// @Success      200         {object}  dto.SendMessageResponse          "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse                "Invalid request"
// @Failure      401         {object}  dto.ErrorResponse                "Unauthorized"
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	var contextInfo *output.MessageContextInfo
	if req.ContextInfo != nil {
		contextInfo = &output.MessageContextInfo{
			StanzaID:    req.ContextInfo.StanzaID,
			Participant: req.ContextInfo.Participant,
		}
	}

	result, err := h.messageService.SendLocationMessage(r.Context(), sessionID, req.Phone, req.Latitude, req.Longitude, req.Name, contextInfo)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "location", req.Name)
	h.writeJSON(w, response)
}

// @Summary      Send contact message
// @Description  Send a single contact card to a WhatsApp contact. VCard is auto-generated if not provided. Supports reply/quote using contextInfo.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                          true  "Session ID"
// @Param        message     body      dto.SendContactMessageRequest   true  "Contact message data"
// @Success      200         {object}  dto.SendMessageResponse         "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse               "Invalid request"
// @Failure      401         {object}  dto.ErrorResponse               "Unauthorized"
// @Failure      404         {object}  dto.ErrorResponse               "Session not found"
// @Failure      412         {object}  dto.ErrorResponse               "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse               "Internal server error"
// @Router       /sessions/{sessionId}/send/message/contact [post]
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	if req.Contact == nil {
		h.writeError(w, http.StatusBadRequest, "validation_error", "contact is required")
		return
	}

	var contextInfo *output.MessageContextInfo
	if req.ContextInfo != nil {
		contextInfo = &output.MessageContextInfo{
			StanzaID:    req.ContextInfo.StanzaID,
			Participant: req.ContextInfo.Participant,
		}
	}

	contactInfo := &input.ContactInfo{
		Name:  req.Contact.Name,
		Phone: req.Contact.Phone,
		VCard: req.Contact.VCard,
	}

	result, err := h.messageService.SendContactMessage(r.Context(), sessionID, req.Phone, contactInfo, contextInfo)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "contact", contactInfo.Name)
	h.writeJSON(w, response)
}

// @Summary      Send reaction message
// @Description  Send a reaction (emoji) to a specific message. Use empty string to remove reaction.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                           true  "Session ID"
// @Param        message     body      dto.SendReactionMessageRequest   true  "Reaction message data"
// @Success      200         {object}  dto.SendMessageResponse          "Reaction sent successfully"
// @Failure      400         {object}  dto.ErrorResponse                "Invalid request"
// @Failure      401         {object}  dto.ErrorResponse                "Unauthorized"
// @Failure      404         {object}  dto.ErrorResponse                "Session not found"
// @Failure      412         {object}  dto.ErrorResponse                "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse                "Internal server error"
// @Router       /sessions/{sessionId}/send/message/reaction [post]
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	if req.MessageID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "messageId is required")
		return
	}

	// Process "me:" prefix in messageID (WuzAPI compatibility)
	// Example: "me:3EB0C767D0D1A6F4FD29" means fromMe=true
	messageID := req.MessageID
	fromMe := false

	if strings.HasPrefix(messageID, "me:") {
		fromMe = true
		messageID = messageID[len("me:"):]
	}

	// Allow explicit fromMe field to override prefix detection
	if req.FromMe != nil {
		fromMe = *req.FromMe
	}

	result, err := h.messageService.SendReactionMessage(r.Context(), sessionID, req.Phone, messageID, req.Reaction, fromMe)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("to", req.Phone).
			Str("message_id", req.MessageID).
			Msg("Failed to send reaction")
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "reaction", req.Reaction)
	h.writeJSON(w, response)
}

// @Summary      Send poll message
// @Description  Send a poll message to a WhatsApp contact with multiple options. Supports single or multiple selection.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                       true  "Session ID"
// @Param        message     body      dto.SendPollMessageRequest   true  "Poll message data"
// @Success      200         {object}  dto.SendMessageResponse      "Poll sent successfully"
// @Failure      400         {object}  dto.ErrorResponse            "Invalid request"
// @Failure      401         {object}  dto.ErrorResponse            "Unauthorized"
// @Failure      404         {object}  dto.ErrorResponse            "Session not found"
// @Failure      412         {object}  dto.ErrorResponse            "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse            "Internal server error"
// @Router       /sessions/{sessionId}/send/message/poll [post]
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "name is required")
		return
	}

	if len(req.Options) == 0 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "options are required")
		return
	}

	result, err := h.messageService.SendPollMessage(r.Context(), sessionID, req.Phone, req.Name, req.Options, req.SelectableOptionsCount)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "poll", req.Name)
	h.writeJSON(w, response)
}

// @Summary      Send sticker message
// @Description  Send a sticker message to a WhatsApp contact. Supports Base64, URL, or file path. Image will be converted to WebP format.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                          true  "Session ID"
// @Param        message     body      dto.SendStickerMessageRequest   true  "Sticker message data"
// @Success      200         {object}  dto.SendMessageResponse         "Sticker sent successfully"
// @Failure      400         {object}  dto.ErrorResponse               "Invalid request or media processing error"
// @Failure      401         {object}  dto.ErrorResponse               "Unauthorized"
// @Failure      404         {object}  dto.ErrorResponse               "Session not found"
// @Failure      412         {object}  dto.ErrorResponse               "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse               "Internal server error"
// @Router       /sessions/{sessionId}/send/message/sticker [post]
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	if req.File == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "file is required")
		return
	}

	// Convert DTO ContextInfo to output ContextInfo
	var contextInfo *output.MessageContextInfo
	if req.ContextInfo != nil {
		contextInfo = &output.MessageContextInfo{
			StanzaID:    req.ContextInfo.StanzaID,
			Participant: req.ContextInfo.Participant,
		}
	}

	// Process media using the new media processor
	mediaProcessor := utils.NewMediaProcessor()
	media, err := mediaProcessor.ProcessMedia(req.File, req.MimeType, req.FileName)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("file", req.File).
			Msg("Failed to process media")
		h.writeError(w, http.StatusBadRequest, "media_processing_error", fmt.Sprintf("Failed to process media: %v", err))
		return
	}

	result, err := h.messageService.SendMediaMessage(r.Context(), sessionID, req.Phone, media, contextInfo)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "sticker", "")
	h.writeJSON(w, response)
}

// @Summary      Send multiple contacts
// @Description  Send multiple contacts in a single message using ContactsArrayMessage. All contacts are sent in one message, not as separate messages. VCard is auto-generated for each contact if not provided.
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                                true  "Session ID"
// @Param        message     body      dto.SendContactsArrayMessageRequest   true  "Contacts array message data"
// @Success      200         {object}  dto.SendMessageResponse               "Contacts sent successfully"
// @Failure      400         {object}  dto.ErrorResponse                     "Invalid request"
// @Failure      401         {object}  dto.ErrorResponse                     "Unauthorized"
// @Failure      404         {object}  dto.ErrorResponse                     "Session not found"
// @Failure      412         {object}  dto.ErrorResponse                     "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse                     "Internal server error"
// @Router       /sessions/{sessionId}/send/message/contacts [post]
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

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "phone is required")
		return
	}

	if len(req.Contacts) == 0 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "contacts array is required and must not be empty")
		return
	}

	// Convert contacts to input format
	contacts := make([]*input.ContactInfo, len(req.Contacts))
	for i, contact := range req.Contacts {
		contacts[i] = &input.ContactInfo{
			Name:  contact.Name,
			Phone: contact.Phone,
			VCard: contact.VCard,
		}
	}

	result, err := h.messageService.SendContactsArrayMessage(r.Context(), sessionID, req.Phone, contacts)
	if err != nil {
		h.handleMessageError(w, err)
		return
	}

	normalizedTo := req.Phone
	if !strings.Contains(normalizedTo, "@") {
		normalizedTo = normalizedTo + "@s.whatsapp.net"
	}

	response := h.buildMessageResponse(result, normalizedTo, "contacts", fmt.Sprintf("%d contacts", len(contacts)))
	h.writeJSON(w, response)
}

// SendTemplate sends a template message
func (h *MessageHandler) SendTemplate(w http.ResponseWriter, r *http.Request) {
	h.writeError(w, http.StatusNotImplemented, "not_implemented", "Template messages not yet implemented")
}

// SendButtons sends a buttons message
func (h *MessageHandler) SendButtons(w http.ResponseWriter, r *http.Request) {
	h.writeError(w, http.StatusNotImplemented, "not_implemented", "Button messages not yet implemented")
}

// SendList sends a list message
func (h *MessageHandler) SendList(w http.ResponseWriter, r *http.Request) {
	h.writeError(w, http.StatusNotImplemented, "not_implemented", "List messages not yet implemented")
}


