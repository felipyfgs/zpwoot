package webhook

import (
	"crypto/rand"
	"encoding/hex"
)

// generateSecretKey gera um secret aleatório de 32 bytes (64 caracteres hex)
func generateSecretKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

