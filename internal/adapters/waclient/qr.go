package waclient

import (
	"context"
	"time"

	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/ports/output"
)

// Constants for QR code operations
const (
	QRWaitTimeout    = 1 * time.Second
	QRRefreshTimeout = 10 * time.Second
	QRCheckInterval  = 500 * time.Millisecond
)

// GetQRCodeForSession retrieves the QR code for a session
func (wac *WAClient) GetQRCodeForSession(ctx context.Context, sessionID string) (*QREvent, error) {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if client.IsConnected() {
		return nil, &output.WhatsAppError{Code: "ALREADY_CONNECTED", Message: "session is already connected"}
	}

	// Return existing valid QR code if available
	if wac.hasValidQRCode(client) {
		return &QREvent{
			Event:     "qr",
			Code:      client.QRCode,
			ExpiresAt: client.QRExpiresAt,
		}, nil
	}

	// Start connection process if not already in progress
	if client.Status != session.StatusQRCode && client.Status != session.StatusConnecting {
		if err := wac.ConnectSession(ctx, sessionID); err != nil {
			return nil, err
		}

		// Wait briefly for QR code generation
		if err := wac.waitForQRCode(ctx, sessionID); err != nil {
			return nil, err
		}
	}

	// Check if QR code is now available
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

// RefreshQRCode forces a refresh of the QR code for a session
func (wac *WAClient) RefreshQRCode(ctx context.Context, sessionID string) (*QREvent, error) {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if client.IsConnected() {
		return nil, &output.WhatsAppError{Code: "ALREADY_CONNECTED", Message: "session is already connected"}
	}

	// Disconnect if not already disconnected
	if client.Status != session.StatusDisconnected {
		if err = wac.DisconnectSession(ctx, sessionID); err != nil {
			wac.logger.Warn().Err(err).Str("session_id", sessionID).Msg("Failed to disconnect session for QR refresh")
		}
		wac.waitBriefly(ctx)
	}

	// Start new connection
	if err = wac.ConnectSession(ctx, sessionID); err != nil {
		return nil, err
	}

	// Wait for QR code generation with timeout
	return wac.waitForQRCodeWithTimeout(ctx, sessionID)
}

// CleanupExpiredQRCodes removes expired QR codes from all sessions
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

// StartQRCleanupRoutine starts a background routine to clean up expired QR codes
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

// Helper functions

// hasValidQRCode checks if the client has a valid, non-expired QR code
func (wac *WAClient) hasValidQRCode(client *Client) bool {
	return client.QRCode != "" && !client.QRExpiresAt.IsZero() && time.Now().Before(client.QRExpiresAt)
}

// isQRCodeExpired checks if the client's QR code is expired
func (wac *WAClient) isQRCodeExpired(client *Client) bool {
	return client.QRCode != "" && !client.QRExpiresAt.IsZero() && time.Now().After(client.QRExpiresAt)
}

// waitForQRCode waits briefly for QR code generation
func (wac *WAClient) waitForQRCode(ctx context.Context, sessionID string) error {
	timer := time.NewTimer(QRWaitTimeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return &output.WhatsAppError{Code: "CONTEXT_CANCELLED", Message: "request cancelled"}
	}
}

// waitBriefly waits for a short period
func (wac *WAClient) waitBriefly(ctx context.Context) {
	timer := time.NewTimer(QRCheckInterval)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-ctx.Done():
	}
}

// waitForQRCodeWithTimeout waits for QR code generation with a longer timeout
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
