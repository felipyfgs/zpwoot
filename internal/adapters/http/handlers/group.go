package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"
)

// GroupHandler gerencia requisições HTTP relacionadas a grupos
type GroupHandler struct {
	groupService input.GroupService
	logger       *logger.Logger
}

// NewGroupHandler cria uma nova instância do GroupHandler
func NewGroupHandler(groupService input.GroupService, logger *logger.Logger) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
		logger:       logger,
	}
}

// ListGroups godoc
// @Summary      List groups
// @Description  List all groups the session is part of
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string  true  "Session ID"
// @Success      200  {object}  dto.ListGroupsResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      412  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups [get]
func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	groups, err := h.groupService.ListGroups(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to list groups")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Int("count", len(groups.Groups)).
		Msg("Groups listed successfully")

	h.writeJSON(w, groups)
}

// GetGroupInfo godoc
// @Summary      Get group info
// @Description  Get detailed information about a specific group
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string  true  "Session ID"
// @Param        groupJid    query     string  true  "Group JID"
// @Success      200  {object}  dto.GroupInfo
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/info [get]
func (h *GroupHandler) GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	groupJID := r.URL.Query().Get("groupJid")
	if groupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	group, err := h.groupService.GetGroupInfo(r.Context(), sessionID, groupJID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", groupJID).
			Msg("Failed to get group info")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", groupJID).
		Msg("Group info retrieved successfully")

	h.writeJSON(w, group)
}

// GetGroupInviteInfo godoc
// @Summary      Get group invite info
// @Description  Get group information from invite code without joining
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                           true  "Session ID"
// @Param        request     body      dto.GetGroupInviteInfoRequest    true  "Invite code"
// @Success      200  {object}  dto.GroupInfo
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/invite-info [post]
func (h *GroupHandler) GetGroupInviteInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.GetGroupInviteInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.Code == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "code is required")
		return
	}

	group, err := h.groupService.GetGroupInviteInfo(r.Context(), sessionID, req.Code)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("code", req.Code).
			Msg("Failed to get group invite info")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", group.JID).
		Msg("Group invite info retrieved successfully")

	h.writeJSON(w, group)
}

// GetGroupInviteLink godoc
// @Summary      Get group invite link
// @Description  Get or reset the invite link for a group
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string  true   "Session ID"
// @Param        groupJid    query     string  true   "Group JID"
// @Param        reset       query     bool    false  "Reset link"
// @Success      200  {object}  dto.GetGroupInviteLinkResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/invite-link [get]
func (h *GroupHandler) GetGroupInviteLink(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	groupJID := r.URL.Query().Get("groupJid")
	if groupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	reset := r.URL.Query().Get("reset") == "true"

	link, err := h.groupService.GetGroupInviteLink(r.Context(), sessionID, groupJID, reset)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", groupJID).
			Bool("reset", reset).
			Msg("Failed to get group invite link")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", groupJID).
		Bool("reset", reset).
		Msg("Group invite link retrieved successfully")

	h.writeJSON(w, dto.GetGroupInviteLinkResponse{
		InviteLink: link,
	})
}

// JoinGroup godoc
// @Summary      Join group
// @Description  Join a group using invite code
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                  true  "Session ID"
// @Param        request     body      dto.JoinGroupRequest    true  "Join request"
// @Success      200  {object}  dto.GroupActionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/join [post]
func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.Code == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "code is required")
		return
	}

	err := h.groupService.JoinGroup(r.Context(), sessionID, req.Code)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("code", req.Code).
			Msg("Failed to join group")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("code", req.Code).
		Msg("Joined group successfully")

	h.writeJSON(w, dto.GroupActionResponse{
		Success: true,
		Message: "Group joined successfully",
	})
}

// CreateGroup godoc
// @Summary      Create group
// @Description  Create a new WhatsApp group
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                     true  "Session ID"
// @Param        request     body      dto.CreateGroupRequest     true  "Group data"
// @Success      200  {object}  dto.GroupInfo
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/create [post]
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "name is required")
		return
	}

	if len(req.Participants) < 1 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "at least one participant is required")
		return
	}

	group, err := h.groupService.CreateGroup(r.Context(), sessionID, req.Name, req.Participants)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_name", req.Name).
			Msg("Failed to create group")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", group.JID).
		Str("group_name", group.Name).
		Msg("Group created successfully")

	h.writeJSON(w, group)
}

// LeaveGroup godoc
// @Summary      Leave group
// @Description  Leave a WhatsApp group
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                   true  "Session ID"
// @Param        request     body      dto.LeaveGroupRequest    true  "Leave request"
// @Success      200  {object}  dto.GroupActionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/leave [post]
func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.LeaveGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.GroupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	err := h.groupService.LeaveGroup(r.Context(), sessionID, req.GroupJID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", req.GroupJID).
			Msg("Failed to leave group")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", req.GroupJID).
		Msg("Left group successfully")

	h.writeJSON(w, dto.GroupActionResponse{
		Success: true,
		Message: "Group left successfully",
	})
}

// UpdateGroupParticipants godoc
// @Summary      Update group participants
// @Description  Add, remove, promote or demote group participants
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                                true  "Session ID"
// @Param        request     body      dto.UpdateGroupParticipantsRequest    true  "Participants update"
// @Success      200  {object}  dto.GroupActionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/participants [post]
func (h *GroupHandler) UpdateGroupParticipants(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.UpdateGroupParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.GroupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	if len(req.Participants) < 1 {
		h.writeError(w, http.StatusBadRequest, "validation_error", "at least one participant is required")
		return
	}

	if req.Action == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "action is required")
		return
	}

	err := h.groupService.UpdateGroupParticipants(r.Context(), sessionID, req.GroupJID, req.Participants, req.Action)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", req.GroupJID).
			Str("action", req.Action).
			Msg("Failed to update group participants")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", req.GroupJID).
		Str("action", req.Action).
		Int("count", len(req.Participants)).
		Msg("Group participants updated successfully")

	h.writeJSON(w, dto.GroupActionResponse{
		Success: true,
		Message: fmt.Sprintf("Group participants %sed successfully", req.Action),
	})
}

// SetGroupName godoc
// @Summary      Set group name
// @Description  Change the name of a group
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                     true  "Session ID"
// @Param        request     body      dto.SetGroupNameRequest    true  "Name update"
// @Success      200  {object}  dto.GroupActionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/name [post]
func (h *GroupHandler) SetGroupName(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SetGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.GroupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "name is required")
		return
	}

	err := h.groupService.SetGroupName(r.Context(), sessionID, req.GroupJID, req.Name)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", req.GroupJID).
			Str("name", req.Name).
			Msg("Failed to set group name")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", req.GroupJID).
		Str("name", req.Name).
		Msg("Group name set successfully")

	h.writeJSON(w, dto.GroupActionResponse{
		Success: true,
		Message: "Group name set successfully",
	})
}

// SetGroupTopic godoc
// @Summary      Set group topic
// @Description  Change the description/topic of a group
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                      true  "Session ID"
// @Param        request     body      dto.SetGroupTopicRequest    true  "Topic update"
// @Success      200  {object}  dto.GroupActionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/topic [post]
func (h *GroupHandler) SetGroupTopic(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SetGroupTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.GroupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	if req.Topic == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "topic is required")
		return
	}

	err := h.groupService.SetGroupTopic(r.Context(), sessionID, req.GroupJID, req.Topic)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", req.GroupJID).
			Msg("Failed to set group topic")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", req.GroupJID).
		Msg("Group topic set successfully")

	h.writeJSON(w, dto.GroupActionResponse{
		Success: true,
		Message: "Group topic set successfully",
	})
}

// SetGroupLocked godoc
// @Summary      Set group locked
// @Description  Lock or unlock group settings (only admins can edit when locked)
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                       true  "Session ID"
// @Param        request     body      dto.SetGroupLockedRequest    true  "Locked setting"
// @Success      200  {object}  dto.GroupActionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/settings/locked [post]
func (h *GroupHandler) SetGroupLocked(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SetGroupLockedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.GroupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	err := h.groupService.SetGroupLocked(r.Context(), sessionID, req.GroupJID, req.Locked)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", req.GroupJID).
			Bool("locked", req.Locked).
			Msg("Failed to set group locked")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", req.GroupJID).
		Bool("locked", req.Locked).
		Msg("Group locked setting updated successfully")

	h.writeJSON(w, dto.GroupActionResponse{
		Success: true,
		Message: "Group locked setting updated successfully",
	})
}

// SetGroupAnnounce godoc
// @Summary      Set group announce
// @Description  Enable or disable announce mode (only admins can send messages when enabled)
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                         true  "Session ID"
// @Param        request     body      dto.SetGroupAnnounceRequest    true  "Announce setting"
// @Success      200  {object}  dto.GroupActionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/settings/announce [post]
func (h *GroupHandler) SetGroupAnnounce(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SetGroupAnnounceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.GroupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	err := h.groupService.SetGroupAnnounce(r.Context(), sessionID, req.GroupJID, req.Announce)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", req.GroupJID).
			Bool("announce", req.Announce).
			Msg("Failed to set group announce")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", req.GroupJID).
		Bool("announce", req.Announce).
		Msg("Group announce setting updated successfully")

	h.writeJSON(w, dto.GroupActionResponse{
		Success: true,
		Message: "Group announce setting updated successfully",
	})
}

// SetDisappearingTimer godoc
// @Summary      Set disappearing timer
// @Description  Configure disappearing messages timer (24h, 7d, 90d, or off)
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                             true  "Session ID"
// @Param        request     body      dto.SetDisappearingTimerRequest    true  "Timer setting"
// @Success      200  {object}  dto.GroupActionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/settings/disappearing [post]
func (h *GroupHandler) SetDisappearingTimer(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SetDisappearingTimerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.GroupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	if req.Duration == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "duration is required")
		return
	}

	err := h.groupService.SetDisappearingTimer(r.Context(), sessionID, req.GroupJID, req.Duration)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", req.GroupJID).
			Str("duration", req.Duration).
			Msg("Failed to set disappearing timer")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", req.GroupJID).
		Str("duration", req.Duration).
		Msg("Disappearing timer set successfully")

	h.writeJSON(w, dto.GroupActionResponse{
		Success: true,
		Message: "Disappearing timer set successfully",
	})
}

// SetGroupPhoto godoc
// @Summary      Set group photo
// @Description  Set the photo of a group (JPEG format required)
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                      true  "Session ID"
// @Param        request     body      dto.SetGroupPhotoRequest    true  "Photo data"
// @Success      200  {object}  dto.SetGroupPhotoResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/photo [post]
func (h *GroupHandler) SetGroupPhoto(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.SetGroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.GroupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	if req.Image == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "image is required")
		return
	}

	// Decodificar imagem Base64
	imageData, err := decodeBase64Image(req.Image)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to decode image")
		h.writeError(w, http.StatusBadRequest, "invalid_image", "Failed to decode image: "+err.Error())
		return
	}

	pictureID, err := h.groupService.SetGroupPhoto(r.Context(), sessionID, req.GroupJID, imageData)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", req.GroupJID).
			Msg("Failed to set group photo")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", req.GroupJID).
		Str("picture_id", pictureID).
		Msg("Group photo set successfully")

	h.writeJSON(w, dto.SetGroupPhotoResponse{
		PictureID: pictureID,
	})
}

// RemoveGroupPhoto godoc
// @Summary      Remove group photo
// @Description  Remove the photo of a group
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                         true  "Session ID"
// @Param        request     body      dto.RemoveGroupPhotoRequest    true  "Remove request"
// @Success      200  {object}  dto.GroupActionResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/photo [delete]
func (h *GroupHandler) RemoveGroupPhoto(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
		return
	}

	var req dto.RemoveGroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	if req.GroupJID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "groupJid is required")
		return
	}

	err := h.groupService.RemoveGroupPhoto(r.Context(), sessionID, req.GroupJID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("group_jid", req.GroupJID).
			Msg("Failed to remove group photo")
		h.handleGroupError(w, err)
		return
	}

	h.logger.Info().
		Str("session_id", sessionID).
		Str("group_jid", req.GroupJID).
		Msg("Group photo removed successfully")

	h.writeJSON(w, dto.GroupActionResponse{
		Success: true,
		Message: "Group photo removed successfully",
	})
}

// writeJSON escreve resposta JSON
func (h *GroupHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// writeError escreve erro JSON
func (h *GroupHandler) writeError(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error:   code,
		Message: message,
	})
}

// handleGroupError trata erros específicos de grupos
func (h *GroupHandler) handleGroupError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	errMsg := err.Error()

	switch {
	case strings.Contains(errMsg, "session not found"):
		h.writeError(w, http.StatusNotFound, "session_not_found", "Session not found")
	case strings.Contains(errMsg, "not connected"):
		h.writeError(w, http.StatusPreconditionFailed, "not_connected", "Session not connected")
	case strings.Contains(errMsg, "invalid"):
		h.writeError(w, http.StatusBadRequest, "invalid_request", errMsg)
	case strings.Contains(errMsg, "required"):
		h.writeError(w, http.StatusBadRequest, "validation_error", errMsg)
	default:
		h.writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
	}
}

// decodeBase64Image decodifica imagem Base64
func decodeBase64Image(imageStr string) ([]byte, error) {
	// Remove data URI prefix se existir
	if strings.HasPrefix(imageStr, "data:") {
		parts := strings.SplitN(imageStr, ",", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid base64 data URI format")
		}
		imageStr = parts[1]
	}

	imageData, err := base64.StdEncoding.DecodeString(imageStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image: %w", err)
	}

	return imageData, nil
}
