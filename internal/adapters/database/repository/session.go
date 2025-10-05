package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/shared"

	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

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

func (r *SessionRepository) GetByJID(ctx context.Context, jid string) (*session.Session, error) {
	query := `
		SELECT "id", "name", "deviceJid", "isConnected", "connectionError",
			   "qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt",
			   "updatedAt", "connectedAt", "lastSeen"
		FROM "zpSessions"
		WHERE "deviceJid" = $1
	`

	var sess session.Session
	err := r.db.GetContext(ctx, &sess, query, jid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by JID: %w", err)
	}

	return &sess, nil
}

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

func (r *SessionRepository) List(ctx context.Context, limit, offset int) ([]*session.Session, error) {
	query := `
		SELECT "id", "name", "deviceJid", "isConnected", "connectionError",
			   "qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt",
			   "updatedAt", "connectedAt", "lastSeen"
		FROM "zpSessions"
		ORDER BY "createdAt" DESC
		LIMIT $1 OFFSET $2
	`

	var sessions []*session.Session
	err := r.db.SelectContext(ctx, &sessions, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	return sessions, nil
}

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

func (r *SessionRepository) UpdateQRCode(ctx context.Context, id string, qrCode string) error {
	query := `
		UPDATE "zpSessions"
		SET "qrCode" = $2, "updatedAt" = NOW()
		WHERE "id" = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, qrCode)
	if err != nil {
		return fmt.Errorf("failed to update session QR code: %w", err)
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
