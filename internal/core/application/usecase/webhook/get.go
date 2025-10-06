package webhook

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
)


type GetUseCase struct {
	webhookRepo webhook.Repository
}


func NewGetUseCase(webhookRepo webhook.Repository) *GetUseCase {
	return &GetUseCase{
		webhookRepo: webhookRepo,
	}
}


func (uc *GetUseCase) Execute(ctx context.Context, sessionID string) (*dto.WebhookResponse, error) {
	wh, err := uc.webhookRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}

	return &dto.WebhookResponse{
		ID:        wh.ID,
		SessionID: wh.SessionID,
		URL:       wh.URL,
		Events:    wh.Events,
		Enabled:   wh.Enabled,
		CreatedAt: wh.CreatedAt,
		UpdatedAt: wh.UpdatedAt,
	}, nil
}
