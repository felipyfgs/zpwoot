package waclient

import (
	"context"
	"encoding/json"
	"time"

	"zpwoot/internal/adapters/logger"

	"go.mau.fi/whatsmeow/types/events"
)

// DefaultEventHandler implements the EventHandler interface
type DefaultEventHandler struct {
	logger        *logger.Logger
	webhookSender WebhookSender
}

// NewDefaultEventHandler creates a new default event handler
func NewDefaultEventHandler(logger *logger.Logger, webhookSender WebhookSender) *DefaultEventHandler {
	return &DefaultEventHandler{
		logger:        logger,
		webhookSender: webhookSender,
	}
}

// HandleEvent processes WhatsApp events
func (eh *DefaultEventHandler) HandleEvent(client *Client, event interface{}) error {
	switch evt := event.(type) {
	case *events.Message:
		return eh.handleMessage(client, evt)
	case *events.Receipt:
		return eh.handleReceipt(client, evt)
	case *events.Presence:
		return eh.handlePresence(client, evt)
	case *events.ChatPresence:
		return eh.handleChatPresence(client, evt)
	case *events.HistorySync:
		return eh.handleHistorySync(client, evt)
	case *events.AppStateSyncComplete:
		return eh.handleAppStateSyncComplete(client, evt)
	case *events.PushNameSetting:
		return eh.handlePushNameSetting(client, evt)
	case *events.BlocklistChange:
		return eh.handleBlocklistChange(client, evt)
	case *events.GroupInfo:
		return eh.handleGroupInfo(client, evt)
	case *events.JoinedGroup:
		return eh.handleJoinedGroup(client, evt)
	case *events.OfflineSyncPreview:
		return eh.handleOfflineSyncPreview(client, evt)
	default:
		eh.logger.Debugf("Unhandled event type: %T", event)
		return nil
	}
}

// handleMessage processes incoming messages
func (eh *DefaultEventHandler) handleMessage(client *Client, evt *events.Message) error {
	eh.logger.Debugf("Message event in session %s: %s from %s", 
		client.Name, evt.Info.ID, evt.Info.Sender.String())

	// Create message info
	messageInfo := &MessageInfo{
		ID:        evt.Info.ID,
		Chat:      evt.Info.Chat.String(),
		Sender:    evt.Info.Sender.String(),
		PushName:  evt.Info.PushName,
		Timestamp: evt.Info.Timestamp,
		FromMe:    evt.Info.IsFromMe,
		Type:      getMessageType(evt.Message),
		IsGroup:   evt.Info.IsGroup,
	}

	// Send webhook if configured
	if eh.shouldSendWebhook(client, EventMessage) {
		webhookData := map[string]interface{}{
			"messageInfo": messageInfo,
			"message":     evt.Message,
		}

		return eh.sendWebhook(client, EventMessage, webhookData)
	}

	return nil
}

// handleReceipt processes message receipts (read receipts, delivery receipts)
func (eh *DefaultEventHandler) handleReceipt(client *Client, evt *events.Receipt) error {
	eh.logger.Debugf("Receipt event in session %s: %s", client.Name, evt.MessageIDs)

	if eh.shouldSendWebhook(client, EventReadReceipt) {
		return eh.sendWebhook(client, EventReadReceipt, evt)
	}

	return nil
}

// handlePresence processes user presence updates
func (eh *DefaultEventHandler) handlePresence(client *Client, evt *events.Presence) error {
	eh.logger.Debugf("Presence event in session %s: %s",
		client.Name, evt.From.String())

	if eh.shouldSendWebhook(client, EventPresence) {
		return eh.sendWebhook(client, EventPresence, evt)
	}

	return nil
}

// handleChatPresence processes chat-specific presence (typing, recording, etc.)
func (eh *DefaultEventHandler) handleChatPresence(client *Client, evt *events.ChatPresence) error {
	eh.logger.Debugf("Chat presence event in session %s in %s",
		client.Name, evt.Chat.String())

	if eh.shouldSendWebhook(client, EventChatPresence) {
		return eh.sendWebhook(client, EventChatPresence, evt)
	}

	return nil
}

// handleHistorySync processes history sync events
func (eh *DefaultEventHandler) handleHistorySync(client *Client, evt *events.HistorySync) error {
	eh.logger.Infof("History sync event in session %s: %d conversations",
		client.Name, len(evt.Data.Conversations))

	if eh.shouldSendWebhook(client, EventHistorySync) {
		// Create a simplified version for webhook (full data might be too large)
		syncInfo := map[string]interface{}{
			"type":              evt.Data.SyncType,
			"conversationCount": len(evt.Data.Conversations),
		}

		return eh.sendWebhook(client, EventHistorySync, syncInfo)
	}

	return nil
}

// handleAppStateSyncComplete processes app state sync completion
func (eh *DefaultEventHandler) handleAppStateSyncComplete(client *Client, evt *events.AppStateSyncComplete) error {
	eh.logger.Infof("App state sync complete in session %s: %s", client.Name, evt.Name)
	return nil
}

// handlePushNameSetting processes push name setting events
func (eh *DefaultEventHandler) handlePushNameSetting(client *Client, evt *events.PushNameSetting) error {
	eh.logger.Debugf("Push name setting in session %s", client.Name)
	return nil
}

// handleBlocklistChange processes blocklist changes
func (eh *DefaultEventHandler) handleBlocklistChange(client *Client, evt *events.BlocklistChange) error {
	eh.logger.Debugf("Blocklist change in session %s", client.Name)
	return nil
}

// handleGroupInfo processes group info events
func (eh *DefaultEventHandler) handleGroupInfo(client *Client, evt *events.GroupInfo) error {
	eh.logger.Debugf("Group info event in session %s: %s", client.Name, evt.JID.String())
	return nil
}

// handleJoinedGroup processes joined group events
func (eh *DefaultEventHandler) handleJoinedGroup(client *Client, evt *events.JoinedGroup) error {
	eh.logger.Infof("Joined group in session %s: %s", client.Name, evt.GroupInfo.JID.String())
	return nil
}

// handleOfflineSyncPreview processes offline sync preview events
func (eh *DefaultEventHandler) handleOfflineSyncPreview(client *Client, evt *events.OfflineSyncPreview) error {
	eh.logger.Debugf("Offline sync preview in session %s", client.Name)
	return nil
}

// shouldSendWebhook checks if a webhook should be sent for the event type
func (eh *DefaultEventHandler) shouldSendWebhook(client *Client, eventType EventType) bool {
	if eh.webhookSender == nil || client.WebhookURL == "" {
		return false
	}

	// Check if event type is in the client's subscribed events
	if len(client.Events) == 0 {
		return true // Send all events if no specific events configured
	}

	for _, subscribedEvent := range client.Events {
		if subscribedEvent == eventType {
			return true
		}
	}

	return false
}

// sendWebhook sends a webhook for an event
func (eh *DefaultEventHandler) sendWebhook(client *Client, eventType EventType, eventData interface{}) error {
	webhookEvent := &WebhookEvent{
		Type:      eventType,
		SessionID: client.SessionID,
		Event:     eventData,
		Timestamp: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return eh.webhookSender.SendWebhook(ctx, webhookEvent)
}

// getMessageType determines the message type from the message content
func getMessageType(msg interface{}) string {
	if msg == nil {
		return "unknown"
	}

	// Use reflection or type assertion to determine message type
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return "unknown"
	}

	var msgMap map[string]interface{}
	if err := json.Unmarshal(msgJSON, &msgMap); err != nil {
		return "unknown"
	}

	// Check for different message types
	if _, ok := msgMap["conversation"]; ok {
		return "text"
	}
	if _, ok := msgMap["extendedTextMessage"]; ok {
		return "text"
	}
	if _, ok := msgMap["imageMessage"]; ok {
		return "image"
	}
	if _, ok := msgMap["audioMessage"]; ok {
		return "audio"
	}
	if _, ok := msgMap["videoMessage"]; ok {
		return "video"
	}
	if _, ok := msgMap["documentMessage"]; ok {
		return "document"
	}
	if _, ok := msgMap["stickerMessage"]; ok {
		return "sticker"
	}
	if _, ok := msgMap["locationMessage"]; ok {
		return "location"
	}
	if _, ok := msgMap["contactMessage"]; ok {
		return "contact"
	}
	if _, ok := msgMap["contactsArrayMessage"]; ok {
		return "contacts"
	}
	if _, ok := msgMap["liveLocationMessage"]; ok {
		return "liveLocation"
	}
	if _, ok := msgMap["buttonsMessage"]; ok {
		return "buttons"
	}
	if _, ok := msgMap["listMessage"]; ok {
		return "list"
	}
	if _, ok := msgMap["templateMessage"]; ok {
		return "template"
	}

	return "unknown"
}

// EventFilter allows filtering events before processing
type EventFilter struct {
	AllowedEvents []EventType
	BlockedChats  []string
	BlockedSenders []string
}

// ShouldProcess checks if an event should be processed based on filters
func (ef *EventFilter) ShouldProcess(eventType EventType, chat, sender string) bool {
	// Check allowed events
	if len(ef.AllowedEvents) > 0 {
		allowed := false
		for _, allowedEvent := range ef.AllowedEvents {
			if allowedEvent == eventType {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}

	// Check blocked chats
	for _, blockedChat := range ef.BlockedChats {
		if blockedChat == chat {
			return false
		}
	}

	// Check blocked senders
	for _, blockedSender := range ef.BlockedSenders {
		if blockedSender == sender {
			return false
		}
	}

	return true
}
