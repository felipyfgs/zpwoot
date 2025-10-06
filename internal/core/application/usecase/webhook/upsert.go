package webhook

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
)


type UpsertUseCase struct {
	webhookRepo    webhook.Repository
	webhookService *webhook.Service
}


func NewUpsertUseCase(
	webhookRepo webhook.Repository,
	webhookService *webhook.Service,
) *UpsertUseCase {
	return &UpsertUseCase{
		webhookRepo:    webhookRepo,
		webhookService: webhookService,
	}
}


func (uc *UpsertUseCase) Execute(
	ctx context.Context,
	sessionID string,
	request *dto.CreateWebhookRequest,
) (*dto.WebhookResponse, error) {

	if err := uc.webhookService.ValidateURL(request.URL); err != nil {
		return nil, fmt.Errorf("invalid webhook URL: %w", err)
	}


	if err := uc.webhookService.ValidateEvents(request.Events); err != nil {
		return nil, fmt.Errorf("invalid events: %w", err)
	}


	if request.Secret != nil && *request.Secret != "" {
		if err := uc.webhookService.ValidateSecret(*request.Secret); err != nil {
			return nil, fmt.Errorf("invalid secret: %w", err)
		}
	}


	existingWebhook, err := uc.webhookRepo.GetBySessionID(ctx, sessionID)
	if err != nil && err.Error() != "webhook not found" {
		return nil, fmt.Errorf("failed to check existing webhook: %w", err)
	}

	var wh *webhook.Webhook

	if existingWebhook != nil {

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

		wh = webhook.NewWebhook(sessionID, request.URL, request.Events)


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
