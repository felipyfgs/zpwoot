package output

import (
	"context"
	"time"
)

type WhatsAppClient interface {
	CreateSession(ctx context.Context, sessionID string) error
	GetSessionStatus(ctx context.Context, sessionID string) (*SessionStatus, error)
	DeleteSession(ctx context.Context, sessionID string) error

	ConnectSession(ctx context.Context, sessionID string) error
	DisconnectSession(ctx context.Context, sessionID string) error
	LogoutSession(ctx context.Context, sessionID string) error
	IsConnected(ctx context.Context, sessionID string) bool
	IsLoggedIn(ctx context.Context, sessionID string) bool

	GetQRCode(ctx context.Context, sessionID string) (*QRCodeInfo, error)

	SendTextMessage(ctx context.Context, sessionID, to, text string) (*MessageResult, error)
	SendMediaMessage(ctx context.Context, sessionID, to string, media *MediaData) (*MessageResult, error)
	SendLocationMessage(ctx context.Context, sessionID, to string, location *Location) (*MessageResult, error)
	SendContactMessage(ctx context.Context, sessionID, to string, contact *ContactInfo) (*MessageResult, error)
}

type SessionStatus struct {
	SessionID   string    `json:"sessionId"`
	Connected   bool      `json:"connected"`
	LoggedIn    bool      `json:"loggedIn"`
	DeviceJID   string    `json:"deviceJid,omitempty"`
	PushName    string    `json:"pushName,omitempty"`
	ConnectedAt time.Time `json:"connectedAt,omitempty"`
	LastSeen    time.Time `json:"lastSeen,omitempty"`
}

type QRCodeInfo struct {
	Code      string    `json:"code"`
	Base64    string    `json:"base64"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type MessageResult struct {
	MessageID string    `json:"messageId"`
	Status    string    `json:"status"`
	SentAt    time.Time `json:"sentAt"`
}

type MediaData struct {
	MimeType string `json:"mimeType"`
	Data     []byte `json:"data"`
	FileName string `json:"fileName,omitempty"`
	Caption  string `json:"caption,omitempty"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

type ContactInfo struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

type WhatsAppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *WhatsAppError) Error() string {
	return e.Message
}

var (
	ErrSessionNotFound     = &WhatsAppError{Code: "SESSION_NOT_FOUND", Message: "Session not found"}
	ErrSessionNotConnected = &WhatsAppError{Code: "SESSION_NOT_CONNECTED", Message: "Session not connected"}
	ErrAlreadyConnected    = &WhatsAppError{Code: "ALREADY_CONNECTED", Message: "Session already connected"}
	ErrInvalidJID          = &WhatsAppError{Code: "INVALID_JID", Message: "Invalid JID format"}
	ErrQRCodeExpired       = &WhatsAppError{Code: "QR_CODE_EXPIRED", Message: "QR code expired"}
	ErrConnectionFailed    = &WhatsAppError{Code: "CONNECTION_FAILED", Message: "Failed to connect to WhatsApp"}
	ErrSendMessageFailed   = &WhatsAppError{Code: "SEND_MESSAGE_FAILED", Message: "Failed to send message"}
)
