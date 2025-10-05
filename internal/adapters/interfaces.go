package adapters

import (
	"context"
	"github.com/zpwoot/internal/domain/session"
)

// WhatsAppClient interface for WhatsApp operations
type WhatsAppClient interface {
	CreateSession(ctx context.Context, sessionID string) error
	ConnectSession(ctx context.Context, sessionID string) (string, error)
	DisconnectSession(sessionID string) error
	GetSessionStatus(sessionID string) (session.Status, error)
	GetSessionJID(sessionID string) (string, error)
	RemoveSession(sessionID string) error
	SendMessage(ctx context.Context, sessionID, to, message string) error
	ListActiveSessions() []string
	Close() error
}

// DatabaseAdapter interface for database operations
type DatabaseAdapter interface {
	Connect() error
	Close() error
	Migrate() error
	Health() error
}
