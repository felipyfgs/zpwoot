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
│   │   │   ├── usecase/session/         # Session use cases
│   │   │   │   └── container.go         # Use cases container
│   │   │   └── validators/              # Input validation
│   │   └── ports/                       # Interfaces (Hexagonal Architecture)
│   │       ├── input/                   # Use case interfaces
│   │       └── output/                  # Dependency interfaces
│   ├── adapters/                        # INFRASTRUCTURE
│   │   ├── database/                    # PostgreSQL + migrations
│   │   │   └── repository/
│   │   │       └── adapter.go          # Repository adapter
│   │   ├── http/                        # REST API handlers
│   │   │   └── router/
│   │   │       └── routes.go           # HTTP routes configuration
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
│   └── container/container.go           # Dependency injection
├── docs/swagger/
│   └── spec.go                         # Swagger specification
├── examples/
│   └── logger.go                       # Logger usage example
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
- **WaClient client** (implements output port)
- **Logger** (implements output port)

## Key Principles

### Clean Architecture
- **Dependency Rule**: Dependencies point inward (toward domain)
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: High-level modules don't depend on low-level modules

### Hexagonal Architecture (Ports & Adapters)
- **Input Ports**: Use case interfaces (for HTTP handlers to use)
- **Output Ports**: Dependency interfaces (for external services)
- **Adapters**: Implement port interfaces

## Dependency Rules

```
✅ Allowed:
cmd/                    → internal/container, internal/config
internal/container/     → internal/core, internal/adapters, internal/config
internal/adapters/      → internal/core (domain, application, ports)
internal/core/application/ → internal/core/domain, internal/core/ports
internal/core/ports/    → internal/core/domain, internal/core/application/dto
internal/core/domain/   → NONE (pure business logic, stdlib only)

❌ Forbidden:
internal/core/domain/      → adapters, application, ports, external frameworks
internal/core/application/ → adapters (must use ports interfaces)
internal/core/ports/       → adapters (only interfaces, no implementations)
```

## Technology Stack

- **Language**: Go 1.24+
- **Web Framework**: Chi router
- **Database**: PostgreSQL with migrations
- **Logging**: Zerolog (structured logging)
- **WaClient**: whatsmeow library
- **Configuration**: Environment variables
- **Development**: Air (hot reload), Docker Compose

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