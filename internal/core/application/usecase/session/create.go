package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"
	"zpwoot/internal/core/ports/output"
)

type CreateUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
}

func NewCreateUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
) *CreateUseCase {
	return &CreateUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
	}
}

func (uc *CreateUseCase) Execute(ctx context.Context, req *dto.CreateRequest) (*dto.CreateSessionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	domainSession, err := uc.sessionService.Create(ctx, req.Name)
	if err != nil {
		if errors.Is(err, shared.ErrSessionAlreadyExists) {
			return nil, dto.ErrSessionAlreadyExists
		}

		return nil, fmt.Errorf("failed to create session in domain: %w", err)
	}

	sessionID := domainSession.ID

	err = uc.whatsappClient.CreateSession(ctx, sessionID)
	if err != nil {
		if rollbackErr := uc.sessionService.Delete(ctx, sessionID); rollbackErr != nil {

		}

		var waErr *output.WhatsAppError
		if errors.As(err, &waErr) {
			switch waErr.Code {
			case "SESSION_ALREADY_EXISTS":
				return nil, dto.ErrSessionAlreadyExists
			default:
				return nil, fmt.Errorf("whatsapp client error: %w", err)
			}
		}

		return nil, fmt.Errorf("failed to create WhatsApp session: %w", err)
	}

	if req.GenerateQRCode {
		err = uc.whatsappClient.ConnectSession(ctx, sessionID)
		if err != nil {

		} else {
			time.Sleep(2 * time.Second)

			qrInfo, qrErr := uc.whatsappClient.GetQRCode(ctx, sessionID)
			if qrErr == nil && qrInfo.Code != "" {
				domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)

				_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusQRCode)
			}
		}
	}

	response := dto.ToCreateResponse(domainSession)

	return response, nil
}
