package webhook

import (
	"fmt"
	"net/url"
	"strings"
)

// Service contém a lógica de negócio relacionada a webhooks
// REGRA: Apenas stdlib, sem dependências externas
type Service struct{}

// NewService cria uma nova instância do serviço de webhook
func NewService() *Service {
	return &Service{}
}

// ValidateURL valida se a URL do webhook é válida
func (s *Service) ValidateURL(webhookURL string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	// Deve ser HTTP ou HTTPS
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("webhook URL must use http or https scheme")
	}

	// Deve ter um host
	if parsedURL.Host == "" {
		return fmt.Errorf("webhook URL must have a valid host")
	}

	// Não permitir localhost em produção (opcional, pode ser configurável)
	if strings.Contains(parsedURL.Host, "localhost") || strings.Contains(parsedURL.Host, "127.0.0.1") {
		// Você pode querer permitir isso em desenvolvimento
		// return fmt.Errorf("webhook URL cannot point to localhost")
	}

	return nil
}

// ValidateEvents valida se os eventos fornecidos são válidos
func (s *Service) ValidateEvents(events []string) error {
	if len(events) == 0 {
		// Vazio significa "todos os eventos"
		return nil
	}

	validEvents := s.GetValidEventTypes()
	validEventMap := make(map[string]bool)
	for _, e := range validEvents {
		validEventMap[e] = true
	}

	for _, event := range events {
		if !validEventMap[event] {
			return fmt.Errorf("invalid event type: %s", event)
		}
	}

	return nil
}

// GetValidEventTypes retorna a lista de tipos de eventos válidos
func (s *Service) GetValidEventTypes() []string {
	return []string{
		"Message",
		"MessageRevoked",
		"MessageReaction",
		"Connected",
		"Disconnected",
		"QRCode",
		"PairSuccess",
		"LoggedOut",
		"HistorySync",
		"Receipt",
		"ChatPresence",
		"GroupInfo",
		"JoinedGroup",
		"Picture",
		"IdentityChange",
		"PrivacySettings",
		"OfflineSyncPreview",
		"OfflineSyncCompleted",
		"AppState",
		"KeepAliveTimeout",
		"KeepAliveRestored",
		"Blocklist",
		"MediaRetry",
		"CallOffer",
		"CallAccept",
		"CallPreAccept",
		"CallTransport",
		"CallOfferNotice",
		"CallRelayLatency",
		"CallTerminate",
		"UnknownCallEvent",
		"NewsletterJoin",
		"NewsletterLeave",
		"NewsletterMuteChange",
		"NewsletterLiveUpdate",
		"NewsletterMessageMeta",
	}
}

// GetEventCategories retorna os eventos agrupados por categoria
func (s *Service) GetEventCategories() map[string][]string {
	return map[string][]string{
		"Messages": {
			"Message",
			"MessageRevoked",
			"MessageReaction",
			"Receipt",
		},
		"Connection": {
			"Connected",
			"Disconnected",
			"QRCode",
			"PairSuccess",
			"LoggedOut",
			"KeepAliveTimeout",
			"KeepAliveRestored",
		},
		"Groups": {
			"GroupInfo",
			"JoinedGroup",
		},
		"User": {
			"Picture",
			"IdentityChange",
			"PrivacySettings",
			"Blocklist",
			"ChatPresence",
		},
		"Sync": {
			"HistorySync",
			"OfflineSyncPreview",
			"OfflineSyncCompleted",
			"AppState",
		},
		"Calls": {
			"CallOffer",
			"CallAccept",
			"CallPreAccept",
			"CallTransport",
			"CallOfferNotice",
			"CallRelayLatency",
			"CallTerminate",
			"UnknownCallEvent",
		},
		"Newsletter": {
			"NewsletterJoin",
			"NewsletterLeave",
			"NewsletterMuteChange",
			"NewsletterLiveUpdate",
			"NewsletterMessageMeta",
		},
		"Media": {
			"MediaRetry",
		},
	}
}

// ValidateSecret valida se o secret fornecido é adequado
func (s *Service) ValidateSecret(secret string) error {
	if secret == "" {
		return fmt.Errorf("webhook secret cannot be empty")
	}

	if len(secret) < 16 {
		return fmt.Errorf("webhook secret must be at least 16 characters long")
	}

	return nil
}

