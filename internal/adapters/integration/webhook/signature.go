package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateSignature(payload []byte, secret string) string {
	if secret == "" {
		return ""
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	signature := h.Sum(nil)

	return hex.EncodeToString(signature)
}
func ValidateSignature(payload []byte, secret string, signature string) bool {
	expectedSignature := GenerateSignature(payload, secret)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
