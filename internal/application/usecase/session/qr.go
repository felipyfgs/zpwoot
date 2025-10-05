package session

import (
	"context"
	"fmt"
	"time"

	"zpwoot/internal/application/dto"
	"zpwoot/internal/application/interfaces"
	"zpwoot/internal/domain/session"
	"zpwoot/internal/domain/shared"
)


type QRUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}


func NewQRUseCase(
	sessionService *session.Service,
	whatsappClient interfaces.WhatsAppClient,
	notificationSvc interfaces.NotificationService,
) *QRUseCase {
	return &QRUseCase{
		sessionService:  sessionService,
		whatsappClient:  whatsappClient,
		notificationSvc: notificationSvc,
	}
}


func (uc *QRUseCase) GetQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {

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


	if domainSession.IsConnected {
		return nil, fmt.Errorf("session is already connected")
	}


	qrInfo, err := uc.whatsappClient.GetQRCode(ctx, sessionID)
	if err != nil {
		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				return nil, dto.ErrSessionNotFound
			case "ALREADY_CONNECTED":
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
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusQRCode)
	}

	return dto.NewQRCodeResponse(
		qrInfo.Base64,
		qrInfo.ExpiresAt,
		string(domainSession.GetStatus()),
	), nil
}


func (uc *QRUseCase) RefreshQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {

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


	if domainSession.IsConnected {
		return nil, fmt.Errorf("session is already connected")
	}


	qrInfo, err := uc.whatsappClient.GetQRCode(ctx, sessionID)
	if err != nil {
		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				return nil, dto.ErrSessionNotFound
			case "ALREADY_CONNECTED":
				return nil, fmt.Errorf("session is already connected")
			default:
				return nil, fmt.Errorf("whatsapp QR refresh error: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to get QR code: %w", err)
	}


	if qrInfo.Code != "" {
		domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusQRCode)


		if uc.notificationSvc != nil {
			go func() {
				_ = uc.notificationSvc.NotifyQRCodeGenerated(ctx, sessionID, qrInfo.Code, qrInfo.ExpiresAt)
			}()
		}
	}

	return dto.NewQRCodeResponse(
		qrInfo.Base64,
		qrInfo.ExpiresAt,
		string(domainSession.GetStatus()),
	), nil
}


func (uc *QRUseCase) CheckQRCodeStatus(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {

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


	if domainSession.IsConnected {
		return dto.NewQRCodeResponse(
			"",
			time.Time{},
			string(session.StatusConnected),
		), nil
	}


	if domainSession.QRCode == "" {
		return dto.NewQRCodeResponse(
			"",
			time.Time{},
			string(session.StatusDisconnected),
		), nil
	}


	if domainSession.QRCodeExpiresAt != nil && time.Now().After(*domainSession.QRCodeExpiresAt) {

		domainSession.ClearQRCode()
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)

		return dto.NewQRCodeResponse(
			"",
			time.Time{},
			string(session.StatusDisconnected),
		), nil
	}


	expiresAt := time.Time{}
	if domainSession.QRCodeExpiresAt != nil {
		expiresAt = *domainSession.QRCodeExpiresAt
	}

	return dto.NewQRCodeResponse(
		domainSession.QRCode,
		expiresAt,
		string(domainSession.GetStatus()),
	), nil
}
