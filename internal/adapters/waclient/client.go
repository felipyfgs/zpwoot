package waclient

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"zpwoot/internal/adapters/logger"

	"github.com/jmoiron/sqlx"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

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

type SessionRepository interface {
	GetSession(ctx context.Context, sessionID string) (*SessionInfo, error)
	GetSessionByName(ctx context.Context, name string) (*SessionInfo, error)
	CreateSession(ctx context.Context, session *SessionInfo) error
	UpdateSession(ctx context.Context, session *SessionInfo) error
	DeleteSession(ctx context.Context, sessionID string) error
	ListSessions(ctx context.Context) ([]*SessionInfo, error)
}

func NewWAClient(container *sqlstore.Container, logger *logger.Logger, sessionRepo SessionRepository) *WAClient {
	wac := &WAClient{
		sessions:    make(map[string]*Client),
		container:   container,
		logger:      logger,
		sessionRepo: sessionRepo,
	}
	go wac.loadSessionsFromDatabase()
	return wac
}

func (wac *WAClient) loadSessionsFromDatabase() {
	ctx := context.Background()
	sessions, err := wac.sessionRepo.ListSessions(ctx)
	if err != nil {
		wac.logger.Error().Err(err).Msg("Failed to load sessions from database")
		return
	}

	wac.logger.Debug().Int("count", len(sessions)).Msg("Loading sessions from database")

	for _, sessionInfo := range sessions {
		deviceStore := wac.getOrCreateDeviceStore(ctx, sessionInfo.DeviceJID)
		client := wac.createClient(ctx, sessionInfo, deviceStore)

		wac.sessionsMutex.Lock()
		wac.sessions[sessionInfo.ID] = client
		wac.sessionsMutex.Unlock()

		wac.logger.Debug().Str("session_id", sessionInfo.ID).Str("name", sessionInfo.Name).Msg("Loaded session")

		if sessionInfo.Connected && sessionInfo.DeviceJID != "" {
			wac.logger.Debug().Str("session_id", sessionInfo.ID).Msg("Auto-reconnecting")
			go wac.autoReconnect(client)
		}
	}
}

func (wac *WAClient) autoReconnect(client *Client) {
	time.Sleep(2 * time.Second)

	if err := client.WAClient.Connect(); err != nil {
		wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Failed to auto-reconnect session")
		client.Status = StatusDisconnected
	} else {
		wac.logger.Debug().Str("session_id", client.SessionID).Msg("Auto-reconnected")
		client.Status = StatusConnected
		client.ConnectedAt = time.Now()
		client.LastSeen = time.Now()
	}
	wac.updateSessionStatus(context.Background(), client)
}

func (wac *WAClient) SetEventHandler(handler EventHandler) {
	wac.eventHandler = handler
}

func (wac *WAClient) SetMediaProcessor(processor MediaProcessor) {
	wac.mediaProcessor = processor
}

func (wac *WAClient) SetWebhookSender(sender WebhookSender) {
	wac.webhookSender = sender
}

func (wac *WAClient) createClient(ctx context.Context, sessionInfo *SessionInfo, deviceStore *store.Device) *Client {
	waClient := whatsmeow.NewClient(deviceStore, waLog.Noop)
	clientCtx, cancel := context.WithCancel(ctx)

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

	client.EventHandler = waClient.AddEventHandler(wac.createEventHandler(client))
	return client
}

func (wac *WAClient) getOrCreateDeviceStore(ctx context.Context, deviceJID string) *store.Device {
	if deviceJID != "" {
		jid, err := types.ParseJID(deviceJID)
		if err == nil {
			deviceStore, err := wac.container.GetDevice(ctx, jid)
			if err == nil {
				return deviceStore
			}
			wac.logger.Warn().Err(err).Str("jid", deviceJID).Msg("Failed to get device by JID, creating new one")
		} else {
			wac.logger.Warn().Err(err).Str("jid", deviceJID).Msg("Failed to parse JID, creating new device")
		}
	}
	return wac.container.NewDevice()
}

func (wac *WAClient) CreateSession(ctx context.Context, config *SessionConfig) (*Client, error) {
	wac.sessionsMutex.Lock()
	defer wac.sessionsMutex.Unlock()

	if _, exists := wac.sessions[config.SessionID]; exists {
		return nil, ErrSessionExists
	}

	if existingSession, err := wac.sessionRepo.GetSessionByName(ctx, config.Name); err == nil && existingSession != nil {
		return nil, ErrSessionExists
	}

	deviceStore := wac.container.NewDevice()
	wac.logger.Debug().
		Str("session_id", config.SessionID).
		Str("session_name", config.Name).
		Str("device_id", fmt.Sprintf("%p", deviceStore)).
		Msg("Created new device for session")

	sessionInfo := &SessionInfo{
		ID:        config.SessionID,
		Name:      config.Name,
		Status:    StatusDisconnected,
		Connected: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	client := wac.createClient(ctx, sessionInfo, deviceStore)
	client.Config = config
	client.Events = config.Events
	client.WebhookURL = config.WebhookURL

	wac.sessions[config.SessionID] = client

	if err := wac.sessionRepo.CreateSession(ctx, sessionInfo); err != nil {
		delete(wac.sessions, config.SessionID)
		return nil, fmt.Errorf("failed to persist session: %w", err)
	}

	wac.logger.Info().Str("session_id", config.SessionID).Msg("Session created")
	return client, nil
}

func (wac *WAClient) GetSession(ctx context.Context, sessionID string) (*Client, error) {
	wac.sessionsMutex.RLock()
	defer wac.sessionsMutex.RUnlock()

	client, exists := wac.sessions[sessionID]
	if !exists {
		return nil, ErrSessionNotFound
	}
	return client, nil
}

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

func (wac *WAClient) ListSessions(ctx context.Context) ([]*Client, error) {
	wac.sessionsMutex.RLock()
	defer wac.sessionsMutex.RUnlock()

	sessions := make([]*Client, 0, len(wac.sessions))
	for _, client := range wac.sessions {
		sessions = append(sessions, client)
	}
	return sessions, nil
}

func (wac *WAClient) ConnectSession(ctx context.Context, sessionID string) error {
	client, err := wac.GetSession(ctx, sessionID)
	if err == ErrSessionNotFound {
		sessionInfo, dbErr := wac.sessionRepo.GetSession(ctx, sessionID)
		if dbErr != nil {
			return fmt.Errorf("session not found: %w", dbErr)
		}
		client, err = wac.recreateClient(ctx, sessionInfo)
		if err != nil {
			return fmt.Errorf("failed to recreate client: %w", err)
		}
	} else if err != nil {
		return err
	}

	if client.Status == StatusConnected {
		return nil
	}

	client.Status = StatusConnecting
	wac.updateSessionStatus(ctx, client)

	if client.WAClient.Store.ID == nil {
		qrChan, err := client.WAClient.GetQRChannel(context.Background())
		if err != nil {
			client.Status = StatusError
			wac.updateSessionStatus(context.Background(), client)
			return fmt.Errorf("failed to get QR channel: %w", err)
		}

		if err = client.WAClient.Connect(); err != nil {
			client.Status = StatusError
			wac.updateSessionStatus(context.Background(), client)
			return fmt.Errorf("failed to connect: %w", err)
		}

		go wac.handleQRCode(client.ctx, client, qrChan)
	} else {
		if err = client.WAClient.Connect(); err != nil {
			client.Status = StatusError
			wac.updateSessionStatus(context.Background(), client)
			return fmt.Errorf("failed to reconnect: %w", err)
		}
	}

	return nil
}

func (wac *WAClient) recreateClient(ctx context.Context, sessionInfo *SessionInfo) (*Client, error) {
	wac.sessionsMutex.Lock()
	defer wac.sessionsMutex.Unlock()

	if existing, exists := wac.sessions[sessionInfo.ID]; exists {
		return existing, nil
	}

	deviceStore := wac.container.NewDevice()
	wac.logger.Debug().
		Str("session_id", sessionInfo.ID).
		Str("session_name", sessionInfo.Name).
		Msg("Recreating client with new device")

	client := wac.createClient(ctx, sessionInfo, deviceStore)
	wac.sessions[sessionInfo.ID] = client

	return client, nil
}

func (wac *WAClient) DisconnectSession(ctx context.Context, sessionID string) error {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if client.Status == StatusDisconnected {
		return nil
	}

	client.WAClient.Disconnect()
	client.Status = StatusDisconnected
	client.cancel()

	wac.updateSessionStatus(ctx, client)
	wac.logger.Info().Str("session_id", client.SessionID).Msg("Session disconnected (credentials kept for reconnection)")

	return nil
}

func (wac *WAClient) LogoutSession(ctx context.Context, sessionID string) error {
	wac.sessionsMutex.Lock()
	defer wac.sessionsMutex.Unlock()

	client, exists := wac.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}

	if client.WAClient.Store.ID != nil {
		if err := client.WAClient.Logout(ctx); err != nil {
			wac.logger.Warn().Err(err).Str("session_id", client.SessionID).Msg("Logout request failed, but continuing with local cleanup")
		}
	} else {
		client.WAClient.Disconnect()
	}

	client.Status = StatusDisconnected
	client.cancel()
	client.WAClient.RemoveEventHandler(client.EventHandler)
	delete(wac.sessions, sessionID)

	sessionInfo := &SessionInfo{
		ID:          client.SessionID,
		Name:        client.Name,
		DeviceJID:   "",
		Status:      StatusDisconnected,
		Connected:   false,
		QRCode:      "",
		QRExpiresAt: time.Time{},
		UpdatedAt:   time.Now(),
		LastSeen:    time.Now(),
	}

	if err := wac.sessionRepo.UpdateSession(ctx, sessionInfo); err != nil {
		wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Failed to update session after logout")
	}

	wac.logger.Info().Str("session_id", client.SessionID).Msg("Session logged out (device unlinked from WhatsApp)")
	return nil
}

func (wac *WAClient) DeleteSession(ctx context.Context, sessionID string) error {
	wac.sessionsMutex.Lock()
	defer wac.sessionsMutex.Unlock()

	client, exists := wac.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}

	if client.Status != StatusDisconnected {
		if client.WAClient.Store.ID != nil {
			if err := client.WAClient.Logout(ctx); err != nil {
				wac.logger.Warn().Err(err).Str("session_id", sessionID).Msg("Logout failed during delete, continuing anyway")
			}
		} else {
			client.WAClient.Disconnect()
		}
		client.cancel()
	}

	client.WAClient.RemoveEventHandler(client.EventHandler)
	delete(wac.sessions, sessionID)

	if err := wac.sessionRepo.DeleteSession(ctx, sessionID); err != nil {
		wac.logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to delete session from database")
	}

	wac.logger.Info().Str("session_id", sessionID).Msg("Session deleted")
	return nil
}

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
			if wac.eventHandler != nil {
				if err := wac.eventHandler.HandleEvent(client, evt); err != nil {
					wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Event handler error")
				}
			}
		}
	}
}

func (wac *WAClient) handleConnected(client *Client, evt *events.Connected) {
	client.Status = StatusConnected
	client.ConnectedAt = time.Now()
	client.LastSeen = time.Now()

	if client.WAClient.Store.ID != nil {
		wac.logger.Info().Str("session_id", client.SessionID).Msg("Connected")
	}

	wac.updateSessionStatus(context.Background(), client)
	wac.sendWebhook(client, EventConnected, evt)
}

func (wac *WAClient) handleDisconnected(client *Client, evt *events.Disconnected) {
	client.Status = StatusDisconnected
	client.LastSeen = time.Now()

	wac.logger.Debug().Str("session_id", client.SessionID).Msg("Disconnected")
	wac.updateSessionStatus(context.Background(), client)
	wac.sendWebhook(client, EventDisconnected, evt)
}

func (wac *WAClient) handleLoggedOut(client *Client, evt *events.LoggedOut) {
	client.Status = StatusDisconnected
	client.LastSeen = time.Now()

	wac.logger.Info().Str("session_id", client.SessionID).Msg("Logged out")
	wac.updateSessionStatus(context.Background(), client)
	wac.sendWebhook(client, EventLoggedOut, evt)
}

func (wac *WAClient) handleQR(client *Client, evt *events.QR) {
	client.Status = StatusQRCode
	client.QRCode = evt.Codes[0]
	client.QRExpiresAt = time.Now().Add(2 * time.Minute)

	wac.logger.Debug().Str("session_id", client.SessionID).Msg("QR generated")
	wac.updateSessionStatus(context.Background(), client)

	qrEvent := &QREvent{
		Event:     "qr",
		Code:      evt.Codes[0],
		ExpiresAt: client.QRExpiresAt,
	}
	wac.sendWebhook(client, EventQR, qrEvent)
}

func (wac *WAClient) handleMessage(client *Client, evt *events.Message) {
	client.LastSeen = time.Now()

	wac.logger.Debug().Str("session_id", client.SessionID).Str("msg_id", evt.Info.ID).Msg("Message")
	wac.sendWebhook(client, EventMessage, evt)

	if wac.eventHandler != nil {
		if err := wac.eventHandler.HandleEvent(client, evt); err != nil {
			wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Event handler error for message")
		}
	}
}

func (wac *WAClient) sendWebhook(client *Client, eventType EventType, event interface{}) {
	if wac.webhookSender != nil && client.WebhookURL != "" {
		webhookEvent := &WebhookEvent{
			Type:      eventType,
			SessionID: client.SessionID,
			Event:     event,
			Timestamp: time.Now(),
		}
		go wac.webhookSender.SendWebhook(context.Background(), webhookEvent)
	}
}

func (wac *WAClient) handleQRCode(ctx context.Context, client *Client, qrChan <-chan whatsmeow.QRChannelItem) {
	for evt := range qrChan {
		switch evt.Event {
		case "code":
			wac.handleQRCodeGenerated(client, evt.Code)
		case "success":
			wac.handleQRSuccess(client)
		case "timeout":
			wac.handleQRTimeout(client)
		case "error":
			wac.handleQRError(client, evt.Event)
		default:
			wac.logger.Debug().Str("event", evt.Event).Str("session_id", client.SessionID).Msg("QR event")
		}
	}

	if client.QRCode != "" && client.Status != StatusConnected {
		wac.logger.Info().Str("session_id", client.SessionID).Msg("QR code channel closed - clearing QR code")
		wac.clearQRCode(client)
		if client.Status == StatusQRCode {
			client.Status = StatusDisconnected
		}
		wac.updateSessionStatus(context.Background(), client)
	}
}

func (wac *WAClient) handleQRCodeGenerated(client *Client, code string) {
	client.Status = StatusQRCode
	client.QRCode = code
	client.QRExpiresAt = time.Now().Add(2 * time.Minute)

	wac.logger.Info().Str("session_id", client.SessionID).Msg("QR code ready - scan to connect")

	qrterminal.GenerateHalfBlock(code, qrterminal.L, os.Stdout)
	fmt.Printf("\n📱 Scan QR code | Session: %s | Expires: %s\n\n",
		client.SessionID, client.QRExpiresAt.Format("15:04:05"))

	wac.updateSessionStatus(context.Background(), client)

	qrEvent := &QREvent{
		Event:     "code",
		Code:      code,
		ExpiresAt: client.QRExpiresAt,
	}
	wac.sendWebhook(client, EventQR, qrEvent)
}

func (wac *WAClient) handleQRSuccess(client *Client) {
	wac.logger.Info().Str("session_id", client.SessionID).Msg("QR scanned - connected")
	wac.clearQRCode(client)
	client.Status = StatusConnected
	client.ConnectedAt = time.Now()
	wac.updateSessionStatus(context.Background(), client)
}

func (wac *WAClient) handleQRTimeout(client *Client) {
	wac.logger.Info().Str("session_id", client.SessionID).Msg("QR expired")
	wac.clearQRCode(client)
	client.Status = StatusDisconnected
	wac.updateSessionStatus(context.Background(), client)
}

func (wac *WAClient) handleQRError(client *Client, event string) {
	wac.logger.Warn().Str("event", event).Str("session_id", client.SessionID).Msg("QR code error")
	wac.clearQRCode(client)
	client.Status = StatusError
	wac.updateSessionStatus(context.Background(), client)
}

func (wac *WAClient) clearQRCode(client *Client) {
	client.QRCode = ""
	client.QRExpiresAt = time.Time{}
}

// updateSessionStatus updates session status in database
func (wac *WAClient) updateSessionStatus(ctx context.Context, client *Client) {
	deviceJID := ""
	if client.WAClient.Store.ID != nil {
		deviceJID = client.WAClient.Store.ID.String()
	}

	// Debug log to track deviceJID changes
	wac.logger.Debug().
		Str("session_id", client.SessionID).
		Str("session_name", client.Name).
		Str("device_jid", deviceJID).
		Str("status", string(client.Status)).
		Msg("Updating session status")

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
		wac.logger.Error().Err(err).Str("session_id", client.SessionID).Str("device_jid", deviceJID).Msg("Failed to update session status")
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
