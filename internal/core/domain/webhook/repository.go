package webhook

import "context"

type Repository interface {
	Create(ctx context.Context, webhook *Webhook) error
	GetByID(ctx context.Context, id string) (*Webhook, error)
	GetBySessionID(ctx context.Context, sessionID string) (*Webhook, error)
	Update(ctx context.Context, webhook *Webhook) error
	Delete(ctx context.Context, id string) error
	DeleteBySessionID(ctx context.Context, sessionID string) error
	List(ctx context.Context, limit, offset int) ([]*Webhook, error)
}
