package message

import (
	"context"
	"fmt"
	"time"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"
)

type ReceiveUseCase struct {
	sessionService *session.Service
}

func NewReceiveUseCase(
	sessionService *session.Service,
) *ReceiveUseCase {
	return &ReceiveUseCase{
		sessionService: sessionService,
	}
}

func (uc *ReceiveUseCase) ProcessIncomingMessage(ctx context.Context, req *dto.ReceiveMessageRequest) error {
	if req.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	if req.Message.ID == "" {
		return fmt.Errorf("message ID is required")
	}

	_, err := uc.sessionService.Get(ctx, req.SessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return dto.ErrSessionNotFound
		}

		return fmt.Errorf("failed to get session: %w", err)
	}

	go func(ctx context.Context) {
		if err := uc.sessionService.UpdateStatus(ctx, req.SessionID, session.StatusConnected); err != nil {
			fmt.Printf("Failed to update session status: %v\n", err)
		}
	}(ctx)

	return nil
}

func (uc *ReceiveUseCase) ProcessIncomingMessageBatch(ctx context.Context, sessionID string, messages []dto.MessageInfo) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	if len(messages) == 0 {
		return nil
	}

	_, err := uc.sessionService.Get(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return dto.ErrSessionNotFound
		}

		return fmt.Errorf("failed to get session: %w", err)
	}

	go func(ctx context.Context) {
		if err := uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusConnected); err != nil {
			fmt.Printf("Failed to update session status: %v\n", err)
		}
	}(ctx)

	for _, message := range messages {
		req := &dto.ReceiveMessageRequest{
			SessionID: sessionID,
			Message:   message,
		}

		go func(ctx context.Context, msgReq *dto.ReceiveMessageRequest) {
			if err := uc.ProcessIncomingMessage(ctx, msgReq); err != nil {
				fmt.Printf("Failed to process incoming message: %v\n", err)
			}
		}(ctx, req)
	}

	return nil
}

func (uc *ReceiveUseCase) ValidateMessage(message *dto.MessageInfo) error {
	if message.ID == "" {
		return fmt.Errorf("message ID is required")
	}

	if message.Chat == "" {
		return fmt.Errorf("chat JID is required")
	}

	if message.Sender == "" && !message.FromMe {
		return fmt.Errorf("sender is required for incoming messages")
	}

	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	if message.Type == "" {
		message.Type = "text"
	}

	return nil
}

func (uc *ReceiveUseCase) CreateMessageInfo(
	messageID, chatJID, senderJID, pushName, messageType, content string,
	fromMe, isGroup bool,
	timestamp time.Time,
) *dto.MessageInfo {
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	return &dto.MessageInfo{
		ID:        messageID,
		Chat:      chatJID,
		Sender:    senderJID,
		PushName:  pushName,
		Timestamp: timestamp,
		FromMe:    fromMe,
		Type:      messageType,
		IsGroup:   isGroup,
		Content:   content,
	}
}
