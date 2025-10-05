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

// QRUseCase handles QR code operations
type QRUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}

// NewQRUseCase creates a new QR code use case
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

// GetQRCode retrieves the current QR code for a session
func (uc *QRUseCase) GetQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {
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

	// Check if session is already connected
	if domainSession.IsConnected {
		return nil, fmt.Errorf("session is already connected")
	}

	// Get QR code from WhatsApp client
	qrInfo, err := uc.whatsappClient.GetQRCode(ctx, sessionID)
	if err != nil {
		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				return nil, dto.ErrSessionNotFound
			case "ALREADY_CONNECTED":
				return nil, fmt.Errorf("session is already connected")
			case "QR_CODE_EXPIRED":
				// Try to refresh QR code
				return uc.RefreshQRCode(ctx, sessionID)
			default:
				return nil, fmt.Errorf("whatsapp QR code error: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to get QR code: %w", err)
	}

	// Update domain session with QR code
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

// RefreshQRCode generates a new QR code for a session
func (uc *QRUseCase) RefreshQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {
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

	// Check if session is already connected
	if domainSession.IsConnected {
		return nil, fmt.Errorf("session is already connected")
	}

	// Get fresh QR code via WhatsApp client
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

	// Update domain session with new QR code
	if qrInfo.Code != "" {
		domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusQRCode)

		// Send notification
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

// CheckQRCodeStatus checks if QR code is still valid
func (uc *QRUseCase) CheckQRCodeStatus(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {
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

	// Check if session is already connected
	if domainSession.IsConnected {
		return dto.NewQRCodeResponse(
			"",
			time.Time{},
			string(session.StatusConnected),
		), nil
	}

	// Check if QR code exists and is valid
	if domainSession.QRCode == "" {
		return dto.NewQRCodeResponse(
			"",
			time.Time{},
			string(session.StatusDisconnected),
		), nil
	}

	// Check if QR code is expired
	if domainSession.QRCodeExpiresAt != nil && time.Now().After(*domainSession.QRCodeExpiresAt) {
		// Clear expired QR code
		domainSession.ClearQRCode()
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)

		return dto.NewQRCodeResponse(
			"",
			time.Time{},
			string(session.StatusDisconnected),
		), nil
	}

	// Return current QR code
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
