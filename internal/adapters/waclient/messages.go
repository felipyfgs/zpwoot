package waclient

import (
	"context"
	"fmt"
	"strings"

	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type MessageSenderImpl struct {
	waClient *WAClient
}

func NewMessageSender(waClient *WAClient) *MessageSenderImpl {
	return &MessageSenderImpl{waClient: waClient}
}

func (ms *MessageSenderImpl) SendTextMessage(ctx context.Context, sessionID string, to string, text string) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	message := &waE2E.Message{Conversation: proto.String(text)}

	if _, err = client.WAClient.SendMessage(ctx, recipientJID, message); err != nil {
		return fmt.Errorf("failed to send text message: %w", err)
	}

	return nil
}

func (ms *MessageSenderImpl) SendMediaMessage(ctx context.Context, sessionID string, to string, media *output.MediaData) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	message := &waE2E.Message{
		Conversation: proto.String(fmt.Sprintf("Media message: %s (%s)", media.FileName, media.MimeType)),
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return fmt.Errorf("failed to send media message: %w", err)
	}

	return nil
}

func (ms *MessageSenderImpl) SendLocationMessage(ctx context.Context, sessionID string, to string, lat, lng float64, name string) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

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

func (ms *MessageSenderImpl) SendContactMessage(ctx context.Context, sessionID string, to string, contact *ContactInfo) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	vcard := contact.VCard
	if vcard == "" {
		vcard = ms.generateVCard(contact.Name, contact.Phone)
	}

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

func (ms *MessageSenderImpl) getConnectedClient(ctx context.Context, sessionID string) (*Client, error) {
	client, err := ms.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	return client, nil
}

func (ms *MessageSenderImpl) generateVCard(name, phone string) string {
	return fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s\nTEL:%s\nEND:VCARD", name, phone)
}

func parseJID(jidStr string) (types.JID, error) {
	if jidStr == "" {
		return types.JID{}, ErrInvalidJID
	}

	if !strings.Contains(jidStr, "@") {
		phone := normalizePhoneNumber(jidStr)
		jidStr = phone + "@s.whatsapp.net"
	}

	jid, err := types.ParseJID(jidStr)
	if err != nil {
		return types.JID{}, ErrInvalidJID
	}

	return jid, nil
}

func normalizePhoneNumber(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.TrimPrefix(phone, "+")
	return phone
}

func (ms *MessageSenderImpl) GetChatInfo(ctx context.Context, sessionID string, chatJID string) (*ChatInfo, error) {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	jid, err := parseJID(chatJID)
	if err != nil {
		return nil, ErrInvalidJID
	}

	chatInfo := &ChatInfo{
		JID:     jid.String(),
		IsGroup: jid.Server == types.GroupServer,
	}

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

type ChatInfo struct {
	JID              string `json:"jid"`
	Name             string `json:"name,omitempty"`
	Topic            string `json:"topic,omitempty"`
	IsGroup          bool   `json:"isGroup"`
	ParticipantCount int    `json:"participantCount,omitempty"`
}

func (ms *MessageSenderImpl) GetContacts(ctx context.Context, sessionID string) ([]*ContactInfo, error) {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return nil, err
	}

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

func (ms *MessageSenderImpl) SendContactMessageFromInput(ctx context.Context, sessionID string, to string, contact *input.ContactInfo) error {

	internalContact := &ContactInfo{
		Name:  contact.Name,
		Phone: contact.Phone,
		VCard: contact.VCard,
	}
	return ms.SendContactMessage(ctx, sessionID, to, internalContact)
}

func (ms *MessageSenderImpl) GetChatInfoAsInput(ctx context.Context, sessionID, chatJID string) (*input.ChatInfo, error) {
	chatInfo, err := ms.GetChatInfo(ctx, sessionID, chatJID)
	if err != nil {
		return nil, err
	}

	return &input.ChatInfo{
		JID:              chatInfo.JID,
		Name:             chatInfo.Name,
		Topic:            chatInfo.Topic,
		IsGroup:          chatInfo.IsGroup,
		ParticipantCount: chatInfo.ParticipantCount,
	}, nil
}

func (ms *MessageSenderImpl) GetContactsAsInput(ctx context.Context, sessionID string) ([]*input.ContactInfo, error) {
	contacts, err := ms.GetContacts(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	var inputContacts []*input.ContactInfo
	for _, contact := range contacts {
		inputContacts = append(inputContacts, &input.ContactInfo{
			Name:  contact.Name,
			Phone: contact.Phone,
			VCard: contact.VCard,
		})
	}

	return inputContacts, nil
}

func (ms *MessageSenderImpl) GetChatsAsInput(ctx context.Context, sessionID string) ([]*input.ChatInfo, error) {
	chats, err := ms.GetChats(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	var inputChats []*input.ChatInfo
	for _, chat := range chats {
		inputChats = append(inputChats, &input.ChatInfo{
			JID:              chat.JID,
			Name:             chat.Name,
			Topic:            chat.Topic,
			IsGroup:          chat.IsGroup,
			ParticipantCount: chat.ParticipantCount,
		})
	}

	return inputChats, nil
}

type MessageServiceWrapper struct {
	*MessageSenderImpl
}

func NewMessageServiceWrapper(messageSender *MessageSenderImpl) input.MessageService {
	return &MessageServiceWrapper{
		MessageSenderImpl: messageSender,
	}
}

func (w *MessageServiceWrapper) SendContactMessage(ctx context.Context, sessionID string, to string, contact *input.ContactInfo) error {
	return w.SendContactMessageFromInput(ctx, sessionID, to, contact)
}

func (w *MessageServiceWrapper) GetChatInfo(ctx context.Context, sessionID, chatJID string) (*input.ChatInfo, error) {
	return w.GetChatInfoAsInput(ctx, sessionID, chatJID)
}

func (w *MessageServiceWrapper) GetContacts(ctx context.Context, sessionID string) ([]*input.ContactInfo, error) {
	return w.GetContactsAsInput(ctx, sessionID)
}

func (w *MessageServiceWrapper) GetChats(ctx context.Context, sessionID string) ([]*input.ChatInfo, error) {
	return w.GetChatsAsInput(ctx, sessionID)
}

func (ms *MessageSenderImpl) GetChats(ctx context.Context, sessionID string) ([]*ChatInfo, error) {
	_, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	var chatList []*ChatInfo

	return chatList, nil
}
