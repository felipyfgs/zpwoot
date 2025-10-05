# zpwoot - WhatsApp API Gateway Architecture

## Overview

zpwoot is a comprehensive WhatsApp Business API built with Go, following Clean Architecture principles. It provides endpoints for session management, messaging, contacts, groups, media handling, and integrations with Chatwoot.

## Architecture Principles

- **Clean Architecture**: Separation of concerns with clear boundaries between layers
- **Domain-Driven Design**: Business logic encapsulated in domain entities and services
- **Dependency Injection**: Loose coupling through dependency inversion
- **Database Migrations**: Automated schema management with version control
- **Configuration Management**: Environment-based configuration with sensible defaults

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    HTTP Adapters                                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Chi Router    │  │   Middleware    │  │   Handlers      │ │
│  │   (routing)     │  │   (CORS, auth)  │  │  (REST API)     │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Application Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Use Cases     │  │      DTOs       │  │   Interfaces    │ │
│  │ (CreateSession) │  │  (SessionDTO)   │  │ (for adapters)  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Domain Layer                               │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │    Entities     │  │    Services     │  │  Repository     │ │
│  │   (Session)     │  │  (Domain Svc)   │  │  (Interfaces)   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                ▲
                                │ (implements)
┌─────────────────────────────────────────────────────────────────┐
│                    Database Adapters                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   PostgreSQL    │  │   Migrations    │  │   Repository    │ │
│  │  (Connection)   │  │   (Auto-run)    │  │ (Implementation)│ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    External Adapters                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   WhatsApp      │  │   Chatwoot      │  │   Webhooks      │ │
│  │  (whatsmeow)    │  │   (HTTP API)    │  │  (HTTP calls)   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘

Infrastructure Adapters (internal/adapters/):
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│     Config      │  │     Logger      │  │   Container     │
│  (Environment)  │  │  (Structured)   │  │ (Bootstrap/DI)  │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

## Project Structure

```
zpwoot/
├── cmd/                          # Application entry points
│   └── zpwoot/
│       └── main.go              # Main application entry point
├── internal/                     # Private application code
│   ├── domain/                  # Domain layer (business logic)
│   │   ├── session/            # Session aggregate
│   │   │   ├── entity.go       # Session entity
│   │   │   ├── repository.go   # Repository interface
│   │   │   └── service.go      # Domain service
│   │   ├── message/            # Message aggregate (future)
│   │   ├── webhook/            # Webhook aggregate (future)
│   │   └── shared/             # Shared domain concepts
│   │       └── errors.go       # Domain errors
│   ├── application/             # Application layer (use cases)
│   │   ├── session/            # Session use cases
│   │   │   ├── create.go       # Create session use case
│   │   │   ├── list.go         # List sessions use case
│   │   │   └── connect.go      # Connect session use case
│   │   ├── message/            # Message use cases (future)
│   │   └── dto/                # Data Transfer Objects
│   │       └── session.go      # Session DTOs
│   └── adapters/                # Infrastructure adapters
│       ├── database/           # Database implementations
│       │   ├── connection.go   # Database connection management
│       │   ├── migrator.go     # Database migration system
│       │   ├── migrations/     # SQL migration files
│       │   │   ├── 001_initial_schema.up.sql
│       │   │   ├── 001_initial_schema.down.sql
│       │   │   ├── 002_add_indexes.up.sql
│       │   │   └── 002_add_indexes.down.sql
│       │   └── repository/     # Repository implementations
│       │       ├── session.go  # Session repository implementation
│       │       ├── webhook.go  # Webhook repository implementation
│       │       ├── chatwoot.go # Chatwoot repository implementation
│       │       └── message.go  # Message repository implementation
│       ├── http/               # HTTP adapters
│       │   ├── router/         # HTTP router
│       │   │   └── router.go   # HTTP router setup
│       │   ├── middleware/     # HTTP middleware
│       │   └── handlers/       # HTTP handlers
│       │       └── session.go  # Session HTTP handlers
│       ├── whatsapp/           # WhatsApp client adapter (future)
│       ├── chatwoot/           # Chatwoot client adapter (future)
│       ├── config/             # Configuration adapter
│       │   └── config.go       # Environment configuration loader
│       ├── container/          # Dependency injection container
│       │   └── container.go    # DI container and bootstrap
│       └── logger/             # Logging adapter
│           └── logger.go       # Structured logging implementation
└── docs/                        # Documentation
    └── api/                    # API documentation (future)
```

## Layer Architecture

### 1. Domain Layer (`internal/domain/`)
- **Purpose**: Pure business logic and domain entities
- **Dependencies**: None (completely isolated)
- **Rules**:
  - No external dependencies
  - No framework imports
  - Only standard library and domain code
- **Components**:
  - **Entities**: Core business objects (Session, Message, etc.)
  - **Services**: Domain business logic
  - **Repositories**: Data access interfaces (contracts only)
  - **Shared**: Common domain concepts and errors

### 2. Application Layer (`internal/application/`)
- **Purpose**: Orchestrates domain objects and implements use cases
- **Dependencies**: Domain layer only
- **Rules**:
  - Can import domain layer
  - Cannot import adapters or external frameworks
  - Defines interfaces for external dependencies
- **Components**:
  - **Use Cases**: Application-specific business rules
  - **DTOs**: Data transfer objects for API boundaries
  - **Interfaces**: Contracts for external services

### 3. Adapters Layer (`internal/adapters/`)
- **Purpose**: Implements external integrations and infrastructure
- **Dependencies**: Application and Domain layers
- **Rules**:
  - Implements interfaces defined in application layer
  - Contains all external framework code
  - Handles data transformation between external and internal formats
- **Components**:
  - **Database**: Connection management, migrations, repository implementations
  - **HTTP**: Web server, routes, handlers, middleware
  - **External APIs**: WhatsApp, Chatwoot, etc.
  - **Config**: Environment configuration loading (infrastructure concern)
  - **Container**: Dependency injection and bootstrap (infrastructure concern)
  - **Logger**: Logging implementation (infrastructure concern)

## Design Decisions

### Why Config, Container, and Logger are in Adapters?

#### **Config (`internal/adapters/config/`)**
- **Justification**: Configuration loading is an infrastructure concern
- **Responsibility**: Reads from environment variables, files, external config services
- **Nature**: Adapter that translates external configuration into internal structures
- **Dependencies**: External libraries (godotenv), OS environment

#### **Container (`internal/adapters/container/`)**
- **Justification**: Dependency injection is infrastructure/bootstrap concern
- **Responsibility**: Wires up all dependencies, manages application lifecycle
- **Nature**: Orchestrates the entire application, knows about all layers
- **Dependencies**: All layers (by design, as it's the composition root)

#### **Logger (`internal/adapters/logger/`)**
- **Justification**: Logging implementation is infrastructure concern
- **Responsibility**: Implements structured logging, handles log output destinations
- **Nature**: Adapter that translates internal log calls to external logging libraries
- **Dependencies**: External libraries (zerolog), OS file system

### Alternative Considered: `pkg/` directory
- **Rejected because**: `pkg/` implies public packages that could be imported by external projects
- **Our components**: Are application-specific infrastructure, not reusable libraries
- **Clean Architecture**: Infrastructure belongs in adapters layer

## Repository Pattern Implementation

### Database Structure (`internal/adapters/database/`)

#### **Connection Management (`connection.go`)**
- Database connection setup and configuration
- Connection pooling and health checks
- Transaction management utilities

#### **Migration System (`migrator.go`)**
- Automatic migration execution on startup
- Version tracking and rollback support
- Embedded SQL files management

#### **Repository Implementations (`repository/`)**
Each repository file implements the corresponding domain repository interface:

- **`session.go`**: Implements `domain.SessionRepository`
  - CRUD operations for sessions
  - Connection status management
  - QR code handling

- **`webhook.go`**: Implements `domain.WebhookRepository`
  - Webhook configuration management
  - Event subscription handling

- **`chatwoot.go`**: Implements `domain.ChatwootRepository`
  - Chatwoot integration settings
  - Account and inbox management

- **`message.go`**: Implements `domain.MessageRepository`
  - Message synchronization tracking
  - WhatsApp ↔ Chatwoot mapping

### Repository Pattern Benefits
- **Testability**: Domain layer doesn't depend on database implementation
- **Flexibility**: Can easily switch database implementations
- **Separation**: Database logic separated from business logic
- **Interface Segregation**: Each repository has focused responsibilities

## Database Schema

### Core Tables

#### zpSessions
- **Purpose**: WhatsApp session management
- **Key Fields**: id, name, deviceJid, isConnected, qrCode
- **Features**: Auto-timestamps, connection tracking, QR code management

#### zpWebhooks
- **Purpose**: Event notification configuration
- **Key Fields**: sessionId, url, events, enabled
- **Features**: Per-session or global webhooks

#### zpChatwoot
- **Purpose**: Chatwoot integration configuration
- **Key Fields**: sessionId, url, token, accountId, inboxId
- **Features**: One-to-one with sessions, rich configuration options

#### zpMessage
- **Purpose**: WhatsApp ↔ Chatwoot message mapping
- **Key Fields**: zpMessageId, cwMessageId, syncStatus
- **Features**: Bidirectional sync tracking

#### zpMigrations
- **Purpose**: Database migration tracking
- **Key Fields**: version, name, appliedAt
- **Features**: Automatic migration management

## Configuration Management

### Environment Variables
```bash
# Application
PORT=8080
SERVER_HOST=0.0.0.0
LOG_LEVEL=info
ZP_API_KEY=your-api-key

# Database
DATABASE_URL=postgres://user:pass@host:port/db

# PostgreSQL (Docker)
POSTGRES_DB=zpwoot
POSTGRES_USER=zpwoot
POSTGRES_PASSWORD=zpwoot123

# WhatsApp
WA_LOG_LEVEL=INFO

# Webhooks
GLOBAL_WEBHOOK_URL=https://your-domain.com/webhooks
```

## Database Migrations

### Automatic Migration System
- **Embedded Migrations**: SQL files embedded in binary
- **Version Control**: Sequential numbering (001, 002, ...)
- **Automatic Execution**: Runs on application startup
- **Rollback Support**: Down migrations for rollback capability
- **Transaction Safety**: Each migration runs in a transaction

### Migration Files
- `XXX_name.up.sql`: Forward migration
- `XXX_name.down.sql`: Rollback migration
- Embedded using `//go:embed migrations`

## Dependency Injection

### Container Pattern
```go
type Container struct {
    config   *config.Config
    logger   *logger.Logger
    database *database.Database
    migrator *database.Migrator
}
```

### Initialization Flow
```
main.go
   │
   ├─► Load Config (.env)
   │
   ├─► Initialize Logger (zerolog)
   │
   ├─► Create Container (DI)
   │
   ├─► Container.Initialize()
   │    │
   │    ├─► Connect Database (PostgreSQL)
   │    │
   │    ├─► Create Migrator
   │    │
   │    ├─► Run Migrations (automatic)
   │    │    │
   │    │    ├─► Create zpMigrations table
   │    │    ├─► Load embedded SQL files
   │    │    ├─► Check applied migrations
   │    │    └─► Execute pending migrations
   │    │
   │    └─► Initialize Services
   │
   ├─► Setup HTTP Router (Chi)
   │
   ├─► Start HTTP Server
   │
   └─► Wait for Shutdown Signal
```

## HTTP Layer

### Router Configuration
- **Framework**: Chi router
- **Middleware**: CORS, logging, recovery, request ID
- **Endpoints**:
  - `GET /` - API information
  - `GET /health` - Health check with database verification

### Future API Structure (Planned)
```
/api/v1/
├── /sessions/          # Session management
├── /messages/          # Message operations
├── /contacts/          # Contact management
├── /groups/            # Group operations
├── /media/             # Media handling
├── /webhooks/          # Webhook management
└── /chatwoot/          # Chatwoot integration
```

## Development Workflow

### Local Development
```bash
# Start services
make up

# Run application
go run cmd/zpwoot/main.go
# or
make dev

# Build
make build
```

### Docker Development
```bash
# Development services
docker-compose -f docker-compose.dev.yml up -d

# Full application
docker-compose up -d
```

## Key Features

### 1. Automatic Database Migrations
- Migrations run automatically on startup
- No manual intervention required
- Version tracking and rollback support

### 2. Clean Architecture
- Clear separation of concerns
- Testable business logic
- Dependency inversion

### 3. Configuration Management
- Environment-based configuration
- Sensible defaults
- Docker-friendly

### 4. Structured Logging
- JSON structured logs
- Configurable log levels
- Request tracing

### 5. Health Monitoring
- Database health checks
- Application status endpoints
- Docker health checks

## Technology Stack

- **Language**: Go 1.24+
- **Web Framework**: Chi router
- **Database**: PostgreSQL 15+
- **Logging**: Zerolog
- **Configuration**: Environment variables with godotenv
- **Containerization**: Docker & Docker Compose
- **Development**: Air (hot reload), golangci-lint

## Refactoring Plan

### Phase 1: Restructure Directories
```bash
# Current → Target
platform/database/database.go    → internal/adapters/database/connection.go
platform/database/migrator.go    → internal/adapters/database/migrator.go
platform/database/migrations/    → internal/adapters/database/migrations/
platform/container/              → internal/adapters/container/
platform/config/                 → internal/adapters/config/
platform/logger/                 → internal/adapters/logger/
internal/infra/http/              → internal/adapters/http/
internal/domain/session/errors.go → internal/domain/shared/errors.go

# New directories to create
internal/adapters/database/repository/session.go
internal/adapters/database/repository/webhook.go
internal/adapters/database/repository/chatwoot.go
internal/adapters/database/repository/message.go
```

### Phase 2: Implement Clean Boundaries
1. **Move domain errors to shared**
2. **Create application layer with use cases**
3. **Implement repository pattern correctly**
4. **Add DTOs for API boundaries**
5. **Create proper dependency injection**

### Phase 3: Add Missing Components
1. **Repository implementations**
   - Implement domain repository interfaces
   - Add proper error handling and transactions
   - Include query optimization and indexing

2. **HTTP handlers with proper error handling**
   - RESTful API endpoints
   - Request validation and sanitization
   - Proper HTTP status codes and error responses

3. **Application use cases**
   - Session management use cases
   - Message handling use cases
   - Webhook configuration use cases

4. **Proper middleware stack**
   - Authentication and authorization
   - Request logging and tracing
   - Rate limiting and CORS

## Architecture Rules

### ✅ **Allowed Dependencies**
```
cmd/           → internal/adapters (container, config only)
adapters/      → internal/application, internal/domain, external libraries
application/   → internal/domain (pure business logic only)
domain/        → NONE (pure business logic, standard library only)
```

### ❌ **Forbidden Dependencies**
```
domain/        → adapters, application, external frameworks
application/   → adapters (must use interfaces)
adapters/      → other adapters (should be independent)
```

### 📁 **Directory Responsibilities**

#### `internal/domain/`
- ✅ Business entities and value objects
- ✅ Domain services (business rules)
- ✅ Repository interfaces (contracts)
- ✅ Domain events
- ❌ Database implementations
- ❌ HTTP handlers
- ❌ External API calls

#### `internal/application/`
- ✅ Use cases (application business rules)
- ✅ DTOs for API boundaries
- ✅ Interfaces for external services
- ✅ Application services
- ❌ Database implementations
- ❌ HTTP routing
- ❌ External API implementations

#### `internal/adapters/`
- ✅ Database connection and migrations
- ✅ Repository implementations (session, webhook, chatwoot, message)
- ✅ HTTP handlers and routing
- ✅ External API clients
- ✅ Configuration loading
- ✅ Dependency injection container
- ✅ Logging implementation
- ✅ Message queue implementations
- ❌ Business logic
- ❌ Use cases
- ❌ Domain entities

## Future Enhancements

### Planned Features
1. **WhatsApp Integration**: whatsmeow client integration
2. **Message Handling**: Send/receive messages
3. **Media Support**: Image, audio, video, document handling
4. **Contact Management**: Contact sync and management
5. **Group Operations**: Group creation and management
6. **Webhook System**: Event notifications
7. **Chatwoot Integration**: Full bidirectional sync
8. **API Documentation**: Swagger/OpenAPI integration
9. **Authentication**: API key and JWT support
10. **Rate Limiting**: Request throttling
11. **Metrics**: Prometheus metrics
12. **Testing**: Comprehensive test suite

## Getting Started

1. **Clone the repository**
2. **Copy `.env.example` to `.env`**
3. **Start services**: `make up`
4. **Run application**: `go run cmd/zpwoot/main.go`
5. **Access API**: `http://localhost:8080`

The application will automatically:
- Load configuration
- Connect to database
- Run migrations
- Start HTTP server

## Contributing

Follow the established architecture patterns:
- Keep domain logic pure
- Use dependency injection
- Write migrations for schema changes
- Add proper logging
- Follow Go conventions
