package validators

import (
	"fmt"
	"regexp"
)

// Message validation rules
const (
	MessageTextMaxLength     = 4096  // WhatsApp limit
	MessageCaptionMaxLength  = 1024
	PhoneNumberMinLength     = 10
	PhoneNumberMaxLength     = 15
	MediaFileNameMaxLength   = 255
	LocationNameMaxLength    = 100
	LocationAddressMaxLength = 200
	ContactNameMaxLength     = 100
)

var (
	// PhoneNumberRegex validates phone numbers (digits only, with optional + prefix)
	PhoneNumberRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	
	// JIDRegex validates WhatsApp JID format (phone@s.whatsapp.net or phone@g.us)
	JIDRegex = regexp.MustCompile(`^[0-9]+@(s\.whatsapp\.net|g\.us)$`)
)

// ValidatePhoneNumber validates a phone number
func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	if !PhoneNumberRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone number format (must be 10-15 digits, optional + prefix)")
	}

	return nil
}

// ValidateJID validates a WhatsApp JID
func ValidateJID(jid string) error {
	if jid == "" {
		return fmt.Errorf("JID cannot be empty")
	}

	if !JIDRegex.MatchString(jid) {
		return fmt.Errorf("invalid JID format (must be phone@s.whatsapp.net or phone@g.us)")
	}

	return nil
}

// ValidateMessageText validates message text
func ValidateMessageText(text string) error {
	if text == "" {
		return fmt.Errorf("message text cannot be empty")
	}

	if len(text) > MessageTextMaxLength {
		return fmt.Errorf("message text exceeds maximum length of %d characters", MessageTextMaxLength)
	}

	return nil
}

// ValidateCaption validates media caption
func ValidateCaption(caption string) error {
	if caption == "" {
		return nil // Empty caption is allowed
	}

	if len(caption) > MessageCaptionMaxLength {
		return fmt.Errorf("caption exceeds maximum length of %d characters", MessageCaptionMaxLength)
	}

	return nil
}

// ValidateFileName validates file name
func ValidateFileName(fileName string) error {
	if fileName == "" {
		return nil // Empty file name is allowed
	}

	if len(fileName) > MediaFileNameMaxLength {
		return fmt.Errorf("file name exceeds maximum length of %d characters", MediaFileNameMaxLength)
	}

	// Check for invalid characters
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if contains(fileName, char) {
			return fmt.Errorf("file name contains invalid character: %s", char)
		}
	}

	return nil
}

// ValidateLatitude validates latitude
func ValidateLatitude(lat float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	return nil
}

// ValidateLongitude validates longitude
func ValidateLongitude(lon float64) error {
	if lon < -180 || lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}
	return nil
}

// ValidateLocationName validates location name
func ValidateLocationName(name string) error {
	if name == "" {
		return nil // Empty name is allowed
	}

	if len(name) > LocationNameMaxLength {
		return fmt.Errorf("location name exceeds maximum length of %d characters", LocationNameMaxLength)
	}

	return nil
}

// ValidateLocationAddress validates location address
func ValidateLocationAddress(address string) error {
	if address == "" {
		return nil // Empty address is allowed
	}

	if len(address) > LocationAddressMaxLength {
		return fmt.Errorf("location address exceeds maximum length of %d characters", LocationAddressMaxLength)
	}

	return nil
}

// ValidateContactName validates contact name
func ValidateContactName(name string) error {
	if name == "" {
		return fmt.Errorf("contact name cannot be empty")
	}

	if len(name) > ContactNameMaxLength {
		return fmt.Errorf("contact name exceeds maximum length of %d characters", ContactNameMaxLength)
	}

	return nil
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

