package dto

import (
	"encoding/base64"
	"time"

	"github.com/skip2/go-qrcode"
	"zpwoot/internal/core/application/validators"
	"zpwoot/internal/core/domain/session"
)

type ProxySettings struct {
	Enabled bool   `json:"enabled" example:"true" description:"Enable proxy"`
	Type    string `json:"type,omitempty" example:"http" enums:"http,https,socks5" description:"Proxy type (http, https, socks5)"`
	Host    string `json:"host,omitempty" example:"proxy.example.com" description:"Proxy host"`
	Port    string `json:"port,omitempty" example:"8080" description:"Proxy port"`
	User    string `json:"user,omitempty" example:"proxyUser123" description:"Proxy username (optional)"`
	Pass    string `json:"pass,omitempty" example:"proxyPass123" description:"Proxy password (optional)"`
} //@name ProxySettings

type WebhookSettings struct {
	Enabled bool     `json:"enabled" example:"true" description:"Enable webhook"`
	URL     string   `json:"url,omitempty" example:"https://api.example.com/webhook" validate:"omitempty,url" description:"Webhook URL"`
	Events  []string `json:"events,omitempty" example:"Message,Receipt,Connected" description:"Events to subscribe (Message, Receipt, Connected, Disconnected, CallOffer, Presence, NewsletterJoin, All)"`
	Secret  string   `json:"secret,omitempty" example:"supersecrettoken123" description:"Webhook secret for validation (optional)"`
} //@name WebhookSettings

type SessionSettings struct {
	Proxy   *ProxySettings   `json:"proxy,omitempty" description:"Proxy configuration"`
	Webhook *WebhookSettings `json:"webhook,omitempty" description:"Webhook configuration"`
} //@name SessionSettings

type CreateRequest struct {
	Name           string           `json:"name" example:"my-session" validate:"required,min=1,max=100" description:"Session name for identification"`
	Settings       *SessionSettings `json:"settings,omitempty" description:"Session settings (proxy, webhook)"`
	GenerateQRCode bool             `json:"qrCode,omitempty" example:"true" description:"Auto-generate QR code after creation (default: false)"`
} //@name CreateSessionRequest

type UpdateRequest struct {
	Name     *string          `json:"name,omitempty" example:"updated-session" validate:"omitempty,min=1,max=100" description:"Session name for identification"`
	Settings *SessionSettings `json:"settings,omitempty" description:"Session settings (proxy, webhook)"`
}

type SessionResponse struct {
	SessionID       string           `json:"sessionId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Unique session identifier"`
	Name            string           `json:"name" example:"my-session" description:"Session name"`
	Status          string           `json:"status" example:"connected" description:"Current session status (disconnected, connecting, connected, qr_code, error)"`
	Connected       bool             `json:"connected" example:"true" description:"Whether session is connected"`
	DeviceJID       string           `json:"deviceJid,omitempty" example:"5511999999999@s.whatsapp.net" description:"WhatsApp device JID when connected"`
	Settings        *SessionSettings `json:"settings,omitempty" description:"Session settings (proxy, webhook)"`
	QRCode          string           `json:"qrCode,omitempty" example:"2@abc123..." description:"QR code string (original from WhatsApp)"`
	QRCodeBase64    string           `json:"qrCodeBase64,omitempty" example:"data:image/png;base64,iVBORw0KGgo..." description:"QR code as base64 image"`
	QRCodeExpiresAt *time.Time       `json:"qrCodeExpiresAt,omitempty" example:"2025-01-15T10:35:00Z" description:"QR code expiration time"`
	CreatedAt       time.Time        `json:"createdAt" example:"2025-01-15T10:30:00Z" description:"Session creation timestamp"`
	UpdatedAt       time.Time        `json:"updatedAt" example:"2025-01-15T10:35:00Z" description:"Last update timestamp"`
	ConnectedAt     *time.Time       `json:"connectedAt,omitempty" example:"2025-01-15T10:32:00Z" description:"Connection timestamp"`
	LastSeen        *time.Time       `json:"lastSeen,omitempty" example:"2025-01-15T10:35:00Z" description:"Last activity timestamp"`
} //@name SessionResponse

type SessionListInfo struct {
	SessionID   string           `json:"sessionId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Unique session identifier"`
	Name        string           `json:"name" example:"my-session" description:"Session name"`
	Status      string           `json:"status" example:"connected" description:"Current session status (disconnected, connecting, connected, qr_code, error)"`
	Connected   bool             `json:"connected" example:"true" description:"Whether session is connected"`
	DeviceJID   string           `json:"deviceJid,omitempty" example:"5511999999999@s.whatsapp.net" description:"WhatsApp device JID when connected"`
	Settings    *SessionSettings `json:"settings,omitempty" description:"Session settings (proxy, webhook)"`
	CreatedAt   time.Time        `json:"createdAt" example:"2025-01-15T10:30:00Z" description:"Session creation timestamp"`
	UpdatedAt   time.Time        `json:"updatedAt" example:"2025-01-15T10:35:00Z" description:"Last update timestamp"`
	ConnectedAt *time.Time       `json:"connectedAt,omitempty" example:"2025-01-15T10:32:00Z" description:"Connection timestamp"`
	LastSeen    *time.Time       `json:"lastSeen,omitempty" example:"2025-01-15T10:35:00Z" description:"Last activity timestamp"`
} //@name SessionListInfo

type SessionListResponse struct {
	Sessions []SessionListInfo `json:"sessions" description:"List of sessions (without QR codes)"`
	Total    int               `json:"total" example:"5" description:"Total number of sessions"`
} //@name SessionListResponse

type SessionActionResponse struct {
	SessionID string `json:"sessionId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Session ID"`
	Action    string `json:"action" example:"connect" description:"Action performed (connect, disconnect, logout, delete)"`
	Status    string `json:"status" example:"success" description:"Action status"`
	Message   string `json:"message,omitempty" example:"Session connected successfully" description:"Action message"`
}

type QRCodeResponse struct {
	QRCode       string `json:"qrCode" example:"2@abc123..." description:"QR code string (original from WhatsApp)"`
	QRCodeBase64 string `json:"qrCodeBase64" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..." description:"Base64 encoded QR code image"`
	ExpiresAt    string `json:"expiresAt" example:"2025-01-15T10:35:00Z" description:"QR code expiration time"`
	Status       string `json:"status" example:"generated" description:"QR code status"`
} //@name QRCodeResponse

type CreateSessionResponse struct {
	ID              string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" description:"Session identifier"`
	Name            string     `json:"name" example:"My WhatsApp Session" description:"Session name"`
	Status          string     `json:"status" example:"disconnected" description:"Initial session status"`
	Connected       bool       `json:"connected" example:"false" description:"Whether session is connected"`
	QRCode          string     `json:"qrCode,omitempty" example:"2@abc123..." description:"QR code string (if generated)"`
	QRCodeBase64    string     `json:"qrCodeBase64,omitempty" example:"data:image/png;base64,iVBORw0KGgo..." description:"QR code as base64 image (if generated)"`
	QRCodeExpiresAt *time.Time `json:"qrCodeExpiresAt,omitempty" example:"2025-01-15T10:35:00Z" description:"QR code expiration time (if generated)"`
	CreatedAt       time.Time  `json:"createdAt" example:"2025-01-15T10:30:00Z" description:"Session creation timestamp"`
}

type SessionDetailResponse struct {
	ID              string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" description:"Session identifier"`
	Name            string     `json:"name" example:"My WhatsApp Session" description:"Session name"`
	DeviceJID       string     `json:"deviceJid,omitempty" example:"5511999999999@s.whatsapp.net" description:"WhatsApp device JID when connected"`
	Status          string     `json:"status" example:"connected" description:"Current session status"`
	Connected       bool       `json:"connected" example:"true" description:"Whether session is connected"`
	ConnectionError string     `json:"connectionError,omitempty" example:"Connection timeout" description:"Connection error message if any"`
	QRCode          string     `json:"qrCode,omitempty" description:"QR code for authentication"`
	QRCodeExpiresAt *time.Time `json:"qrCodeExpiresAt,omitempty" example:"2025-01-15T10:35:00Z" description:"QR code expiration time"`
	ProxyConfig     *string    `json:"proxyConfig,omitempty" description:"Proxy configuration as JSON string"`
	CreatedAt       time.Time  `json:"createdAt" example:"2025-01-15T10:30:00Z" description:"Session creation timestamp"`
	UpdatedAt       time.Time  `json:"updatedAt" example:"2025-01-15T10:35:00Z" description:"Last update timestamp"`
	ConnectedAt     *time.Time `json:"connectedAt,omitempty" example:"2025-01-15T10:32:00Z" description:"Connection timestamp"`
	LastSeen        *time.Time `json:"lastSeen,omitempty" example:"2025-01-15T10:35:00Z" description:"Last activity timestamp"`
}

type SessionStatusResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" description:"Session identifier"`
	Status    string `json:"status" example:"connected" description:"Current session status"`
	Connected bool   `json:"connected" example:"true" description:"Whether session is connected"`
	Message   string `json:"message,omitempty" example:"Session is already connected" description:"Additional status message"`
} //@name SessionStatusResponse

func (r *CreateRequest) Validate() error {

	if err := validators.ValidateSessionName(r.Name); err != nil {
		return NewValidationError("name", err.Error())
	}

	if r.Settings != nil && r.Settings.Webhook != nil && r.Settings.Webhook.Enabled {
		if err := validators.ValidateWebhookURL(r.Settings.Webhook.URL); err != nil {
			return NewValidationError("settings.webhook.url", err.Error())
		}
	}

	if r.Settings != nil && r.Settings.Proxy != nil && r.Settings.Proxy.Enabled {
		if r.Settings.Proxy.Host == "" {
			return NewValidationError("settings.proxy.host", "proxy host is required when proxy is enabled")
		}
		if r.Settings.Proxy.Port == "" {
			return NewValidationError("settings.proxy.port", "proxy port is required when proxy is enabled")
		}
	}

	return nil
}

func (r *CreateRequest) ToDomain() *session.Session {
	return session.NewSession(r.Name)
}

func FromDomain(s *session.Session) *SessionResponse {
	response := &SessionResponse{
		SessionID:   s.ID,
		Name:        s.Name,
		Status:      string(s.GetStatus()),
		Connected:   s.IsConnected,
		DeviceJID:   s.DeviceJID,
		QRCode:      s.QRCode,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		ConnectedAt: s.ConnectedAt,
		LastSeen:    s.LastSeen,
	}

	if s.QRCode != "" {
		response.QRCodeBase64 = QRBase64(s.QRCode)
		if s.QRCodeExpiresAt != nil {
			response.QRCodeExpiresAt = s.QRCodeExpiresAt
		}
	}

	return response
}

func (s *SessionResponse) ToListInfo() *SessionListInfo {
	return &SessionListInfo{
		SessionID:   s.SessionID,
		Name:        s.Name,
		Status:      s.Status,
		Connected:   s.Connected,
		DeviceJID:   s.DeviceJID,
		Settings:    s.Settings,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		ConnectedAt: s.ConnectedAt,
		LastSeen:    s.LastSeen,
	}
}

func NewQRResponse(qrCode string, expiresAt time.Time, status string) *QRCodeResponse {
	expiresAtStr := ""
	if !expiresAt.IsZero() {
		expiresAtStr = expiresAt.Format(time.RFC3339)
	}

	return &QRCodeResponse{
		QRCode:       qrCode,
		QRCodeBase64: QRBase64(qrCode),
		ExpiresAt:    expiresAtStr,
		Status:       status,
	}
}

func ToDetailResponse(s *session.Session) *SessionDetailResponse {
	return &SessionDetailResponse{
		ID:              s.ID,
		Name:            s.Name,
		DeviceJID:       s.DeviceJID,
		Status:          string(s.GetStatus()),
		Connected:       s.IsConnected,
		ConnectionError: s.ConnectionError,
		QRCode:          s.QRCode,
		QRCodeExpiresAt: s.QRCodeExpiresAt,
		ProxyConfig:     s.ProxyConfig,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		ConnectedAt:     s.ConnectedAt,
		LastSeen:        s.LastSeen,
	}
}

func ToCreateResponse(s *session.Session) *CreateSessionResponse {
	response := &CreateSessionResponse{
		ID:        s.ID,
		Name:      s.Name,
		Status:    string(s.GetStatus()),
		Connected: s.IsConnected,
		CreatedAt: s.CreatedAt,
	}

	if s.QRCode != "" {
		response.QRCode = s.QRCode
		response.QRCodeBase64 = QRBase64(s.QRCode)
		if s.QRCodeExpiresAt != nil {
			response.QRCodeExpiresAt = s.QRCodeExpiresAt
		}
	}

	return response
}

func ToStatusResponse(s *session.Session) *SessionStatusResponse {
	return &SessionStatusResponse{
		ID:        s.ID,
		Status:    string(s.GetStatus()),
		Connected: s.IsConnected,
	}
}

func ToListResponse(s *session.Session) *SessionResponse {
	return FromDomain(s)
}

func ToListInfo(s *session.Session) *SessionListInfo {
	return &SessionListInfo{
		SessionID:   s.ID,
		Name:        s.Name,
		Status:      string(s.GetStatus()),
		Connected:   s.IsConnected,
		DeviceJID:   s.DeviceJID,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		ConnectedAt: s.ConnectedAt,
		LastSeen:    s.LastSeen,
	}
}

func QRBase64(qrString string) string {
	if qrString == "" {
		return ""
	}

	qrImage, err := qrcode.Encode(qrString, qrcode.Medium, 256)
	if err != nil {
		return ""
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrImage)
}

// PairPhoneRequest representa a requisição de pareamento por telefone
type PairPhoneRequest struct {
	Phone string `json:"phone" validate:"required" example:"5511999999999" description:"Phone number with country code"`
} //@name PairPhoneRequest

// PairPhoneResponse representa a resposta do pareamento por telefone
type PairPhoneResponse struct {
	LinkingCode string `json:"linkingCode" example:"ABCD-EFGH" description:"8-character linking code to enter on phone"`
} //@name PairPhoneResponse

var (
	ErrInvalidSessionName   = NewValidationError("name", "Session name is required")
	ErrSessionNameTooLong   = NewValidationError("name", "Session name must be less than 100 characters")
	ErrSessionNotFound      = &ErrorInfo{Code: "SESSION_NOT_FOUND", Message: "Session not found"}
	ErrSessionAlreadyExists = &ErrorInfo{Code: "SESSION_ALREADY_EXISTS", Message: "Session already exists"}
)
