package message

import (
	"context"
	"fmt"
	"time"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/application/interfaces"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"
)


type ReceiveUseCase struct {
	sessionService  *session.Service
	notificationSvc interfaces.NotificationService
}


func NewReceiveUseCase(
	sessionService *session.Service,
	notificationSvc interfaces.NotificationService,
) *ReceiveUseCase {
	return &ReceiveUseCase{
		sessionService:  sessionService,
		notificationSvc: notificationSvc,
	}
}


func (uc *ReceiveUseCase) ProcessIncomingMessage(ctx context.Context, req *dto.ReceiveMessageRequest) error {

	if req.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	if req.Message.ID == "" {
		return fmt.Errorf("message ID is required")
	}


	_, err := uc.sessionService.GetSession(ctx, req.SessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return dto.ErrSessionNotFound
		}
		return fmt.Errorf("failed to get session: %w", err)
	}


	go func() {
		_ = uc.sessionService.UpdateSessionStatus(ctx, req.SessionID, session.StatusConnected)
	}()


	if uc.notificationSvc != nil {
		messageEvent := &interfaces.MessageEvent{
			ID:        req.Message.ID,
			Chat:      req.Message.Chat,
			Sender:    req.Message.Sender,
			PushName:  req.Message.PushName,
			Timestamp: req.Message.Timestamp,
			FromMe:    req.Message.FromMe,
			Type:      req.Message.Type,
			IsGroup:   req.Message.IsGroup,
			Content:   req.Message.Content,
		}

		go func() {
			_ = uc.notificationSvc.NotifyMessageReceived(ctx, req.SessionID, messageEvent)
		}()
	}

	return nil
}


func (uc *ReceiveUseCase) ProcessIncomingMessageBatch(ctx context.Context, sessionID string, messages []dto.MessageInfo) error {

	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	if len(messages) == 0 {
		return nil
	}


	_, err := uc.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return dto.ErrSessionNotFound
		}
		return fmt.Errorf("failed to get session: %w", err)
	}


	go func() {
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusConnected)
	}()


	for _, message := range messages {
		req := &dto.ReceiveMessageRequest{
			SessionID: sessionID,
			Message:   message,
		}


		go func(msgReq *dto.ReceiveMessageRequest) {
			_ = uc.ProcessIncomingMessage(ctx, msgReq)
		}(req)
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
