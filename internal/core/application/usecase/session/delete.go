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

type DeleteUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
}

func NewDeleteUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
) *DeleteUseCase {
	return &DeleteUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
	}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, sessionID string) error {
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

	if domainSession.IsConnected {
		_ = uc.whatsappClient.DisconnectSession(ctx, sessionID)
	}

	err = uc.whatsappClient.DeleteSession(ctx, sessionID)
	if err != nil {
		var waErr *output.WhatsAppError
		if errors.As(err, &waErr) {
			switch waErr.Code {
			case sessionNotFoundCode:
				break
			default:
				break
			}
		}
	}

	err = uc.sessionService.Delete(ctx, sessionID)
	if err != nil {
		if errors.Is(err, shared.ErrSessionNotFound) {
			return dto.ErrSessionNotFound
		}

		return fmt.Errorf("failed to delete session from domain: %w", err)
	}

	return nil
}

func (uc *DeleteUseCase) ExecuteForce(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	_ = uc.whatsappClient.DisconnectSession(ctx, sessionID)
	_ = uc.whatsappClient.DeleteSession(ctx, sessionID)

	err := uc.sessionService.Delete(ctx, sessionID)
	if err != nil && !errors.Is(err, shared.ErrSessionNotFound) {
		return fmt.Errorf("failed to delete session from domain: %w", err)
	}

	return nil
}

func (uc *DeleteUseCase) ExecuteWithValidation(ctx context.Context, sessionID string, force bool) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	_, err := uc.sessionService.Get(ctx, sessionID)
	if err != nil {
		if errors.Is(err, shared.ErrSessionNotFound) {
			return dto.ErrSessionNotFound
		}

		return fmt.Errorf("failed to validate session: %w", err)
	}

	if force {
		return uc.ExecuteForce(ctx, sessionID)
	}

	return uc.Execute(ctx, sessionID)
}
