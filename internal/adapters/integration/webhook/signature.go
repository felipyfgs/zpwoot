package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateSignature gera uma assinatura HMAC-SHA256 para o payload
// usando o secret fornecido
func GenerateSignature(payload []byte, secret string) string {
	if secret == "" {
		return ""
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	signature := h.Sum(nil)

	return hex.EncodeToString(signature)
}

// ValidateSignature valida se a assinatura fornecida corresponde ao payload
func ValidateSignature(payload []byte, secret string, signature string) bool {
	expectedSignature := GenerateSignature(payload, secret)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
