package webhook

import (
	"time"

	"github.com/google/uuid"
)



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


func (w *Webhook) HasEvent(eventType string) bool {
	if !w.Enabled {
		return false
	}


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


func (w *Webhook) IsActive() bool {
	return w.Enabled
}


func (w *Webhook) Enable() {
	w.Enabled = true
	w.UpdatedAt = time.Now()
}


func (w *Webhook) Disable() {
	w.Enabled = false
	w.UpdatedAt = time.Now()
}


func (w *Webhook) SetSecret(secret string) {
	w.Secret = &secret
	w.UpdatedAt = time.Now()
}


func (w *Webhook) UpdateURL(url string) {
	w.URL = url
	w.UpdatedAt = time.Now()
}


func (w *Webhook) UpdateEvents(events []string) {
	w.Events = events
	w.UpdatedAt = time.Now()
}


func (w *Webhook) Update(url string, events []string, secret *string) {
	w.URL = url
	w.Events = events
	if secret != nil {
		w.Secret = secret
	}
	w.UpdatedAt = time.Now()
}
