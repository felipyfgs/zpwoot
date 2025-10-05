package message

import (
	"context"
	"fmt"

	"zpwoot/internal/application/dto"
	"zpwoot/internal/application/interfaces"
	"zpwoot/internal/domain/session"
	"zpwoot/internal/domain/shared"

	"github.com/google/uuid"
)


type SendUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}


func NewSendUseCase(
	sessionService *session.Service,
	whatsappClient interfaces.WhatsAppClient,
	notificationSvc interfaces.NotificationService,
) *SendUseCase {
	return &SendUseCase{
		sessionService:  sessionService,
		whatsappClient:  whatsappClient,
		notificationSvc: notificationSvc,
	}
}


func (uc *SendUseCase) Execute(ctx context.Context, sessionID string, req *dto.SendMessageRequest) (*dto.SendMessageResponse, error) {

	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}


	domainSession, err := uc.sessionService.GetSession(ctx, sessionID)
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


	var messageResult *interfaces.MessageResult
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
		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
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


	go func() {
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusConnected)
	}()


	if uc.notificationSvc != nil {
		go func() {
			_ = uc.notificationSvc.NotifyMessageSent(ctx, sessionID, messageResult.MessageID)
		}()
	}


	finalMessageID := messageID
	if messageResult.MessageID != "" {
		finalMessageID = messageResult.MessageID
	}

	return &dto.SendMessageResponse{
		MessageID: finalMessageID,
		Status:    messageResult.Status,
		SentAt:    messageResult.SentAt,
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
