package waclient

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// MessageSenderImpl implements the MessageSender interface
type MessageSenderImpl struct {
	waClient *WAClient
}

// NewMessageSender creates a new message sender
func NewMessageSender(waClient *WAClient) *MessageSenderImpl {
	return &MessageSenderImpl{
		waClient: waClient,
	}
}

// SendTextMessage sends a text message
func (ms *MessageSenderImpl) SendTextMessage(ctx context.Context, sessionID string, to string, text string) error {
	client, err := ms.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if client.Status != StatusConnected {
		return ErrNotConnected
	}

	// Parse recipient JID
	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	// Create text message
	message := &waE2E.Message{
		Conversation: proto.String(text),
	}

	// Send message
	_, err = client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return fmt.Errorf("failed to send text message: %w", err)
	}

	return nil
}

// SendMediaMessage sends a media message (placeholder - would need full media processing)
func (ms *MessageSenderImpl) SendMediaMessage(ctx context.Context, sessionID string, to string, media *MediaData) error {
	client, err := ms.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if client.Status != StatusConnected {
		return ErrNotConnected
	}

	// Parse recipient JID
	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	// For now, just send a placeholder text message
	// Full media implementation would require upload and proper message creation
	message := &waE2E.Message{
		Conversation: proto.String(fmt.Sprintf("Media message: %s (%s)", media.FileName, media.MimeType)),
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return fmt.Errorf("failed to send media message: %w", err)
	}

	return nil
}

// SendLocationMessage sends a location message
func (ms *MessageSenderImpl) SendLocationMessage(ctx context.Context, sessionID string, to string, lat, lng float64, name string) error {
	client, err := ms.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if client.Status != StatusConnected {
		return ErrNotConnected
	}

	// Parse recipient JID
	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	// Create location message
	message := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  proto.Float64(lat),
			DegreesLongitude: proto.Float64(lng),
			Name:             proto.String(name),
		},
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return fmt.Errorf("failed to send location message: %w", err)
	}

	return nil
}

// SendContactMessage sends a contact message
func (ms *MessageSenderImpl) SendContactMessage(ctx context.Context, sessionID string, to string, contact *ContactInfo) error {
	client, err := ms.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if client.Status != StatusConnected {
		return ErrNotConnected
	}

	// Parse recipient JID
	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	// Create vCard if not provided
	vcard := contact.VCard
	if vcard == "" {
		vcard = fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s\nTEL:%s\nEND:VCARD", contact.Name, contact.Phone)
	}

	// Create contact message
	message := &waE2E.Message{
		ContactMessage: &waE2E.ContactMessage{
			DisplayName: proto.String(contact.Name),
			Vcard:       proto.String(vcard),
		},
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return fmt.Errorf("failed to send contact message: %w", err)
	}

	return nil
}

// parseJID parses a string into a WhatsApp JID
func parseJID(jidStr string) (types.JID, error) {
	if jidStr == "" {
		return types.JID{}, ErrInvalidJID
	}

	// If it's just a phone number, add @s.whatsapp.net
	if !strings.Contains(jidStr, "@") {
		// Remove any non-numeric characters except +
		phone := strings.ReplaceAll(jidStr, " ", "")
		phone = strings.ReplaceAll(phone, "-", "")
		phone = strings.ReplaceAll(phone, "(", "")
		phone = strings.ReplaceAll(phone, ")", "")
		
		// Remove + if present
		if strings.HasPrefix(phone, "+") {
			phone = phone[1:]
		}

		jidStr = phone + "@s.whatsapp.net"
	}

	jid, err := types.ParseJID(jidStr)
	if err != nil {
		return types.JID{}, ErrInvalidJID
	}

	return jid, nil
}

// GetChatInfo gets information about a chat
func (ms *MessageSenderImpl) GetChatInfo(ctx context.Context, sessionID string, chatJID string) (*ChatInfo, error) {
	client, err := ms.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if client.Status != StatusConnected {
		return nil, ErrNotConnected
	}

	jid, err := parseJID(chatJID)
	if err != nil {
		return nil, ErrInvalidJID
	}

	// Get basic chat info
	chatInfo := &ChatInfo{
		JID:     jid.String(),
		IsGroup: jid.Server == types.GroupServer,
	}

	// For groups, get group info
	if chatInfo.IsGroup {
		groupInfo, err := client.WAClient.GetGroupInfo(jid)
		if err == nil {
			chatInfo.Name = groupInfo.Name
			chatInfo.Topic = groupInfo.Topic
			chatInfo.ParticipantCount = len(groupInfo.Participants)
		}
	}

	return chatInfo, nil
}

// ChatInfo represents chat information
type ChatInfo struct {
	JID              string `json:"jid"`
	Name             string `json:"name,omitempty"`
	Topic            string `json:"topic,omitempty"`
	IsGroup          bool   `json:"isGroup"`
	ParticipantCount int    `json:"participantCount,omitempty"`
}

// GetContacts gets the contact list
func (ms *MessageSenderImpl) GetContacts(ctx context.Context, sessionID string) ([]*ContactInfo, error) {
	client, err := ms.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if client.Status != StatusConnected {
		return nil, ErrNotConnected
	}

	// Get contacts from WhatsApp
	contacts, err := client.WAClient.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	var contactList []*ContactInfo
	for jid, contact := range contacts {
		contactInfo := &ContactInfo{
			Phone: jid.User,
			Name:  contact.PushName,
		}
		if contact.BusinessName != "" {
			contactInfo.Name = contact.BusinessName
		}
		contactList = append(contactList, contactInfo)
	}

	return contactList, nil
}

// GetChats gets the chat list
func (ms *MessageSenderImpl) GetChats(ctx context.Context, sessionID string) ([]*ChatInfo, error) {
	client, err := ms.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if client.Status != StatusConnected {
		return nil, ErrNotConnected
	}

	// For now, return empty list - would need proper implementation
	// based on whatsmeow's actual API for getting chats
	var chatList []*ChatInfo

	return chatList, nil
}
