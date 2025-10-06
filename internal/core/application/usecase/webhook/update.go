package webhook

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
)

// UpdateUseCase implementa o caso de uso de atualizar webhook
type UpdateUseCase struct {
	webhookRepo    webhook.Repository
	webhookService *webhook.Service
}

// NewUpdateUseCase cria uma nova instância do use case
func NewUpdateUseCase(
	webhookRepo webhook.Repository,
	webhookService *webhook.Service,
) *UpdateUseCase {
	return &UpdateUseCase{
		webhookRepo:    webhookRepo,
		webhookService: webhookService,
	}
}

// Execute executa o caso de uso
func (uc *UpdateUseCase) Execute(
	ctx context.Context,
	sessionID string,
	request *dto.CreateWebhookRequest,
) (*dto.WebhookResponse, error) {
	// Validar URL
	if err := uc.webhookService.ValidateURL(request.URL); err != nil {
		return nil, fmt.Errorf("invalid webhook URL: %w", err)
	}

	// Validar eventos
	if err := uc.webhookService.ValidateEvents(request.Events); err != nil {
		return nil, fmt.Errorf("invalid events: %w", err)
	}

	// Validar secret se fornecido
	if request.Secret != nil && *request.Secret != "" {
		if err := uc.webhookService.ValidateSecret(*request.Secret); err != nil {
			return nil, fmt.Errorf("invalid secret: %w", err)
		}
	}

	// Buscar webhook existente
	existingWebhook, err := uc.webhookRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("webhook not found for session %s: %w", sessionID, err)
	}

	// Atualizar campos
	existingWebhook.UpdateURL(request.URL)
	existingWebhook.UpdateEvents(request.Events)

	if request.Secret != nil && *request.Secret != "" {
		existingWebhook.SetSecret(*request.Secret)
	}

	// Salvar alterações
	if err := uc.webhookRepo.Update(ctx, existingWebhook); err != nil {
		return nil, fmt.Errorf("failed to update webhook: %w", err)
	}

	// Converter para DTO
	return &dto.WebhookResponse{
		ID:        existingWebhook.ID,
		SessionID: existingWebhook.SessionID,
		URL:       existingWebhook.URL,
		Events:    existingWebhook.Events,
		Enabled:   existingWebhook.Enabled,
		CreatedAt: existingWebhook.CreatedAt,
		UpdatedAt: existingWebhook.UpdatedAt,
	}, nil
}
