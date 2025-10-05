package input

import (
	"context"
	"zpwoot/internal/core/application/dto"
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
