package message

import (
	"context"
	"errors"
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
	logger         output.Logger
}

func NewSendUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
	logger output.Logger,
) *SendUseCase {
	return &SendUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
		logger:         logger,
	}
}

func (uc *SendUseCase) Execute(ctx context.Context, sessionID string, req *dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	if err := uc.validateSendRequest(sessionID, req); err != nil {
		return nil, err
	}

	_, err := uc.getValidatedSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	messageResult, err := uc.sendMessageByType(ctx, sessionID, req)
	if err != nil {
		return nil, err
	}

	uc.updateSessionStatusAsync(ctx, sessionID)

	return uc.buildSendResponse(messageResult), nil
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

func (uc *SendUseCase) validateSendRequest(sessionID string, req *dto.SendMessageRequest) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func (uc *SendUseCase) getValidatedSession(ctx context.Context, sessionID string) (*session.Session, error) {
	domainSession, err := uc.sessionService.Get(ctx, sessionID)
	if err != nil {
		if errors.Is(err, shared.ErrSessionNotFound) {
			return nil, dto.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if !domainSession.IsConnected {
		return nil, fmt.Errorf("session is not connected")
	}

	return domainSession, nil
}

func (uc *SendUseCase) sendMessageByType(ctx context.Context, sessionID string, req *dto.SendMessageRequest) (*output.MessageResult, error) {
	var messageResult *output.MessageResult
	var err error

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
		return nil, uc.handleWhatsAppError(err)
	}

	return messageResult, nil
}

func (uc *SendUseCase) handleWhatsAppError(err error) error {
	var waErr *output.WhatsAppError
	if errors.As(err, &waErr) {
		switch waErr.Code {
		case "SESSION_NOT_FOUND":
			return dto.ErrSessionNotFound
		case "SESSION_NOT_CONNECTED":
			return fmt.Errorf("session is not connected")
		case "INVALID_JID":
			return fmt.Errorf("invalid recipient")
		default:
			return fmt.Errorf("whatsapp send error: %w", err)
		}
	}

	return fmt.Errorf("failed to send message: %w", err)
}

func (uc *SendUseCase) updateSessionStatusAsync(ctx context.Context, sessionID string) {
	go func(ctx context.Context) {
		if err := uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusConnected); err != nil {
			uc.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to update session status")
		}
	}(ctx)
}

func (uc *SendUseCase) buildSendResponse(messageResult *output.MessageResult) *dto.SendMessageResponse {
	messageID := uuid.New().String()
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
	}
}
