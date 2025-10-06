package message

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"
	"zpwoot/internal/core/ports/output"

	"github.com/google/uuid"
)

type SendUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
}

func NewSendUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
) *SendUseCase {
	return &SendUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
	}
}

func (uc *SendUseCase) Execute(ctx context.Context, sessionID string, req *dto.SendMessageRequest) (*dto.SendMessageResponse, error) {

	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	domainSession, err := uc.sessionService.Get(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return nil, dto.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if !domainSession.IsConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	messageID := uuid.New().String()

	var messageResult *output.MessageResult
	switch req.Type {
	case "text":
		messageResult, err = uc.whatsappClient.SendTextMessage(ctx, sessionID, req.To, req.Text)
	case "media":
		messageResult, err = uc.whatsappClient.SendMediaMessage(ctx, sessionID, req.To, req.Media.ToInterfacesMediaData())
	case "location":
		messageResult, err = uc.whatsappClient.SendLocationMessage(ctx, sessionID, req.To, req.Location.ToInterfacesLocation())
	case "contact":
		messageResult, err = uc.whatsappClient.SendContactMessage(ctx, sessionID, req.To, req.Contact.ToInterfacesContactInfo())
	default:
		return nil, fmt.Errorf("unsupported message type: %s", req.Type)
	}

	if err != nil {
		if waErr, ok := err.(*output.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				return nil, dto.ErrSessionNotFound
			case "SESSION_NOT_CONNECTED":
				return nil, fmt.Errorf("session is not connected")
			case "INVALID_JID":
				return nil, fmt.Errorf("invalid recipient")
			default:
				return nil, fmt.Errorf("whatsapp send error: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	go func(ctx context.Context) {

		if err := uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusConnected); err != nil {

			fmt.Printf("Failed to update session status: %v\n", err)
		}
	}(ctx)

	finalMessageID := messageID
	if messageResult.MessageID != "" {
		finalMessageID = messageResult.MessageID
	}

	return &dto.SendMessageResponse{
		Success:   true,
		ID:        finalMessageID,
		To:        "",
		Type:      "",
		Timestamp: messageResult.SentAt.Unix(),
		Status:    messageResult.Status,
	}, nil
}

func (uc *SendUseCase) SendText(ctx context.Context, sessionID, to, text string) (*dto.SendMessageResponse, error) {
	req := &dto.SendMessageRequest{
		To:   to,
		Type: "text",
		Text: text,
	}
	return uc.Execute(ctx, sessionID, req)
}

func (uc *SendUseCase) SendMedia(ctx context.Context, sessionID, to string, media *dto.MediaData) (*dto.SendMessageResponse, error) {
	req := &dto.SendMessageRequest{
		To:    to,
		Type:  "media",
		Media: media,
	}
	return uc.Execute(ctx, sessionID, req)
}

func (uc *SendUseCase) SendLocation(ctx context.Context, sessionID, to string, location *dto.Location) (*dto.SendMessageResponse, error) {
	req := &dto.SendMessageRequest{
		To:       to,
		Type:     "location",
		Location: location,
	}
	return uc.Execute(ctx, sessionID, req)
}

func (uc *SendUseCase) SendContact(ctx context.Context, sessionID, to string, contact *dto.ContactInfo) (*dto.SendMessageResponse, error) {
	req := &dto.SendMessageRequest{
		To:      to,
		Type:    "contact",
		Contact: contact,
	}
	return uc.Execute(ctx, sessionID, req)
}
