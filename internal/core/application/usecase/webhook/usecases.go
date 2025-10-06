package webhook

import (
	"context"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
	"zpwoot/internal/core/ports/input"
)

type WebhookUseCases struct {
	create     *CreateUseCase
	update     *UpdateUseCase
	upsert     *UpsertUseCase
	get        *GetUseCase
	delete     *DeleteUseCase
	listEvents *ListEventsUseCase
}

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
func (w *WebhookUseCases) Create(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error) {
	return w.create.Execute(ctx, sessionID, request)
}
func (w *WebhookUseCases) Update(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error) {
	return w.update.Execute(ctx, sessionID, request)
}
func (w *WebhookUseCases) Upsert(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error) {
	return w.upsert.Execute(ctx, sessionID, request)
}
func (w *WebhookUseCases) Get(ctx context.Context, sessionID string) (*dto.WebhookResponse, error) {
	return w.get.Execute(ctx, sessionID)
}
func (w *WebhookUseCases) Delete(ctx context.Context, sessionID string) error {
	return w.delete.Execute(ctx, sessionID)
}
func (w *WebhookUseCases) ListEvents(ctx context.Context) (*dto.ListEventsResponse, error) {
	return w.listEvents.Execute(ctx)
}
