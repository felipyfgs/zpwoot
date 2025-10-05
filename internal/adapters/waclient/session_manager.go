package waclient

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zpwoot/internal/domain/shared"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBSessionManager struct {
	db *sqlx.DB
}

func NewDBSessionManager(db *sqlx.DB) *DBSessionManager {
	return &DBSessionManager{db: db}
}

func scanSessionRow(scanner interface {
	Scan(dest ...interface{}) error
}) (*SessionInfo, error) {
	var session SessionInfo
	var deviceJid, connectionError, qrCode sql.NullString
	var qrExpiresAt, connectedAt, lastSeen sql.NullTime
	var proxyConfig sql.NullString

	err := scanner.Scan(
		&session.ID,
		&session.Name,
		&deviceJid,
		&session.Connected,
		&connectionError,
		&qrCode,
		&qrExpiresAt,
		&proxyConfig,
		&session.CreatedAt,
		&session.UpdatedAt,
		&connectedAt,
		&lastSeen,
	)
	if err != nil {
		return nil, err
	}

	if deviceJid.Valid {
		session.DeviceJID = deviceJid.String
	}
	if qrCode.Valid {
		session.QRCode = qrCode.String
	}
	if qrExpiresAt.Valid {
		session.QRExpiresAt = qrExpiresAt.Time
	}
	if connectedAt.Valid {
		session.ConnectedAt = connectedAt.Time
	}
	if lastSeen.Valid {
		session.LastSeen = lastSeen.Time
	}

	if session.Connected {
		session.Status = StatusConnected
	} else {
		session.Status = StatusDisconnected
	}

	return &session, nil
}

func toNullableFields(session *SessionInfo) (deviceJid, qrCode sql.NullString, qrExpiresAt, connectedAt, lastSeen sql.NullTime) {
	if session.DeviceJID != "" {
		deviceJid = sql.NullString{String: session.DeviceJID, Valid: true}
	}
	if session.QRCode != "" {
		qrCode = sql.NullString{String: session.QRCode, Valid: true}
	}
	if !session.QRExpiresAt.IsZero() {
		qrExpiresAt = sql.NullTime{Time: session.QRExpiresAt, Valid: true}
	}
	if !session.ConnectedAt.IsZero() {
		connectedAt = sql.NullTime{Time: session.ConnectedAt, Valid: true}
	}
	if !session.LastSeen.IsZero() {
		lastSeen = sql.NullTime{Time: session.LastSeen, Valid: true}
	}
	return
}

func (sm *DBSessionManager) GetSession(ctx context.Context, sessionID string) (*SessionInfo, error) {
	query := `
		SELECT "id", "name", "deviceJid", "isConnected", "connectionError",
			   "qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt",
			   "updatedAt", "connectedAt", "lastSeen"
		FROM "zpSessions"
		WHERE "id" = $1
	`

	row := sm.db.QueryRowContext(ctx, query, sessionID)
	session, err := scanSessionRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (sm *DBSessionManager) GetSessionByName(ctx context.Context, name string) (*SessionInfo, error) {
	query := `
		SELECT "id", "name", "deviceJid", "isConnected", "connectionError",
			   "qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt",
			   "updatedAt", "connectedAt", "lastSeen"
		FROM "zpSessions"
		WHERE "name" = $1
	`

	row := sm.db.QueryRowContext(ctx, query, name)
	session, err := scanSessionRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by name: %w", err)
	}

	return session, nil
}

func (sm *DBSessionManager) CreateSession(ctx context.Context, session *SessionInfo) error {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	query := `
		INSERT INTO "zpSessions" (
			"id", "name", "deviceJid", "isConnected", "connectionError",
			"qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt",
			"updatedAt", "connectedAt", "lastSeen"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`

	deviceJid, qrCode, qrExpiresAt, connectedAt, lastSeen := toNullableFields(session)

	_, err := sm.db.ExecContext(ctx, query,
		session.ID,
		session.Name,
		deviceJid,
		session.Connected,
		nil,
		qrCode,
		qrExpiresAt,
		nil,
		session.CreatedAt,
		session.UpdatedAt,
		connectedAt,
		lastSeen,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (sm *DBSessionManager) UpdateSession(ctx context.Context, session *SessionInfo) error {
	query := `
		UPDATE "zpSessions" SET
			"name" = $2,
			"deviceJid" = $3,
			"isConnected" = $4,
			"connectionError" = $5,
			"qrCode" = $6,
			"qrCodeExpiresAt" = $7,
			"updatedAt" = $8,
			"connectedAt" = $9,
			"lastSeen" = $10
		WHERE "id" = $1
	`

	deviceJid, qrCode, qrExpiresAt, connectedAt, lastSeen := toNullableFields(session)

	result, err := sm.db.ExecContext(ctx, query,
		session.ID,
		session.Name,
		deviceJid,
		session.Connected,
		nil,
		qrCode,
		qrExpiresAt,
		time.Now(),
		connectedAt,
		lastSeen,
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

func (sm *DBSessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	result, err := sm.db.ExecContext(ctx, `DELETE FROM "zpSessions" WHERE "id" = $1`, sessionID)
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

func (sm *DBSessionManager) ListSessions(ctx context.Context) ([]*SessionInfo, error) {
	query := `
		SELECT "id", "name", "deviceJid", "isConnected", "connectionError",
			   "qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt",
			   "updatedAt", "connectedAt", "lastSeen"
		FROM "zpSessions"
		ORDER BY "createdAt" DESC
	`

	rows, err := sm.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*SessionInfo
	for rows.Next() {
		session, err := scanSessionRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sessions: %w", err)
	}

	return sessions, nil
}
