package webhook

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func ValidateURL(webhookURL string) error {
	if webhookURL == "" {
		return errors.New("webhook URL cannot be empty")
	}

	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must have a scheme (http or https)")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https, got: %s", parsedURL.Scheme)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	if isLocalhost(parsedURL.Host) {
		return fmt.Errorf("localhost URLs are not allowed for webhooks")
	}

	if len(webhookURL) > 2048 {
		return fmt.Errorf("URL too long, maximum length is 2048 characters")
	}

	return nil
}
func isLocalhost(host string) bool {
	if colonIndex := strings.LastIndex(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}

	localhosts := []string{
		"localhost",
		"127.0.0.1",
		"::1",
		"0.0.0.0",
	}

	for _, localhost := range localhosts {
		if host == localhost {
			return true
		}
	}

	return false
}
func ValidateSecret(secret string) error {
	if secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}

	if len(secret) < 8 {
		return fmt.Errorf("secret must be at least 8 characters long")
	}

	if len(secret) > 255 {
		return fmt.Errorf("secret too long, maximum length is 255 characters")
	}

	return nil
}
func ValidateEvents(events []string, validEvents []string) error {
	if len(events) == 0 {
		return nil
	}

	validEventMap := make(map[string]bool)
	for _, event := range validEvents {
		validEventMap[event] = true
	}

	for _, event := range events {
		if event == "" {
			return fmt.Errorf("event name cannot be empty")
		}

		if !validEventMap[event] {
			return fmt.Errorf("invalid event type: %s", event)
		}
	}

	eventSet := make(map[string]bool)
	for _, event := range events {
		if eventSet[event] {
			return fmt.Errorf("duplicate event type: %s", event)
		}

		eventSet[event] = true
	}

	return nil
}
