package session

import (
	"context"
	"fmt"

	"zpwoot/internal/application/dto"
	"zpwoot/internal/application/interfaces"
	"zpwoot/internal/domain/session"
)

// ListUseCase handles listing sessions
type ListUseCase struct {
	sessionService *session.Service
	whatsappClient interfaces.WhatsAppClient
}

// NewListUseCase creates a new list sessions use case
func NewListUseCase(
	sessionService *session.Service,
	whatsappClient interfaces.WhatsAppClient,
) *ListUseCase {
	return &ListUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
	}
}

// Execute retrieves a paginated list of sessions
func (uc *ListUseCase) Execute(ctx context.Context, pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	// Validate and set defaults for pagination
	if pagination == nil {
		pagination = &dto.PaginationRequest{Limit: 20, Offset: 0}
	}
	pagination.Validate()

	// Get sessions from domain layer
	domainSessions, err := uc.sessionService.ListSessions(ctx, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions from domain: %w", err)
	}

	// Convert to response DTOs and sync with WhatsApp status
	sessionResponses := make([]dto.SessionResponse, len(domainSessions))
	for i, domainSession := range domainSessions {
		// Convert to DTO
		sessionResponse := dto.SessionToListResponse(domainSession)

		// Get current WhatsApp status for each session
		if uc.whatsappClient != nil {
			waStatus, err := uc.whatsappClient.GetSessionStatus(ctx, domainSession.ID)
			if err == nil && waStatus != nil {
				sessionResponse.Connected = waStatus.Connected
				sessionResponse.DeviceJID = waStatus.DeviceJID
				if waStatus.Connected {
					sessionResponse.Status = "connected"
				} else if waStatus.LoggedIn {
					sessionResponse.Status = "disconnected"
				} else {
					sessionResponse.Status = "qr_code"
				}
				if !waStatus.ConnectedAt.IsZero() {
					sessionResponse.ConnectedAt = &waStatus.ConnectedAt
				}
				if !waStatus.LastSeen.IsZero() {
					sessionResponse.LastSeen = &waStatus.LastSeen
				}
			}
		}

		sessionResponses[i] = *sessionResponse
	}

	// Calculate total count (this could be optimized with a separate count query)
	totalCount := len(sessionResponses)
	hasMore := len(sessionResponses) == pagination.Limit

	// Create paginated response
	response := &dto.PaginationResponse{
		Items:   sessionResponses,
		Total:   totalCount,
		Limit:   pagination.Limit,
		Offset:  pagination.Offset,
		HasMore: hasMore,
	}

	return response, nil
}

// ExecuteSimple retrieves all sessions without pagination
func (uc *ListUseCase) ExecuteSimple(ctx context.Context) ([]*dto.SessionListResponse, error) {
	// Get all sessions (with a reasonable limit)
	pagination := &dto.PaginationRequest{Limit: 100, Offset: 0}

	result, err := uc.Execute(ctx, pagination)
	if err != nil {
		return nil, err
	}

	// Extract sessions from pagination response
	sessions, ok := result.Items.([]*dto.SessionListResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}

	return sessions, nil
}

// ExecuteWithFilter retrieves sessions with filtering (future enhancement)
func (uc *ListUseCase) ExecuteWithFilter(ctx context.Context, filter *SessionFilter, pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	// For now, just call the basic Execute method
	// In the future, this could support filtering by status, name, etc.
	return uc.Execute(ctx, pagination)
}

// SessionFilter represents filtering options for sessions
type SessionFilter struct {
	Status    string `json:"status,omitempty"`
	Connected *bool  `json:"connected,omitempty"`
	Name      string `json:"name,omitempty"`
}
