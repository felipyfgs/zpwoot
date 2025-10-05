package message

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	domainSession "zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"
)

// MessageUseCases implements the MessageUseCases interface
type MessageUseCases struct {
	send    *SendUseCase
	receive *ReceiveUseCase
}

// NewMessageUseCases creates a new MessageUseCases implementation
func NewMessageUseCases(sessionService *domainSession.Service, whatsappClient output.WhatsAppClient) input.MessageUseCases {
	return &MessageUseCases{
		send:    NewSendUseCase(sessionService, whatsappClient),
		receive: NewReceiveUseCase(sessionService),
	}
}

// SendTextMessage sends a text message
func (m *MessageUseCases) SendTextMessage(ctx context.Context, sessionID, to, text string) error {
	req := &dto.SendMessageRequest{
		To:   to,
		Type: "text",
		Text: text,
	}
	_, err := m.send.Execute(ctx, sessionID, req)
	return err
}

// SendMediaMessage sends a media message
func (m *MessageUseCases) SendMediaMessage(ctx context.Context, sessionID, to string, media *dto.MediaData) error {
	req := &dto.SendMessageRequest{
		To:    to,
		Type:  "media",
		Media: media,
	}
	_, err := m.send.Execute(ctx, sessionID, req)
	return err
}

// SendLocationMessage sends a location message
func (m *MessageUseCases) SendLocationMessage(ctx context.Context, sessionID, to string, latitude, longitude float64, name string) error {
	req := &dto.SendMessageRequest{
		To:   to,
		Type: "location",
		Location: &dto.Location{
			Latitude:  latitude,
			Longitude: longitude,
			Name:      name,
		},
	}
	_, err := m.send.Execute(ctx, sessionID, req)
	return err
}

// SendContactMessage sends a contact message
func (m *MessageUseCases) SendContactMessage(ctx context.Context, sessionID, to string, contact *dto.ContactInfo) error {
	req := &dto.SendMessageRequest{
		To:      to,
		Type:    "contact",
		Contact: contact,
	}
	_, err := m.send.Execute(ctx, sessionID, req)
	return err
}

// Execute processes an incoming message
func (m *MessageUseCases) Execute(ctx context.Context, req *dto.ReceiveMessageRequest) error {
	return m.receive.ProcessIncomingMessage(ctx, req)
}

// GetChatInfo retrieves chat information (placeholder)
func (m *MessageUseCases) GetChatInfo(ctx context.Context, sessionID, chatJID string) (interface{}, error) {
	return nil, fmt.Errorf("GetChatInfo not implemented yet")
}

// GetContacts retrieves contacts (placeholder)
func (m *MessageUseCases) GetContacts(ctx context.Context, sessionID string) (interface{}, error) {
	return nil, fmt.Errorf("GetContacts not implemented yet")
}

// GetChats retrieves chats (placeholder)
func (m *MessageUseCases) GetChats(ctx context.Context, sessionID string) (interface{}, error) {
	return nil, fmt.Errorf("GetChats not implemented yet")
}
