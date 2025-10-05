package session

import (
	"context"
)

// Repository defines the interface for session data access
type Repository interface {
	// Create creates a new session
	Create(ctx context.Context, session *Session) error

	// GetByID retrieves a session by ID
	GetByID(ctx context.Context, id string) (*Session, error)

	// GetByJID retrieves a session by JID
	GetByJID(ctx context.Context, jid string) (*Session, error)

	// Update updates an existing session
	Update(ctx context.Context, session *Session) error

	// Delete deletes a session by ID
	Delete(ctx context.Context, id string) error

	// List retrieves all sessions with pagination
	List(ctx context.Context, limit, offset int) ([]*Session, error)

	// UpdateStatus updates the session status
	UpdateStatus(ctx context.Context, id string, status Status) error

	// UpdateQRCode updates the session QR code
	UpdateQRCode(ctx context.Context, id string, qrCode string) error
}
