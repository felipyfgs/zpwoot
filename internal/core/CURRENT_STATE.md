# Estado Atual do Projeto - Análise de Camadas

Este documento mapeia **o que cada camada implementa atualmente** no projeto zpwoot.

---

## 📊 Mapeamento Completo

### 1. **DOMAIN** (`internal/core/domain/`)

#### ✅ **O que existe:**

```
internal/core/domain/
├── session/
│   ├── entity.go          # Entidade Session com campos e métodos
│   ├── repository.go      # Interface Repository (8 métodos)
│   └── service.go         # Service com lógica de negócio
└── shared/
    └── errors.go          # Erros de domínio (ErrSessionNotFound, etc.)
```

#### **Detalhes:**

**`session/entity.go`**
- ✅ Struct `Session` com 12 campos
- ✅ Type `Status` (StatusDisconnected, StatusConnecting, StatusConnected, etc.)
- ✅ Métodos: `NewSession()`, `SetConnected()`, `SetDisconnected()`, `SetQRCode()`, `GetStatus()`

**`session/repository.go`**
- ✅ Interface `Repository` com 8 métodos:
  - `Create()`
  - `GetByID()`
  - `GetByJID()`
  - `Update()`
  - `Delete()`
  - `List()`
  - `UpdateStatus()`
  - `UpdateQRCode()`

**`session/service.go`**
- ✅ Struct `Service` com campo `repo Repository`
- ✅ Métodos de negócio:
  - `CreateSession()`
  - `GetSession()`
  - `UpdateSessionStatus()`
  - `UpdateQRCode()`
  - `ListSessions()`
  - `DeleteSession()`

**`shared/errors.go`**
- ✅ Erros de domínio:
  - `ErrSessionNotFound`
  - `ErrInvalidStatus`
  - `ErrSessionAlreadyExists`

#### ❌ **O que falta:**

- Agregado `Message` (futuro)
- Agregado `Webhook` (futuro)
- Value Objects mais ricos (ex: `PhoneNumber`, `JID`)
- Domain Events (ex: `SessionCreated`, `MessageSent`)

---

### 2. **APPLICATION** (`internal/core/application/`)

#### ✅ **O que existe:**

```
internal/core/application/
├── dto/
│   ├── common.go          # Response, ErrorResponse, PaginationRequest
│   ├── session.go         # 8 DTOs de sessão
│   └── message.go         # 3 DTOs de mensagem
├── interfaces/            # ⚠️ DEVERIA ESTAR EM ports/output/
│   ├── whatsapp.go       # Interface WhatsAppClient (10 métodos)
│   └── notification.go   # Interface NotificationService (5 métodos)
├── usecase/
│   ├── session/
│   │   ├── create.go     # CreateUseCase
│   │   ├── connect.go    # ConnectUseCase
│   │   ├── disconnect.go # DisconnectUseCase
│   │   ├── logout.go     # LogoutUseCase
│   │   ├── get.go        # GetUseCase
│   │   ├── list.go       # ListUseCase
│   │   ├── delete.go     # DeleteUseCase
│   │   ├── qr.go         # QRUseCase
│   │   └── usecases.go   # Container de use cases
│   └── message/
│       ├── send.go       # SendUseCase
│       └── receive.go    # ReceiveUseCase
└── validators/
    ├── session.go        # ValidateCreateSession, ValidateSendMessage
    └── message.go        # ValidateMessageRequest
```

#### **Detalhes:**

**DTOs (`dto/`)**

`common.go`:
- ✅ `Response[T]` - Wrapper genérico de resposta
- ✅ `ErrorResponse` - Estrutura de erro
- ✅ `PaginationRequest` - Paginação

`session.go`:
- ✅ `CreateSessionRequest`
- ✅ `SessionResponse`
- ✅ `SessionListResponse`
- ✅ `QRCodeResponse`
- ✅ `ConnectSessionRequest`
- ✅ `DisconnectSessionRequest`
- ✅ `DeleteSessionRequest`
- ✅ Funções de conversão: `NewSessionResponse()`, `NewQRCodeResponse()`

`message.go`:
- ✅ `SendMessageRequest`
- ✅ `MessageResponse`
- ✅ `ReceiveMessageRequest`

**Interfaces (`interfaces/`)** - ⚠️ MOVER PARA `ports/output/`

`whatsapp.go`:
- ✅ Interface `WhatsAppClient` com 10 métodos
- ✅ Tipos relacionados: `SessionStatus`, `QRCodeInfo`, `MessageResult`, `MediaData`, `Location`, `ContactInfo`
- ✅ Erros: `WhatsAppError`, `ErrSessionNotFound`, `ErrSessionNotConnected`, etc.

`notification.go`:
- ✅ Interface `NotificationService` com 5 métodos
- ✅ Tipos de eventos: `WebhookEvent`, `MessageEvent`, `SessionEvent`, `QRCodeEvent`
- ✅ Constantes: `EventTypeSessionConnected`, `EventTypeMessageReceived`, etc.

**Use Cases (`usecase/`)**

Todos os use cases seguem o padrão:
```go
type XxxUseCase struct {
    sessionService  *session.Service
    whatsappClient  interfaces.WhatsAppClient
    notificationSvc interfaces.NotificationService
}

func (uc *XxxUseCase) Execute(ctx context.Context, ...) (..., error) {
    // 1. Validação
    // 2. Chamada ao domain service
    // 3. Chamada aos adapters (WhatsApp, Notification)
    // 4. Retorno de DTO
}
```

**Session Use Cases:**
- ✅ `CreateUseCase` - Cria sessão + inicializa WhatsApp
- ✅ `ConnectUseCase` - Conecta ao WhatsApp + gera QR
- ✅ `DisconnectUseCase` - Desconecta graciosamente
- ✅ `LogoutUseCase` - Faz logout e limpa dados
- ✅ `GetUseCase` - Obtém detalhes + sincroniza status
- ✅ `ListUseCase` - Lista sessões com paginação
- ✅ `DeleteUseCase` - Deleta sessão + cleanup
- ✅ `QRUseCase` - Gerencia QR Code (get + refresh)

**Message Use Cases:**
- ✅ `SendUseCase` - Envia mensagens (text, media, location, contact)
- ✅ `ReceiveUseCase` - Processa mensagens recebidas + notifica

**Validators (`validators/`)**
- ✅ `ValidateCreateSession()` - Valida criação de sessão
- ✅ `ValidateSendMessage()` - Valida envio de mensagem

#### ❌ **O que falta:**

- Mover `interfaces/` para `ports/output/`
- Criar mais validadores (connect, disconnect, etc.)
- Adicionar testes unitários dos use cases

---

### 3. **PORTS** (`internal/core/ports/`)

#### ✅ **O que existe:**

```
internal/core/ports/
└── (vazio - diretório criado mas sem conteúdo)
```

#### ❌ **O que deveria existir:**

```
internal/core/ports/
├── output/                    # Portas de saída
│   ├── whatsapp.go           # ← MOVER de application/interfaces/
│   ├── notification.go       # ← MOVER de application/interfaces/
│   └── logger.go             # ← CRIAR (não existe!)
└── input/                     # Portas de entrada (OPCIONAL)
    ├── session.go            # Interfaces de use cases de sessão
    └── message.go            # Interfaces de use cases de mensagem
```

---

## 🔍 Análise de Dependências Externas

### **Interfaces que o Core precisa (Output Ports):**

#### 1. **WhatsAppClient** ✅ (existe em `application/interfaces/`)

**Localização atual:** `internal/core/application/interfaces/whatsapp.go`  
**Deveria estar em:** `internal/core/ports/output/whatsapp.go`

**Métodos:**
```go
type WhatsAppClient interface {
    // Session Management
    CreateSession(ctx context.Context, sessionID string) error
    GetSessionStatus(ctx context.Context, sessionID string) (*SessionStatus, error)
    DeleteSession(ctx context.Context, sessionID string) error
    
    // Connection
    ConnectSession(ctx context.Context, sessionID string) error
    DisconnectSession(ctx context.Context, sessionID string) error
    LogoutSession(ctx context.Context, sessionID string) error
    IsConnected(ctx context.Context, sessionID string) bool
    IsLoggedIn(ctx context.Context, sessionID string) bool
    
    // QR Code
    GetQRCode(ctx context.Context, sessionID string) (*QRCodeInfo, error)
    
    // Messaging
    SendTextMessage(ctx context.Context, sessionID, to, text string) (*MessageResult, error)
    SendMediaMessage(ctx context.Context, sessionID, to string, media *MediaData) (*MessageResult, error)
    SendLocationMessage(ctx context.Context, sessionID, to string, location *Location) (*MessageResult, error)
    SendContactMessage(ctx context.Context, sessionID, to string, contact *ContactInfo) (*MessageResult, error)
}
```

**Implementação:** `internal/adapters/waclient/whatsapp_adapter.go`

---

#### 2. **NotificationService** ✅ (existe em `application/interfaces/`)

**Localização atual:** `internal/core/application/interfaces/notification.go`  
**Deveria estar em:** `internal/core/ports/output/notification.go`

**Métodos:**
```go
type NotificationService interface {
    SendWebhook(ctx context.Context, event *WebhookEvent) error
    NotifySessionConnected(ctx context.Context, sessionID string, deviceJID string) error
    NotifySessionDisconnected(ctx context.Context, sessionID string) error
    NotifyQRCodeGenerated(ctx context.Context, sessionID string, qrCode string, expiresAt time.Time) error
    NotifyMessageReceived(ctx context.Context, sessionID string, message *MessageEvent) error
    NotifyMessageSent(ctx context.Context, sessionID string, messageID string) error
}
```

**Implementação:** ❌ Não existe! (atualmente é `nil` no container)

---

#### 3. **Logger** ❌ (NÃO EXISTE INTERFACE!)

**Problema:** O logger é usado como implementação concreta em todo o código.

**Uso atual:**
```go
// ❌ Acoplamento direto à implementação
import "zpwoot/internal/adapters/logger"

type Container struct {
    logger *logger.Logger  // Implementação concreta!
}
```

**Deveria ser:**
```go
// ✅ Uso de interface (port)
import "zpwoot/internal/core/ports/output"

type Container struct {
    logger output.Logger  // Interface!
}
```

**Interface proposta:**
```go
// internal/core/ports/output/logger.go
package output

import "context"

type Logger interface {
    // Níveis de log
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
    
    // Context-aware
    WithContext(ctx context.Context) Logger
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    WithError(err error) Logger
    
    // Domain-specific
    WithComponent(component string) Logger
    WithRequestID(requestID string) Logger
    WithSessionID(sessionID string) Logger
}

type Field struct {
    Key   string
    Value interface{}
}
```

**Implementação:** `internal/adapters/logger/logger_adapter.go` (criar)

---

## 📋 Checklist de Migração para Ports

### **Fase 1: Criar estrutura de ports**

- [ ] Criar `internal/core/ports/output/`
- [ ] Criar `internal/core/ports/input/` (opcional)
- [ ] Criar `internal/core/ports/README.md` ✅ (já criado)

### **Fase 2: Mover interfaces existentes**

- [ ] Mover `application/interfaces/whatsapp.go` → `ports/output/whatsapp.go`
- [ ] Mover `application/interfaces/notification.go` → `ports/output/notification.go`
- [ ] Atualizar imports em todos os arquivos

### **Fase 3: Criar Logger Port**

- [ ] Criar `ports/output/logger.go` com interface
- [ ] Criar `adapters/logger/logger_adapter.go` implementando a interface
- [ ] Atualizar `container.go` para usar a interface
- [ ] Atualizar todos os use cases que usam logger (se houver)

### **Fase 4: Criar Input Ports (opcional)**

- [ ] Criar `ports/input/session.go` com interfaces de use cases
- [ ] Criar `ports/input/message.go` com interfaces de use cases
- [ ] Fazer use cases implementarem as interfaces
- [ ] Atualizar handlers HTTP para usar as interfaces

### **Fase 5: Cleanup**

- [ ] Remover `application/interfaces/` (após migração completa)
- [ ] Atualizar `ARCHITECTURE.md`
- [ ] Atualizar `README.md` do projeto
- [ ] Executar `go mod tidy`
- [ ] Executar testes

---

## 🎯 Próximos Passos Recomendados

1. **Criar estrutura de ports:**
   ```bash
   mkdir -p internal/core/ports/output
   mkdir -p internal/core/ports/input  # opcional
   ```

2. **Mover WhatsApp interface:**
   ```bash
   mv internal/core/application/interfaces/whatsapp.go internal/core/ports/output/whatsapp.go
   ```

3. **Mover Notification interface:**
   ```bash
   mv internal/core/application/interfaces/notification.go internal/core/ports/output/notification.go
   ```

4. **Criar Logger port:**
   - Criar `internal/core/ports/output/logger.go`
   - Criar `internal/adapters/logger/logger_adapter.go`

5. **Atualizar imports:**
   ```bash
   find . -name "*.go" -type f | xargs sed -i 's|zpwoot/internal/core/application/interfaces|zpwoot/internal/core/ports/output|g'
   ```

6. **Remover diretório vazio:**
   ```bash
   rmdir internal/core/application/interfaces
   ```

7. **Verificar compilação:**
   ```bash
   go mod tidy
   go build ./...
   ```

---

## ✅ Conclusão

O projeto já está **bem estruturado** seguindo Clean Architecture, mas precisa de alguns ajustes:

1. ✅ **Domain** está correto e independente
2. ✅ **Application** está bem organizado com use cases claros
3. ⚠️ **Ports** precisa ser criado e populado
4. ⚠️ **Interfaces** estão no lugar errado (`application/interfaces/` → `ports/output/`)
5. ❌ **Logger** não tem interface (acoplamento direto)

Após a migração, o projeto estará 100% alinhado com Clean Architecture e Hexagonal Architecture.

