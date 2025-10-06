# Mapeamento das Rotas de Mensagens - WhatsApp API

## Status Atual das Rotas

### ✅ Rotas Implementadas e Ativas

#### Rotas Gerais
| Rota | Método | Handler | Status | Descrição |
|------|--------|---------|--------|-----------|
| `POST /sessions/{sessionId}/messages` | SendMessage | h.Message.SendMessage | ✅ Ativo | Envio genérico de mensagens |
| `GET /sessions/{sessionId}/chats` | GetChats | h.Message.GetChats | ✅ Ativo | Listar conversas |
| `GET /sessions/{sessionId}/contacts` | GetContacts | h.Message.GetContacts | ✅ Ativo | Listar contatos |
| `GET /sessions/{sessionId}/chat-info` | GetChatInfo | h.Message.GetChatInfo | ✅ Ativo | Informações do chat |

#### Rotas de Mensagens Básicas
| Rota | Método | Handler | Status | Whatsmeow Method |
|------|--------|---------|--------|------------------|
| `POST /sessions/{sessionId}/send/message/text` | SendText | h.Message.SendText | ✅ Ativo | SendMessage + Conversation |
| `POST /sessions/{sessionId}/send/message/image` | SendImage | h.Message.SendImage | ✅ Ativo | Upload + SendMessage + ImageMessage |
| `POST /sessions/{sessionId}/send/message/audio` | SendAudio | h.Message.SendAudio | ✅ Ativo | Upload + SendMessage + AudioMessage |
| `POST /sessions/{sessionId}/send/message/video` | SendVideo | h.Message.SendVideo | ✅ Ativo | Upload + SendMessage + VideoMessage |
| `POST /sessions/{sessionId}/send/message/document` | SendDocument | h.Message.SendDocument | ✅ Ativo | Upload + SendMessage + DocumentMessage |
| `POST /sessions/{sessionId}/send/message/sticker` | SendSticker | h.Message.SendSticker | ✅ Ativo | Upload + SendMessage + StickerMessage |
| `POST /sessions/{sessionId}/send/message/location` | SendLocation | h.Message.SendLocation | ✅ Ativo | SendMessage + LocationMessage |
| `POST /sessions/{sessionId}/send/message/contact` | SendContact | h.Message.SendContact | ✅ Ativo | SendMessage + ContactMessage |
| `POST /sessions/{sessionId}/send/message/contacts` | SendContactsArray | h.Message.SendContactsArray | ✅ Ativo | SendMessage + ContactsArrayMessage |

#### Rotas de Mensagens Avançadas
| Rota | Método | Handler | Status | Whatsmeow Method |
|------|--------|---------|--------|------------------|
| `POST /sessions/{sessionId}/send/message/reaction` | SendReaction | h.Message.SendReaction | ✅ Ativo | BuildReaction + SendMessage |
| `POST /sessions/{sessionId}/send/message/poll` | SendPoll | h.Message.SendPoll | ✅ Ativo | BuildPollCreation + SendMessage |
| `POST /sessions/{sessionId}/send/message/buttons` | SendButtons | h.Message.SendButtons | ✅ Ativo | SendMessage + ButtonsMessage |
| `POST /sessions/{sessionId}/send/message/list` | SendList | h.Message.SendList | ✅ Ativo | SendMessage + ListMessage |
| `POST /sessions/{sessionId}/send/message/template` | SendTemplate | h.Message.SendTemplate | ✅ Ativo | SendMessage + TemplateMessage |
| `POST /sessions/{sessionId}/send/message/viewonce` | SendViewOnce | h.Message.SendViewOnce | ✅ Ativo | Upload + SendMessage + ViewOnce flag |

### 📊 Estatísticas
- **Total de Rotas**: 19
- **Rotas Ativas**: 19 (100%)
- **Rotas Inativas**: 0 (0%)
- **Última Atualização**: 2025-10-06

## Métodos whatsmeow Necessários

### Métodos Principais de Envio

| Método | Descrição | Uso |
|--------|-----------|-----|
| `SendMessage(ctx, to, message, extra)` | Método principal para envio | Todos os tipos de mensagem |
| `Upload(ctx, plaintext, appInfo)` | Upload de mídia | Imagens, áudios, vídeos, documentos |
| `BuildReaction(chat, sender, id, reaction)` | Construir reação | Reações a mensagens |
| `BuildPollCreation(name, options, selectableCount)` | Construir enquete | Enquetes |
| `BuildPollVote(ctx, pollInfo, optionNames)` | Votar em enquete | Votos em enquetes |
| `GenerateMessageID()` | Gerar ID de mensagem | Todos os envios |

### Tipos de Mensagem (waE2E.Message)

| Tipo | Campo | Descrição |
|------|-------|-----------|
| Texto | `Conversation` | Mensagem de texto simples |
| Texto Estendido | `ExtendedTextMessage` | Texto com formatação/links |
| Imagem | `ImageMessage` | Mensagem de imagem |
| Áudio | `AudioMessage` | Mensagem de áudio/voice note |
| Vídeo | `VideoMessage` | Mensagem de vídeo |
| Documento | `DocumentMessage` | Arquivo/documento |
| Sticker | `StickerMessage` | Sticker/adesivo |
| Localização | `LocationMessage` | Coordenadas GPS |
| Contato | `ContactMessage` | Cartão de contato |
| Contatos | `ContactsArrayMessage` | Múltiplos contatos |
| Template | `TemplateMessage` | Mensagem template |
| Lista | `ListMessage` | Lista interativa |
| Botões | `ButtonsMessage` | Botões interativos |
| Enquete | `PollCreationMessage` | Enquete/poll |

## Estruturas de Dados Necessárias

### DTOs de Request

```go
// Texto
type SendTextMessageRequest struct {
    To   string `json:"to"`
    Text string `json:"text"`
}

// Imagem
type SendImageMessageRequest struct {
    To      string     `json:"to"`
    Image   *MediaData `json:"image"`
    Caption string     `json:"caption,omitempty"`
}

// Áudio
type SendAudioMessageRequest struct {
    To    string     `json:"to"`
    Audio *MediaData `json:"audio"`
}

// Localização
type SendLocationMessageRequest struct {
    To        string  `json:"to"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Name      string  `json:"name,omitempty"`
    Address   string  `json:"address,omitempty"`
}

// Contato
type SendContactMessageRequest struct {
    To      string       `json:"to"`
    Contact *ContactInfo `json:"contact"`
}

// Reação
type SendReactionMessageRequest struct {
    To        string `json:"to"`
    MessageID string `json:"messageId"`
    Reaction  string `json:"reaction"`
}

// Enquete
type SendPollMessageRequest struct {
    To                     string   `json:"to"`
    Name                   string   `json:"name"`
    Options                []string `json:"options"`
    SelectableOptionsCount int      `json:"selectableOptionsCount"`
}
```

### MediaData Structure

```go
type MediaData struct {
    Data     string `json:"data"`     // Base64 encoded data
    MimeType string `json:"mimeType"`
    FileName string `json:"fileName,omitempty"`
}
```

## Exemplos de Uso

### Mensagem de Texto
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/text \
  -H "Content-Type: application/json" \
  -d '{
    "to": "5511999999999@s.whatsapp.net",
    "text": "Olá, mundo!"
  }'
```

### Mensagem de Imagem
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/image \
  -H "Content-Type: application/json" \
  -d '{
    "to": "5511999999999@s.whatsapp.net",
    "image": {
      "data": "base64_encoded_image_data",
      "mimeType": "image/jpeg",
      "fileName": "foto.jpg"
    },
    "caption": "Minha foto"
  }'
```

### Reação a Mensagem
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/reaction \
  -H "Content-Type: application/json" \
  -d '{
    "to": "5511999999999@s.whatsapp.net",
    "messageId": "3EB0XXXXXXXXXXXXX",
    "reaction": "👍"
  }'
```

### Enquete (Poll)
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/poll \
  -H "Content-Type: application/json" \
  -d '{
    "to": "5511999999999@s.whatsapp.net",
    "name": "Qual sua cor favorita?",
    "options": ["Azul", "Verde", "Vermelho", "Amarelo"],
    "selectableOptionsCount": 1
  }'
```

### Mensagem com Botões
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/buttons \
  -H "Content-Type: application/json" \
  -d '{
    "to": "5511999999999@s.whatsapp.net",
    "text": "Escolha uma opção:",
    "buttons": [
      {"id": "1", "text": "Opção 1"},
      {"id": "2", "text": "Opção 2"},
      {"id": "3", "text": "Opção 3"}
    ]
  }'
```

### Lista Interativa
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/list \
  -H "Content-Type: application/json" \
  -d '{
    "to": "5511999999999@s.whatsapp.net",
    "text": "Escolha um produto",
    "title": "Nossos Produtos",
    "sections": [
      {
        "title": "Eletrônicos",
        "rows": [
          {"id": "1", "title": "Smartphone", "description": "R$ 1.000"},
          {"id": "2", "title": "Notebook", "description": "R$ 3.000"}
        ]
      }
    ]
  }'
```

### View Once (Visualização Única)
```bash
curl -X POST http://localhost:8080/sessions/my-session/send/message/viewonce \
  -H "Content-Type: application/json" \
  -d '{
    "to": "5511999999999@s.whatsapp.net",
    "media": {
      "data": "base64_encoded_image_data",
      "mimeType": "image/jpeg"
    },
    "caption": "Esta imagem só pode ser vista uma vez"
  }'
```

## Status da Implementação

### ✅ Completo
- [x] Análise profunda da biblioteca whatsmeow
- [x] Implementação de todos os métodos no waclient
- [x] Implementação de todos os handlers HTTP
- [x] Validação de DTOs
- [x] Documentação atualizada
- [x] Compilação bem-sucedida

### 🚧 Próximos Passos Recomendados
1. **Testes Unitários** - Criar testes para cada handler
2. **Testes de Integração** - Testar fluxo completo de envio
3. **Swagger/OpenAPI** - Gerar documentação automática da API
4. **Tratamento de Erros Avançado** - Melhorar mensagens de erro
5. **Rate Limiting** - Implementar controle de taxa de envio
6. **Webhooks** - Sistema de notificações de eventos

## Referências

- [whatsmeow Documentation](https://pkg.go.dev/go.mau.fi/whatsmeow)
- [WuzAPI Example](https://github.com/asternic/wuzapi)
- Código existente em `internal/adapters/waclient/`
