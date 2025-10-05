package input

import (
	"context"
	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/output"
)

// MessageSender defines the interface for sending different types of messages
type MessageSender interface {
	SendTextMessage(ctx context.Context, sessionID, to, text string) error
	SendMediaMessage(ctx context.Context, sessionID, to string, media *dto.MediaData) error
	SendLocationMessage(ctx context.Context, sessionID, to string, latitude, longitude float64, name string) error
	SendContactMessage(ctx context.Context, sessionID, to string, contact *dto.ContactInfo) error
}

// MessageReceiver defines the interface for receiving messages
type MessageReceiver interface {
	Execute(ctx context.Context, req *dto.ReceiveMessageRequest) error
}

// ChatInfoGetter defines the interface for getting chat information
type ChatInfoGetter interface {
	GetChatInfo(ctx context.Context, sessionID, chatJID string) (interface{}, error)
}

// ContactsGetter defines the interface for getting contacts
type ContactsGetter interface {
	GetContacts(ctx context.Context, sessionID string) (interface{}, error)
}

// ChatsGetter defines the interface for getting chats
type ChatsGetter interface {
	GetChats(ctx context.Context, sessionID string) (interface{}, error)
}

// MessageUseCases combines all message-related use case interfaces
type MessageUseCases interface {
	MessageSender
	MessageReceiver
	ChatInfoGetter
	ContactsGetter
	ChatsGetter
}

// MessageService defines the interface for message operations at the adapter level
type MessageService interface {
	SendTextMessage(ctx context.Context, sessionID, to, text string) error
	SendMediaMessage(ctx context.Context, sessionID, to string, media *output.MediaData) error
	SendLocationMessage(ctx context.Context, sessionID, to string, latitude, longitude float64, name string) error
	SendContactMessage(ctx context.Context, sessionID, to string, contact *ContactInfo) error
	GetChatInfo(ctx context.Context, sessionID, chatJID string) (*ChatInfo, error)
	GetContacts(ctx context.Context, sessionID string) ([]*ContactInfo, error)
	GetChats(ctx context.Context, sessionID string) ([]*ChatInfo, error)
}

// ContactInfo represents contact information for input layer
type ContactInfo struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	VCard string `json:"vcard,omitempty"`
}

// ChatInfo represents chat information for input layer
type ChatInfo struct {
	JID              string `json:"jid"`
	Name             string `json:"name,omitempty"`
	Topic            string `json:"topic,omitempty"`
	IsGroup          bool   `json:"isGroup"`
	ParticipantCount int    `json:"participantCount,omitempty"`
}
