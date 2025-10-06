package webhook

import (
	"context"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
)

// ListEventsUseCase implementa o caso de uso de listar eventos disponíveis
type ListEventsUseCase struct {
	webhookService *webhook.Service
}

// NewListEventsUseCase cria uma nova instância do use case
func NewListEventsUseCase(webhookService *webhook.Service) *ListEventsUseCase {
	return &ListEventsUseCase{
		webhookService: webhookService,
	}
}

// Execute executa o caso de uso
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

