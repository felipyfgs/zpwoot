package session

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/application/interfaces"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"
)


type DisconnectUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}


func NewDisconnectUseCase(
	sessionService *session.Service,
	whatsappClient interfaces.WhatsAppClient,
	notificationSvc interfaces.NotificationService,
) *DisconnectUseCase {
	return &DisconnectUseCase{
		sessionService:  sessionService,
		whatsappClient:  whatsappClient,
		notificationSvc: notificationSvc,
	}
}


func (uc *DisconnectUseCase) Execute(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {

	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}


	domainSession, err := uc.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return nil, dto.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session from domain: %w", err)
	}


	if !domainSession.IsConnected {
		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(domainSession.GetStatus()),
			Connected: false,
		}, nil
	}


	err = uc.whatsappClient.DisconnectSession(ctx, sessionID)
	if err != nil {
		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":

				domainSession.SetDisconnected()
				_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)
				return &dto.SessionStatusResponse{
					ID:        sessionID,
					Status:    string(session.StatusDisconnected),
					Connected: false,
				}, nil
			case "ALREADY_DISCONNECTED":

				domainSession.SetDisconnected()
				_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)
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
	err = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)
	if err != nil {

	}


	if uc.notificationSvc != nil {
		go func() {
			_ = uc.notificationSvc.NotifySessionDisconnected(ctx, sessionID)
		}()
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


	domainSession, err := uc.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return nil, dto.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session from domain: %w", err)
	}


	_ = uc.whatsappClient.DisconnectSession(ctx, sessionID)


	domainSession.SetDisconnected()
	_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)


	if uc.notificationSvc != nil {
		go func() {
			_ = uc.notificationSvc.NotifySessionDisconnected(ctx, sessionID)
		}()
	}

	return &dto.SessionStatusResponse{
		ID:        sessionID,
		Status:    string(session.StatusDisconnected),
		Connected: false,
	}, nil
}
