package waclient

import (
	"context"
	"time"

	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/ports/output"

	"go.mau.fi/whatsmeow"
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
	return c.WAClient != nil && c.WAClient.IsConnected()
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

type WebhookSender interface {
	SendWebhook(ctx context.Context, event *WebhookEvent) error
}

type ContactInfo struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	VCard string `json:"vcard,omitempty"`
}

var (
	ErrSessionNotFound  = output.ErrSessionNotFound
	ErrSessionExists    = &output.WhatsAppError{Code: "SESSION_EXISTS", Message: "session already exists"}
	ErrNotConnected     = output.ErrSessionNotConnected
	ErrInvalidJID       = output.ErrInvalidJID
	ErrConnectionFailed = output.ErrConnectionFailed
	ErrAlreadyPaired    = &output.WhatsAppError{Code: "ALREADY_PAIRED", Message: "session is already paired"}
)
