package session

import (
	"time"
)

// Session represents a WhatsApp session
type Session struct {
	ID          string    `json:"id" db:"id"`
	JID         string    `json:"jid" db:"jid"`
	Name        string    `json:"name" db:"name"`
	Status      Status    `json:"status" db:"status"`
	QRCode      string    `json:"qr_code,omitempty" db:"qr_code"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	ConnectedAt *time.Time `json:"connected_at,omitempty" db:"connected_at"`
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
	return &Session{
		Name:      name,
		Status:    StatusDisconnected,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
