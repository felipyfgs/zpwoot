# Análise Profunda da Biblioteca whatsmeow

## 📚 Documentação Oficial
- **Repositório**: https://github.com/tulir/whatsmeow
- **Go Docs**: https://pkg.go.dev/go.mau.fi/whatsmeow
- **Licença**: MPL-2.0

---

## 🔑 Conceitos Fundamentais

### 1. Client.GenerateMessageID()

**Assinatura:**
```go
func (cli *Client) GenerateMessageID() types.MessageID
```

**Descrição:**
- Gera um ID único para mensagens
- Retorna `types.MessageID` (que é um `string`)
- **IMPORTANTE**: O whatsmeow gera automaticamente um ID se não for fornecido
- Usado opcionalmente em `SendRequestExtra.ID`

**Quando usar:**
- Quando você precisa rastrear a mensagem antes de enviá-la
- Para deduplicação de mensagens
- Para correlacionar request/response

**Exemplo do WuzAPI:**
```go
if t.Id == "" {
    msgid = clientManager.GetWhatsmeowClient(txtid).GenerateMessageID()
} else {
    msgid = t.Id
}
```

---

### 2. Client.SendMessage()

**Assinatura:**
```go
func (cli *Client) SendMessage(
    ctx context.Context,
    to types.JID,
    message *waE2E.Message,
    extra ...SendRequestExtra
) (SendResponse, error)
```

**Parâmetros:**
- `ctx`: Context para cancelamento/timeout
- `to`: JID do destinatário (types.JID)
- `message`: Mensagem protobuf (*waE2E.Message)
- `extra`: Parâmetros opcionais (SendRequestExtra)

**Retorno: SendResponse**
```go
type SendResponse struct {
    Timestamp     time.Time           // Timestamp do servidor
    ID            types.MessageID     // ID da mensagem enviada
    ServerID      types.MessageServerID // ID do servidor (newsletters)
    DebugTimings  MessageDebugTimings // Métricas de debug
    Sender        types.JID           // JID do remetente (LID ou PN)
}
```

**IMPORTANTE:**
- O `SendResponse.ID` é o ID REAL da mensagem retornado pelo WhatsApp
- Este é o ID que deve ser retornado ao cliente da API
- O `SendResponse.Timestamp` é o timestamp oficial do servidor

---

### 3. SendRequestExtra (Parâmetros Opcionais)

```go
type SendRequestExtra struct {
    ID            types.MessageID    // ID customizado (opcional)
    InlineBotJID  types.JID         // Bot JID (opcional)
    Peer          bool              // Mensagem peer (para próprios devices)
    Timeout       time.Duration     // Timeout (padrão: 75s)
    MediaHandle   string            // Handle de mídia (newsletters)
    Meta          *types.MsgMetaInfo // Metadados
    AdditionalNodes *[]waBinary.Node // Nós adicionais (avançado)
}
```

**Uso:**
```go
// Sem parâmetros extras
resp, err := client.SendMessage(ctx, to, message)

// Com ID customizado
resp, err := client.SendMessage(ctx, to, message, whatsmeow.SendRequestExtra{
    ID: client.GenerateMessageID(),
})
```

---

### 4. Client.Upload()

**Assinatura:**
```go
func (cli *Client) Upload(
    ctx context.Context,
    plaintext []byte,
    appInfo MediaType
) (UploadResponse, error)
```

**Parâmetros:**
- `plaintext`: Dados binários do arquivo ([]byte)
- `appInfo`: Tipo de mídia (MediaType)

**MediaType Constants:**
```go
const (
    MediaImage         MediaType = "WhatsApp Image Keys"
    MediaVideo         MediaType = "WhatsApp Video Keys"
    MediaAudio         MediaType = "WhatsApp Audio Keys"
    MediaDocument      MediaType = "WhatsApp Document Keys"
    MediaHistory       MediaType = "WhatsApp History Keys"
    MediaAppState      MediaType = "WhatsApp App State Keys"
    MediaLinkThumbnail MediaType = "WhatsApp Link Thumbnail Keys"
)
```

**Retorno: UploadResponse**
```go
type UploadResponse struct {
    URL           string  // URL do arquivo
    DirectPath    string  // Path direto
    Handle        string  // Handle (newsletters)
    ObjectID      string  // Object ID
    MediaKey      []byte  // Chave de criptografia
    FileEncSHA256 []byte  // SHA256 do arquivo criptografado
    FileSHA256    []byte  // SHA256 do arquivo original
    FileLength    uint64  // Tamanho do arquivo
}
```

**Fluxo de Upload:**
1. Receber dados em base64 do cliente
2. Decodificar para []byte
3. Chamar `client.Upload(ctx, fileData, mediaType)`
4. Usar `UploadResponse` para construir a mensagem protobuf

---

### 5. Client.BuildReaction()

**Assinatura:**
```go
func (cli *Client) BuildReaction(
    chat types.JID,
    sender types.JID,
    id types.MessageID,
    reaction string
) *waE2E.Message
```

**Parâmetros:**
- `chat`: JID do chat (grupo ou contato)
- `sender`: JID do remetente da mensagem original
- `id`: ID da mensagem para reagir
- `reaction`: Emoji da reação (string vazia para remover)

**Retorno:**
- Retorna `*waE2E.Message` pronta para enviar

**Uso:**
```go
reactionMsg := client.BuildReaction(chatJID, senderJID, messageID, "👍")
resp, err := client.SendMessage(ctx, chatJID, reactionMsg)
```

---

### 6. Client.BuildPollCreation()

**Assinatura:**
```go
func (cli *Client) BuildPollCreation(
    name string,
    optionNames []string,
    selectableOptionCount int
) *waE2E.Message
```

**Parâmetros:**
- `name`: Título da enquete
- `optionNames`: Array de opções (2-12 opções)
- `selectableOptionCount`: Quantas opções podem ser selecionadas

**Retorno:**
- Retorna `*waE2E.Message` pronta para enviar

**Uso:**
```go
pollMsg := client.BuildPollCreation(
    "Qual sua cor favorita?",
    []string{"Azul", "Verde", "Vermelho"},
    1, // Apenas uma opção
)
resp, err := client.SendMessage(ctx, recipientJID, pollMsg)
```

---

## 📦 Estruturas de Mensagem waE2E.Message

### Mensagem de Texto Simples
```go
message := &waE2E.Message{
    Conversation: proto.String("Olá, mundo!"),
}
```

### Mensagem de Imagem
```go
message := &waE2E.Message{
    ImageMessage: &waE2E.ImageMessage{
        URL:           proto.String(uploaded.URL),
        DirectPath:    proto.String(uploaded.DirectPath),
        MediaKey:      uploaded.MediaKey,
        Mimetype:      proto.String("image/jpeg"),
        FileEncSHA256: uploaded.FileEncSHA256,
        FileSHA256:    uploaded.FileSHA256,
        FileLength:    proto.Uint64(fileLength),
        Caption:       proto.String("Legenda da imagem"),
        ViewOnce:      proto.Bool(false), // View once
    },
}
```

### Mensagem de Áudio
```go
ptt := true // Push-to-talk (voice note)
message := &waE2E.Message{
    AudioMessage: &waE2E.AudioMessage{
        URL:           proto.String(uploaded.URL),
        DirectPath:    proto.String(uploaded.DirectPath),
        MediaKey:      uploaded.MediaKey,
        Mimetype:      proto.String("audio/ogg; codecs=opus"),
        FileEncSHA256: uploaded.FileEncSHA256,
        FileSHA256:    uploaded.FileSHA256,
        FileLength:    proto.Uint64(fileLength),
        PTT:           &ptt,
    },
}
```

### Mensagem de Vídeo
```go
message := &waE2E.Message{
    VideoMessage: &waE2E.VideoMessage{
        URL:           proto.String(uploaded.URL),
        DirectPath:    proto.String(uploaded.DirectPath),
        MediaKey:      uploaded.MediaKey,
        Mimetype:      proto.String("video/mp4"),
        FileEncSHA256: uploaded.FileEncSHA256,
        FileSHA256:    uploaded.FileSHA256,
        FileLength:    proto.Uint64(fileLength),
        Caption:       proto.String("Legenda do vídeo"),
        ViewOnce:      proto.Bool(false),
    },
}
```

### Mensagem de Documento
```go
message := &waE2E.Message{
    DocumentMessage: &waE2E.DocumentMessage{
        URL:           proto.String(uploaded.URL),
        DirectPath:    proto.String(uploaded.DirectPath),
        MediaKey:      uploaded.MediaKey,
        Mimetype:      proto.String("application/pdf"),
        FileEncSHA256: uploaded.FileEncSHA256,
        FileSHA256:    uploaded.FileSHA256,
        FileLength:    proto.Uint64(fileLength),
        FileName:      proto.String("documento.pdf"),
        Caption:       proto.String("Legenda do documento"),
    },
}
```

---

## ✅ Padrão Correto de Implementação

### Fluxo Completo de Envio de Mensagem

```go
func SendMessage(ctx context.Context, sessionID, to string, ...) error {
    // 1. Obter cliente conectado
    client, err := getConnectedClient(ctx, sessionID)
    if err != nil {
        return err
    }

    // 2. Parse JID do destinatário
    recipientJID, err := parseJID(to)
    if err != nil {
        return ErrInvalidJID
    }

    // 3. Construir mensagem protobuf
    message := &waE2E.Message{
        // ... campos da mensagem
    }

    // 4. Enviar mensagem (whatsmeow gera ID automaticamente)
    resp, err := client.WAClient.SendMessage(ctx, recipientJID, message)
    if err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }

    // 5. Retornar dados do SendResponse
    // resp.ID - ID real da mensagem
    // resp.Timestamp - Timestamp do servidor
    // resp.Sender - JID do remetente
    
    return nil
}
```

---

## 🎯 Conclusões Importantes

1. **NÃO gerar IDs manualmente** - O whatsmeow faz isso automaticamente
2. **Usar SendResponse.ID** - Este é o ID real retornado pelo WhatsApp
3. **Usar SendResponse.Timestamp** - Timestamp oficial do servidor
4. **Upload antes de enviar** - Mídia deve ser enviada via Upload() primeiro
5. **Usar proto.String(), proto.Uint64(), etc** - Para campos protobuf
6. **BuildReaction e BuildPollCreation** - Métodos helper do whatsmeow
7. **MediaType correto** - Usar constantes do whatsmeow para upload

---

## 📖 Exemplos Reais do WuzAPI

### SendImage do WuzAPI
```go
func SendImage(w http.ResponseWriter, r *http.Request) {
    txtid := chi.URLParam(r, "phone")
    var t SendImageRequest
    json.NewDecoder(r.Body).Decode(&t)

    // Decode base64
    rawDecodedText, _ := base64.StdEncoding.DecodeString(t.Image)

    // Upload
    uploaded, err := clientManager.GetWhatsmeowClient(txtid).Upload(
        context.Background(),
        rawDecodedText,
        whatsmeow.MediaImage,
    )

    // Build message
    msg := &waE2E.Message{
        ImageMessage: &waE2E.ImageMessage{
            Caption:       proto.String(t.Caption),
            URL:           proto.String(uploaded.URL),
            DirectPath:    proto.String(uploaded.DirectPath),
            MediaKey:      uploaded.MediaKey,
            Mimetype:      proto.String(http.DetectContentType(rawDecodedText)),
            FileEncSHA256: uploaded.FileEncSHA256,
            FileSHA256:    uploaded.FileSHA256,
            FileLength:    proto.Uint64(uint64(len(rawDecodedText))),
        },
    }

    // Send
    resp, err := clientManager.GetWhatsmeowClient(txtid).SendMessage(
        context.Background(),
        recipient,
        msg,
    )

    // Return response with real ID
    json.NewEncoder(w).Encode(map[string]interface{}{
        "Details": "Sent",
        "Timestamp": resp.Timestamp,
        "Id": resp.ID,
    })
}
```

### SendAudio do WuzAPI
```go
func SendAudio(w http.ResponseWriter, r *http.Request) {
    // ... decode base64 audio ...

    uploaded, err := clientManager.GetWhatsmeowClient(txtid).Upload(
        context.Background(),
        rawDecodedText,
        whatsmeow.MediaAudio,
    )

    msg := &waE2E.Message{
        AudioMessage: &waE2E.AudioMessage{
            URL:           proto.String(uploaded.URL),
            DirectPath:    proto.String(uploaded.DirectPath),
            MediaKey:      uploaded.MediaKey,
            Mimetype:      proto.String("audio/ogg; codecs=opus"),
            FileEncSHA256: uploaded.FileEncSHA256,
            FileSHA256:    uploaded.FileSHA256,
            FileLength:    proto.Uint64(uint64(len(rawDecodedText))),
            PTT:           proto.Bool(true), // Voice note
        },
    }

    resp, err := clientManager.GetWhatsmeowClient(txtid).SendMessage(
        context.Background(),
        recipient,
        msg,
    )
}
```

### SendPoll do WuzAPI
```go
func SendPoll(w http.ResponseWriter, r *http.Request) {
    txtid := chi.URLParam(r, "phone")
    var t SendPollRequest
    json.NewDecoder(r.Body).Decode(&t)

    recipient, _ := types.ParseJID(t.Phone)

    // Use BuildPollCreation
    msg := clientManager.GetWhatsmeowClient(txtid).BuildPollCreation(
        t.Name,
        t.Options,
        t.Selectable,
    )

    resp, err := clientManager.GetWhatsmeowClient(txtid).SendMessage(
        context.Background(),
        recipient,
        msg,
    )

    json.NewEncoder(w).Encode(map[string]interface{}{
        "Details": "Sent",
        "Timestamp": resp.Timestamp,
        "Id": resp.ID,
    })
}
```

---

## 📝 Próximos Passos

1. ✅ Análise profunda da biblioteca whatsmeow completa
2. ✅ Documentação de conceitos fundamentais criada
3. ⏭️ Revisar código implementado na Fase 2
4. ⏭️ Corrigir métodos para retornar SendResponse corretamente
5. ⏭️ Implementar handlers HTTP com integração adequada
6. ⏭️ Testar compilação e funcionalidade

