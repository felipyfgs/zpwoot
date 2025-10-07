package dto

import "time"

type CreateWebhookRequest struct {
	URL    string   `json:"url" validate:"required,url"`
	Secret *string  `json:"secret,omitempty"`
	Events []string `json:"events,omitempty"`
} // @name CreateWebhookRequest
type WebhookResponse struct {
	ID        string    `json:"id"`
	SessionID string    `json:"sessionId"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
} // @name WebhookResponse
type EventCategoryResponse struct {
	Category string   `json:"category"`
	Events   []string `json:"events"`
} // @name EventCategoryResponse
type ListEventsResponse struct {
	Categories []EventCategoryResponse `json:"categories"`
	AllEvents  []string                `json:"allEvents"`
} // @name ListEventsResponse
type WebhookEventPayload struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}
