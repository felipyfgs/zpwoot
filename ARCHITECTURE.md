# zpwoot - Clean Architecture Implementation

## Overview

zpwoot is a WhatsApp Business API built with Go following **Clean Architecture** principles. It provides session management and messaging capabilities with a clean separation between business logic and infrastructure.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP HANDLERS                         │
│                    (Adapters - Input)                      │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │          Uses input.SessionUseCases                     │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    INPUT PORTS                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  SessionUseCases, MessageUseCases (Interfaces)         │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              ▲
                              │ implements
┌─────────────────────────────────────────────────────────────┐
│                      USE CASES                              │
│                   (Application Layer)                       │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │    CreateUseCase, ConnectUseCase, etc.                 │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ uses
┌─────────────────────────────────────────────────────────────┐
│                    OUTPUT PORTS                             │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  WhatsAppClient, Logger (Interfaces)                   │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              ▲
                              │ implements
┌─────────────────────────────────────────────────────────────┐
│                      ADAPTERS                               │
│                   (Infrastructure)                          │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  WhatsAppAdapter, LoggerAdapter, DatabaseAdapter       │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                        DOMAIN                               │
│                   (Business Logic)                          │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │    Session Entity, Domain Services, Repository         │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
zpwoot/
├── cmd/zpwoot/main.go                     # Application entry point
├── internal/
│   ├── core/                             # CORE LAYER
│   │   ├── domain/                       # Business Logic (Pure)
│   │   │   ├── session/
│   │   │   │   ├── entity.go            # Session entity
│   │   │   │   ├── repository.go        # Repository interface
│   │   │   │   └── service.go           # Domain service
│   │   │   └── shared/errors.go         # Domain errors
│   │   ├── application/                  # Use Cases
│   │   │   ├── dto/                     # Data Transfer Objects
│   │   │   │   ├── responses.go         # API response structures
│   │   │   │   ├── session.go           # Session DTOs
│   │   │   │   └── message.go           # Message DTOs
│   │   │   ├── usecase/                 # Use cases implementation
│   │   │   ├── session/             # Session use cases
│   │   │   └── message/             # Message use cases
│   │   └── validators/              # Input validation
│   │   └── ports/                       # Interfaces (Hexagonal Architecture)
│   │       ├── input/                   # Use case interfaces
│   │       └── output/                  # Dependency interfaces
│   ├── adapters/                        # INFRASTRUCTURE
│   │   ├── database/                    # PostgreSQL + migrations
│   │   │   └── repository/
│   │   │       └── adapter.go          # Repository adapter
│   │   ├── http/                        # REST API layer
│   │   │   ├── handlers/               # HTTP request handlers
│   │   │   ├── middleware/             # HTTP middleware (auth, CORS, etc.)
│   │   │   └── router/                 # Route configuration
│   │   ├── logger/
│   │   │   └── adapter.go              # Logger adapter
│   │   └── waclient/                    # WhatsApp integration
│   │       ├── adapter.go              # WhatsApp adapter (implements port)
│   │       ├── manager.go              # Session management
│   │       ├── events.go               # Event handling
│   │       ├── messages.go             # Message sending
│   │       ├── qr.go                   # QR code management
│   │       └── types.go                # Type definitions
│   ├── config/config.go                 # Configuration
│   └── container/container.go           # Dependency Injection Container
├── docs/
│   ├── API.md                          # API documentation
│   └── swagger/                        # OpenAPI/Swagger documentation
│       ├── docs.go                     # Generated Go documentation
│       ├── swagger.json                # OpenAPI JSON specification
│       └── swagger.yaml                # OpenAPI YAML specification
└── [docker files, etc.]
```

## Layers

### 1. Domain (`internal/core/domain/`)
- **Pure business logic** - zero external dependencies
- **Session entity** with business rules
- **Repository interface** (contract only)
- **Domain errors** and value objects

### 2. Application (`internal/core/application/`)
- **Use cases** that orchestrate domain objects
- **DTOs** for data transfer
- **Validators** for input validation
- **Depends only on**: Domain + Ports

### 3. Ports (`internal/core/ports/`)
- **Interfaces only** - no implementations
- **Input ports**: Use case interfaces (for HTTP handlers)
- **Output ports**: Dependency interfaces (WhatsApp, Logger, etc.)

### 4. Adapters (`internal/adapters/`)
- **Infrastructure implementations**
- **HTTP handlers** (implement input ports)
- **Database repository** (implements domain interface)
- **WhatsApp client** (implements output port)
- **Logger** (implements output port)

### 5. Container (`internal/container/`)
- **Dependency Injection Container** - Single responsibility
- **Lifecycle management** (start/stop)
- **Wires dependencies** together at startup
- **Simple getters** for accessing dependencies
- **No business logic** - pure DI container

## Key Principles

### Clean Architecture
- **Dependency Rule**: Dependencies point inward (toward domain)
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: High-level modules don't depend on low-level modules

### Hexagonal Architecture (Ports & Adapters)
- **Input Ports**: Use case interfaces (for HTTP handlers to use)
- **Output Ports**: Dependency interfaces (for external services)
- **Adapters**: Implement port interfaces

## Container Responsibilities

The **Dependency Injection Container** (`internal/container/`) has a **single responsibility**:

### ✅ What the Container DOES:
- **Dependency Injection**: Wires all dependencies together
- **Lifecycle Management**: Handles application start/stop
- **Configuration**: Manages app configuration
- **Simple Getters**: Provides access to initialized dependencies
- **Infrastructure Setup**: Database, logger, migrations

### ❌ What the Container DOES NOT:
- **Business Logic**: No domain rules or use case logic
- **HTTP Handling**: No request/response processing
- **Data Processing**: No message or session processing
- **External Communication**: No direct WhatsApp/API calls
- **Complex Abstractions**: No unnecessary interfaces or patterns

### Container Structure:
```go
type Container struct {
    // Infrastructure only
    config   *config.Config
    logger   *logger.Logger
    database *database.Database
    migrator *database.Migrator

    // Wired dependencies
    sessionService  *domainSession.Service
    whatsappClient  output.WhatsAppClient
    sessionUseCases input.SessionUseCases
    messageUseCases input.MessageUseCases
}
```

## Dependency Rules

```
✅ Allowed:
cmd/                       → internal/container, internal/config
internal/container/        → internal/core, internal/adapters, internal/config
internal/adapters/         → internal/core (domain, application, ports)
internal/core/application/ → internal/core/domain, internal/core/ports
internal/core/ports/       → internal/core/domain, internal/core/application/dto
internal/core/domain/      → NONE (pure business logic, stdlib only)

❌ Forbidden:
internal/core/domain/      → adapters, application, ports, external frameworks
internal/core/application/ → adapters (must use ports interfaces)
internal/core/ports/       → adapters (only interfaces, no implementations)
internal/container/        → business logic (only DI and lifecycle)
```

## Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Chi router with middleware
- **Database**: PostgreSQL with automated migrations
- **Logging**: Zerolog (structured logging)
- **WhatsApp**: whatsmeow library
- **Configuration**: Environment variables with validation
- **Documentation**: Swagger/OpenAPI 3.0
- **Security**: API key authentication, CORS, security headers
- **Development**: Air (hot reload), Docker Compose
- **Architecture**: Clean Architecture + Hexagonal Architecture

## Getting Started

1. **Clone and setup**:
   ```bash
   git clone <repo>
   cp .env.example .env
   ```

2. **Start services**:
   ```bash
   make up      # Start PostgreSQL
   make dev     # Start app with hot reload
   ```

3. **API available at**: `http://localhost:8080`

## Development Commands

```bash
# Development
make dev         # Start with hot reload
make run         # Start normally
make build       # Build binary

# Database
make migrate     # Run migrations
make db-reset    # Reset database

# Documentation
make swagger     # Generate API docs
make docs        # View documentation

# Quality
make test        # Run tests
make lint        # Run linter
make fmt         # Format code

# Docker
make up          # Start services
make down        # Stop services
```

## Security Features

### Authentication & Authorization
- **API Key Authentication**: Required for all endpoints
- **Environment-based Configuration**: No hardcoded credentials
- **Secure Random Generation**: Crypto-secure random for IDs

### HTTP Security
- **CORS Configuration**: Restrictive origins in production
- **Security Headers**: XSS protection, content sniffing prevention
- **HTTPS Support**: TLS configuration ready
- **Request Validation**: Input sanitization and validation

### Database Security
- **SSL Connections**: Required in production
- **Connection Pooling**: Secure connection management
- **Migration Safety**: Automated, reversible migrations

### Documentation
- **Swagger UI**: Interactive API documentation
- **Security Schemas**: API key authentication documented
- **Clean Model Names**: Professional API documentation