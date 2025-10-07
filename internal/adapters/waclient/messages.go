package waclient

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"

	"go.mau.fi/whatsmeow"
	waCommon "go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type Sender struct {
	waClient *WAClient
}

func NewSender(waClient *WAClient) *Sender {
	return &Sender{waClient: waClient}
}

func (ms *Sender) SendTextMessage(ctx context.Context, sessionID string, to string, text string, contextInfo *output.MessageContextInfo) (*whatsmeow.SendResponse, error) {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if err := ms.waitForConnection(client, 10*time.Second); err != nil {
		return nil, err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return nil, ErrInvalidJID
	}

	var message *waE2E.Message

	if contextInfo != nil && contextInfo.StanzaID != "" {
		message = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        proto.String(text),
				ContextInfo: buildContextInfo(contextInfo),
			},
		}
	} else {
		message = &waE2E.Message{Conversation: proto.String(text)}
	}

	resp, err := client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send text message: %w", err)
	}

	return &resp, nil
}

func (ms *Sender) SendMediaMessage(ctx context.Context, sessionID string, to string, media *output.MediaData) (*whatsmeow.SendResponse, error) {
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
		imgMsg := &waE2E.ImageMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(fileData))),
			Caption:       proto.String(media.Caption),
		}
		if media.ViewOnce {
			imgMsg.ViewOnce = proto.Bool(true)
		}

		message = &waE2E.Message{
			ImageMessage: imgMsg,
		}
	case whatsmeow.MediaVideo:
		vidMsg := &waE2E.VideoMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(fileData))),
			Caption:       proto.String(media.Caption),
		}
		if media.ViewOnce {
			vidMsg.ViewOnce = proto.Bool(true)
		}

		message = &waE2E.Message{
			VideoMessage: vidMsg,
		}
	case whatsmeow.MediaAudio:
		ptt := true

		audioMsg := &waE2E.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(fileData))),
			PTT:           &ptt,
		}
		if media.ViewOnce {
			audioMsg.ViewOnce = proto.Bool(true)
		}

		message = &waE2E.Message{
			AudioMessage: audioMsg,
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

func (ms *Sender) SendLocationMessage(ctx context.Context, sessionID string, to string, lat, lng float64, name string) error {
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

func (ms *Sender) SendContactMessage(ctx context.Context, sessionID string, to string, contact *ContactInfo) error {
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

func (ms *Sender) getConnectedClient(ctx context.Context, sessionID string) (*Client, error) {
	client, err := ms.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (ms *Sender) waitForConnection(client *Client, timeout time.Duration) error {
	if client.WAClient.IsConnected() {
		return nil
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ticker.C:
			if client.WAClient.IsConnected() {
				return nil
			}
		case <-timeoutTimer.C:
			return fmt.Errorf("timeout waiting for websocket connection")
		}
	}
}

func (ms *Sender) generateVCard(name, phone string) string {
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

func (ms *Sender) GetChatInfo(ctx context.Context, sessionID string, chatJID string) (*ChatInfo, error) {
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

func (ms *Sender) GetContacts(ctx context.Context, sessionID string) ([]*ContactInfo, error) {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	contacts, err := client.WAClient.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	contactList := make([]*ContactInfo, 0, len(contacts))

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

func (ms *Sender) SendContactMessageFromInput(ctx context.Context, sessionID string, to string, contact *input.ContactInfo) error {
	internalContact := &ContactInfo{
		Name:  contact.Name,
		Phone: contact.Phone,
		VCard: contact.VCard,
	}

	return ms.SendContactMessage(ctx, sessionID, to, internalContact)
}

func (ms *Sender) SendContactsArrayMessage(ctx context.Context, sessionID string, to string, contacts []*input.ContactInfo) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	contactMessages := make([]*waE2E.ContactMessage, len(contacts))

	for i, contact := range contacts {
		vcard := contact.VCard
		if vcard == "" {
			vcard = ms.generateVCard(contact.Name, contact.Phone)
		}

		contactMessages[i] = &waE2E.ContactMessage{
			DisplayName: proto.String(contact.Name),
			Vcard:       proto.String(vcard),
		}
	}

	message := &waE2E.Message{
		ContactsArrayMessage: &waE2E.ContactsArrayMessage{
			DisplayName: proto.String("Contacts"),
			Contacts:    contactMessages,
		},
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, message)
	if err != nil {
		return fmt.Errorf("failed to send contacts array message: %w", err)
	}

	return nil
}

func (ms *Sender) GetChatInfoAsInput(ctx context.Context, sessionID, chatJID string) (*input.ChatInfo, error) {
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

func (ms *Sender) GetContactsAsInput(ctx context.Context, sessionID string) ([]*input.ContactInfo, error) {
	contacts, err := ms.GetContacts(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	inputContacts := make([]*input.ContactInfo, 0, len(contacts))
	for _, contact := range contacts {
		inputContacts = append(inputContacts, &input.ContactInfo{
			Name:  contact.Name,
			Phone: contact.Phone,
			VCard: contact.VCard,
		})
	}

	return inputContacts, nil
}

func (ms *Sender) GetChatsAsInput(ctx context.Context, sessionID string) ([]*input.ChatInfo, error) {
	chats, err := ms.GetChats(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	inputChats := make([]*input.ChatInfo, 0, len(chats))
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

type MessageService struct {
	*Sender
}

func NewMessageService(messageSender *Sender) input.MessageService {
	return &MessageService{
		Sender: messageSender,
	}
}

func (w *MessageService) SendTextMessage(ctx context.Context, sessionID string, to string, text string, contextInfo *output.MessageContextInfo) (*output.MessageResult, error) {
	resp, err := w.Sender.SendTextMessage(ctx, sessionID, to, text, contextInfo)
	if err != nil {
		return nil, err
	}

	return &output.MessageResult{
		MessageID: resp.ID,
		Status:    "sent",
		SentAt:    resp.Timestamp,
	}, nil
}

func (w *MessageService) SendMediaMessage(ctx context.Context, sessionID string, to string, media *output.MediaData, contextInfo *output.MessageContextInfo) (*output.MessageResult, error) {
	resp, err := w.Sender.SendMediaMessage(ctx, sessionID, to, media)
	if err != nil {
		return nil, err
	}

	return &output.MessageResult{
		MessageID: resp.ID,
		Status:    "sent",
		SentAt:    resp.Timestamp,
	}, nil
}

func (w *MessageService) SendLocationMessage(ctx context.Context, sessionID, to string, latitude, longitude float64, name string, contextInfo *output.MessageContextInfo) (*output.MessageResult, error) {
	err := w.Sender.SendLocationMessage(ctx, sessionID, to, latitude, longitude, name)
	if err != nil {
		return nil, err
	}

	return &output.MessageResult{
		MessageID: fallbackID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}, nil
}

func (w *MessageService) SendContactMessage(ctx context.Context, sessionID string, to string, contact *input.ContactInfo, contextInfo *output.MessageContextInfo) (*output.MessageResult, error) {
	err := w.SendContactMessageFromInput(ctx, sessionID, to, contact)
	if err != nil {
		return nil, err
	}

	return &output.MessageResult{
		MessageID: fallbackID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}, nil
}

func (w *MessageService) SendContactsArrayMessage(ctx context.Context, sessionID, to string, contacts []*input.ContactInfo) (*output.MessageResult, error) {
	err := w.Sender.SendContactsArrayMessage(ctx, sessionID, to, contacts)
	if err != nil {
		return nil, err
	}

	return &output.MessageResult{
		MessageID: fallbackID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}, nil
}

func (w *MessageService) SendReactionMessage(ctx context.Context, sessionID, to, messageID, reaction string, fromMe bool) (*output.MessageResult, error) {
	err := w.Sender.SendReactionMessage(ctx, sessionID, to, messageID, reaction, fromMe)
	if err != nil {
		return nil, err
	}

	return &output.MessageResult{
		MessageID: fallbackID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}, nil
}

func (w *MessageService) SendPollMessage(ctx context.Context, sessionID, to, name string, options []string, selectableCount int) (*output.MessageResult, error) {
	err := w.Sender.SendPollMessage(ctx, sessionID, to, name, options, selectableCount)
	if err != nil {
		return nil, err
	}

	return &output.MessageResult{
		MessageID: fallbackID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}, nil
}

func (w *MessageService) SendButtonsMessage(ctx context.Context, sessionID, to, text string, buttons []input.ButtonInfo) (*output.MessageResult, error) {
	waButtons := make([]ButtonInfo, 0, len(buttons))
	for _, btn := range buttons {
		waButtons = append(waButtons, ButtonInfo{
			ID:   btn.ID,
			Text: btn.Text,
		})
	}

	err := w.Sender.SendButtonsMessage(ctx, sessionID, to, text, waButtons)
	if err != nil {
		return nil, err
	}

	return &output.MessageResult{
		MessageID: fallbackID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}, nil
}

func (w *MessageService) SendListMessage(ctx context.Context, sessionID, to, text, title string, sections []input.ListSectionInfo) (*output.MessageResult, error) {
	waSections := make([]ListSection, 0, len(sections))

	for _, section := range sections {
		rows := make([]ListRow, 0, len(section.Rows))
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

	err := w.Sender.SendListMessage(ctx, sessionID, to, text, title, "", waSections)
	if err != nil {
		return nil, err
	}

	return &output.MessageResult{
		MessageID: fallbackID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}, nil
}

func (w *MessageService) SendTemplateMessage(ctx context.Context, sessionID, to string, template input.TemplateInfo) (*output.MessageResult, error) {
	waTemplate := TemplateInfo{
		Content: template.Content,
		Footer:  template.Footer,
	}

	err := w.Sender.SendTemplateMessage(ctx, sessionID, to, waTemplate)
	if err != nil {
		return nil, err
	}

	return &output.MessageResult{
		MessageID: fallbackID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}, nil
}

func (w *MessageService) GetChatInfo(ctx context.Context, sessionID, chatJID string) (*input.ChatInfo, error) {
	return w.GetChatInfoAsInput(ctx, sessionID, chatJID)
}

func (w *MessageService) GetContacts(ctx context.Context, sessionID string) ([]*input.ContactInfo, error) {
	return w.GetContactsAsInput(ctx, sessionID)
}

func (w *MessageService) GetChats(ctx context.Context, sessionID string) ([]*input.ChatInfo, error) {
	return w.GetChatsAsInput(ctx, sessionID)
}

func (w *MessageService) DeleteMessage(ctx context.Context, sessionID, phone, messageID string) error {
	return w.Sender.DeleteMessage(ctx, sessionID, phone, messageID)
}

func (w *MessageService) EditMessage(ctx context.Context, sessionID, phone, messageID, text string) error {
	return w.Sender.EditMessage(ctx, sessionID, phone, messageID, text)
}

func (w *MessageService) MarkRead(ctx context.Context, sessionID, phone string, messageIDs []string) error {
	return w.Sender.MarkRead(ctx, sessionID, phone, messageIDs)
}

func (w *MessageService) RequestHistorySync(ctx context.Context, sessionID string, count int) error {
	return w.Sender.RequestHistorySync(ctx, sessionID, count)
}

func (ms *Sender) GetChats(ctx context.Context, sessionID string) ([]*ChatInfo, error) {
	_, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	var chatList []*ChatInfo

	return chatList, nil
}

func (ms *Sender) SendReactionMessage(ctx context.Context, sessionID string, to string, messageID string, reaction string, fromMe bool) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(to)
	if err != nil {
		return ErrInvalidJID
	}

	if reaction == "remove" {
		reaction = ""
	}

	reactionMsg := &waE2E.Message{
		ReactionMessage: &waE2E.ReactionMessage{
			Key: &waCommon.MessageKey{
				RemoteJID: proto.String(recipientJID.String()),
				FromMe:    proto.Bool(fromMe),
				ID:        proto.String(messageID),
			},
			Text:              proto.String(reaction),
			GroupingKey:       proto.String(reaction),
			SenderTimestampMS: proto.Int64(time.Now().UnixMilli()),
		},
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, reactionMsg)
	if err != nil {
		return fmt.Errorf("failed to send reaction: %w", err)
	}

	return nil
}

func (ms *Sender) SendPollMessage(ctx context.Context, sessionID string, to string, name string, options []string, selectableCount int) error {
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

func (ms *Sender) SendButtonsMessage(ctx context.Context, sessionID string, to string, text string, buttons []ButtonInfo) error {
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

func (ms *Sender) SendListMessage(ctx context.Context, sessionID string, to string, text string, title string, buttonText string, sections []ListSection) error {
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

func (ms *Sender) SendTemplateMessage(ctx context.Context, sessionID string, to string, template TemplateInfo) error {
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

func fallbackID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return strings.ToUpper(hex.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano()))))
	}

	return strings.ToUpper(hex.EncodeToString(bytes))
}

func buildContextInfo(contextInfo *output.MessageContextInfo) *waE2E.ContextInfo {
	if contextInfo == nil {
		return nil
	}

	ctx := &waE2E.ContextInfo{}

	if contextInfo.StanzaID != "" {
		ctx.StanzaID = proto.String(contextInfo.StanzaID)
	}

	if contextInfo.Participant != "" {
		ctx.Participant = proto.String(contextInfo.Participant)
	}

	return ctx
}
func (ms *Sender) DeleteMessage(ctx context.Context, sessionID string, phone string, messageID string) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(phone)
	if err != nil {
		return ErrInvalidJID
	}

	// Use BuildRevoke instead of deprecated RevokeMessage
	revokeMsg := client.WAClient.BuildRevoke(recipientJID, types.EmptyJID, messageID)
	_, err = client.WAClient.SendMessage(context.Background(), recipientJID, revokeMsg)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}
func (ms *Sender) EditMessage(ctx context.Context, sessionID string, phone string, messageID string, text string) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(phone)
	if err != nil {
		return ErrInvalidJID
	}

	_, err = client.WAClient.SendMessage(ctx, recipientJID, &waE2E.Message{
		EditedMessage: &waE2E.FutureProofMessage{
			Message: &waE2E.Message{
				Conversation: proto.String(text),
			},
		},
		ProtocolMessage: &waE2E.ProtocolMessage{
			Key: &waCommon.MessageKey{
				ID:        proto.String(messageID),
				RemoteJID: proto.String(recipientJID.String()),
			},
			Type:          waE2E.ProtocolMessage_MESSAGE_EDIT.Enum(),
			EditedMessage: &waE2E.Message{Conversation: proto.String(text)},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to edit message: %w", err)
	}

	return nil
}
func (ms *Sender) MarkRead(ctx context.Context, sessionID string, phone string, messageIDs []string) error {
	client, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	recipientJID, err := parseJID(phone)
	if err != nil {
		return ErrInvalidJID
	}

	err = client.WAClient.MarkRead(messageIDs, time.Now(), recipientJID, recipientJID)
	if err != nil {
		return fmt.Errorf("failed to mark as read: %w", err)
	}

	return nil
}
func (ms *Sender) RequestHistorySync(ctx context.Context, sessionID string, count int) error {
	_, err := ms.getConnectedClient(ctx, sessionID)
	if err != nil {
		return err
	}

	if count <= 0 {
		count = 50
	}

	// TODO: Implement history sync request when available in whatsmeow
	// For now, just validate the parameters
	_ = count

	return nil
}
