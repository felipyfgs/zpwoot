package waclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"zpwoot/internal/adapters/logger"

	"github.com/jmoiron/sqlx"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// WAClient is the main WhatsApp client manager
type WAClient struct {
	sessions      map[string]*Client
	sessionsMutex sync.RWMutex
	container     *sqlstore.Container
	logger        *logger.Logger
	eventHandler  EventHandler
	mediaProcessor MediaProcessor
	webhookSender WebhookSender
	sessionRepo   SessionRepository
}

// SessionRepository defines the interface for session persistence
type SessionRepository interface {
	GetSession(ctx context.Context, sessionID string) (*SessionInfo, error)
	GetSessionByName(ctx context.Context, name string) (*SessionInfo, error)
	CreateSession(ctx context.Context, session *SessionInfo) error
	UpdateSession(ctx context.Context, session *SessionInfo) error
	DeleteSession(ctx context.Context, sessionID string) error
	ListSessions(ctx context.Context) ([]*SessionInfo, error)
}

// NewWAClient creates a new WhatsApp client manager
func NewWAClient(
	container *sqlstore.Container,
	logger *logger.Logger,
	sessionRepo SessionRepository,
) *WAClient {
	return &WAClient{
		sessions:    make(map[string]*Client),
		container:   container,
		logger:      logger,
		sessionRepo: sessionRepo,
	}
}

// SetEventHandler sets the event handler
func (wac *WAClient) SetEventHandler(handler EventHandler) {
	wac.eventHandler = handler
}

// SetMediaProcessor sets the media processor
func (wac *WAClient) SetMediaProcessor(processor MediaProcessor) {
	wac.mediaProcessor = processor
}

// SetWebhookSender sets the webhook sender
func (wac *WAClient) SetWebhookSender(sender WebhookSender) {
	wac.webhookSender = sender
}

// CreateSession creates a new WhatsApp session
func (wac *WAClient) CreateSession(ctx context.Context, config *SessionConfig) (*Client, error) {
	wac.sessionsMutex.Lock()
	defer wac.sessionsMutex.Unlock()

	// Check if session already exists
	if _, exists := wac.sessions[config.SessionID]; exists {
		return nil, ErrSessionExists
	}

	// Check if session name already exists
	if existingSession, err := wac.sessionRepo.GetSessionByName(ctx, config.Name); err == nil && existingSession != nil {
		return nil, ErrSessionExists
	}

	// Create device store
	deviceStore, err := wac.container.GetFirstDevice(ctx)
	if err != nil {
		deviceStore = wac.container.NewDevice()
	}

	// Create WhatsApp client
	waClient := whatsmeow.NewClient(deviceStore, waLog.Noop)

	// Create client context
	clientCtx, cancel := context.WithCancel(ctx)

	// Create client instance
	client := &Client{
		SessionID:   config.SessionID,
		Name:        config.Name,
		WAClient:    waClient,
		Status:      StatusDisconnected,
		Config:      config,
		Events:      config.Events,
		WebhookURL:  config.WebhookURL,
		ctx:         clientCtx,
		cancel:      cancel,
	}

	// Register event handler
	client.EventHandler = waClient.AddEventHandler(wac.createEventHandler(client))

	// Store session
	wac.sessions[config.SessionID] = client

	// Persist session to database
	sessionInfo := &SessionInfo{
		ID:        config.SessionID,
		Name:      config.Name,
		Status:    StatusDisconnected,
		Connected: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := wac.sessionRepo.CreateSession(ctx, sessionInfo); err != nil {
		delete(wac.sessions, config.SessionID)
		return nil, fmt.Errorf("failed to persist session: %w", err)
	}

	wac.logger.Infof("Created WhatsApp session: %s (%s)", config.Name, config.SessionID)
	return client, nil
}

// GetSession retrieves a session by ID
func (wac *WAClient) GetSession(ctx context.Context, sessionID string) (*Client, error) {
	wac.sessionsMutex.RLock()
	defer wac.sessionsMutex.RUnlock()

	client, exists := wac.sessions[sessionID]
	if !exists {
		return nil, ErrSessionNotFound
	}

	return client, nil
}

// GetSessionByName retrieves a session by name
func (wac *WAClient) GetSessionByName(ctx context.Context, name string) (*Client, error) {
	wac.sessionsMutex.RLock()
	defer wac.sessionsMutex.RUnlock()

	for _, client := range wac.sessions {
		if client.Name == name {
			return client, nil
		}
	}

	return nil, ErrSessionNotFound
}

// ListSessions returns all active sessions
func (wac *WAClient) ListSessions(ctx context.Context) ([]*Client, error) {
	wac.sessionsMutex.RLock()
	defer wac.sessionsMutex.RUnlock()

	sessions := make([]*Client, 0, len(wac.sessions))
	for _, client := range wac.sessions {
		sessions = append(sessions, client)
	}

	return sessions, nil
}

// ConnectSession connects a session to WhatsApp
func (wac *WAClient) ConnectSession(ctx context.Context, sessionID string) error {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if client.Status == StatusConnected {
		return nil // Already connected
	}

	client.Status = StatusConnecting
	wac.updateSessionStatus(ctx, client)

	// Connect to WhatsApp
	if client.WAClient.Store.ID == nil {
		// First time connection - need QR code
		qrChan, err := client.WAClient.GetQRChannel(ctx)
		if err != nil {
			client.Status = StatusError
			wac.updateSessionStatus(ctx, client)
			return fmt.Errorf("failed to get QR channel: %w", err)
		}

		err = client.WAClient.Connect()
		if err != nil {
			client.Status = StatusError
			wac.updateSessionStatus(ctx, client)
			return fmt.Errorf("failed to connect: %w", err)
		}

		// Handle QR code
		go wac.handleQRCode(ctx, client, qrChan)
	} else {
		// Reconnection
		err = client.WAClient.Connect()
		if err != nil {
			client.Status = StatusError
			wac.updateSessionStatus(ctx, client)
			return fmt.Errorf("failed to reconnect: %w", err)
		}
	}

	return nil
}

// DisconnectSession disconnects a session
func (wac *WAClient) DisconnectSession(ctx context.Context, sessionID string) error {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if client.Status == StatusDisconnected {
		return nil // Already disconnected
	}

	client.WAClient.Disconnect()
	client.Status = StatusDisconnected
	client.cancel()

	wac.updateSessionStatus(ctx, client)
	wac.logger.Infof("Disconnected WhatsApp session: %s", client.Name)

	return nil
}

// DeleteSession deletes a session
func (wac *WAClient) DeleteSession(ctx context.Context, sessionID string) error {
	wac.sessionsMutex.Lock()
	defer wac.sessionsMutex.Unlock()

	client, exists := wac.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}

	// Disconnect if connected
	if client.Status != StatusDisconnected {
		client.WAClient.Disconnect()
		client.cancel()
	}

	// Remove event handler
	client.WAClient.RemoveEventHandler(client.EventHandler)

	// Delete from memory
	delete(wac.sessions, sessionID)

	// Delete from database
	if err := wac.sessionRepo.DeleteSession(ctx, sessionID); err != nil {
		wac.logger.Errorf("Failed to delete session from database: %v", err)
	}

	wac.logger.Infof("Deleted WhatsApp session: %s", client.Name)
	return nil
}

// createEventHandler creates an event handler for a client
func (wac *WAClient) createEventHandler(client *Client) func(interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Connected:
			wac.handleConnected(client, v)
		case *events.Disconnected:
			wac.handleDisconnected(client, v)
		case *events.LoggedOut:
			wac.handleLoggedOut(client, v)
		case *events.QR:
			wac.handleQR(client, v)
		case *events.Message:
			wac.handleMessage(client, v)
		default:
			// Handle other events through the event handler
			if wac.eventHandler != nil {
				if err := wac.eventHandler.HandleEvent(client, evt); err != nil {
					wac.logger.Errorf("Event handler error: %v", err)
				}
			}
		}
	}
}

// Event handlers
func (wac *WAClient) handleConnected(client *Client, evt *events.Connected) {
	client.Status = StatusConnected
	client.ConnectedAt = time.Now()
	client.LastSeen = time.Now()

	if client.WAClient.Store.ID != nil {
		wac.logger.Infof("Session connected: %s (%s)", client.Name, client.WAClient.Store.ID.String())
	}

	wac.updateSessionStatus(context.Background(), client)

	// Send webhook if configured
	if wac.webhookSender != nil && client.WebhookURL != "" {
		webhookEvent := &WebhookEvent{
			Type:      EventConnected,
			SessionID: client.SessionID,
			Event:     evt,
			Timestamp: time.Now(),
		}
		go wac.webhookSender.SendWebhook(context.Background(), webhookEvent)
	}
}

func (wac *WAClient) handleDisconnected(client *Client, evt *events.Disconnected) {
	client.Status = StatusDisconnected
	client.LastSeen = time.Now()

	wac.logger.Infof("Session disconnected: %s", client.Name)
	wac.updateSessionStatus(context.Background(), client)

	// Send webhook if configured
	if wac.webhookSender != nil && client.WebhookURL != "" {
		webhookEvent := &WebhookEvent{
			Type:      EventDisconnected,
			SessionID: client.SessionID,
			Event:     evt,
			Timestamp: time.Now(),
		}
		go wac.webhookSender.SendWebhook(context.Background(), webhookEvent)
	}
}

func (wac *WAClient) handleLoggedOut(client *Client, evt *events.LoggedOut) {
	client.Status = StatusDisconnected
	client.LastSeen = time.Now()

	wac.logger.Infof("Session logged out: %s", client.Name)
	wac.updateSessionStatus(context.Background(), client)

	// Send webhook if configured
	if wac.webhookSender != nil && client.WebhookURL != "" {
		webhookEvent := &WebhookEvent{
			Type:      EventLoggedOut,
			SessionID: client.SessionID,
			Event:     evt,
			Timestamp: time.Now(),
		}
		go wac.webhookSender.SendWebhook(context.Background(), webhookEvent)
	}
}

func (wac *WAClient) handleQR(client *Client, evt *events.QR) {
	client.Status = StatusQRCode
	client.QRCode = evt.Codes[0]
	client.QRExpiresAt = time.Now().Add(2 * time.Minute) // QR expires in 2 minutes

	wac.logger.Infof("QR code generated for session: %s", client.Name)
	wac.updateSessionStatus(context.Background(), client)

	// Send webhook if configured
	if wac.webhookSender != nil && client.WebhookURL != "" {
		qrEvent := &QREvent{
			Event:     "qr",
			Code:      evt.Codes[0],
			ExpiresAt: client.QRExpiresAt,
		}

		webhookEvent := &WebhookEvent{
			Type:      EventQR,
			SessionID: client.SessionID,
			Event:     qrEvent,
			Timestamp: time.Now(),
		}
		go wac.webhookSender.SendWebhook(context.Background(), webhookEvent)
	}
}

func (wac *WAClient) handleMessage(client *Client, evt *events.Message) {
	client.LastSeen = time.Now()

	wac.logger.Debugf("Message received in session %s: %s", client.Name, evt.Info.ID)

	// Send webhook if configured
	if wac.webhookSender != nil && client.WebhookURL != "" {
		webhookEvent := &WebhookEvent{
			Type:      EventMessage,
			SessionID: client.SessionID,
			Event:     evt,
			Timestamp: time.Now(),
		}
		go wac.webhookSender.SendWebhook(context.Background(), webhookEvent)
	}
}

func (wac *WAClient) handleQRCode(ctx context.Context, client *Client, qrChan <-chan whatsmeow.QRChannelItem) {
	for evt := range qrChan {
		if evt.Event == "code" {
			client.Status = StatusQRCode
			client.QRCode = evt.Code
			client.QRExpiresAt = time.Now().Add(2 * time.Minute)

			wac.logger.Infof("QR code updated for session: %s", client.Name)
			wac.updateSessionStatus(ctx, client)

			// Send webhook if configured
			if wac.webhookSender != nil && client.WebhookURL != "" {
				qrEvent := &QREvent{
					Event:     evt.Event,
					Code:      evt.Code,
					ExpiresAt: client.QRExpiresAt,
				}

				webhookEvent := &WebhookEvent{
					Type:      EventQR,
					SessionID: client.SessionID,
					Event:     qrEvent,
					Timestamp: time.Now(),
				}
				go wac.webhookSender.SendWebhook(ctx, webhookEvent)
			}
		} else {
			wac.logger.Infof("QR channel event: %s", evt.Event)
		}
	}
}

// updateSessionStatus updates session status in database
func (wac *WAClient) updateSessionStatus(ctx context.Context, client *Client) {
	deviceJID := ""
	if client.WAClient.Store.ID != nil {
		deviceJID = client.WAClient.Store.ID.String()
	}

	sessionInfo := &SessionInfo{
		ID:          client.SessionID,
		Name:        client.Name,
		DeviceJID:   deviceJID,
		Status:      client.Status,
		Connected:   client.Status == StatusConnected,
		QRCode:      client.QRCode,
		QRExpiresAt: client.QRExpiresAt,
		ConnectedAt: client.ConnectedAt,
		LastSeen:    client.LastSeen,
		UpdatedAt:   time.Now(),
	}

	if err := wac.sessionRepo.UpdateSession(ctx, sessionInfo); err != nil {
		wac.logger.Errorf("Failed to update session status: %v", err)
	}
}

// NewWAStoreContainer creates a new WhatsApp store container
func NewWAStoreContainer(db *sqlx.DB, logger *logger.Logger) *sqlstore.Container {
	// Get database URL from connection
	dbURL := "postgres://zpwoot:zpwoot123@localhost:5432/zpwoot?sslmode=disable"

	// Create WhatsApp store container
	container, err := sqlstore.New(context.Background(), "postgres", dbURL, waLog.Noop)
	if err != nil {
		logger.Errorf("Failed to create WhatsApp store container: %v", err)
		return nil
	}

	return container
}
