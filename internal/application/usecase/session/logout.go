package session

import (
	"context"
	"fmt"

	"zpwoot/internal/application/dto"
	"zpwoot/internal/application/interfaces"
	"zpwoot/internal/domain/session"
	"zpwoot/internal/domain/shared"
)


type LogoutUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}


func NewLogoutUseCase(
	sessionService *session.Service,
	whatsappClient interfaces.WhatsAppClient,
	notificationSvc interfaces.NotificationService,
) *LogoutUseCase {
	return &LogoutUseCase{
		sessionService:  sessionService,
		whatsappClient:  whatsappClient,
		notificationSvc: notificationSvc,
	}
}


func (uc *LogoutUseCase) Execute(ctx context.Context, sessionID string) error {

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


	err = uc.whatsappClient.LogoutSession(ctx, sessionID)
	if err != nil {
		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":

				break
			default:
				return fmt.Errorf("whatsapp logout error: %w", err)
			}
		} else {
			return fmt.Errorf("failed to logout WhatsApp session: %w", err)
		}
	}


	domainSession.SetDisconnected()
	domainSession.DeviceJID = ""
	domainSession.QRCode = ""
	domainSession.QRCodeExpiresAt = nil


	err = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)
	if err != nil {


	}


	if uc.notificationSvc != nil {
		go func() {
			_ = uc.notificationSvc.NotifySessionDisconnected(ctx, sessionID)
		}()
	}

	return nil
}

