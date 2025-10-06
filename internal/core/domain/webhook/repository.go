package webhook

import "context"

// Repository define o contrato para persistência de webhooks
// REGRA: Interface no domain, implementação nos adapters
type Repository interface {
	// Create cria um novo webhook
	Create(ctx context.Context, webhook *Webhook) error

	// GetByID busca um webhook pelo ID
	GetByID(ctx context.Context, id string) (*Webhook, error)

	// GetBySessionID busca um webhook pela sessão
	GetBySessionID(ctx context.Context, sessionID string) (*Webhook, error)

	// Update atualiza um webhook existente
	Update(ctx context.Context, webhook *Webhook) error

	// Delete remove um webhook pelo ID
	Delete(ctx context.Context, id string) error

	// DeleteBySessionID remove um webhook pela sessão
	DeleteBySessionID(ctx context.Context, sessionID string) error

	// List lista todos os webhooks (opcional, para admin)
	List(ctx context.Context, limit, offset int) ([]*Webhook, error)
}

