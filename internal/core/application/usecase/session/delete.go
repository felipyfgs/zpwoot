package session

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/application/interfaces"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"
)


type DeleteUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}


func NewDeleteUseCase(
	sessionService *session.Service,
	whatsappClient interfaces.WhatsAppClient,
	notificationSvc interfaces.NotificationService,
) *DeleteUseCase {
	return &DeleteUseCase{
		sessionService:  sessionService,
		whatsappClient:  whatsappClient,
		notificationSvc: notificationSvc,
	}
}


func (uc *DeleteUseCase) Execute(ctx context.Context, sessionID string) error {

	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}


	domainSession, err := uc.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return dto.ErrSessionNotFound
		}
		return fmt.Errorf("failed to get session from domain: %w", err)
	}


	if domainSession.IsConnected {

		_ = uc.whatsappClient.DisconnectSession(ctx, sessionID)
	}


	err = uc.whatsappClient.DeleteSession(ctx, sessionID)
	if err != nil {
		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":

				break
			default:


				break
			}
		}

	}


	err = uc.sessionService.DeleteSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return dto.ErrSessionNotFound
		}
		return fmt.Errorf("failed to delete session from domain: %w", err)
	}


	if uc.notificationSvc != nil {
		go func() {
			_ = uc.notificationSvc.NotifySessionDisconnected(ctx, sessionID)
		}()
	}

	return nil
}


func (uc *DeleteUseCase) ExecuteForce(ctx context.Context, sessionID string) error {

	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}


	_ = uc.whatsappClient.DisconnectSession(ctx, sessionID)
	_ = uc.whatsappClient.DeleteSession(ctx, sessionID)


	err := uc.sessionService.DeleteSession(ctx, sessionID)
	if err != nil && err != shared.ErrSessionNotFound {
		return fmt.Errorf("failed to delete session from domain: %w", err)
	}


	if uc.notificationSvc != nil {
		go func() {
			_ = uc.notificationSvc.NotifySessionDisconnected(ctx, sessionID)
		}()
	}

	return nil
}


func (uc *DeleteUseCase) ExecuteWithValidation(ctx context.Context, sessionID string, force bool) error {

	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}


	_, err := uc.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return dto.ErrSessionNotFound
		}
		return fmt.Errorf("failed to validate session: %w", err)
	}


	if force {
		return uc.ExecuteForce(ctx, sessionID)
	}

	return uc.Execute(ctx, sessionID)
}
