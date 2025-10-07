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

type ConnectUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
}

func NewConnectUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
) *ConnectUseCase {
	return &ConnectUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
	}
}

func (uc *ConnectUseCase) Execute(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
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
		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(domainSession.GetStatus()),
			Connected: true,
			Message:   "Session is already connected",
		}, nil
	}

	err = uc.whatsappClient.ConnectSession(ctx, sessionID)
	if err != nil {
		domainSession.SetError(err.Error())

		if updateErr := uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusError); updateErr != nil {
			fmt.Printf("Failed to update session status: %v\n", updateErr)
		}

		var waErr *output.WhatsAppError
		if errors.As(err, &waErr) {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				return nil, dto.ErrSessionNotFound
			case "ALREADY_CONNECTED":
				domainSession.SetConnected(domainSession.DeviceJID)

				_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusConnected)

				return &dto.SessionStatusResponse{
					ID:        sessionID,
					Status:    string(session.StatusConnected),
					Connected: true,
					Message:   "Session was already connected at WhatsApp level",
				}, nil
			default:
				return nil, fmt.Errorf("whatsapp connection error: %w", err)
			}
		}

		return nil, fmt.Errorf("failed to connect WhatsApp session: %w", err)
	}

	err = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusConnecting)
	if err != nil {

		fmt.Printf("Failed to update session status to connecting: %v\n", err)
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

		_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusConnected)
	} else if !waStatus.LoggedIn {
		qrInfo, err := uc.whatsappClient.GetQRCode(ctx, sessionID)
		if err == nil && qrInfo.Code != "" {
			domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)

			_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusQRCode)
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
