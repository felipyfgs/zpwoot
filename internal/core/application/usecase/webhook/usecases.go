package webhook

import (
	"context"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
	"zpwoot/internal/core/ports/input"
)

// WebhookUseCases implementa a interface WebhookUseCases conectando todos os use cases individuais
// REGRA: Implementa ports/input/WebhookUseCases usando composição dos use cases
type WebhookUseCases struct {
	create     *CreateUseCase
	update     *UpdateUseCase
	upsert     *UpsertUseCase
	get        *GetUseCase
	delete     *DeleteUseCase
	listEvents *ListEventsUseCase
}

// NewWebhookUseCases cria uma nova instância do container de use cases
func NewWebhookUseCases(
	webhookRepo webhook.Repository,
	webhookService *webhook.Service,
) input.WebhookUseCases {
	return &WebhookUseCases{
		create:     NewCreateUseCase(webhookRepo, webhookService),
		update:     NewUpdateUseCase(webhookRepo, webhookService),
		upsert:     NewUpsertUseCase(webhookRepo, webhookService),
		get:        NewGetUseCase(webhookRepo),
		delete:     NewDeleteUseCase(webhookRepo),
		listEvents: NewListEventsUseCase(webhookService),
	}
}

// Create cria um novo webhook para uma sessão
func (w *WebhookUseCases) Create(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error) {
	return w.create.Execute(ctx, sessionID, request)
}

// Update atualiza um webhook existente
func (w *WebhookUseCases) Update(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error) {
	return w.update.Execute(ctx, sessionID, request)
}

// Upsert cria ou atualiza um webhook para uma sessão
func (w *WebhookUseCases) Upsert(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error) {
	return w.upsert.Execute(ctx, sessionID, request)
}

// Get busca o webhook configurado para uma sessão
func (w *WebhookUseCases) Get(ctx context.Context, sessionID string) (*dto.WebhookResponse, error) {
	return w.get.Execute(ctx, sessionID)
}

// Delete remove o webhook de uma sessão
func (w *WebhookUseCases) Delete(ctx context.Context, sessionID string) error {
	return w.delete.Execute(ctx, sessionID)
}

// ListEvents retorna a lista de eventos disponíveis
func (w *WebhookUseCases) ListEvents(ctx context.Context) (*dto.ListEventsResponse, error) {
	return w.listEvents.Execute(ctx)
}
