# 📚 Swagger Documentation - Message Routes Summary

## ✅ All Message Routes Documented

All message routes have been fully documented with Swagger/OpenAPI comments. Below is the complete list of documented endpoints:

---

## 📋 Message Routes (19 endpoints)

### 1. Generic Message
```
POST /sessions/{sessionId}/messages
```
**Summary**: Send a message (generic)  
**Description**: Send a message of any type (text, media, location, contact) using a generic endpoint  
**Request Body**: `dto.SendMessageRequest`

---

### 2. Text Message
```
POST /sessions/{sessionId}/send/message/text
```
**Summary**: Send text message  
**Description**: Send a simple text message to a WhatsApp contact  
**Request Body**: `dto.SendTextMessageRequest`

---

### 3. Image Message
```
POST /sessions/{sessionId}/send/message/image
```
**Summary**: Send image message  
**Description**: Send an image message to a WhatsApp contact with optional caption  
**Request Body**: `dto.SendImageMessageRequest`

---

### 4. Audio Message
```
POST /sessions/{sessionId}/send/message/audio
```
**Summary**: Send audio message  
**Description**: Send an audio/voice message to a WhatsApp contact  
**Request Body**: `dto.SendAudioMessageRequest`

---

### 5. Video Message
```
POST /sessions/{sessionId}/send/message/video
```
**Summary**: Send video message  
**Description**: Send a video message to a WhatsApp contact with optional caption  
**Request Body**: `dto.SendVideoMessageRequest`

---

### 6. Document Message
```
POST /sessions/{sessionId}/send/message/document
```
**Summary**: Send document message  
**Description**: Send a document/file message to a WhatsApp contact  
**Request Body**: `dto.SendDocumentMessageRequest`

---

### 7. Sticker Message
```
POST /sessions/{sessionId}/send/message/sticker
```
**Summary**: Send sticker message  
**Description**: Send a sticker message to a WhatsApp contact  
**Request Body**: `dto.SendStickerMessageRequest`

---

### 8. Location Message
```
POST /sessions/{sessionId}/send/message/location
```
**Summary**: Send location message  
**Description**: Send a location message with GPS coordinates to a WhatsApp contact  
**Request Body**: `dto.SendLocationMessageRequest`

---

### 9. Contact Message
```
POST /sessions/{sessionId}/send/message/contact
```
**Summary**: Send contact message  
**Description**: Send a contact card to a WhatsApp contact  
**Request Body**: `dto.SendContactMessageRequest`

---

### 10. Multiple Contacts
```
POST /sessions/{sessionId}/send/message/contacts
```
**Summary**: Send multiple contacts  
**Description**: Send multiple contact cards to a WhatsApp contact  
**Request Body**: `dto.SendContactsArrayMessageRequest`

---

### 11. Reaction Message
```
POST /sessions/{sessionId}/send/message/reaction
```
**Summary**: Send reaction to message  
**Description**: React to a message with an emoji  
**Request Body**: `dto.SendReactionMessageRequest`

---

### 12. Poll Message
```
POST /sessions/{sessionId}/send/message/poll
```
**Summary**: Send poll message  
**Description**: Send a poll/survey with multiple options to a WhatsApp contact (2-12 options)  
**Request Body**: `dto.SendPollMessageRequest`

---

### 13. Buttons Message
```
POST /sessions/{sessionId}/send/message/buttons
```
**Summary**: Send buttons message  
**Description**: Send an interactive message with buttons (max 3 buttons)  
**Request Body**: `dto.SendButtonsMessageRequest`

---

### 14. List Message
```
POST /sessions/{sessionId}/send/message/list
```
**Summary**: Send list message  
**Description**: Send an interactive list message with sections and rows  
**Request Body**: `dto.SendListMessageRequest`

---

### 15. Template Message
```
POST /sessions/{sessionId}/send/message/template
```
**Summary**: Send template message  
**Description**: Send a WhatsApp Business template message  
**Request Body**: `dto.SendTemplateMessageRequest`

---

### 16. View Once Message
```
POST /sessions/{sessionId}/send/message/viewonce
```
**Summary**: Send view once message  
**Description**: Send a media message that can only be viewed once (disappears after viewing)  
**Request Body**: `dto.SendViewOnceMessageRequest`

---

## 📊 Common Response Schemas

### Success Response
```json
{
  "messageId": "msg_123456789",
  "status": "sent",
  "sentAt": "2025-10-06T10:30:00Z"
}
```
**Schema**: `dto.SendMessageResponse`

### Error Response
```json
{
  "error": "SESSION_NOT_FOUND",
  "message": "Session not found"
}
```
**Schema**: `dto.ErrorResponse`

---

## 🔐 Common HTTP Status Codes

| Code | Description |
|------|-------------|
| 200  | Message sent successfully |
| 400  | Invalid request (validation error) |
| 404  | Session not found |
| 412  | Session not connected |
| 500  | Internal server error |

---

## 🏷️ Swagger Tags

All message routes are grouped under the **`Messages`** tag for easy navigation in Swagger UI.

---

## 🚀 How to Generate Swagger Documentation

1. Install swag CLI:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate Swagger docs:
```bash
swag init -g cmd/zpwoot/main.go -o docs/swagger
```

3. Access Swagger UI:
```
http://localhost:8080/swagger/index.html
```

---

## 📝 Example Swagger Annotation

```go
// SendText godoc
// @Summary      Send text message
// @Description  Send a simple text message to a WhatsApp contact
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Param        sessionId   path      string                        true  "Session ID"
// @Param        message     body      dto.SendTextMessageRequest    true  "Text message data"
// @Success      200         {object}  dto.SendMessageResponse       "Message sent successfully"
// @Failure      400         {object}  dto.ErrorResponse             "Invalid request"
// @Failure      404         {object}  dto.ErrorResponse             "Session not found"
// @Failure      412         {object}  dto.ErrorResponse             "Session not connected"
// @Failure      500         {object}  dto.ErrorResponse             "Internal server error"
// @Router       /sessions/{sessionId}/send/message/text [post]
func (h *MessageHandler) SendText(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

---

## ✅ Documentation Status

- **Total Routes**: 19
- **Documented**: 19 (100%) ✅
- **Tags**: Messages
- **Schemas**: All DTOs documented
- **Status**: Complete and ready for Swagger generation

---

**Last Updated**: 2025-10-06  
**File**: `internal/adapters/http/handlers/message.go`

