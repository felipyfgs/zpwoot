package session

import (
	"time"
)


type Session struct {
	ID              string     `json:"id" db:"id"`
	Name            string     `json:"name" db:"name"`
	DeviceJID       string     `json:"device_jid,omitempty" db:"deviceJid"`
	IsConnected     bool       `json:"is_connected" db:"isConnected"`
	ConnectionError string     `json:"connection_error,omitempty" db:"connectionError"`
	QRCode          string     `json:"qr_code,omitempty" db:"qrCode"`
	QRCodeExpiresAt *time.Time `json:"qr_code_expires_at,omitempty" db:"qrCodeExpiresAt"`
	ProxyConfig     string     `json:"proxy_config,omitempty" db:"proxyConfig"`
	CreatedAt       time.Time  `json:"created_at" db:"createdAt"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updatedAt"`
	ConnectedAt     *time.Time `json:"connected_at,omitempty" db:"connectedAt"`
	LastSeen        *time.Time `json:"last_seen,omitempty" db:"lastSeen"`
}


type Status string

const (
	StatusDisconnected Status = "disconnected"
	StatusConnecting   Status = "connecting"
	StatusConnected    Status = "connected"
	StatusQRCode       Status = "qr_code"
	StatusError        Status = "error"
)


func (s Status) IsValid() bool {
	switch s {
	case StatusDisconnected, StatusConnecting, StatusConnected, StatusQRCode, StatusError:
		return true
	default:
		return false
	}
}


func NewSession(name string) *Session {
	now := time.Now()
	return &Session{
		Name:        name,
		IsConnected: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}


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


func (s *Session) SetQRCode(qrCode string, expiresAt time.Time) {
	s.QRCode = qrCode
	s.QRCodeExpiresAt = &expiresAt
	s.UpdatedAt = time.Now()
}


func (s *Session) ClearQRCode() {
	s.QRCode = ""
	s.QRCodeExpiresAt = nil
	s.UpdatedAt = time.Now()
}


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


func (s *Session) SetDisconnected() {
	s.IsConnected = false
	now := time.Now()
	s.LastSeen = &now
	s.UpdatedAt = now
	s.ClearQRCode()
}


func (s *Session) SetError(error string) {
	s.IsConnected = false
	s.ConnectionError = error
	now := time.Now()
	s.LastSeen = &now
	s.UpdatedAt = now
	s.ClearQRCode()
}


func (s *Session) UpdateLastSeen() {
	now := time.Now()
	s.LastSeen = &now
	s.UpdatedAt = now
}
