package waclient

import (
	"context"
	"fmt"
	"strings"

	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"

	"go.mau.fi/whatsmeow"
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

func (ms *MessageSenderImpl) SendTextMessage(ctx context.Context, sessionID string, to string, text string) (*whatsmeow.SendResponse, error) {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return nil, ErrInvalidJID
	}

	message := &waE2E.Message{Conversation: proto.String(text)}

	resp, err := client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send text message: %w", err)
	}

	return &resp, nil
}

func (ms *MessageSenderImpl) SendMediaMessage(ctx context.Context, sessionID string, to string, media *output.MediaData) (*whatsmeow.SendResponse, error) {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return nil, ErrInvalidJID
	}

	var fileData []byte
	if len(media.Data) > 0 {
		fileData = media.Data
	} else {
		return nil, fmt.Errorf("media data is required")
	}

	var mediaType whatsmeow.MediaType
	mimeType := media.MimeType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	switch {
	case strings.HasPrefix(mimeType, "image/"):
		mediaType = whatsmeow.MediaImage
	case strings.HasPrefix(mimeType, "video/"):
		mediaType = whatsmeow.MediaVideo
	case strings.HasPrefix(mimeType, "audio/"):
		mediaType = whatsmeow.MediaAudio
	default:
		mediaType = whatsmeow.MediaDocument
	}

	uploaded, err := client.WAClient.Upload(ctx, fileData, mediaType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload media: %w", err)
	}

	var message *waE2E.Message
	switch mediaType {
	case whatsmeow.MediaImage:
		message = &waE2E.Message{
			ImageMessage: &waE2E.ImageMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mimeType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(fileData))),
				Caption:       proto.String(media.Caption),
			},
		}
	case whatsmeow.MediaVideo:
		message = &waE2E.Message{
			VideoMessage: &waE2E.VideoMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mimeType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(fileData))),
				Caption:       proto.String(media.Caption),
			},
		}
	case whatsmeow.MediaAudio:
		ptt := true
		message = &waE2E.Message{
			AudioMessage: &waE2E.AudioMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mimeType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(fileData))),
				PTT:           &ptt,
			},
		}
	case whatsmeow.MediaDocument:
		fileName := media.FileName
		if fileName == "" {
			fileName = "document"
		}
		message = &waE2E.Message{
			DocumentMessage: &waE2E.DocumentMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mimeType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(fileData))),
				FileName:      proto.String(fileName),
				Caption:       proto.String(media.Caption),
			},
		}
	}

	resp, err := client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send media message: %w", err)
	}

	return &resp, nil
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

func (w *MessageServiceWrapper) SendTextMessage(ctx context.Context, sessionID string, to string, text string) error {
	_, err := w.MessageSenderImpl.SendTextMessage(ctx, sessionID, to, text)
	return err
}

func (w *MessageServiceWrapper) SendMediaMessage(ctx context.Context, sessionID string, to string, media *output.MediaData) error {
	_, err := w.MessageSenderImpl.SendMediaMessage(ctx, sessionID, to, media)
	return err
}

func (w *MessageServiceWrapper) SendContactMessage(ctx context.Context, sessionID string, to string, contact *input.ContactInfo) error {
	return w.SendContactMessageFromInput(ctx, sessionID, to, contact)
}

func (w *MessageServiceWrapper) SendReactionMessage(ctx context.Context, sessionID, to, messageID, reaction string) error {
	return w.MessageSenderImpl.SendReactionMessage(ctx, sessionID, to, messageID, reaction)
}

func (w *MessageServiceWrapper) SendPollMessage(ctx context.Context, sessionID, to, name string, options []string, selectableCount int) error {
	return w.MessageSenderImpl.SendPollMessage(ctx, sessionID, to, name, options, selectableCount)
}

func (w *MessageServiceWrapper) SendButtonsMessage(ctx context.Context, sessionID, to, text string, buttons []input.ButtonInfo) error {

	var waButtons []ButtonInfo
	for _, btn := range buttons {
		waButtons = append(waButtons, ButtonInfo{
			ID:   btn.ID,
			Text: btn.Text,
		})
	}
	return w.MessageSenderImpl.SendButtonsMessage(ctx, sessionID, to, text, waButtons)
}

func (w *MessageServiceWrapper) SendListMessage(ctx context.Context, sessionID, to, text, title string, sections []input.ListSectionInfo) error {

	var waSections []ListSection
	for _, section := range sections {
		var rows []ListRow
		for _, row := range section.Rows {
			rows = append(rows, ListRow{
				ID:          row.ID,
				Title:       row.Title,
				Description: row.Description,
			})
		}
		waSections = append(waSections, ListSection{
			Title: section.Title,
			Rows:  rows,
		})
	}
	return w.MessageSenderImpl.SendListMessage(ctx, sessionID, to, text, title, "", waSections)
}

func (w *MessageServiceWrapper) SendTemplateMessage(ctx context.Context, sessionID, to string, template input.TemplateInfo) error {
	waTemplate := TemplateInfo{
		Content: template.Content,
		Footer:  template.Footer,
	}
	return w.MessageSenderImpl.SendTemplateMessage(ctx, sessionID, to, waTemplate)
}

func (w *MessageServiceWrapper) SendViewOnceMessage(ctx context.Context, sessionID, to string, media *output.MediaData) error {
	return w.MessageSenderImpl.SendViewOnceMessage(ctx, sessionID, to, media)
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

func (ms *MessageSenderImpl) SendReactionMessage(ctx context.Context, sessionID string, to string, messageID string, reaction string) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	senderJID := recipientJID
	if recipientJID.Server == types.GroupServer {

		senderJID = recipientJID
	}

	reactionMsg := client.WAClient.BuildReaction(recipientJID, senderJID, messageID, reaction)

	_, err = client.WAClient.SendMessage(ctx, recipientJID, reactionMsg)
	if err != nil {
		return fmt.Errorf("failed to send reaction: %w", err)
	}

	return nil
}

func (ms *MessageSenderImpl) SendPollMessage(ctx context.Context, sessionID string, to string, name string, options []string, selectableCount int) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	pollMsg := client.WAClient.BuildPollCreation(name, options, selectableCount)

	_, err = client.WAClient.SendMessage(ctx, recipientJID, pollMsg)
	if err != nil {
		return fmt.Errorf("failed to send poll: %w", err)
	}

	return nil
}

func (ms *MessageSenderImpl) SendButtonsMessage(ctx context.Context, sessionID string, to string, text string, buttons []ButtonInfo) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	var waButtons []*waE2E.ButtonsMessage_Button
	for _, btn := range buttons {
		waButtons = append(waButtons, &waE2E.ButtonsMessage_Button{
			ButtonID: proto.String(btn.ID),
			ButtonText: &waE2E.ButtonsMessage_Button_ButtonText{
				DisplayText: proto.String(btn.Text),
			},
			Type: waE2E.ButtonsMessage_Button_RESPONSE.Enum(),
		})
	}

	message := &waE2E.Message{
		ButtonsMessage: &waE2E.ButtonsMessage{
			ContentText: proto.String(text),
			Buttons:     waButtons,
		},
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return fmt.Errorf("failed to send buttons message: %w", err)
	}

	return nil
}

func (ms *MessageSenderImpl) SendListMessage(ctx context.Context, sessionID string, to string, text string, title string, buttonText string, sections []ListSection) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	var waSections []*waE2E.ListMessage_Section
	for _, section := range sections {
		var rows []*waE2E.ListMessage_Row
		for _, row := range section.Rows {
			rows = append(rows, &waE2E.ListMessage_Row{
				RowID:       proto.String(row.ID),
				Title:       proto.String(row.Title),
				Description: proto.String(row.Description),
			})
		}
		waSections = append(waSections, &waE2E.ListMessage_Section{
			Title: proto.String(section.Title),
			Rows:  rows,
		})
	}

	message := &waE2E.Message{
		ListMessage: &waE2E.ListMessage{
			Title:       proto.String(title),
			Description: proto.String(text),
			ButtonText:  proto.String(buttonText),
			ListType:    waE2E.ListMessage_SINGLE_SELECT.Enum(),
			Sections:    waSections,
		},
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return fmt.Errorf("failed to send list message: %w", err)
	}

	return nil
}

func (ms *MessageSenderImpl) SendTemplateMessage(ctx context.Context, sessionID string, to string, template TemplateInfo) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	message := &waE2E.Message{
		TemplateMessage: &waE2E.TemplateMessage{
			HydratedTemplate: &waE2E.TemplateMessage_HydratedFourRowTemplate{
				HydratedContentText: proto.String(template.Content),
				HydratedFooterText:  proto.String(template.Footer),
			},
		},
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return fmt.Errorf("failed to send template message: %w", err)
	}

	return nil
}

func (ms *MessageSenderImpl) SendViewOnceMessage(ctx context.Context, sessionID string, to string, media *output.MediaData) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	var fileData []byte
	if len(media.Data) > 0 {
		fileData = media.Data
	} else {
		return fmt.Errorf("media data is required")
	}

	var mediaType whatsmeow.MediaType
	mimeType := media.MimeType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if strings.HasPrefix(mimeType, "image/") {
		mediaType = whatsmeow.MediaImage
	} else if strings.HasPrefix(mimeType, "video/") {
		mediaType = whatsmeow.MediaVideo
	} else {
		return fmt.Errorf("view once only supports image and video")
	}

	uploaded, err := client.WAClient.Upload(ctx, fileData, mediaType)
	if err != nil {
		return fmt.Errorf("failed to upload media: %w", err)
	}

	var message *waE2E.Message
	viewOnce := true

	if mediaType == whatsmeow.MediaImage {
		message = &waE2E.Message{
			ImageMessage: &waE2E.ImageMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mimeType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(fileData))),
				Caption:       proto.String(media.Caption),
				ViewOnce:      &viewOnce,
			},
		}
	} else {
		message = &waE2E.Message{
			VideoMessage: &waE2E.VideoMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mimeType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(fileData))),
				Caption:       proto.String(media.Caption),
				ViewOnce:      &viewOnce,
			},
		}
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return fmt.Errorf("failed to send view once message: %w", err)
	}

	return nil
}

type ButtonInfo struct {
	ID   string
	Text string
}

type ListSection struct {
	Title string
	Rows  []ListRow
}

type ListRow struct {
	ID          string
	Title       string
	Description string
}

type TemplateInfo struct {
	Content string
	Footer  string
}
