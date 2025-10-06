package input

import (
	"context"

	"zpwoot/internal/core/application/dto"
)

type WebhookUseCases interface {
	Create(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error)
	Update(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error)
	Upsert(ctx context.Context, sessionID string, request *dto.CreateWebhookRequest) (*dto.WebhookResponse, error)
	Get(ctx context.Context, sessionID string) (*dto.WebhookResponse, error)
	Delete(ctx context.Context, sessionID string) error
	ListEvents(ctx context.Context) (*dto.ListEventsResponse, error)
}
