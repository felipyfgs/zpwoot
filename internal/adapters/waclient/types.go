package waclient

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
)

// SessionStatus represents the status of a WhatsApp session
type SessionStatus string

const (
	StatusDisconnected SessionStatus = "disconnected"
	StatusConnecting   SessionStatus = "connecting"
	StatusConnected    SessionStatus = "connected"
	StatusQRCode       SessionStatus = "qr_code"
	StatusError        SessionStatus = "error"
)

// MediaType represents the type of media
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeVideo    MediaType = "video"
	MediaTypeDocument MediaType = "document"
	MediaTypeSticker  MediaType = "sticker"
)

// EventType represents WhatsApp event types
type EventType string

const (
	EventMessage      EventType = "Message"
	EventConnected    EventType = "Connected"
	EventDisconnected EventType = "Disconnected"
	EventQR           EventType = "QR"
	EventReadReceipt  EventType = "ReadReceipt"
	EventPresence     EventType = "Presence"
	EventChatPresence EventType = "ChatPresence"
	EventHistorySync  EventType = "HistorySync"
	EventLoggedOut    EventType = "LoggedOut"
)

// SupportedEventTypes lists all supported event types
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

// QREvent represents a QR code event
type QREvent struct {
	Event     string    `json:"event"`
	Code      string    `json:"code,omitempty"`
	Base64    string    `json:"qrCodeBase64,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

// MediaData represents processed media data
type MediaData struct {
	Base64   string `json:"base64,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	FileName string `json:"fileName,omitempty"`
	Size     int64  `json:"size,omitempty"`
}

// MessageInfo represents WhatsApp message information
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

// WebhookEvent represents an event to be sent via webhook
type WebhookEvent struct {
	Type      EventType   `json:"type"`
	SessionID string      `json:"sessionId"`
	Event     interface{} `json:"event"`
	Timestamp time.Time   `json:"timestamp"`
}

// SessionConfig represents session configuration
type SessionConfig struct {
	SessionID     string            `json:"sessionId"`
	Name          string            `json:"name"`
	DeviceJID     string            `json:"deviceJid,omitempty"`
	ProxyConfig   map[string]string `json:"proxyConfig,omitempty"`
	WebhookURL    string            `json:"webhookUrl,omitempty"`
	Events        []EventType       `json:"events,omitempty"`
	AutoReconnect bool              `json:"autoReconnect"`
}

// Client represents a WhatsApp client instance
type Client struct {
	SessionID    string
	Name         string
	WAClient     *whatsmeow.Client
	EventHandler uint32
	Status       SessionStatus
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

// EventHandler defines the interface for handling WhatsApp events
type EventHandler interface {
	HandleEvent(client *Client, event interface{}) error
}

// MediaProcessor defines the interface for processing media
type MediaProcessor interface {
	ProcessMedia(ctx context.Context, client *Client, media interface{}) (*MediaData, error)
}

// WebhookSender defines the interface for sending webhooks
type WebhookSender interface {
	SendWebhook(ctx context.Context, event *WebhookEvent) error
}

// SessionManager defines the interface for managing sessions
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

// MessageSender defines the interface for sending messages
type MessageSender interface {
	SendTextMessage(ctx context.Context, sessionID string, to string, text string) error
	SendMediaMessage(ctx context.Context, sessionID string, to string, media *MediaData) error
	SendLocationMessage(ctx context.Context, sessionID string, to string, lat, lng float64, name string) error
	SendContactMessage(ctx context.Context, sessionID string, to string, contact *ContactInfo) error
}

// ContactInfo represents contact information
type ContactInfo struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	VCard string `json:"vcard,omitempty"`
}

// SendMessageRequest represents a message sending request
type SendMessageRequest struct {
	SessionID string       `json:"sessionId"`
	To        string       `json:"to"`
	Type      string       `json:"type"`
	Text      string       `json:"text,omitempty"`
	Media     *MediaData   `json:"media,omitempty"`
	Location  *Location    `json:"location,omitempty"`
	Contact   *ContactInfo `json:"contact,omitempty"`
}

// Location represents location data
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// MessageResponse represents a message sending response
type MessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId,omitempty"`
	Error     string `json:"error,omitempty"`
}

// SessionInfo represents internal session information (for database and internal use)
// Note: For API responses, use dto.SessionResponse or dto.SessionListInfo instead
type SessionInfo struct {
	ID          string
	Name        string
	DeviceJID   string
	Status      SessionStatus
	Connected   bool
	QRCode      string
	QRExpiresAt time.Time
	ConnectedAt time.Time
	LastSeen    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Error types
type WAError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *WAError) Error() string {
	return e.Message
}

// Common errors
var (
	ErrSessionNotFound  = &WAError{Code: "SESSION_NOT_FOUND", Message: "session not found"}
	ErrSessionExists    = &WAError{Code: "SESSION_EXISTS", Message: "session already exists"}
	ErrNotConnected     = &WAError{Code: "NOT_CONNECTED", Message: "session not connected"}
	ErrInvalidJID       = &WAError{Code: "INVALID_JID", Message: "invalid JID format"}
	ErrQRExpired        = &WAError{Code: "QR_EXPIRED", Message: "QR code expired"}
	ErrConnectionFailed = &WAError{Code: "CONNECTION_FAILED", Message: "failed to connect to WhatsApp"}
	ErrInvalidMedia     = &WAError{Code: "INVALID_MEDIA", Message: "invalid media data"}
	ErrWebhookFailed    = &WAError{Code: "WEBHOOK_FAILED", Message: "failed to send webhook"}
)
