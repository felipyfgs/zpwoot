package waclient

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"zpwoot/internal/core/session"
	"zpwoot/platform/logger"
)

type Gateway struct {
	container    *sqlstore.Container
	logger       *logger.Logger
	db           *sqlx.DB
	eventHandler session.EventHandler
	validator    *Validator
	mapper       *Mapper
	qrGenerator  *QRGenerator

	sessions     map[uuid.UUID]*MyClient
	sessionUUIDs map[uuid.UUID]uuid.UUID
	mutex        sync.RWMutex

	killChannels map[uuid.UUID]chan bool
}

func NewGateway(container *sqlstore.Container, logger *logger.Logger) *Gateway {
	return &Gateway{
		container:    container,
		logger:       logger,
		validator:    NewValidator(logger),
		mapper:       NewMapper(logger),
		qrGenerator:  NewQRGenerator(logger).(*QRGenerator),
		sessions:     make(map[uuid.UUID]*MyClient),
		sessionUUIDs: make(map[uuid.UUID]uuid.UUID),
		killChannels: make(map[uuid.UUID]chan bool),
	}
}

func (g *Gateway) SetDatabase(db *sqlx.DB) {
	g.db = db
}

func (g *Gateway) SetEventHandler(handler session.EventHandler) {
	g.eventHandler = handler
}

func (g *Gateway) CreateSession(ctx context.Context, sessionId uuid.UUID) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.logger.InfoWithFields("Creating WhatsApp session", map[string]interface{}{
		"session_id": sessionId.String(),
	})

	if _, exists := g.sessions[sessionId]; exists {
		return fmt.Errorf("session already exists: %s", sessionId.String())
	}

	deviceStore := g.container.NewDevice()
	if deviceStore == nil {
		return fmt.Errorf("failed to create device store for session: %s", sessionId.String())
	}

	waLogger := NewWhatsmeowLogger(g.logger)
	client := whatsmeow.NewClient(deviceStore, waLogger)
	if client == nil {
		return fmt.Errorf("failed to create WhatsApp client for session: %s", sessionId.String())
	}

	// Create MyClient and register it
	myClient := NewMyClient(sessionId, client, g.db, g, g.logger)
	g.sessions[sessionId] = myClient
	g.sessionUUIDs[sessionId] = sessionId
	g.killChannels[sessionId] = make(chan bool, 1)

	g.logger.InfoWithFields("WhatsApp session created successfully", map[string]interface{}{
		"session_id": sessionId.String(),
	})

	return nil
}

func (g *Gateway) ConnectSession(ctx context.Context, sessionId uuid.UUID) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.logger.InfoWithFields("Connecting WhatsApp session", map[string]interface{}{
		"session_id": sessionId.String(),
	})

	go g.startClient(sessionId, sessionId)

	return nil
}

func (g *Gateway) DisconnectSession(ctx context.Context, sessionId uuid.UUID) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.logger.InfoWithFields("Disconnecting WhatsApp session", map[string]interface{}{
		"session_name": sessionId,
	})

	if killChan, exists := g.killChannels[sessionId]; exists {
		select {
		case killChan <- true:
			g.logger.InfoWithFields("Kill signal sent to session", map[string]interface{}{
				"session_name": sessionId,
			})
		default:
			g.logger.WarnWithFields("Kill channel full or closed", map[string]interface{}{
				"session_name": sessionId,
			})
		}
	}

	delete(g.sessions, sessionId)
	delete(g.killChannels, sessionId)

	return nil
}

func (g *Gateway) DeleteSession(ctx context.Context, sessionId uuid.UUID) error {
	g.logger.InfoWithFields("Deleting WhatsApp session", map[string]interface{}{
		"session_name": sessionId,
	})

	err := g.DisconnectSession(ctx, sessionId)
	if err != nil {
		g.logger.WarnWithFields("Error disconnecting session during delete", map[string]interface{}{
			"session_name": sessionId,
			"error":        err.Error(),
		})
	}

	g.mutex.Lock()
	delete(g.sessionUUIDs, sessionId)
	g.mutex.Unlock()

	return nil
}

func (g *Gateway) RestoreSession(ctx context.Context, sessionId uuid.UUID) error {
	g.logger.InfoWithFields("Restoring WhatsApp session", map[string]interface{}{
		"session_name": sessionId,
	})

	// Register session UUID if not exists
	g.mutex.Lock()
	if _, exists := g.sessionUUIDs[sessionId]; !exists {
		g.sessionUUIDs[sessionId] = sessionId
	}
	g.mutex.Unlock()

	var deviceJID string
	query := `SELECT COALESCE("deviceJid", '') FROM "zpSessions" WHERE id = $1`
	err := g.db.Get(&deviceJID, query, sessionId.String())
	if err != nil {
		g.logger.ErrorWithFields("Failed to get device JID from database", map[string]interface{}{
			"session_id":   sessionId.String(),
			"session_name": sessionId,
			"error":        err.Error(),
		})
		return fmt.Errorf("failed to get device JID: %w", err)
	}

	if deviceJID != "" {

		jid, err := g.validator.ParseJID(deviceJID)
		if err != nil {
			g.logger.ErrorWithFields("Invalid device JID in database", map[string]interface{}{
				"session_id":   sessionId.String(),
				"session_name": sessionId,
				"device_jid":   deviceJID,
				"error":        err.Error(),
			})
			return fmt.Errorf("invalid device JID: %w", err)
		}

		_, err = g.container.GetDevice(ctx, jid)
		if err != nil {
			g.logger.ErrorWithFields("Failed to get device store", map[string]interface{}{
				"session_id":   sessionId.String(),
				"session_name": sessionId,
				"device_jid":   deviceJID,
				"error":        err.Error(),
			})
			return fmt.Errorf("failed to get device store: %w", err)
		}

		g.logger.InfoWithFields("Device store validated", map[string]interface{}{
			"session_id":   sessionId.String(),
			"session_name": sessionId,
			"device_jid":   deviceJID,
		})

	}

	return nil
}

func (g *Gateway) RestoreAllSessions(ctx context.Context, sessionIDs []uuid.UUID) error {
	g.logger.InfoWithFields("Restoring all WhatsApp sessions", map[string]interface{}{
		"session_count": len(sessionIDs),
	})

	var errors []error
	for _, sessionID := range sessionIDs {
		err := g.RestoreSession(ctx, sessionID)
		if err != nil {
			g.logger.ErrorWithFields("Failed to restore session", map[string]interface{}{
				"session_id": sessionID.String(),
				"error":      err.Error(),
			})
			errors = append(errors, fmt.Errorf("session %s: %w", sessionID.String(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to restore %d sessions", len(errors))
	}

	return nil
}

func (g *Gateway) RegisterSessionUUID(sessionId, sessionUUID string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	sessionID, err := uuid.Parse(sessionUUID)
	if err != nil {
		g.logger.ErrorWithFields("Invalid session UUID", map[string]interface{}{
			"session_name": sessionId,
			"session_uuid": sessionUUID,
			"error":        err.Error(),
		})
		return
	}

	g.sessionUUIDs[sessionID] = sessionID
	g.logger.InfoWithFields("Session UUID registered for WhatsApp connection", map[string]interface{}{
		"session_name": sessionId,
		"session_id":   sessionID.String(),
	})
}

func (g *Gateway) SessionExists(sessionId uuid.UUID) bool {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	_, exists := g.sessions[sessionId]
	if exists {
		return true
	}

	_, uuidExists := g.sessionUUIDs[sessionId]
	return uuidExists
}

func (g *Gateway) IsSessionConnected(ctx context.Context, sessionId uuid.UUID) (bool, error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	client, exists := g.sessions[sessionId]
	if !exists {

		_, uuidExists := g.sessionUUIDs[sessionId]
		if uuidExists {

			return false, nil
		}
		return false, fmt.Errorf("session not found: %s", sessionId)
	}

	return client.IsConnected(), nil
}

func (g *Gateway) GetSessionInfo(ctx context.Context, sessionId uuid.UUID) (*session.DeviceInfo, error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	client, exists := g.sessions[sessionId]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionId)
	}

	if client.WAClient == nil || client.WAClient.Store.ID == nil {
		return nil, fmt.Errorf("session not initialized: %s", sessionId)
	}

	deviceInfo := g.mapper.MapDeviceInfoFromWhatsmeow(
		*client.WAClient.Store.ID,
		client.WAClient.Store.PushName,
		client.IsConnected(),
	)

	return deviceInfo, nil
}

func (g *Gateway) GenerateQRCode(ctx context.Context, sessionId uuid.UUID) (*session.QRCodeResponse, error) {
	g.logger.InfoWithFields("Generating QR code for session", map[string]interface{}{
		"session_name": sessionId,
	})

	return g.qrGenerator.Generate(ctx, sessionId)
}

func (g *Gateway) Stop(ctx context.Context) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.logger.InfoWithFields("Stopping WhatsApp gateway", map[string]interface{}{
		"active_sessions": len(g.sessions),
	})

	for sessionId, killChan := range g.killChannels {
		select {
		case killChan <- true:
			g.logger.DebugWithFields("Kill signal sent", map[string]interface{}{
				"session_name": sessionId,
			})
		default:
			g.logger.WarnWithFields("Kill channel full or closed", map[string]interface{}{
				"session_name": sessionId,
			})
		}
	}

	g.sessions = make(map[uuid.UUID]*MyClient)
	g.sessionUUIDs = make(map[uuid.UUID]uuid.UUID)
	g.killChannels = make(map[uuid.UUID]chan bool)

	return nil
}

func (g *Gateway) startClient(sessionID uuid.UUID, sessionId uuid.UUID) {
	g.logger.InfoWithFields("Starting WhatsApp client", map[string]interface{}{
		"session_id":   sessionID.String(),
		"session_name": sessionId,
	})

	var deviceStore *store.Device
	var err error

	deviceJID, err := g.getDeviceJIDFromDB(sessionID)
	if err != nil {
		g.logger.ErrorWithFields("Failed to get device JID from database", map[string]interface{}{
			"session_id": sessionID.String(),
			"error":      err.Error(),
		})
		deviceStore = g.container.NewDevice()
	} else if deviceJID != "" {

		jid, err := g.validator.ParseJID(deviceJID)
		if err != nil {
			g.logger.ErrorWithFields("Invalid device JID", map[string]interface{}{
				"session_id": sessionID.String(),
				"device_jid": deviceJID,
				"error":      err.Error(),
			})
			deviceStore = g.container.NewDevice()
		} else {
			deviceStore, err = g.container.GetDevice(context.Background(), jid)
			if err != nil {
				g.logger.ErrorWithFields("Failed to get device store", map[string]interface{}{
					"session_id": sessionID.String(),
					"device_jid": deviceJID,
					"error":      err.Error(),
				})
				deviceStore = g.container.NewDevice()
			}
		}
	} else {
		g.logger.InfoWithFields("No device JID found, creating new device", map[string]interface{}{
			"session_id": sessionID.String(),
		})
		deviceStore = g.container.NewDevice()
	}

	if deviceStore == nil {
		g.logger.ErrorWithFields("Failed to create device store", map[string]interface{}{
			"session_id": sessionID.String(),
		})
		return
	}

	waLogger := NewWhatsmeowLogger(g.logger)
	client := whatsmeow.NewClient(deviceStore, waLogger)
	if client == nil {
		g.logger.ErrorWithFields("Failed to create WhatsApp client", map[string]interface{}{
			"session_id": sessionID.String(),
		})
		return
	}

	myClient := NewMyClient(sessionID, client, g.db, g, g.logger)

	g.mutex.Lock()
	g.sessions[sessionID] = myClient
	g.mutex.Unlock()

	clientManager := GetClientManager(g.logger)
	clientManager.SetMyClient(sessionID, myClient)
	clientManager.SetWhatsmeowClient(sessionID, client)

	if client.Store.ID == nil {

		g.handleQRCodePairing(myClient)
	} else {

		g.logger.InfoWithFields("Device already paired, connecting", map[string]interface{}{
			"session_id": sessionID.String(),
		})
		err = client.Connect()
		if err != nil {
			g.logger.ErrorWithFields("Failed to connect existing session", map[string]interface{}{
				"session_id": sessionID.String(),
				"error":      err.Error(),
			})
			myClient.SetConnectionError(err.Error())
			return
		}
	}

	g.keepClientAlive(sessionID, sessionId, myClient)
}

func (g *Gateway) handleQRCodePairing(myClient *MyClient) {
	qrChan, err := myClient.WAClient.GetQRChannel(context.Background())
	if err != nil {
		if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			g.logger.ErrorWithFields("Failed to get QR channel", map[string]interface{}{
				"session_id": myClient.sessionID.String(),
				"error":      err.Error(),
			})
			return
		}
	} else {

		err = myClient.WAClient.Connect()
		if err != nil {
			g.logger.ErrorWithFields("Failed to connect for QR generation", map[string]interface{}{
				"session_id": myClient.sessionID.String(),
				"error":      err.Error(),
			})
			return
		}

		for evt := range qrChan {
			switch evt.Event {
			case "code":
				g.logger.InfoWithFields("QR code received", map[string]interface{}{
					"session_id": myClient.sessionID.String(),
				})
				err = myClient.handleQRCode(evt.Code)
				if err != nil {
					g.logger.ErrorWithFields("Failed to handle QR code", map[string]interface{}{
						"session_id": myClient.sessionID.String(),
						"error":      err.Error(),
					})
				}

			case "timeout":
				g.logger.WarnWithFields("QR code timeout", map[string]interface{}{
					"session_id": myClient.sessionID.String(),
				})
				myClient.ClearQRCode()

				clientManager := GetClientManager(g.logger)
				clientManager.DeleteMyClient(myClient.sessionID)
				clientManager.DeleteWhatsmeowClient(myClient.sessionID)

				if killChan, exists := g.killChannels[myClient.sessionID]; exists {
					select {
					case killChan <- true:
					default:
					}
				}
				return

			case "success":
				g.logger.InfoWithFields("QR pairing successful", map[string]interface{}{
					"session_id": myClient.sessionID.String(),
				})
				myClient.ClearQRCode()

			default:
				g.logger.DebugWithFields("QR event", map[string]interface{}{
					"session_id": myClient.sessionID.String(),
					"event":      evt.Event,
				})
			}
		}
	}
}

func (g *Gateway) keepClientAlive(sessionID uuid.UUID, sessionId uuid.UUID, myClient *MyClient) {
	killChan, exists := g.killChannels[sessionID]
	if !exists {
		g.logger.ErrorWithFields("No kill channel found for session", map[string]interface{}{
			"session_id":   sessionID.String(),
			"session_name": sessionId,
		})
		return
	}

	g.logger.InfoWithFields("Client keep-alive started", map[string]interface{}{
		"session_id":   sessionID.String(),
		"session_name": sessionId,
	})

	for {
		select {
		case <-killChan:
			g.logger.InfoWithFields("Received kill signal", map[string]interface{}{
				"session_id":   sessionID.String(),
				"session_name": sessionId,
			})

			if myClient.WAClient != nil {
				myClient.WAClient.Disconnect()
			}

			clientManager := GetClientManager(g.logger)
			clientManager.DeleteMyClient(sessionID)
			clientManager.DeleteWhatsmeowClient(sessionID)
			clientManager.DeleteHTTPClient(sessionID)

			myClient.UpdateConnectionStatus(false)
			myClient.ClearQRCode()

			return

		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (g *Gateway) getDeviceJIDFromDB(sessionID uuid.UUID) (string, error) {
	var deviceJID string
	query := `SELECT COALESCE("deviceJid", '') FROM "zpSessions" WHERE id = $1`
	err := g.db.Get(&deviceJID, query, sessionID.String())
	return deviceJID, err
}

func (g *Gateway) SetProxy(ctx context.Context, sessionId uuid.UUID, proxy *session.ProxyConfig) error {
	g.logger.InfoWithFields("Setting proxy for session", map[string]interface{}{
		"session_name": sessionId,
		"proxy_type":   proxy.Type,
		"proxy_host":   proxy.Host,
		"proxy_port":   proxy.Port,
	})

	return fmt.Errorf("proxy configuration not implemented yet")
}

func (g *Gateway) SendTextMessage(ctx context.Context, sessionId uuid.UUID, to, content string) (*session.MessageSendResult, error) {
	g.logger.InfoWithFields("Send text message requested", map[string]interface{}{
		"session_id": sessionId.String(),
		"to":         to,
		"content_len": len(content),
	})

	return nil, fmt.Errorf("text message sending not implemented yet - focus is on connection only")
}

func (g *Gateway) SendMediaMessage(ctx context.Context, sessionId uuid.UUID, to, mediaURL, caption, mediaType string) (*session.MessageSendResult, error) {
	g.logger.InfoWithFields("Send media message requested", map[string]interface{}{
		"session_id":   sessionId.String(),
		"to":           to,
		"media_type":   mediaType,
		"has_caption":  caption != "",
	})

	return nil, fmt.Errorf("media message sending not implemented yet - focus is on connection only")
}

func (g *Gateway) SendLocationMessage(ctx context.Context, sessionId uuid.UUID, to string, latitude, longitude float64, address string) (*session.MessageSendResult, error) {
	g.logger.InfoWithFields("Send location message requested", map[string]interface{}{
		"session_id": sessionId.String(),
		"to":         to,
		"latitude":   latitude,
		"longitude":  longitude,
		"address":    address,
	})

	return nil, fmt.Errorf("location message sending not implemented yet - focus is on connection only")
}

func (g *Gateway) SendContactMessage(ctx context.Context, sessionId uuid.UUID, to, contactName, contactPhone string) (*session.MessageSendResult, error) {
	g.logger.InfoWithFields("Send contact message requested", map[string]interface{}{
		"session_id":    sessionId.String(),
		"to":            to,
		"contact_name":  contactName,
		"contact_phone": contactPhone,
	})

	return nil, fmt.Errorf("contact message sending not implemented yet - focus is on connection only")
}
