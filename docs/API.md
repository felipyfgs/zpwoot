# zpwoot API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication

Todos os endpoints (exceto `/`, `/health` e `/swagger/*`) requerem autenticação via API Key.

**Header:**
```
Authorization: YOUR_API_KEY
```
ou
```
X-API-Key: YOUR_API_KEY
```

---

## 📋 Índice

- [Health & Info](#health--info)
- [Sessions](#sessions)
- [Messages](#messages)

---

## Health & Info

### GET `/`
Informações básicas do serviço.

**Autenticação:** Não requerida

**Response:**
```json
{
  "message": "zpwoot WhatsApp API is running",
  "version": "1.0.0",
  "service": "zpwoot"
}
```

### GET `/health`
Verifica saúde do serviço e banco de dados.

**Autenticação:** Não requerida

**Response:**
```json
{
  "status": "ok",
  "service": "zpwoot",
  "version": "1.0.0"
}
```

### GET `/swagger/*`
Documentação Swagger/OpenAPI interativa.

**Autenticação:** Não requerida

**URL:** `http://localhost:8080/swagger/index.html`

---

## Sessions

### POST `/sessions/create`
Cria uma nova sessão WhatsApp.

**Body:**
```json
{
  "name": "my-session"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "sessionId": "uuid-here",
    "name": "my-session",
    "status": "disconnected"
  }
}
```

### GET `/sessions/list`
Lista todas as sessões.

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "sessionId": "uuid-here",
      "name": "my-session",
      "connected": false,
      "loggedIn": false
    }
  ]
}
```

### GET `/sessions/{sessionId}/info`
Obtém informações de uma sessão específica.

**Response:**
```json
{
  "success": true,
  "data": {
    "sessionId": "uuid-here",
    "name": "my-session",
    "connected": true,
    "loggedIn": true,
    "deviceJid": "5511999999999@s.whatsapp.net"
  }
}
```

### DELETE `/sessions/{sessionId}/delete`
Deleta uma sessão.

**Response:**
```json
{
  "success": true,
  "message": "Session deleted successfully"
}
```

### POST `/sessions/{sessionId}/connect`
Conecta uma sessão WhatsApp.

**Response:**
```json
{
  "success": true,
  "data": {
    "sessionId": "uuid-here",
    "connected": true
  }
}
```

### POST `/sessions/{sessionId}/disconnect`
Desconecta uma sessão WhatsApp.

**Response:**
```json
{
  "success": true,
  "data": {
    "sessionId": "uuid-here",
    "connected": false
  }
}
```

### POST `/sessions/{sessionId}/logout`
Faz logout de uma sessão WhatsApp.

**Response:**
```json
{
  "success": true,
  "message": "Session logged out successfully"
}
```

### GET `/sessions/{sessionId}/qr`
Obtém QR Code para autenticação.

**Response:**
```json
{
  "success": true,
  "data": {
    "qrCode": "base64-image-data",
    "code": "qr-code-string"
  }
}
```

---

## Messages

### POST `/sessions/{sessionId}/send/message/text`
Envia mensagem de texto.

**Body:**
```json
{
  "phone": "5511999999999",
  "text": "Hello, World!",
  "contextInfo": {
    "stanzaId": "message-id-to-reply",
    "participant": "5511888888888@s.whatsapp.net"
  }
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/image`
Envia imagem. Suporta Base64, URL ou caminho de arquivo.

**Body:**
```json
{
  "phone": "5511999999999",
  "file": "https://example.com/image.jpg",
  "caption": "Check this out!",
  "viewOnce": false,
  "mimeType": "image/jpeg",
  "fileName": "image.jpg",
  "contextInfo": {
    "stanzaId": "message-id-to-reply"
  }
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/video`
Envia vídeo. Suporta Base64, URL ou caminho de arquivo.

**Body:**
```json
{
  "phone": "5511999999999",
  "file": "https://example.com/video.mp4",
  "caption": "Watch this!",
  "viewOnce": false,
  "mimeType": "video/mp4",
  "fileName": "video.mp4",
  "contextInfo": {
    "stanzaId": "message-id-to-reply"
  }
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/audio`
Envia áudio/voice note. Suporta Base64, URL ou caminho de arquivo.

**Body:**
```json
{
  "phone": "5511999999999",
  "file": "https://example.com/audio.mp3",
  "viewOnce": false,
  "mimeType": "audio/mpeg",
  "fileName": "audio.mp3",
  "contextInfo": {
    "stanzaId": "message-id-to-reply"
  }
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/document`
Envia documento. Suporta Base64, URL ou caminho de arquivo.

**Body:**
```json
{
  "phone": "5511999999999",
  "file": "https://example.com/document.pdf",
  "caption": "Important document",
  "mimeType": "application/pdf",
  "fileName": "document.pdf",
  "contextInfo": {
    "stanzaId": "message-id-to-reply"
  }
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/sticker`
Envia sticker. Suporta Base64, URL ou caminho de arquivo.

**Body:**
```json
{
  "phone": "5511999999999",
  "file": "https://example.com/sticker.webp",
  "mimeType": "image/webp",
  "fileName": "sticker.webp",
  "contextInfo": {
    "stanzaId": "message-id-to-reply"
  }
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/location`
Envia localização.

**Body:**
```json
{
  "phone": "5511999999999",
  "latitude": -23.550520,
  "longitude": -46.633308,
  "name": "São Paulo",
  "address": "Av. Paulista, 1578",
  "contextInfo": {
    "stanzaId": "message-id-to-reply"
  }
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/contact`
Envia contato.

**Body:**
```json
{
  "phone": "5511999999999",
  "contact": {
    "name": "John Doe",
    "phone": "5511888888888"
  },
  "contextInfo": {
    "stanzaId": "message-id-to-reply"
  }
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/contacts`
Envia múltiplos contatos.

**Body:**
```json
{
  "phone": "5511999999999",
  "contacts": [
    {
      "name": "John Doe",
      "phone": "5511888888888"
    },
    {
      "name": "Jane Doe",
      "phone": "5511777777777"
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/reaction`
Envia reação a uma mensagem.

**Body:**
```json
{
  "phone": "5511999999999",
  "messageId": "3EB0C767D0D1A6F4FD29",
  "reaction": "👍",
  "fromMe": false
}
```

**Nota:** Use `fromMe: true` se a mensagem foi enviada por você, ou prefixe o messageId com `me:` (ex: `"me:3EB0C767D0D1A6F4FD29"`).

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/poll`
Envia enquete.

**Body:**
```json
{
  "phone": "5511999999999",
  "name": "What's your favorite color?",
  "options": ["Red", "Blue", "Green", "Yellow"],
  "selectableOptionsCount": 1
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/buttons`
Envia mensagem com botões.

**Body:**
```json
{
  "phone": "5511999999999",
  "text": "Choose an option:",
  "buttons": [
    {
      "id": "1",
      "text": "Option 1"
    },
    {
      "id": "2",
      "text": "Option 2"
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/list`
Envia mensagem com lista de opções.

**Body:**
```json
{
  "phone": "5511999999999",
  "text": "Select an option",
  "title": "Menu",
  "buttonText": "View Options",
  "sections": [
    {
      "title": "Section 1",
      "rows": [
        {
          "id": "1",
          "title": "Option 1",
          "description": "Description 1"
        },
        {
          "id": "2",
          "title": "Option 2",
          "description": "Description 2"
        }
      ]
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "id": "message-id",
  "to": "5511999999999@s.whatsapp.net",
  "status": "sent"
}
```

### POST `/sessions/{sessionId}/send/message/template`
Envia mensagem template.

**Status:** Não implementado

**Response:**
```json
{
  "error": "not_implemented",
  "message": "Template messages not yet implemented"
}
```

---

## Exemplos Práticos

### Enviar Imagem ViewOnce
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/image \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "file": "https://picsum.photos/800/600",
    "caption": "Esta imagem desaparecerá após visualização",
    "viewOnce": true
  }'
```

### Enviar Texto com Resposta
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/text \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "text": "Respondendo sua mensagem",
    "contextInfo": {
      "stanzaId": "3EB0A9253FA64269E11C9D"
    }
  }'
```

### Criar e Conectar Sessão
```bash
# 1. Criar sessão
curl -X POST http://localhost:8080/sessions/create \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "my-session"}'

# 2. Obter QR Code
curl -X GET http://localhost:8080/sessions/SESSION_ID/qr \
  -H "Authorization: YOUR_API_KEY"

# 3. Conectar (após escanear QR)
curl -X POST http://localhost:8080/sessions/SESSION_ID/connect \
  -H "Authorization: YOUR_API_KEY"
```

---

## Códigos de Erro

| Código | Descrição |
|--------|-----------|
| 400 | Bad Request - Requisição inválida |
| 401 | Unauthorized - API Key inválida ou ausente |
| 404 | Not Found - Recurso não encontrado |
| 409 | Conflict - Conflito (ex: sessão já existe) |
| 412 | Precondition Failed - Sessão não conectada |
| 500 | Internal Server Error - Erro interno |
| 503 | Service Unavailable - Serviço indisponível |

---

## Notas Importantes

### Formato de Telefone
- Use o formato internacional sem `+`: `5511999999999`
- O sistema adiciona automaticamente `@s.whatsapp.net` se necessário

### Tipos de Arquivo Suportados
- **Imagem:** jpg, jpeg, png, gif, webp
- **Vídeo:** mp4, avi, mov, mkv
- **Áudio:** mp3, ogg, wav, m4a
- **Documento:** pdf, doc, docx, xls, xlsx, txt, etc.
- **Sticker:** webp

### ViewOnce
- Disponível para: imagem, vídeo, áudio
- Mensagem desaparece após visualização
- Use `"viewOnce": true` no body

### ContextInfo (Respostas)
- `stanzaId`: ID da mensagem a ser respondida
- `participant`: Necessário apenas em grupos (JID do participante)

---

## Swagger UI

Para documentação interativa completa, acesse:
```
http://localhost:8080/swagger/index.html
```
