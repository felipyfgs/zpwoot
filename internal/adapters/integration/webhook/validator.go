package webhook

import (
	"fmt"
	"net/url"
	"strings"
)

// ValidateURL valida se uma URL é válida para webhook
func ValidateURL(webhookURL string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	// Parse da URL
	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Verificar se tem scheme
	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must have a scheme (http or https)")
	}

	// Verificar se é HTTP ou HTTPS
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https, got: %s", parsedURL.Scheme)
	}

	// Verificar se tem host
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	// Verificar se não é localhost em produção (opcional - pode ser removido se necessário)
	if isLocalhost(parsedURL.Host) {
		// Em desenvolvimento, permitir localhost
		// Em produção, você pode querer bloquear isso
		// Por enquanto, vamos permitir para facilitar testes
	}

	// Verificar comprimento máximo
	if len(webhookURL) > 2048 {
		return fmt.Errorf("URL too long, maximum length is 2048 characters")
	}

	return nil
}

// isLocalhost verifica se o host é localhost
func isLocalhost(host string) bool {
	// Remove porta se existir
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

// ValidateSecret valida se um secret é válido
func ValidateSecret(secret string) error {
	if secret == "" {
		return fmt.Errorf("secret cannot be empty")
	}

	// Verificar comprimento mínimo
	if len(secret) < 8 {
		return fmt.Errorf("secret must be at least 8 characters long")
	}

	// Verificar comprimento máximo
	if len(secret) > 255 {
		return fmt.Errorf("secret too long, maximum length is 255 characters")
	}

	return nil
}

// ValidateEvents valida se a lista de eventos é válida
func ValidateEvents(events []string, validEvents []string) error {
	if len(events) == 0 {
		// Se não especificar eventos, aceita todos (comportamento padrão)
		return nil
	}

	// Criar mapa dos eventos válidos para busca rápida
	validEventMap := make(map[string]bool)
	for _, event := range validEvents {
		validEventMap[event] = true
	}

	// Verificar se todos os eventos são válidos
	for _, event := range events {
		if event == "" {
			return fmt.Errorf("event name cannot be empty")
		}

		if !validEventMap[event] {
			return fmt.Errorf("invalid event type: %s", event)
		}
	}

	// Verificar duplicatas
	eventSet := make(map[string]bool)
	for _, event := range events {
		if eventSet[event] {
			return fmt.Errorf("duplicate event type: %s", event)
		}
		eventSet[event] = true
	}

	return nil
}
