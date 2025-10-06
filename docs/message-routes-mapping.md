# Mapeamento das Rotas de Mensagens - WhatsApp API

## Status Atual das Rotas

### ✅ Rotas Implementadas (Ativas)
| Rota | Método | Handler | Status |
|------|--------|---------|--------|
| `POST /sessions/{sessionId}/messages` | SendMessage | h.Message.SendMessage | ✅ Ativo |
| `GET /sessions/{sessionId}/chats` | GetChats | h.Message.GetChats | ✅ Ativo |
| `GET /sessions/{sessionId}/contacts` | GetContacts | h.Message.GetContacts | ✅ Ativo |
| `GET /sessions/{sessionId}/chat-info` | GetChatInfo | h.Message.GetChatInfo | ✅ Ativo |

### ⚠️ Rotas Definidas mas Inativas
| Rota | Método | Handler | Status |
|------|--------|---------|--------|
| `POST /sessions/{sessionId}/send/message/text` | SendText | h.Message.SendText | ⚠️ Definida mas não registrada |
| `POST /sessions/{sessionId}/send/message/image` | SendImage | h.Message.SendImage | ⚠️ Definida mas não registrada |
| `POST /sessions/{sessionId}/send/message/audio` | SendAudio | h.Message.SendAudio | ⚠️ Definida mas não registrada |
| `POST /sessions/{sessionId}/send/message/video` | SendVideo | h.Message.SendVideo | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/document` | SendDocument | h.Message.SendDocument | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/sticker` | SendSticker | h.Message.SendSticker | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/location` | SendLocation | h.Message.SendLocation | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/contact` | SendContact | h.Message.SendContact | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/contacts` | SendContactsArray | h.Message.SendContactsArray | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/reaction` | SendReaction | h.Message.SendReaction | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/template` | SendTemplate | h.Message.SendTemplate | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/buttons` | SendButtons | h.Message.SendButtons | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/list` | SendList | h.Message.SendList | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/poll` | SendPoll | h.Message.SendPoll | ❌ Handler não implementado |
| `POST /sessions/{sessionId}/send/message/viewonce` | SendViewOnce | h.Message.SendViewOnce | ❌ Handler não implementado |

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

## Próximos Passos

1. **Implementar handlers básicos** - Começar com texto, imagem, áudio
2. **Criar DTOs de request/response** - Estruturas de dados padronizadas
3. **Implementar upload de mídia** - Sistema de upload para arquivos
4. **Adicionar validações** - Validar dados de entrada
5. **Implementar handlers avançados** - Reações, enquetes, listas
6. **Testes** - Criar testes unitários e de integração
7. **Documentação** - Swagger/OpenAPI specs

## Referências

- [whatsmeow Documentation](https://pkg.go.dev/go.mau.fi/whatsmeow)
- [WuzAPI Example](https://github.com/asternic/wuzapi)
- Código existente em `internal/adapters/waclient/`
