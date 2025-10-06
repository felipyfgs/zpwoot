package webhook

import (
	"time"

	"github.com/google/uuid"
)

// Webhook representa a configuração de webhook para uma sessão WhatsApp
// REGRA: Zero dependências externas, apenas stdlib e uuid
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

// HasEvent verifica se o webhook está configurado para receber um tipo específico de evento
func (w *Webhook) HasEvent(eventType string) bool {
	if !w.Enabled {
		return false
	}

	// Se Events estiver vazio, aceita todos os eventos
	if len(w.Events) == 0 {
		return true
	}

	for _, e := range w.Events {
		if e == eventType {
			return true
		}
	}
	return false
}

// IsActive verifica se o webhook está ativo
func (w *Webhook) IsActive() bool {
	return w.Enabled
}

// Enable ativa o webhook
func (w *Webhook) Enable() {
	w.Enabled = true
	w.UpdatedAt = time.Now()
}

// Disable desativa o webhook
func (w *Webhook) Disable() {
	w.Enabled = false
	w.UpdatedAt = time.Now()
}

// SetSecret define o secret do webhook
func (w *Webhook) SetSecret(secret string) {
	w.Secret = &secret
	w.UpdatedAt = time.Now()
}

// UpdateURL atualiza a URL do webhook
func (w *Webhook) UpdateURL(url string) {
	w.URL = url
	w.UpdatedAt = time.Now()
}

// UpdateEvents atualiza a lista de eventos do webhook
func (w *Webhook) UpdateEvents(events []string) {
	w.Events = events
	w.UpdatedAt = time.Now()
}

// Update atualiza múltiplos campos do webhook
func (w *Webhook) Update(url string, events []string, secret *string) {
	w.URL = url
	w.Events = events
	if secret != nil {
		w.Secret = secret
	}
	w.UpdatedAt = time.Now()
}

