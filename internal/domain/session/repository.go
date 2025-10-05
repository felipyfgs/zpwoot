package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines session persistence interface
type Repository interface {
	// Basic operations
	Save(ctx context.Context, session *Session) error
	FindByID(ctx context.Context, id uuid.UUID) (*Session, error)
	FindByName(ctx context.Context, name string) (*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	FindAll(ctx context.Context) ([]*Session, error)
	FindConnected(ctx context.Context) ([]*Session, error)
	FindDisconnected(ctx context.Context) ([]*Session, error)

	// Existence checks
	ExistsByName(ctx context.Context, name string) (bool, error)

	// Statistics
	Count(ctx context.Context) (int64, error)
	CountConnected(ctx context.Context) (int64, error)
}

// WhatsAppClient defines WhatsApp operations interface
type WhatsAppClient interface {
	// Connection management
	Connect(ctx context.Context, sessionID uuid.UUID, proxyConfig *ProxyConfig) error
	Disconnect(ctx context.Context, sessionID uuid.UUID) error
	IsConnected(ctx context.Context, sessionID uuid.UUID) (bool, error)

	// QR Code operations
	GenerateQRCode(ctx context.Context, sessionID uuid.UUID) (qrCode string, expiresAt time.Time, err error)
	GetQRCode(ctx context.Context, sessionID uuid.UUID) (qrCode string, expiresAt time.Time, err error)

	// Device operations
	GetDeviceJID(ctx context.Context, sessionID uuid.UUID) (string, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error

	// Status operations
	Ping(ctx context.Context, sessionID uuid.UUID) error
}
