# Melhorias na Documentação Swagger

Este documento descreve as melhorias implementadas na documentação Swagger da API zpwoot.

## 📋 Resumo das Alterações

### ✅ Problemas Corrigidos

1. **Nomes de Modelos Limpos**
   - ❌ Antes: `zpwoot_internal_core_application_dto.SendTextMessageRequest`
   - ✅ Depois: `SendTextMessageRequest`
   
2. **Exemplos Adicionados em Todos os DTOs**
   - Todos os campos agora possuem exemplos práticos
   - Valores realistas e úteis para testes

3. **Documentação Consistente**
   - Todos os modelos seguem o mesmo padrão
   - Anotações `@name` adicionadas para controlar nomes no Swagger

---

## 🔧 Mudanças Técnicas

### DTOs de Mensagem (`internal/core/application/dto/message.go`)

Adicionadas anotações `@name` e tags `example` em todos os tipos:

```go
type SendTextMessageRequest struct {
    To   string `json:"to" validate:"required" example:"5511999999999"`
    Text string `json:"text" validate:"required" example:"Hello! This is a test message from zpwoot API."`
} //@name SendTextMessageRequest
```

#### Tipos Atualizados:

- ✅ `MediaData` - Exemplos de URL e Base64
- ✅ `Location` - Coordenadas de São Paulo
- ✅ `ContactInfo` - Informações de contato completas
- ✅ `SendMessageResponse` - Resposta padrão de envio
- ✅ `SendTextMessageRequest` - Mensagem de texto
- ✅ `SendImageMessageRequest` - Mensagem de imagem
- ✅ `SendAudioMessageRequest` - Mensagem de áudio
- ✅ `SendVideoMessageRequest` - Mensagem de vídeo
- ✅ `SendDocumentMessageRequest` - Mensagem de documento
- ✅ `SendStickerMessageRequest` - Mensagem de sticker
- ✅ `SendLocationMessageRequest` - Mensagem de localização
- ✅ `SendContactMessageRequest` - Mensagem de contato
- ✅ `SendContactsArrayMessageRequest` - Múltiplos contatos
- ✅ `SendReactionMessageRequest` - Reação com emoji
- ✅ `SendPollMessageRequest` - Enquete
- ✅ `SendButtonsMessageRequest` - Botões interativos
- ✅ `SendListMessageRequest` - Lista de opções
- ✅ `SendTemplateMessageRequest` - Template do WhatsApp Business
- ✅ `SendViewOnceMessageRequest` - Mensagem que desaparece
- ✅ `Button` - Botão individual
- ✅ `ListRow` - Linha de lista
- ✅ `ListSection` - Seção de lista
- ✅ `TemplateParameter` - Parâmetro de template
- ✅ `TemplateComponent` - Componente de template
- ✅ `TemplateMessage` - Mensagem de template

### DTOs de Health (`internal/adapters/http/handlers/health.go`)

```go
type HealthResponse struct {
    Status  string `json:"status" example:"ok"`
    Service string `json:"service" example:"zpwoot"`
    Version string `json:"version,omitempty" example:"1.0.0"`
} //@name HealthResponse

type InfoResponse struct {
    Message string `json:"message" example:"zpwoot WhatsApp API is running"`
    Version string `json:"version" example:"1.0.0"`
    Service string `json:"service" example:"zpwoot"`
} //@name InfoResponse
```

---

## 📊 Exemplos de Modelos no Swagger

### SendTextMessageRequest

```json
{
  "to": "5511999999999",
  "text": "Hello! This is a test message from zpwoot API."
}
```

### SendImageMessageRequest

```json
{
  "to": "5511999999999",
  "image": {
    "url": "https://example.com/image.jpg",
    "mimeType": "image/jpeg",
    "fileName": "image.jpg",
    "caption": "Check out this image!"
  },
  "caption": "Check out this beautiful image!"
}
```

### MediaData

```json
{
  "url": "https://example.com/image.jpg",
  "base64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
  "fileName": "image.jpg",
  "mimeType": "image/jpeg",
  "caption": "Check out this image!"
}
```

### SendLocationMessageRequest

```json
{
  "to": "5511999999999",
  "latitude": -23.550520,
  "longitude": -46.633308,
  "name": "São Paulo",
  "address": "Av. Paulista, 1578 - Bela Vista, São Paulo - SP"
}
```

### SendPollMessageRequest

```json
{
  "to": "5511999999999",
  "name": "What's your favorite color?",
  "options": ["Red", "Blue", "Green", "Yellow"],
  "selectableOptionsCount": 1
}
```

### SendButtonsMessageRequest

```json
{
  "to": "5511999999999",
  "text": "Please choose an option:",
  "buttons": [
    {
      "id": "btn_1",
      "text": "Click Me"
    }
  ]
}
```

### SendListMessageRequest

```json
{
  "to": "5511999999999",
  "text": "Please select an option from the list",
  "title": "Menu Options",
  "sections": [
    {
      "title": "Section 1",
      "rows": [
        {
          "id": "row_1",
          "title": "Option 1",
          "description": "Description for option 1"
        }
      ]
    }
  ]
}
```

### SendMessageResponse

```json
{
  "messageId": "msg_123456789",
  "status": "sent",
  "sentAt": "2024-01-15T10:30:00Z"
}
```

---

## 🎯 Benefícios

1. **Melhor Experiência do Desenvolvedor**
   - Exemplos práticos e realistas
   - Fácil de testar diretamente no Swagger UI
   - Valores de exemplo prontos para copiar e colar

2. **Documentação Mais Clara**
   - Nomes de modelos limpos e legíveis
   - Estrutura consistente em todos os endpoints
   - Fácil navegação entre modelos relacionados

3. **Redução de Erros**
   - Desenvolvedores veem exatamente o formato esperado
   - Exemplos mostram valores válidos
   - Menos tentativa e erro ao integrar

4. **Facilita Testes**
   - Botão "Try it out" no Swagger já vem preenchido
   - Valores de exemplo funcionam imediatamente
   - Testes rápidos sem precisar criar dados manualmente

---

## 🔗 Acessando a Documentação

### Swagger UI
```
http://localhost:8080/swagger/index.html
```

### Swagger JSON
```
http://localhost:8080/swagger/doc.json
```

### Exemplos de Uso
Veja o arquivo [API_EXAMPLES.md](./API_EXAMPLES.md) para exemplos completos de uso via cURL.

---

## 📝 Como Adicionar Exemplos em Novos DTOs

Ao criar novos DTOs, siga este padrão:

```go
type MyNewRequest struct {
    Field1 string `json:"field1" validate:"required" example:"example_value"`
    Field2 int    `json:"field2" example:"123"`
    Field3 bool   `json:"field3,omitempty" example:"true"`
} //@name MyNewRequest
```

**Importante:**
- Sempre adicione a tag `example` com valores realistas
- Sempre adicione a anotação `//@name` no final do tipo
- Use valores que façam sentido no contexto da aplicação
- Para números de telefone, use formato brasileiro: `5511999999999`
- Para datas, use formato ISO 8601: `2024-01-15T10:30:00Z`

---

## 🔄 Regenerando o Swagger

Após fazer alterações nos DTOs ou handlers:

```bash
make swagger
```

Ou para regenerar e iniciar o servidor:

```bash
make swagger-serve
```

---

## ✅ Checklist de Qualidade

- [x] Todos os DTOs têm anotação `@name`
- [x] Todos os campos têm tag `example`
- [x] Exemplos são realistas e úteis
- [x] Nomes de modelos estão limpos (sem prefixo de diretório)
- [x] Documentação Swagger gerada e testada
- [x] Exemplos de uso documentados em API_EXAMPLES.md
- [x] Todos os endpoints de mensagem documentados
- [x] Respostas de erro documentadas

---

**zpwoot** - Making WhatsApp Business API integration simple and powerful! 🚀

