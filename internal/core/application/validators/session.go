package validators

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

const (
	SessionNameMinLength = 1
	SessionNameMaxLength = 100
	SessionIDLength      = 36
)

var (
	SessionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)

	SessionIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

func ValidateSessionName(name string) error {
	if name == "" {
		return fmt.Errorf("session name cannot be empty")
	}

	length := utf8.RuneCountInString(name)
	if length < SessionNameMinLength {
		return fmt.Errorf("session name must be at least %d characters", SessionNameMinLength)
	}
	if length > SessionNameMaxLength {
		return fmt.Errorf("session name must not exceed %d characters", SessionNameMaxLength)
	}

	if !SessionNameRegex.MatchString(name) {
		return fmt.Errorf("session name contains invalid characters (only alphanumeric, spaces, hyphens, and underscores allowed)")
	}

	return nil
}

func ValidateSessionID(id string) error {
	if id == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	if len(id) != SessionIDLength {
		return fmt.Errorf("session ID must be %d characters long", SessionIDLength)
	}

	if !SessionIDRegex.MatchString(id) {
		return fmt.Errorf("session ID must be a valid UUID")
	}

	return nil
}

func ValidateWebhookURL(url string) error {
	if url == "" {
		return nil
	}

	if len(url) < 10 {
		return fmt.Errorf("webhook URL is too short")
	}

	if len(url) > 2048 {
		return fmt.Errorf("webhook URL is too long (max 2048 characters)")
	}

	if len(url) < 7 || (url[:7] != "http://" && url[:8] != "https://") {
		return fmt.Errorf("webhook URL must start with http:// or https://")
	}

	return nil
}
