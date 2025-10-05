package waclient

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"zpwoot/internal/adapters/logger"

	"go.mau.fi/whatsmeow/types/events"
)


type DefaultEventHandler struct {
	logger        *logger.Logger
	webhookSender WebhookSender
}


func NewDefaultEventHandler(logger *logger.Logger, webhookSender WebhookSender) *DefaultEventHandler {
	return &DefaultEventHandler{
		logger:        logger,
		webhookSender: webhookSender,
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
		eh.logger.Debug().
			Str("event_type", fmt.Sprintf("%T", event)).
			Msg("Unhandled event type")
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


	if eh.shouldSendWebhook(client, EventMessage) {
		webhookData := map[string]interface{}{
			"messageInfo": messageInfo,
			"message":     evt.Message,
		}

		return eh.sendWebhook(client, EventMessage, webhookData)
	}

	return nil
}


func (eh *DefaultEventHandler) handleReceipt(client *Client, evt *events.Receipt) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Interface("message_ids", evt.MessageIDs).
		Msg("Receipt event in session")

	if eh.shouldSendWebhook(client, EventReadReceipt) {
		return eh.sendWebhook(client, EventReadReceipt, evt)
	}

	return nil
}


func (eh *DefaultEventHandler) handlePresence(client *Client, evt *events.Presence) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Str("from", evt.From.String()).
		Msg("Presence event in session")

	if eh.shouldSendWebhook(client, EventPresence) {
		return eh.sendWebhook(client, EventPresence, evt)
	}

	return nil
}


func (eh *DefaultEventHandler) handleChatPresence(client *Client, evt *events.ChatPresence) error {
	eh.logger.Debug().
		Str("session_name", client.Name).
		Str("chat", evt.Chat.String()).
		Msg("Chat presence event in session")

	if eh.shouldSendWebhook(client, EventChatPresence) {
		return eh.sendWebhook(client, EventChatPresence, evt)
	}

	return nil
}


func (eh *DefaultEventHandler) handleHistorySync(client *Client, evt *events.HistorySync) error {
	eh.logger.Info().
		Str("session_name", client.Name).
		Int("conversations_count", len(evt.Data.Conversations)).
		Msg("History sync event in session")

	if eh.shouldSendWebhook(client, EventHistorySync) {

		syncInfo := map[string]interface{}{
			"type":              evt.Data.SyncType,
			"conversationCount": len(evt.Data.Conversations),
		}

		return eh.sendWebhook(client, EventHistorySync, syncInfo)
	}

	return nil
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


func (eh *DefaultEventHandler) shouldSendWebhook(client *Client, eventType EventType) bool {
	if eh.webhookSender == nil || client.WebhookURL == "" {
		return false
	}


	if len(client.Events) == 0 {
		return true
	}

	for _, subscribedEvent := range client.Events {
		if subscribedEvent == eventType {
			return true
		}
	}

	return false
}


func (eh *DefaultEventHandler) sendWebhook(client *Client, eventType EventType, eventData interface{}) error {
	webhookEvent := &WebhookEvent{
		Type:      eventType,
		SessionID: client.SessionID,
		Event:     eventData,
		Timestamp: time.Now(),
	}


	ctx, cancel := context.WithTimeout(client.ctx, 30*time.Second)
	defer cancel()

	return eh.webhookSender.SendWebhook(ctx, webhookEvent)
}


func getMessageType(msg interface{}) string {
	if msg == nil {
		return "unknown"
	}


	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return "unknown"
	}

	var msgMap map[string]interface{}
	if err := json.Unmarshal(msgJSON, &msgMap); err != nil {
		return "unknown"
	}


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
