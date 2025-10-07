package session

import (
	"context"
	"errors"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"
	"zpwoot/internal/core/ports/output"
)

type LogoutUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
	logger         output.Logger
}

func NewLogoutUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
	logger output.Logger,
) *LogoutUseCase {
	return &LogoutUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
		logger:         logger,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	domainSession, err := uc.sessionService.Get(ctx, sessionID)
	if err != nil {
		if errors.Is(err, shared.ErrSessionNotFound) {
			return dto.ErrSessionNotFound
		}

		return fmt.Errorf("failed to get session from domain: %w", err)
	}

	if !domainSession.IsConnected && domainSession.DeviceJID == "" {
		return fmt.Errorf("session is already logged out")
	}

	err = uc.whatsappClient.LogoutSession(ctx, sessionID)
	if err != nil {
		var waErr *output.WhatsAppError
		if errors.As(err, &waErr) {
			switch waErr.Code {
			case sessionNotFoundCode:
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

	err = uc.sessionService.Update(ctx, domainSession)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}
