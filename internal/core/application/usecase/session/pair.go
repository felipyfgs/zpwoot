package session

import (
	"context"
	"fmt"
	"strings"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/output"
)

// PairUseCase implementa o pareamento por telefone
type PairUseCase struct {
	whatsappClient output.WhatsAppClient
}

// NewPairUseCase cria uma nova instância do PairUseCase
func NewPairUseCase(whatsappClient output.WhatsAppClient) *PairUseCase {
	return &PairUseCase{
		whatsappClient: whatsappClient,
	}
}

// Execute realiza o pareamento por telefone
func (uc *PairUseCase) Execute(ctx context.Context, sessionID string, phone string) (*dto.PairPhoneResponse, error) {
	// Validar sessionID
	if sessionID == "" {
		return nil, fmt.Errorf("sessionID is required")
	}

	// Validar phone
	if phone == "" {
		return nil, fmt.Errorf("phone number is required")
	}

	// Limpar o número de telefone (remover caracteres não numéricos)
	cleanPhone := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	if len(cleanPhone) < 10 {
		return nil, fmt.Errorf("invalid phone number")
	}

	// Chamar o cliente WhatsApp para fazer o pareamento
	linkingCode, err := uc.whatsappClient.PairPhone(ctx, sessionID, cleanPhone)
	if err != nil {
		return nil, fmt.Errorf("failed to pair phone: %w", err)
	}

	return &dto.PairPhoneResponse{
		LinkingCode: linkingCode,
	}, nil
}
