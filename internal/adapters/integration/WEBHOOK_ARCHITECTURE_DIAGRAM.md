# Diagrama de Arquitetura - Webhook Integration

## Visão Geral da Arquitetura Clean + Hexagonal

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          EXTERNAL WORLD                                 │
│  ┌──────────────────┐                        ┌──────────────────┐      │
│  │  HTTP Clients    │                        │  External APIs   │      │
│  │  (Postman, etc)  │                        │  (Webhook URLs)  │      │
│  └────────┬─────────┘                        └────────▲─────────┘      │
└───────────┼──────────────────────────────────────────┼─────────────────┘
            │                                           │
            │ HTTP Request                              │ HTTP POST
            ▼                                           │
┌─────────────────────────────────────────────────────────────────────────┐
│                    ADAPTERS LAYER (Infrastructure)                      │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  HTTP Handlers (adapters/http/handlers/webhook.go)              │  │
│  │  - SetWebhook(w, r)                                              │  │
│  │  - GetWebhook(w, r)                                              │  │
│  │  - ListEvents(w, r)                                              │  │
│  └────────┬─────────────────────────────────────────────────────────┘  │
│           │ uses                                                        │
│           ▼                                                             │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  WebhookSender Adapter (adapters/webhook/sender.go)             │  │
│  │  - SendWebhook(ctx, url, secret, event)                         │  │
│  │  - Retry logic                                                   │  │
│  │  - HMAC signature                                                │  │
│  └────────────────────────────────────────────────┬─────────────────┘  │
│                                                    │                    │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  Database Repository (adapters/database/repository/webhook.go)  │  │
│  │  - Create(ctx, webhook)                                          │  │
│  │  - GetBySessionID(ctx, sessionID)                                │  │
│  │  - Update(ctx, webhook)                                          │  │
│  └────────┬─────────────────────────────────────────────────────────┘  │
└───────────┼──────────────────────────────────────────────────────────┬─┘
            │ implements                                               │
            ▼                                                          │
┌─────────────────────────────────────────────────────────────────────┼──┐
│                    PORTS LAYER (Interfaces)                         │  │
│  ┌──────────────────────────────────────────────────────────────┐  │  │
│  │  INPUT PORTS (ports/input/webhook.go)                        │  │  │
│  │  interface WebhookUseCases {                                 │  │  │
│  │    CreateOrUpdate(ctx, sessionID, req) (*Response, error)    │  │  │
│  │    Get(ctx, sessionID) (*Response, error)                    │  │  │
│  │    Delete(ctx, sessionID) error                              │  │  │
│  │    ListEvents(ctx) (*ListEventsResponse, error)              │  │  │
│  │  }                                                            │  │  │
│  └────────▲─────────────────────────────────────────────────────┘  │  │
│           │ implements                                              │  │
│           │                                                          │  │
│  ┌──────────────────────────────────────────────────────────────┐  │  │
│  │  OUTPUT PORTS (ports/output/webhook_sender.go)               │  │  │
│  │  interface WebhookSender {                                   │  │  │
│  │    SendWebhook(ctx, url, secret, event) error                │  │  │
│  │  }                                                            │  │  │
│  └────────▲─────────────────────────────────────────────────────┘  │  │
└───────────┼──────────────────────────────────────────────────────┼──┘
            │                                                       │
            │ used by                                               │ implements
            ▼                                                       │
┌─────────────────────────────────────────────────────────────────┼──┐
│                    APPLICATION LAYER (Use Cases)                │  │
│  ┌──────────────────────────────────────────────────────────────┼┐ │
│  │  Use Cases (application/usecase/webhook/)                    ││ │
│  │  ┌────────────────────────────────────────────────────────┐  ││ │
│  │  │  CreateOrUpdateUseCase                                  │  ││ │
│  │  │  - Execute(ctx, sessionID, req) (*Response, error)      │  ││ │
│  │  │  - Validates URL                                        │  ││ │
│  │  │  - Validates events                                     │  ││ │
│  │  │  - Calls domain service                                 │  ││ │
│  │  │  - Calls repository                                     │  ││ │
│  │  └────────────────────────────────────────────────────────┘  ││ │
│  │  ┌────────────────────────────────────────────────────────┐  ││ │
│  │  │  GetUseCase                                             │  ││ │
│  │  │  - Execute(ctx, sessionID) (*Response, error)           │  ││ │
│  │  └────────────────────────────────────────────────────────┘  ││ │
│  │  ┌────────────────────────────────────────────────────────┐  ││ │
│  │  │  ListEventsUseCase                                      │  ││ │
│  │  │  - Execute(ctx) (*ListEventsResponse, error)            │  ││ │
│  │  └────────────────────────────────────────────────────────┘  ││ │
│  └──────────────────────────────────────────────────────────────┘│ │
│           │ uses                                                  │ │
│           ▼                                                       │ │
│  ┌──────────────────────────────────────────────────────────────┐│ │
│  │  DTOs (application/dto/webhook.go)                           ││ │
│  │  - CreateWebhookRequest                                      ││ │
│  │  - WebhookResponse                                           ││ │
│  │  - EventCategoryResponse                                     ││ │
│  └──────────────────────────────────────────────────────────────┘│ │
└───────────┼────────────────────────────────────────────────────────┘
            │ uses
            ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    DOMAIN LAYER (Business Logic)                        │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  Webhook Entity (domain/webhook/entity.go)                       │  │
│  │  ┌────────────────────────────────────────────────────────────┐  │  │
│  │  │  type Webhook struct {                                      │  │  │
│  │  │    ID, SessionID, URL, Secret, Events, Enabled, ...        │  │  │
│  │  │  }                                                          │  │  │
│  │  │                                                             │  │  │
│  │  │  Business Rules:                                            │  │  │
│  │  │  - NewWebhook(sessionID, url, events) *Webhook             │  │  │
│  │  │  - SetSecret(secret)                                        │  │  │
│  │  │  - Enable() / Disable()                                     │  │  │
│  │  │  - HasEvent(eventType) bool                                 │  │  │
│  │  │  - IsActive() bool                                          │  │  │
│  │  └────────────────────────────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  Repository Interface (domain/webhook/repository.go)             │  │
│  │  interface Repository {                                          │  │
│  │    Create(ctx, webhook) error                                    │  │
│  │    GetBySessionID(ctx, sessionID) (*Webhook, error)              │  │
│  │    Update(ctx, webhook) error                                    │  │
│  │    Delete(ctx, id) error                                         │  │
│  │  }                                                                │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  Domain Service (domain/webhook/service.go)                      │  │
│  │  - ValidateURL(url) error                                        │  │
│  │  - ValidateEvents(events) error                                  │  │
│  │  - Business validation logic                                     │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  Event Types (domain/webhook/events.go)                          │  │
│  │  - EventMessage, EventConnected, EventDisconnected, ...          │  │
│  │  - GetAllEventTypes() []EventType                                │  │
│  │  - GetAllEventCategories() []EventCategory                       │  │
│  └──────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

## Fluxo de Dados - Criar/Atualizar Webhook

```
1. HTTP Request
   POST /sessions/{sessionId}/webhook/set
   Body: { "url": "https://...", "events": ["Message", "Connected"] }
   │
   ▼
2. WebhookHandler.SetWebhook(w, r)
   - Extrai sessionId da URL
   - Decodifica JSON do body
   │
   ▼
3. webhookUseCases.CreateOrUpdate(ctx, sessionID, req)
   [Interface - ports/input/webhook.go]
   │
   ▼
4. CreateOrUpdateUseCase.Execute(ctx, sessionID, req)
   [Use Case - application/usecase/webhook/]
   │
   ├─▶ webhookService.ValidateURL(req.URL)
   │   [Domain Service - domain/webhook/service.go]
   │
   ├─▶ webhookService.ValidateEvents(req.Events)
   │   [Domain Service - domain/webhook/service.go]
   │
   ├─▶ webhookRepo.GetBySessionID(ctx, sessionID)
   │   [Repository Interface - domain/webhook/repository.go]
   │   │
   │   ▼
   │   PostgresWebhookRepository.GetBySessionID(ctx, sessionID)
   │   [Adapter - adapters/database/repository/webhook.go]
   │   │
   │   ▼
   │   SELECT FROM "zpWebhooks" WHERE "sessionId" = $1
   │
   ├─▶ Se existe: webhook.Update(...)
   │   Se não: webhook.NewWebhook(sessionID, url, events)
   │   [Domain Entity - domain/webhook/entity.go]
   │
   ├─▶ webhook.SetSecret(secret)
   │   [Domain Entity - domain/webhook/entity.go]
   │
   └─▶ webhookRepo.Create(ctx, webhook) ou Update(ctx, webhook)
       [Repository Interface - domain/webhook/repository.go]
       │
       ▼
       PostgresWebhookRepository.Create/Update(ctx, webhook)
       [Adapter - adapters/database/repository/webhook.go]
       │
       ▼
       INSERT/UPDATE "zpWebhooks" ...
       │
       ▼
5. Return WebhookResponse DTO
   │
   ▼
6. HTTP Response
   Status: 200 OK
   Body: { "id": "...", "sessionId": "...", "url": "...", ... }
```

## Fluxo de Dados - Enviar Evento via Webhook

```
1. WhatsApp Event Received
   [whatsmeow library]
   │
   ▼
2. EventHandler.HandleEvent(client, event)
   [adapters/waclient/events.go]
   │
   ├─▶ Identifica tipo de evento (Message, Connected, etc)
   │
   ├─▶ Verifica se webhook está configurado
   │   client.WebhookURL != ""
   │
   ├─▶ Verifica se evento está na lista
   │   client.Events contém o tipo de evento
   │
   └─▶ webhookSender.SendWebhook(ctx, url, secret, event)
       [Interface - ports/output/webhook_sender.go]
       │
       ▼
3. HTTPWebhookSender.SendWebhook(ctx, url, secret, event)
   [Adapter - adapters/webhook/sender.go]
   │
   ├─▶ Serializa evento para JSON
   │
   ├─▶ Gera assinatura HMAC-SHA256
   │   signature.Generate(payload, secret)
   │   [adapters/webhook/signature.go]
   │
   ├─▶ Prepara headers
   │   X-Webhook-Signature: <hmac>
   │   X-Webhook-Timestamp: <timestamp>
   │   X-Webhook-Event: <event-type>
   │   X-Session-ID: <session-id>
   │
   ├─▶ HTTP POST para webhook URL
   │   Timeout: 30 segundos
   │
   ├─▶ Se falhar: Retry
   │   - Tentativa 1: Imediato
   │   - Tentativa 2: Após 5s
   │   - Tentativa 3: Após 15s
   │
   └─▶ Log resultado (sucesso/falha)
       │
       ▼
4. External Webhook Endpoint
   Recebe evento e processa
```

## Regras de Dependência (Setas indicam "depende de")

```
HTTP Handlers ──────────────────────────────────┐
                                                 │
WebhookSender Adapter ──────────────────────┐   │
                                             │   │
Database Repository Adapter ────────────┐   │   │
                                         │   │   │
                                         ▼   ▼   ▼
                                    ┌────────────────┐
                                    │  PORTS         │
                                    │  (Interfaces)  │
                                    └────────────────┘
                                         │   │   │
                                         ▼   ▼   ▼
                                    ┌────────────────┐
                                    │  USE CASES     │
                                    └────────────────┘
                                             │
                                             ▼
                                    ┌────────────────┐
                                    │  DOMAIN        │
                                    │  (Pure Logic)  │
                                    └────────────────┘
```

## Resumo das Camadas

| Camada | Localização | Dependências | Responsabilidade |
|--------|-------------|--------------|------------------|
| **Domain** | `internal/core/domain/webhook/` | Nenhuma (stdlib) | Regras de negócio puras |
| **Application** | `internal/core/application/` | Domain + Ports | Orquestração de use cases |
| **Ports** | `internal/core/ports/` | Domain + DTOs | Interfaces (contratos) |
| **Adapters** | `internal/adapters/` | Tudo | Implementações de infraestrutura |

## Benefícios desta Arquitetura

1. **Testabilidade**: Domain e Application podem ser testados sem infraestrutura
2. **Manutenibilidade**: Mudanças em adapters não afetam o core
3. **Flexibilidade**: Fácil trocar implementações (ex: PostgreSQL → MongoDB)
4. **Clareza**: Separação clara de responsabilidades
5. **Independência**: Core não depende de frameworks externos

