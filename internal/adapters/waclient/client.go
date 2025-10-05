package waclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"sync"
	"time"

	"zpwoot/internal/adapters/logger"

	"github.com/jmoiron/sqlx"
	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// WAClient is the main WhatsApp client manager
type WAClient struct {
	sessions       map[string]*Client
	sessionsMutex  sync.RWMutex
	container      *sqlstore.Container
	logger         *logger.Logger
	eventHandler   EventHandler
	mediaProcessor MediaProcessor
	webhookSender  WebhookSender
	sessionRepo    SessionRepository
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
	wac := &WAClient{
		sessions:    make(map[string]*Client),
		container:   container,
		logger:      logger,
		sessionRepo: sessionRepo,
	}

	// Load existing sessions from database
	go wac.loadSessionsFromDatabase()

	return wac
}

// loadSessionsFromDatabase loads all sessions from database into memory
func (wac *WAClient) loadSessionsFromDatabase() {
	ctx := context.Background()
	sessions, err := wac.sessionRepo.ListSessions(ctx)
	if err != nil {
		wac.logger.Error().
			Err(err).
			Msg("Failed to load sessions from database")
		return
	}

	wac.logger.Info().
		Int("count", len(sessions)).
		Msg("Loading sessions from database")

	for _, sessionInfo := range sessions {
		// Get device store by JID if available
		var deviceStore *store.Device
		var err error

		if sessionInfo.DeviceJID != "" {
			// Try to get existing device by JID
			jid, parseErr := types.ParseJID(sessionInfo.DeviceJID)
			if parseErr == nil {
				deviceStore, err = wac.container.GetDevice(ctx, jid)
				if err != nil {
					wac.logger.Warn().
						Err(err).
						Str("jid", sessionInfo.DeviceJID).
						Msg("Failed to get device by JID, creating new one")
					deviceStore = wac.container.NewDevice()
				}
			} else {
				wac.logger.Warn().
					Err(parseErr).
					Str("jid", sessionInfo.DeviceJID).
					Msg("Failed to parse JID, creating new device")
				deviceStore = wac.container.NewDevice()
			}
		} else {
			deviceStore = wac.container.NewDevice()
		}

		// Create WhatsApp client
		waClient := whatsmeow.NewClient(deviceStore, waLog.Noop)

		// Create client context
		clientCtx, cancel := context.WithCancel(ctx)

		// Create client instance
		client := &Client{
			SessionID:   sessionInfo.ID,
			Name:        sessionInfo.Name,
			WAClient:    waClient,
			Status:      sessionInfo.Status,
			QRCode:      sessionInfo.QRCode,
			QRExpiresAt: sessionInfo.QRExpiresAt,
			ConnectedAt: sessionInfo.ConnectedAt,
			LastSeen:    sessionInfo.LastSeen,
			Config: &SessionConfig{
				SessionID: sessionInfo.ID,
				Name:      sessionInfo.Name,
			},
			ctx:    clientCtx,
			cancel: cancel,
		}

		// Register event handler
		client.EventHandler = waClient.AddEventHandler(wac.createEventHandler(client))

		// Store session in memory
		wac.sessionsMutex.Lock()
		wac.sessions[sessionInfo.ID] = client
		wac.sessionsMutex.Unlock()

		wac.logger.Info().
			Str("session_id", sessionInfo.ID).
			Str("name", sessionInfo.Name).
			Str("status", string(sessionInfo.Status)).
			Msg("Loaded session from database")

		// Auto-reconnect if session was connected
		if sessionInfo.Connected && sessionInfo.DeviceJID != "" {
			wac.logger.Info().
				Str("session_id", sessionInfo.ID).
				Str("name", sessionInfo.Name).
				Str("device_jid", sessionInfo.DeviceJID).
				Msg("Auto-reconnecting session")

			// Reconnect in background
			go func(c *Client) {
				time.Sleep(2 * time.Second) // Wait a bit for initialization

				err := c.WAClient.Connect()
				if err != nil {
					wac.logger.Error().
						Err(err).
						Str("session_id", c.SessionID).
						Msg("Failed to auto-reconnect session")
					c.Status = StatusDisconnected
					wac.updateSessionStatus(context.Background(), c)
				} else {
					wac.logger.Info().
						Str("session_id", c.SessionID).
						Str("name", c.Name).
						Msg("Session auto-reconnected successfully")
					c.Status = StatusConnected
					c.ConnectedAt = time.Now()
					c.LastSeen = time.Now()
					wac.updateSessionStatus(context.Background(), c)
				}
			}(client)
		}
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
		SessionID:  config.SessionID,
		Name:       config.Name,
		WAClient:   waClient,
		Status:     StatusDisconnected,
		Config:     config,
		Events:     config.Events,
		WebhookURL: config.WebhookURL,
		ctx:        clientCtx,
		cancel:     cancel,
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

	wac.logger.Info().
		Str("name", config.Name).
		Str("session_id", config.SessionID).
		Msg("Created WhatsApp session")
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
		qrChan, err := client.WAClient.GetQRChannel(context.Background())
		if err != nil {
			client.Status = StatusError
			wac.updateSessionStatus(context.Background(), client)
			return fmt.Errorf("failed to get QR channel: %w", err)
		}

		err = client.WAClient.Connect()
		if err != nil {
			client.Status = StatusError
			wac.updateSessionStatus(context.Background(), client)
			return fmt.Errorf("failed to connect: %w", err)
		}

		// Handle QR code using client's context (not HTTP request context)
		go wac.handleQRCode(client.ctx, client, qrChan)
	} else {
		// Reconnection
		err = client.WAClient.Connect()
		if err != nil {
			client.Status = StatusError
			wac.updateSessionStatus(context.Background(), client)
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
	wac.logger.Info().
		Str("name", client.Name).
		Msg("Disconnected WhatsApp session")

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
		wac.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to delete session from database")
	}

	wac.logger.Info().
		Str("name", client.Name).
		Msg("Deleted WhatsApp session")
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
					wac.logger.Error().
						Err(err).
						Str("session_id", client.SessionID).
						Msg("Event handler error")
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
		wac.logger.Info().
			Str("name", client.Name).
			Str("store_id", client.WAClient.Store.ID.String()).
			Msg("Session connected")
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

	wac.logger.Info().
		Str("name", client.Name).
		Msg("Session disconnected")
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

	wac.logger.Info().
		Str("name", client.Name).
		Msg("Session logged out")
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

	wac.logger.Info().
		Str("name", client.Name).
		Msg("QR code generated for session")
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

	// Log message details
	wac.logger.Info().
		Str("session_name", client.Name).
		Str("session_id", client.SessionID).
		Str("message_id", evt.Info.ID).
		Str("from", evt.Info.Sender.String()).
		Str("chat", evt.Info.Chat.String()).
		Bool("from_me", evt.Info.IsFromMe).
		Bool("is_group", evt.Info.IsGroup).
		Str("push_name", evt.Info.PushName).
		Time("timestamp", evt.Info.Timestamp).
		Msg("📨 Message received")

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

	// Call event handler if configured
	if wac.eventHandler != nil {
		if err := wac.eventHandler.HandleEvent(client, evt); err != nil {
			wac.logger.Error().
				Err(err).
				Str("session_id", client.SessionID).
				Msg("Event handler error for message")
		}
	}
}

func (wac *WAClient) handleQRCode(ctx context.Context, client *Client, qrChan <-chan whatsmeow.QRChannelItem) {
	for evt := range qrChan {
		if evt.Event == "code" {
			client.Status = StatusQRCode
			client.QRCode = evt.Code
			client.QRExpiresAt = time.Now().Add(2 * time.Minute)

			// Display QR code in terminal
			wac.logger.Info().
				Str("name", client.Name).
				Str("session_id", client.SessionID).
				Msg("QR code generated for session")

			// Print QR code to terminal for easy scanning
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			fmt.Printf("\n🔗 QR Code for session '%s' (ID: %s)\n", client.Name, client.SessionID)
			fmt.Printf("📱 Scan this QR code with WhatsApp to connect\n")
			fmt.Printf("⏰ Expires at: %s\n\n", client.QRExpiresAt.Format("2006-01-02 15:04:05"))

			// Generate base64 QR code image for API/webhook
			qrImage, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
			if err != nil {
				wac.logger.Error().
					Err(err).
					Str("session_id", client.SessionID).
					Msg("Failed to generate QR code image")
			} else {
				// Store base64 encoded QR code
				base64QR := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrImage)
				client.QRCode = base64QR // Store base64 version for API

				wac.logger.Info().
					Str("name", client.Name).
					Str("session_id", client.SessionID).
					Msg("QR code image generated and encoded to base64")
			}

			// Update session status in database
			wac.updateSessionStatus(context.Background(), client)

			// Send webhook if configured
			if wac.webhookSender != nil && client.WebhookURL != "" {
				qrEvent := &QREvent{
					Event:     evt.Event,
					Code:      evt.Code,
					Base64:    client.QRCode, // Send base64 version
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
		} else if evt.Event == "success" {
			wac.logger.Info().
				Str("name", client.Name).
				Str("session_id", client.SessionID).
				Msg("QR code scanned successfully - session connected")

			// Clear QR code after successful pairing
			client.QRCode = ""
			client.QRExpiresAt = time.Time{}
			client.Status = StatusConnected
			client.ConnectedAt = time.Now()

			wac.updateSessionStatus(context.Background(), client)
		} else if evt.Event == "timeout" {
			wac.logger.Warn().
				Str("name", client.Name).
				Str("session_id", client.SessionID).
				Msg("QR code expired - please request a new one")

			// Clear expired QR code
			client.QRCode = ""
			client.QRExpiresAt = time.Time{}
			client.Status = StatusDisconnected

			wac.updateSessionStatus(context.Background(), client)
		} else {
			wac.logger.Info().
				Str("event", evt.Event).
				Str("session_id", client.SessionID).
				Msg("QR channel event")
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
		wac.logger.Error().
			Err(err).
			Str("session_id", client.SessionID).
			Msg("Failed to update session status")
	}
}

// NewWAStoreContainer creates a new WhatsApp store container
func NewWAStoreContainer(db *sqlx.DB, logger *logger.Logger) *sqlstore.Container {
	// Get database URL from connection
	dbURL := "postgres://zpwoot:zpwoot123@localhost:5432/zpwoot?sslmode=disable"

	// Create WhatsApp store container
	container, err := sqlstore.New(context.Background(), "postgres", dbURL, waLog.Noop)
	if err != nil {
		logger.Error().
			Err(err).
			Str("db_url", dbURL).
			Msg("Failed to create WhatsApp store container")
		return nil
	}

	return container
}
