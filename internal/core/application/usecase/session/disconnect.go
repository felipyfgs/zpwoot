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

type DisconnectUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
}

func NewDisconnectUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
) *DisconnectUseCase {
	return &DisconnectUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
	}
}

func (uc *DisconnectUseCase) Execute(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	domainSession, err := uc.sessionService.Get(ctx, sessionID)
	if err != nil {
		if errors.Is(err, shared.ErrSessionNotFound) {
			return nil, dto.ErrSessionNotFound
		}

		return nil, fmt.Errorf("failed to get session from domain: %w", err)
	}

	if !domainSession.IsConnected {
		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(domainSession.GetStatus()),
			Connected: false,
			Message:   "Session is already disconnected",
		}, nil
	}

	err = uc.whatsappClient.DisconnectSession(ctx, sessionID)
	if err != nil {
		var waErr *output.WhatsAppError
		if errors.As(err, &waErr) {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				domainSession.SetDisconnected()

				_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusDisconnected)

				return &dto.SessionStatusResponse{
					ID:        sessionID,
					Status:    string(session.StatusDisconnected),
					Connected: false,
				}, nil
			case "ALREADY_DISCONNECTED":
				domainSession.SetDisconnected()

				_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusDisconnected)

				return &dto.SessionStatusResponse{
					ID:        sessionID,
					Status:    string(session.StatusDisconnected),
					Connected: false,
				}, nil
			default:
				return nil, fmt.Errorf("whatsapp disconnection error: %w", err)
			}
		}

		return nil, fmt.Errorf("failed to disconnect WhatsApp session: %w", err)
	}

	domainSession.SetDisconnected()

	err = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusDisconnected)
	if err != nil {

	}

	return &dto.SessionStatusResponse{
		ID:        sessionID,
		Status:    string(session.StatusDisconnected),
		Connected: false,
	}, nil
}

func (uc *DisconnectUseCase) ExecuteForce(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	domainSession, err := uc.sessionService.Get(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return nil, dto.ErrSessionNotFound
		}

		return nil, fmt.Errorf("failed to get session from domain: %w", err)
	}

	_ = uc.whatsappClient.DisconnectSession(ctx, sessionID)

	domainSession.SetDisconnected()

	_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusDisconnected)

	return &dto.SessionStatusResponse{
		ID:        sessionID,
		Status:    string(session.StatusDisconnected),
		Connected: false,
	}, nil
}
