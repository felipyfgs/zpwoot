package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"zpwoot/internal/core/domain/webhook"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)


type WebhookRepository struct {
	db *sqlx.DB
}


func NewWebhookRepository(db *sqlx.DB) *WebhookRepository {
	return &WebhookRepository{
		db: db,
	}
}


func (r *WebhookRepository) Create(ctx context.Context, wh *webhook.Webhook) error {
	eventsJSON, err := json.Marshal(wh.Events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	query := `
		INSERT INTO "zpWebhooks" (
			"id", "sessionId", "url", "secret", "events", 
			"enabled", "createdAt", "updatedAt"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err = r.db.ExecContext(ctx, query,
		wh.ID,
		wh.SessionID,
		wh.URL,
		wh.Secret,
		eventsJSON,
		wh.Enabled,
		wh.CreatedAt,
		wh.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return fmt.Errorf("webhook already exists for session")
			}
		}
		return fmt.Errorf("failed to create webhook: %w", err)
	}

	return nil
}


func (r *WebhookRepository) GetByID(ctx context.Context, id string) (*webhook.Webhook, error) {
	query := `
		SELECT "id", "sessionId", "url", "secret", "events", 
		       "enabled", "createdAt", "updatedAt"
		FROM "zpWebhooks"
		WHERE "id" = $1
	`

	var wh webhookDB
	err := r.db.GetContext(ctx, &wh, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("webhook not found")
		}
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}

	return wh.toDomain()
}


func (r *WebhookRepository) GetBySessionID(ctx context.Context, sessionID string) (*webhook.Webhook, error) {
	query := `
		SELECT "id", "sessionId", "url", "secret", "events", 
		       "enabled", "createdAt", "updatedAt"
		FROM "zpWebhooks"
		WHERE "sessionId" = $1
	`

	var wh webhookDB
	err := r.db.GetContext(ctx, &wh, query, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("webhook not found")
		}
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}

	return wh.toDomain()
}


func (r *WebhookRepository) Update(ctx context.Context, wh *webhook.Webhook) error {
	eventsJSON, err := json.Marshal(wh.Events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	query := `
		UPDATE "zpWebhooks" SET
			"url" = $2,
			"secret" = $3,
			"events" = $4,
			"enabled" = $5,
			"updatedAt" = $6
		WHERE "id" = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		wh.ID,
		wh.URL,
		wh.Secret,
		eventsJSON,
		wh.Enabled,
		wh.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("webhook not found")
	}

	return nil
}


func (r *WebhookRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "zpWebhooks" WHERE "id" = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("webhook not found")
	}

	return nil
}


func (r *WebhookRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	query := `DELETE FROM "zpWebhooks" WHERE "sessionId" = $1`

	result, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("webhook not found")
	}

	return nil
}


func (r *WebhookRepository) List(ctx context.Context, limit, offset int) ([]*webhook.Webhook, error) {
	query := `
		SELECT "id", "sessionId", "url", "secret", "events", 
		       "enabled", "createdAt", "updatedAt"
		FROM "zpWebhooks"
		ORDER BY "createdAt" DESC
		LIMIT $1 OFFSET $2
	`

	var webhooksDB []webhookDB
	err := r.db.SelectContext(ctx, &webhooksDB, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}

	webhooks := make([]*webhook.Webhook, 0, len(webhooksDB))
	for _, whDB := range webhooksDB {
		wh, err := whDB.toDomain()
		if err != nil {
			return nil, err
		}
		webhooks = append(webhooks, wh)
	}

	return webhooks, nil
}


type webhookDB struct {
	ID        string       `db:"id"`
	SessionID string       `db:"sessionId"`
	URL       string       `db:"url"`
	Secret    *string      `db:"secret"`
	Events    []byte       `db:"events"`
	Enabled   bool         `db:"enabled"`
	CreatedAt sql.NullTime `db:"createdAt"`
	UpdatedAt sql.NullTime `db:"updatedAt"`
}


func (wh *webhookDB) toDomain() (*webhook.Webhook, error) {
	var events []string
	if len(wh.Events) > 0 {
		if err := json.Unmarshal(wh.Events, &events); err != nil {
			return nil, fmt.Errorf("failed to unmarshal events: %w", err)
		}
	}

	return &webhook.Webhook{
		ID:        wh.ID,
		SessionID: wh.SessionID,
		URL:       wh.URL,
		Secret:    wh.Secret,
		Events:    events,
		Enabled:   wh.Enabled,
		CreatedAt: wh.CreatedAt.Time,
		UpdatedAt: wh.UpdatedAt.Time,
	}, nil
}
