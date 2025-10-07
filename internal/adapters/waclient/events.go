package waclient

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"zpwoot/internal/adapters/logger"
	"zpwoot/internal/core/domain/webhook"
	"zpwoot/internal/core/ports/output"

	"github.com/google/uuid"
	"go.mau.fi/whatsmeow/types/events"
)

const (
	unknownMessageType = "unknown"
)

type DefaultEventHandler struct {
	logger        *logger.Logger
	webhookSender output.WebhookSender
	webhookRepo   webhook.Repository
}

func NewDefaultEventHandler(logger *logger.Logger, webhookSender output.WebhookSender, webhookRepo webhook.Repository) *DefaultEventHandler {
	return &DefaultEventHandler{
		logger:        logger,
		webhookSender: webhookSender,
		webhookRepo:   webhookRepo,
	}
}

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
		eh.logUnhandledEvent(event)
		return nil
	}
}

func (eh *DefaultEventHandler) handleMessage(client *Client, evt *events.Message) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Str("message_id", evt.Info.ID).
		Str("sender", evt.Info.Sender.String()).
		Msg("Message event in session")

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

	webhookData := map[string]interface{}{
		"messageInfo": messageInfo,
		"message":     evt.Message,
	}

	return eh.sendWebhookIfEnabled(client, EventMessage, webhookData)
}

func (eh *DefaultEventHandler) handleReceipt(client *Client, evt *events.Receipt) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Interface("message_ids", evt.MessageIDs).
		Msg("Receipt event in session")

	return eh.sendWebhookIfEnabled(client, EventReadReceipt, evt)
}

func (eh *DefaultEventHandler) handlePresence(client *Client, evt *events.Presence) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Str("from", evt.From.String()).
		Msg("Presence event in session")

	return eh.sendWebhookIfEnabled(client, EventPresence, evt)
}

func (eh *DefaultEventHandler) handleChatPresence(client *Client, evt *events.ChatPresence) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Str("chat", evt.Chat.String()).
		Msg("Chat presence event in session")

	return eh.sendWebhookIfEnabled(client, EventChatPresence, evt)
}

func (eh *DefaultEventHandler) handleHistorySync(client *Client, evt *events.HistorySync) error {
	eh.logger.Info().
		Str("session_name", client.Name).
		Int("conversations_count", len(evt.Data.Conversations)).
		Msg("History sync event in session")

	syncInfo := map[string]interface{}{
		"type":              evt.Data.SyncType,
		"conversationCount": len(evt.Data.Conversations),
	}

	return eh.sendWebhookIfEnabled(client, EventHistorySync, syncInfo)
}

func (eh *DefaultEventHandler) handleAppStateSyncComplete(client *Client, evt *events.AppStateSyncComplete) error {
	eh.logger.Info().
		Str("session_name", client.Name).
		Str("sync_name", string(evt.Name)).
		Msg("App state sync complete in session")

	return nil
}

func (eh *DefaultEventHandler) handlePushNameSetting(client *Client, _ *events.PushNameSetting) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Msg("Push name setting in session")

	return nil
}

func (eh *DefaultEventHandler) handleBlocklistChange(client *Client, _ *events.BlocklistChange) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Msg("Blocklist change in session")

	return nil
}

func (eh *DefaultEventHandler) handleGroupInfo(client *Client, evt *events.GroupInfo) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Str("group_jid", evt.JID.String()).
		Msg("Group info event in session")

	return nil
}

func (eh *DefaultEventHandler) handleJoinedGroup(client *Client, evt *events.JoinedGroup) error {
	eh.logger.Info().
		Str("session_name", client.Name).
		Str("group_jid", evt.GroupInfo.JID.String()).
		Msg("Joined group in session")

	return nil
}

func (eh *DefaultEventHandler) handleOfflineSyncPreview(client *Client, _ *events.OfflineSyncPreview) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Msg("Offline sync preview in session")

	return nil
}

func (eh *DefaultEventHandler) logUnhandledEvent(event interface{}) {
	eh.logger.Debug().
		Str("event_type", fmt.Sprintf("%T", event)).
		Msg("Unhandled event type")
}

func (eh *DefaultEventHandler) sendWebhookIfEnabled(client *Client, eventType EventType, eventData interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	webhookConfig, err := eh.webhookRepo.GetBySessionID(ctx, client.SessionID)
	if err != nil {
		if err.Error() == "webhook not found" {
			return nil
		}

		eh.logger.Error().Err(err).Str("session_id", client.SessionID).Msg("Failed to load webhook config")

		return nil
	}

	if !eh.shouldSendWebhook(webhookConfig, eventType) {
		return nil
	}

	return eh.sendWebhook(webhookConfig, eventType, eventData, client.SessionID)
}

func (eh *DefaultEventHandler) shouldSendWebhook(webhookConfig *webhook.Webhook, eventType EventType) bool {
	if eh.webhookSender == nil || !webhookConfig.Enabled {
		return false
	}

	if len(webhookConfig.Events) == 0 {
		return true
	}

	eventTypeStr := string(eventType)
	for _, subscribedEvent := range webhookConfig.Events {
		if subscribedEvent == eventTypeStr {
			return true
		}
	}

	return false
}

func (eh *DefaultEventHandler) sendWebhook(webhookConfig *webhook.Webhook, eventType EventType, eventData interface{}, sessionID string) error {
	var data map[string]interface{}
	if mapData, ok := eventData.(map[string]interface{}); ok {
		data = mapData
	} else {
		jsonData, err := json.Marshal(eventData)
		if err != nil {
			eh.logger.Error().Err(err).Msg("Failed to marshal event data")
			return err
		}

		if err := json.Unmarshal(jsonData, &data); err != nil {
			eh.logger.Error().Err(err).Msg("Failed to unmarshal event data")
			return err
		}
	}

	webhookEvent := &output.WebhookEvent{
		ID:        uuid.New().String(),
		Type:      string(eventType),
		SessionID: sessionID,
		Timestamp: time.Now(),
		Data:      data,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return eh.webhookSender.SendWebhook(ctx, webhookConfig.URL, webhookConfig.Secret, webhookEvent)
}

func getMessageType(msg interface{}) string {
	msgMap, ok := convertMessageToMap(msg)
	if !ok {
		return unknownMessageType
	}

	return detectMessageType(msgMap)
}

func convertMessageToMap(msg interface{}) (map[string]interface{}, bool) {
	if msg == nil {
		return nil, false
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, false
	}

	var msgMap map[string]interface{}
	if err := json.Unmarshal(msgJSON, &msgMap); err != nil {
		return nil, false
	}

	return msgMap, true
}

func detectMessageType(msgMap map[string]interface{}) string {
	if hasMessageField(msgMap, "conversation", "extendedTextMessage") {
		return "text"
	}

	if hasMessageField(msgMap, "imageMessage") {
		return "image"
	}
	if hasMessageField(msgMap, "audioMessage") {
		return "audio"
	}
	if hasMessageField(msgMap, "videoMessage") {
		return "video"
	}
	if hasMessageField(msgMap, "documentMessage") {
		return "document"
	}
	if hasMessageField(msgMap, "stickerMessage") {
		return "sticker"
	}

	if hasMessageField(msgMap, "locationMessage") {
		return "location"
	}
	if hasMessageField(msgMap, "liveLocationMessage") {
		return "liveLocation"
	}

	if hasMessageField(msgMap, "contactMessage") {
		return "contact"
	}
	if hasMessageField(msgMap, "contactsArrayMessage") {
		return "contacts"
	}

	if hasMessageField(msgMap, "buttonsMessage") {
		return "buttons"
	}
	if hasMessageField(msgMap, "listMessage") {
		return "list"
	}
	if hasMessageField(msgMap, "templateMessage") {
		return "template"
	}

	return unknownMessageType
}

func hasMessageField(msgMap map[string]interface{}, fields ...string) bool {
	for _, field := range fields {
		if _, ok := msgMap[field]; ok {
			return true
		}
	}
	return false
}

type EventFilter struct {
	AllowedEvents  []EventType
	BlockedChats   []string
	BlockedSenders []string
}

func (ef *EventFilter) ShouldProcess(eventType EventType, chat, sender string) bool {
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

	for _, blockedChat := range ef.BlockedChats {
		if blockedChat == chat {
			return false
		}
	}

	for _, blockedSender := range ef.BlockedSenders {
		if blockedSender == sender {
			return false
		}
	}

	return true
}
