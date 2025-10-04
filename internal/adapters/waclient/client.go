package waclient

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.mau.fi/whatsmeow"

	"zpwoot/platform/logger"
)

type MyClient struct {
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
	sessionID      uuid.UUID
	subscriptions  []string
	db             *sqlx.DB
	gateway        *Gateway
	logger         *logger.Logger

	isConnected   bool
	connectionMux sync.RWMutex
}

type ClientManager struct {
	clients     map[uuid.UUID]*MyClient
	httpClients map[uuid.UUID]interface{}
	mutex       sync.RWMutex
	logger      *logger.Logger
}

var (
	clientManager *ClientManager
	once          sync.Once
)

func GetClientManager(logger *logger.Logger) *ClientManager {
	once.Do(func() {
		clientManager = &ClientManager{
			clients:     make(map[uuid.UUID]*MyClient),
			httpClients: make(map[uuid.UUID]interface{}),
			logger:      logger,
		}
	})
	return clientManager
}

func (cm *ClientManager) SetWhatsmeowClient(sessionID uuid.UUID, client *whatsmeow.Client) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if mycli, exists := cm.clients[sessionID]; exists {
		mycli.WAClient = client
	}
}

func (cm *ClientManager) SetMyClient(sessionID uuid.UUID, client *MyClient) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.clients[sessionID] = client
}

func (cm *ClientManager) GetMyClient(sessionID uuid.UUID) (*MyClient, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	client, exists := cm.clients[sessionID]
	return client, exists
}

func (cm *ClientManager) DeleteMyClient(sessionID uuid.UUID) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.clients, sessionID)
}

func (cm *ClientManager) DeleteWhatsmeowClient(sessionID uuid.UUID) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if mycli, exists := cm.clients[sessionID]; exists {
		mycli.WAClient = nil
	}
}

func (cm *ClientManager) SetHTTPClient(sessionID uuid.UUID, client interface{}) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.httpClients[sessionID] = client
}

func (cm *ClientManager) DeleteHTTPClient(sessionID uuid.UUID) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.httpClients, sessionID)
}

func NewMyClient(sessionID uuid.UUID, client *whatsmeow.Client, db *sqlx.DB, gateway *Gateway, logger *logger.Logger) *MyClient {
	mycli := &MyClient{
		WAClient:      client,
		sessionID:     sessionID,
		db:            db,
		gateway:       gateway,
		logger:        logger,
		subscriptions: []string{"Connected", "Disconnected", "QR", "PairSuccess", "LoggedOut"},
	}

	if client != nil {
		mycli.eventHandlerID = client.AddEventHandler(mycli.myEventHandler)
	}

	return mycli
}

func (mc *MyClient) IsConnected() bool {
	mc.connectionMux.RLock()
	defer mc.connectionMux.RUnlock()
	return mc.isConnected
}

func (mc *MyClient) SetConnected(connected bool) {
	mc.connectionMux.Lock()
	defer mc.connectionMux.Unlock()
	mc.isConnected = connected
}

func (mc *MyClient) GetSessionID() uuid.UUID {
	return mc.sessionID
}

func (mc *MyClient) GetSessionName() string {
	return mc.sessionName
}

func (mc *MyClient) Disconnect() {
	if mc.WAClient != nil {
		mc.WAClient.Disconnect()
		mc.SetConnected(false)
	}
}

func (mc *MyClient) Connect() error {
	if mc.WAClient == nil {
		return fmt.Errorf("WhatsApp client not initialized")
	}

	err := mc.WAClient.Connect()
	if err != nil {
		mc.logger.ErrorWithFields("Failed to connect WhatsApp client", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
		return err
	}

	return nil
}

func (mc *MyClient) UpdateConnectionStatus(connected bool) error {
	query := `UPDATE "zpSessions" SET "isConnected" = $1, "updatedAt" = NOW() WHERE id = $2`
	_, err := mc.db.Exec(query, connected, mc.sessionID.String())
	if err != nil {
		mc.logger.ErrorWithFields("Failed to update connection status in database", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"connected":  connected,
			"error":      err.Error(),
		})
		return err
	}

	mc.SetConnected(connected)
	return nil
}

func (mc *MyClient) UpdateDeviceJID(deviceJID string) error {
	query := `UPDATE "zpSessions" SET "deviceJid" = $1, "isConnected" = true, "connectedAt" = NOW(), "updatedAt" = NOW() WHERE id = $2`
	_, err := mc.db.Exec(query, deviceJID, mc.sessionID.String())
	if err != nil {
		mc.logger.ErrorWithFields("Failed to update device JID in database", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"device_jid": deviceJID,
			"error":      err.Error(),
		})
		return err
	}

	mc.SetConnected(true)
	return nil
}

func (mc *MyClient) ClearQRCode() error {
	query := `UPDATE "zpSessions" SET "qrCode" = NULL, "qrCodeExpiresAt" = NULL, "updatedAt" = NOW() WHERE id = $1`
	_, err := mc.db.Exec(query, mc.sessionID.String())
	if err != nil {
		mc.logger.ErrorWithFields("Failed to clear QR code in database", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
		return err
	}

	return nil
}

func (mc *MyClient) UpdateQRCode(qrCode string, expiresAt time.Time) error {
	query := `UPDATE "zpSessions" SET "qrCode" = $1, "qrCodeExpiresAt" = $2, "updatedAt" = NOW() WHERE id = $3`
	_, err := mc.db.Exec(query, qrCode, expiresAt, mc.sessionID.String())
	if err != nil {
		mc.logger.ErrorWithFields("Failed to update QR code in database", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
		return err
	}

	return nil
}

func (mc *MyClient) SetConnectionError(errorMsg string) error {
	query := `UPDATE "zpSessions" SET "connectionError" = $1, "isConnected" = false, "updatedAt" = NOW() WHERE id = $2`
	_, err := mc.db.Exec(query, errorMsg, mc.sessionID.String())
	if err != nil {
		mc.logger.ErrorWithFields("Failed to set connection error in database", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
		return err
	}

	mc.SetConnected(false)
	return nil
}

func (mc *MyClient) GetDeviceJID() (string, error) {
	var deviceJID sql.NullString
	query := `SELECT "deviceJid" FROM "zpSessions" WHERE id = $1`
	err := mc.db.Get(&deviceJID, query, mc.sessionID.String())
	if err != nil {
		return "", err
	}

	if deviceJID.Valid {
		return deviceJID.String, nil
	}
	return "", nil
}
