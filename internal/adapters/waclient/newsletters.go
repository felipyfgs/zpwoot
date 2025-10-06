package waclient

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"

	"go.mau.fi/whatsmeow/types"
)

type NewsletterService struct {
	waClient *WAClient
}

func NewNewsletterService(waClient *WAClient) input.NewsletterService {
	return &NewsletterService{
		waClient: waClient,
	}
}

// ListNewsletters lista todos os newsletters que a sessão segue
func (ns *NewsletterService) ListNewsletters(ctx context.Context, sessionID string) (*dto.ListNewslettersResponse, error) {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	newsletters, err := client.WAClient.GetSubscribedNewsletters()
	if err != nil {
		return nil, fmt.Errorf("failed to get subscribed newsletters: %w", err)
	}

	response := &dto.ListNewslettersResponse{
		Newsletters: make([]dto.NewsletterInfo, 0, len(newsletters)),
	}

	for _, newsletter := range newsletters {
		info := dto.NewsletterInfo{
			JID:             newsletter.JID.String(),
			Name:            newsletter.Name,
			Description:     newsletter.Description,
			SubscriberCount: int(newsletter.SubscriberCount),
			IsOwner:         newsletter.ViewerMetadata != nil && newsletter.ViewerMetadata.Role == "OWNER",
			IsFollowing:     true, // Se está na lista, está seguindo
			IsMuted:         newsletter.ViewerMetadata != nil && newsletter.ViewerMetadata.Mute == "ON",
			CreatedAt:       newsletter.CreationTime,
		}

		response.Newsletters = append(response.Newsletters, info)
	}

	return response, nil
}

// GetNewsletterInfo obtém informações detalhadas de um newsletter
func (ns *NewsletterService) GetNewsletterInfo(ctx context.Context, sessionID string, newsletterJID string) (*dto.NewsletterInfo, error) {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	jid, err := parseJID(newsletterJID)
	if err != nil {
		return nil, fmt.Errorf("invalid newsletter JID: %w", err)
	}

	newsletter, err := client.WAClient.GetNewsletterInfo(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter info: %w", err)
	}

	return &dto.NewsletterInfo{
		JID:             newsletter.JID.String(),
		Name:            newsletter.Name,
		Description:     newsletter.Description,
		SubscriberCount: int(newsletter.SubscriberCount),
		IsOwner:         newsletter.ViewerMetadata != nil && newsletter.ViewerMetadata.Role == "OWNER",
		IsFollowing:     newsletter.ViewerMetadata != nil,
		IsMuted:         newsletter.ViewerMetadata != nil && newsletter.ViewerMetadata.Mute == "ON",
		CreatedAt:       newsletter.CreationTime,
	}, nil
}

// GetNewsletterInfoWithInvite obtém informações de um newsletter via código de convite
func (ns *NewsletterService) GetNewsletterInfoWithInvite(ctx context.Context, sessionID string, req *dto.NewsletterInfoWithInviteRequest) (*dto.NewsletterInfo, error) {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	if req.InviteKey == "" {
		return nil, fmt.Errorf("invite key is required")
	}

	newsletter, err := client.WAClient.GetNewsletterInfoWithInvite(req.InviteKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter info with invite: %w", err)
	}

	return &dto.NewsletterInfo{
		JID:             newsletter.JID.String(),
		Name:            newsletter.Name,
		Description:     newsletter.Description,
		SubscriberCount: int(newsletter.SubscriberCount),
		IsOwner:         newsletter.ViewerMetadata != nil && newsletter.ViewerMetadata.Role == "OWNER",
		IsFollowing:     newsletter.ViewerMetadata != nil,
		IsMuted:         newsletter.ViewerMetadata != nil && newsletter.ViewerMetadata.Mute == "ON",
		CreatedAt:       newsletter.CreationTime,
	}, nil
}

// CreateNewsletter cria um novo newsletter
func (ns *NewsletterService) CreateNewsletter(ctx context.Context, sessionID string, req *dto.CreateNewsletterRequest) (*dto.NewsletterInfo, error) {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	if req.Name == "" {
		return nil, fmt.Errorf("newsletter name is required")
	}

	// TODO: Implementar CreateNewsletter quando disponível na whatsmeow
	// Por enquanto, retorna erro não implementado
	return nil, fmt.Errorf("create newsletter not yet implemented in whatsmeow")
}

// FollowNewsletter segue um newsletter
func (ns *NewsletterService) FollowNewsletter(ctx context.Context, sessionID string, req *dto.FollowNewsletterRequest) error {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	if req.NewsletterJID != "" {
		// Seguir por JID
		jid, err := parseJID(req.NewsletterJID)
		if err != nil {
			return fmt.Errorf("invalid newsletter JID: %w", err)
		}

		err = client.WAClient.FollowNewsletter(jid)
		if err != nil {
			return fmt.Errorf("failed to follow newsletter: %w", err)
		}
	} else if req.InviteCode != "" {
		// Seguir por código de convite
		// Primeiro obter info do newsletter
		newsletter, err := client.WAClient.GetNewsletterInfoWithInvite(req.InviteCode)
		if err != nil {
			return fmt.Errorf("failed to get newsletter info with invite: %w", err)
		}

		// Depois seguir
		err = client.WAClient.FollowNewsletter(newsletter.JID)
		if err != nil {
			return fmt.Errorf("failed to follow newsletter: %w", err)
		}
	} else {
		return fmt.Errorf("either newsletter_jid or invite_code is required")
	}

	return nil
}

// UnfollowNewsletter deixa de seguir um newsletter
func (ns *NewsletterService) UnfollowNewsletter(ctx context.Context, sessionID string, newsletterJID string) error {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	jid, err := parseJID(newsletterJID)
	if err != nil {
		return fmt.Errorf("invalid newsletter JID: %w", err)
	}

	err = client.WAClient.UnfollowNewsletter(jid)
	if err != nil {
		return fmt.Errorf("failed to unfollow newsletter: %w", err)
	}

	return nil
}

// GetMessages obtém mensagens de um newsletter
func (ns *NewsletterService) GetMessages(ctx context.Context, sessionID string, newsletterJID string, req *dto.GetNewsletterMessagesRequest) (*dto.ListNewsletterMessagesResponse, error) {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	jid, err := parseJID(newsletterJID)
	if err != nil {
		return nil, fmt.Errorf("invalid newsletter JID: %w", err)
	}

	// Preparar parâmetros
	params := &types.GetNewsletterMessagesParams{}
	if req.Count > 0 {
		params.Count = int32(req.Count)
	} else {
		params.Count = 50 // Padrão
	}
	if req.Before != "" {
		params.Before = req.Before
	}

	messages, err := client.WAClient.GetNewsletterMessages(jid, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter messages: %w", err)
	}

	response := &dto.ListNewsletterMessagesResponse{
		Messages: make([]dto.NewsletterMessage, 0, len(messages.Messages)),
		HasMore:  messages.HasMore,
		Cursor:   messages.Cursor,
	}

	for _, msg := range messages.Messages {
		message := dto.NewsletterMessage{
			ID:        msg.MessageID.String(),
			ServerID:  msg.ServerID.String(),
			Content:   "", // TODO: Extrair conteúdo baseado no tipo
			Type:      "text", // TODO: Determinar tipo da mensagem
			Timestamp: msg.Timestamp.Unix(),
			ViewCount: int(msg.ViewsCount),
		}

		response.Messages = append(response.Messages, message)
	}

	return response, nil
}

// MarkViewed marca mensagens como visualizadas
func (ns *NewsletterService) MarkViewed(ctx context.Context, sessionID string, newsletterJID string, req *dto.NewsletterMarkViewedRequest) error {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	jid, err := parseJID(newsletterJID)
	if err != nil {
		return fmt.Errorf("invalid newsletter JID: %w", err)
	}

	// Converter server IDs para tipos.MessageServerID
	serverIDs := make([]types.MessageServerID, len(req.ServerIDs))
	for i, serverID := range req.ServerIDs {
		serverIDs[i] = types.MessageServerID(serverID)
	}

	err = client.WAClient.NewsletterMarkViewed(jid, serverIDs)
	if err != nil {
		return fmt.Errorf("failed to mark messages as viewed: %w", err)
	}

	return nil
}

// SendReaction envia reação a uma mensagem do newsletter
func (ns *NewsletterService) SendReaction(ctx context.Context, sessionID string, newsletterJID string, req *dto.NewsletterReactionRequest) error {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	jid, err := parseJID(newsletterJID)
	if err != nil {
		return fmt.Errorf("invalid newsletter JID: %w", err)
	}

	serverID := types.MessageServerID(req.ServerID)
	messageID := types.MessageID(req.MessageID)

	err = client.WAClient.NewsletterSendReaction(jid, serverID, req.Reaction, messageID)
	if err != nil {
		return fmt.Errorf("failed to send reaction: %w", err)
	}

	return nil
}

// ToggleMute silencia ou dessilencia um newsletter
func (ns *NewsletterService) ToggleMute(ctx context.Context, sessionID string, newsletterJID string, req *dto.NewsletterMuteRequest) error {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	jid, err := parseJID(newsletterJID)
	if err != nil {
		return fmt.Errorf("invalid newsletter JID: %w", err)
	}

	err = client.WAClient.NewsletterToggleMute(jid, req.Mute)
	if err != nil {
		return fmt.Errorf("failed to toggle mute: %w", err)
	}

	return nil
}
