package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"
)

// CommunityHandler gerencia requisições HTTP relacionadas a comunidades
type CommunityHandler struct {
	communityService input.CommunityService
	logger           *logger.Logger
}

// NewCommunityHandler cria uma nova instância do CommunityHandler
func NewCommunityHandler(communityService input.CommunityService, logger *logger.Logger) *CommunityHandler {
	return &CommunityHandler{
		communityService: communityService,
		logger:           logger,
	}
}

// ListCommunities lista todas as comunidades que a sessão participa
// @Summary Lista comunidades
// @Description Lista todas as comunidades que a sessão participa
// @Tags Comunidades
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} dto.ListCommunitiesResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/communities [get]
func (h *CommunityHandler) ListCommunities(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	communities, err := h.communityService.ListCommunities(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to list communities")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(communities)
}

// GetCommunityInfo obtém informações detalhadas de uma comunidade
// @Summary Obter informações da comunidade
// @Description Obtém informações detalhadas de uma comunidade específica
// @Tags Comunidades
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param communityJid query string true "JID da comunidade" example:"123456789@g.us"
// @Success 200 {object} dto.CommunityInfo
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/communities/info [get]
func (h *CommunityHandler) GetCommunityInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	communityJID := r.URL.Query().Get("communityJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if communityJID == "" {
		h.logger.Error().Msg("Community JID is required")
		http.Error(w, "Community JID is required", http.StatusBadRequest)
		return
	}

	community, err := h.communityService.GetCommunityInfo(r.Context(), sessionID, communityJID)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("community_jid", communityJID).Msg("Failed to get community info")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(community)
}

// CreateCommunity cria uma nova comunidade
// @Summary Criar comunidade
// @Description Cria uma nova comunidade WhatsApp
// @Tags Comunidades
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param request body dto.CreateCommunityRequest true "Dados da comunidade"
// @Success 201 {object} dto.CommunityInfo
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/communities [post]
func (h *CommunityHandler) CreateCommunity(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	var req dto.CreateCommunityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	community, err := h.communityService.CreateCommunity(r.Context(), sessionID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to create community")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(community)
}

// LinkGroup vincula um grupo a uma comunidade
// @Summary Vincular grupo à comunidade
// @Description Vincula um grupo existente a uma comunidade
// @Tags Comunidades
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param communityJid path string true "JID da comunidade"
// @Param request body dto.LinkGroupRequest true "Dados do grupo a vincular"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/communities/{communityJid}/link [post]
func (h *CommunityHandler) LinkGroup(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	communityJID := chi.URLParam(r, "communityJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if communityJID == "" {
		h.logger.Error().Msg("Community JID is required")
		http.Error(w, "Community JID is required", http.StatusBadRequest)
		return
	}

	var req dto.LinkGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.communityService.LinkGroup(r.Context(), sessionID, communityJID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("community_jid", communityJID).Str("group_jid", req.GroupJID).Msg("Failed to link group")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Group linked to community successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UnlinkGroup desvincula um grupo de uma comunidade
// @Summary Desvincular grupo da comunidade
// @Description Desvincula um grupo de uma comunidade
// @Tags Comunidades
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param communityJid path string true "JID da comunidade"
// @Param request body dto.UnlinkGroupRequest true "Dados do grupo a desvincular"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/communities/{communityJid}/unlink [post]
func (h *CommunityHandler) UnlinkGroup(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	communityJID := chi.URLParam(r, "communityJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if communityJID == "" {
		h.logger.Error().Msg("Community JID is required")
		http.Error(w, "Community JID is required", http.StatusBadRequest)
		return
	}

	var req dto.UnlinkGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.communityService.UnlinkGroup(r.Context(), sessionID, communityJID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("community_jid", communityJID).Str("group_jid", req.GroupJID).Msg("Failed to unlink group")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Group unlinked from community successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSubGroups obtém todos os subgrupos de uma comunidade
// @Summary Listar subgrupos da comunidade
// @Description Lista todos os subgrupos vinculados a uma comunidade
// @Tags Comunidades
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param communityJid path string true "JID da comunidade"
// @Success 200 {object} dto.ListCommunitySubGroupsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/communities/{communityJid}/groups [get]
func (h *CommunityHandler) GetSubGroups(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	communityJID := chi.URLParam(r, "communityJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if communityJID == "" {
		h.logger.Error().Msg("Community JID is required")
		http.Error(w, "Community JID is required", http.StatusBadRequest)
		return
	}

	subGroups, err := h.communityService.GetSubGroups(r.Context(), sessionID, communityJID)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("community_jid", communityJID).Msg("Failed to get sub groups")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subGroups)
}

// GetParticipants obtém todos os participantes de uma comunidade
// @Summary Listar participantes da comunidade
// @Description Lista todos os participantes de uma comunidade (de todos os grupos vinculados)
// @Tags Comunidades
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param communityJid path string true "JID da comunidade"
// @Success 200 {object} dto.ListCommunityParticipantsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/communities/{communityJid}/participants [get]
func (h *CommunityHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	communityJID := chi.URLParam(r, "communityJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if communityJID == "" {
		h.logger.Error().Msg("Community JID is required")
		http.Error(w, "Community JID is required", http.StatusBadRequest)
		return
	}

	participants, err := h.communityService.GetParticipants(r.Context(), sessionID, communityJID)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("community_jid", communityJID).Msg("Failed to get participants")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participants)
}
