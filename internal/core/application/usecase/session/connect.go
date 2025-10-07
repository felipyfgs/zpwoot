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

type ConnectUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
	logger         output.Logger
}

func NewConnectUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
	logger output.Logger,
) *ConnectUseCase {
	return &ConnectUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
		logger:         logger,
	}
}

func (uc *ConnectUseCase) Execute(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
	domainSession, err := uc.validateSessionForConnection(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if response := uc.checkExistingConnection(sessionID, domainSession); response != nil {
		return response, nil
	}

	if response, err := uc.performWhatsAppConnection(ctx, sessionID, domainSession); err != nil {
		return response, err
	} else if response != nil {

		return response, nil
	}

	uc.updateConnectionStatus(ctx, sessionID)

	time.Sleep(1 * time.Second)

	updatedSession, err := uc.sessionService.Get(ctx, sessionID)
	if err == nil && updatedSession != nil {
		domainSession = updatedSession
	}

	return uc.buildConnectionResponse(ctx, sessionID, domainSession)
}

func (uc *ConnectUseCase) ExecuteWithAutoReconnect(ctx context.Context, sessionID string, autoReconnect bool) (*dto.SessionStatusResponse, error) {
	return uc.Execute(ctx, sessionID)
}

func (uc *ConnectUseCase) validateSessionForConnection(ctx context.Context, sessionID string) (*session.Session, error) {
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

	return domainSession, nil
}

func (uc *ConnectUseCase) checkExistingConnection(sessionID string, domainSession *session.Session) *dto.SessionStatusResponse {
	if domainSession.IsConnected {
		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(domainSession.GetStatus()),
			Connected: true,
			Message:   "Session is already connected",
		}
	}
	return nil
}

func (uc *ConnectUseCase) performWhatsAppConnection(ctx context.Context, sessionID string, domainSession *session.Session) (*dto.SessionStatusResponse, error) {
	err := uc.whatsappClient.ConnectSession(ctx, sessionID)
	if err != nil {
		domainSession.SetError(err.Error())

		if updateErr := uc.sessionService.Update(ctx, domainSession); updateErr != nil {
			uc.logger.Error().Err(updateErr).Str("session_id", sessionID).Msg("Failed to update session status")
		}

		var waErr *output.WhatsAppError
		if errors.As(err, &waErr) {
			switch waErr.Code {
			case "SESSION_NOT_FOUND":
				return nil, dto.ErrSessionNotFound
			case "ALREADY_CONNECTED":
				domainSession.SetConnected(domainSession.DeviceJID)
				_ = uc.sessionService.Update(ctx, domainSession)

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

	return nil, nil
}

func (uc *ConnectUseCase) updateConnectionStatus(ctx context.Context, sessionID string) {
	err := uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusConnecting)
	if err != nil {
		uc.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to update session status to connecting")
	}
}

func (uc *ConnectUseCase) buildConnectionResponse(ctx context.Context, sessionID string, domainSession *session.Session) (*dto.SessionStatusResponse, error) {
	waStatus, err := uc.whatsappClient.GetSessionStatus(ctx, sessionID)
	if err != nil {
		return &dto.SessionStatusResponse{
			ID:        sessionID,
			Status:    string(session.StatusConnecting),
			Connected: false,
		}, err
	}

	if waStatus.Connected {
		domainSession.SetConnected(waStatus.DeviceJID)
		_ = uc.sessionService.Update(ctx, domainSession)
	} else if !waStatus.LoggedIn {
		qrInfo, err := uc.whatsappClient.GetQRCode(ctx, sessionID)
		if err == nil && qrInfo.Code != "" {
			domainSession.SetQRCode(qrInfo.Code, qrInfo.ExpiresAt)
			_ = uc.sessionService.Update(ctx, domainSession)
		}
	}

	return dto.ToStatusResponse(domainSession), nil
}
