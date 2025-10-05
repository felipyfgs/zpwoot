package interfaces

import (
	"context"
	"time"
)


type NotificationService interface {

	SendWebhook(ctx context.Context, event *WebhookEvent) error


	NotifySessionConnected(ctx context.Context, sessionID string, deviceJID string) error
	NotifySessionDisconnected(ctx context.Context, sessionID string) error
	NotifyQRCodeGenerated(ctx context.Context, sessionID string, qrCode string, expiresAt time.Time) error
	NotifyMessageReceived(ctx context.Context, sessionID string, message *MessageEvent) error
	NotifyMessageSent(ctx context.Context, sessionID string, messageID string) error
}


type WebhookEvent struct {
	Type      string      `json:"type"`
	SessionID string      `json:"sessionId"`
	Event     interface{} `json:"event"`
	Timestamp time.Time   `json:"timestamp"`
}


type MessageEvent struct {
	ID        string    `json:"id"`
	Chat      string    `json:"chat"`
	Sender    string    `json:"sender"`
	PushName  string    `json:"pushName"`
	Timestamp time.Time `json:"timestamp"`
	FromMe    bool      `json:"fromMe"`
	Type      string    `json:"type"`
	IsGroup   bool      `json:"isGroup"`
	Content   string    `json:"content,omitempty"`
}


type SessionEvent struct {
	SessionID   string    `json:"sessionId"`
	Status      string    `json:"status"`
	DeviceJID   string    `json:"deviceJid,omitempty"`
	ConnectedAt time.Time `json:"connectedAt,omitempty"`
}


type QRCodeEvent struct {
	SessionID string    `json:"sessionId"`
	QRCode    string    `json:"qrCode"`
	ExpiresAt time.Time `json:"expiresAt"`
}


const (
	EventTypeSessionConnected    = "session.connected"
	EventTypeSessionDisconnected = "session.disconnected"
	EventTypeQRCodeGenerated     = "qr.generated"
	EventTypeMessageReceived     = "message.received"
	EventTypeMessageSent         = "message.sent"
)
