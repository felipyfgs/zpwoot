package webhook

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
)

// CreateUseCase implementa o caso de uso de criar webhook
type CreateUseCase struct {
	webhookRepo    webhook.Repository
	webhookService *webhook.Service
}

// NewCreateUseCase cria uma nova instância do use case
func NewCreateUseCase(
	webhookRepo webhook.Repository,
	webhookService *webhook.Service,
) *CreateUseCase {
	return &CreateUseCase{
		webhookRepo:    webhookRepo,
		webhookService: webhookService,
	}
}

// Execute executa o caso de uso
func (uc *CreateUseCase) Execute(
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
	if err == nil && existingWebhook != nil {
		return nil, fmt.Errorf("webhook already exists for session %s", sessionID)
	}

	// Criar novo webhook
	wh := webhook.NewWebhook(sessionID, request.URL, request.Events)

	// Gerar secret se não fornecido
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

// generateSecret gera um secret aleatório de 32 bytes (64 caracteres hex)
func generateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
