package message

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	domainSession "zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"
)

type UseCases struct {
	send    *SendUseCase
	receive *ReceiveUseCase
}

func NewUseCases(sessionService *domainSession.Service, whatsappClient output.WhatsAppClient, logger output.Logger) input.MessageUseCases {
	return &UseCases{
		send:    NewSendUseCase(sessionService, whatsappClient, logger),
		receive: NewReceiveUseCase(sessionService, logger),
	}
}

func (m *UseCases) SendTextMessage(ctx context.Context, sessionID, to, text string) error {
	req := &dto.SendMessageRequest{
		To:   to,
		Type: "text",
		Text: text,
	}
	_, err := m.send.Execute(ctx, sessionID, req)

	return err
}

func (m *UseCases) SendMediaMessage(ctx context.Context, sessionID, to string, media *dto.MediaData) error {
	req := &dto.SendMessageRequest{
		To:    to,
		Type:  "media",
		Media: media,
	}
	_, err := m.send.Execute(ctx, sessionID, req)

	return err
}

func (m *UseCases) SendLocationMessage(ctx context.Context, sessionID, to string, latitude, longitude float64, name string) error {
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

func (m *UseCases) SendContactMessage(ctx context.Context, sessionID, to string, contact *dto.ContactInfo) error {
	req := &dto.SendMessageRequest{
		To:      to,
		Type:    "contact",
		Contact: contact,
	}
	_, err := m.send.Execute(ctx, sessionID, req)

	return err
}

func (m *UseCases) Execute(ctx context.Context, req *dto.ReceiveMessageRequest) error {
	return m.receive.ProcessIncomingMessage(ctx, req)
}

func (m *UseCases) GetChatInfo(ctx context.Context, sessionID, chatJID string) (interface{}, error) {
	return nil, fmt.Errorf("GetChatInfo not implemented yet")
}

func (m *UseCases) GetContacts(ctx context.Context, sessionID string) (interface{}, error) {
	return nil, fmt.Errorf("GetContacts not implemented yet")
}

func (m *UseCases) GetChats(ctx context.Context, sessionID string) (interface{}, error) {
	return nil, fmt.Errorf("GetChats not implemented yet")
}
