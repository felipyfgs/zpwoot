package session

import (
	"context"
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

		if waErr, ok := err.(*output.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_ALREADY_EXISTS":
				return nil, dto.ErrSessionAlreadyExists
			default:
				return nil, fmt.Errorf("whatsapp client error: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to create WhatsApp session: %w", err)
	}

	// If QR code generation is requested, connect the session and wait for QR code
	if req.GenerateQRCode {
		// Connect the session to generate QR code
		err = uc.whatsappClient.ConnectSession(ctx, sessionID)
		if err != nil {
			// Don't fail the creation, just log the error
			// The session is created but QR code generation failed
		} else {
			// Wait a bit for QR code to be generated via events, then try to get it
			time.Sleep(2 * time.Second)

			// Try to get the QR code (it should be available now via events)
			qrInfo, qrErr := uc.whatsappClient.GetQRCode(ctx, sessionID)
			if qrErr == nil && qrInfo.Code != "" {
				domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)
				_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusQRCode)
			}
		}
	}

	response := dto.SessionToCreateResponse(domainSession)

	return response, nil
}
