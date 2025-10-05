package session

import (
	"context"
	"fmt"

	"zpwoot/internal/application/dto"
	"zpwoot/internal/application/interfaces"
	"zpwoot/internal/domain/session"
	"zpwoot/internal/domain/shared"
)

// GetUseCase handles getting session details
type GetUseCase struct {
	sessionService *session.Service
	whatsappClient interfaces.WhatsAppClient
}

// NewGetUseCase creates a new get session use case
func NewGetUseCase(
	sessionService *session.Service,
	whatsappClient interfaces.WhatsAppClient,
) *GetUseCase {
	return &GetUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
	}
}

// Execute retrieves session details by ID
func (uc *GetUseCase) Execute(ctx context.Context, sessionID string) (*dto.SessionDetailResponse, error) {
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

	// Get current WhatsApp session status to sync
	waStatus, err := uc.whatsappClient.GetSessionStatus(ctx, sessionID)
	if err != nil {
		// If WhatsApp session not found, just return domain session
		if waErr, ok := err.(*interfaces.WhatsAppError); ok && waErr.Code == "SESSION_NOT_FOUND" {
			return dto.SessionToDetailResponse(domainSession), nil
		}
		// For other errors, log but don't fail - return domain session
		return dto.SessionToDetailResponse(domainSession), nil
	}

	// Sync domain session with WhatsApp session status
	if waStatus != nil {
		// Update connection status
		if waStatus.Connected && !domainSession.IsConnected {
			domainSession.SetConnected(waStatus.DeviceJID)
		} else if !waStatus.Connected && domainSession.IsConnected {
			domainSession.SetDisconnected()
		}

		// Update device JID
		if waStatus.DeviceJID != "" {
			domainSession.DeviceJID = waStatus.DeviceJID
		}

		// Update last seen
		if !waStatus.LastSeen.IsZero() {
			domainSession.UpdateLastSeen()
		}

		// Update session status in domain (fire and forget)
		go func() {
			if waStatus.Connected {
				_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusConnected)
			} else {
				_ = uc.sessionService.UpdateSessionStatus(ctx, sessionID, session.StatusDisconnected)
			}
		}()
	}

	// Convert to response DTO
	response := dto.SessionToDetailResponse(domainSession)

	return response, nil
}

// ExecuteWithSync retrieves session details and forces sync with WhatsApp
func (uc *GetUseCase) ExecuteWithSync(ctx context.Context, sessionID string) (*dto.SessionDetailResponse, error) {
	// Get session details
	response, err := uc.Execute(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Force sync with WhatsApp client
	waStatus, err := uc.whatsappClient.GetSessionStatus(ctx, sessionID)
	if err != nil {
		// Return response even if sync fails
		return response, nil
	}

	// Update response with latest WhatsApp info
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
