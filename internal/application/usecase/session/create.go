package session

import (
	"context"
	"fmt"

	"zpwoot/internal/application/dto"
	"zpwoot/internal/application/interfaces"
	"zpwoot/internal/domain/session"
	"zpwoot/internal/domain/shared"
)

// CreateUseCase handles session creation
type CreateUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}

// NewCreateUseCase creates a new create session use case
func NewCreateUseCase(
	sessionService *session.Service,
	whatsappClient interfaces.WhatsAppClient,
	notificationSvc interfaces.NotificationService,
) *CreateUseCase {
	return &CreateUseCase{
		sessionService:  sessionService,
		whatsappClient:  whatsappClient,
		notificationSvc: notificationSvc,
	}
}

// Execute creates a new WhatsApp session
func (uc *CreateUseCase) Execute(ctx context.Context, req *dto.CreateSessionRequest) (*dto.CreateSessionResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create domain session first
	domainSession, err := uc.sessionService.CreateSession(ctx, req.Name)
	if err != nil {
		if err == shared.ErrSessionAlreadyExists {
			return nil, dto.ErrSessionAlreadyExists
		}
		return nil, fmt.Errorf("failed to create session in domain: %w", err)
	}

	sessionID := domainSession.ID

	// Create WhatsApp client session
	err = uc.whatsappClient.CreateSession(ctx, sessionID)
	if err != nil {
		// Rollback domain session creation
		if rollbackErr := uc.sessionService.DeleteSession(ctx, sessionID); rollbackErr != nil {
			// Log rollback error but don't override original error
		}

		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_ALREADY_EXISTS":
				return nil, dto.ErrSessionAlreadyExists
			default:
				return nil, fmt.Errorf("whatsapp client error: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to create WhatsApp session: %w", err)
	}

	// Send notification if service is available
	if uc.notificationSvc != nil {
		go func() {
			_ = uc.notificationSvc.NotifySessionConnected(ctx, sessionID, "")
		}()
	}

	// Convert to response DTO
	response := dto.SessionToCreateResponse(domainSession)

	return response, nil
}
