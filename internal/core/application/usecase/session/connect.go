package session

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/application/interfaces"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"
)


type ConnectUseCase struct {
	sessionService  *session.Service
	whatsappClient  interfaces.WhatsAppClient
	notificationSvc interfaces.NotificationService
}


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


func (uc *ConnectUseCase) Execute(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {

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
		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(domainSession.GetStatus()),
			Connected: true,
		}, nil
	}


	err = uc.whatsappClient.ConnectSession(ctx, sessionID)
	if err != nil {

		domainSession.SetError(err.Error())
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusError)

		if waErr, ok := err.(*interfaces.WhatsAppError); ok {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				return nil, dto.ErrSessionNotFound
			case "ALREADY_CONNECTED":

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


	err = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusConnecting)
	if err != nil {

	}


	waStatus, err := uc.whatsappClient.GetSessionStatus(ctx, sessionID)
	if err != nil {

		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(session.StatusConnecting),
			Connected: false,
		}, nil
	}


	if waStatus.Connected {
		domainSession.SetConnected(waStatus.DeviceJID)
		_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusConnected)


		if uc.notificationSvc != nil {
			go func() {
				_ = uc.notificationSvc.NotifySessionConnected(ctx, sessionID, waStatus.DeviceJID)
			}()
		}
	} else if !waStatus.LoggedIn {

		qrInfo, err := uc.whatsappClient.GetQRCode(ctx, sessionID)
		if err == nil && qrInfo.Code != "" {
			domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)
			_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusQRCode)


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


func (uc *ConnectUseCase) ExecuteWithAutoReconnect(ctx context.Context, sessionID string, autoReconnect bool) (*dto.SessionStatusResponse, error) {


	return uc.Execute(ctx, sessionID)
}
