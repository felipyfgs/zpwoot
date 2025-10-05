package session

import (
	"context"
	"fmt"

	"zpwoot/internal/application/dto"
	"zpwoot/internal/application/interfaces"
	"zpwoot/internal/domain/session"
	"zpwoot/internal/domain/shared"
)

// DeleteUseCase handles session deletion
type DeleteUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}

// NewDeleteUseCase creates a new delete session use case
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

// Execute deletes a session completely
func (uc *DeleteUseCase) Execute(ctx context.Context, sessionID string) error {
	// Validate input
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	// Get session from domain layer to check if it exists
	domainSession, err := uc.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return dto.ErrSessionNotFound
		}
		return fmt.Errorf("failed to get session from domain: %w", err)
	}

	// First, disconnect the session if it's connected
	if domainSession.IsConnected {
		// Try to disconnect gracefully
		_ = uc.whatsappClient.DisconnectSession(ctx, sessionID)
	}

	// Delete from WhatsApp client
	err = uc.whatsappClient.DeleteSession(ctx, sessionID)
	if err != nil {
		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				// Session not found in WhatsApp client, continue with domain deletion
				break
			default:
				// Log error but continue with domain deletion
				// In production, this should be logged properly
				break
			}
		}
		// Don't fail the deletion if WhatsApp client fails
	}

	// Delete from domain layer
	err = uc.sessionService.DeleteSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return dto.ErrSessionNotFound
		}
		return fmt.Errorf("failed to delete session from domain: %w", err)
	}

	// Send notification
	if uc.notificationSvc != nil {
		go func() {
			_ = uc.notificationSvc.NotifySessionDisconnected(ctx, sessionID)
		}()
	}

	return nil
}

// ExecuteForce forcefully deletes a session (cleanup mode)
func (uc *DeleteUseCase) ExecuteForce(ctx context.Context, sessionID string) error {
	// Validate input
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	// Force disconnect and delete from WhatsApp client (ignore all errors)
	_ = uc.whatsappClient.DisconnectSession(ctx, sessionID)
	_ = uc.whatsappClient.DeleteSession(ctx, sessionID)

	// Delete from domain layer (ignore not found errors)
	err := uc.sessionService.DeleteSession(ctx, sessionID)
	if err != nil && err != shared.ErrSessionNotFound {
		return fmt.Errorf("failed to delete session from domain: %w", err)
	}

	// Send notification
	if uc.notificationSvc != nil {
		go func() {
			_ = uc.notificationSvc.NotifySessionDisconnected(ctx, sessionID)
		}()
	}

	return nil
}

// ExecuteWithValidation deletes a session with additional validation
func (uc *DeleteUseCase) ExecuteWithValidation(ctx context.Context, sessionID string, force bool) error {
	// Validate input
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	// Check if session exists
	_, err := uc.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		if err == shared.ErrSessionNotFound {
			return dto.ErrSessionNotFound
		}
		return fmt.Errorf("failed to validate session: %w", err)
	}

	// Use appropriate deletion method
	if force {
		return uc.ExecuteForce(ctx, sessionID)
	}

	return uc.Execute(ctx, sessionID)
}
