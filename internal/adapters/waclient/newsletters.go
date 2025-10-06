package waclient

import (
	"context"
	"fmt"
	"strconv"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"

	"go.mau.fi/whatsmeow"
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
			JID:             newsletter.ID.String(),
			Name:            newsletter.ThreadMeta.Name.Text,
			Description:     newsletter.ThreadMeta.Description.Text,
			SubscriberCount: newsletter.ThreadMeta.SubscriberCount,
			IsOwner:         newsletter.ViewerMeta != nil && newsletter.ViewerMeta.Role == "owner",
			IsFollowing:     true,
			IsMuted:         newsletter.ViewerMeta != nil && newsletter.ViewerMeta.Mute == "on",
			CreatedAt:       newsletter.ThreadMeta.CreationTime.Unix(),
		}

		response.Newsletters = append(response.Newsletters, info)
	}

	return response, nil
}
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
		JID:             newsletter.ID.String(),
		Name:            newsletter.ThreadMeta.Name.Text,
		Description:     newsletter.ThreadMeta.Description.Text,
		SubscriberCount: newsletter.ThreadMeta.SubscriberCount,
		IsOwner:         newsletter.ViewerMeta != nil && newsletter.ViewerMeta.Role == "owner",
		IsFollowing:     newsletter.ViewerMeta != nil,
		IsMuted:         newsletter.ViewerMeta != nil && newsletter.ViewerMeta.Mute == "on",
		CreatedAt:       newsletter.ThreadMeta.CreationTime.Unix(),
	}, nil
}
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
		JID:             newsletter.ID.String(),
		Name:            newsletter.ThreadMeta.Name.Text,
		Description:     newsletter.ThreadMeta.Description.Text,
		SubscriberCount: newsletter.ThreadMeta.SubscriberCount,
		IsOwner:         newsletter.ViewerMeta != nil && newsletter.ViewerMeta.Role == "owner",
		IsFollowing:     newsletter.ViewerMeta != nil,
		IsMuted:         newsletter.ViewerMeta != nil && newsletter.ViewerMeta.Mute == "on",
		CreatedAt:       newsletter.ThreadMeta.CreationTime.Unix(),
	}, nil
}
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

	return nil, fmt.Errorf("create newsletter not yet implemented in whatsmeow")
}
func (ns *NewsletterService) FollowNewsletter(ctx context.Context, sessionID string, req *dto.FollowNewsletterRequest) error {
	client, err := ns.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	if req.NewsletterJID != "" {

		jid, err := parseJID(req.NewsletterJID)
		if err != nil {
			return fmt.Errorf("invalid newsletter JID: %w", err)
		}

		err = client.WAClient.FollowNewsletter(jid)
		if err != nil {
			return fmt.Errorf("failed to follow newsletter: %w", err)
		}
	} else if req.InviteCode != "" {
		newsletter, err := client.WAClient.GetNewsletterInfoWithInvite(req.InviteCode)
		if err != nil {
			return fmt.Errorf("failed to get newsletter info with invite: %w", err)
		}
		err = client.WAClient.FollowNewsletter(newsletter.ID)
		if err != nil {
			return fmt.Errorf("failed to follow newsletter: %w", err)
		}
	} else {
		return fmt.Errorf("either newsletter_jid or invite_code is required")
	}

	return nil
}
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
	params := &whatsmeow.GetNewsletterMessagesParams{
		Count: 50,
	}
	if req.Count > 0 {
		params.Count = req.Count
	}
	if req.Before != "" {

		if beforeID, err := strconv.Atoi(req.Before); err == nil {
			params.Before = beforeID
		}
	}

	messages, err := client.WAClient.GetNewsletterMessages(jid, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletter messages: %w", err)
	}

	response := &dto.ListNewsletterMessagesResponse{
		Messages: make([]dto.NewsletterMessage, 0, len(messages)),
		HasMore:  false,
		Cursor:   "",
	}

	for _, msg := range messages {
		message := dto.NewsletterMessage{
			ID:        msg.MessageID,
			ServerID:  string(msg.MessageServerID),
			Content:   "",
			Type:      msg.Type,
			Timestamp: msg.Timestamp.Unix(),
			ViewCount: msg.ViewsCount,
		}

		response.Messages = append(response.Messages, message)
	}

	return response, nil
}
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
	serverIDs := make([]types.MessageServerID, len(req.ServerIDs))
	for i, serverIDStr := range req.ServerIDs {

		if serverID, err := strconv.Atoi(serverIDStr); err == nil {
			serverIDs[i] = serverID
		}
	}

	err = client.WAClient.NewsletterMarkViewed(jid, serverIDs)
	if err != nil {
		return fmt.Errorf("failed to mark messages as viewed: %w", err)
	}

	return nil
}
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
	serverIDInt, err := strconv.Atoi(req.ServerID)
	if err != nil {
		return fmt.Errorf("invalid server ID: %w", err)
	}
	serverID := serverIDInt
	messageID := req.MessageID

	err = client.WAClient.NewsletterSendReaction(jid, serverID, req.Reaction, messageID)
	if err != nil {
		return fmt.Errorf("failed to send reaction: %w", err)
	}

	return nil
}
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
