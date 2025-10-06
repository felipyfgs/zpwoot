package webhook

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
)

// UpsertUseCase implementa o caso de uso de criar ou atualizar webhook (upsert)
type UpsertUseCase struct {
	webhookRepo    webhook.Repository
	webhookService *webhook.Service
}

// NewUpsertUseCase cria uma nova instância do use case
func NewUpsertUseCase(
	webhookRepo webhook.Repository,
	webhookService *webhook.Service,
) *UpsertUseCase {
	return &UpsertUseCase{
		webhookRepo:    webhookRepo,
		webhookService: webhookService,
	}
}

// Execute executa o caso de uso (cria se não existe, atualiza se existe)
func (uc *UpsertUseCase) Execute(
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

	// Verificar se já existe webhook para esta sessão
	existingWebhook, err := uc.webhookRepo.GetBySessionID(ctx, sessionID)
	if err != nil && err.Error() != "webhook not found" {
		return nil, fmt.Errorf("failed to check existing webhook: %w", err)
	}

	var wh *webhook.Webhook

	if existingWebhook != nil {
		// Update
		existingWebhook.UpdateURL(request.URL)
		existingWebhook.UpdateEvents(request.Events)

		if request.Secret != nil && *request.Secret != "" {
			existingWebhook.SetSecret(*request.Secret)
		}

		if err := uc.webhookRepo.Update(ctx, existingWebhook); err != nil {
			return nil, fmt.Errorf("failed to update webhook: %w", err)
		}
		wh = existingWebhook
	} else {
		// Create
		wh = webhook.NewWebhook(sessionID, request.URL, request.Events)

		// Gerar secret se não fornecido
		secret := request.Secret
		if secret == nil || *secret == "" {
			generated, err := generateSecretKey()
			if err != nil {
				return nil, fmt.Errorf("failed to generate secret: %w", err)
			}
			secret = &generated
		}
		wh.SetSecret(*secret)

		if err := uc.webhookRepo.Create(ctx, wh); err != nil {
			return nil, fmt.Errorf("failed to create webhook: %w", err)
		}
	}

	// Converter para DTO
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
