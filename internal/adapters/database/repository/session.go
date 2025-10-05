package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zpwoot/internal/domain/session"
	"zpwoot/internal/domain/shared"

	"github.com/jmoiron/sqlx"
)

// SessionRepository implements the session.Repository interface
type SessionRepository struct {
	db *sqlx.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

// Create creates a new session
func (r *SessionRepository) Create(ctx context.Context, sess *session.Session) error {
	query := `
		INSERT INTO "zpSessions" (
			"id", "name", "deviceJid", "isConnected", "connectionError", 
			"qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt", 
			"updatedAt", "connectedAt", "lastSeen"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		sess.ID,
		sess.Name,
		sess.DeviceJID,
		sess.IsConnected,
		sess.ConnectionError,
		sess.QRCode,
		sess.QRCodeExpiresAt,
		sess.ProxyConfig,
		sess.CreatedAt,
		sess.UpdatedAt,
		sess.ConnectedAt,
		sess.LastSeen,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetByID retrieves a session by ID
func (r *SessionRepository) GetByID(ctx context.Context, id string) (*session.Session, error) {
	query := `
		SELECT "id", "name", "deviceJid", "isConnected", "connectionError", 
			   "qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt", 
			   "updatedAt", "connectedAt", "lastSeen"
		FROM "zpSessions" 
		WHERE "id" = $1
	`

	var sess session.Session
	err := r.db.GetContext(ctx, &sess, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &sess, nil
}

// GetByName retrieves a session by name
func (r *SessionRepository) GetByName(ctx context.Context, name string) (*session.Session, error) {
	query := `
		SELECT "id", "name", "deviceJid", "isConnected", "connectionError", 
			   "qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt", 
			   "updatedAt", "connectedAt", "lastSeen"
		FROM "zpSessions" 
		WHERE "name" = $1
	`

	var sess session.Session
	err := r.db.GetContext(ctx, &sess, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by name: %w", err)
	}

	return &sess, nil
}

// List retrieves all sessions
func (r *SessionRepository) List(ctx context.Context) ([]*session.Session, error) {
	query := `
		SELECT "id", "name", "deviceJid", "isConnected", "connectionError", 
			   "qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt", 
			   "updatedAt", "connectedAt", "lastSeen"
		FROM "zpSessions" 
		ORDER BY "createdAt" DESC
	`

	var sessions []*session.Session
	err := r.db.SelectContext(ctx, &sessions, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	return sessions, nil
}

// Update updates a session
func (r *SessionRepository) Update(ctx context.Context, sess *session.Session) error {
	query := `
		UPDATE "zpSessions" SET
			"name" = $2,
			"deviceJid" = $3,
			"isConnected" = $4,
			"connectionError" = $5,
			"qrCode" = $6,
			"qrCodeExpiresAt" = $7,
			"proxyConfig" = $8,
			"updatedAt" = $9,
			"connectedAt" = $10,
			"lastSeen" = $11
		WHERE "id" = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		sess.ID,
		sess.Name,
		sess.DeviceJID,
		sess.IsConnected,
		sess.ConnectionError,
		sess.QRCode,
		sess.QRCodeExpiresAt,
		sess.ProxyConfig,
		time.Now(),
		sess.ConnectedAt,
		sess.LastSeen,
	)

	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return shared.ErrSessionNotFound
	}

	return nil
}

// UpdateStatus updates only the session status
func (r *SessionRepository) UpdateStatus(ctx context.Context, id string, status session.Status) error {
	query := `
		UPDATE "zpSessions" SET
			"isConnected" = $2,
			"updatedAt" = NOW()
		WHERE "id" = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, status == session.StatusConnected)
	if err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return shared.ErrSessionNotFound
	}

	return nil
}

// Delete deletes a session
func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "zpSessions" WHERE "id" = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return shared.ErrSessionNotFound
	}

	return nil
}
