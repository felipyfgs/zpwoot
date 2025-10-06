package input

import (
	"context"
	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/output"
)

type MessageSender interface {
	SendTextMessage(ctx context.Context, sessionID, to, text string) error
	SendMediaMessage(ctx context.Context, sessionID, to string, media *dto.MediaData) error
	SendLocationMessage(ctx context.Context, sessionID, to string, latitude, longitude float64, name string) error
	SendContactMessage(ctx context.Context, sessionID, to string, contact *dto.ContactInfo) error
}

type MessageReceiver interface {
	Execute(ctx context.Context, req *dto.ReceiveMessageRequest) error
}

type ChatInfoGetter interface {
	GetChatInfo(ctx context.Context, sessionID, chatJID string) (interface{}, error)
}

type ContactsGetter interface {
	GetContacts(ctx context.Context, sessionID string) (interface{}, error)
}

type ChatsGetter interface {
	GetChats(ctx context.Context, sessionID string) (interface{}, error)
}

type MessageUseCases interface {
	MessageSender
	MessageReceiver
	ChatInfoGetter
	ContactsGetter
	ChatsGetter
}

type MessageService interface {
	SendTextMessage(ctx context.Context, sessionID, to, text string, contextInfo *output.MessageContextInfo) (*output.MessageResult, error)
	SendMediaMessage(ctx context.Context, sessionID, to string, media *output.MediaData, contextInfo *output.MessageContextInfo) (*output.MessageResult, error)
	SendLocationMessage(ctx context.Context, sessionID, to string, latitude, longitude float64, name string, contextInfo *output.MessageContextInfo) (*output.MessageResult, error)
	SendContactMessage(ctx context.Context, sessionID, to string, contact *ContactInfo, contextInfo *output.MessageContextInfo) (*output.MessageResult, error)
	SendReactionMessage(ctx context.Context, sessionID, to, messageID, reaction string, fromMe bool) (*output.MessageResult, error)
	SendPollMessage(ctx context.Context, sessionID, to, name string, options []string, selectableCount int) (*output.MessageResult, error)
	SendButtonsMessage(ctx context.Context, sessionID, to, text string, buttons []ButtonInfo) (*output.MessageResult, error)
	SendListMessage(ctx context.Context, sessionID, to, text, title string, sections []ListSectionInfo) (*output.MessageResult, error)
	SendTemplateMessage(ctx context.Context, sessionID, to string, template TemplateInfo) (*output.MessageResult, error)
	SendViewOnceMessage(ctx context.Context, sessionID, to string, media *output.MediaData) (*output.MessageResult, error)
	GetChatInfo(ctx context.Context, sessionID, chatJID string) (*ChatInfo, error)
	GetContacts(ctx context.Context, sessionID string) ([]*ContactInfo, error)
	GetChats(ctx context.Context, sessionID string) ([]*ChatInfo, error)
}

type ButtonInfo struct {
	ID   string
	Text string
}

type ListSectionInfo struct {
	Title string
	Rows  []ListRowInfo
}

type ListRowInfo struct {
	ID          string
	Title       string
	Description string
}

type TemplateInfo struct {
	Content string
	Footer  string
}

type ContactInfo struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	VCard string `json:"vcard,omitempty"`
}

type ChatInfo struct {
	JID              string `json:"jid"`
	Name             string `json:"name,omitempty"`
	Topic            string `json:"topic,omitempty"`
	IsGroup          bool   `json:"isGroup"`
	ParticipantCount int    `json:"participantCount,omitempty"`
}
