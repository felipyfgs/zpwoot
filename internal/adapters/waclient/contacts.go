package waclient

import (
	"context"
	"fmt"

	"zpwoot/internal/core/ports/input"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type ContactService struct {
	waClient *WAClient
}

func NewContactService(waClient *WAClient) input.ContactService {
	return &ContactService{
		waClient: waClient,
	}
}
func (cs *ContactService) CheckUser(ctx context.Context, sessionID string, phones []string) ([]input.UserCheckResult, error) {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	resp, err := client.WAClient.IsOnWhatsApp(phones)
	if err != nil {
		return nil, fmt.Errorf("failed to check users: %w", err)
	}

	results := make([]input.UserCheckResult, 0, len(resp))

	for _, item := range resp {
		result := input.UserCheckResult{
			Query:        item.Query,
			IsInWhatsApp: item.IsIn,
			JID:          item.JID.String(),
		}

		if item.VerifiedName != nil {
			result.VerifiedName = item.VerifiedName.Details.GetVerifiedName()
		}

		results = append(results, result)
	}

	return results, nil
}
func (cs *ContactService) GetUser(ctx context.Context, sessionID string, phone string) (*input.UserDetail, error) {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	jid, err := parseJID(phone)
	if err != nil {
		return nil, ErrInvalidJID
	}

	info, err := client.WAClient.GetUserInfo([]types.JID{jid})
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	detail := &input.UserDetail{
		JID: jid.String(),
	}

	if userInfo, ok := info[jid]; ok {
		if userInfo.VerifiedName != nil {
			detail.VerifiedName = userInfo.VerifiedName.Details.GetVerifiedName()
		}

		if userInfo.Status != "" {
			detail.Status = userInfo.Status
		}

		if userInfo.PictureID != "" {
			detail.PictureID = userInfo.PictureID
		}
	}

	return detail, nil
}
func (cs *ContactService) GetAvatar(ctx context.Context, sessionID string, phone string, preview bool) (*input.AvatarInfo, error) {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	jid, err := parseJID(phone)
	if err != nil {
		return nil, ErrInvalidJID
	}

	pic, err := client.WAClient.GetProfilePictureInfo(jid, &whatsmeow.GetProfilePictureParams{
		Preview: preview,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar: %w", err)
	}

	if pic == nil {
		return nil, fmt.Errorf("no avatar found")
	}

	avatarInfo := &input.AvatarInfo{
		URL:       pic.URL,
		ID:        pic.ID,
		Type:      pic.Type,
		DirectURL: pic.DirectPath,
	}

	return avatarInfo, nil
}
func (cs *ContactService) GetContacts(ctx context.Context, sessionID string) ([]input.Contact, error) {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	contacts, err := client.WAClient.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	results := make([]input.Contact, 0, len(contacts))

	for jid, contact := range contacts {
		c := input.Contact{
			JID:    jid.String(),
			Name:   contact.FullName,
			Notify: contact.PushName,
		}

		if contact.BusinessName != "" {
			c.BusinessName = contact.BusinessName
		}

		results = append(results, c)
	}

	return results, nil
}
func (cs *ContactService) SendPresence(ctx context.Context, sessionID string, presence string) error {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	var presenceType types.Presence

	switch presence {
	case "available":
		presenceType = types.PresenceAvailable
	case "unavailable":
		presenceType = types.PresenceUnavailable
	default:
		return fmt.Errorf("invalid presence type: %s", presence)
	}

	err = client.WAClient.SendPresence(presenceType)
	if err != nil {
		return fmt.Errorf("failed to send presence: %w", err)
	}

	return nil
}
func (cs *ContactService) ChatPresence(ctx context.Context, sessionID string, phone string, presence string, media string) error {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	jid, err := parseJID(phone)
	if err != nil {
		return ErrInvalidJID
	}

	var chatPresence types.ChatPresence

	switch presence {
	case "composing":
		chatPresence = types.ChatPresenceComposing
	case "paused":
		chatPresence = types.ChatPresencePaused
	default:
		return fmt.Errorf("invalid chat presence: %s", presence)
	}

	var mediaType types.ChatPresenceMedia

	if media != "" {
		switch media {
		case "text":
			mediaType = types.ChatPresenceMediaText
		case "audio":
			mediaType = types.ChatPresenceMediaAudio
		default:
			mediaType = types.ChatPresenceMediaText
		}
	} else {
		mediaType = types.ChatPresenceMediaText
	}

	err = client.WAClient.SendChatPresence(jid, chatPresence, mediaType)
	if err != nil {
		return fmt.Errorf("failed to send chat presence: %w", err)
	}

	return nil
}
