package handlers

import (
	"encoding/json"
	"net/http"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"

	"github.com/go-chi/chi/v5"
)

type ContactHandler struct {
	contactService input.ContactService
	logger         output.Logger
}

func NewContactHandler(contactService input.ContactService, logger output.Logger) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
		logger:         logger,
	}
}

// CheckUser godoc
// @Summary      Check if users are on WhatsApp
// @Description  Verify if phone numbers are registered on WhatsApp
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                    true  "Session ID"
// @Param        request     body      dto.CheckUserRequest      true  "Phone numbers to check"
// @Success      200  {object}  dto.CheckUserResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/contacts/check [post]
func (h *ContactHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	var req dto.CheckUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeBadRequest, "Invalid JSON body")
		return
	}

	if len(req.Phones) == 0 {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "phones is required")
		return
	}

	results, err := h.contactService.CheckUser(r.Context(), sessionID, req.Phones)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to check users")

		h.writeError(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "Failed to check users")
		return
	}

	users := make([]dto.WhatsAppUserInfo, 0, len(results))
	for _, result := range results {
		users = append(users, dto.WhatsAppUserInfo{
			Query:        result.Query,
			IsInWhatsApp: result.IsInWhatsApp,
			JID:          result.JID,
			VerifiedName: result.VerifiedName,
		})
	}

	response := &dto.CheckUserResponse{
		Users: users,
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// GetUser godoc
// @Summary      Get user information
// @Description  Get detailed information about a WhatsApp user
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                 true  "Session ID"
// @Param        request     body      dto.GetUserRequest     true  "Phone number"
// @Success      200  {object}  dto.GetUserResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/contacts/user [post]
func (h *ContactHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	var req dto.GetUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeBadRequest, "Invalid JSON body")
		return
	}

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "phone is required")
		return
	}

	userDetail, err := h.contactService.GetUser(r.Context(), sessionID, req.Phone)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("phone", req.Phone).
			Msg("Failed to get user")

		h.writeError(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "Failed to get user")
		return
	}

	response := &dto.GetUserResponse{
		JID:          userDetail.JID,
		VerifiedName: userDetail.VerifiedName,
		Status:       userDetail.Status,
		PictureID:    userDetail.PictureID,
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// GetAvatar godoc
// @Summary      Get user avatar
// @Description  Get profile picture of a WhatsApp user
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                  true  "Session ID"
// @Param        request     body      dto.GetAvatarRequest    true  "Avatar request"
// @Success      200  {object}  dto.GetAvatarResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/contacts/avatar [post]
func (h *ContactHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	var req dto.GetAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeBadRequest, "Invalid JSON body")
		return
	}

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "phone is required")
		return
	}

	avatarInfo, err := h.contactService.GetAvatar(r.Context(), sessionID, req.Phone, req.Preview)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("phone", req.Phone).
			Msg("Failed to get avatar")

		h.writeError(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "Failed to get avatar")
		return
	}

	response := &dto.GetAvatarResponse{
		URL:       avatarInfo.URL,
		ID:        avatarInfo.ID,
		Type:      avatarInfo.Type,
		DirectURL: avatarInfo.DirectURL,
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// GetContacts godoc
// @Summary      Get all contacts
// @Description  List all contacts from the WhatsApp account
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string    true  "Session ID"
// @Success      200  {object}  dto.GetContactsResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/contacts [get]
func (h *ContactHandler) GetContacts(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	contacts, err := h.contactService.GetContacts(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to get contacts")

		h.writeError(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "Failed to get contacts")
		return
	}

	contactDetails := make([]dto.ContactDetail, 0, len(contacts))
	for _, contact := range contacts {
		contactDetails = append(contactDetails, dto.ContactDetail{
			JID:          contact.JID,
			Name:         contact.Name,
			Notify:       contact.Notify,
			VerifiedName: contact.VerifiedName,
			BusinessName: contact.BusinessName,
		})
	}

	response := &dto.GetContactsResponse{
		Contacts: contactDetails,
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// SendPresence godoc
// @Summary      Send presence
// @Description  Send general presence status (available/unavailable)
// @Tags         Presence
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                      true  "Session ID"
// @Param        request     body      dto.SendPresenceRequest     true  "Presence request"
// @Success      200  {object}  dto.SendPresenceResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/presence/send [post]
func (h *ContactHandler) SendPresence(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	var req dto.SendPresenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeBadRequest, "Invalid JSON body")
		return
	}

	if req.Presence == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "presence is required")
		return
	}

	err := h.contactService.SendPresence(r.Context(), sessionID, req.Presence)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("presence", req.Presence).
			Msg("Failed to send presence")

		h.writeError(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "Failed to send presence")
		return
	}

	response := &dto.SendPresenceResponse{
		Success: true,
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// ChatPresence godoc
// @Summary      Send chat presence
// @Description  Send presence in specific chat (typing/recording)
// @Tags         Presence
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                     true  "Session ID"
// @Param        request     body      dto.ChatPresenceRequest    true  "Chat presence request"
// @Success      200  {object}  dto.ChatPresenceResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/presence/chat [post]
func (h *ContactHandler) ChatPresence(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "sessionId is required")
		return
	}

	var req dto.ChatPresenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeBadRequest, "Invalid JSON body")
		return
	}

	if req.Phone == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "phone is required")
		return
	}

	if req.Presence == "" {
		h.writeError(w, http.StatusBadRequest, dto.ErrorCodeValidation, "presence is required")
		return
	}

	err := h.contactService.ChatPresence(r.Context(), sessionID, req.Phone, req.Presence, req.Media)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("phone", req.Phone).
			Str("presence", req.Presence).
			Msg("Failed to send chat presence")

		h.writeError(w, http.StatusInternalServerError, dto.ErrorCodeInternalError, "Failed to send chat presence")
		return
	}

	response := &dto.ChatPresenceResponse{
		Success: true,
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

func (h *ContactHandler) writeSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := dto.NewSuccessResponse(data)
	json.NewEncoder(w).Encode(response)
}

func (h *ContactHandler) writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	})
}
