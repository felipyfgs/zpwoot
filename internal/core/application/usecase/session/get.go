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

type GetUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
}

func NewGetUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
) *GetUseCase {
	return &GetUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
	}
}

func (uc *GetUseCase) Execute(ctx context.Context, sessionID string) (*dto.SessionDetailResponse, error) {
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

	waStatus, err := uc.whatsappClient.GetSessionStatus(ctx, sessionID)
	if err != nil {
		if waErr, ok := err.(*output.WhatsAppError); ok && waErr.Code == "SESSION_NOT_FOUND" {
			return dto.ToDetailResponse(domainSession), nil
		}

		return dto.ToDetailResponse(domainSession), nil
	}

	if waStatus != nil {
		if waStatus.Connected && !domainSession.IsConnected {
			domainSession.SetConnected(waStatus.DeviceJID)
		} else if !waStatus.Connected && domainSession.IsConnected {
			domainSession.SetDisconnected()
		}

		if waStatus.DeviceJID != "" {
			domainSession.DeviceJID = waStatus.DeviceJID
		}

		if !waStatus.LastSeen.IsZero() {
			domainSession.UpdateLastSeen()
		}

		go func(ctx context.Context) {
			if waStatus.Connected {
				_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusConnected)
			} else {
				_ = uc.sessionService.UpdateStatus(ctx, sessionID, session.StatusDisconnected)
			}
		}(ctx)
	}

	response := dto.ToDetailResponse(domainSession)

	return response, nil
}

func (uc *GetUseCase) ExecuteWithSync(ctx context.Context, sessionID string) (*dto.SessionDetailResponse, error) {
	response, err := uc.Execute(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	waStatus, err := uc.whatsappClient.GetSessionStatus(ctx, sessionID)
	if err != nil {
		return response, nil
	}

	if waStatus != nil {
		response.DeviceJID = waStatus.DeviceJID
		response.Connected = waStatus.Connected

		if waStatus.Connected {
			response.Status = "connected"
		} else if waStatus.LoggedIn {
			response.Status = "disconnected"
		} else {
			response.Status = "qr_code"
		}

		if !waStatus.ConnectedAt.IsZero() {
			response.ConnectedAt = &waStatus.ConnectedAt
		}

		if !waStatus.LastSeen.IsZero() {
			response.LastSeen = &waStatus.LastSeen
		}
	}

	return response, nil
}
