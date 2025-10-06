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

---

## 📋 Índice

- [Health & Info](#health--info)
- [Sessions](#sessions)
- [Messages](#messages)

---

## Health & Info

### GET `/`
Informações básicas do serviço.

**Autenticação:** ❌ Não requerida

**Exemplo:**
```bash
curl http://localhost:8080/
```

**Response:**
```json
{
  "message": "zpwoot WhatsApp API is running",
  "version": "1.0.0",
  "service": "zpwoot"
}
```

---

### GET `/health`
Verifica saúde do serviço e banco de dados.

**Autenticação:** ❌ Não requerida

**Exemplo:**
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "service": "zpwoot",
  "version": "1.0.0"
}
```

---

### GET `/swagger/*`
Documentação Swagger/OpenAPI interativa.

**Autenticação:** ❌ Não requerida

**URL:** `http://localhost:8080/swagger/index.html`

---

## Sessions

### POST `/sessions/create`
Cria uma nova sessão WhatsApp.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "name": "my-session"
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/create \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "my-session"}'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "my-session",
    "status": "disconnected",
    "connected": false,
    "createdAt": "2025-10-06T10:30:00Z"
  },
  "timestamp": "2025-10-06T10:30:00Z"
}
```

---

### GET `/sessions/list`
Lista todas as sessões.

**Autenticação:** ✅ Requerida

**Exemplo:**
```bash
curl http://localhost:8080/sessions/list \
  -H "Authorization: YOUR_API_KEY"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "sessions": [
      {
        "sessionId": "550e8400-e29b-41d4-a716-446655440000",
        "name": "my-session",
        "status": "connected",
        "connected": true,
        "deviceJid": "5511999999999@s.whatsapp.net",
        "createdAt": "2025-10-06T10:30:00Z",
        "updatedAt": "2025-10-06T10:35:00Z"
      }
    ],
    "total": 1
  },
  "timestamp": "2025-10-06T10:40:00Z"
}
```

---

### GET `/sessions/{sessionId}/info`
Obtém informações de uma sessão específica.

**Autenticação:** ✅ Requerida

**Exemplo:**
```bash
curl http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/info \
  -H "Authorization: YOUR_API_KEY"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "sessionId": "550e8400-e29b-41d4-a716-446655440000",
    "name": "my-session",
    "status": "connected",
    "connected": true,
    "deviceJid": "5511999999999@s.whatsapp.net",
    "createdAt": "2025-10-06T10:30:00Z",
    "updatedAt": "2025-10-06T10:35:00Z",
    "connectedAt": "2025-10-06T10:32:00Z"
  },
  "timestamp": "2025-10-06T10:40:00Z"
}
```

---

### DELETE `/sessions/{sessionId}/delete`
Deleta uma sessão.

**Autenticação:** ✅ Requerida

**Exemplo:**
```bash
curl -X DELETE http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/delete \
  -H "Authorization: YOUR_API_KEY"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "sessionId": "550e8400-e29b-41d4-a716-446655440000",
    "action": "delete",
    "status": "success",
    "message": "Session deleted successfully"
  },
  "timestamp": "2025-10-06T10:40:00Z"
}
```

---

### POST `/sessions/{sessionId}/connect`
Conecta uma sessão WhatsApp.

**Autenticação:** ✅ Requerida

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/connect \
  -H "Authorization: YOUR_API_KEY"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "sessionId": "550e8400-e29b-41d4-a716-446655440000",
    "action": "connect",
    "status": "success"
  },
  "timestamp": "2025-10-06T10:40:00Z"
}
```

---

### POST `/sessions/{sessionId}/disconnect`
Desconecta uma sessão WhatsApp.

**Autenticação:** ✅ Requerida

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/disconnect \
  -H "Authorization: YOUR_API_KEY"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "sessionId": "550e8400-e29b-41d4-a716-446655440000",
    "action": "disconnect",
    "status": "success"
  },
  "timestamp": "2025-10-06T10:40:00Z"
}
```

---

### POST `/sessions/{sessionId}/logout`
Faz logout de uma sessão WhatsApp.

**Autenticação:** ✅ Requerida

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/logout \
  -H "Authorization: YOUR_API_KEY"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "sessionId": "550e8400-e29b-41d4-a716-446655440000",
    "action": "logout",
    "status": "success",
    "message": "Session logged out successfully"
  },
  "timestamp": "2025-10-06T10:40:00Z"
}
```

---

### GET `/sessions/{sessionId}/qr`
Obtém QR Code para autenticação.

**Autenticação:** ✅ Requerida

**Exemplo:**
```bash
curl http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/qr \
  -H "Authorization: YOUR_API_KEY"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "qrCode": "2@abc123def456...",
    "qrCodeBase64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
    "expiresAt": "2025-10-06T10:35:00Z",
    "status": "qr_code"
  },
  "timestamp": "2025-10-06T10:30:00Z"
}
```

---

## Messages

### POST `/sessions/{sessionId}/send/message/text`
Envia mensagem de texto.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "text": "Hello, World!",
  "contextInfo": {
    "stanzaId": "3EB0A9253FA64269E11C9D"
  }
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/text \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "text": "Hello, World!"
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0A9253FA64269E11C9D",
  "to": "5511999999999@s.whatsapp.net",
  "type": "text",
  "content": "Hello, World!",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/image`
Envia imagem. Suporta Base64, URL ou caminho de arquivo.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "file": "https://picsum.photos/800/600",
  "caption": "Check this out!",
  "viewOnce": false,
  "mimeType": "image/jpeg",
  "fileName": "image.jpg",
  "contextInfo": {
    "stanzaId": "3EB0A9253FA64269E11C9D"
  }
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/image \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "file": "https://picsum.photos/800/600",
    "caption": "Check this out!"
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11C9E",
  "to": "5511999999999@s.whatsapp.net",
  "type": "image",
  "content": "Check this out!",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/video`
Envia vídeo. Suporta Base64, URL ou caminho de arquivo.

**Autenticação:** ✅ Requerida

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
    "stanzaId": "3EB0A9253FA64269E11C9D"
  }
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/video \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "file": "https://example.com/video.mp4",
    "caption": "Watch this!"
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11C9F",
  "to": "5511999999999@s.whatsapp.net",
  "type": "video",
  "content": "Watch this!",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/audio`
Envia áudio/voice note. Suporta Base64, URL ou caminho de arquivo.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "file": "https://example.com/audio.mp3",
  "viewOnce": false,
  "mimeType": "audio/mpeg",
  "fileName": "audio.mp3",
  "contextInfo": {
    "stanzaId": "3EB0A9253FA64269E11C9D"
  }
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/audio \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "file": "https://example.com/audio.mp3"
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11CA0",
  "to": "5511999999999@s.whatsapp.net",
  "type": "audio",
  "content": "",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/document`
Envia documento. Suporta Base64, URL ou caminho de arquivo.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "file": "https://example.com/document.pdf",
  "caption": "Important document",
  "mimeType": "application/pdf",
  "fileName": "document.pdf",
  "contextInfo": {
    "stanzaId": "3EB0A9253FA64269E11C9D"
  }
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/document \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "file": "https://example.com/document.pdf",
    "fileName": "document.pdf"
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11CA1",
  "to": "5511999999999@s.whatsapp.net",
  "type": "document",
  "content": "Important document",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/sticker`
Envia sticker. Suporta Base64, URL ou caminho de arquivo.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "file": "https://example.com/sticker.webp",
  "mimeType": "image/webp",
  "fileName": "sticker.webp",
  "contextInfo": {
    "stanzaId": "3EB0A9253FA64269E11C9D"
  }
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/sticker \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "file": "https://example.com/sticker.webp"
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11CA2",
  "to": "5511999999999@s.whatsapp.net",
  "type": "sticker",
  "content": "",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/location`
Envia localização.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "latitude": -23.550520,
  "longitude": -46.633308,
  "name": "Avenida Paulista",
  "address": "Av. Paulista, 1578 - São Paulo, SP",
  "contextInfo": {
    "stanzaId": "3EB0A9253FA64269E11C9D"
  }
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/location \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "latitude": -23.550520,
    "longitude": -46.633308,
    "name": "Avenida Paulista"
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11CA3",
  "to": "5511999999999@s.whatsapp.net",
  "type": "location",
  "content": "Avenida Paulista",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/contact`
Envia contato.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "contact": {
    "name": "João Silva",
    "phone": "5511888888888"
  },
  "contextInfo": {
    "stanzaId": "3EB0A9253FA64269E11C9D"
  }
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/contact \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "contact": {
      "name": "João Silva",
      "phone": "5511888888888"
    }
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11CA4",
  "to": "5511999999999@s.whatsapp.net",
  "type": "contact",
  "content": "João Silva",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/contacts`
Envia múltiplos contatos.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "contacts": [
    {
      "name": "João Silva",
      "phone": "5511888888888"
    },
    {
      "name": "Maria Santos",
      "phone": "5511777777777"
    }
  ]
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/contacts \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "contacts": [
      {"name": "João Silva", "phone": "5511888888888"},
      {"name": "Maria Santos", "phone": "5511777777777"}
    ]
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11CA5",
  "to": "5511999999999@s.whatsapp.net",
  "type": "contacts",
  "content": "2 contacts",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/reaction`
Envia reação a uma mensagem.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "messageId": "3EB0C767D0D1A6F4FD29",
  "reaction": "👍",
  "fromMe": false
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/reaction \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "messageId": "3EB0C767D0D1A6F4FD29",
    "reaction": "👍"
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11CA6",
  "to": "5511999999999@s.whatsapp.net",
  "type": "reaction",
  "content": "👍",
  "timestamp": 1696570882,
  "status": "sent"
}
```

**Nota:** Use `fromMe: true` se a mensagem foi enviada por você.

---

### POST `/sessions/{sessionId}/send/message/poll`
Envia enquete.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "name": "Qual sua cor favorita?",
  "options": ["Vermelho", "Azul", "Verde", "Amarelo"],
  "selectableOptionsCount": 1
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/poll \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "name": "Qual sua cor favorita?",
    "options": ["Vermelho", "Azul", "Verde", "Amarelo"],
    "selectableOptionsCount": 1
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11CA7",
  "to": "5511999999999@s.whatsapp.net",
  "type": "poll",
  "content": "Qual sua cor favorita?",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/buttons`
Envia mensagem com botões.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "text": "Escolha uma opção:",
  "buttons": [
    {
      "id": "1",
      "text": "Opção 1"
    },
    {
      "id": "2",
      "text": "Opção 2"
    }
  ]
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/buttons \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "text": "Escolha uma opção:",
    "buttons": [
      {"id": "1", "text": "Opção 1"},
      {"id": "2", "text": "Opção 2"}
    ]
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11CA8",
  "to": "5511999999999@s.whatsapp.net",
  "type": "buttons",
  "content": "Escolha uma opção:",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/list`
Envia mensagem com lista de opções.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "phone": "5511999999999",
  "text": "Selecione uma opção",
  "title": "Menu",
  "buttonText": "Ver Opções",
  "sections": [
    {
      "title": "Seção 1",
      "rows": [
        {
          "id": "1",
          "title": "Opção 1",
          "description": "Descrição 1"
        },
        {
          "id": "2",
          "title": "Opção 2",
          "description": "Descrição 2"
        }
      ]
    }
  ]
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/send/message/list \
  -H "Authorization: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "text": "Selecione uma opção",
    "title": "Menu",
    "buttonText": "Ver Opções",
    "sections": [{
      "title": "Seção 1",
      "rows": [
        {"id": "1", "title": "Opção 1", "description": "Descrição 1"}
      ]
    }]
  }'
```

**Response:**
```json
{
  "success": true,
  "id": "3EB0B1234FA64269E11CA9",
  "to": "5511999999999@s.whatsapp.net",
  "type": "list",
  "content": "Selecione uma opção",
  "timestamp": 1696570882,
  "status": "sent"
}
```

---

### POST `/sessions/{sessionId}/send/message/template`
Envia mensagem template.

**Autenticação:** ✅ Requerida

**Status:** ⚠️ Não implementado

**Response:**
```json
{
  "error": "not_implemented",
  "message": "Template messages not yet implemented"
}
```

---

## Respostas de Erro

Todos os endpoints retornam erros no seguinte formato:

```json
{
  "error": "error_code",
  "message": "Human readable error message"
}
```

### Códigos de Erro Comuns

| Código HTTP | Error Code | Descrição |
|-------------|------------|-----------|
| 400 | `validation_error` | Dados inválidos na requisição |
| 400 | `invalid_request` | JSON inválido |
| 400 | `media_processing_error` | Erro ao processar mídia |
| 400 | `invalid_jid` | Número de telefone inválido |
| 401 | `unauthorized` | API Key inválida ou ausente |
| 404 | `session_not_found` | Sessão não encontrada |
| 409 | `session_already_exists` | Sessão já existe |
| 412 | `not_connected` | Sessão não conectada |
| 500 | `internal_error` | Erro interno do servidor |
| 500 | `whatsapp_error` | Erro do WhatsApp |
| 501 | `not_implemented` | Funcionalidade não implementada |

---

## Notas Importantes

### 📱 Formato de Telefone
- Use formato internacional sem `+`: `5511999999999`
- Código do país + DDD + número
- Sistema adiciona `@s.whatsapp.net` automaticamente

### 📁 Tipos de Arquivo Suportados

**Formatos de entrada:**
- **URL**: `https://example.com/file.jpg`
- **Base64**: `data:image/jpeg;base64,/9j/4AAQ...`
- **Caminho local**: `/path/to/file.jpg`

**Tipos de mídia:**
- **Imagem**: jpg, jpeg, png, gif, webp
- **Vídeo**: mp4, avi, mov, mkv
- **Áudio**: mp3, ogg, wav, m4a, opus
- **Documento**: pdf, doc, docx, xls, xlsx, txt, zip, etc.
- **Sticker**: webp

### 👁️ ViewOnce (Visualização Única)
- Disponível para: **imagem**, **vídeo**, **áudio**
- Mensagem desaparece após visualização
- Adicione `"viewOnce": true` no body
- Pode ser combinado com `contextInfo` (respostas)

### 💬 ContextInfo (Respostas/Citações)
- `stanzaId`: ID da mensagem a ser respondida (obrigatório)
- `participant`: JID do participante (apenas para grupos)
- Exemplo: `{"stanzaId": "3EB0A9253FA64269E11C9D"}`

### 🔄 Status da Sessão
- `disconnected`: Sessão criada mas não conectada
- `connecting`: Conectando ao WhatsApp
- `qr_code`: Aguardando scan do QR Code
- `connected`: Conectado e pronto para uso
- `error`: Erro na conexão

---

## Swagger UI

Para documentação interativa completa com exemplos e testes:

```
http://localhost:8080/swagger/index.html
```
