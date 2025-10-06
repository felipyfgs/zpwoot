package session

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"
	"zpwoot/internal/core/ports/output"
)

type LogoutUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
}

func NewLogoutUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
) *LogoutUseCase {
	return &LogoutUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
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

	if !domainSession.IsConnected && domainSession.DeviceJID == "" {
		return fmt.Errorf("session is already logged out")
	}

	err = uc.whatsappClient.LogoutSession(ctx, sessionID)
	if err != nil {
		if waErr, ok := err.(*output.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":

				break
			case "ALREADY_LOGGED_OUT":

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

	return nil
}
