# 🏗️ Arquitetura zpwoot - Clean Architecture Pragmática

## 📋 Visão Geral

O **zpwoot** é uma **API Gateway para WhatsApp Business** que implementa uma **Clean Architecture Pragmática** rigorosamente estruturada. A arquitetura combina os princípios da Clean Architecture com padrões modernos como DDD, CQRS e Event-Driven Architecture, garantindo alta performance, escalabilidade e manutenibilidade.

## 🎯 Princípios Fundamentais

### 1. **Separação Rigorosa de Responsabilidades**
- Cada camada tem uma responsabilidade específica e bem definida
- Dependências fluem sempre para dentro (Dependency Inversion Principle)
- Lógica de negócio completamente isolada de detalhes de implementação
- Zero dependências externas no core domain

### 2. **Domain-Driven Design (DDD)**
- Bounded contexts bem definidos (Session, Messaging, Group, Contact)
- Rich domain models com comportamento encapsulado
- Value objects para conceitos de negócio
- Domain services para lógica que não pertence a entidades

### 3. **Testabilidade e Qualidade**
- 100% das interfaces mockáveis para testes unitários
- Dependency injection em todas as camadas
- Lógica de negócio testável isoladamente
- Cobertura completa: unitários, integração e end-to-end

## 🏛️ Estrutura de Camadas

```
zpwoot/
├── internal/                    # 🏠 Application Core
│   ├── core/                    # 🎯 Domain Layer (Business Logic)
│   │   ├── session/             # Session domain (WhatsApp sessions)
│   │   ├── messaging/           # Messaging domain (Messages & sync)
│   │   ├── group/               # Group domain (WhatsApp groups)
│   │   ├── contact/             # Contact domain (Contact management)
│   │   └── shared/              # Shared domain concepts
│   ├── usecases/                # 🔧 Application Layer (Use Cases)
│   │   ├── session_usecase.go   # Session orchestration
│   │   ├── message_usecase.go   # Message orchestration
│   │   ├── group_usecase.go     # Group orchestration
│   │   └── shared/              # Shared application services
│   └── adapters/                # 🔌 Infrastructure Layer
│       ├── repository/          # Data persistence implementations
│       ├── server/              # HTTP server, routing & handlers
│       └── waclient/            # WhatsApp client integration
├── platform/                   # 🏗️ Platform Layer (Cross-cutting)
│   ├── config/                  # Configuration management
│   ├── database/                # Database connection & migrations
│   ├── logger/                  # Structured logging
│   └── container/               # Dependency injection container
└── cmd/                        # 🚀 Entry Points
    └── server/                  # HTTP server entry point
```

## 📁 Detalhamento das Camadas

### 🎯 **CORE - Domain Layer (Camada de Domínio)**

**Responsabilidade:** Contém a lógica de negócio pura, entidades ricas, value objects e contratos. É o coração do sistema, completamente isolado de detalhes externos.

**Estrutura Detalhada:**
```
internal/core/
├── session/                    # 📱 Session Domain
│   ├── models.go               # Session entity, ProxyConfig value object
│   ├── service.go              # Session business rules & orchestration
│   └── contracts.go            # Repository & WhatsAppGateway interfaces
├── messaging/                  # 💬 Messaging Domain
│   ├── models.go               # Message entity, MessageType enums
│   ├── service.go              # Message business rules & validation
│   └── contracts.go            # Repository & MessageGateway interfaces
├── group/                      # 👥 Group Domain
│   ├── models.go               # Group entity, GroupSettings value object
│   ├── service.go              # Group permissions & business rules
│   └── contracts.go            # Repository & GroupGateway interfaces
├── contact/                    # 📞 Contact Domain
│   ├── models.go               # Contact entity, ContactInfo value object
│   ├── service.go              # Contact validation & business rules
│   └── contracts.go            # Repository interface
└── shared/                     # 🔗 Shared Domain Concepts
    ├── errors/                 # Domain-specific errors
    │   └── errors.go           # ErrSessionNotFound, ErrInvalidMessage, etc.
    ├── events/                 # Domain events
    │   └── events.go           # SessionCreated, MessageSent, etc.
    └── types/                  # Common value objects
        └── types.go            # JID, PhoneNumber, Timestamp types
```

### **Domínios Implementados (Bounded Contexts)**

#### **1. Session Domain** 📱
**Responsabilidade**: Gerenciamento completo do ciclo de vida de sessões WhatsApp
- **Entidades**: Session (com rich behavior)
- **Value Objects**: ProxyConfig, SessionStatus
- **Business Rules**: Validação de conexão, QR code lifecycle, proxy configuration
- **Interfaces**: Repository, WhatsAppGateway, EventHandler

#### **2. Messaging Domain** 💬
**Responsabilidade**: Sistema completo de mensagens com sincronização Chatwoot
- **Entidades**: Message (com sync status)
- **Value Objects**: MessageType, SyncStatus
- **Business Rules**: Validação de conteúdo, sync logic, delivery status
- **Interfaces**: Repository, MessageGateway

#### **3. Group Domain** 👥
**Responsabilidade**: Gerenciamento de grupos WhatsApp e permissões
- **Entidades**: Group, GroupParticipant
- **Value Objects**: GroupSettings, ParticipantRole
- **Business Rules**: Permissões de admin, validação de participantes
- **Interfaces**: Repository, GroupGateway

#### **4. Contact Domain** 📞
**Responsabilidade**: Gerenciamento de contatos e verificação de números
- **Entidades**: Contact
- **Value Objects**: ContactInfo, PhoneNumber
- **Business Rules**: Validação de números, verificação WhatsApp
- **Interfaces**: Repository

#### **5. Shared Domain** 🔗
**Responsabilidade**: Conceitos compartilhados entre domínios
- **Errors**: Domain-specific errors tipados
- **Events**: Domain events para comunicação assíncrona
- **Types**: Value objects comuns (JID, Timestamp, etc.)

**Padrões Implementados:**
- **Rich Domain Models**: Entidades com comportamento encapsulado
- **Value Objects**: ProxyConfig, SessionStatus, MessageType
- **Factory Methods**: NewSession(), NewMessage() para criação consistente
- **Repository Pattern**: Interfaces para persistência abstrata
- **Gateway Pattern**: Interfaces para integrações externas
- **Domain Services**: Lógica que não pertence a uma entidade específica
- **Domain Events**: Eventos de negócio para comunicação assíncrona

**Regras Rigorosas:**
- ✅ **PODE:** Definir entidades, value objects, regras de negócio
- ✅ **PODE:** Definir interfaces para dependências externas (Repository, Gateway)
- ✅ **PODE:** Usar outros módulos do core (session pode usar shared/errors)
- ✅ **PODE:** Implementar validações de domínio
- ❌ **NÃO PODE:** Importar adapters, services ou platform
- ❌ **NÃO PODE:** Conhecer detalhes de HTTP, SQL, APIs externas
- ❌ **NÃO PODE:** Ter dependências de frameworks externos
- ❌ **NÃO PODE:** Conter lógica de infraestrutura

**Exemplo de Implementação:**
```go
// ✅ CORRETO - Rich Domain Model
type Session struct {
    ID              uuid.UUID
    Name            string
    IsConnected     bool
    ConnectionError *string
    ProxyConfig     *ProxyConfig // Value Object
}

// ✅ CORRETO - Business Method
func (s *Session) UpdateConnectionStatus(connected bool, error string) {
    s.IsConnected = connected
    if !connected {
        s.ConnectionError = &error
    } else {
        s.ConnectionError = nil
    }
    s.UpdatedAt = time.Now()
}

// ✅ CORRETO - Domain Service
type Service struct {
    repository Repository        // Interface
    gateway    WhatsAppGateway  // Interface
}

func (s *Service) CreateSession(req *CreateSessionRequest) (*Session, error) {
    // Validações de domínio
    if len(req.Name) == 0 {
        return nil, ErrInvalidSessionName
    }

    // Lógica de negócio
    session := NewSession(req.Name)
    return session, s.repository.Create(ctx, session)
}
```

**Imports Válidos/Inválidos:**
```go
// ✅ PERMITIDO
import "zpwoot/internal/core/session"
import "zpwoot/internal/core/messaging"
import "github.com/google/uuid"           // Bibliotecas básicas OK

// ❌ PROIBIDO
import "zpwoot/internal/adapters/repository"
import "zpwoot/internal/usecases"
import "zpwoot/platform/database"
import "github.com/gin-gonic/gin"         // Frameworks externos
import "github.com/jmoiron/sqlx"          // Bibliotecas de infraestrutura
```

### 🔧 **USECASES - Application Layer (Camada de Aplicação)**

**Responsabilidade:** Orquestra use cases, coordena operações entre múltiplos domínios, gerencia transações e implementa a lógica de aplicação. Atua como uma fachada entre a interface externa e o core domain.

**Estrutura Detalhada:**
```
internal/usecases/
├── session_usecase.go          # 📱 Session use cases orchestration
├── message_usecase.go          # 💬 Message use cases orchestration
├── group_usecase.go            # 👥 Group use cases orchestration
├── chatwoot_usecase.go         # 🔗 Chatwoot integration orchestration
└── shared/                     # 🔗 Shared application services
    ├── validation/             # Input validation logic
    │   └── validator.go        # Struct validation, custom rules
    ├── mapping/                # DTO ↔ Domain mapping
    │   └── mapper.go           # Conversion utilities
    └── contracts/              # Application DTOs & contracts
        ├── session_contracts.go # Session DTOs
        ├── message_contracts.go # Message DTOs
        └── response_contracts.go # Common response formats
```

### **Use Cases Implementados (Application Services)**

#### **1. SessionUseCase** 📱
**Responsabilidade**: Orquestração completa de use cases de sessão
- **Use Cases**: CreateSession, ConnectSession, DisconnectSession, GetSession
- **Coordenação**: Session domain + WhatsApp gateway
- **Validação**: Session name, proxy config, connection parameters
- **Observabilidade**: Logging estruturado de todas as operações

#### **2. MessageUseCase** 💬
**Responsabilidade**: Orquestração de mensagens e sincronização
- **Use Cases**: SendMessage, GetMessages, SyncWithChatwoot
- **Coordenação**: Messaging domain + Session domain + Chatwoot integration
- **Tipos Suportados**: Text, Media, Document, Interactive messages
- **Sync Logic**: Bidirectional sync com Chatwoot

#### **3. GroupUseCase** 👥
**Responsabilidade**: Orquestração de operações de grupo
- **Use Cases**: CreateGroup, ManageParticipants, UpdateSettings
- **Coordenação**: Group domain + Session domain
- **Validação**: Permissions, participant limits, group settings
- **Business Logic**: Admin permissions, participant management

#### **4. ChatwootUseCase** 🔗
**Responsabilidade**: Integração completa com Chatwoot
- **Use Cases**: ConfigureIntegration, SyncMessages, ManageConversations
- **Coordenação**: Messaging domain + external Chatwoot API
- **Features**: Inbox creation, conversation mapping, webhook handling
- **Sync Strategy**: Real-time bidirectional synchronization

**Padrões Implementados:**
- **Application Service Pattern**: Orquestração de use cases complexos
- **DTO Pattern**: Data Transfer Objects para comunicação externa
- **Validation Pattern**: Validação centralizada de entrada
- **Mapping Pattern**: Conversão entre DTOs e Domain Objects
- **Transaction Script**: Coordenação de operações transacionais
- **Facade Pattern**: Interface simplificada para operações complexas

**Responsabilidades Específicas:**
- **Orquestração**: Coordena operações entre múltiplos domínios
- **Validação**: Valida dados de entrada usando validators
- **Mapeamento**: Converte DTOs ↔ Domain Objects
- **Logging**: Observabilidade e auditoria de operações
- **Error Handling**: Tratamento e propagação consistente de erros
- **Transaction Management**: Gerencia transações cross-domain

**Regras da Camada:**
- ✅ **PODE:** Importar e usar todos os módulos do core
- ✅ **PODE:** Definir DTOs, requests e responses específicos
- ✅ **PODE:** Coordenar múltiplas operações do core em uma transação
- ✅ **PODE:** Implementar validações de entrada complexas
- ✅ **PODE:** Fazer logging e observabilidade
- ✅ **PODE:** Gerenciar estado de aplicação (não de domínio)
- ❌ **NÃO PODE:** Conter lógica de negócio (deve estar no core)
- ❌ **NÃO PODE:** Acessar adapters diretamente (usar via DI)
- ❌ **NÃO PODE:** Conhecer detalhes de HTTP, SQL ou APIs externas
- ❌ **NÃO PODE:** Importar platform ou adapters

**Exemplo de Implementação:**
```go
// ✅ CORRETO - Application Use Case
type SessionService struct {
    sessionCore *session.Service     // Core domain service
    logger      *logger.Logger       // Observability
    validator   *validation.Validator // Input validation
}

func (s *SessionService) CreateSession(ctx context.Context, req *contracts.CreateSessionRequest) (*contracts.CreateSessionResponse, error) {
    // 1. Logging de entrada
    s.logger.InfoWithFields("Creating session", map[string]interface{}{
        "name": req.Name,
        "auto_connect": req.AutoConnect,
    })

    // 2. Validação de entrada
    if err := s.validator.ValidateStruct(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 3. Mapeamento DTO -> Domain
    coreReq := &session.CreateSessionRequest{
        Name:        req.Name,
        AutoConnect: req.AutoConnect,
        ProxyConfig: s.mapProxyConfig(req.ProxyConfig),
    }

    // 4. Delegação para Core Domain
    session, err := s.sessionCore.CreateSession(ctx, coreReq)
    if err != nil {
        s.logger.ErrorWithFields("Failed to create session", map[string]interface{}{
            "error": err.Error(),
            "name": req.Name,
        })
        return nil, fmt.Errorf("failed to create session: %w", err)
    }

    // 5. Mapeamento Domain -> DTO
    response := &contracts.CreateSessionResponse{
        ID:          session.ID.String(),
        Name:        session.Name,
        IsConnected: session.IsConnected,
        CreatedAt:   session.CreatedAt,
    }

    // 6. Logging de sucesso
    s.logger.InfoWithFields("Session created successfully", map[string]interface{}{
        "session_id": session.ID.String(),
        "name": session.Name,
    })

    return response, nil
}

// ✅ CORRETO - Coordenação multi-domain
func (s *SessionService) SendMessage(ctx context.Context, sessionID, to, content string) error {
    // 1. Validar sessão existe e está conectada
    session, err := s.sessionCore.GetByID(ctx, uuid.MustParse(sessionID))
    if err != nil {
        return fmt.Errorf("session not found: %w", err)
    }

    // 2. Coordenar com messaging domain
    return s.sessionCore.SendTextMessage(ctx, session.ID, to, content)
}
```

### 🔌 **ADAPTERS - Infrastructure Layer (Camada de Infraestrutura)**

**Responsabilidade:** Implementa as interfaces definidas no core domain, conecta o sistema com o mundo externo (banco de dados, APIs, protocolos de rede). Contém todos os detalhes de implementação específicos de tecnologia.

**Estrutura Detalhada:**
```
internal/adapters/
├── repository/                 # 💾 Data Persistence Implementations
│   ├── session_repository.go   # PostgreSQL session persistence
│   ├── message_repository.go   # PostgreSQL message persistence
│   ├── group_repository.go     # PostgreSQL group persistence
│   └── contact_repository.go   # PostgreSQL contact persistence
├── server/                     # 🌐 HTTP Server Infrastructure
│   ├── server.go               # HTTP server setup & lifecycle
│   ├── router/                 # Request routing
│   │   ├── router.go           # Main router setup
│   │   ├── session_routes.go   # Session endpoints
│   │   ├── message_routes.go   # Message endpoints
│   │   └── group_routes.go     # Group endpoints
│   ├── handlers/               # HTTP request handlers
│   │   ├── session_handler.go  # Session HTTP handlers
│   │   ├── message_handler.go  # Message HTTP handlers
│   │   └── chatwoot_handler.go # Chatwoot webhook handlers
│   ├── middleware/             # HTTP middleware
│   │   ├── auth.go             # API key authentication
│   │   ├── cors.go             # CORS handling
│   │   ├── logging.go          # Request logging
│   │   └── recovery.go         # Panic recovery
│   └── shared/                 # Shared HTTP utilities
│       ├── response.go         # Standard response formats
│       └── validation.go       # HTTP validation helpers
└── waclient/                   # 📱 WhatsApp Integration
    ├── gateway.go              # WhatsApp gateway implementation
    ├── client.go               # WhatsApp client management
    ├── events.go               # WhatsApp event handling
    ├── mapper.go               # WhatsApp ↔ Domain mapping
    └── validator.go            # WhatsApp data validation
```

### **Adapters Implementados (Infrastructure Layer)**

#### **1. Repository Adapters** 💾
**Implementações PostgreSQL para persistência:**

**SessionRepository:**
- **Interface**: `session.Repository` do core
- **Funcionalidades**: CRUD completo, queries otimizadas, error mapping
- **Features**: Connection pooling, prepared statements, transaction support
- **Tabela**: `zpSessions` com campos otimizados

**MessageRepository:**
- **Interface**: `messaging.Repository` do core
- **Funcionalidades**: Persistência de mensagens, sync status, Chatwoot mapping
- **Features**: Bulk operations, indexação otimizada, full-text search
- **Tabela**: `zpMessage` com relacionamento para sessões

#### **2. Server Adapters** 🌐
**Infraestrutura HTTP completa:**

**HTTP Server:**
- **Framework**: Chi router v5 (alta performance)
- **Features**: Graceful shutdown, middleware chain, CORS support
- **Authentication**: API key based com middleware customizado
- **Documentation**: Swagger/OpenAPI automático

**Handlers & Routes:**
- **Session Handlers**: CRUD de sessões, QR code, connection management
- **Message Handlers**: Send messages, get history, media upload
- **Group Handlers**: Group management, participant operations
- **Chatwoot Handlers**: Webhook processing, configuration

**Middleware Chain:**
- **Auth Middleware**: API key validation e context injection
- **CORS Middleware**: Cross-origin request handling
- **Logging Middleware**: Request/response logging estruturado
- **Recovery Middleware**: Panic recovery e error handling

#### **3. WhatsApp Client Adapter** 📱
**Integração completa com WhatsApp:**

**Gateway Implementation:**
- **Interface**: `session.WhatsAppGateway` do core
- **Library**: go.mau.fi/whatsmeow (oficial Meta)
- **Features**: Multi-session support, event handling, QR generation
- **Connection Management**: Pool de clientes, reconnection logic

**Client Management:**
- **MyClient**: Wrapper customizado do whatsmeow.Client
- **ClientManager**: Singleton para gerenciar múltiplas sessões
- **Event Processing**: Real-time event handling e propagação
- **State Management**: Persistent session state com SQLite

**Features Avançadas:**
- **QR Code Generation**: Geração automática para pareamento
- **Event Streaming**: Processamento assíncrono de eventos WhatsApp
- **Message Types**: Suporte completo a text, media, documents, interactive
- **Error Handling**: Retry logic, circuit breaker, graceful degradation

**Padrões Implementados:**
- **Adapter Pattern**: Adapta interfaces externas para contratos do core
- **Repository Pattern**: Implementações concretas de persistência
- **Gateway Pattern**: Implementações de integrações externas
- **Data Mapper**: Conversão entre Domain Objects e modelos externos
- **Connection Pooling**: Gerenciamento eficiente de conexões
- **Circuit Breaker**: Proteção contra falhas de serviços externos
- **Middleware Chain**: Pipeline de processamento de requests

**Responsabilidades Específicas:**

**Repository Adapters:**
- Implementam interfaces Repository do core
- Gerenciam conexões com banco de dados
- Fazem mapeamento Domain ↔ Database Models
- Tratam erros específicos de persistência
- Implementam queries otimizadas

**Server Adapters:**
- Implementam API REST com Chi router
- Gerenciam autenticação e autorização
- Fazem serialização/deserialização JSON
- Implementam middleware de logging e recovery
- Tratam erros HTTP específicos

**WhatsApp Client Adapter:**
- Implementa WhatsAppGateway interface
- Gerencia conexões com WhatsApp via whatsmeow
- Processa eventos em tempo real
- Mantém estado de sessões ativas
- Implementa retry logic e circuit breaker

**Regras Rigorosas:**
- ✅ **PODE:** Implementar interfaces definidas no core
- ✅ **PODE:** Usar qualquer framework ou biblioteca externa
- ✅ **PODE:** Conhecer detalhes de protocolos (HTTP, SQL, WebSocket)
- ✅ **PODE:** Importar services via dependency injection
- ✅ **PODE:** Fazer logging e métricas específicas
- ✅ **PODE:** Implementar cache, retry, circuit breaker
- ❌ **NÃO PODE:** Conter lógica de negócio (deve estar no core)
- ❌ **NÃO PODE:** Importar outros adapters diretamente
- ❌ **NÃO PODE:** Modificar entidades do core
- ❌ **NÃO PODE:** Tomar decisões de negócio

**Exemplo de Implementação:**
```go
// ✅ CORRETO - Repository Implementation
type SessionRepository struct {
    db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) session.Repository {
    return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, sess *session.Session) error {
    // 1. Mapeamento Domain -> Database Model
    model, err := r.toModel(sess)
    if err != nil {
        return fmt.Errorf("failed to convert session to model: %w", err)
    }

    // 2. SQL específico com prepared statement
    query := `
        INSERT INTO "zpSessions" (
            id, name, "deviceJid", "isConnected", "connectionError",
            "qrCode", "qrCodeExpiresAt", "proxyConfig", "createdAt",
            "updatedAt", "connectedAt", "lastSeen"
        ) VALUES (
            :id, :name, :deviceJid, :isConnected, :connectionError,
            :qrCode, :qrCodeExpiresAt, :proxyConfig, :createdAt,
            :updatedAt, :connectedAt, :lastSeen
        )
    `

    // 3. Execução com context
    _, err = r.db.NamedExecContext(ctx, query, model)
    if err != nil {
        // 4. Mapeamento de erros PostgreSQL -> Domain errors
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code {
            case "23505": // unique_violation
                if pqErr.Constraint == "zpSessions_name_key" {
                    return errors.ErrSessionAlreadyExists
                }
            }
        }
        return fmt.Errorf("failed to create session: %w", err)
    }

    return nil
}

// ✅ CORRETO - WhatsApp Gateway Implementation
type Gateway struct {
    container    *sqlstore.Container
    logger       *logger.Logger
    sessions     map[uuid.UUID]*MyClient
    mutex        sync.RWMutex
}

func (g *Gateway) CreateSession(ctx context.Context, sessionID uuid.UUID) error {
    // 1. Criar whatsmeow device store
    deviceStore := g.container.NewDevice()

    // 2. Inicializar WhatsApp client
    client := whatsmeow.NewClient(deviceStore, waLog.Noop)

    // 3. Setup event handlers
    myClient := NewMyClient(sessionID, client, g.db, g, g.logger)

    // 4. Registrar na sessão map (thread-safe)
    g.mutex.Lock()
    g.sessions[sessionID] = myClient
    g.mutex.Unlock()

    // 5. Iniciar processo de conexão
    if client.Store.ID == nil {
        g.handleQRCodePairing(myClient)
    } else {
        err := client.Connect()
        if err != nil {
            return fmt.Errorf("failed to connect existing session: %w", err)
        }
    }

    return nil
}

// ✅ CORRETO - HTTP Handler
type SessionHandler struct {
    *shared.BaseHandler
    sessionService *services.SessionService
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
    h.LogRequest(r, "create session")

    // 1. Parse request body
    var req contracts.CreateSessionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", "INVALID_JSON")
        return
    }

    // 2. Delegate to application service
    response, err := h.sessionService.CreateSession(r.Context(), &req)
    if err != nil {
        h.HandleServiceError(w, err)
        return
    }

    // 3. Write success response
    h.WriteSuccessResponse(w, http.StatusCreated, response)
}
```

### 🏗️ **PLATFORM - Platform Layer (Camada de Plataforma)**

**Responsabilidade:** Fornece infraestrutura transversal, configuração, logging, monitoramento, database management e dependency injection. Suporta todas as outras camadas com serviços de plataforma.

**Estrutura Detalhada:**
```
platform/
├── config/                     # 🔧 Configuration Management
│   ├── config.go               # Configuration structs & loading
│   └── validation.go           # Configuration validation
├── logger/                     # 📝 Structured Logging System
│   ├── logger.go               # Logger interface & implementation
│   └── fields.go               # Structured logging fields
├── database/                   # 💾 Database Infrastructure
│   ├── database.go             # Database connection & pooling
│   ├── migrator.go             # Database migration system
│   └── health.go               # Database health checks
└── container/                  # 🔗 Dependency Injection Container
    ├── container.go            # Main DI container
    └── wire.go                 # Dependency wiring (optional)
```

**Responsabilidades Específicas:**

**Configuration Management:**
- Carregamento de configurações via environment variables
- Validação de configurações na inicialização
- Hot-reload de configurações (quando aplicável)
- Configurações tipadas e type-safe

**Logging System:**
- Logging estruturado com zerolog
- Múltiplos outputs (console, file, syslog)
- Log levels configuráveis
- Context-aware logging com fields

**Database Infrastructure:**
- Connection pooling otimizado
- Sistema de migrações automáticas
- Health checks e monitoring
- Transaction management

**Dependency Injection:**
- Container centralizado para todas as dependências
- Inicialização ordenada de componentes
- Lifecycle management (start/stop)
- Interface-based dependency resolution

**Regras da Camada:**
- ✅ **PODE:** Configurar e inicializar todas as dependências
- ✅ **PODE:** Implementar cross-cutting concerns (logging, metrics)
- ✅ **PODE:** Usar qualquer biblioteca externa necessária
- ✅ **PODE:** Gerenciar lifecycle de componentes
- ✅ **PODE:** Implementar health checks e monitoring
- ❌ **NÃO PODE:** Conter lógica de negócio específica
- ❌ **NÃO PODE:** Conhecer detalhes específicos do domínio
- ❌ **NÃO PODE:** Implementar use cases ou regras de negócio

### 🚀 **CMD** - Entry Points

**Responsabilidade:** Pontos de entrada da aplicação.

**Estrutura:**
```
cmd/
├── server/                     # HTTP Server
├── worker/                     # Background workers
└── cli/                        # CLI tools
```

**Regras:**
- ✅ **PODE:** Inicializar aplicação
- ✅ **PODE:** Configurar dependency injection
- ✅ **PODE:** Importar qualquer camada
- ❌ **NÃO PODE:** Conter lógica de negócio
- ❌ **NÃO PODE:** Conter lógica de aplicação

## 🔄 Fluxo de Dependências

```
cmd → platform → adapters → services → core
 ↑                                        ↓
 └────────── interfaces ←─────────────────┘
```

**Regra de Ouro:** Dependências sempre fluem para dentro (core não depende de nada externo).

## 📋 Regras de Import

### ✅ **Imports Permitidos**

```go
// Core pode importar outros módulos do core
import "zpwoot/internal/core/session"
import "zpwoot/internal/core/messaging"

// UseCases pode importar
import "zpwoot/internal/core/session"
import "zpwoot/internal/core/messaging"

// Adapters pode importar
import "zpwoot/internal/core/session"
import "zpwoot/internal/usecases"
import "github.com/gin-gonic/gin"

// Platform pode importar
import "zpwoot/internal/adapters"
import "zpwoot/internal/usecases"
import "zpwoot/internal/core"

// CMD pode importar
import "zpwoot/platform"
import "zpwoot/internal/adapters"
import "zpwoot/internal/usecases"
import "zpwoot/internal/core"
```

### ❌ **Imports Proibidos**

```go
// Core NÃO pode importar
import "zpwoot/internal/usecases"        // ❌
import "zpwoot/internal/adapters"        // ❌
import "zpwoot/platform"                 // ❌
import "github.com/gin-gonic/gin"        // ❌

// UseCases NÃO pode importar
import "zpwoot/internal/adapters"        // ❌
import "zpwoot/platform"                 // ❌

// Adapters NÃO pode importar
import "zpwoot/platform"                 // ❌ (exceto para DI)
```

## 🧪 Estratégia de Testes

### **Testes Unitários**
- **Core:** Testa lógica de negócio isoladamente
- **UseCases:** Testa orquestração com mocks
- **Adapters:** Testa implementações específicas

### **Testes de Integração**
- Testa interação entre camadas
- Usa banco de dados de teste
- Testa APIs externas com mocks

### **Testes End-to-End**
- Testa fluxos completos
- Usa ambiente similar à produção

## 🔧 Dependency Injection

### **Container de DI**
```go
type Container struct {
    // Core
    SessionService *session.Service
    
    // UseCases
    SessionAppService *usecases.SessionService

    // Adapters
    SessionRepo session.Repository
    WhatsAppGW  session.WhatsAppGateway
}
```

### **Inicialização**
```go
func NewContainer() *Container {
    // 1. Criar adapters
    sessionRepo := postgres.NewSessionRepository(db)
    whatsappGW := whatsapp.NewGateway(client)
    
    // 2. Criar core services
    sessionCore := session.NewService(sessionRepo, whatsappGW)
    
    // 3. Criar application services
    sessionApp := usecases.NewSessionService(sessionCore)
    
    return &Container{...}
}
```

## 📊 Métricas de Qualidade

### **Acoplamento**
- Core: 0 dependências externas
- UseCases: Apenas core
- Adapters: Core + UseCases (via DI)

### **Coesão**
- Cada módulo tem responsabilidade única
- Funcionalidades relacionadas agrupadas

### **Testabilidade**
- 100% das interfaces mockáveis
- Lógica de negócio testável isoladamente

## 🚨 Violações Comuns

### ❌ **Anti-Patterns a Evitar**

1. **Core importando adapters**
```go
// ❌ ERRADO
import "zpwoot/adapters/database"
```

2. **Lógica de negócio em adapters**
```go
// ❌ ERRADO
func (h *HTTPHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
    if session.Status == "active" && session.LastSeen > time.Now() {
        // ❌ Lógica de negócio no adapter
    }
}
```

3. **UseCases acessando adapters diretamente**
```go
// ❌ ERRADO
func (s *SessionService) CreateSession() {
    db := postgres.Connect() // ❌ Acesso direto
}
```

## ✅ **Checklist de Conformidade**

- [ ] Core não importa nenhuma camada externa
- [ ] UseCases só importa core
- [ ] Adapters implementam interfaces do core
- [ ] Todas as dependências são injetadas
- [ ] Lógica de negócio está no core
- [ ] Testes cobrem todas as camadas
- [ ] Documentação está atualizada

## 🔍 Ferramentas de Validação

### **Linter Customizado**
```bash
# Verificar violações de arquitetura
make arch-lint

# Verificar imports proibidos
make import-check

# Verificar cobertura de testes por camada
make test-coverage
```

### **Scripts de Validação**
```bash
#!/bin/bash
# scripts/validate-architecture.sh

# Verificar se core não importa camadas externas
if grep -r "zpwoot/adapters\|zpwoot/services\|zpwoot/platform" core/; then
    echo "❌ ERRO: Core não pode importar camadas externas"
    exit 1
fi

echo "✅ Arquitetura válida"
```

## 📚 Exemplos Práticos

### **Adicionando Nova Feature**

1. **Criar módulo no core:**
```go
// core/newsletter/models.go
type Newsletter struct {
    ID      string
    Title   string
    Content string
}

// core/newsletter/service.go
type Service struct {
    repo Repository
}

// core/newsletter/contracts.go
type Repository interface {
    Save(ctx context.Context, newsletter *Newsletter) error
}
```

2. **Criar use case de aplicação:**
```go
// usecases/newsletter_usecase.go
type NewsletterService struct {
    newsletterCore *newsletter.Service
}
```

3. **Implementar adapters:**
```go
// adapters/database/newsletter_repository.go
type PostgresNewsletterRepository struct {
    db *sql.DB
}

// adapters/http/newsletter_handler.go
type NewsletterHandler struct {
    service *services.NewsletterService
}
```

### **Refatorando Código Existente**

1. **Identificar lógica de negócio**
2. **Extrair para core**
3. **Criar interfaces**
4. **Implementar adapters**
5. **Criar use cases**
6. **Atualizar DI**

## 🎓 Treinamento da Equipe

### **Onboarding Checklist**
- [ ] Ler este documento
- [ ] Entender fluxo de dependências
- [ ] Praticar com feature simples
- [ ] Code review com arquiteto
- [ ] Validar com ferramentas

### **Code Review Guidelines**
- Verificar imports
- Validar responsabilidades
- Checar testes
- Confirmar interfaces

---

**Esta arquitetura garante:**
- 🎯 **Foco no negócio** - lógica isolada e testável
- 🔧 **Flexibilidade** - fácil troca de implementações
- 🚀 **Performance** - sem overhead desnecessário
- 📈 **Escalabilidade** - crescimento organizado
- 🛡️ **Manutenibilidade** - código limpo e estruturado

## 🌐 **API REST - Endpoints Disponíveis**

### **Session Management** 📱
```
POST   /sessions                    # Criar nova sessão
GET    /sessions                    # Listar todas as sessões
GET    /sessions/{id}               # Obter sessão específica
PUT    /sessions/{id}               # Atualizar sessão
DELETE /sessions/{id}               # Deletar sessão
POST   /sessions/{id}/connect       # Conectar sessão
POST   /sessions/{id}/disconnect    # Desconectar sessão
GET    /sessions/{id}/qr            # Obter QR code
GET    /sessions/{id}/status        # Status da sessão
```

### **Message Operations** 💬
```
POST   /sessions/{id}/messages/text        # Enviar mensagem de texto
POST   /sessions/{id}/messages/media       # Enviar mídia
POST   /sessions/{id}/messages/document    # Enviar documento
POST   /sessions/{id}/messages/interactive # Enviar mensagem interativa
GET    /sessions/{id}/messages             # Histórico de mensagens
PUT    /sessions/{id}/messages/{msgId}     # Editar mensagem
DELETE /sessions/{id}/messages/{msgId}     # Revogar mensagem
POST   /sessions/{id}/messages/{msgId}/react # Reagir à mensagem
```

### **Group Management** 👥
```
POST   /sessions/{id}/groups               # Criar grupo
GET    /sessions/{id}/groups               # Listar grupos
GET    /sessions/{id}/groups/{groupId}     # Obter grupo específico
PUT    /sessions/{id}/groups/{groupId}     # Atualizar grupo
DELETE /sessions/{id}/groups/{groupId}     # Deletar grupo
POST   /sessions/{id}/groups/{groupId}/participants # Adicionar participante
DELETE /sessions/{id}/groups/{groupId}/participants/{jid} # Remover participante
POST   /sessions/{id}/groups/{groupId}/admins/{jid} # Promover a admin
```

### **Contact Operations** 📞
```
GET    /sessions/{id}/contacts             # Listar contatos
GET    /sessions/{id}/contacts/{jid}       # Obter contato específico
POST   /sessions/{id}/contacts/check       # Verificar número no WhatsApp
GET    /sessions/{id}/contacts/{jid}/avatar # Obter avatar do contato
```

### **Chatwoot Integration** 🔗
```
POST   /sessions/{id}/chatwoot/config      # Configurar integração
GET    /sessions/{id}/chatwoot/config      # Obter configuração
PUT    /sessions/{id}/chatwoot/config      # Atualizar configuração
DELETE /sessions/{id}/chatwoot/config      # Remover integração
POST   /chatwoot/webhook                   # Webhook do Chatwoot
```

### **System Endpoints** ⚙️
```
GET    /health                             # Health check
GET    /swagger/*                          # Documentação API
```

### **Authentication** 🔐
Todas as rotas (exceto `/health`, `/swagger`, `/chatwoot/webhook`) requerem autenticação via:
- **Header**: `Authorization: Bearer {api_key}` ou `X-API-Key: {api_key}`
- **API Key**: Configurada via environment variable `API_KEY`

## 📊 **Métricas de Qualidade Arquitetural Atual**

### **Conformidade com Clean Architecture**
| Aspecto | Status | Observação |
|---------|--------|------------|
| **Dependency Rule** | ✅ Excelente | Dependências fluem sempre para dentro |
| **Interface Segregation** | ✅ Muito Bom | Interfaces específicas e coesas |
| **Single Responsibility** | ✅ Muito Bom | Cada camada tem responsabilidade clara |
| **Open/Closed Principle** | ✅ Bom | Extensível via interfaces |
| **Testability** | ✅ Excelente | 100% das interfaces mockáveis |

### **Acoplamento e Coesão**
- **Core**: 0 dependências externas ✅
- **Services**: Apenas core + platform ✅
- **Adapters**: Core + Services via DI ✅
- **Platform**: Independente de domínio ✅

### **Modularidade por Domínio**
- **Session Domain**: Gerenciamento completo de sessões WhatsApp
- **Messaging Domain**: Sistema de mensagens com sync Chatwoot
- **Group Domain**: Gerenciamento de grupos WhatsApp
- **Contact Domain**: Gerenciamento de contatos
- **Shared Domain**: Erros, eventos e tipos compartilhados

## 🔄 **Padrões Arquiteturais Implementados**

### **1. Clean Architecture**
- **Dependency Rule**: Dependências fluem sempre para dentro
- **Interface Segregation**: Interfaces específicas e coesas
- **Dependency Inversion**: Abstrações não dependem de detalhes

### **2. Domain-Driven Design (DDD)**
- **Bounded Contexts**: Módulos bem definidos (session, messaging, group, contact)
- **Entities**: Objetos com identidade e ciclo de vida
- **Value Objects**: Objetos imutáveis (ProxyConfig, SessionStatus)
- **Domain Services**: Lógica que não pertence a entidades específicas
- **Repository Pattern**: Abstração de persistência

### **3. Application Service Pattern**
- **Use Case Orchestration**: Services coordenam operações complexas
- **DTO Pattern**: Contratos de entrada/saída específicos
- **Validation Pattern**: Validação centralizada e reutilizável
- **Mapping Pattern**: Conversão entre DTOs e Domain Objects

### **4. Infrastructure Patterns**
- **Adapter Pattern**: Adapta tecnologias externas para interfaces do core
- **Gateway Pattern**: Integrações com sistemas externos (WhatsApp, Chatwoot)
- **Repository Pattern**: Implementações concretas de persistência
- **Data Mapper**: Conversão Domain ↔ Database Models

## 🚀 **Stack Tecnológica Atual**

### **Core Technologies**
- **Language**: Go 1.21+
- **HTTP Router**: Chi v5 (alta performance)
- **Database**: PostgreSQL com SQLX
- **WhatsApp**: go.mau.fi/whatsmeow (oficial)
- **Logging**: Zerolog (structured logging)

### **Development & Operations**
- **Hot Reload**: Air para desenvolvimento
- **Documentation**: Swagger/OpenAPI automático
- **Database Migrations**: Sistema embarcado
- **Docker**: Ambiente completo com docker-compose
- **Testing**: Testify para testes unitários

## ⚙️ **Configuração e Deployment**

### **Environment Variables**
```bash
# Server Configuration
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30
SERVER_IDLE_TIMEOUT=120

# Database Configuration
DATABASE_URL=postgres://user:pass@localhost:5432/zpwoot
DATABASE_MAX_OPEN_CONNS=25
DATABASE_MAX_IDLE_CONNS=5
DATABASE_CONN_MAX_LIFETIME=300

# Security
API_KEY=your-secure-api-key-here

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# WhatsApp Configuration
WHATSAPP_LOG_LEVEL=ERROR
WHATSAPP_STORE_PATH=./store

# Chatwoot Integration (Optional)
CHATWOOT_BASE_URL=https://your-chatwoot.com
CHATWOOT_API_ACCESS_TOKEN=your-token
```

### **Docker Deployment**
```yaml
# docker-compose.yml
version: '3.8'
services:
  zpwoot:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://zpwoot:password@postgres:5432/zpwoot
      - API_KEY=your-secure-api-key
    depends_on:
      - postgres
    volumes:
      - ./store:/app/store

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: zpwoot
      POSTGRES_USER: zpwoot
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### **Production Deployment**
```bash
# Build
make build

# Run migrations
./zpwoot migrate

# Start server
./zpwoot server

# Health check
curl http://localhost:8080/health
```

### **Integrations**
- **Chatwoot**: Integração bidirecional opcional
- **Webhook Support**: Sistema de webhooks flexível
- **API Authentication**: API Key based

## 📈 **Roadmap de Evolução Arquitetural**

### **Fase 1: Performance Optimizations (Próxima)**
- [ ] Object pooling para reduzir GC pressure
- [ ] Multi-level caching (L1 in-memory + L2 Redis)
- [ ] Batch operations para alta throughput
- [ ] Connection pooling otimizado

### **Fase 2: Observability & Monitoring**
- [ ] Distributed tracing com Jaeger
- [ ] Métricas com Prometheus
- [ ] Health checks avançados
- [ ] Performance profiling

### **Fase 3: Scalability Enhancements**
- [ ] Event-driven architecture com NATS/Kafka
- [ ] CQRS para separação read/write
- [ ] Horizontal scaling support
- [ ] Load balancing strategies

### **Fase 4: Additional Protocols**
- [ ] gRPC API para alta performance
- [ ] WebSocket para real-time updates
- [ ] GraphQL para flexible queries
- [ ] Message queues para async processing

## 🎯 **Conclusão**

O **zpwoot** implementa uma arquitetura exemplar que combina:
- **Clean Architecture** para separação rigorosa de responsabilidades
- **DDD** para modelagem rica de domínio
- **Padrões modernos** para alta qualidade de código
- **Tecnologias robustas** para performance e confiabilidade

A arquitetura atual está **sólida e bem estruturada**, pronta para evoluir com otimizações de performance mantendo a excelente qualidade arquitetural existente.

---

**Versão:** 2.0
**Última atualização:** 2025-01-04
**Responsável:** Equipe zpwoot
