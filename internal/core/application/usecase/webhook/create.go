package webhook

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
)

type CreateUseCase struct {
	webhookRepo    webhook.Repository
	webhookService *webhook.Service
}

func NewCreateUseCase(
	webhookRepo webhook.Repository,
	webhookService *webhook.Service,
) *CreateUseCase {
	return &CreateUseCase{
		webhookRepo:    webhookRepo,
		webhookService: webhookService,
	}
}
func (uc *CreateUseCase) Execute(
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
	if err == nil && existingWebhook != nil {
		return nil, fmt.Errorf("webhook already exists for session %s", sessionID)
	}

	wh := webhook.NewWebhook(sessionID, request.URL, request.Events)

	secret := request.Secret
	if secret == nil || *secret == "" {
		generated, err := generateSecret()
		if err != nil {
			return nil, fmt.Errorf("failed to generate secret: %w", err)
		}

		secret = &generated
	}

	wh.SetSecret(*secret)

	if err := uc.webhookRepo.Create(ctx, wh); err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
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
func generateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
