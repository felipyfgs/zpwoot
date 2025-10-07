package session

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/ports/output"
)

type ListUseCase struct {
	sessionService *session.Service
	whatsappClient output.WhatsAppClient
	logger         output.Logger
}

func NewListUseCase(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
	logger output.Logger,
) *ListUseCase {
	return &ListUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
		logger:         logger,
	}
}

func (uc *ListUseCase) Execute(ctx context.Context, pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	if pagination == nil {
		pagination = &dto.PaginationRequest{Limit: 20, Offset: 0}
	}

	_ = pagination.Validate()

	domainSessions, err := uc.sessionService.List(ctx, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions from domain: %w", err)
	}

	sessionResponses := make([]dto.SessionListInfo, len(domainSessions))

	for i, domainSession := range domainSessions {
		sessionResponse := dto.ToListInfo(domainSession)

		if uc.whatsappClient != nil {
			uc.updateSessionResponseFromWAClient(ctx, domainSession.ID, sessionResponse)
		}

		sessionResponses[i] = *sessionResponse
	}

	totalCount := len(sessionResponses)
	hasMore := len(sessionResponses) == pagination.Limit

	response := &dto.PaginationResponse{
		Items:   sessionResponses,
		Total:   totalCount,
		Limit:   pagination.Limit,
		Offset:  pagination.Offset,
		HasMore: hasMore,
	}

	return response, nil
}

func (uc *ListUseCase) ExecuteSimple(ctx context.Context) ([]*dto.SessionListResponse, error) {
	pagination := &dto.PaginationRequest{Limit: 100, Offset: 0}

	result, err := uc.Execute(ctx, pagination)
	if err != nil {
		return nil, err
	}

	sessions, ok := result.Items.([]*dto.SessionListResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}

	return sessions, nil
}

func (uc *ListUseCase) ExecuteWithFilter(ctx context.Context, filter *SessionFilter, pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	return uc.Execute(ctx, pagination)
}

type SessionFilter struct {
	Status    string `json:"status,omitempty"`
	Connected *bool  `json:"connected,omitempty"`
	Name      string `json:"name,omitempty"`
}

func (uc *ListUseCase) updateSessionResponseFromWAClient(ctx context.Context, sessionID string, sessionResponse *dto.SessionListInfo) {
	waStatus, err := uc.whatsappClient.GetSessionStatus(ctx, sessionID)
	if err != nil || waStatus == nil {
		return
	}

	sessionResponse.Connected = waStatus.Connected
	sessionResponse.DeviceJID = waStatus.DeviceJID

	switch {
	case waStatus.Connected:
		sessionResponse.Status = "connected"
	case waStatus.LoggedIn:
		sessionResponse.Status = "disconnected"
	default:
		sessionResponse.Status = "qr_code"
	}

	if !waStatus.ConnectedAt.IsZero() {
		sessionResponse.ConnectedAt = &waStatus.ConnectedAt
	}

	if !waStatus.LastSeen.IsZero() {
		sessionResponse.LastSeen = &waStatus.LastSeen
	}
}
