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

type CreateUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
	logger         output.Logger
}

func NewCreateUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
	logger output.Logger,
) *CreateUseCase {
	return &CreateUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
		logger:         logger,
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
			uc.logger.Error().Err(rollbackErr).Str("session_id", sessionID).Msg("Failed to rollback session creation")
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
		// Connect and get QR code directly using WhatsApp Meow's QR channel
		qrInfo, err := uc.whatsappClient.ConnectAndGetQRCode(ctx, sessionID)
		if err != nil {
			uc.logger.Warn().Err(err).Str("session_id", sessionID).Msg("Failed to connect and get QR code")
		} else if qrInfo != nil && qrInfo.Code != "" {
			// Set QR code in domain session
			domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)

			// Save to database
			if updateErr := uc.sessionService.Update(ctx, domainSession); updateErr != nil {
				uc.logger.Error().Err(updateErr).Str("session_id", sessionID).Msg("Failed to update session with QR code")
			} else {
				uc.logger.Info().Str("session_id", sessionID).Msg("QR code successfully generated and stored")
			}
		}
	}

	response := dto.ToCreateResponse(domainSession)

	return response, nil
}
