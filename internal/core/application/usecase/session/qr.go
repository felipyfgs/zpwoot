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

type QRUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
	logger         output.Logger
}

func NewQRUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
	logger output.Logger,
) *QRUseCase {
	return &QRUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
		logger:         logger,
	}
}

func (uc *QRUseCase) GetQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {
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

	if domainSession.IsConnected {
		return nil, fmt.Errorf("session is already connected")
	}

	qrInfo, err := uc.whatsappClient.GetQRCode(ctx, sessionID)
	if err != nil {
		var waErr *output.WhatsAppError
		if errors.As(err, &waErr) {
			switch waErr.Code {
			case sessionNotFoundCode:
				return nil, dto.ErrSessionNotFound
			case alreadyConnectedCode:
				return nil, fmt.Errorf("session is already connected")
			case "QR_CODE_EXPIRED":
				return uc.RefreshQRCode(ctx, sessionID)
			default:
				return nil, fmt.Errorf("whatsapp QR code error: %w", err)
			}
		}

		return nil, fmt.Errorf("failed to get QR code: %w", err)
	}

	if qrInfo.Code != "" {
		domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)

		_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusQRCode)
	}

	return dto.NewQRResponse(
		qrInfo.Base64,
		qrInfo.ExpiresAt,
		string(domainSession.GetStatus()),
	), nil
}

func (uc *QRUseCase) RefreshQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {
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

	if domainSession.IsConnected {
		return nil, fmt.Errorf("session is already connected")
	}

	qrInfo, err := uc.whatsappClient.GetQRCode(ctx, sessionID)
	if err != nil {
		var waErr *output.WhatsAppError
		if errors.As(err, &waErr) {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				return nil, dto.ErrSessionNotFound
			case alreadyConnectedCode:
				return nil, fmt.Errorf("session is already connected")
			default:
				return nil, fmt.Errorf("whatsapp QR refresh error: %w", err)
			}
		}

		return nil, fmt.Errorf("failed to get QR code: %w", err)
	}

	if qrInfo.Code != "" {
		domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)

		_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusQRCode)
	}

	return dto.NewQRResponse(
		qrInfo.Base64,
		qrInfo.ExpiresAt,
		string(domainSession.GetStatus()),
	), nil
}

func (uc *QRUseCase) CheckQRCodeStatus(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {
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

	if domainSession.IsConnected {
		return dto.NewQRResponse(
			"",
			time.Time{},
			string(session.StatusConnected),
		), nil
	}

	if domainSession.QRCode == "" {
		return dto.NewQRResponse(
			"",
			time.Time{},
			string(session.StatusDisconnected),
		), nil
	}

	if domainSession.QRCodeExpiresAt != nil && time.Now().After(*domainSession.QRCodeExpiresAt) {
		domainSession.ClearQRCode()

		_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusDisconnected)

		return dto.NewQRResponse(
			"",
			time.Time{},
			string(session.StatusDisconnected),
		), nil
	}

	expiresAt := time.Time{}
	if domainSession.QRCodeExpiresAt != nil {
		expiresAt = *domainSession.QRCodeExpiresAt
	}

	return dto.NewQRResponse(
		domainSession.QRCode,
		expiresAt,
		string(domainSession.GetStatus()),
	), nil
}
