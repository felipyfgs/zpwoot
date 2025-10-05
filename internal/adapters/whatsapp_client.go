package adapters

import (
	"context"
	"fmt"
	"sync"

	"github.com/zpwoot/internal/domain/session"
	"github.com/zpwoot/internal/infra/waclient"
	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// WhatsAppClientAdapter adapts the WhatsApp client for domain use
type WhatsAppClientAdapter struct {
	clients map[string]*waclient.WAClient
	mutex   sync.RWMutex
	logger  waLog.Logger
	dbPath  string
}

// NewWhatsAppClientAdapter creates a new WhatsApp client adapter
func NewWhatsAppClientAdapter(dbPath string, logger waLog.Logger) *WhatsAppClientAdapter {
	return &WhatsAppClientAdapter{
		clients: make(map[string]*waclient.WAClient),
		logger:  logger,
		dbPath:  dbPath,
	}
}

// CreateSession creates a new WhatsApp session
func (w *WhatsAppClientAdapter) CreateSession(ctx context.Context, sessionID string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Check if session already exists
	if _, exists := w.clients[sessionID]; exists {
		return fmt.Errorf("session %s already exists", sessionID)
	}

	// Create new WhatsApp client
	client, err := waclient.NewWAClient(w.dbPath, w.logger)
	if err != nil {
		return fmt.Errorf("failed to create WhatsApp client: %w", err)
	}

	// Initialize session
	_, err = client.CreateSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to initialize session: %w", err)
	}

	w.clients[sessionID] = client
	return nil
}

// ConnectSession connects a WhatsApp session and returns QR code if needed
func (w *WhatsAppClientAdapter) ConnectSession(ctx context.Context, sessionID string) (string, error) {
	w.mutex.RLock()
	client, exists := w.clients[sessionID]
	w.mutex.RUnlock()

	if !exists {
		return "", fmt.Errorf("session %s not found", sessionID)
	}

	// Try to connect
	if err := client.Connect(ctx); err != nil {
		return "", fmt.Errorf("failed to connect session: %w", err)
	}

	// If not authenticated, get QR code
	if !client.IsConnected() {
		qrCode, err := client.GetQRCode(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get QR code: %w", err)
		}
		return qrCode, nil
	}

	return "", nil // Already connected
}

// DisconnectSession disconnects a WhatsApp session
func (w *WhatsAppClientAdapter) DisconnectSession(sessionID string) error {
	w.mutex.RLock()
	client, exists := w.clients[sessionID]
	w.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	client.Disconnect()
	return nil
}

// GetSessionStatus returns the current status of a session
func (w *WhatsAppClientAdapter) GetSessionStatus(sessionID string) (session.Status, error) {
	w.mutex.RLock()
	client, exists := w.clients[sessionID]
	w.mutex.RUnlock()

	if !exists {
		return session.StatusDisconnected, fmt.Errorf("session %s not found", sessionID)
	}

	if client.IsConnected() {
		return session.StatusConnected, nil
	}

	return session.StatusDisconnected, nil
}

// GetSessionJID returns the JID of a connected session
func (w *WhatsAppClientAdapter) GetSessionJID(sessionID string) (string, error) {
	w.mutex.RLock()
	client, exists := w.clients[sessionID]
	w.mutex.RUnlock()

	if !exists {
		return "", fmt.Errorf("session %s not found", sessionID)
	}

	jid := client.GetJID()
	return jid.String(), nil
}

// RemoveSession removes a session from memory
func (w *WhatsAppClientAdapter) RemoveSession(sessionID string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	client, exists := w.clients[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	// Disconnect and close
	client.Disconnect()
	if err := client.Close(); err != nil {
		return fmt.Errorf("failed to close session: %w", err)
	}

	delete(w.clients, sessionID)
	return nil
}

// SendMessage sends a message through a session
func (w *WhatsAppClientAdapter) SendMessage(ctx context.Context, sessionID, to, message string) error {
	w.mutex.RLock()
	client, exists := w.clients[sessionID]
	w.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	if !client.IsConnected() {
		return fmt.Errorf("session %s is not connected", sessionID)
	}

	// Parse recipient JID
	jid, err := waclient.MapStringToJID(to)
	if err != nil {
		return fmt.Errorf("invalid recipient JID: %w", err)
	}

	return client.SendMessage(ctx, jid, message)
}

// ListActiveSessions returns all active session IDs
func (w *WhatsAppClientAdapter) ListActiveSessions() []string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	sessions := make([]string, 0, len(w.clients))
	for sessionID := range w.clients {
		sessions = append(sessions, sessionID)
	}

	return sessions
}

// Close closes all sessions and cleans up resources
func (w *WhatsAppClientAdapter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	var lastErr error
	for sessionID, client := range w.clients {
		client.Disconnect()
		if err := client.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close session %s: %w", sessionID, err)
		}
	}

	w.clients = make(map[string]*waclient.WAClient)
	return lastErr
}
