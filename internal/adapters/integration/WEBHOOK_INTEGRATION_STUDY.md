# Estudo de Implementação de Webhooks - zpwoot

## 📋 Visão Geral

Este documento apresenta um estudo completo de como implementar a integração de webhooks no projeto zpwoot, **seguindo rigorosamente a Clean Architecture + Hexagonal Architecture** já estabelecida no projeto.

## 🎯 Objetivos

1. Permitir que usuários configurem webhooks para receber eventos do WhatsApp em tempo real
2. Suportar múltiplos tipos de eventos (mensagens, conexão, grupos, etc.)
3. Garantir entrega confiável com retry e logging
4. Implementar segurança com assinatura de webhooks (HMAC)
5. Fornecer APIs REST para gerenciamento de webhooks

## 🏗️ Arquitetura Proposta (Clean Architecture)

### Estrutura de Pastas Correta
```
zpwoot/
├── internal/
│   ├── core/                                    # CORE LAYER
│   │   ├── domain/                              # Business Logic (Pure)
│   │   │   ├── webhook/                         # ✅ NOVO: Domínio de Webhook
│   │   │   │   ├── entity.go                    # Entidade Webhook
│   │   │   │   ├── repository.go                # Interface do repositório
│   │   │   │   └── service.go                   # Serviço de domínio
│   │   │   ├── session/                         # Existente
│   │   │   └── shared/errors.go                 # Erros compartilhados
│   │   │
│   │   ├── application/                         # Use Cases
│   │   │   ├── dto/
│   │   │   │   └── webhook.go                   # ✅ NOVO: DTOs de Webhook
│   │   │   └── usecase/
│   │   │       └── webhook/                     # ✅ NOVO: Use Cases de Webhook
│   │   │           ├── create_or_update.go      # Upsert webhook
│   │   │           ├── get.go                   # Buscar webhook
│   │   │           ├── list_events.go           # Listar eventos disponíveis
│   │   │           └── delete.go                # Deletar webhook
│   │   │
│   │   └── ports/                               # Interfaces
│   │       ├── input/
│   │       │   └── webhook.go                   # ✅ NOVO: Interface de use cases
│   │       └── output/
│   │           └── webhook_sender.go            # ✅ NOVO: Interface para envio
│   │
│   └── adapters/                                # INFRASTRUCTURE
│       ├── database/
│       │   └── repository/
│       │       └── webhook.go                   # ✅ NOVO: Implementação do repositório
│       │
│       ├── http/
│       │   └── handlers/
│       │       └── webhook.go                   # ✅ NOVO: HTTP handlers
│       │
│       └── webhook/                             # ✅ NOVO: Adapter de webhook
│           ├── sender.go                        # Implementação do WebhookSender
│           ├── signature.go                     # Geração de assinaturas HMAC
│           └── validator.go                     # Validação de URLs
```

## 📊 Modelo de Dados

### Tabela zpWebhooks (já existe no schema)
```sql
CREATE TABLE IF NOT EXISTS "zpWebhooks" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "sessionId" UUID REFERENCES "zpSessions"("id") ON DELETE CASCADE,
    "url" VARCHAR(2048) NOT NULL,
    "secret" VARCHAR(255),
    "events" JSONB NOT NULL DEFAULT '[]',
    "enabled" BOOLEAN NOT NULL DEFAULT true,
    "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

### Entidade Webhook (Go)
```go
type Webhook struct {
    ID        string    `json:"id" db:"id"`
    SessionID string    `json:"sessionId" db:"sessionId"`
    URL       string    `json:"url" db:"url"`
    Secret    *string   `json:"secret,omitempty" db:"secret"`
    Events    []string  `json:"events" db:"events"`
    Enabled   bool      `json:"enabled" db:"enabled"`
    CreatedAt time.Time `json:"createdAt" db:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}
```

## 🎭 Tipos de Eventos Suportados

### Mensagens e Comunicação
- `Message` - Nova mensagem recebida
- `UndecryptableMessage` - Mensagem que não pôde ser descriptografada
- `Receipt` - Confirmação de entrega
- `MediaRetry` - Retry de download de mídia
- `ReadReceipt` - Confirmação de leitura

### Grupos e Contatos
- `GroupInfo` - Informações de grupo atualizadas
- `JoinedGroup` - Entrou em um grupo
- `Picture` - Foto de perfil alterada
- `BlocklistChange` - Mudança na lista de bloqueio
- `Blocklist` - Lista de bloqueio completa

### Conexão e Sessão
- `Connected` - Conectado ao WhatsApp
- `Disconnected` - Desconectado
- `ConnectFailure` - Falha na conexão
- `KeepAliveRestored` - Keep-alive restaurado
- `KeepAliveTimeout` - Timeout do keep-alive
- `LoggedOut` - Sessão deslogada
- `ClientOutdated` - Cliente desatualizado
- `TemporaryBan` - Ban temporário
- `StreamError` - Erro no stream
- `StreamReplaced` - Stream substituído
- `PairSuccess` - Pareamento bem-sucedido
- `PairError` - Erro no pareamento
- `QR` - QR Code gerado
- `QRScannedWithoutMultidevice` - QR escaneado sem multidevice

### Privacidade e Configurações
- `PrivacySettings` - Configurações de privacidade
- `PushNameSetting` - Nome de exibição alterado
- `UserAbout` - Status/sobre alterado

### Sincronização e Estado
- `AppState` - Estado do app
- `AppStateSyncComplete` - Sincronização completa
- `HistorySync` - Sincronização de histórico
- `OfflineSyncCompleted` - Sincronização offline completa
- `OfflineSyncPreview` - Preview de sincronização offline

### Chamadas
- `CallOffer` - Oferta de chamada
- `CallAccept` - Chamada aceita
- `CallTerminate` - Chamada terminada
- `CallOfferNotice` - Notificação de oferta
- `CallRelayLatency` - Latência do relay

### Presença e Atividade
- `Presence` - Presença do usuário (online/offline)
- `ChatPresence` - Presença em chat (digitando, gravando)

### Identidade
- `IdentityChange` - Mudança de identidade

### Erros
- `CATRefreshError` - Erro ao atualizar CAT

### Newsletter (Canais do WhatsApp)
- `NewsletterJoin` - Entrou em canal
- `NewsletterLeave` - Saiu de canal
- `NewsletterMuteChange` - Mudança de mute em canal
- `NewsletterLiveUpdate` - Atualização ao vivo de canal

### Facebook/Meta Bridge
- `FBMessage` - Mensagem do Facebook

### Especial
- `All` - Recebe todos os eventos

## 🔌 APIs REST Propostas

### 1. Upsert Webhook Configuration
```
POST /sessions/{sessionId}/webhook/set
```

**Request Body:**
```json
{
  "url": "https://example.com/webhook",
  "secret": "my-secret-key",
  "events": ["Message", "Receipt", "Connected"],
  "enabled": true
}
```

**Response:**
```json
{
  "id": "uuid",
  "sessionId": "session-uuid",
  "url": "https://example.com/webhook",
  "events": ["Message", "Receipt", "Connected"],
  "enabled": true,
  "createdAt": "2025-10-06T10:00:00Z",
  "updatedAt": "2025-10-06T10:00:00Z"
}
```

### 2. Get Webhook Configuration
```
GET /sessions/{sessionId}/webhook/find
```

**Response:**
```json
{
  "id": "uuid",
  "sessionId": "session-uuid",
  "url": "https://example.com/webhook",
  "events": ["Message", "Receipt", "Connected"],
  "enabled": true,
  "createdAt": "2025-10-06T10:00:00Z",
  "updatedAt": "2025-10-06T10:00:00Z"
}
```

### 3. List Available Events
```
GET /sessions/{sessionId}/webhook/events
```

**Response:**
```json
{
  "events": [
    {
      "type": "Message",
      "category": "Messages and Communication",
      "description": "Nova mensagem recebida"
    },
    {
      "type": "Connected",
      "category": "Connection and Session",
      "description": "Conectado ao WhatsApp"
    }
  ]
}
```

## 🔐 Segurança - Assinatura de Webhooks

### Geração de Assinatura HMAC-SHA256
```go
func GenerateSignature(payload []byte, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write(payload)
    return hex.EncodeToString(h.Sum(nil))
}
```

### Headers Enviados
```
X-Webhook-Signature: <hmac-sha256-signature>
X-Webhook-Timestamp: <unix-timestamp>
X-Webhook-Event: <event-type>
X-Session-ID: <session-id>
Content-Type: application/json
```

### Payload do Webhook
```json
{
  "id": "event-uuid",
  "type": "Message",
  "sessionId": "session-uuid",
  "timestamp": "2025-10-06T10:00:00Z",
  "data": {
    "messageInfo": {
      "id": "msg-id",
      "chat": "5511999999999@s.whatsapp.net",
      "sender": "5511888888888@s.whatsapp.net",
      "pushName": "John Doe",
      "timestamp": "2025-10-06T10:00:00Z",
      "fromMe": false,
      "type": "text",
      "isGroup": false
    },
    "message": {
      "conversation": "Hello World"
    }
  }
}
```

## 🔄 Fluxo de Funcionamento (Clean Architecture)

### 1. Configuração do Webhook
```
HTTP Request
    ↓
HTTP Handler (adapters/http/handlers/webhook.go)
    ↓
WebhookUseCases Interface (ports/input/webhook.go)
    ↓
CreateOrUpdateUseCase (application/usecase/webhook/)
    ↓
Webhook Domain Service (domain/webhook/service.go)
    ↓
Webhook Repository Interface (domain/webhook/repository.go)
    ↓
PostgreSQL Repository Adapter (adapters/database/repository/webhook.go)
    ↓
Database
```

### 2. Recebimento de Evento WhatsApp
```
WhatsApp Event
    ↓
EventHandler (adapters/waclient/events.go)
    ↓
WebhookSender Interface (ports/output/webhook_sender.go)
    ↓
WebhookSender Adapter (adapters/webhook/sender.go)
    ↓
HTTP POST com HMAC signature
    ↓
Cliente Externo
```

### 3. Retry Strategy
- Tentativa 1: Imediato
- Tentativa 2: Após 5 segundos
- Tentativa 3: Após 15 segundos
- Timeout por requisição: 30 segundos

## 📐 Regras de Dependência (Clean Architecture)

### ✅ Permitido:
```
adapters/http/handlers/webhook.go     → ports/input/webhook.go
adapters/database/repository/webhook.go → domain/webhook/repository.go
adapters/webhook/sender.go            → ports/output/webhook_sender.go
application/usecase/webhook/          → domain/webhook/, ports/output/
ports/input/webhook.go                → application/dto/webhook.go
ports/output/webhook_sender.go        → application/dto/webhook.go
domain/webhook/                       → NADA (apenas stdlib)
```

### ❌ Proibido:
```
domain/webhook/                       → adapters/ (NUNCA!)
domain/webhook/                       → application/ (NUNCA!)
application/usecase/webhook/          → adapters/ (usar ports!)
ports/                                → adapters/ (apenas interfaces!)
```

## 📦 Dependências Necessárias

### Bibliotecas Go
- `crypto/hmac` - Geração de assinaturas
- `crypto/sha256` - Hash SHA256
- `encoding/hex` - Encoding hexadecimal
- `net/http` - Cliente HTTP
- `time` - Timeouts e timestamps
- `context` - Controle de contexto
- `encoding/json` - Serialização JSON

### Bibliotecas Externas (já no projeto)
- `go.mau.fi/whatsmeow` - Cliente WhatsApp
- `go.mau.fi/whatsmeow/types/events` - Tipos de eventos
- `github.com/jmoiron/sqlx` - Database
- `github.com/google/uuid` - Geração de UUIDs

## 🧪 Casos de Teste

### Testes Unitários
1. Validação de URL de webhook
2. Geração de assinatura HMAC
3. Serialização de eventos
4. Filtro de eventos
5. Validação de configuração

### Testes de Integração
1. Criar webhook via API
2. Atualizar webhook via API
3. Buscar webhook via API
4. Listar eventos disponíveis
5. Enviar evento para webhook
6. Retry em caso de falha
7. Desabilitar webhook

### Testes End-to-End
1. Configurar webhook → Receber mensagem → Verificar entrega
2. Configurar múltiplos eventos → Verificar filtro
3. Webhook com secret → Verificar assinatura
4. Webhook inválido → Verificar erro

## 🚀 Plano de Implementação (Ordem Correta - Clean Architecture)

### Fase 1: Domain Layer (internal/core/domain/webhook/)
**Regra: Zero dependências externas, apenas stdlib**
- [ ] `entity.go` - Entidade Webhook com regras de negócio
- [ ] `repository.go` - Interface do repositório (contrato)
- [ ] `service.go` - Serviço de domínio (validações de negócio)
- [ ] `events.go` - Tipos de eventos e categorias

### Fase 2: Application DTOs (internal/core/application/dto/)
**Regra: Pode depender apenas de domain**
- [ ] `webhook.go` - DTOs para requests/responses
  - CreateWebhookRequest
  - WebhookResponse
  - EventCategoryResponse
  - ListEventsResponse

### Fase 3: Ports (internal/core/ports/)
**Regra: Apenas interfaces, sem implementações**
- [ ] `input/webhook.go` - Interface WebhookUseCases
  - CreateOrUpdate(ctx, sessionID, request) (*WebhookResponse, error)
  - Get(ctx, sessionID) (*WebhookResponse, error)
  - Delete(ctx, sessionID) error
  - ListEvents(ctx) (*ListEventsResponse, error)

- [ ] `output/webhook_sender.go` - Interface WebhookSender
  - SendWebhook(ctx, event *WebhookEvent) error

### Fase 4: Use Cases (internal/core/application/usecase/webhook/)
**Regra: Depende de domain + ports (não de adapters!)**
- [ ] `create_or_update.go` - Use case de criação/atualização
- [ ] `get.go` - Use case de busca
- [ ] `delete.go` - Use case de deleção
- [ ] `list_events.go` - Use case de listagem de eventos

### Fase 5: Database Adapter (internal/adapters/database/repository/)
**Regra: Implementa interface do domain**
- [ ] `webhook.go` - PostgresWebhookRepository
  - Implementa domain/webhook/repository.go
  - CRUD operations
  - Queries otimizadas

### Fase 6: Webhook Sender Adapter (internal/adapters/webhook/)
**Regra: Implementa interface do port output**
- [ ] `sender.go` - HTTPWebhookSender
  - Implementa ports/output/webhook_sender.go
  - Cliente HTTP com retry
  - Timeout handling

- [ ] `signature.go` - Geração de assinaturas HMAC
- [ ] `validator.go` - Validação de URLs

### Fase 7: HTTP Handlers (internal/adapters/http/handlers/)
**Regra: Usa ports/input (não use cases diretamente!)**
- [ ] `webhook.go` - WebhookHandler
  - SetWebhook(w, r) - POST /sessions/{id}/webhook/set
  - GetWebhook(w, r) - GET /sessions/{id}/webhook/find
  - ListEvents(w, r) - GET /sessions/{id}/webhook/events
  - DeleteWebhook(w, r) - DELETE /sessions/{id}/webhook

### Fase 8: Router Integration (internal/adapters/http/router/)
- [ ] Adicionar rotas de webhook ao router
- [ ] Aplicar middlewares (auth, logging)

### Fase 9: Container Integration (internal/container/)
**Regra: Apenas DI, sem lógica de negócio**
- [ ] Instanciar WebhookRepository
- [ ] Instanciar WebhookService
- [ ] Instanciar WebhookSender
- [ ] Instanciar WebhookUseCases
- [ ] Injetar no EventHandler existente

### Fase 10: Event Handler Integration (internal/adapters/waclient/)
- [ ] Atualizar events.go para usar WebhookSender
- [ ] Carregar configuração de webhook por sessão
- [ ] Filtrar eventos conforme configuração

### Fase 11: Testing
- [ ] Testes unitários de domain
- [ ] Testes unitários de use cases
- [ ] Testes de integração de repository
- [ ] Testes de integração de sender
- [ ] Testes E2E de handlers

### Fase 12: Documentation
- [ ] Swagger/OpenAPI specs
- [ ] README com exemplos
- [ ] Diagramas de arquitetura

## 📚 Referências

- [whatsmeow Documentation](https://pkg.go.dev/go.mau.fi/whatsmeow)
- [whatsmeow Types](https://pkg.go.dev/go.mau.fi/whatsmeow/types)
- [whatsmeow Events](https://pkg.go.dev/go.mau.fi/whatsmeow/types/events)
- [Webhook Best Practices](https://webhooks.fyi/)
- [HMAC Authentication](https://www.okta.com/identity-101/hmac/)

