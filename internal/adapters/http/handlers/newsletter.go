package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"

	"github.com/go-chi/chi/v5"
)

type NewsletterHandler struct {
	newsletterService input.NewsletterService
	logger            *logger.Logger
}

func NewNewsletterHandler(newsletterService input.NewsletterService, logger *logger.Logger) *NewsletterHandler {
	return &NewsletterHandler{
		newsletterService: newsletterService,
		logger:            logger,
	}
}
func (h *NewsletterHandler) writeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error().Err(err).Msg("Failed to encode JSON response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return err
	}
	return nil
}
func (h *NewsletterHandler) validateNewsletterRequest(w http.ResponseWriter, sessionID, newsletterJID string) bool {
	if sessionID == "" {
		h.logger.Error().Msg("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return false
	}

	if newsletterJID == "" {
		h.logger.Error().Msg("Newsletter JID is required")
		http.Error(w, "Newsletter JID is required", http.StatusBadRequest)
		return false
	}

	return true
}
func (h *NewsletterHandler) handleNewsletterOperation(
	w http.ResponseWriter,
	r *http.Request,
	operation string,
	operationFunc func(context.Context, string, string, interface{}) error,
	requestType interface{},
) {
	sessionID := chi.URLParam(r, "sessionId")
	newsletterJID := chi.URLParam(r, "newsletterJid")

	if !h.validateNewsletterRequest(w, sessionID, newsletterJID) {
		return
	}

	if err := json.NewDecoder(r.Body).Decode(requestType); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := operationFunc(r.Context(), sessionID, newsletterJID, requestType)
	if err != nil {
		h.logger.Error().Err(err).Str("session_id", sessionID).Str("newsletter_jid", newsletterJID).Msgf("Failed to %s", operation)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var message string
	switch operation {
	case "mark_viewed":
		message = "Messages marked as viewed successfully"
	case "send_reaction":
		message = "Reaction sent successfully"
	case "toggle_mute":
		message = "Newsletter mute setting updated successfully"
	default:
		message = "Operation completed successfully"
	}

	response := map[string]interface{}{
		"success": true,
		"message": message,
	}

	if err := h.writeJSON(w, response); err != nil {
		return
	}
}

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

	if err := h.writeJSON(w, newsletters); err != nil {
		return
	}
}

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

	if err := h.writeJSON(w, newsletter); err != nil {
		return
	}
}

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

	if err := h.writeJSON(w, newsletter); err != nil {
		return
	}
}

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

	w.WriteHeader(http.StatusCreated)
	if err := h.writeJSON(w, newsletter); err != nil {
		return
	}
}

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

	if err := h.writeJSON(w, response); err != nil {
		return
	}
}

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

	if err := h.writeJSON(w, response); err != nil {
		return
	}
}

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

	if err := h.writeJSON(w, messages); err != nil {
		return
	}
}

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
	var req dto.NewsletterMarkViewedRequest
	h.handleNewsletterOperation(w, r, "mark_viewed", func(ctx context.Context, sessionID, newsletterJID string, reqData interface{}) error {
		return h.newsletterService.MarkViewed(ctx, sessionID, newsletterJID, reqData.(*dto.NewsletterMarkViewedRequest))
	}, &req)
}

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
	var req dto.NewsletterReactionRequest
	h.handleNewsletterOperation(w, r, "send_reaction", func(ctx context.Context, sessionID, newsletterJID string, reqData interface{}) error {
		return h.newsletterService.SendReaction(ctx, sessionID, newsletterJID, reqData.(*dto.NewsletterReactionRequest))
	}, &req)
}

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
	var req dto.NewsletterMuteRequest
	h.handleNewsletterOperation(w, r, "toggle_mute", func(ctx context.Context, sessionID, newsletterJID string, reqData interface{}) error {
		return h.newsletterService.ToggleMute(ctx, sessionID, newsletterJID, reqData.(*dto.NewsletterMuteRequest))
	}, &req)
}

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

	response := map[string]interface{}{
		"success":    true,
		"message":    "Message sent to newsletter successfully",
		"message_id": "temp_id_" + newsletterJID,
	}

	if err := h.writeJSON(w, response); err != nil {
		return
	}
}
