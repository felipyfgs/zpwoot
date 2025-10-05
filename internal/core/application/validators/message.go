package validators

import (
	"fmt"
	"regexp"
)

const (
	MessageTextMaxLength     = 4096
	MessageCaptionMaxLength  = 1024
	PhoneNumberMinLength     = 10
	PhoneNumberMaxLength     = 15
	MediaFileNameMaxLength   = 255
	LocationNameMaxLength    = 100
	LocationAddressMaxLength = 200
	ContactNameMaxLength     = 100
)

var (
	PhoneNumberRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

	JIDRegex = regexp.MustCompile(`^[0-9]+@(s\.whatsapp\.net|g\.us)$`)
)

func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	if !PhoneNumberRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone number format (must be 10-15 digits, optional + prefix)")
	}

	return nil
}

func ValidateJID(jid string) error {
	if jid == "" {
		return fmt.Errorf("JID cannot be empty")
	}

	if !JIDRegex.MatchString(jid) {
		return fmt.Errorf("invalid JID format (must be phone@s.whatsapp.net or phone@g.us)")
	}

	return nil
}

func ValidateMessageText(text string) error {
	if text == "" {
		return fmt.Errorf("message text cannot be empty")
	}

	if len(text) > MessageTextMaxLength {
		return fmt.Errorf("message text exceeds maximum length of %d characters", MessageTextMaxLength)
	}

	return nil
}

func ValidateCaption(caption string) error {
	if caption == "" {
		return nil
	}

	if len(caption) > MessageCaptionMaxLength {
		return fmt.Errorf("caption exceeds maximum length of %d characters", MessageCaptionMaxLength)
	}

	return nil
}

func ValidateFileName(fileName string) error {
	if fileName == "" {
		return nil
	}

	if len(fileName) > MediaFileNameMaxLength {
		return fmt.Errorf("file name exceeds maximum length of %d characters", MediaFileNameMaxLength)
	}

	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if contains(fileName, char) {
			return fmt.Errorf("file name contains invalid character: %s", char)
		}
	}

	return nil
}

func ValidateLatitude(lat float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	return nil
}

func ValidateLongitude(lon float64) error {
	if lon < -180 || lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}
	return nil
}

func ValidateLocationName(name string) error {
	if name == "" {
		return nil
	}

	if len(name) > LocationNameMaxLength {
		return fmt.Errorf("location name exceeds maximum length of %d characters", LocationNameMaxLength)
	}

	return nil
}

func ValidateLocationAddress(address string) error {
	if address == "" {
		return nil
	}

	if len(address) > LocationAddressMaxLength {
		return fmt.Errorf("location address exceeds maximum length of %d characters", LocationAddressMaxLength)
	}

	return nil
}

func ValidateContactName(name string) error {
	if name == "" {
		return fmt.Errorf("contact name cannot be empty")
	}

	if len(name) > ContactNameMaxLength {
		return fmt.Errorf("contact name exceeds maximum length of %d characters", ContactNameMaxLength)
	}

	return nil
}

func contains(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
