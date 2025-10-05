package session

import (
	"context"
	"fmt"

	"zpwoot/internal/application/dto"
	"zpwoot/internal/application/interfaces"
	"zpwoot/internal/domain/session"
	"zpwoot/internal/domain/shared"
)

// DisconnectUseCase handles session disconnection
type DisconnectUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}

// NewDisconnectUseCase creates a new disconnect session use case
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

// Execute disconnects a session from WhatsApp
func (uc *DisconnectUseCase) Execute(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
	// Validate input
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	// Get session from domain layer
	domainSession, err := uc.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return nil, dto.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session from domain: %w", err)
	}

	// Check if already disconnected
	if !domainSession.IsConnected {
		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(domainSession.GetStatus()),
			Connected: false,
		}, nil
	}

	// Disconnect via WhatsApp client
	err = uc.whatsappClient.DisconnectSession(ctx, sessionID)
	if err != nil {
		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				// Session not found in WhatsApp client, update domain to disconnected
				domainSession.SetDisconnected()
				_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)
				return &dto.SessionStatusResponse{
					ID:        sessionID,
					Status:    string(session.StatusDisconnected),
					Connected: false,
				}, nil
			case "ALREADY_DISCONNECTED":
				// Already disconnected in WhatsApp client, sync domain
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

	// Update domain session to disconnected
	domainSession.SetDisconnected()
	err = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)
	if err != nil {
		// Log error but don't fail the disconnection process
	}

	// Send notification
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

// ExecuteForce forcefully disconnects a session (cleanup)
func (uc *DisconnectUseCase) ExecuteForce(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
	// Validate input
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	// Get session from domain layer
	domainSession, err := uc.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return nil, dto.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session from domain: %w", err)
	}

	// Force disconnect via WhatsApp client (ignore errors)
	_ = uc.whatsappClient.DisconnectSession(ctx, sessionID)

	// Always update domain session to disconnected
	domainSession.SetDisconnected()
	_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)

	// Send notification
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
