package webhook

import (
	"context"
	"fmt"

	"zpwoot/internal/core/domain/webhook"
)

// DeleteUseCase implementa o caso de uso de deletar webhook
type DeleteUseCase struct {
	webhookRepo webhook.Repository
}

// NewDeleteUseCase cria uma nova instância do use case
func NewDeleteUseCase(webhookRepo webhook.Repository) *DeleteUseCase {
	return &DeleteUseCase{
		webhookRepo: webhookRepo,
	}
}

// Execute executa o caso de uso
func (uc *DeleteUseCase) Execute(ctx context.Context, sessionID string) error {
	if err := uc.webhookRepo.DeleteBySessionID(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	return nil
}

