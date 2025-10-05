package waclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"zpwoot/internal/adapters/logger"

	"github.com/skip2/go-qrcode"
)

// QRManager handles QR code generation and management
type QRManager struct {
	logger *logger.Logger
}

// NewQRManager creates a new QR code manager
func NewQRManager(logger *logger.Logger) *QRManager {
	return &QRManager{
		logger: logger,
	}
}

// GenerateQRCodeBase64 generates a base64 encoded QR code image
func (qm *QRManager) GenerateQRCodeBase64(qrText string) (string, error) {
	// Generate QR code as PNG bytes
	qrBytes, err := qrcode.Encode(qrText, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Encode to base64
	base64QR := base64.StdEncoding.EncodeToString(qrBytes)
	return base64QR, nil
}

// GenerateQREvent creates a QR event with base64 encoded image
func (qm *QRManager) GenerateQREvent(qrText string) (*QREvent, error) {
	base64QR, err := qm.GenerateQRCodeBase64(qrText)
	if err != nil {
		return nil, err
	}

	return &QREvent{
		Event:     "qr",
		Code:      qrText,
		Base64:    base64QR,
		ExpiresAt: time.Now().Add(2 * time.Minute), // QR codes expire in 2 minutes
	}, nil
}

// IsQRExpired checks if a QR code has expired
func (qm *QRManager) IsQRExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

// GetQRCodeForSession retrieves the current QR code for a session
func (wac *WAClient) GetQRCodeForSession(ctx context.Context, sessionID string) (*QREvent, error) {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if session is already connected
	if client.Status == StatusConnected {
		return nil, &WAError{
			Code:    "ALREADY_CONNECTED",
			Message: "session is already connected",
		}
	}

	// Check if QR code exists and is not expired
	if client.QRCode != "" && !client.QRExpiresAt.IsZero() {
		qrManager := NewQRManager(wac.logger)
		if !qrManager.IsQRExpired(client.QRExpiresAt) {
			// Generate base64 for existing QR code
			base64QR, err := qrManager.GenerateQRCodeBase64(client.QRCode)
			if err != nil {
				wac.logger.Error().
					Err(err).
					Str("session_id", sessionID).
					Msg("Failed to generate base64 for existing QR")
			}

			return &QREvent{
				Event:     "qr",
				Code:      client.QRCode,
				Base64:    base64QR,
				ExpiresAt: client.QRExpiresAt,
			}, nil
		}
	}

	// If no valid QR code exists, try to generate a new one
	if client.Status != StatusQRCode && client.Status != StatusConnecting {
		// Start connection process to generate QR
		err := wac.ConnectSession(ctx, sessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to start connection for QR generation: %w", err)
		}

		// Wait a bit for QR code to be generated
		time.Sleep(1 * time.Second)

		// Try to get the updated client
		client, err = wac.GetSession(ctx, sessionID)
		if err != nil {
			return nil, err
		}

		if client.QRCode != "" {
			qrManager := NewQRManager(wac.logger)
			base64QR, err := qrManager.GenerateQRCodeBase64(client.QRCode)
			if err != nil {
				wac.logger.Error().
					Err(err).
					Str("session_id", sessionID).
					Msg("Failed to generate base64 for new QR")
			}

			return &QREvent{
				Event:     "qr",
				Code:      client.QRCode,
				Base64:    base64QR,
				ExpiresAt: client.QRExpiresAt,
			}, nil
		}
	}

	return nil, &WAError{
		Code:    "QR_NOT_AVAILABLE",
		Message: "QR code not available, try connecting the session first",
	}
}

// RefreshQRCode forces a refresh of the QR code for a session
func (wac *WAClient) RefreshQRCode(ctx context.Context, sessionID string) (*QREvent, error) {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if session is already connected
	if client.Status == StatusConnected {
		return nil, &WAError{
			Code:    "ALREADY_CONNECTED",
			Message: "session is already connected",
		}
	}

	// Disconnect and reconnect to force new QR generation
	if client.Status != StatusDisconnected {
		err = wac.DisconnectSession(ctx, sessionID)
		if err != nil {
			wac.logger.Warn().
				Err(err).
				Str("session_id", sessionID).
				Msg("Failed to disconnect session for QR refresh")
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Start connection to generate new QR
	err = wac.ConnectSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect session for QR refresh: %w", err)
	}

	// Wait for QR code generation
	maxWait := 10 * time.Second
	checkInterval := 500 * time.Millisecond
	elapsed := time.Duration(0)

	for elapsed < maxWait {
		time.Sleep(checkInterval)
		elapsed += checkInterval

		// Get updated client
		client, err = wac.GetSession(ctx, sessionID)
		if err != nil {
			continue
		}

		if client.QRCode != "" && client.Status == StatusQRCode {
			qrManager := NewQRManager(wac.logger)
			base64QR, err := qrManager.GenerateQRCodeBase64(client.QRCode)
			if err != nil {
				wac.logger.Error().
					Err(err).
					Str("session_id", sessionID).
					Msg("Failed to generate base64 for refreshed QR")
			}

			return &QREvent{
				Event:     "qr",
				Code:      client.QRCode,
				Base64:    base64QR,
				ExpiresAt: client.QRExpiresAt,
			}, nil
		}
	}

	return nil, &WAError{
		Code:    "QR_GENERATION_TIMEOUT",
		Message: "timeout waiting for QR code generation",
	}
}

// CleanupExpiredQRCodes removes expired QR codes from sessions
func (wac *WAClient) CleanupExpiredQRCodes(ctx context.Context) error {
	sessions, err := wac.ListSessions(ctx)
	if err != nil {
		return fmt.Errorf("failed to list sessions for QR cleanup: %w", err)
	}

	qrManager := NewQRManager(wac.logger)
	cleanedCount := 0

	for _, client := range sessions {
		if client.QRCode != "" && !client.QRExpiresAt.IsZero() {
			if qrManager.IsQRExpired(client.QRExpiresAt) {
				// Clear expired QR code
				client.QRCode = ""
				client.QRExpiresAt = time.Time{}

				// Update in database
				wac.updateSessionStatus(ctx, client)
				cleanedCount++

				wac.logger.Debug().
					Str("session_name", client.Name).
					Msg("Cleaned expired QR code for session")
			}
		}
	}

	if cleanedCount > 0 {
		wac.logger.Info().
			Int("count", cleanedCount).
			Msg("Cleaned expired QR codes")
	}

	return nil
}

// StartQRCleanupRoutine starts a background routine to clean expired QR codes
func (wac *WAClient) StartQRCleanupRoutine(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	wac.logger.Info().Msg("Started QR code cleanup routine")

	for {
		select {
		case <-ctx.Done():
			wac.logger.Info().Msg("Stopping QR code cleanup routine")
			return
		case <-ticker.C:
			if err := wac.CleanupExpiredQRCodes(ctx); err != nil {
				wac.logger.Error().
					Err(err).
					Msg("QR cleanup error")
			}
		}
	}
}
