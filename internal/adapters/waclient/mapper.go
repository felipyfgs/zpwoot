package waclient

import (
	"fmt"
	"time"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	"zpwoot/internal/core/session"
	"zpwoot/platform/logger"
)

type Mapper struct {
	logger *logger.Logger
}

func NewMapper(logger *logger.Logger) *Mapper {
	return &Mapper{
		logger: logger,
	}
}

func (m *Mapper) MapDeviceInfoFromWhatsmeow(jid types.JID, pushName string, connected bool) *session.DeviceInfo {
	deviceInfo := &session.DeviceInfo{
		Platform:    "android",
		DeviceModel: "zpwoot-client",
		OSVersion:   "11",
		AppVersion:  "2.23.24.76",
	}

	m.logger.DebugWithFields("Mapped device info from whatsmeow", map[string]interface{}{
		"jid":       jid.String(),
		"push_name": pushName,
		"connected": connected,
		"platform":  deviceInfo.Platform,
	})

	return deviceInfo
}

func (m *Mapper) MapDeviceInfoFromPairSuccess(evt *events.PairSuccess) *session.DeviceInfo {
	deviceInfo := &session.DeviceInfo{
		Platform:    evt.Platform,
		DeviceModel: "WhatsApp-Client",
		OSVersion:   "11",
		AppVersion:  "2.23.24.76",
	}

	m.logger.DebugWithFields("Mapped device info from pair success", map[string]interface{}{
		"jid":           evt.ID.String(),
		"business_name": evt.BusinessName,
		"platform":      deviceInfo.Platform,
	})

	return deviceInfo
}

func (m *Mapper) MapConnectionStatus(evt interface{}) (session.SessionStatus, string) {
	switch e := evt.(type) {
	case *events.Connected:
		m.logger.DebugWithFields("Mapped connection status", map[string]interface{}{
			"event":  "Connected",
			"status": "connected",
		})
		return session.StatusConnected, "connected"

	case *events.Disconnected:
		m.logger.DebugWithFields("Mapped connection status", map[string]interface{}{
			"event":  "Disconnected",
			"status": "disconnected",
		})
		return session.StatusDisconnected, "disconnected"

	case *events.LoggedOut:
		reason := e.Reason.String()
		m.logger.DebugWithFields("Mapped connection status", map[string]interface{}{
			"event":  "LoggedOut",
			"status": "logged_out",
			"reason": reason,
		})
		return session.StatusLoggedOut, fmt.Sprintf("logged out: %s", reason)

	case *events.ConnectFailure:
		reason := fmt.Sprintf("%+v", e)
		m.logger.DebugWithFields("Mapped connection status", map[string]interface{}{
			"event":  "ConnectFailure",
			"status": "error",
			"reason": reason,
		})
		return session.StatusError, fmt.Sprintf("connection failed: %s", reason)

	default:
		m.logger.WarnWithFields("Unknown event type for connection status mapping", map[string]interface{}{
			"event_type": fmt.Sprintf("%T", evt),
		})
		return session.StatusError, "unknown event"
	}
}

func (m *Mapper) MapQRCodeResponse(qrCode string, base64Image string, expiresAt time.Time) *session.QRCodeResponse {
	timeout := int(time.Until(expiresAt).Seconds())
	if timeout < 0 {
		timeout = 0
	}

	response := &session.QRCodeResponse{
		QRCode:      qrCode,
		QRCodeImage: base64Image,
		ExpiresAt:   expiresAt,
		Timeout:     timeout,
	}

	m.logger.DebugWithFields("Mapped QR code response", map[string]interface{}{
		"expires_at": expiresAt,
		"timeout":    timeout,
		"has_image":  base64Image != "",
	})

	return response
}

func (m *Mapper) MapEventTypeToString(evt interface{}) string {
	switch evt.(type) {
	case *events.Connected:
		return "Connected"
	case *events.Disconnected:
		return "Disconnected"
	case *events.LoggedOut:
		return "LoggedOut"
	case *events.PairSuccess:
		return "PairSuccess"
	case *events.ConnectFailure:
		return "ConnectFailure"
	case *events.StreamReplaced:
		return "StreamReplaced"
	case *events.PushNameSetting:
		return "PushNameSetting"
	case *events.AppStateSyncComplete:
		return "AppStateSyncComplete"
	default:
		eventType := fmt.Sprintf("%T", evt)
		m.logger.DebugWithFields("Mapped unknown event type", map[string]interface{}{
			"event_type": eventType,
		})
		return eventType
	}
}

func (m *Mapper) MapJIDToString(jid types.JID) string {
	if jid.IsEmpty() {
		return ""
	}
	return jid.String()
}

func (m *Mapper) MapStringToJID(jidStr string) (types.JID, error) {
	if jidStr == "" {
		return types.JID{}, fmt.Errorf("JID string cannot be empty")
	}

	jid, err := types.ParseJID(jidStr)
	if err != nil {
		m.logger.ErrorWithFields("Failed to parse JID string", map[string]interface{}{
			"jid_string": jidStr,
			"error":      err.Error(),
		})
		return types.JID{}, fmt.Errorf("failed to parse JID: %w", err)
	}

	return jid, nil
}

func (m *Mapper) MapErrorToConnectionError(err error) string {
	if err == nil {
		return ""
	}

	errorMsg := err.Error()
	m.logger.DebugWithFields("Mapped error to connection error", map[string]interface{}{
		"original_error": errorMsg,
	})

	return errorMsg
}

func (m *Mapper) MapSessionStatusToString(status session.SessionStatus) string {
	switch status {
	case session.StatusCreated:
		return "created"
	case session.StatusConnecting:
		return "connecting"
	case session.StatusConnected:
		return "connected"
	case session.StatusDisconnected:
		return "disconnected"
	case session.StatusError:
		return "error"
	case session.StatusLoggedOut:
		return "logged_out"
	default:
		return "unknown"
	}
}

func (m *Mapper) MapStringToSessionStatus(statusStr string) session.SessionStatus {
	switch statusStr {
	case "created":
		return session.StatusCreated
	case "connecting":
		return session.StatusConnecting
	case "connected":
		return session.StatusConnected
	case "disconnected":
		return session.StatusDisconnected
	case "error":
		return session.StatusError
	case "logged_out":
		return session.StatusLoggedOut
	default:
		m.logger.WarnWithFields("Unknown status string", map[string]interface{}{
			"status_string": statusStr,
		})
		return session.StatusError
	}
}

func (m *Mapper) IsConnectionEvent(evt interface{}) bool {
	switch evt.(type) {
	case *events.Connected, *events.Disconnected, *events.LoggedOut,
		*events.ConnectFailure, *events.PairSuccess:
		return true
	default:
		return false
	}
}

func (m *Mapper) IsQREvent(evt interface{}) bool {

	return false
}

func (m *Mapper) GetEventPriority(evt interface{}) string {
	switch evt.(type) {
	case *events.Connected, *events.PairSuccess:
		return "high"
	case *events.Disconnected, *events.LoggedOut, *events.ConnectFailure:
		return "high"
	case *events.StreamReplaced:
		return "medium"
	case *events.PushNameSetting, *events.AppStateSyncComplete:
		return "low"
	default:
		return "low"
	}
}

func (m *Mapper) CreateEventContext(evt interface{}, sessionID string, sessionName string) map[string]interface{} {
	context := map[string]interface{}{
		"session_id":   sessionID,
		"session_name": sessionName,
		"event_type":   m.MapEventTypeToString(evt),
		"priority":     m.GetEventPriority(evt),
		"timestamp":    time.Now(),
	}

	switch e := evt.(type) {
	case *events.PairSuccess:
		context["device_jid"] = e.ID.String()
		context["business_name"] = e.BusinessName
		context["platform"] = e.Platform
	case *events.LoggedOut:
		context["reason"] = e.Reason.String()
	case *events.ConnectFailure:
		context["failure_reason"] = fmt.Sprintf("%+v", e)
	}

	return context
}
