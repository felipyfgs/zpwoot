package dto

import "time"

// CreateWebhookRequest representa a requisição para criar/atualizar um webhook
type CreateWebhookRequest struct {
	URL    string   `json:"url" validate:"required,url"`
	Secret *string  `json:"secret,omitempty"`
	Events []string `json:"events,omitempty"` // Se vazio, aceita todos os eventos
} //@name CreateWebhookRequest

// WebhookResponse representa a resposta com os dados do webhook
type WebhookResponse struct {
	ID        string    `json:"id"`
	SessionID string    `json:"sessionId"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
} //@name WebhookResponse

// EventCategoryResponse representa uma categoria de eventos
type EventCategoryResponse struct {
	Category string   `json:"category"`
	Events   []string `json:"events"`
} //@name EventCategoryResponse

// ListEventsResponse representa a resposta com a lista de eventos disponíveis
type ListEventsResponse struct {
	Categories []EventCategoryResponse `json:"categories"`
	AllEvents  []string                `json:"allEvents"`
} //@name ListEventsResponse

// WebhookEventPayload representa o payload enviado para o webhook
type WebhookEventPayload struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}
