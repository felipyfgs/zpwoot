# zpwoot API - Exemplos de Uso

Este documento contém exemplos práticos de como usar a API zpwoot para enviar mensagens pelo WhatsApp.

## 📋 Índice

- [Autenticação](#autenticação)
- [Gerenciamento de Sessões](#gerenciamento-de-sessões)
- [Envio de Mensagens](#envio-de-mensagens)
  - [Mensagem de Texto](#mensagem-de-texto)
  - [Mensagem de Imagem](#mensagem-de-imagem)
  - [Mensagem de Áudio](#mensagem-de-áudio)
  - [Mensagem de Vídeo](#mensagem-de-vídeo)
  - [Mensagem de Documento](#mensagem-de-documento)
  - [Mensagem de Localização](#mensagem-de-localização)
  - [Mensagem de Contato](#mensagem-de-contato)
  - [Mensagem de Reação](#mensagem-de-reação)
  - [Mensagem de Enquete](#mensagem-de-enquete)
  - [Mensagem com Botões](#mensagem-com-botões)
  - [Mensagem com Lista](#mensagem-com-lista)

---

## 🔐 Autenticação

Todas as requisições (exceto `/health` e `/`) requerem autenticação via API Key no header:

```bash
Authorization: YOUR_API_KEY_HERE
```

Configure a API Key na variável de ambiente `ZP_API_KEY`.

---

## 📱 Gerenciamento de Sessões

### Criar uma Nova Sessão

```bash
curl -X POST http://localhost:8080/sessions/create \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "name": "my-whatsapp-session",
    "qrCode": true
  }'
```

**Resposta:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "my-whatsapp-session",
  "status": "disconnected",
  "connected": false,
  "qrCode": "2@abc123...",
  "qrCodeBase64": "data:image/png;base64,iVBORw0KGgo...",
  "qrCodeExpiresAt": "2025-01-15T10:35:00Z",
  "createdAt": "2025-01-15T10:30:00Z"
}
```

### Listar Todas as Sessões

```bash
curl -X GET http://localhost:8080/sessions \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba"
```

### Obter QR Code de uma Sessão

```bash
curl -X GET http://localhost:8080/sessions/my-session/qr \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba"
```

### Conectar uma Sessão

```bash
curl -X POST http://localhost:8080/sessions/my-session/connect \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba"
```

---

## 💬 Envio de Mensagens

### Mensagem de Texto

Envie uma mensagem de texto simples para um contato.

```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/text \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "text": "Hello! This is a test message from zpwoot API."
  }'
```

**Resposta:**
```json
{
  "messageId": "msg_123456789",
  "status": "sent",
  "sentAt": "2024-01-15T10:30:00Z"
}
```

---

### Mensagem de Imagem

Envie uma imagem com legenda opcional.

**Usando URL:**
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/image \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "image": {
      "url": "https://example.com/image.jpg",
      "mimeType": "image/jpeg",
      "fileName": "photo.jpg"
    },
    "caption": "Check out this beautiful image!"
  }'
```

**Usando Base64:**
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/image \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "image": {
      "base64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
      "mimeType": "image/png",
      "fileName": "image.png"
    },
    "caption": "Image sent via base64"
  }'
```

---

### Mensagem de Áudio

Envie um arquivo de áudio ou nota de voz.

```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/audio \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "audio": {
      "url": "https://example.com/audio.mp3",
      "mimeType": "audio/mpeg",
      "fileName": "voice_note.mp3"
    }
  }'
```

---

### Mensagem de Vídeo

Envie um vídeo com legenda opcional.

```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/video \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "video": {
      "url": "https://example.com/video.mp4",
      "mimeType": "video/mp4",
      "fileName": "video.mp4"
    },
    "caption": "Watch this amazing video!"
  }'
```

---

### Mensagem de Documento

Envie um documento (PDF, DOC, XLS, etc).

```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/document \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "document": {
      "url": "https://example.com/document.pdf",
      "mimeType": "application/pdf",
      "fileName": "report.pdf"
    },
    "caption": "Important document attached"
  }'
```

---

### Mensagem de Localização

Envie uma localização GPS.

```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/location \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "latitude": -23.550520,
    "longitude": -46.633308,
    "name": "São Paulo",
    "address": "Av. Paulista, 1578 - Bela Vista, São Paulo - SP"
  }'
```

---

### Mensagem de Contato

Envie um cartão de contato.

```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/contact \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "contact": {
      "name": "John Doe",
      "phone": "5511888888888"
    }
  }'
```

---

### Mensagem de Reação

Reaja a uma mensagem com um emoji.

```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/reaction \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "messageId": "3EB0C767D0D1A6F4FD29",
    "reaction": "👍"
  }'
```

---

### Mensagem de Enquete

Crie uma enquete com múltiplas opções.

```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/poll \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "name": "What'\''s your favorite color?",
    "options": ["Red", "Blue", "Green", "Yellow"],
    "selectableOptionsCount": 1
  }'
```

---

### Mensagem com Botões

Envie uma mensagem interativa com até 3 botões.

```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/buttons \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "text": "Please choose an option:",
    "buttons": [
      {
        "id": "btn_1",
        "text": "Option 1"
      },
      {
        "id": "btn_2",
        "text": "Option 2"
      },
      {
        "id": "btn_3",
        "text": "Option 3"
      }
    ]
  }'
```

---

### Mensagem com Lista

Envie uma mensagem com lista de opções organizadas em seções.

```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/list \
  -H "Content-Type: application/json" \
  -H "Authorization: a0b1125a0eb3364d98e2c49ec6f7d6ba" \
  -d '{
    "to": "5511999999999",
    "text": "Please select an option from the list",
    "title": "Menu Options",
    "sections": [
      {
        "title": "Main Dishes",
        "rows": [
          {
            "id": "row_1",
            "title": "Pizza",
            "description": "Delicious Italian pizza"
          },
          {
            "id": "row_2",
            "title": "Burger",
            "description": "Juicy beef burger"
          }
        ]
      },
      {
        "title": "Drinks",
        "rows": [
          {
            "id": "row_3",
            "title": "Soda",
            "description": "Refreshing soda"
          }
        ]
      }
    ]
  }'
```

---

## 📊 Respostas de Erro

Todas as respostas de erro seguem o formato padrão:

```json
{
  "success": false,
  "error": {
    "code": "validation_error",
    "message": "Validation failed"
  },
  "timestamp": "2025-01-15T10:30:00Z"
}
```

### Códigos de Erro Comuns

- `400` - Bad Request (dados inválidos)
- `404` - Session Not Found (sessão não encontrada)
- `412` - Session Not Connected (sessão não conectada)
- `500` - Internal Server Error (erro interno)

---

## 🔗 Links Úteis

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health
- **API Info**: http://localhost:8080/

---

**zpwoot** - Making WhatsApp Business API integration simple and powerful! 🚀

