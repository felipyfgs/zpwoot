package waclient

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	"zpwoot/internal/core/session"
)

// myEventHandler handles WhatsApp events, adapted from wuzapi
func (mc *MyClient) myEventHandler(rawEvt interface{}) {
	sessionID := mc.sessionID.String()
	
	mc.logger.DebugWithFields("Received WhatsApp event", map[string]interface{}{
		"session_id": sessionID,
		"event_type": fmt.Sprintf("%T", rawEvt),
	})

	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		mc.handleAppStateSyncComplete(evt)
		
	case *events.Connected:
		mc.handleConnected(evt)
		
	case *events.PushNameSetting:
		mc.handlePushNameSetting(evt)
		
	case *events.PairSuccess:
		mc.handlePairSuccess(evt)
		
	case *events.StreamReplaced:
		mc.handleStreamReplaced(evt)
		
	case *events.LoggedOut:
		mc.handleLoggedOut(evt)
		
	case *events.Disconnected:
		mc.handleDisconnected(evt)
		
	case *events.ConnectFailure:
		mc.handleConnectFailure(evt)
		
	default:
		mc.logger.DebugWithFields("Unhandled event type", map[string]interface{}{
			"session_id": sessionID,
			"event_type": fmt.Sprintf("%T", rawEvt),
		})
	}
}

// handleAppStateSyncComplete handles app state sync completion
func (mc *MyClient) handleAppStateSyncComplete(evt *events.AppStateSyncComplete) {
	if len(mc.WAClient.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
		err := mc.WAClient.SendPresence(types.PresenceAvailable)
		if err != nil {
			mc.logger.WarnWithFields("Failed to send available presence", map[string]interface{}{
				"session_id": mc.sessionID.String(),
				"error":      err.Error(),
			})
		} else {
			mc.logger.InfoWithFields("Marked self as available", map[string]interface{}{
				"session_id": mc.sessionID.String(),
			})
		}
	}
}

// handleConnected handles connection events
func (mc *MyClient) handleConnected(evt *events.Connected) {
	mc.logger.InfoWithFields("WhatsApp connected", map[string]interface{}{
		"session_id": mc.sessionID.String(),
	})

	// Update connection status in database
	err := mc.UpdateConnectionStatus(true)
	if err != nil {
		mc.logger.ErrorWithFields("Failed to update connection status", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
	}

	// Clear QR code since we're connected
	err = mc.ClearQRCode()
	if err != nil {
		mc.logger.ErrorWithFields("Failed to clear QR code", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
	}

	// Send presence available
	if len(mc.WAClient.Store.PushName) > 0 {
		err = mc.WAClient.SendPresence(types.PresenceAvailable)
		if err != nil {
			mc.logger.WarnWithFields("Failed to send available presence", map[string]interface{}{
				"session_id": mc.sessionID.String(),
				"error":      err.Error(),
			})
		}
	}

	// Notify gateway about connection
	if mc.gateway != nil && mc.gateway.eventHandler != nil {
		deviceInfo := &session.DeviceInfo{
			Platform:    "android",
			DeviceModel: "zpwoot-client",
			OSVersion:   "11",
			AppVersion:  "2.23.24.76",
		}
		mc.gateway.eventHandler.OnSessionConnected(mc.sessionName, deviceInfo)
	}
}

// handlePushNameSetting handles push name setting events
func (mc *MyClient) handlePushNameSetting(evt *events.PushNameSetting) {
	pushName := ""
	if evt.Action != nil && evt.Action.Name != nil {
		pushName = *evt.Action.Name
	}

	mc.logger.InfoWithFields("Push name setting received", map[string]interface{}{
		"session_id": mc.sessionID.String(),
		"push_name":  pushName,
	})

	// Send presence available when pushname is set
	if len(pushName) > 0 {
		err := mc.WAClient.SendPresence(types.PresenceAvailable)
		if err != nil {
			mc.logger.WarnWithFields("Failed to send available presence", map[string]interface{}{
				"session_id": mc.sessionID.String(),
				"error":      err.Error(),
			})
		}
	}
}

// handlePairSuccess handles successful pairing
func (mc *MyClient) handlePairSuccess(evt *events.PairSuccess) {
	mc.logger.InfoWithFields("QR Pair Success", map[string]interface{}{
		"session_id":    mc.sessionID.String(),
		"session_name":  mc.sessionName,
		"jid":           evt.ID.String(),
		"business_name": evt.BusinessName,
		"platform":      evt.Platform,
	})

	// Update device JID in database
	err := mc.UpdateDeviceJID(evt.ID.String())
	if err != nil {
		mc.logger.ErrorWithFields("Failed to update device JID", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"device_jid": evt.ID.String(),
			"error":      err.Error(),
		})
	}

	// Clear QR code
	err = mc.ClearQRCode()
	if err != nil {
		mc.logger.ErrorWithFields("Failed to clear QR code after pairing", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
	}

	// Notify gateway about successful pairing
	if mc.gateway != nil && mc.gateway.eventHandler != nil {
		deviceInfo := &session.DeviceInfo{
			Platform:    evt.Platform,
			DeviceModel: "WhatsApp-Client",
			OSVersion:   "11",
			AppVersion:  "2.23.24.76",
		}
		mc.gateway.eventHandler.OnSessionConnected(mc.sessionName, deviceInfo)
	}
}

// handleStreamReplaced handles stream replacement
func (mc *MyClient) handleStreamReplaced(evt *events.StreamReplaced) {
	mc.logger.InfoWithFields("Stream replaced", map[string]interface{}{
		"session_id": mc.sessionID.String(),
	})
	// Just log for now, no action needed
}

// handleLoggedOut handles logout events
func (mc *MyClient) handleLoggedOut(evt *events.LoggedOut) {
	mc.logger.InfoWithFields("Logged out from WhatsApp", map[string]interface{}{
		"session_id": mc.sessionID.String(),
		"reason":     evt.Reason.String(),
	})

	// Update connection status
	err := mc.UpdateConnectionStatus(false)
	if err != nil {
		mc.logger.ErrorWithFields("Failed to update connection status on logout", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
	}

	// Clear QR code
	err = mc.ClearQRCode()
	if err != nil {
		mc.logger.ErrorWithFields("Failed to clear QR code on logout", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
	}

	// Notify gateway about disconnection
	if mc.gateway != nil && mc.gateway.eventHandler != nil {
		mc.gateway.eventHandler.OnSessionDisconnected(mc.sessionName, evt.Reason.String())
	}

	// Remove from client manager
	clientManager := GetClientManager(mc.logger)
	clientManager.DeleteMyClient(mc.sessionID)
	clientManager.DeleteWhatsmeowClient(mc.sessionID)
	clientManager.DeleteHTTPClient(mc.sessionID)
}

// handleDisconnected handles disconnection events
func (mc *MyClient) handleDisconnected(evt *events.Disconnected) {
	mc.logger.InfoWithFields("Disconnected from WhatsApp", map[string]interface{}{
		"session_id": mc.sessionID.String(),
	})

	// Update connection status
	err := mc.UpdateConnectionStatus(false)
	if err != nil {
		mc.logger.ErrorWithFields("Failed to update connection status on disconnect", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
	}

	// Notify gateway about disconnection
	if mc.gateway != nil && mc.gateway.eventHandler != nil {
		mc.gateway.eventHandler.OnSessionDisconnected(mc.sessionName, "disconnected")
	}
}

// handleConnectFailure handles connection failure events
func (mc *MyClient) handleConnectFailure(evt *events.ConnectFailure) {
	mc.logger.ErrorWithFields("Failed to connect to WhatsApp", map[string]interface{}{
		"session_id": mc.sessionID.String(),
		"reason":     fmt.Sprintf("%+v", evt),
	})

	// Set connection error
	err := mc.SetConnectionError(fmt.Sprintf("Connection failed: %+v", evt))
	if err != nil {
		mc.logger.ErrorWithFields("Failed to set connection error", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
	}

	// Notify gateway about connection error
	if mc.gateway != nil && mc.gateway.eventHandler != nil {
		mc.gateway.eventHandler.OnConnectionError(mc.sessionName, fmt.Errorf("connection failed: %+v", evt))
	}
}

// handleQRCode handles QR code generation (called from gateway)
func (mc *MyClient) handleQRCode(qrCode string) error {
	mc.logger.InfoWithFields("QR code generated", map[string]interface{}{
		"session_id": mc.sessionID.String(),
	})

	// Generate QR code image
	image, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		mc.logger.ErrorWithFields("Failed to encode QR code", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
		return err
	}

	// Create base64 encoded QR code
	base64QRCode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)

	// Update QR code in database
	expiresAt := time.Now().Add(2 * time.Minute) // QR codes typically expire in 2 minutes
	err = mc.UpdateQRCode(base64QRCode, expiresAt)
	if err != nil {
		mc.logger.ErrorWithFields("Failed to update QR code in database", map[string]interface{}{
			"session_id": mc.sessionID.String(),
			"error":      err.Error(),
		})
		return err
	}

	// Display QR code in terminal for debugging
	qrterminal.GenerateHalfBlock(qrCode, qrterminal.L, os.Stdout)
	fmt.Println("QR code:", qrCode)

	// Notify gateway about QR code generation
	if mc.gateway != nil && mc.gateway.eventHandler != nil {
		mc.gateway.eventHandler.OnQRCodeGenerated(mc.sessionName, base64QRCode, expiresAt)
	}

	return nil
}
