# zpwoot API Examples

## Session Management

### Create Session

#### Minimal Request
```json
{
  "name": "my-simple-session"
}
```

#### With Webhook Configuration
```json
{
  "name": "webhook-session",
  "settings": {
    "webhook": {
      "enabled": true,
      "url": "https://api.example.com/webhook",
      "events": [
        "Message",
        "Receipt",
        "Connected",
        "Disconnected"
      ],
      "secret": "my-webhook-secret"
    }
  }
}
```

#### With Proxy Configuration
```json
{
  "name": "proxy-session",
  "settings": {
    "proxy": {
      "enabled": true,
      "type": "http",
      "host": "proxy.example.com",
      "port": "8080",
      "user": "proxyUser123",
      "pass": "proxyPass123"
    }
  }
}
```

#### Complete Configuration
```json
{
  "name": "my-session",
  "settings": {
    "proxy": {
      "enabled": true,
      "type": "http",
      "host": "proxy.example.com",
      "port": "8080",
      "user": "proxyUser123",
      "pass": "proxyPass123"
    },
    "webhook": {
      "enabled": true,
      "url": "https://api.example.com/webhook",
      "events": [
        "Message",
        "Receipt",
        "Connected",
        "Disconnected",
        "CallOffer",
        "Presence",
        "NewsletterJoin",
        "All"
      ],
      "secret": "supersecrettoken123"
    }
  }
}
```

### Available Webhook Events

- `Message` - New message received
- `Receipt` - Message receipt/read status
- `Connected` - Session connected
- `Disconnected` - Session disconnected
- `CallOffer` - Incoming call
- `Presence` - Contact presence update
- `NewsletterJoin` - Newsletter join event
- `All` - Subscribe to all events

### Proxy Types

- `http` - HTTP proxy
- `https` - HTTPS proxy
- `socks5` - SOCKS5 proxy

## cURL Examples

### Create Session (Minimal)
```bash
curl -X POST http://localhost:8080/sessions/create \
  -H "Authorization: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-session"
  }'
```

### Create Session (With Webhook)
```bash
curl -X POST http://localhost:8080/sessions/create \
  -H "Authorization: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "webhook-session",
    "settings": {
      "webhook": {
        "enabled": true,
        "url": "https://api.example.com/webhook",
        "events": ["Message", "Receipt"],
        "secret": "my-secret"
      }
    }
  }'
```

### List Sessions
```bash
curl -X GET http://localhost:8080/sessions/list \
  -H "Authorization: your-api-key"
```

### Get Session Info
```bash
curl -X GET http://localhost:8080/sessions/{sessionId}/info \
  -H "Authorization: your-api-key"
```

### Connect Session
```bash
curl -X POST http://localhost:8080/sessions/{sessionId}/connect \
  -H "Authorization: your-api-key"
```

### Get QR Code
```bash
curl -X GET http://localhost:8080/sessions/{sessionId}/qr \
  -H "Authorization: your-api-key"
```

### Send Message
```bash
curl -X POST http://localhost:8080/sessions/{sessionId}/messages \
  -H "Authorization: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "5511999999999",
    "message": "Hello from zpwoot!"
  }'
```

### Delete Session
```bash
curl -X DELETE http://localhost:8080/sessions/{sessionId}/delete \
  -H "Authorization: your-api-key"
```

## Response Examples

### Session Response
```json
{
  "sessionId": "550e8400-e29b-41d4-a716-446655440000",
  "name": "my-session",
  "status": "connected",
  "connected": true,
  "deviceJid": "5511999999999@s.whatsapp.net",
  "settings": {
    "webhook": {
      "enabled": true,
      "url": "https://api.example.com/webhook",
      "events": ["Message", "Receipt"],
      "secret": "my-secret"
    },
    "proxy": {
      "enabled": true,
      "type": "http",
      "host": "proxy.example.com",
      "port": "8080",
      "user": "proxyUser123"
    }
  },
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T10:35:00Z",
  "connectedAt": "2025-01-15T10:32:00Z",
  "lastSeen": "2025-01-15T10:35:00Z"
}
```

### QR Code Response
```json
{
  "qrCode": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
  "expiresAt": "2025-01-15T10:35:00Z",
  "status": "generated"
}
```

### Error Response
```json
{
  "error": "validation_error",
  "message": "Invalid request",
  "details": {
    "field": "name",
    "message": "session name is required"
  },
  "timestamp": "2025-01-15T10:30:00Z"
}
```

## Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid API key
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Authentication

All API endpoints (except `/`, `/health`, and `/swagger/*`) require authentication using an API key.

### Header Format

```
Authorization: your-api-key
```

Or alternatively:

```
X-API-Key: your-api-key
```

The API key is configured in the `.env` file as `ZP_API_KEY`.

