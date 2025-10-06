package input

import (
	"context"

	"zpwoot/internal/core/application/dto"
)

// WebhookUseCases define a interface para os casos de uso de webhook
// REGRA: Apenas interfaces, sem implementações
type WebhookUseCases interface {
	// Create cria um novo webhook para uma sessão
	Create(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error)

	// Update atualiza um webhook existente
	Update(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error)

	// Upsert cria ou atualiza um webhook para uma sessão
	Upsert(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error)

	// Get busca o webhook configurado para uma sessão
	Get(ctx context.Context, sessionID string) (*dto.WebhookResponse, error)

	// Delete remove o webhook de uma sessão
	Delete(ctx context.Context, sessionID string) error

	// ListEvents retorna a lista de eventos disponíveis
	ListEvents(ctx context.Context) (*dto.ListEventsResponse, error)
}

