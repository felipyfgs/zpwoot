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
	if err := uc.validateSessionID(sessionID); err != nil {
		return nil, err
	}

	domainSession, err := uc.getDomainSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if response := uc.checkAlreadyDisconnected(sessionID, domainSession); response != nil {
		return response, nil
	}

	if response, err := uc.performWhatsAppDisconnection(ctx, sessionID, domainSession); err != nil || response != nil {
		return response, err
	}

	uc.finalizeDisconnection(ctx, sessionID, domainSession)

	return uc.buildDisconnectedResponse(sessionID), nil
}

func (uc *DisconnectUseCase) validateSessionID(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	return nil
}

func (uc *DisconnectUseCase) getDomainSession(ctx context.Context, sessionID string) (*session.Session, error) {
	domainSession, err := uc.sessionService.Get(ctx, sessionID)
	if err != nil {
		if errors.Is(err, shared.ErrSessionNotFound) {
			return nil, dto.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session from domain: %w", err)
	}
	return domainSession, nil
}

func (uc *DisconnectUseCase) checkAlreadyDisconnected(sessionID string, domainSession *session.Session) *dto.SessionStatusResponse {
	if !domainSession.IsConnected {
		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(domainSession.GetStatus()),
			Connected: false,
			Message:   "Session is already disconnected",
		}
	}
	return nil
}

func (uc *DisconnectUseCase) performWhatsAppDisconnection(ctx context.Context, sessionID string, domainSession *session.Session) (*dto.SessionStatusResponse, error) {
	err := uc.whatsappClient.DisconnectSession(ctx, sessionID)
	if err != nil {
		return uc.handleWhatsAppDisconnectionError(ctx, sessionID, domainSession, err)
	}
	return nil, fmt.Errorf("disconnection completed successfully")
}

func (uc *DisconnectUseCase) handleWhatsAppDisconnectionError(ctx context.Context, sessionID string, domainSession *session.Session, err error) (*dto.SessionStatusResponse, error) {
	var waErr *output.WhatsAppError
	if errors.As(err, &waErr) {
		switch waErr.Code {
		case sessionNotFoundCode, "ALREADY_DISCONNECTED":
			return uc.handleSessionNotFoundOrAlreadyDisconnected(ctx, sessionID, domainSession), nil
		default:
			return nil, fmt.Errorf("whatsapp disconnection error: %w", err)
		}
	}
	return nil, fmt.Errorf("failed to disconnect WhatsApp session: %w", err)
}

func (uc *DisconnectUseCase) handleSessionNotFoundOrAlreadyDisconnected(ctx context.Context, sessionID string, domainSession *session.Session) *dto.SessionStatusResponse {
	domainSession.SetDisconnected()
	_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusDisconnected)

	return &dto.SessionStatusResponse{
		ID:        sessionID,
		Status:    string(session.StatusDisconnected),
		Connected: false,
	}
}

func (uc *DisconnectUseCase) finalizeDisconnection(ctx context.Context, sessionID string, domainSession *session.Session) {
	domainSession.SetDisconnected()

	if err := uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusDisconnected); err != nil {
		fmt.Printf("Failed to update session status to disconnected: %v\n", err)
	}
}

func (uc *DisconnectUseCase) buildDisconnectedResponse(sessionID string) *dto.SessionStatusResponse {
	return &dto.SessionStatusResponse{
		ID:        sessionID,
		Status:    string(session.StatusDisconnected),
		Connected: false,
	}
}

func (uc *DisconnectUseCase) ExecuteForce(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
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

	_ = uc.whatsappClient.DisconnectSession(ctx, sessionID)

	domainSession.SetDisconnected()

	_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusDisconnected)

	return &dto.SessionStatusResponse{
		ID:        sessionID,
		Status:    string(session.StatusDisconnected),
		Connected: false,
	}, nil
}
