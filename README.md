# zpwoot - WhatsApp API Gateway

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A comprehensive WhatsApp Business API built with Go, following Clean Architecture principles. Provides endpoints for session management, messaging, contacts, groups, media handling, and integrations with Chatwoot.

## ✨ Features

- 🚀 **Clean Architecture** - Maintainable and testable codebase
- 🔄 **Automatic Migrations** - Database schema managed automatically
- 📊 **PostgreSQL Integration** - Robust data persistence
- 🐳 **Docker Ready** - Full containerization support
- 📝 **Structured Logging** - JSON logs with configurable levels
- 🔧 **Environment Configuration** - Easy deployment configuration
- 🏥 **Health Checks** - Built-in health monitoring
- 🔌 **Chatwoot Integration** - Ready for customer support integration

## 🏗️ Architecture

zpwoot follows Clean Architecture principles with clear separation of concerns:

```
┌─── HTTP Layer ────┐
│   Chi Router      │
│   Middleware      │
└───────────────────┘
         │
┌─── Application ───┐
│   Use Cases       │
│   Services        │
└───────────────────┘
         │
┌─── Domain Layer ──┐
│   Entities        │
│   Business Logic  │
└───────────────────┘
         │
┌─── Infrastructure ┐
│   Database        │
│   External APIs   │
└───────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL 15+ (or use Docker)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-org/zpwoot.git
   cd zpwoot
   ```

2. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start services**
   ```bash
   make up
   ```

4. **Run the application**
   ```bash
   go run cmd/zpwoot/main.go
   ```

5. **Access the API**
   ```bash
   curl http://localhost:8080/health
   ```

### Docker Development

```bash
# Start all services including database
make up

# Run application in development mode with hot reload
make dev

# Build and run with Docker
docker-compose up -d
```

## 📊 Database Schema

The application automatically creates and manages the following tables:

- **zpSessions** - WhatsApp session management
- **zpWebhooks** - Event notification configuration  
- **zpChatwoot** - Chatwoot integration settings
- **zpMessage** - Message synchronization tracking
- **zpMigrations** - Database version control

All migrations run automatically on startup - no manual intervention required!

## 🔧 Configuration

Configure the application using environment variables:

```bash
# Application
PORT=8080
SERVER_HOST=0.0.0.0
LOG_LEVEL=info
ZP_API_KEY=your-api-key

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/zpwoot

# WhatsApp
WA_LOG_LEVEL=INFO

# Webhooks
GLOBAL_WEBHOOK_URL=https://your-domain.com/webhooks
```

See `.env.example` for all available options.

## 🛠️ Development

### Available Commands

```bash
# Development
make dev          # Run with hot reload
make build        # Build binary
make test         # Run tests
make lint         # Run linter

# Database
make up           # Start PostgreSQL and services
make down         # Stop services
make db-reset     # Reset database

# Docker
make docker-build # Build Docker image
make docker-run   # Run in Docker

# Chatwoot Integration
make up-cw        # Start Chatwoot services
make setup-chatwoot SESSION=my-session  # Configure Chatwoot
```

### Project Structure

```
zpwoot/
├── cmd/zpwoot/           # Application entry point
├── internal/
│   ├── domain/           # Business logic
│   ├── application/      # Use cases (planned)
│   └── infra/           # Infrastructure
├── platform/
│   ├── config/          # Configuration
│   ├── database/        # Database & migrations
│   ├── logger/          # Logging
│   └── container/       # Dependency injection
└── docs/                # Documentation
```

## 📡 API Endpoints

### Current Endpoints

- `GET /` - API information
- `GET /health` - Health check with database verification

### Planned Endpoints

- `POST /api/v1/sessions` - Create WhatsApp session
- `GET /api/v1/sessions` - List sessions
- `POST /api/v1/messages` - Send message
- `GET /api/v1/messages` - Get messages
- `POST /api/v1/webhooks` - Configure webhooks
- `POST /api/v1/chatwoot` - Setup Chatwoot integration

## 🔌 Integrations

### Chatwoot
Full bidirectional integration with Chatwoot for customer support:

```bash
# Start Chatwoot
make up-cw

# Configure integration
make setup-chatwoot SESSION=my-session

# Access Chatwoot
open http://localhost:3001
```

### WhatsApp (Planned)
Integration with whatsmeow for WhatsApp Business API functionality.

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific tests
go test ./internal/domain/session/...
```

## 📚 Documentation

- [Architecture Guide](ARCHITECTURE.md) - Detailed architecture documentation
- [API Documentation](docs/api.md) - API reference (planned)
- [Deployment Guide](docs/deployment.md) - Production deployment (planned)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Follow the architecture patterns
4. Write tests for new functionality
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines

- Follow Clean Architecture principles
- Keep domain logic pure (no external dependencies)
- Use dependency injection
- Write database migrations for schema changes
- Add structured logging
- Follow Go conventions and best practices

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [whatsmeow](https://github.com/tulir/whatsmeow) - WhatsApp Web multidevice API
- [Chatwoot](https://github.com/chatwoot/chatwoot) - Customer engagement platform
- [Chi](https://github.com/go-chi/chi) - Lightweight HTTP router
- [Zerolog](https://github.com/rs/zerolog) - Structured logging

## 📞 Support

- 📧 Email: support@zpwoot.com
- 🐛 Issues: [GitHub Issues](https://github.com/your-org/zpwoot/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/your-org/zpwoot/discussions)

---

**zpwoot** - Making WhatsApp Business API integration simple and powerful! 🚀
