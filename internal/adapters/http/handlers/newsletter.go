package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"
)

// NewsletterHandler gerencia requisições HTTP relacionadas a newsletters
type NewsletterHandler struct {
	newsletterService input.NewsletterService
	logger            *logger.Logger
}

// NewNewsletterHandler cria uma nova instância do NewsletterHandler
func NewNewsletterHandler(newsletterService input.NewsletterService, logger *logger.Logger) *NewsletterHandler {
	return &NewsletterHandler{
		newsletterService: newsletterService,
		logger:            logger,
	}
}

// ListNewsletters lista todos os newsletters que a sessão segue
// @Summary Lista newsletters
// @Description Lista todos os newsletters que a sessão segue
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} dto.ListNewslettersResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters [get]
func (h *NewsletterHandler) ListNewsletters(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	newsletters, err := h.newsletterService.ListNewsletters(r.Context(), sessionID)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to list newsletters")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsletters)
}

// GetNewsletterInfo obtém informações detalhadas de um newsletter
// @Summary Obter informações do newsletter
// @Description Obtém informações detalhadas de um newsletter específico
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param newsletterJid query string true "JID do newsletter" example:"123456789@newsletter"
// @Success 200 {object} dto.NewsletterInfo
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters/info [get]
func (h *NewsletterHandler) GetNewsletterInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	newsletterJID := r.URL.Query().Get("newsletterJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if newsletterJID == "" {
		h.logger.Error().Msg("Newsletter JID is required")
		http.Error(w, "Newsletter JID is required", http.StatusBadRequest)
		return
	}

	newsletter, err := h.newsletterService.GetNewsletterInfo(r.Context(), sessionID, newsletterJID)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("newsletter_jid", newsletterJID).Msg("Failed to get newsletter info")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsletter)
}

// GetNewsletterInfoWithInvite obtém informações de um newsletter via código de convite
// @Summary Obter informações do newsletter via convite
// @Description Obtém informações de um newsletter usando código de convite
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param request body dto.NewsletterInfoWithInviteRequest true "Código de convite"
// @Success 200 {object} dto.NewsletterInfo
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters/info-invite [post]
func (h *NewsletterHandler) GetNewsletterInfoWithInvite(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	var req dto.NewsletterInfoWithInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newsletter, err := h.newsletterService.GetNewsletterInfoWithInvite(r.Context(), sessionID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to get newsletter info with invite")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsletter)
}

// CreateNewsletter cria um novo newsletter
// @Summary Criar newsletter
// @Description Cria um novo newsletter WhatsApp
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param request body dto.CreateNewsletterRequest true "Dados do newsletter"
// @Success 201 {object} dto.NewsletterInfo
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters [post]
func (h *NewsletterHandler) CreateNewsletter(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	var req dto.CreateNewsletterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newsletter, err := h.newsletterService.CreateNewsletter(r.Context(), sessionID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to create newsletter")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newsletter)
}

// FollowNewsletter segue um newsletter
// @Summary Seguir newsletter
// @Description Segue um newsletter por JID ou código de convite
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param request body dto.FollowNewsletterRequest true "Dados para seguir newsletter"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters/follow [post]
func (h *NewsletterHandler) FollowNewsletter(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	var req dto.FollowNewsletterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.newsletterService.FollowNewsletter(r.Context(), sessionID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to follow newsletter")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Newsletter followed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UnfollowNewsletter deixa de seguir um newsletter
// @Summary Deixar de seguir newsletter
// @Description Para de seguir um newsletter específico
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param newsletterJid path string true "JID do newsletter"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters/{newsletterJid}/unfollow [post]
func (h *NewsletterHandler) UnfollowNewsletter(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	newsletterJID := chi.URLParam(r, "newsletterJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if newsletterJID == "" {
		h.logger.Error().Msg("Newsletter JID is required")
		http.Error(w, "Newsletter JID is required", http.StatusBadRequest)
		return
	}

	err := h.newsletterService.UnfollowNewsletter(r.Context(), sessionID, newsletterJID)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("newsletter_jid", newsletterJID).Msg("Failed to unfollow newsletter")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Newsletter unfollowed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMessages obtém mensagens de um newsletter
// @Summary Obter mensagens do newsletter
// @Description Lista mensagens de um newsletter com paginação
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param newsletterJid path string true "JID do newsletter"
// @Param count query int false "Número de mensagens (padrão: 50)"
// @Param before query string false "Cursor para paginação"
// @Success 200 {object} dto.ListNewsletterMessagesResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters/{newsletterJid}/messages [get]
func (h *NewsletterHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	newsletterJID := chi.URLParam(r, "newsletterJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if newsletterJID == "" {
		h.logger.Error().Msg("Newsletter JID is required")
		http.Error(w, "Newsletter JID is required", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	req := &dto.GetNewsletterMessagesRequest{}
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil {
			req.Count = count
		}
	}
	if before := r.URL.Query().Get("before"); before != "" {
		req.Before = before
	}

	messages, err := h.newsletterService.GetMessages(r.Context(), sessionID, newsletterJID, req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("newsletter_jid", newsletterJID).Msg("Failed to get newsletter messages")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// MarkViewed marca mensagens como visualizadas
// @Summary Marcar mensagens como visualizadas
// @Description Marca mensagens específicas do newsletter como visualizadas
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param newsletterJid path string true "JID do newsletter"
// @Param request body dto.NewsletterMarkViewedRequest true "IDs das mensagens"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters/{newsletterJid}/mark-viewed [post]
func (h *NewsletterHandler) MarkViewed(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	newsletterJID := chi.URLParam(r, "newsletterJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if newsletterJID == "" {
		h.logger.Error().Msg("Newsletter JID is required")
		http.Error(w, "Newsletter JID is required", http.StatusBadRequest)
		return
	}

	var req dto.NewsletterMarkViewedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.newsletterService.MarkViewed(r.Context(), sessionID, newsletterJID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("newsletter_jid", newsletterJID).Msg("Failed to mark messages as viewed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Messages marked as viewed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SendReaction envia reação a uma mensagem do newsletter
// @Summary Reagir a mensagem do newsletter
// @Description Envia uma reação (emoji) a uma mensagem específica do newsletter
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param newsletterJid path string true "JID do newsletter"
// @Param request body dto.NewsletterReactionRequest true "Dados da reação"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters/{newsletterJid}/react [post]
func (h *NewsletterHandler) SendReaction(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	newsletterJID := chi.URLParam(r, "newsletterJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if newsletterJID == "" {
		h.logger.Error().Msg("Newsletter JID is required")
		http.Error(w, "Newsletter JID is required", http.StatusBadRequest)
		return
	}

	var req dto.NewsletterReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.newsletterService.SendReaction(r.Context(), sessionID, newsletterJID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("newsletter_jid", newsletterJID).Msg("Failed to send reaction")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Reaction sent successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ToggleMute silencia ou dessilencia um newsletter
// @Summary Silenciar/dessilenciar newsletter
// @Description Alterna o estado de silenciamento de um newsletter
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param newsletterJid path string true "JID do newsletter"
// @Param request body dto.NewsletterMuteRequest true "Estado de silenciamento"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters/{newsletterJid}/mute [post]
func (h *NewsletterHandler) ToggleMute(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	newsletterJID := chi.URLParam(r, "newsletterJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if newsletterJID == "" {
		h.logger.Error().Msg("Newsletter JID is required")
		http.Error(w, "Newsletter JID is required", http.StatusBadRequest)
		return
	}

	var req dto.NewsletterMuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.newsletterService.ToggleMute(r.Context(), sessionID, newsletterJID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("newsletter_jid", newsletterJID).Msg("Failed to toggle mute")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	action := "muted"
	if !req.Mute {
		action = "unmuted"
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Newsletter " + action + " successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SendMessage envia uma mensagem para um newsletter
// @Summary Enviar mensagem para newsletter
// @Description Envia uma mensagem de texto para um newsletter (apenas para owners/admins)
// @Tags Newsletters
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param newsletterJid path string true "JID do newsletter"
// @Param request body map[string]interface{} true "Dados da mensagem"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sessions/{sessionId}/newsletters/{newsletterJid}/send [post]
func (h *NewsletterHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	newsletterJID := chi.URLParam(r, "newsletterJid")

	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	if newsletterJID == "" {
		h.logger.Error().Msg("Newsletter JID is required")
		http.Error(w, "Newsletter JID is required", http.StatusBadRequest)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Implementar envio de mensagem para newsletter
	// Por enquanto, simular sucesso
	response := map[string]interface{}{
		"success":    true,
		"message":    "Message sent to newsletter successfully",
		"message_id": "temp_id_" + newsletterJID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
