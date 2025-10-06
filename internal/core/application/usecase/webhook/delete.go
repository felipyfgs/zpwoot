package webhook

import (
	"context"
	"fmt"

	"zpwoot/internal/core/domain/webhook"
)

type DeleteUseCase struct {
	webhookRepo webhook.Repository
}

func NewDeleteUseCase(webhookRepo webhook.Repository) *DeleteUseCase {
	return &DeleteUseCase{
		webhookRepo: webhookRepo,
	}
}
func (uc *DeleteUseCase) Execute(ctx context.Context, sessionID string) error {
	if err := uc.webhookRepo.DeleteBySessionID(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	return nil
}
