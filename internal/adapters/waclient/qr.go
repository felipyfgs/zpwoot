package waclient

import (
	"context"
	"time"

	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/ports/output"
)

const (
	QRWaitTimeout    = 1 * time.Second
	QRRefreshTimeout = 10 * time.Second
	QRCheckInterval  = 500 * time.Millisecond
)

func (wac *WAClient) GetQRCodeForSession(ctx context.Context, sessionID string) (*QREvent, error) {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if client.IsConnected() {
		return nil, &output.WhatsAppError{Code: "ALREADY_CONNECTED", Message: "session is already connected"}
	}

	if wac.hasValidQRCode(client) {
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

		if err := wac.waitForQRCode(ctx); err != nil {
			return nil, err
		}
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

	return nil, &output.WhatsAppError{Code: "QR_NOT_AVAILABLE", Message: "QR code not available, try connecting the session first"}
}

func (wac *WAClient) RefreshQRCode(ctx context.Context, sessionID string) (*QREvent, error) {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if client.IsConnected() {
		return nil, &output.WhatsAppError{Code: "ALREADY_CONNECTED", Message: "session is already connected"}
	}

	if client.Status != session.StatusDisconnected {
		if err = wac.DisconnectSession(ctx, sessionID); err != nil {
			wac.logger.Warn().Err(err).Str("session_id", sessionID).Msg("Failed to disconnect session for QR refresh")
		}

		wac.waitBriefly(ctx)
	}

	if err = wac.ConnectSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return wac.waitForQRCodeWithTimeout(ctx, sessionID)
}

func (wac *WAClient) CleanupExpiredQRCodes(ctx context.Context) error {
	sessions, err := wac.ListSessions(ctx)
	if err != nil {
		return err
	}

	cleanedCount := 0

	for _, client := range sessions {
		if wac.isQRCodeExpired(client) {
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

func (wac *WAClient) hasValidQRCode(client *Client) bool {
	return client.QRCode != "" && !client.QRExpiresAt.IsZero() && time.Now().Before(client.QRExpiresAt)
}

func (wac *WAClient) isQRCodeExpired(client *Client) bool {
	return client.QRCode != "" && !client.QRExpiresAt.IsZero() && time.Now().After(client.QRExpiresAt)
}

func (wac *WAClient) waitForQRCode(ctx context.Context) error {
	timer := time.NewTimer(QRWaitTimeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return &output.WhatsAppError{Code: "CONTEXT_CANCELLED", Message: "request cancelled"}
	}
}

func (wac *WAClient) waitBriefly(ctx context.Context) {
	timer := time.NewTimer(QRCheckInterval)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-ctx.Done():
	}
}

func (wac *WAClient) waitForQRCodeWithTimeout(ctx context.Context, sessionID string) (*QREvent, error) {
	ticker := time.NewTicker(QRCheckInterval)
	defer ticker.Stop()

	timeout := time.After(QRRefreshTimeout)

	for {
		select {
		case <-ticker.C:
			client, err := wac.GetSession(ctx, sessionID)
			if err != nil {
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
			return nil, &output.WhatsAppError{Code: "QR_GENERATION_TIMEOUT", Message: "timeout waiting for QR code generation"}

		case <-ctx.Done():
			return nil, &output.WhatsAppError{Code: "CONTEXT_CANCELLED", Message: "request cancelled"}
		}
	}
}
