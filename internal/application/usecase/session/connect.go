package session

import (
	"context"
	"fmt"

	"zpwoot/internal/application/dto"
	"zpwoot/internal/application/interfaces"
	"zpwoot/internal/domain/session"
	"zpwoot/internal/domain/shared"
)

// ConnectUseCase handles session connection
type ConnectUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}

// NewConnectUseCase creates a new connect session use case
func NewConnectUseCase(
	sessionService *session.Service,
	whatsappClient interfaces.WhatsAppClient,
	notificationSvc interfaces.NotificationService,
) *ConnectUseCase {
	return &ConnectUseCase{
		sessionService:  sessionService,
		whatsappClient:  whatsappClient,
		notificationSvc: notificationSvc,
	}
}

// Execute connects a session to WhatsApp
func (uc *ConnectUseCase) Execute(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
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

	// Check if already connected
	if domainSession.IsConnected {
		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(domainSession.GetStatus()),
			Connected: true,
		}, nil
	}

	// Connect via WhatsApp client
	err = uc.whatsappClient.ConnectSession(ctx, sessionID)
	if err != nil {
		// Update session with error status
		domainSession.SetError(err.Error())
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusError)

		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				return nil, dto.ErrSessionNotFound
			case "ALREADY_CONNECTED":
				// Update domain session to reflect connected state
				domainSession.SetConnected(domainSession.DeviceJID)
				_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusConnected)
				return &dto.SessionStatusResponse{
					ID:        sessionID,
					Status:    string(session.StatusConnected),
					Connected: true,
				}, nil
			default:
				return nil, fmt.Errorf("whatsapp connection error: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to connect WhatsApp session: %w", err)
	}

	// Update session status to connecting
	err = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusConnecting)
	if err != nil {
		// Log error but don't fail the connection process
	}

	// Get updated session status from WhatsApp client
	waStatus, err := uc.whatsappClient.GetSessionStatus(ctx, sessionID)
	if err != nil {
		// Connection might be in progress, return connecting status
		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(session.StatusConnecting),
			Connected: false,
		}, nil
	}

	// Update domain session with WhatsApp status
	if waStatus.Connected {
		domainSession.SetConnected(waStatus.DeviceJID)
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusConnected)

		// Send notification
		if uc.notificationSvc != nil {
			go func() {
				_ = uc.notificationSvc.NotifySessionConnected(ctx, sessionID, waStatus.DeviceJID)
			}()
		}
	} else if !waStatus.LoggedIn {
		// Try to get QR code if not logged in
		qrInfo, err := uc.whatsappClient.GetQRCode(ctx, sessionID)
		if err == nil && qrInfo.Code != "" {
			domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)
			_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusQRCode)

			// Send QR code notification
			if uc.notificationSvc != nil {
				go func() {
					_ = uc.notificationSvc.NotifyQRCodeGenerated(ctx, sessionID, qrInfo.Code, qrInfo.ExpiresAt)
				}()
			}
		}
	}

	return &dto.SessionStatusResponse{
		ID:        sessionID,
		Status:    string(domainSession.GetStatus()),
		Connected: domainSession.IsConnected,
	}, nil
}

// ExecuteWithAutoReconnect connects a session with auto-reconnect enabled
func (uc *ConnectUseCase) ExecuteWithAutoReconnect(ctx context.Context, sessionID string, autoReconnect bool) (*dto.SessionStatusResponse, error) {
	// For now, just call the basic Execute method
	// In the future, this could configure auto-reconnect settings
	return uc.Execute(ctx, sessionID)
}
