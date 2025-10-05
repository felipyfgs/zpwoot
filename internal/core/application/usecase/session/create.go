package session

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/application/interfaces"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"
)


type CreateUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}


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


func (uc *CreateUseCase) Execute(ctx context.Context, req *dto.CreateSessionRequest) (*dto.CreateSessionResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}


	domainSession, err := uc.sessionService.CreateSession(ctx, req.Name)
	if err != nil {
		if err == shared.ErrSessionAlreadyExists {
			return nil, dto.ErrSessionAlreadyExists
		}
		return nil, fmt.Errorf("failed to create session in domain: %w", err)
	}

	sessionID := domainSession.ID


	err = uc.whatsappClient.CreateSession(ctx, sessionID)
	if err != nil {

		if rollbackErr := uc.sessionService.DeleteSession(ctx, sessionID); rollbackErr != nil {

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


	if uc.notificationSvc != nil {
		go func() {
			_ = uc.notificationSvc.NotifySessionConnected(ctx, sessionID, "")
		}()
	}


	response := dto.SessionToCreateResponse(domainSession)

	return response, nil
}
