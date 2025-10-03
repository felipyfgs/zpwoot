package waclient

import (
	"fmt"
	"regexp"
	"strings"

	"go.mau.fi/whatsmeow/types"
	"zpwoot/platform/logger"
)

// Validator provides validation functions for WhatsApp data
type Validator struct {
	logger *logger.Logger
}

// NewValidator creates a new validator instance
func NewValidator(logger *logger.Logger) *Validator {
	return &Validator{
		logger: logger,
	}
}

// parseJID parses and validates a JID string, based on wuzapi implementation
func parseJID(arg string) (types.JID, bool) {
	// Remove leading + if present
	if len(arg) > 0 && arg[0] == '+' {
		arg = arg[1:]
	}

	// If no @ symbol, assume it's a phone number for default user server
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	}

	// Parse as full JID
	recipient, err := types.ParseJID(arg)
	if err != nil {
		return recipient, false
	}

	// Validate that user part is not empty
	if recipient.User == "" {
		return recipient, false
	}

	return recipient, true
}

// ParseJID parses and validates a JID string (public method)
func (v *Validator) ParseJID(jid string) (types.JID, error) {
	parsedJID, valid := parseJID(jid)
	if !valid {
		v.logger.ErrorWithFields("Invalid JID format", map[string]interface{}{
			"jid": jid,
		})
		return types.JID{}, fmt.Errorf("invalid JID format: %s", jid)
	}

	v.logger.DebugWithFields("JID parsed successfully", map[string]interface{}{
		"original_jid": jid,
		"parsed_jid":   parsedJID.String(),
		"user":         parsedJID.User,
		"server":       parsedJID.Server,
	})

	return parsedJID, nil
}

// ValidatePhoneNumber validates a phone number format
func (v *Validator) ValidatePhoneNumber(phone string) error {
	// Remove common prefixes and formatting
	cleaned := strings.TrimSpace(phone)
	cleaned = strings.TrimPrefix(cleaned, "+")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Check if it's all digits
	if !regexp.MustCompile(`^\d+$`).MatchString(cleaned) {
		return fmt.Errorf("phone number must contain only digits: %s", phone)
	}

	// Check minimum length (international format)
	if len(cleaned) < 7 {
		return fmt.Errorf("phone number too short: %s", phone)
	}

	// Check maximum length
	if len(cleaned) > 15 {
		return fmt.Errorf("phone number too long: %s", phone)
	}

	v.logger.DebugWithFields("Phone number validated", map[string]interface{}{
		"original": phone,
		"cleaned":  cleaned,
	})

	return nil
}

// NormalizePhoneNumber normalizes a phone number to WhatsApp format
func (v *Validator) NormalizePhoneNumber(phone string) (string, error) {
	err := v.ValidatePhoneNumber(phone)
	if err != nil {
		return "", err
	}

	// Clean the phone number
	cleaned := strings.TrimSpace(phone)
	cleaned = strings.TrimPrefix(cleaned, "+")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	return cleaned, nil
}

// IsValidJID checks if a JID string is valid
func (v *Validator) IsValidJID(jid string) bool {
	_, valid := parseJID(jid)
	return valid
}

// IsGroupJID checks if a JID represents a group
func (v *Validator) IsGroupJID(jid string) bool {
	parsedJID, valid := parseJID(jid)
	if !valid {
		return false
	}
	return parsedJID.Server == types.GroupServer
}

// IsUserJID checks if a JID represents a user
func (v *Validator) IsUserJID(jid string) bool {
	parsedJID, valid := parseJID(jid)
	if !valid {
		return false
	}
	return parsedJID.Server == types.DefaultUserServer
}

// IsBroadcastJID checks if a JID represents a broadcast list
func (v *Validator) IsBroadcastJID(jid string) bool {
	parsedJID, valid := parseJID(jid)
	if !valid {
		return false
	}
	return parsedJID.Server == types.BroadcastServer
}

// ValidateSessionName validates a session name format
func (v *Validator) ValidateSessionName(name string) error {
	if name == "" {
		return fmt.Errorf("session name cannot be empty")
	}

	if len(name) > 100 {
		return fmt.Errorf("session name too long (max 100 characters): %s", name)
	}

	// Check for valid characters (alphanumeric, underscore, hyphen)
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
		return fmt.Errorf("session name contains invalid characters (only alphanumeric, underscore, and hyphen allowed): %s", name)
	}

	v.logger.DebugWithFields("Session name validated", map[string]interface{}{
		"session_name": name,
	})

	return nil
}

// ExtractPhoneFromJID extracts phone number from a JID
func (v *Validator) ExtractPhoneFromJID(jid string) (string, error) {
	parsedJID, valid := parseJID(jid)
	if !valid {
		return "", fmt.Errorf("invalid JID: %s", jid)
	}

	// For user JIDs, the user part is the phone number
	if parsedJID.Server == types.DefaultUserServer {
		return parsedJID.User, nil
	}

	return "", fmt.Errorf("JID is not a user JID: %s", jid)
}

// FormatJIDForDisplay formats a JID for display purposes
func (v *Validator) FormatJIDForDisplay(jid string) string {
	parsedJID, valid := parseJID(jid)
	if !valid {
		return jid // Return original if invalid
	}

	switch parsedJID.Server {
	case types.DefaultUserServer:
		// Format phone number with +
		return "+" + parsedJID.User
	case types.GroupServer:
		// Return group JID as-is
		return parsedJID.String()
	case types.BroadcastServer:
		// Return broadcast JID as-is
		return parsedJID.String()
	default:
		return parsedJID.String()
	}
}

// ValidateMessageContent validates message content
func (v *Validator) ValidateMessageContent(content string) error {
	if content == "" {
		return fmt.Errorf("message content cannot be empty")
	}

	// Check maximum length (WhatsApp limit is around 65536 characters)
	if len(content) > 65536 {
		return fmt.Errorf("message content too long (max 65536 characters)")
	}

	return nil
}

// SanitizeInput sanitizes user input to prevent injection attacks
func (v *Validator) SanitizeInput(input string) string {
	// Remove null bytes
	sanitized := strings.ReplaceAll(input, "\x00", "")
	
	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)
	
	return sanitized
}

// ValidateQRCode validates QR code format
func (v *Validator) ValidateQRCode(qrCode string) error {
	if qrCode == "" {
		return fmt.Errorf("QR code cannot be empty")
	}

	// Basic validation - QR codes should be reasonable length
	if len(qrCode) < 10 {
		return fmt.Errorf("QR code too short")
	}

	if len(qrCode) > 2048 {
		return fmt.Errorf("QR code too long")
	}

	return nil
}
