package webhook

import (
	"context"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
)


type ListEventsUseCase struct {
	webhookService *webhook.Service
}


func NewListEventsUseCase(webhookService *webhook.Service) *ListEventsUseCase {
	return &ListEventsUseCase{
		webhookService: webhookService,
	}
}


func (uc *ListEventsUseCase) Execute(ctx context.Context) (*dto.ListEventsResponse, error) {
	categories := uc.webhookService.GetEventCategories()
	allEvents := uc.webhookService.GetValidEventTypes()

	var categoryResponses []dto.EventCategoryResponse
	for category, events := range categories {
		categoryResponses = append(categoryResponses, dto.EventCategoryResponse{
			Category: category,
			Events:   events,
		})
	}

	return &dto.ListEventsResponse{
		Categories: categoryResponses,
		AllEvents:  allEvents,
	}, nil
}
