package waclient

import (
	"context"
	"time"

	"zpwoot/internal/core/domain/session"
)

func (wac *WAClient) GetQRCodeForSession(ctx context.Context, sessionID string) (*QREvent, error) {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if client.Status == session.StatusConnected {
		return nil, &WAError{Code: "ALREADY_CONNECTED", Message: "session is already connected"}
	}

	if client.QRCode != "" && !client.QRExpiresAt.IsZero() && time.Now().Before(client.QRExpiresAt) {
		return &QREvent{
			Event:     "qr",
			Code:      client.QRCode,
			ExpiresAt: client.QRExpiresAt,
		}, nil
	}

	if client.Status != session.StatusQRCode && client.Status != session.StatusConnecting {
		if err := wac.ConnectSession(ctx, sessionID); err != nil {
			return nil, err
		}

		timer := time.NewTimer(1 * time.Second)
		defer timer.Stop()

		select {
		case <-timer.C:

		case <-ctx.Done():
			return nil, &WAError{Code: "CONTEXT_CANCELLED", Message: "request cancelled"}
		}

		if client, err = wac.GetSession(ctx, sessionID); err != nil {
			return nil, err
		}

		if client.QRCode != "" {
			return &QREvent{
				Event:     "qr",
				Code:      client.QRCode,
				ExpiresAt: client.QRExpiresAt,
			}, nil
		}
	}

	return nil, &WAError{Code: "QR_NOT_AVAILABLE", Message: "QR code not available, try connecting the session first"}
}

func (wac *WAClient) RefreshQRCode(ctx context.Context, sessionID string) (*QREvent, error) {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if client.Status == session.StatusConnected {
		return nil, &WAError{Code: "ALREADY_CONNECTED", Message: "session is already connected"}
	}

	if client.Status != session.StatusDisconnected {
		if err = wac.DisconnectSession(ctx, sessionID); err != nil {
			wac.logger.Warn().Err(err).Str("session_id", sessionID).Msg("Failed to disconnect session for QR refresh")
		}

		timer := time.NewTimer(500 * time.Millisecond)
		defer timer.Stop()

		select {
		case <-timer.C:

		case <-ctx.Done():
			return nil, &WAError{Code: "CONTEXT_CANCELLED", Message: "request cancelled"}
		}
	}

	if err = wac.ConnectSession(ctx, sessionID); err != nil {
		return nil, err
	}

	maxWait := 10 * time.Second
	checkInterval := 500 * time.Millisecond
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	timeout := time.After(maxWait)

	for {
		select {
		case <-ticker.C:
			if client, err = wac.GetSession(ctx, sessionID); err != nil {
				continue
			}

			if client.QRCode != "" && client.Status == session.StatusQRCode {
				return &QREvent{
					Event:     "qr",
					Code:      client.QRCode,
					ExpiresAt: client.QRExpiresAt,
				}, nil
			}

		case <-timeout:
			return nil, &WAError{Code: "QR_GENERATION_TIMEOUT", Message: "timeout waiting for QR code generation"}

		case <-ctx.Done():
			return nil, &WAError{Code: "CONTEXT_CANCELLED", Message: "request cancelled"}
		}
	}
}

func (wac *WAClient) CleanupExpiredQRCodes(ctx context.Context) error {
	sessions, err := wac.ListSessions(ctx)
	if err != nil {
		return err
	}

	cleanedCount := 0
	for _, client := range sessions {
		if client.QRCode != "" && !client.QRExpiresAt.IsZero() && time.Now().After(client.QRExpiresAt) {
			wac.clearQRCode(client)
			wac.updateSessionStatus(ctx, client)
			cleanedCount++
			wac.logger.Debug().Str("session_name", client.Name).Msg("Cleaned expired QR code for session")
		}
	}

	if cleanedCount > 0 {
		wac.logger.Info().Int("count", cleanedCount).Msg("Cleaned expired QR codes")
	}

	return nil
}

func (wac *WAClient) StartQRCleanupRoutine(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	wac.logger.Info().Msg("Started QR code cleanup routine")

	for {
		select {
		case <-ctx.Done():
			wac.logger.Info().Msg("Stopping QR code cleanup routine")
			return
		case <-ticker.C:
			if err := wac.CleanupExpiredQRCodes(ctx); err != nil {
				wac.logger.Error().Err(err).Msg("QR cleanup error")
			}
		}
	}
}
