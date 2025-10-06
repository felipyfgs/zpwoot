package session

import (
	"context"
	"fmt"
	"strings"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/output"
)

type PairUseCase struct {
	whatsappClient output.WhatsAppClient
}

func NewPairUseCase(whatsappClient output.WhatsAppClient) *PairUseCase {
	return &PairUseCase{
		whatsappClient: whatsappClient,
	}
}
func (uc *PairUseCase) Execute(ctx context.Context, sessionID string, phone string) (*dto.PairPhoneResponse, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("sessionID is required")
	}

	if phone == "" {
		return nil, fmt.Errorf("phone number is required")
	}

	cleanPhone := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}

		return -1
	}, phone)

	if len(cleanPhone) < 10 {
		return nil, fmt.Errorf("invalid phone number")
	}

	linkingCode, err := uc.whatsappClient.PairPhone(ctx, sessionID, cleanPhone)
	if err != nil {
		return nil, fmt.Errorf("failed to pair phone: %w", err)
	}

	return &dto.PairPhoneResponse{
		LinkingCode: linkingCode,
	}, nil
}
