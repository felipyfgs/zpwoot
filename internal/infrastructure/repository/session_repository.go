package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"zpwoot/internal/domain/session"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) session.Repository {
	return &SessionRepository{
		db: db,
	}
}

type sessionModel struct {
	ID              string         `db:"id"`
	Name            string         `db:"name"`
	DeviceJID       sql.NullString `db:"deviceJid"`
	IsConnected     bool           `db:"isConnected"`
	ConnectionError sql.NullString `db:"connectionError"`
	QRCode          sql.NullString `db:"qrCode"`
	QRCodeExpiresAt sql.NullTime   `db:"qrCodeExpiresAt"`
	ProxyConfig     sql.NullString `db:"proxyConfig"`
	CreatedAt       time.Time      `db:"createdAt"`
	UpdatedAt       time.Time      `db:"updatedAt"`
	ConnectedAt     sql.NullTime   `db:"connectedAt"`
	LastSeen        sql.NullTime   `db:"lastSeen"`
}

func (r *SessionRepository) Create(ctx context.Context, sess *session.Session) error {
	model, err := r.toModel(sess)
	if err != nil {
		return fmt.Errorf("failed to convert session to model: %w", err)
	}

	query := `
		INSERT INTO "zpSessions" (
			id, name, "deviceJid", "isConnected", "connectionError",
			"qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt",
			"updatedAt", "connectedAt", "lastSeen"
		) VALUES (
			:id, :name, :deviceJid, :isConnected, :connectionError,
			:qrCode, :qrCodeExpiresAt, :proxyConfig, :createdAt,
			:updatedAt, :connectedAt, :lastSeen
		)
	`

	_, err = r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				if pqErr.Constraint == "zpSessions_name_key" {
					return session.ErrSessionAlreadyExists
				}
			}
		}
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*session.Session, error) {
	var model sessionModel
	query := `SELECT * FROM "zpSessions" WHERE id = $1`

	err := r.db.GetContext(ctx, &model, query, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, session.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by ID: %w", err)
	}

	return r.fromModel(&model)
}

func (r *SessionRepository) Update(ctx context.Context, sess *session.Session) error {
	model, err := r.toModel(sess)
	if err != nil {
		return fmt.Errorf("failed to convert session to model: %w", err)
	}

	query := `
		UPDATE "zpSessions" SET
			name = :name,
			"deviceJid" = :deviceJid,
			"isConnected" = :isConnected,
			"connectionError" = :connectionError,
			"qrCode" = :qrCode,
			"qrCodeExpiresAt" = :qrCodeExpiresAt,
			"proxyConfig" = :proxyConfig,
			"updatedAt" = :updatedAt,
			"connectedAt" = :connectedAt,
			"lastSeen" = :lastSeen
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				if pqErr.Constraint == "zpSessions_name_key" {
					return session.ErrSessionAlreadyExists
				}
			}
		}
		return fmt.Errorf("failed to update session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM "zpSessions" WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

func (r *SessionRepository) List(ctx context.Context, limit, offset int) ([]*session.Session, error) {
	var models []sessionModel
	query := `
		SELECT * FROM "zpSessions"
		ORDER BY "createdAt" DESC
		LIMIT $1 OFFSET $2
	`

	err := r.db.SelectContext(ctx, &models, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		sess, err := r.fromModel(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to session: %w", err)
		}
		sessions[i] = sess
	}

	return sessions, nil
}

func (r *SessionRepository) ListConnected(ctx context.Context) ([]*session.Session, error) {
	var models []sessionModel
	query := `SELECT * FROM "zpSessions" WHERE "isConnected" = true ORDER BY "connectedAt" DESC`

	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list connected sessions: %w", err)
	}

	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		sess, err := r.fromModel(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to session: %w", err)
		}
		sessions[i] = sess
	}

	return sessions, nil
}

func (r *SessionRepository) ListByStatus(ctx context.Context, connected bool) ([]*session.Session, error) {
	var models []sessionModel
	query := `SELECT * FROM "zpSessions" WHERE "isConnected" = $1 ORDER BY "updatedAt" DESC`

	err := r.db.SelectContext(ctx, &models, query, connected)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions by status: %w", err)
	}

	sessions := make([]*session.Session, len(models))
	for i, model := range models {
		sess, err := r.fromModel(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to session: %w", err)
		}
		sessions[i] = sess
	}

	return sessions, nil
}

func (r *SessionRepository) UpdateConnectionStatus(ctx context.Context, id uuid.UUID, connected bool) error {
	query := `
		UPDATE "zpSessions" SET
			"isConnected" = $2,
			"connectedAt" = CASE WHEN $2 = true THEN NOW() ELSE "connectedAt" END,
			"lastSeen" = CASE WHEN $2 = true THEN NOW() ELSE "lastSeen" END,
			"updatedAt" = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id.String(), connected)
	if err != nil {
		return fmt.Errorf("failed to update connection status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

func (r *SessionRepository) UpdateLastSeen(ctx context.Context, id uuid.UUID, lastSeen time.Time) error {
	query := `UPDATE "zpSessions" SET "lastSeen" = $2, "updatedAt" = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String(), lastSeen)
	if err != nil {
		return fmt.Errorf("failed to update last seen: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

func (r *SessionRepository) UpdateQRCode(ctx context.Context, id uuid.UUID, qrCode string, expiresAt time.Time) error {
	query := `
		UPDATE "zpSessions" SET
			"qrCode" = $2,
			"qrCodeExpiresAt" = $3,
			"updatedAt" = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id.String(), qrCode, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to update QR code: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

func (r *SessionRepository) UpdateDeviceJID(ctx context.Context, id uuid.UUID, deviceJID string) error {
	query := `
		UPDATE "zpSessions" SET
			"deviceJid" = $2,
			"isConnected" = true,
			"connectedAt" = NOW(),
			"lastSeen" = NOW(),
			"qrCode" = NULL,
			"qrCodeExpiresAt" = NULL,
			"updatedAt" = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id.String(), deviceJID)
	if err != nil {
		return fmt.Errorf("failed to update device JID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

func (r *SessionRepository) ClearQRCode(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE "zpSessions" SET
			"qrCode" = NULL,
			"qrCodeExpiresAt" = NULL,
			"updatedAt" = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to clear QR code: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

func (r *SessionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM "zpSessions" WHERE name = $1)`

	err := r.db.GetContext(ctx, &exists, query, name)
	if err != nil {
		return false, fmt.Errorf("failed to check if session exists by name: %w", err)
	}

	return exists, nil
}

func (r *SessionRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM "zpSessions"`

	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	return count, nil
}

func (r *SessionRepository) toModel(sess *session.Session) (*sessionModel, error) {
	model := &sessionModel{
		ID:          sess.ID().String(),
		Name:        sess.Name(),
		IsConnected: sess.IsConnected(),
		CreatedAt:   sess.CreatedAt(),
		UpdatedAt:   sess.UpdatedAt(),
	}

	if sess.DeviceJID() != nil {
		model.DeviceJID = sql.NullString{String: *sess.DeviceJID(), Valid: true}
	}

	if sess.ConnectionError() != nil {
		model.ConnectionError = sql.NullString{String: *sess.ConnectionError(), Valid: true}
	}

	if sess.QRCode() != nil {
		model.QRCode = sql.NullString{String: *sess.QRCode(), Valid: true}
	}

	if sess.QRCodeExpiresAt() != nil {
		model.QRCodeExpiresAt = sql.NullTime{Time: *sess.QRCodeExpiresAt(), Valid: true}
	}

	if sess.ProxyConfig() != nil {
		proxyJSON, err := json.Marshal(sess.ProxyConfig())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal proxy config: %w", err)
		}
		model.ProxyConfig = sql.NullString{String: string(proxyJSON), Valid: true}
	}

	if sess.ConnectedAt() != nil {
		model.ConnectedAt = sql.NullTime{Time: *sess.ConnectedAt(), Valid: true}
	}

	if sess.LastSeen() != nil {
		model.LastSeen = sql.NullTime{Time: *sess.LastSeen(), Valid: true}
	}

	return model, nil
}

func (r *SessionRepository) fromModel(model *sessionModel) (*session.Session, error) {
	id, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session ID: %w", err)
	}

	var deviceJID *string
	if model.DeviceJID.Valid {
		deviceJID = &model.DeviceJID.String
	}

	var connectionError *string
	if model.ConnectionError.Valid {
		connectionError = &model.ConnectionError.String
	}

	var qrCode *string
	if model.QRCode.Valid {
		qrCode = &model.QRCode.String
	}

	var qrCodeExpiresAt *time.Time
	if model.QRCodeExpiresAt.Valid {
		qrCodeExpiresAt = &model.QRCodeExpiresAt.Time
	}

	var proxyConfig *session.ProxyConfig
	if model.ProxyConfig.Valid {
		var pc session.ProxyConfig
		if err := json.Unmarshal([]byte(model.ProxyConfig.String), &pc); err != nil {
			return nil, fmt.Errorf("failed to unmarshal proxy config: %w", err)
		}
		proxyConfig = &pc
	}

	var connectedAt *time.Time
	if model.ConnectedAt.Valid {
		connectedAt = &model.ConnectedAt.Time
	}

	var lastSeen *time.Time
	if model.LastSeen.Valid {
		lastSeen = &model.LastSeen.Time
	}

	sess := session.RestoreSession(
		id, model.Name, model.IsConnected, deviceJID,
		connectionError, qrCode, qrCodeExpiresAt,
		proxyConfig, model.CreatedAt, model.UpdatedAt,
		connectedAt, lastSeen,
	)

	return sess, nil
}

// CountConnected returns the number of connected sessions
func (r *SessionRepository) CountConnected(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM "zpSessions" WHERE "isConnected" = true`

	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count connected sessions: %w", err)
	}

	return count, nil
}

// FindByID finds a session by ID
func (r *SessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*session.Session, error) {
	var model sessionModel
	query := `SELECT * FROM "zpSessions" WHERE "id" = $1`

	err := r.db.GetContext(ctx, &model, query, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, session.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to find session by ID: %w", err)
	}

	return r.fromModel(&model)
}

// FindAll returns all sessions
func (r *SessionRepository) FindAll(ctx context.Context) ([]*session.Session, error) {
	var sessions []*sessionModel
	query := `SELECT * FROM "zpSessions" ORDER BY "createdAt" DESC`

	err := r.db.SelectContext(ctx, &sessions, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all sessions: %w", err)
	}

	result := make([]*session.Session, len(sessions))
	for i, s := range sessions {
		sess, err := r.toDomainSession(s)
		if err != nil {
			return nil, fmt.Errorf("failed to convert session %s: %w", s.ID, err)
		}
		result[i] = sess
	}

	return result, nil
}
