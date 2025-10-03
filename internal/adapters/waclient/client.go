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

// MyClient represents a WhatsApp client instance based on wuzapi design
type MyClient struct {
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
	sessionID      uuid.UUID
	sessionName    string
	subscriptions  []string
	db             *sqlx.DB
	gateway        *Gateway
	logger         *logger.Logger
	
	// Connection state
	isConnected    bool
	connectionMux  sync.RWMutex
}

// ClientManager manages multiple WhatsApp clients similar to wuzapi
type ClientManager struct {
	clients    map[uuid.UUID]*MyClient
	httpClients map[uuid.UUID]interface{} // For future HTTP client management
	mutex      sync.RWMutex
	logger     *logger.Logger
}

// Global client manager instance
var (
	clientManager *ClientManager
	once          sync.Once
)

// GetClientManager returns the singleton client manager
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

// SetWhatsmeowClient stores a whatsmeow client for a session
func (cm *ClientManager) SetWhatsmeowClient(sessionID uuid.UUID, client *whatsmeow.Client) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if mycli, exists := cm.clients[sessionID]; exists {
		mycli.WAClient = client
	}
}

// SetMyClient stores a MyClient instance
func (cm *ClientManager) SetMyClient(sessionID uuid.UUID, client *MyClient) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.clients[sessionID] = client
}

// GetMyClient retrieves a MyClient instance
func (cm *ClientManager) GetMyClient(sessionID uuid.UUID) (*MyClient, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	client, exists := cm.clients[sessionID]
	return client, exists
}

// DeleteMyClient removes a MyClient instance
func (cm *ClientManager) DeleteMyClient(sessionID uuid.UUID) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.clients, sessionID)
}

// DeleteWhatsmeowClient removes a whatsmeow client
func (cm *ClientManager) DeleteWhatsmeowClient(sessionID uuid.UUID) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if mycli, exists := cm.clients[sessionID]; exists {
		mycli.WAClient = nil
	}
}

// SetHTTPClient stores an HTTP client (for future use)
func (cm *ClientManager) SetHTTPClient(sessionID uuid.UUID, client interface{}) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.httpClients[sessionID] = client
}

// DeleteHTTPClient removes an HTTP client
func (cm *ClientManager) DeleteHTTPClient(sessionID uuid.UUID) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.httpClients, sessionID)
}

// NewMyClient creates a new MyClient instance
func NewMyClient(sessionID uuid.UUID, sessionName string, client *whatsmeow.Client, db *sqlx.DB, gateway *Gateway, logger *logger.Logger) *MyClient {
	mycli := &MyClient{
		WAClient:    client,
		sessionID:   sessionID,
		sessionName: sessionName,
		db:          db,
		gateway:     gateway,
		logger:      logger,
		subscriptions: []string{"Connected", "Disconnected", "QR", "PairSuccess", "LoggedOut"},
	}

	// Register event handler
	if client != nil {
		mycli.eventHandlerID = client.AddEventHandler(mycli.myEventHandler)
	}

	return mycli
}

// IsConnected returns the connection status
func (mc *MyClient) IsConnected() bool {
	mc.connectionMux.RLock()
	defer mc.connectionMux.RUnlock()
	return mc.isConnected
}

// SetConnected updates the connection status
func (mc *MyClient) SetConnected(connected bool) {
	mc.connectionMux.Lock()
	defer mc.connectionMux.Unlock()
	mc.isConnected = connected
}

// GetSessionID returns the session ID
func (mc *MyClient) GetSessionID() uuid.UUID {
	return mc.sessionID
}

// GetSessionName returns the session name
func (mc *MyClient) GetSessionName() string {
	return mc.sessionName
}

// Disconnect disconnects the WhatsApp client
func (mc *MyClient) Disconnect() {
	if mc.WAClient != nil {
		mc.WAClient.Disconnect()
		mc.SetConnected(false)
	}
}

// Connect connects the WhatsApp client
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

// UpdateConnectionStatus updates the connection status in database
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

// UpdateDeviceJID updates the device JID in database
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

// ClearQRCode clears the QR code from database
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

// UpdateQRCode updates the QR code in database
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

// SetConnectionError sets connection error in database
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

// GetDeviceJID retrieves the device JID from database
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
