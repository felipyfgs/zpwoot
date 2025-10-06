package waclient

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/ports/output"
)

type EventType string

const (
	EventMessage      EventType = "Message"
	EventConnected    EventType = "Connected"
	EventDisconnected EventType = "Disconnected"
	EventQR           EventType = "QR"
	EventReceipt      EventType = "Receipt"
	EventReadReceipt  EventType = "ReadReceipt"
	EventPresence     EventType = "Presence"
	EventChatPresence EventType = "ChatPresence"
	EventHistorySync  EventType = "HistorySync"
	EventLoggedOut    EventType = "LoggedOut"
)

var SupportedEventTypes = []EventType{
	EventMessage,
	EventConnected,
	EventDisconnected,
	EventQR,
	EventReadReceipt,
	EventPresence,
	EventChatPresence,
	EventHistorySync,
	EventLoggedOut,
}

type QREvent struct {
	Event     string    `json:"event"`
	Code      string    `json:"code,omitempty"`
	Base64    string    `json:"qrCodeBase64,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

type MessageInfo struct {
	ID        string    `json:"id"`
	Chat      string    `json:"chat"`
	Sender    string    `json:"sender"`
	PushName  string    `json:"pushName"`
	Timestamp time.Time `json:"timestamp"`
	FromMe    bool      `json:"fromMe"`
	Type      string    `json:"type"`
	IsGroup   bool      `json:"isGroup"`
}

type WebhookEvent struct {
	Type      EventType   `json:"type"`
	SessionID string      `json:"sessionId"`
	Event     interface{} `json:"event"`
	Timestamp time.Time   `json:"timestamp"`
}

type SessionConfig struct {
	SessionID     string            `json:"sessionId"`
	Name          string            `json:"name"`
	DeviceJID     string            `json:"deviceJid,omitempty"`
	ProxyConfig   map[string]string `json:"proxyConfig,omitempty"`
	WebhookURL    string            `json:"webhookUrl,omitempty"`
	Events        []EventType       `json:"events,omitempty"`
	AutoReconnect bool              `json:"autoReconnect"`
}

type Client struct {
	SessionID    string
	Name         string
	WAClient     *whatsmeow.Client
	EventHandler uint32
	Status       session.Status
	QRCode       string
	QRExpiresAt  time.Time
	ConnectedAt  time.Time
	LastSeen     time.Time
	Config       *SessionConfig
	Events       []EventType
	WebhookURL   string
	ctx          context.Context
	cancel       context.CancelFunc
}

func (c *Client) IsConnected() bool {

	if c.WAClient == nil {
		return false
	}
	return c.WAClient.IsConnected()
}

func (c *Client) IsLoggedIn() bool {
	return c.WAClient != nil && c.WAClient.Store.ID != nil
}

func (c *Client) GetDeviceJID() string {
	if c.IsLoggedIn() {
		return c.WAClient.Store.ID.String()
	}
	return ""
}

type EventHandler interface {
	HandleEvent(client *Client, event interface{}) error
}

type MediaProcessor interface {
	ProcessMedia(ctx context.Context, client *Client, media interface{}) (*output.MediaData, error)
}

type WebhookSender interface {
	SendWebhook(ctx context.Context, event *WebhookEvent) error
}

type SessionManager interface {
	CreateSession(ctx context.Context, config *SessionConfig) (*Client, error)
	GetSession(ctx context.Context, sessionID string) (*Client, error)
	GetSessionByName(ctx context.Context, name string) (*Client, error)
	ListSessions(ctx context.Context) ([]*Client, error)
	UpdateSession(ctx context.Context, client *Client) error
	DeleteSession(ctx context.Context, sessionID string) error
	ConnectSession(ctx context.Context, sessionID string) error
	DisconnectSession(ctx context.Context, sessionID string) error
}

type MessageSender interface {
	SendTextMessage(ctx context.Context, sessionID string, to string, text string) error
	SendMediaMessage(ctx context.Context, sessionID string, to string, media *output.MediaData) error
	SendLocationMessage(ctx context.Context, sessionID string, to string, lat, lng float64, name string) error
	SendContactMessage(ctx context.Context, sessionID string, to string, contact *output.ContactInfo) error
}

type SendMessageRequest struct {
	SessionID string              `json:"sessionId"`
	To        string              `json:"to"`
	Type      string              `json:"type"`
	Text      string              `json:"text,omitempty"`
	Media     *output.MediaData   `json:"media,omitempty"`
	Location  *output.Location    `json:"location,omitempty"`
	Contact   *output.ContactInfo `json:"contact,omitempty"`
}

type MessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId,omitempty"`
	Error     string `json:"error,omitempty"`
}

var (
	ErrSessionNotFound  = output.ErrSessionNotFound
	ErrSessionExists    = &output.WhatsAppError{Code: "SESSION_EXISTS", Message: "session already exists"}
	ErrNotConnected     = output.ErrSessionNotConnected
	ErrInvalidJID       = output.ErrInvalidJID
	ErrQRExpired        = output.ErrQRCodeExpired
	ErrConnectionFailed = output.ErrConnectionFailed
	ErrInvalidMedia     = &output.WhatsAppError{Code: "INVALID_MEDIA", Message: "invalid media data"}
	ErrWebhookFailed    = &output.WhatsAppError{Code: "WEBHOOK_FAILED", Message: "failed to send webhook"}
)
