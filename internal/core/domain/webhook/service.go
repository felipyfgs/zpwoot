package webhook

import (
	"fmt"
	"net/url"
	"strings"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
func (s *Service) ValidateURL(webhookURL string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("webhook URL must use http or https scheme")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("webhook URL must have a valid host")
	}

	if strings.Contains(parsedURL.Host, "localhost") || strings.Contains(parsedURL.Host, "127.0.0.1") {
	}

	return nil
}
func (s *Service) ValidateEvents(events []string) error {
	if len(events) == 0 {
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
func (s *Service) ValidateSecret(secret string) error {
	if secret == "" {
		return fmt.Errorf("webhook secret cannot be empty")
	}

	if len(secret) < 16 {
		return fmt.Errorf("webhook secret must be at least 16 characters long")
	}

	return nil
}
