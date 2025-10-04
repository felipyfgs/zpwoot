package waclient

import (
	"fmt"
	"regexp"
	"strings"

	"go.mau.fi/whatsmeow/types"
	"zpwoot/platform/logger"
)

type Validator struct {
	logger *logger.Logger
}

func NewValidator(logger *logger.Logger) *Validator {
	return &Validator{
		logger: logger,
	}
}

func parseJID(arg string) (types.JID, bool) {

	if len(arg) > 0 && arg[0] == '+' {
		arg = arg[1:]
	}

	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	}

	recipient, err := types.ParseJID(arg)
	if err != nil {
		return recipient, false
	}

	if recipient.User == "" {
		return recipient, false
	}

	return recipient, true
}

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

func (v *Validator) ValidatePhoneNumber(phone string) error {

	cleaned := strings.TrimSpace(phone)
	cleaned = strings.TrimPrefix(cleaned, "+")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	if !regexp.MustCompile(`^\d+$`).MatchString(cleaned) {
		return fmt.Errorf("phone number must contain only digits: %s", phone)
	}

	if len(cleaned) < 7 {
		return fmt.Errorf("phone number too short: %s", phone)
	}

	if len(cleaned) > 15 {
		return fmt.Errorf("phone number too long: %s", phone)
	}

	v.logger.DebugWithFields("Phone number validated", map[string]interface{}{
		"original": phone,
		"cleaned":  cleaned,
	})

	return nil
}

func (v *Validator) NormalizePhoneNumber(phone string) (string, error) {
	err := v.ValidatePhoneNumber(phone)
	if err != nil {
		return "", err
	}

	cleaned := strings.TrimSpace(phone)
	cleaned = strings.TrimPrefix(cleaned, "+")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	return cleaned, nil
}

func (v *Validator) IsValidJID(jid string) bool {
	_, valid := parseJID(jid)
	return valid
}

func (v *Validator) IsGroupJID(jid string) bool {
	parsedJID, valid := parseJID(jid)
	if !valid {
		return false
	}
	return parsedJID.Server == types.GroupServer
}

func (v *Validator) IsUserJID(jid string) bool {
	parsedJID, valid := parseJID(jid)
	if !valid {
		return false
	}
	return parsedJID.Server == types.DefaultUserServer
}

func (v *Validator) IsBroadcastJID(jid string) bool {
	parsedJID, valid := parseJID(jid)
	if !valid {
		return false
	}
	return parsedJID.Server == types.BroadcastServer
}

func (v *Validator) ValidatesessionID(name string) error {
	if name == "" {
		return fmt.Errorf("session name cannot be empty")
	}

	if len(name) > 100 {
		return fmt.Errorf("session name too long (max 100 characters): %s", name)
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
		return fmt.Errorf("session name contains invalid characters (only alphanumeric, underscore, and hyphen allowed): %s", name)
	}

	v.logger.DebugWithFields("Session name validated", map[string]interface{}{
		"session_name": name,
	})

	return nil
}

func (v *Validator) ExtractPhoneFromJID(jid string) (string, error) {
	parsedJID, valid := parseJID(jid)
	if !valid {
		return "", fmt.Errorf("invalid JID: %s", jid)
	}

	if parsedJID.Server == types.DefaultUserServer {
		return parsedJID.User, nil
	}

	return "", fmt.Errorf("JID is not a user JID: %s", jid)
}

func (v *Validator) FormatJIDForDisplay(jid string) string {
	parsedJID, valid := parseJID(jid)
	if !valid {
		return jid
	}

	switch parsedJID.Server {
	case types.DefaultUserServer:

		return "+" + parsedJID.User
	case types.GroupServer:

		return parsedJID.String()
	case types.BroadcastServer:

		return parsedJID.String()
	default:
		return parsedJID.String()
	}
}

func (v *Validator) ValidateMessageContent(content string) error {
	if content == "" {
		return fmt.Errorf("message content cannot be empty")
	}

	if len(content) > 65536 {
		return fmt.Errorf("message content too long (max 65536 characters)")
	}

	return nil
}

func (v *Validator) SanitizeInput(input string) string {

	sanitized := strings.ReplaceAll(input, "\x00", "")

	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

func (v *Validator) ValidateQRCode(qrCode string) error {
	if qrCode == "" {
		return fmt.Errorf("QR code cannot be empty")
	}

	if len(qrCode) < 10 {
		return fmt.Errorf("QR code too short")
	}

	if len(qrCode) > 2048 {
		return fmt.Errorf("QR code too long")
	}

	return nil
}
