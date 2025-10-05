# Core Layer - Clean Architecture

O diretório `internal/core` contém a **lógica central do sistema**, independente de frameworks, banco de dados ou APIs externas. Seguindo a **Clean Architecture**, ele é dividido em três camadas principais: **Domain**, **Application** e **Ports**.

---

## 🎯 Visão Geral

```
┌─────────────────────────────────────────────────────────────┐
│                    ADAPTERS (Infraestrutura)                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   HTTP       │  │   Database   │  │  WhatsApp    │      │
│  │   Handlers   │  │  Repository  │  │   Client     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ↓ ↑ implementa/usa
┌─────────────────────────────────────────────────────────────┐
│                    PORTS (Interfaces)                        │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Input Ports (Use Cases)  │  Output Ports (Services) │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            ↓ ↑ define/usa
┌─────────────────────────────────────────────────────────────┐
│                    APPLICATION (Use Cases)                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Session    │  │   Message    │  │     DTOs     │      │
│  │  Use Cases   │  │  Use Cases   │  │  Validators  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ↓ ↑ orquestra
┌─────────────────────────────────────────────────────────────┐
│                    DOMAIN (Business Logic)                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Entities   │  │   Services   │  │ Repositories │      │
│  │   (Session)  │  │  (Business)  │  │ (Interfaces) │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

---

## 📂 Estrutura

```
internal/core/
├── domain/              # Camada de Domínio (Regras de Negócio)
│   ├── session/        # Agregado de Sessão
│   │   ├── entity.go
│   │   ├── repository.go
│   │   └── service.go
│   └── shared/         # Código compartilhado do domínio
│       └── errors.go
│
├── application/        # Camada de Aplicação (Use Cases)
│   ├── dto/           # Data Transfer Objects
│   ├── usecase/       # Casos de Uso
│   │   ├── session/
│   │   └── message/
│   └── validators/    # Validadores de entrada
│
└── ports/             # Camada de Portas (Interfaces)
    ├── input/         # Portas de Entrada (opcional)
    └── output/        # Portas de Saída
```

---

## 1️⃣ Domain Layer (`internal/core/domain/`)

### **Responsabilidade**
Contém a **lógica de negócio pura** e as **regras de domínio**. É o coração da aplicação.

### **O que DEVE conter:**
- ✅ **Entities** (Entidades de domínio)
- ✅ **Value Objects** (Objetos de valor)
- ✅ **Domain Services** (Serviços de domínio)
- ✅ **Repository Interfaces** (Contratos de persistência)
- ✅ **Domain Events** (Eventos de domínio)
- ✅ **Domain Errors** (Erros de negócio)

### **O que NÃO DEVE conter:**
- ❌ Implementações de banco de dados
- ❌ Chamadas HTTP ou APIs externas
- ❌ Frameworks externos (exceto stdlib)
- ❌ DTOs de API
- ❌ Lógica de apresentação

### **Regras de Dependência:**
- ✅ Pode depender APENAS de: `stdlib` do Go
- ❌ NÃO pode depender de: `application`, `ports`, `adapters`

### **Exemplo:**
```go
// internal/core/domain/session/entity.go
package session

type Session struct {
    ID              string
    Name            string
    DeviceJID       string
    IsConnected     bool
    QRCode          string
    QRCodeExpiresAt *time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

func NewSession(name string) *Session {
    return &Session{
        ID:        uuid.New().String(),
        Name:      name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}
```

---

## 2️⃣ Application Layer (`internal/core/application/`)

### **Responsabilidade**
Implementa os **casos de uso do sistema**, ou seja, **como o sistema deve se comportar** em resposta a ações externas.

### **O que DEVE conter:**
- ✅ **Use Cases** - Implementação dos casos de uso
- ✅ **DTOs** - Data Transfer Objects para API
- ✅ **Validators** - Validação de entrada de dados

### **O que NÃO DEVE conter:**
- ❌ Lógica de negócio pura (vai no domain)
- ❌ Implementações de infraestrutura
- ❌ Detalhes de HTTP, Database, etc.

### **Regras de Dependência:**
- ✅ Pode depender de: `domain`, `ports`
- ❌ NÃO pode depender de: `adapters`

### **Exemplo:**
```go
// internal/core/application/usecase/session/create.go
package session

type CreateUseCase struct {
    sessionService  *session.Service
    whatsappClient  output.WhatsAppClient
    notificationSvc output.NotificationService
}

func (uc *CreateUseCase) Execute(ctx context.Context, req *dto.CreateSessionRequest) (*dto.SessionResponse, error) {
    // 1. Validar entrada
    // 2. Criar sessão no domínio
    // 3. Criar sessão no WhatsApp (via port)
    // 4. Notificar criação (via port)
    // 5. Retornar DTO
}
```

---

## 3️⃣ Ports Layer (`internal/core/ports/`)

### **Responsabilidade**
Define as **interfaces (contratos)** entre o Core e o mundo externo (adapters).

### **Tipos de Ports:**

#### **A. Output Ports** (`output/`)
Interfaces que o **Core define** e os **Adapters implementam**.

**Exemplos:**
- `WhatsAppClient` - Comunicação com WhatsApp
- `NotificationService` - Envio de notificações
- `Logger` - Logging estruturado

#### **B. Input Ports** (`input/`) - OPCIONAL
Interfaces que definem os **Use Cases**.

**Exemplos:**
- `SessionCreator` - Criar sessão
- `MessageSender` - Enviar mensagem

### **Exemplo:**
```go
// internal/core/ports/output/whatsapp.go
package output

type WhatsAppClient interface {
    CreateSession(ctx context.Context, sessionID string) error
    SendTextMessage(ctx context.Context, sessionID, to, text string) (*MessageResult, error)
}
```

**Veja mais detalhes em:** [`ports/README.md`](./ports/README.md)

---

## 📋 Regras de Dependência

```
ADAPTERS  →  PORTS  →  APPLICATION  →  DOMAIN
   ↓          ↓           ↓              ↓
Implementa  Define    Orquestra     Regras de
Interfaces  Contratos  Use Cases     Negócio
```

### **Regra de Ouro:**
> **As dependências sempre apontam PARA DENTRO (em direção ao domínio)**

---

## 🔄 Fluxo de Execução

```
POST /sessions
   ↓
SessionHandler.CreateSession()  (Adapter)
   ↓
CreateUseCase.Execute()         (Application)
   ↓
session.Service.CreateSession() (Domain)
   ↓
session.Repository.Create()     (Domain Interface)
   ↑
SessionRepository.Create()      (Adapter Implementation)
   ↓
PostgreSQL
```

---

## ✅ Boas Práticas

### **Domain:**
1. Mantenha o domínio **puro** (sem dependências externas)
2. Use **Value Objects** para conceitos importantes
3. Defina **erros de domínio** específicos

### **Application:**
1. Use cases devem ser **pequenos e focados**
2. Sempre **valide entrada** antes de processar
3. Use **DTOs** para comunicação externa
4. Implemente **rollback** em caso de falha

### **Ports:**
1. Interfaces devem ser **pequenas e coesas**
2. Use **tipos de domínio** ou DTOs nos parâmetros
3. Documente **contratos** claramente

---

## 📚 Referências

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture - Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design - Eric Evans](https://www.domainlanguage.com/ddd/)

