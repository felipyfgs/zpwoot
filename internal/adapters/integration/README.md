# Integration Adapters - zpwoot

## 📋 Visão Geral

Este diretório (`internal/adapters/integration/`) contém **adapters de integração com sistemas externos**, seguindo a Clean Architecture do projeto zpwoot.

## 🎯 Propósito

Centralizar a lógica de integração com serviços externos como:
- **Webhooks** - Envio de eventos para URLs externas
- **Chatwoot** - Integração com plataforma de atendimento (futuro)
- **Outras integrações** - Qualquer serviço externo que precise se comunicar com o zpwoot

## 🏗️ Estrutura Atual

```
internal/adapters/integration/
├── README.md                           # Este arquivo
├── WEBHOOK_INTEGRATION_STUDY.md        # Estudo completo de webhooks
├── WEBHOOK_ARCHITECTURE_DIAGRAM.md     # Diagramas de arquitetura
└── webhook/                            # Adapter de Webhook
    ├── sender.go                       # Implementação do envio HTTP
    ├── signature.go                    # Geração de assinaturas HMAC
    └── validator.go                    # Validação de URLs
```

## 🔄 Como Funciona (Clean Architecture)

### Camadas e Responsabilidades

```
┌─────────────────────────────────────────────────────────────┐
│  DOMAIN (internal/core/domain/webhook/)                     │
│  - Entidade Webhook                                         │
│  - Regras de negócio puras                                  │
│  - Interface Repository                                     │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │ usa
┌─────────────────────────────────────────────────────────────┐
│  APPLICATION (internal/core/application/usecase/webhook/)   │
│  - Use Cases (CreateOrUpdate, Get, Delete, ListEvents)      │
│  - Orquestração de lógica de negócio                        │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │ usa
┌─────────────────────────────────────────────────────────────┐
│  PORTS (internal/core/ports/)                               │
│  - input/webhook.go: Interface WebhookUseCases              │
│  - output/webhook_sender.go: Interface WebhookSender        │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │ implementa
┌─────────────────────────────────────────────────────────────┐
│  ADAPTERS (internal/adapters/)                              │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  integration/webhook/                                 │  │
│  │  - sender.go: Implementa WebhookSender                │  │
│  │  - Envia HTTP POST com retry                          │  │
│  │  - Gera assinaturas HMAC                              │  │
│  │  - Valida URLs                                        │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  database/repository/webhook.go                       │  │
│  │  - Implementa domain/webhook/Repository               │  │
│  │  - CRUD no PostgreSQL                                 │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  http/handlers/webhook.go                             │  │
│  │  - Usa ports/input/WebhookUseCases                    │  │
│  │  - Endpoints REST                                     │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 📦 Exemplo: Webhook Integration

### 1. Domain Layer
```go
// internal/core/domain/webhook/entity.go
type Webhook struct {
    ID        string
    SessionID string
    URL       string
    Secret    *string
    Events    []string
    Enabled   bool
}

func (w *Webhook) HasEvent(eventType string) bool {
    // Regra de negócio
}
```

### 2. Port Layer
```go
// internal/core/ports/output/webhook_sender.go
type WebhookSender interface {
    SendWebhook(ctx context.Context, url, secret string, event *WebhookEvent) error
}
```

### 3. Integration Adapter (AQUI!)
```go
// internal/adapters/integration/webhook/sender.go
type HTTPWebhookSender struct {
    httpClient *http.Client
    logger     *logger.Logger
}

func (s *HTTPWebhookSender) SendWebhook(ctx context.Context, url, secret string, event *WebhookEvent) error {
    // 1. Serializa evento
    payload, _ := json.Marshal(event)
    
    // 2. Gera assinatura HMAC
    signature := generateHMAC(payload, secret)
    
    // 3. Prepara request
    req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
    req.Header.Set("X-Webhook-Signature", signature)
    req.Header.Set("X-Webhook-Event", event.Type)
    
    // 4. Envia com retry
    return s.sendWithRetry(req)
}
```

## 🔌 Integrações Planejadas

### Webhook (Em Implementação)
- ✅ Estudo completo realizado
- ⏳ Implementação pendente
- 📍 Localização: `internal/adapters/integration/webhook/`

### Chatwoot (Futuro)
- 📋 Integração bidirecional com Chatwoot
- 📍 Localização futura: `internal/adapters/integration/chatwoot/`
- Estrutura proposta:
  ```
  internal/adapters/integration/chatwoot/
  ├── client.go           # Cliente HTTP para API do Chatwoot
  ├── message_sync.go     # Sincronização de mensagens
  ├── contact_sync.go     # Sincronização de contatos
  └── webhook_handler.go  # Recebe webhooks do Chatwoot
  ```

### Outras Integrações
- **N8N** - Automação de workflows
- **Make/Zapier** - Integrações no-code
- **Custom APIs** - APIs customizadas dos clientes

## 📐 Regras de Arquitetura

### ✅ O que os Integration Adapters PODEM fazer:
1. **Implementar interfaces** definidas em `ports/output/`
2. **Comunicar com APIs externas** (HTTP, gRPC, etc)
3. **Usar bibliotecas externas** (http client, SDKs)
4. **Fazer retry e error handling** de comunicação
5. **Logging e métricas** de integrações
6. **Transformar dados** entre formatos internos e externos

### ❌ O que os Integration Adapters NÃO PODEM fazer:
1. **Conter regras de negócio** (isso fica no domain!)
2. **Acessar banco de dados diretamente** (usar repository!)
3. **Ser chamados diretamente** por handlers HTTP (usar ports!)
4. **Depender de outros adapters** (usar interfaces!)

## 🔄 Fluxo de Dados Típico

### Envio de Evento via Webhook
```
WhatsApp Event
    ↓
EventHandler (waclient/events.go)
    ↓
WebhookSender Interface (ports/output/)
    ↓
HTTPWebhookSender (integration/webhook/sender.go) ← AQUI!
    ↓
HTTP POST → External API
```

### Configuração de Webhook
```
HTTP Request
    ↓
WebhookHandler (http/handlers/webhook.go)
    ↓
WebhookUseCases Interface (ports/input/)
    ↓
CreateOrUpdateUseCase (application/usecase/webhook/)
    ↓
Webhook Domain Service (domain/webhook/)
    ↓
Repository Interface (domain/webhook/)
    ↓
PostgresRepository (database/repository/webhook.go)
```

## 🧪 Testes

### Testes de Integration Adapters
```go
// internal/adapters/integration/webhook/sender_test.go
func TestHTTPWebhookSender_SendWebhook(t *testing.T) {
    // Mock HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verifica headers
        assert.NotEmpty(t, r.Header.Get("X-Webhook-Signature"))
        
        // Verifica payload
        var event WebhookEvent
        json.NewDecoder(r.Body).Decode(&event)
        assert.Equal(t, "Message", event.Type)
        
        w.WriteHeader(http.StatusOK)
    }))
    defer server.Close()
    
    sender := NewHTTPWebhookSender(http.DefaultClient, logger)
    err := sender.SendWebhook(ctx, server.URL, "secret", &event)
    
    assert.NoError(t, err)
}
```

## 📚 Documentação Adicional

- **[WEBHOOK_INTEGRATION_STUDY.md](./WEBHOOK_INTEGRATION_STUDY.md)** - Estudo completo de implementação
- **[WEBHOOK_ARCHITECTURE_DIAGRAM.md](./WEBHOOK_ARCHITECTURE_DIAGRAM.md)** - Diagramas visuais
- **[ARCHITECTURE.md](../../../ARCHITECTURE.md)** - Arquitetura geral do projeto

## 🚀 Próximos Passos

1. ✅ Estudo de arquitetura completo
2. ⏳ Implementar domain layer (webhook entity, repository interface)
3. ⏳ Implementar application layer (use cases)
4. ⏳ Implementar ports (input/output interfaces)
5. ⏳ Implementar integration adapter (sender, signature, validator)
6. ⏳ Implementar database repository
7. ⏳ Implementar HTTP handlers
8. ⏳ Integrar com EventHandler existente
9. ⏳ Testes completos
10. ⏳ Documentação de API

## 💡 Dicas para Novos Desenvolvedores

1. **Sempre comece pelo Domain** - Defina entidades e regras de negócio primeiro
2. **Use interfaces (Ports)** - Nunca dependa diretamente de implementações
3. **Mantenha adapters simples** - Apenas tradução entre sistemas
4. **Teste isoladamente** - Use mocks para testar sem dependências externas
5. **Siga o fluxo de dependências** - Domain ← Application ← Ports ← Adapters

## 🤝 Contribuindo

Ao adicionar novas integrações:

1. Crie um subdiretório em `internal/adapters/integration/`
2. Defina a interface em `internal/core/ports/output/`
3. Implemente a interface no adapter
4. Adicione testes
5. Documente no README
6. Atualize o container para injeção de dependências

---

**Mantido por**: Equipe zpwoot  
**Última atualização**: 2025-10-06

