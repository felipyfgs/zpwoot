package session

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a WhatsApp session entity
type Session struct {
	id               uuid.UUID
	name             string
	deviceJID        *string
	isConnected      bool
	connectionError  *string
	qrCode           *string
	qrCodeExpiresAt  *time.Time
	proxyConfig      *ProxyConfig
	createdAt        time.Time
	updatedAt        time.Time
	connectedAt      *time.Time
	lastSeen         *time.Time
}

// ProxyConfig represents proxy configuration
type ProxyConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// NewSession creates a new session
func NewSession(name string, proxyConfig *ProxyConfig) (*Session, error) {
	if name == "" {
		return nil, ErrInvalidSessionName
	}

	now := time.Now()
	return &Session{
		id:          uuid.New(),
		name:        name,
		isConnected: false,
		proxyConfig: proxyConfig,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// RestoreSession creates a session from persisted data (for repository use)
func RestoreSession(id uuid.UUID, name string, isConnected bool, deviceJID *string,
	connectionError *string, qrCode *string, qrCodeExpiresAt *time.Time,
	proxyConfig *ProxyConfig, createdAt, updatedAt time.Time,
	connectedAt, lastSeen *time.Time) *Session {

	return &Session{
		id:               id,
		name:             name,
		isConnected:      isConnected,
		deviceJID:        deviceJID,
		connectionError:  connectionError,
		qrCode:           qrCode,
		qrCodeExpiresAt:  qrCodeExpiresAt,
		proxyConfig:      proxyConfig,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
		connectedAt:      connectedAt,
		lastSeen:         lastSeen,
	}
}

// Business methods

// Connect marks session as connected
func (s *Session) Connect(deviceJID string) {
	s.isConnected = true
	s.deviceJID = &deviceJID
	now := time.Now()
	s.connectedAt = &now
	s.updatedAt = now
	s.connectionError = nil
}

// Disconnect marks session as disconnected
func (s *Session) Disconnect(reason string) {
	s.isConnected = false
	s.deviceJID = nil
	s.connectedAt = nil
	s.updatedAt = time.Now()
	if reason != "" {
		s.connectionError = &reason
	}
}

// SetQRCode sets QR code for pairing
func (s *Session) SetQRCode(qrCode string, expiresAt time.Time) {
	s.qrCode = &qrCode
	s.qrCodeExpiresAt = &expiresAt
	s.updatedAt = time.Now()
}

// ClearQRCode removes QR code
func (s *Session) ClearQRCode() {
	s.qrCode = nil
	s.qrCodeExpiresAt = nil
	s.updatedAt = time.Now()
}

// UpdateLastSeen updates last activity
func (s *Session) UpdateLastSeen() {
	now := time.Now()
	s.lastSeen = &now
	s.updatedAt = now
}

// Getters

func (s *Session) ID() uuid.UUID {
	return s.id
}

func (s *Session) Name() string {
	return s.name
}

func (s *Session) DeviceJID() *string {
	return s.deviceJID
}

func (s *Session) IsConnected() bool {
	return s.isConnected
}

func (s *Session) ConnectionError() *string {
	return s.connectionError
}

func (s *Session) QRCode() *string {
	return s.qrCode
}

func (s *Session) QRCodeExpiresAt() *time.Time {
	return s.qrCodeExpiresAt
}

func (s *Session) ProxyConfig() *ProxyConfig {
	return s.proxyConfig
}

func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Session) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *Session) ConnectedAt() *time.Time {
	return s.connectedAt
}

func (s *Session) LastSeen() *time.Time {
	return s.lastSeen
}

// Domain queries

// HasValidQRCode checks if QR code is valid
func (s *Session) HasValidQRCode() bool {
	if s.qrCode == nil || s.qrCodeExpiresAt == nil {
		return false
	}
	return time.Now().Before(*s.qrCodeExpiresAt)
}

// CanConnect checks if session can connect
func (s *Session) CanConnect() bool {
	return !s.isConnected
}

// HasProxy checks if session has proxy config
func (s *Session) HasProxy() bool {
	return s.proxyConfig != nil
}

// IsActive checks if session had recent activity
func (s *Session) IsActive() bool {
	if !s.isConnected {
		return false
	}
	if s.lastSeen == nil {
		return true // Just connected
	}
	// Active if seen within last 5 minutes
	return time.Since(*s.lastSeen) < 5*time.Minute
}
