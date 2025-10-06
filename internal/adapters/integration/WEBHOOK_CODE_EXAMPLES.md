# Exemplos de Código - Webhook Integration (Clean Architecture)

Este documento contém exemplos de código para cada camada da arquitetura, seguindo rigorosamente as regras de dependência.

## 1. Domain Layer (internal/core/domain/webhook/)

### entity.go
```go
package webhook

import (
	"time"
	"github.com/google/uuid"
)

// Webhook representa a configuração de webhook para uma sessão
// REGRA: Zero dependências externas, apenas stdlib
type Webhook struct {
	ID        string
	SessionID string
	URL       string
	Secret    *string
	Events    []string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewWebhook cria uma nova instância de Webhook
func NewWebhook(sessionID, url string, events []string) *Webhook {
	now := time.Now()
	return &Webhook{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		URL:       url,
		Events:    events,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetSecret define o secret do webhook (regra de negócio)
func (w *Webhook) SetSecret(secret string) {
	w.Secret = &secret
	w.UpdatedAt = time.Now()
}

// Enable habilita o webhook
func (w *Webhook) Enable() {
	w.Enabled = true
	w.UpdatedAt = time.Now()
}

// Disable desabilita o webhook
func (w *Webhook) Disable() {
	w.Enabled = false
	w.UpdatedAt = time.Now()
}

// HasEvent verifica se o webhook está inscrito em um evento (regra de negócio)
func (w *Webhook) HasEvent(eventType string) bool {
	if len(w.Events) == 0 {
		return true // Sem filtro = todos os eventos
	}
	
	for _, event := range w.Events {
		if event == "All" || event == eventType {
			return true
		}
	}
	return false
}

// IsActive verifica se o webhook está ativo (regra de negócio)
func (w *Webhook) IsActive() bool {
	return w.Enabled && w.URL != ""
}
```

### repository.go
```go
package webhook

import "context"

// Repository define o contrato para persistência de webhooks
// REGRA: Interface no domain, implementação nos adapters
type Repository interface {
	Create(ctx context.Context, webhook *Webhook) error
	GetByID(ctx context.Context, id string) (*Webhook, error)
	GetBySessionID(ctx context.Context, sessionID string) (*Webhook, error)
	Update(ctx context.Context, webhook *Webhook) error
	Delete(ctx context.Context, id string) error
	DeleteBySessionID(ctx context.Context, sessionID string) error
}
```

### service.go
```go
package webhook

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Service contém lógica de negócio de webhooks
// REGRA: Apenas lógica de domínio, sem dependências externas
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ValidateURL valida se uma URL é válida (regra de negócio)
func (s *Service) ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	return nil
}

// ValidateEvents valida uma lista de eventos (regra de negócio)
func (s *Service) ValidateEvents(events []string) error {
	if len(events) == 0 {
		return nil // Vazio = todos os eventos
	}

	validEvents := GetAllEventTypes()
	for _, event := range events {
		if !contains(validEvents, event) {
			return fmt.Errorf("invalid event type: %s", event)
		}
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
```

---

## 2. Application DTOs (internal/core/application/dto/)

### webhook.go
```go
package dto

import "time"

// CreateWebhookRequest representa o request para criar/atualizar webhook
type CreateWebhookRequest struct {
	URL    string   `json:"url" validate:"required,url"`
	Secret *string  `json:"secret,omitempty"`
	Events []string `json:"events" validate:"dive,required"`
}

// WebhookResponse representa a resposta de webhook
type WebhookResponse struct {
	ID        string    `json:"id"`
	SessionID string    `json:"sessionId"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// EventInfo representa informações sobre um evento
type EventInfo struct {
	Type        string `json:"type"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// EventCategoryResponse representa uma categoria de eventos
type EventCategoryResponse struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Events      []EventInfo `json:"events"`
}

// ListEventsResponse representa a resposta de listagem de eventos
type ListEventsResponse struct {
	Categories []EventCategoryResponse `json:"categories"`
}
```

---

## 3. Ports (internal/core/ports/)

### input/webhook.go
```go
package input

import (
	"context"
	"zpwoot/internal/core/application/dto"
)

// WebhookUseCases define os casos de uso de webhook
// REGRA: Interface para os handlers HTTP usarem
type WebhookUseCases interface {
	CreateOrUpdate(ctx context.Context, sessionID string, req *dto.CreateWebhookRequest) (*dto.WebhookResponse, error)
	Get(ctx context.Context, sessionID string) (*dto.WebhookResponse, error)
	Delete(ctx context.Context, sessionID string) error
	ListEvents(ctx context.Context) (*dto.ListEventsResponse, error)
}
```

### output/webhook_sender.go
```go
package output

import (
	"context"
	"time"
)

// WebhookEvent representa um evento a ser enviado via webhook
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// WebhookSender define a interface para envio de webhooks
// REGRA: Interface no port, implementação no adapter
type WebhookSender interface {
	SendWebhook(ctx context.Context, webhookURL, secret string, event *WebhookEvent) error
}
```

---

## 4. Use Cases (internal/core/application/usecase/webhook/)

### create_or_update.go
```go
package webhook

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/webhook"
)

type CreateOrUpdateUseCase struct {
	webhookService *webhook.Service
	webhookRepo    webhook.Repository
}

func NewCreateOrUpdateUseCase(service *webhook.Service, repo webhook.Repository) *CreateOrUpdateUseCase {
	return &CreateOrUpdateUseCase{
		webhookService: service,
		webhookRepo:    repo,
	}
}

func (uc *CreateOrUpdateUseCase) Execute(ctx context.Context, sessionID string, req *dto.CreateWebhookRequest) (*dto.WebhookResponse, error) {
	// Validar URL (regra de negócio)
	if err := uc.webhookService.ValidateURL(req.URL); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Validar eventos (regra de negócio)
	if err := uc.webhookService.ValidateEvents(req.Events); err != nil {
		return nil, fmt.Errorf("invalid events: %w", err)
	}

	// Buscar webhook existente
	existingWebhook, err := uc.webhookRepo.GetBySessionID(ctx, sessionID)
	
	if err != nil && err != webhook.ErrWebhookNotFound {
		return nil, fmt.Errorf("failed to check existing webhook: %w", err)
	}

	// Upsert
	var wh *webhook.Webhook
	if existingWebhook != nil {
		// Update
		existingWebhook.URL = req.URL
		existingWebhook.Events = req.Events
		if req.Secret != nil {
			existingWebhook.SetSecret(*req.Secret)
		}
		if err := uc.webhookRepo.Update(ctx, existingWebhook); err != nil {
			return nil, fmt.Errorf("failed to update webhook: %w", err)
		}
		wh = existingWebhook
	} else {
		// Create
		wh = webhook.NewWebhook(sessionID, req.URL, req.Events)
		
		secret := req.Secret
		if secret == nil {
			generated, err := generateSecret()
			if err != nil {
				return nil, fmt.Errorf("failed to generate secret: %w", err)
			}
			secret = &generated
		}
		wh.SetSecret(*secret)
		
		if err := uc.webhookRepo.Create(ctx, wh); err != nil {
			return nil, fmt.Errorf("failed to create webhook: %w", err)
		}
	}

	return toDTO(wh), nil
}

func generateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func toDTO(wh *webhook.Webhook) *dto.WebhookResponse {
	return &dto.WebhookResponse{
		ID:        wh.ID,
		SessionID: wh.SessionID,
		URL:       wh.URL,
		Events:    wh.Events,
		Enabled:   wh.Enabled,
		CreatedAt: wh.CreatedAt,
		UpdatedAt: wh.UpdatedAt,
	}
}
```

---

## 5. Adapters

### Database Repository (internal/adapters/database/repository/webhook.go)
```go
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"zpwoot/internal/core/domain/webhook"
)

// PostgresWebhookRepository implementa webhook.Repository
type PostgresWebhookRepository struct {
	db *sqlx.DB
}

func NewPostgresWebhookRepository(db *sqlx.DB) *PostgresWebhookRepository {
	return &PostgresWebhookRepository{db: db}
}

func (r *PostgresWebhookRepository) Create(ctx context.Context, wh *webhook.Webhook) error {
	eventsJSON, _ := json.Marshal(wh.Events)
	
	query := `
		INSERT INTO "zpWebhooks" ("id", "sessionId", "url", "secret", "events", "enabled", "createdAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.ExecContext(ctx, query, wh.ID, wh.SessionID, wh.URL, wh.Secret, eventsJSON, wh.Enabled, wh.CreatedAt, wh.UpdatedAt)
	return err
}

// ... outros métodos
```

### HTTP Handler (internal/adapters/http/handlers/webhook.go)
```go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"
)

type WebhookHandler struct {
	webhookUseCases input.WebhookUseCases
}

func NewWebhookHandler(useCases input.WebhookUseCases) *WebhookHandler {
	return &WebhookHandler{webhookUseCases: useCases}
}

func (h *WebhookHandler) SetWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	
	var req dto.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	
	response, err := h.webhookUseCases.CreateOrUpdate(r.Context(), sessionID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

---

## Resumo das Regras

1. **Domain** (`internal/core/domain/webhook/`)
   - ✅ Apenas stdlib
   - ❌ Sem dependências externas

2. **Application** (`internal/core/application/`)
   - ✅ Pode usar domain + ports
   - ❌ Não pode usar adapters

3. **Ports** (`internal/core/ports/`)
   - ✅ Apenas interfaces
   - ❌ Sem implementações

4. **Adapters** (`internal/adapters/`)
   - ✅ Implementa interfaces dos ports
   - ✅ Pode usar tudo (domain, application, ports)

