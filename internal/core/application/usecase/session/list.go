package session

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/application/interfaces"
	"zpwoot/internal/core/domain/session"
)


type ListUseCase struct {
	sessionService *session.Service
	whatsappClient interfaces.WhatsAppClient
}


func NewListUseCase(
	sessionService *session.Service,
	whatsappClient interfaces.WhatsAppClient,
) *ListUseCase {
	return &ListUseCase{
		sessionService: sessionService,
		whatsappClient: whatsappClient,
	}
}


func (uc *ListUseCase) Execute(ctx context.Context, pagination *dto.PaginationRequest) (*dto.PaginationResponse, error) {

	if pagination == nil {
		pagination = &dto.PaginationRequest{Limit: 20, Offset: 0}
	}
	pagination.Validate()


	domainSessions, err := uc.sessionService.ListSessions(ctx, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions from domain: %w", err)
	}


	sessionResponses := make([]dto.SessionResponse, len(domainSessions))
	for i, domainSession := range domainSessions {

		sessionResponse := dto.SessionToListResponse(domainSession)


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
