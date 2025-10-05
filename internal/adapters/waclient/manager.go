package waclient

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/core/domain/session"

	"github.com/jmoiron/sqlx"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Constants for QR code management and session handling
const (
	QRExpirationBuffer  = 2 * time.Minute
	QRFirstTimeout      = 60 * time.Second
	QRSubsequentTimeout = 20 * time.Second
	AutoReconnectDelay  = 2 * time.Second
)

// WAClient manages WhatsApp client sessions and their lifecycle
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
	GetByID(ctx context.Context, sessionID string) (*session.Session, error)
	GetByName(ctx context.Context, name string) (*session.Session, error)
	Create(ctx context.Context, sess *session.Session) error
	Update(ctx context.Context, sess *session.Session) error
	Delete(ctx context.Context, sessionID string) error
	List(ctx context.Context, limit, offset int) ([]*session.Session, error)
}

// NewWAClient creates a new WhatsApp client manager
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
	sessions, err := wac.sessionRepo.List(ctx, 1000, 0)
	if err != nil {
		wac.logger.Error().Err(err).Msg("Failed to load sessions from database")
		return
	}

	wac.logger.Debug().Int("count", len(sessions)).Msg("Loading sessions from database")

	for _, sess := range sessions {
		deviceStore := wac.getOrCreateDeviceStore(ctx, sess.DeviceJID)
		client := wac.createClient(ctx, sess, deviceStore)

		wac.sessionsMutex.Lock()
		wac.sessions[sess.ID] = client
		wac.sessionsMutex.Unlock()

		wac.logger.Debug().Str("session_id", sess.ID).Str("name", sess.Name).Msg("Loaded session")

		if sess.IsConnected && sess.DeviceJID != "" {
			wac.logger.Debug().Str("session_id", sess.ID).Msg("Auto-reconnecting")
			go wac.autoReconnect(client)
		}
	}
}

// autoReconnect attempts to reconnect a session after a delay
func (wac *WAClient) autoReconnect(client *Client) {
	timer := time.NewTimer(AutoReconnectDelay)
	defer timer.Stop()

	select {
	case <-timer.C:
		// Proceed with reconnection
	case <-client.ctx.Done():
		wac.logger.Debug().Str("session_id", client.SessionID).Msg("Auto-reconnect cancelled")
		return
	}

	if err := client.WAClient.Connect(); err != nil {
		wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Failed to auto-reconnect session")
		client.Status = session.StatusDisconnected
	} else {
		wac.logger.Debug().Str("session_id", client.SessionID).Msg("Auto-reconnected")
		client.Status = session.StatusConnected
		client.ConnectedAt = time.Now()
		client.LastSeen = time.Now()
	}
	wac.updateSessionStatus(client.ctx, client)
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

func (wac *WAClient) createClient(ctx context.Context, sess *session.Session, deviceStore *store.Device) *Client {
	waClient := whatsmeow.NewClient(deviceStore, waLog.Noop)
	clientCtx, cancel := context.WithCancel(ctx)

	client := &Client{
		SessionID:   sess.ID,
		Name:        sess.Name,
		WAClient:    waClient,
		Status:      sess.GetStatus(),
		QRCode:      sess.QRCode,
		QRExpiresAt: getTimeValue(sess.QRCodeExpiresAt),
		ConnectedAt: getTimeValue(sess.ConnectedAt),
		LastSeen:    getTimeValue(sess.LastSeen),
		Config: &SessionConfig{
			SessionID: sess.ID,
			Name:      sess.Name,
		},
		ctx:    clientCtx,
		cancel: cancel,
	}

	client.EventHandler = waClient.AddEventHandler(wac.createEventHandler(client))
	return client
}

func getTimeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
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

	sess, err := wac.sessionRepo.GetByID(ctx, config.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session from database: %w", err)
	}

	deviceStore := wac.container.NewDevice()
	wac.logger.Debug().
		Str("session_id", config.SessionID).
		Str("session_name", sess.Name).
		Str("device_id", fmt.Sprintf("%p", deviceStore)).
		Msg("Created new device for session")

	client := wac.createClient(ctx, sess, deviceStore)
	client.Config = config
	client.Events = config.Events
	client.WebhookURL = config.WebhookURL

	wac.sessions[config.SessionID] = client

	wac.logger.Info().Str("session_id", config.SessionID).Msg("WhatsApp session initialized")
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
		sess, dbErr := wac.sessionRepo.GetByID(ctx, sessionID)
		if dbErr != nil {
			return fmt.Errorf("session not found: %w", dbErr)
		}
		client, err = wac.recreateClient(ctx, sess)
		if err != nil {
			return fmt.Errorf("failed to recreate client: %w", err)
		}
	} else if err != nil {
		return err
	}

	if client.Status == session.StatusConnected {
		return nil
	}

	client.Status = session.StatusConnecting
	wac.updateSessionStatus(ctx, client)

	if client.WAClient.Store.ID == nil {

		if err = client.WAClient.Connect(); err != nil {
			wac.logger.Error().Err(err).Msg("Failed to connect client")
			client.Status = session.StatusError
			wac.updateSessionStatus(ctx, client)
			return fmt.Errorf("failed to connect: %w", err)
		}

		wac.logger.Info().Str("session_id", client.SessionID).Msg("Connected - waiting for QR codes via events")
	} else {
		if err = client.WAClient.Connect(); err != nil {
			client.Status = session.StatusError
			wac.updateSessionStatus(ctx, client)
			return fmt.Errorf("failed to reconnect: %w", err)
		}
	}

	return nil
}

func (wac *WAClient) recreateClient(ctx context.Context, sess *session.Session) (*Client, error) {
	wac.sessionsMutex.Lock()
	defer wac.sessionsMutex.Unlock()

	if existing, exists := wac.sessions[sess.ID]; exists {
		return existing, nil
	}

	deviceStore := wac.container.NewDevice()
	wac.logger.Debug().
		Str("session_id", sess.ID).
		Str("session_name", sess.Name).
		Msg("Recreating client with new device")

	client := wac.createClient(ctx, sess, deviceStore)
	wac.sessions[sess.ID] = client

	return client, nil
}

func (wac *WAClient) DisconnectSession(ctx context.Context, sessionID string) error {
	client, err := wac.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if client.Status == session.StatusDisconnected {
		return nil
	}

	client.WAClient.Disconnect()
	client.Status = session.StatusDisconnected
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

	client.Status = session.StatusDisconnected
	client.cancel()
	client.WAClient.RemoveEventHandler(client.EventHandler)
	delete(wac.sessions, sessionID)

	now := time.Now()
	sess := &session.Session{
		ID:          client.SessionID,
		Name:        client.Name,
		DeviceJID:   "",
		IsConnected: false,
		QRCode:      "",
		UpdatedAt:   now,
		LastSeen:    &now,
	}

	if err := wac.sessionRepo.Update(ctx, sess); err != nil {
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

	if client.Status != session.StatusDisconnected {
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
			wac.handleQREvent(client, v)
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
	client.Status = session.StatusConnected
	client.ConnectedAt = time.Now()
	client.LastSeen = time.Now()

	if client.WAClient.Store.ID != nil {
		wac.logger.Info().Str("session_id", client.SessionID).Msg("Connected")
	}

	wac.updateSessionStatus(client.ctx, client)
	wac.sendWebhook(client, EventConnected, evt)
}

func (wac *WAClient) handleDisconnected(client *Client, evt *events.Disconnected) {
	client.Status = session.StatusDisconnected
	client.LastSeen = time.Now()

	wac.logger.Debug().Str("session_id", client.SessionID).Msg("Disconnected")
	wac.updateSessionStatus(client.ctx, client)
	wac.sendWebhook(client, EventDisconnected, evt)
}

func (wac *WAClient) handleLoggedOut(client *Client, evt *events.LoggedOut) {
	client.Status = session.StatusDisconnected
	client.LastSeen = time.Now()

	wac.logger.Info().Str("session_id", client.SessionID).Msg("Logged out")
	wac.updateSessionStatus(client.ctx, client)
	wac.sendWebhook(client, EventLoggedOut, evt)
}

func (wac *WAClient) handleQREvent(client *Client, evt *events.QR) {
	wac.logger.Info().
		Str("session_id", client.SessionID).
		Int("total_codes", len(evt.Codes)).
		Msg("Received QR codes from WhatsApp via events")

	if len(evt.Codes) == 0 {
		wac.logger.Error().
			Str("session_id", client.SessionID).
			Msg("No QR codes received")
		return
	}

	client.Status = session.StatusQRCode
	wac.updateSessionStatus(client.ctx, client)

	go wac.processAllQRCodesFromEvent(context.Background(), client, evt.Codes)
}

func (wac *WAClient) processAllQRCodesFromEvent(ctx context.Context, client *Client, codes []string) {
	for i, code := range codes {

		if client.Status != session.StatusQRCode {
			wac.logger.Info().
				Str("session_id", client.SessionID).
				Int("qr_number", i+1).
				Msg("QR processing stopped - client status changed")
			return
		}

		wac.logger.Info().
			Str("session_id", client.SessionID).
			Int("qr_number", i+1).
			Int("total_codes", len(codes)).
			Msg("Displaying QR code")

		qrterminal.GenerateHalfBlock(code, qrterminal.L, os.Stdout)

		var timeout time.Duration
		if i == 0 {
			timeout = 60 * time.Second
		} else {
			timeout = 20 * time.Second
		}

		fmt.Printf("\n📱 Scan QR code | Session: %s | QR: %d/%d | Timeout: %d seconds\n\n",
			client.SessionID, i+1, len(codes), int(timeout.Seconds()))

		client.QRCode = code
		client.QRExpiresAt = time.Now().Add(timeout)
		wac.updateSessionStatus(ctx, client)

		qrEvent := &QREvent{
			Event:     "code",
			Code:      code,
			ExpiresAt: client.QRExpiresAt,
		}
		wac.sendWebhook(client, EventQR, qrEvent)

		select {
		case <-time.After(timeout):
			wac.logger.Info().
				Str("session_id", client.SessionID).
				Int("qr_number", i+1).
				Msg("QR code expired - showing next QR code")

			wac.clearQRCode(client)
			wac.updateSessionStatus(ctx, client)

		case <-ctx.Done():
			wac.logger.Info().
				Str("session_id", client.SessionID).
				Msg("QR processing cancelled")
			return
		}
	}

	wac.logger.Info().
		Str("session_id", client.SessionID).
		Int("total_codes", len(codes)).
		Msg("All QR codes expired - disconnecting")

	wac.clearQRCode(client)
	client.Status = session.StatusDisconnected
	wac.updateSessionStatus(ctx, client)
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

		go wac.webhookSender.SendWebhook(client.ctx, webhookEvent)
	}
}

func (wac *WAClient) clearQRCode(client *Client) {
	client.QRCode = ""
	client.QRExpiresAt = time.Time{}
}

func (wac *WAClient) updateSessionStatus(ctx context.Context, client *Client) {
	deviceJID := ""
	if client.WAClient.Store.ID != nil {
		deviceJID = client.WAClient.Store.ID.String()
	}

	wac.logger.Debug().
		Str("session_id", client.SessionID).
		Str("session_name", client.Name).
		Str("device_jid", deviceJID).
		Str("status", string(client.Status)).
		Str("qr_code", client.QRCode).
		Time("qr_expires_at", client.QRExpiresAt).
		Msg("Updating session status")

	now := time.Now()
	sess := &session.Session{
		ID:          client.SessionID,
		Name:        client.Name,
		DeviceJID:   deviceJID,
		IsConnected: client.Status == session.StatusConnected,
		QRCode:      client.QRCode,
		UpdatedAt:   now,
	}

	if !client.QRExpiresAt.IsZero() {
		sess.QRCodeExpiresAt = &client.QRExpiresAt
	}
	if !client.ConnectedAt.IsZero() {
		sess.ConnectedAt = &client.ConnectedAt
	}
	if !client.LastSeen.IsZero() {
		sess.LastSeen = &client.LastSeen
	}

	dbCtx := context.Background()
	if err := wac.sessionRepo.Update(dbCtx, sess); err != nil {
		wac.logger.Error().Err(err).Str("session_id", client.SessionID).Str("device_jid", deviceJID).Msg("Failed to update session status")
	}
}

func NewWAStoreContainer(db *sqlx.DB, logger *logger.Logger) *sqlstore.Container {

	dbURL := "postgres://zpwoot:zpwoot123@localhost:5432/zpwoot?sslmode=disable"

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
