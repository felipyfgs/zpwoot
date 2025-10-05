package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// WhatsAppGateway defines WhatsApp operations interface
type WhatsAppGateway interface {
	CreateSession(ctx context.Context, sessionId uuid.UUID) error
	ConnectSession(ctx context.Context, sessionId uuid.UUID) error
	DisconnectSession(ctx context.Context, sessionId uuid.UUID) error
	DeleteSession(ctx context.Context, sessionId uuid.UUID) error
	RestoreSession(ctx context.Context, sessionId uuid.UUID) error
	RestoreAllSessions(ctx context.Context, sessionIds []uuid.UUID) error
	SessionExists(sessionId uuid.UUID) bool

	IsSessionConnected(ctx context.Context, sessionId uuid.UUID) (bool, error)
	GetSessionInfo(ctx context.Context, sessionId uuid.UUID) (*DeviceInfo, error)

	GenerateQRCode(ctx context.Context, sessionId uuid.UUID) (*QRCodeResponse, error)

	SetProxy(ctx context.Context, sessionId uuid.UUID, proxy *ProxyConfig) error

	SetEventHandler(handler EventHandler)

	SendTextMessage(ctx context.Context, sessionId uuid.UUID, to, content string) (*MessageSendResult, error)
	SendMediaMessage(ctx context.Context, sessionId uuid.UUID, to, mediaURL, caption, mediaType string) (*MessageSendResult, error)
	SendLocationMessage(ctx context.Context, sessionId uuid.UUID, to string, latitude, longitude float64, address string) (*MessageSendResult, error)
	SendContactMessage(ctx context.Context, sessionId uuid.UUID, to, contactName, contactPhone string) (*MessageSendResult, error)
}

type EventHandler interface {
	OnSessionConnected(sessionId uuid.UUID, deviceInfo *DeviceInfo)
	OnSessionDisconnected(sessionId uuid.UUID, reason string)
	OnQRCodeGenerated(sessionId uuid.UUID, qrCode string, expiresAt time.Time)
	OnConnectionError(sessionId uuid.UUID, err error)
	OnMessageReceived(sessionId uuid.UUID, message *WhatsAppMessage)
	OnMessageSent(sessionId uuid.UUID, messageID string, status string)
}

type WhatsAppMessage struct {
	ID        string                 `json:"id"`
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	Chat      string                 `json:"chat"`
	Type      string                 `json:"type"`
	Content   string                 `json:"content,omitempty"`
	MediaURL  string                 `json:"media_url,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	FromMe    bool                   `json:"from_me"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type MessageSendResult struct {
	MessageID string    `json:"message_id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	To        string    `json:"to"`
}

type QRCodeGenerator interface {
	Generate(ctx context.Context, sessionId uuid.UUID) (*QRCodeResponse, error)
	GenerateImage(ctx context.Context, qrCode string) ([]byte, error)
	IsExpired(expiresAt time.Time) bool
}

// DeviceInfo represents WhatsApp device information
type DeviceInfo struct {
	Platform    string `json:"platform"`
	DeviceModel string `json:"deviceModel"`
	OSVersion   string `json:"osVersion"`
	AppVersion  string `json:"appVersion"`
}

// QRCodeResponse represents QR code information
type QRCodeResponse struct {
	QRCode      string    `json:"qrCode"`
	QRCodeImage string    `json:"qrCodeImage,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt"`
	Timeout     int       `json:"timeoutSeconds"`
}
