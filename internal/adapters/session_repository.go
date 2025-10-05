package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zpwoot/internal/domain/session"
)

// SessionRepository implements the session.Repository interface
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

// Create creates a new session
func (r *SessionRepository) Create(ctx context.Context, sess *session.Session) error {
	// Generate ID if not provided
	if sess.ID == "" {
		sess.ID = uuid.New().String()
	}

	query := `
		INSERT INTO sessions (id, jid, name, status, qr_code, created_at, updated_at, connected_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		sess.ID,
		sess.JID,
		sess.Name,
		string(sess.Status),
		sess.QRCode,
		sess.CreatedAt,
		sess.UpdatedAt,
		sess.ConnectedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetByID retrieves a session by ID
func (r *SessionRepository) GetByID(ctx context.Context, id string) (*session.Session, error) {
	query := `
		SELECT id, jid, name, status, qr_code, created_at, updated_at, connected_at
		FROM sessions
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)
	
	sess := &session.Session{}
	var status string
	var connectedAt sql.NullTime

	err := row.Scan(
		&sess.ID,
		&sess.JID,
		&sess.Name,
		&status,
		&sess.QRCode,
		&sess.CreatedAt,
		&sess.UpdatedAt,
		&connectedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, session.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	sess.Status = session.Status(status)
	if connectedAt.Valid {
		sess.ConnectedAt = &connectedAt.Time
	}

	return sess, nil
}

// GetByJID retrieves a session by JID
func (r *SessionRepository) GetByJID(ctx context.Context, jid string) (*session.Session, error) {
	query := `
		SELECT id, jid, name, status, qr_code, created_at, updated_at, connected_at
		FROM sessions
		WHERE jid = ?
	`

	row := r.db.QueryRowContext(ctx, query, jid)
	
	sess := &session.Session{}
	var status string
	var connectedAt sql.NullTime

	err := row.Scan(
		&sess.ID,
		&sess.JID,
		&sess.Name,
		&status,
		&sess.QRCode,
		&sess.CreatedAt,
		&sess.UpdatedAt,
		&connectedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, session.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by JID: %w", err)
	}

	sess.Status = session.Status(status)
	if connectedAt.Valid {
		sess.ConnectedAt = &connectedAt.Time
	}

	return sess, nil
}

// Update updates an existing session
func (r *SessionRepository) Update(ctx context.Context, sess *session.Session) error {
	sess.UpdatedAt = time.Now()

	query := `
		UPDATE sessions
		SET jid = ?, name = ?, status = ?, qr_code = ?, updated_at = ?, connected_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		sess.JID,
		sess.Name,
		string(sess.Status),
		sess.QRCode,
		sess.UpdatedAt,
		sess.ConnectedAt,
		sess.ID,
	)

	if err != nil {
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

// Delete deletes a session by ID
func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
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

// List retrieves all sessions with pagination
func (r *SessionRepository) List(ctx context.Context, limit, offset int) ([]*session.Session, error) {
	query := `
		SELECT id, jid, name, status, qr_code, created_at, updated_at, connected_at
		FROM sessions
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*session.Session

	for rows.Next() {
		sess := &session.Session{}
		var status string
		var connectedAt sql.NullTime

		err := rows.Scan(
			&sess.ID,
			&sess.JID,
			&sess.Name,
			&status,
			&sess.QRCode,
			&sess.CreatedAt,
			&sess.UpdatedAt,
			&connectedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		sess.Status = session.Status(status)
		if connectedAt.Valid {
			sess.ConnectedAt = &connectedAt.Time
		}

		sessions = append(sessions, sess)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sessions: %w", err)
	}

	return sessions, nil
}

// UpdateStatus updates the session status
func (r *SessionRepository) UpdateStatus(ctx context.Context, id string, status session.Status) error {
	var connectedAt *time.Time
	if status == session.StatusConnected {
		now := time.Now()
		connectedAt = &now
	}

	query := `
		UPDATE sessions
		SET status = ?, updated_at = ?, connected_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, string(status), time.Now(), connectedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
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

// UpdateQRCode updates the session QR code
func (r *SessionRepository) UpdateQRCode(ctx context.Context, id string, qrCode string) error {
	query := `
		UPDATE sessions
		SET qr_code = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, qrCode, time.Now(), id)
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
