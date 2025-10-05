package dto

import (
	"encoding/base64"
	"time"

	"github.com/skip2/go-qrcode"
	"zpwoot/internal/adapters/waclient"
	"zpwoot/internal/domain/session"
)

// ProxySettings represents proxy configuration
type ProxySettings struct {
	Enabled bool   `json:"enabled" example:"true" description:"Enable proxy"`
	Type    string `json:"type,omitempty" example:"http" enums:"http,https,socks5" description:"Proxy type (http, https, socks5)"`
	Host    string `json:"host,omitempty" example:"proxy.example.com" description:"Proxy host"`
	Port    string `json:"port,omitempty" example:"8080" description:"Proxy port"`
	User    string `json:"user,omitempty" example:"proxyUser123" description:"Proxy username (optional)"`
	Pass    string `json:"pass,omitempty" example:"proxyPass123" description:"Proxy password (optional)"`
}

// WebhookSettings represents webhook configuration
type WebhookSettings struct {
	Enabled bool     `json:"enabled" example:"true" description:"Enable webhook"`
	URL     string   `json:"url,omitempty" example:"https://api.example.com/webhook" validate:"omitempty,url" description:"Webhook URL"`
	Events  []string `json:"events,omitempty" example:"Message,Receipt,Connected" description:"Events to subscribe (Message, Receipt, Connected, Disconnected, CallOffer, Presence, NewsletterJoin, All)"`
	Secret  string   `json:"secret,omitempty" example:"supersecrettoken123" description:"Webhook secret for validation (optional)"`
}

// SessionSettings represents session configuration
type SessionSettings struct {
	Proxy   *ProxySettings   `json:"proxy,omitempty" description:"Proxy configuration"`
	Webhook *WebhookSettings `json:"webhook,omitempty" description:"Webhook configuration"`
}

// CreateSessionRequest represents a session creation request
type CreateSessionRequest struct {
	Name           string           `json:"name" example:"my-session" validate:"required,min=1,max=100" description:"Session name for identification"`
	Settings       *SessionSettings `json:"settings,omitempty" description:"Session settings (proxy, webhook)"`
	GenerateQRCode bool             `json:"qrCode,omitempty" example:"true" description:"Auto-generate QR code after creation (default: false)"`
}

// UpdateSessionRequest represents a session update request
type UpdateSessionRequest struct {
	Name     *string          `json:"name,omitempty" example:"updated-session" validate:"omitempty,min=1,max=100" description:"Session name for identification"`
	Settings *SessionSettings `json:"settings,omitempty" description:"Session settings (proxy, webhook)"`
}

// SessionResponse represents a unified session information response
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
}

// SessionListInfo represents session information for list/info endpoints (excludes QR code)
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
}

// SessionListResponse represents a list of sessions response
type SessionListResponse struct {
	Sessions []SessionListInfo `json:"sessions" description:"List of sessions (without QR codes)"`
	Total    int               `json:"total" example:"5" description:"Total number of sessions"`
}

// SessionActionResponse represents a session action response (connect, disconnect, logout, delete)
type SessionActionResponse struct {
	SessionID string `json:"sessionId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Session ID"`
	Action    string `json:"action" example:"connect" description:"Action performed (connect, disconnect, logout, delete)"`
	Status    string `json:"status" example:"success" description:"Action status"`
	Message   string `json:"message,omitempty" example:"Session connected successfully" description:"Action message"`
}

// QRCodeResponse represents a QR code response
type QRCodeResponse struct {
	QRCode       string `json:"qrCode" example:"2@abc123..." description:"QR code string (original from WhatsApp)"`
	QRCodeBase64 string `json:"qrCodeBase64" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..." description:"Base64 encoded QR code image"`
	ExpiresAt    string `json:"expiresAt" example:"2025-01-15T10:35:00Z" description:"QR code expiration time"`
	Status       string `json:"status" example:"generated" description:"QR code status"`
}

// Legacy DTOs for backward compatibility with use cases

// CreateSessionResponse represents a session creation response (legacy)
type CreateSessionResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" description:"Session identifier"`
	Name      string    `json:"name" example:"My WhatsApp Session" description:"Session name"`
	Status    string    `json:"status" example:"disconnected" description:"Initial session status"`
	Connected bool      `json:"connected" example:"false" description:"Whether session is connected"`
	CreatedAt time.Time `json:"createdAt" example:"2025-01-15T10:30:00Z" description:"Session creation timestamp"`
}

// SessionDetailResponse represents detailed session information (legacy)
type SessionDetailResponse struct {
	ID              string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" description:"Session identifier"`
	Name            string     `json:"name" example:"My WhatsApp Session" description:"Session name"`
	DeviceJID       string     `json:"deviceJid,omitempty" example:"5511999999999@s.whatsapp.net" description:"WhatsApp device JID when connected"`
	Status          string     `json:"status" example:"connected" description:"Current session status"`
	Connected       bool       `json:"connected" example:"true" description:"Whether session is connected"`
	ConnectionError string     `json:"connectionError,omitempty" example:"Connection timeout" description:"Connection error message if any"`
	QRCode          string     `json:"qrCode,omitempty" description:"QR code for authentication"`
	QRCodeExpiresAt *time.Time `json:"qrCodeExpiresAt,omitempty" example:"2025-01-15T10:35:00Z" description:"QR code expiration time"`
	ProxyConfig     string     `json:"proxyConfig,omitempty" description:"Proxy configuration as JSON string"`
	CreatedAt       time.Time  `json:"createdAt" example:"2025-01-15T10:30:00Z" description:"Session creation timestamp"`
	UpdatedAt       time.Time  `json:"updatedAt" example:"2025-01-15T10:35:00Z" description:"Last update timestamp"`
	ConnectedAt     *time.Time `json:"connectedAt,omitempty" example:"2025-01-15T10:32:00Z" description:"Connection timestamp"`
	LastSeen        *time.Time `json:"lastSeen,omitempty" example:"2025-01-15T10:35:00Z" description:"Last activity timestamp"`
}

// SessionStatusResponse represents session status (legacy)
type SessionStatusResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" description:"Session identifier"`
	Status    string `json:"status" example:"connected" description:"Current session status"`
	Connected bool   `json:"connected" example:"true" description:"Whether session is connected"`
}

// Validation methods

// Validate validates the CreateSessionRequest
func (r *CreateSessionRequest) Validate() error {
	if r.Name == "" {
		return NewValidationError("name", "session name is required")
	}
	if len(r.Name) > 100 {
		return NewValidationError("name", "session name must be less than 100 characters")
	}

	// Validate webhook settings
	if r.Settings != nil && r.Settings.Webhook != nil && r.Settings.Webhook.Enabled {
		if r.Settings.Webhook.URL == "" {
			return NewValidationError("settings.webhook.url", "webhook URL is required when webhook is enabled")
		}
	}

	// Validate proxy settings
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

// Conversion methods

// ToSessionConfig converts CreateSessionRequest to waclient.SessionConfig
func (r *CreateSessionRequest) ToSessionConfig() *waclient.SessionConfig {
	config := &waclient.SessionConfig{
		Name:          r.Name,
		AutoReconnect: true, // Default to true
	}

	// Convert webhook settings
	if r.Settings != nil && r.Settings.Webhook != nil && r.Settings.Webhook.Enabled {
		config.WebhookURL = r.Settings.Webhook.URL

		// Convert event strings to EventType
		events := make([]waclient.EventType, 0, len(r.Settings.Webhook.Events))
		for _, event := range r.Settings.Webhook.Events {
			events = append(events, waclient.EventType(event))
		}
		config.Events = events
	}

	// Convert proxy settings to map
	if r.Settings != nil && r.Settings.Proxy != nil && r.Settings.Proxy.Enabled {
		config.ProxyConfig = map[string]string{
			"enabled": "true",
			"type":    r.Settings.Proxy.Type,
			"host":    r.Settings.Proxy.Host,
			"port":    r.Settings.Proxy.Port,
		}
		if r.Settings.Proxy.User != "" {
			config.ProxyConfig["user"] = r.Settings.Proxy.User
		}
		if r.Settings.Proxy.Pass != "" {
			config.ProxyConfig["pass"] = r.Settings.Proxy.Pass
		}
	}

	return config
}

// ToDomain converts CreateSessionRequest to domain.Session
func (r *CreateSessionRequest) ToDomain() *session.Session {
	return session.NewSession(r.Name)
}

// FromDomainSession converts domain.Session to SessionResponse
func FromDomainSession(s *session.Session) *SessionResponse {
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

	// Generate base64 and set expiration if QR code exists
	if s.QRCode != "" {
		response.QRCodeBase64 = GenerateQRCodeBase64(s.QRCode)
		if s.QRCodeExpiresAt != nil {
			response.QRCodeExpiresAt = s.QRCodeExpiresAt
		}
	}

	// Note: Settings should be populated from session config if available
	// This is a basic conversion from domain model

	return response
}

// FromWAClient converts waclient.Client to SessionResponse
func FromWAClient(client *waclient.Client) *SessionResponse {
	deviceJID := ""
	if client.WAClient.Store.ID != nil {
		deviceJID = client.WAClient.Store.ID.String()
	}

	response := &SessionResponse{
		SessionID: client.SessionID,
		Name:      client.Name,
		Status:    string(client.Status),
		Connected: client.Status == waclient.StatusConnected,
		DeviceJID: deviceJID,
		QRCode:    client.QRCode, // Original QR string
		CreatedAt: time.Now(),    // Should be stored in client
		UpdatedAt: time.Now(),
	}

	// Generate base64 QR code if QR string exists
	if client.QRCode != "" {
		response.QRCodeBase64 = GenerateQRCodeBase64(client.QRCode)
		if !client.QRExpiresAt.IsZero() {
			response.QRCodeExpiresAt = &client.QRExpiresAt
		}
	}

	// Build settings from client config
	if client.Config != nil {
		settings := &SessionSettings{}

		// Webhook settings
		if client.WebhookURL != "" || len(client.Events) > 0 {
			events := make([]string, len(client.Events))
			for i, event := range client.Events {
				events[i] = string(event)
			}

			settings.Webhook = &WebhookSettings{
				Enabled: client.WebhookURL != "",
				URL:     client.WebhookURL,
				Events:  events,
			}
		}

		// Proxy settings
		if client.Config.ProxyConfig != nil && len(client.Config.ProxyConfig) > 0 {
			if enabled, ok := client.Config.ProxyConfig["enabled"]; ok && enabled == "true" {
				settings.Proxy = &ProxySettings{
					Enabled: true,
					Type:    client.Config.ProxyConfig["type"],
					Host:    client.Config.ProxyConfig["host"],
					Port:    client.Config.ProxyConfig["port"],
					User:    client.Config.ProxyConfig["user"],
				}
			}
		}

		if settings.Webhook != nil || settings.Proxy != nil {
			response.Settings = settings
		}
	}

	if !client.QRExpiresAt.IsZero() {
		response.QRCodeExpiresAt = &client.QRExpiresAt
	}
	if !client.ConnectedAt.IsZero() {
		response.ConnectedAt = &client.ConnectedAt
	}
	if !client.LastSeen.IsZero() {
		response.LastSeen = &client.LastSeen
	}

	return response
}

// ToListInfo converts SessionResponse to SessionListInfo (removes QR code)
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

// FromWAClientList converts slice of waclient.Client to SessionListResponse (without QR codes)
func FromWAClientList(clients []*waclient.Client) *SessionListResponse {
	sessions := make([]SessionListInfo, len(clients))
	for i, client := range clients {
		sessionResp := FromWAClient(client)
		sessions[i] = *sessionResp.ToListInfo() // Remove QR code
	}

	return &SessionListResponse{
		Sessions: sessions,
		Total:    len(sessions),
	}
}

// FromQREvent converts waclient.QREvent to QRCodeResponse
func FromQREvent(qrEvent *waclient.QREvent) *QRCodeResponse {
	return &QRCodeResponse{
		QRCode:       qrEvent.Code,                       // Original QR string
		QRCodeBase64: GenerateQRCodeBase64(qrEvent.Code), // Base64 image
		ExpiresAt:    qrEvent.ExpiresAt.Format(time.RFC3339),
		Status:       "generated",
	}
}

// NewQRCodeResponse creates a QRCodeResponse with proper time formatting
func NewQRCodeResponse(qrCode string, expiresAt time.Time, status string) *QRCodeResponse {
	expiresAtStr := ""
	if !expiresAt.IsZero() {
		expiresAtStr = expiresAt.Format(time.RFC3339)
	}

	return &QRCodeResponse{
		QRCode:       qrCode,                       // Original QR string
		QRCodeBase64: GenerateQRCodeBase64(qrCode), // Base64 image
		ExpiresAt:    expiresAtStr,
		Status:       status,
	}
}

// Legacy conversion functions for backward compatibility

// SessionToDetailResponse converts domain.Session to SessionDetailResponse (legacy)
func SessionToDetailResponse(s *session.Session) *SessionDetailResponse {
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

// SessionToCreateResponse converts domain.Session to CreateSessionResponse (legacy)
func SessionToCreateResponse(s *session.Session) *CreateSessionResponse {
	return &CreateSessionResponse{
		ID:        s.ID,
		Name:      s.Name,
		Status:    string(s.GetStatus()),
		Connected: s.IsConnected,
		CreatedAt: s.CreatedAt,
	}
}

// SessionToStatusResponse converts domain.Session to SessionStatusResponse (legacy)
func SessionToStatusResponse(s *session.Session) *SessionStatusResponse {
	return &SessionStatusResponse{
		ID:        s.ID,
		Status:    string(s.GetStatus()),
		Connected: s.IsConnected,
	}
}

// SessionToListResponse converts domain.Session to SessionResponse for list view (legacy)
func SessionToListResponse(s *session.Session) *SessionResponse {
	return FromDomainSession(s)
}

// GenerateQRCodeBase64 converts QR code string to base64 image
func GenerateQRCodeBase64(qrString string) string {
	if qrString == "" {
		return ""
	}

	// Generate QR code image
	qrImage, err := qrcode.Encode(qrString, qrcode.Medium, 256)
	if err != nil {
		return "" // Return empty if encoding fails
	}

	// Convert to base64
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrImage)
}

// Error definitions for sessions
var (
	ErrInvalidSessionName   = NewValidationError("name", "Session name is required")
	ErrSessionNameTooLong   = NewValidationError("name", "Session name must be less than 100 characters")
	ErrSessionNotFound      = &ErrorInfo{Code: "SESSION_NOT_FOUND", Message: "Session not found"}
	ErrSessionAlreadyExists = &ErrorInfo{Code: "SESSION_ALREADY_EXISTS", Message: "Session already exists"}
)
