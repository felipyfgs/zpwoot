package output

import (
	"context"
	"time"
)


type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}



type WebhookSender interface {




	SendWebhook(ctx context.Context, url string, secret *string, event *WebhookEvent) error
}
