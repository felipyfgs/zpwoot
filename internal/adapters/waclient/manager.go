package waclient

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/domain/webhook"
	"zpwoot/internal/core/ports/output"

	"github.com/jmoiron/sqlx"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCompanionReg"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

const (
	AutoReconnectDelay = 2 * time.Second
)

type WAClient struct {
	sessions      map[string]*Client
	sessionsMutex sync.RWMutex
	container     *sqlstore.Container
	logger        *logger.Logger
	eventHandler  EventHandler
	webhookSender WebhookSender
	sessionRepo   SessionRepository
}

type SessionRepository interface {
	GetByID(ctx context.Context, sessionID string) (*session.Session, error)
	GetByName(ctx context.Context, name string) (*session.Session, error)
	Create(ctx context.Context, sess *session.Session) error
	Update(ctx context.Context, sess *session.Session) error
	Delete(ctx context.Context, sessionID string) error
	List(ctx context.Context, limit, offset int) ([]*session.Session, error)
}

func NewWAClient(container *sqlstore.Container, logger *logger.Logger, sessionRepo SessionRepository, webhookSender output.WebhookSender, webhookRepo webhook.Repository) *WAClient {
	store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_UNKNOWN.Enum()
	store.DeviceProps.Os = proto.String(runtime.GOOS)

	wac := &WAClient{
		sessions:    make(map[string]*Client),
		container:   container,
		logger:      logger,
		sessionRepo: sessionRepo,
	}

	if webhookSender != nil && webhookRepo != nil {
		wac.eventHandler = NewDefaultEventHandler(logger, webhookSender, webhookRepo)
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

	for _, sess := range sessions {
		if sess.DeviceJID == "" {
			wac.logger.Debug().Str("session_id", sess.ID).Msg("Skipping session without device JID")
			continue
		}

		jid, parseErr := types.ParseJID(sess.DeviceJID)
		if parseErr != nil {
			wac.logger.Error().Err(parseErr).Str("jid", sess.DeviceJID).Str("session_id", sess.ID).Msg("Failed to parse JID, skipping session")
			continue
		}

		deviceStore, err := wac.container.GetDevice(ctx, jid)
		if err != nil {
			wac.logger.Error().Err(err).Str("jid", sess.DeviceJID).Str("session_id", sess.ID).Msg("Failed to get device store, skipping session")
			continue
		}

		client := wac.createClient(ctx, sess, deviceStore)

		wac.sessionsMutex.Lock()
		wac.sessions[sess.ID] = client
		wac.sessionsMutex.Unlock()

		if sess.IsConnected {
			go wac.autoReconnect(client)
		}
	}
}

func (wac *WAClient) autoReconnect(client *Client) {
	timer := time.NewTimer(AutoReconnectDelay)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-client.ctx.Done():
		wac.logger.Debug().Str("session_id", client.SessionID).Msg("Auto-reconnect canceled")
		return
	}

	if err := client.WAClient.Connect(); err != nil {
		wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Failed to auto-reconnect session")
		client.Status = session.StatusDisconnected
	} else {
		wac.logger.Debug().Str("session_id", client.SessionID).Msg("Auto-reconnected")
	}

	wac.updateSessionStatus(client.ctx, client)
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
	client := wac.createClient(ctx, sess, deviceStore)
	client.Config = config
	client.Events = config.Events
	client.WebhookURL = config.WebhookURL

	wac.sessions[config.SessionID] = client

	wac.logger.Info().Str("session_id", config.SessionID).Msg("WhatsApp session initialized with new device")
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
	sess, err := wac.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	client, err := wac.getOrRecreateClient(ctx, sess)
	if err != nil {
		return err
	}

	wac.logger.Info().Str("session_id", sessionID).Msg("Connecting session")

	if client.WAClient.IsConnected() {
		wac.logger.Info().Str("session_id", sessionID).Msg("Session already connected")
		return nil
	}

	client.Status = session.StatusConnecting
	wac.updateSessionStatus(ctx, client)

	if client.WAClient.Store.ID == nil {
		return wac.connectNewSession(ctx, client)
	}

	return wac.reconnectExistingSession(ctx, client)
}

func (wac *WAClient) getOrRecreateClient(ctx context.Context, sess *session.Session) (*Client, error) {
	client, err := wac.GetSession(ctx, sess.ID)
	if errors.Is(err, ErrSessionNotFound) {
		var deviceStore *store.Device

		if sess.DeviceJID != "" {
			jid, parseErr := types.ParseJID(sess.DeviceJID)
			if parseErr != nil {
				wac.logger.Error().Err(parseErr).Str("jid", sess.DeviceJID).Str("session_id", sess.ID).Msg("Failed to parse JID, creating new device")
				deviceStore = wac.container.NewDevice()
			} else {
				deviceStore, err = wac.container.GetDevice(ctx, jid)
				if err != nil {
					wac.logger.Error().Err(err).Str("jid", sess.DeviceJID).Str("session_id", sess.ID).Msg("Failed to get device, creating new one")
					deviceStore = wac.container.NewDevice()
				} else {
					wac.logger.Info().Str("jid", sess.DeviceJID).Str("session_id", sess.ID).Msg("Loaded existing device store")
				}
			}
		} else {
			wac.logger.Info().Str("session_id", sess.ID).Msg("No device JID, creating new device for QR code")
			deviceStore = wac.container.NewDevice()
		}

		if deviceStore == nil {
			wac.logger.Warn().Str("session_id", sess.ID).Msg("No store found, creating new one")
			deviceStore = wac.container.NewDevice()
		}

		client = wac.createClient(ctx, sess, deviceStore)

		wac.sessionsMutex.Lock()
		wac.sessions[sess.ID] = client
		wac.sessionsMutex.Unlock()
	} else if err != nil {
		return nil, err
	}
	return client, nil
}

func (wac *WAClient) connectNewSession(ctx context.Context, client *Client) error {
	qrChan, err := client.WAClient.GetQRChannel(context.Background())
	if err != nil && !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
		wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Failed to get QR channel")
		client.Status = session.StatusError
		wac.updateSessionStatus(ctx, client)
		return fmt.Errorf("failed to get QR channel: %w", err)
	}

	if err := client.WAClient.Connect(); err != nil {
		wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Failed to connect")
		client.Status = session.StatusError
		wac.updateSessionStatus(ctx, client)
		return fmt.Errorf("failed to connect: %w", err)
	}

	wac.logger.Info().Str("session_id", client.SessionID).Msg("Waiting for QR code")
	go wac.processQRCodes(context.Background(), client, qrChan)
	return nil
}

func (wac *WAClient) reconnectExistingSession(ctx context.Context, client *Client) error {
	if err := client.WAClient.Connect(); err != nil {
		if strings.Contains(err.Error(), "websocket is already connected") {
			client.Status = session.StatusConnected
			wac.updateSessionStatus(ctx, client)
			return nil
		}

		client.Status = session.StatusError
		wac.updateSessionStatus(ctx, client)
		return fmt.Errorf("failed to reconnect: %w", err)
	}
	return nil
}

func (wac *WAClient) processQRCodes(ctx context.Context, client *Client, qrChan <-chan whatsmeow.QRChannelItem) {
	for evt := range qrChan {
		switch evt.Event {
		case "code":
			wac.logger.Info().Str("session_id", client.SessionID).Msg("QR code generated")
			wac.displayQRCode(evt.Code, client.SessionID)
			wac.updateClientWithQRCode(ctx, client, evt.Code)
			wac.sendQRWebhook(client, evt.Code)
		case "timeout":
			wac.logger.Warn().Str("session_id", client.SessionID).Msg("QR code expired")
			wac.clearQRCode(client)
			wac.updateSessionStatus(ctx, client)
		case "success":
			wac.logger.Info().Str("session_id", client.SessionID).Msg("QR code scanned successfully")
		}
	}
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
		case *events.PairSuccess:
			wac.handlePairSuccess(client, v)
		case *events.QR:
			wac.handleQREvent(client, v)
		case *events.Message:
			wac.handleMessage(client, v)
		case *events.Receipt:
			wac.handleReceipt(client, v)
		case *events.Presence:
			wac.handlePresence(client, v)
		case *events.ChatPresence:
			wac.handleChatPresence(client, v)
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

	go func(ctx context.Context) {
		wac.updateSessionStatus(ctx, client)
	}(client.ctx)
	wac.sendWebhook(client, EventConnected, evt)
}

func (wac *WAClient) handleDisconnected(client *Client, evt *events.Disconnected) {
	client.Status = session.StatusDisconnected
	client.LastSeen = time.Now()

	wac.logger.Warn().Str("session_id", client.SessionID).Msg("Disconnected from WhatsApp")

	go func() {
		wac.updateSessionStatus(context.Background(), client)
	}()
	wac.sendWebhook(client, EventDisconnected, evt)
}

func (wac *WAClient) handleLoggedOut(client *Client, evt *events.LoggedOut) {
	client.Status = session.StatusDisconnected
	client.LastSeen = time.Now()

	wac.logger.Info().Str("session_id", client.SessionID).Msg("Logged out")

	go func() {
		wac.updateSessionStatus(context.Background(), client)
	}()
	wac.sendWebhook(client, EventLoggedOut, evt)
}

func (wac *WAClient) handlePairSuccess(client *Client, evt *events.PairSuccess) {
	wac.logger.Info().Str("session_id", client.SessionID).Str("device_jid", evt.ID.String()).Msg("Device paired successfully")

	sess, err := wac.sessionRepo.GetByID(context.Background(), client.SessionID)
	if err != nil {
		wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Failed to get session")
		return
	}

	sess.DeviceJID = evt.ID.String()

	if err := wac.sessionRepo.Update(context.Background(), sess); err != nil {
		wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Failed to update session")
		return
	}

	wac.sendWebhook(client, EventPairSuccess, map[string]interface{}{
		"deviceJID":    evt.ID.String(),
		"businessName": evt.BusinessName,
		"platform":     evt.Platform,
	})
}

func (wac *WAClient) handleQREvent(client *Client, evt *events.QR) {
	wac.logger.Debug().Str("session_id", client.SessionID).Int("codes", len(evt.Codes)).Msg("QR event received (ignored, using channel)")
}

func (wac *WAClient) displayQRCode(code, sessionID string) {
	qrterminal.GenerateHalfBlock(code, qrterminal.L, os.Stdout)
}

func (wac *WAClient) updateClientWithQRCode(ctx context.Context, client *Client, code string) {
	client.QRCode = code
	client.QRExpiresAt = time.Now().Add(60 * time.Second)
	client.Status = session.StatusQRCode
	wac.updateSessionStatus(ctx, client)
}

func (wac *WAClient) sendQRWebhook(client *Client, code string) {
	qrEvent := &QREvent{
		Event:     "code",
		Code:      code,
		ExpiresAt: client.QRExpiresAt,
	}
	wac.sendWebhook(client, EventQR, qrEvent)
}

func (wac *WAClient) handleMessage(client *Client, evt *events.Message) {
	client.LastSeen = time.Now()

	wac.sendWebhook(client, EventMessage, evt)

	if wac.eventHandler != nil {
		if err := wac.eventHandler.HandleEvent(client, evt); err != nil {
			wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Event handler error for message")
		}
	}
}

func (wac *WAClient) handleReceipt(client *Client, evt *events.Receipt) {
	client.LastSeen = time.Now()

	wac.sendWebhook(client, EventReceipt, evt)

	if wac.eventHandler != nil {
		if err := wac.eventHandler.HandleEvent(client, evt); err != nil {
			wac.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Event handler error for receipt")
		}
	}
}

func (wac *WAClient) handlePresence(client *Client, evt *events.Presence) {
	client.LastSeen = time.Now()

	status := "online"
	if evt.Unavailable {
		status = "offline"
	}

	wac.logger.Info().
		Str("session_id", client.SessionID).
		Str("from", evt.From.String()).
		Str("status", status).
		Msg("Presence update")

	wac.sendWebhook(client, EventPresence, evt)
}

func (wac *WAClient) handleChatPresence(client *Client, evt *events.ChatPresence) {
	client.LastSeen = time.Now()

	wac.logger.Info().
		Str("session_id", client.SessionID).
		Str("chat", evt.MessageSource.Chat.String()).
		Str("sender", evt.MessageSource.Sender.String()).
		Str("state", string(evt.State)).
		Str("media", string(evt.Media)).
		Msg("Chat presence update")

	wac.sendWebhook(client, EventChatPresence, evt)
}

func (wac *WAClient) sendWebhook(client *Client, eventType EventType, event interface{}) {
	if wac.webhookSender == nil || client.WebhookURL == "" {
		return
	}

	webhookEvent := &WebhookEvent{
		Type:      eventType,
		SessionID: client.SessionID,
		Event:     event,
		Timestamp: time.Now(),
	}

	go func() {
		if err := wac.webhookSender.SendWebhook(client.ctx, webhookEvent); err != nil {
			wac.logger.Error().Err(err).Msg("Failed to send webhook")
		}
	}()
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

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := wac.sessionRepo.Update(dbCtx, sess); err != nil {
		if !errors.Is(err, context.Canceled) {
			wac.logger.Error().Err(err).Str("session_id", client.SessionID).Str("device_jid", deviceJID).Msg("Failed to update session status")
		}
	}
}

func NewWAStoreContainer(db *sqlx.DB, logger *logger.Logger, dbURL string) *sqlstore.Container {
	container, err := sqlstore.New(context.Background(), "postgres", dbURL, waLog.Noop)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to create WhatsApp store container")

		return nil
	}

	return container
}
