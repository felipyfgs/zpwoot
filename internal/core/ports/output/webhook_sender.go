package output

import (
	"context"
	"time"
)

// WebhookEvent representa um evento a ser enviado via webhook
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// WebhookSender define a interface para envio de webhooks
// REGRA: Interface no port, implementação no adapter
type WebhookSender interface {
	// SendWebhook envia um evento para o webhook configurado
	// url: URL do webhook
	// secret: Secret para assinatura HMAC (pode ser nil)
	// event: Evento a ser enviado
	SendWebhook(ctx context.Context, url string, secret *string, event *WebhookEvent) error
}
