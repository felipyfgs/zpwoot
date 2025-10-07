package waclient

import (
	"context"
	"fmt"
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

	// Check if session is already authenticated (logged in)
	// We should only return "already connected" if the session is fully authenticated,
	// not just connected at WebSocket level. QR codes are needed for authentication.
	if client.IsConnected() && client.IsLoggedIn() {
		return nil, &output.WhatsAppError{Code: "ALREADY_CONNECTED", Message: "session is already connected"}
	}

	// If we already have a valid QR code, return it immediately
	if wac.hasValidQRCode(client) {
		return &QREvent{
			Event:     "qr",
			Code:      client.QRCode,
			ExpiresAt: client.QRExpiresAt,
		}, nil
	}

	// Use the improved waiting method with timeout
	return wac.waitForQRCodeWithTimeout(ctx, sessionID)
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
	// Create a context with timeout for this operation
	timeoutCtx, cancel := context.WithTimeout(ctx, QRRefreshTimeout)
	defer cancel()

	// Check every 50ms for QR code availability
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	wac.logger.Debug().Str("session_id", sessionID).Msg("Waiting for QR code generation...")

	for {
		select {
		case <-ticker.C:
			client, err := wac.GetSession(timeoutCtx, sessionID)
			if err != nil {
				continue
			}

			// Check if QR code is available
			if client.QRCode != "" && client.Status == session.StatusQRCode {
				wac.logger.Debug().Str("session_id", sessionID).Msg("QR code found and ready")
				return &QREvent{
					Event:     "qr",
					Code:      client.QRCode,
					ExpiresAt: client.QRExpiresAt,
				}, nil
			}

			// Check if session became authenticated while waiting
			if client.IsLoggedIn() {
				return nil, &output.WhatsAppError{Code: "ALREADY_CONNECTED", Message: "session became authenticated while waiting for QR code"}
			}

		case <-timeoutCtx.Done():
			if timeoutCtx.Err() == context.DeadlineExceeded {
				wac.logger.Warn().Str("session_id", sessionID).Msg("Timeout waiting for QR code generation")
				return nil, &output.WhatsAppError{Code: "QR_GENERATION_TIMEOUT", Message: "timeout waiting for QR code generation"}
			}
			return nil, &output.WhatsAppError{Code: "CONTEXT_CANCELLED", Message: "request cancelled"}
		}
	}
}

// ConnectAndGetQRCode connects the session and waits for QR code using WhatsApp Meow's QR channel
func (wac *WAClient) ConnectAndGetQRCode(ctx context.Context, sessionID string) (*output.QRCodeInfo, error) {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if session is already authenticated
	if client.IsConnected() && client.IsLoggedIn() {
		return nil, &output.WhatsAppError{Code: "ALREADY_CONNECTED", Message: "session is already connected"}
	}

	// Get QR channel BEFORE connecting (WhatsApp Meow requirement)
	qrChan, err := client.WAClient.GetQRChannel(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get QR channel: %w", err)
	}

	// Now connect the session
	if !client.IsConnected() {
		if err := client.WAClient.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect session: %w", err)
		}
	}

	// Wait for QR code from the channel
	select {
	case qrItem := <-qrChan:
		if qrItem.Event == "code" && qrItem.Code != "" {
			// Update client with QR code
			client.QRCode = qrItem.Code
			client.QRExpiresAt = time.Now().Add(qrItem.Timeout)
			client.Status = session.StatusQRCode

			return &output.QRCodeInfo{
				Code:      qrItem.Code,
				ExpiresAt: client.QRExpiresAt,
			}, nil
		} else if qrItem.Event == "error" {
			return nil, fmt.Errorf("QR generation error: %v", qrItem.Error)
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while waiting for QR code")
	case <-time.After(5 * time.Second):
		return nil, &output.WhatsAppError{Code: "QR_GENERATION_TIMEOUT", Message: "timeout waiting for QR code generation"}
	}

	return nil, &output.WhatsAppError{Code: "QR_NOT_AVAILABLE", Message: "QR code not available"}
}

// GetQRCode returns the current QR code for the session (if available)
func (wac *WAClient) GetQRCode(ctx context.Context, sessionID string) (*output.QRCodeInfo, error) {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if session is already authenticated
	if client.IsConnected() && client.IsLoggedIn() {
		return nil, &output.WhatsAppError{Code: "ALREADY_CONNECTED", Message: "session is already connected"}
	}

	// Return QR code if available
	if client.QRCode != "" && client.Status == session.StatusQRCode {
		return &output.QRCodeInfo{
			Code:      client.QRCode,
			ExpiresAt: client.QRExpiresAt,
		}, nil
	}

	return nil, &output.WhatsAppError{Code: "QR_NOT_AVAILABLE", Message: "QR code not available, try connecting the session first"}
}
