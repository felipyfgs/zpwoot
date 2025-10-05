package session

import (
	"time"
)

// Session represents a WhatsApp session
type Session struct {
	ID                string     `json:"id" db:"id"`
	Name              string     `json:"name" db:"name"`
	DeviceJID         string     `json:"device_jid,omitempty" db:"deviceJid"`
	IsConnected       bool       `json:"is_connected" db:"isConnected"`
	ConnectionError   string     `json:"connection_error,omitempty" db:"connectionError"`
	QRCode            string     `json:"qr_code,omitempty" db:"qrCode"`
	QRCodeExpiresAt   *time.Time `json:"qr_code_expires_at,omitempty" db:"qrCodeExpiresAt"`
	ProxyConfig       string     `json:"proxy_config,omitempty" db:"proxyConfig"` // JSON string
	CreatedAt         time.Time  `json:"created_at" db:"createdAt"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updatedAt"`
	ConnectedAt       *time.Time `json:"connected_at,omitempty" db:"connectedAt"`
	LastSeen          *time.Time `json:"last_seen,omitempty" db:"lastSeen"`
}

// Status represents the session status
type Status string

const (
	StatusDisconnected Status = "disconnected"
	StatusConnecting   Status = "connecting"
	StatusConnected    Status = "connected"
	StatusQRCode       Status = "qr_code"
	StatusError        Status = "error"
)

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusDisconnected, StatusConnecting, StatusConnected, StatusQRCode, StatusError:
		return true
	default:
		return false
	}
}

// NewSession creates a new session
func NewSession(name string) *Session {
	now := time.Now()
	return &Session{
		Name:        name,
		IsConnected: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GetStatus returns the session status based on connection state and QR code
func (s *Session) GetStatus() Status {
	if s.IsConnected {
		return StatusConnected
	}
	if s.QRCode != "" && (s.QRCodeExpiresAt == nil || time.Now().Before(*s.QRCodeExpiresAt)) {
		return StatusQRCode
	}
	if s.ConnectionError != "" {
		return StatusError
	}
	return StatusDisconnected
}

// SetQRCode sets the QR code with expiration time
func (s *Session) SetQRCode(qrCode string, expiresAt time.Time) {
	s.QRCode = qrCode
	s.QRCodeExpiresAt = &expiresAt
	s.UpdatedAt = time.Now()
}

// ClearQRCode clears the QR code
func (s *Session) ClearQRCode() {
	s.QRCode = ""
	s.QRCodeExpiresAt = nil
	s.UpdatedAt = time.Now()
}

// SetConnected marks the session as connected
func (s *Session) SetConnected(deviceJID string) {
	s.IsConnected = true
	s.DeviceJID = deviceJID
	s.ConnectionError = ""
	now := time.Now()
	s.ConnectedAt = &now
	s.LastSeen = &now
	s.UpdatedAt = now
	s.ClearQRCode()
}

// SetDisconnected marks the session as disconnected
func (s *Session) SetDisconnected() {
	s.IsConnected = false
	now := time.Now()
	s.LastSeen = &now
	s.UpdatedAt = now
	s.ClearQRCode()
}

// SetError sets a connection error
func (s *Session) SetError(error string) {
	s.IsConnected = false
	s.ConnectionError = error
	now := time.Now()
	s.LastSeen = &now
	s.UpdatedAt = now
	s.ClearQRCode()
}

// UpdateLastSeen updates the last seen timestamp
func (s *Session) UpdateLastSeen() {
	now := time.Now()
	s.LastSeen = &now
	s.UpdatedAt = now
}
