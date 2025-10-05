# Ports Layer - Hexagonal Architecture

Este diretório contém as **interfaces (ports)** que definem os contratos entre o **Core** (domain + application) e o **mundo externo** (adapters).

---

## 🎯 O que são Ports?

Na **Arquitetura Hexagonal (Ports & Adapters)**, os **Ports** são interfaces que:

1. **Isolam o Core** da infraestrutura externa
2. **Definem contratos** que os adapters devem implementar
3. **Permitem testabilidade** através de mocks/stubs
4. **Facilitam substituição** de implementações

---

## 📂 Estrutura

```
internal/core/ports/
├── output/                    # Portas de Saída (Output Ports)
│   ├── whatsapp.go           # Interface para cliente WhatsApp
│   └── logger.go             # Interface para logging
└── input/                     # Portas de Entrada (Input Ports) - IMPLEMENTADO
    ├── session.go            # Interfaces de use cases de sessão
    └── message.go            # Interfaces de use cases de mensagem
```

---

## 🔌 Output Ports (Portas de Saída)

### **Definição**
Interfaces que o **Core define** e os **Adapters implementam**.

### **Quando criar um Output Port?**
Crie um output port quando o Core precisa:
- ✅ Acessar um serviço externo (API, WhatsApp, etc.)
- ✅ Persistir dados (Database)
- ✅ Enviar notificações (Webhooks, Email, SMS)
- ✅ Fazer logging
- ✅ Acessar cache (Redis, Memcached)
- ✅ Enviar mensagens (Queue, Pub/Sub)

### **Exemplo: WhatsApp Client Port**

```go
// internal/core/ports/output/whatsapp.go
package output

import (
    "context"
    "time"
)

// WhatsAppClient define o contrato para comunicação com WhatsApp
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

// SessionStatus representa o status de uma sessão WhatsApp
type SessionStatus struct {
    SessionID   string    `json:"sessionId"`
    Connected   bool      `json:"connected"`
    LoggedIn    bool      `json:"loggedIn"`
    DeviceJID   string    `json:"deviceJid,omitempty"`
    PushName    string    `json:"pushName,omitempty"`
    ConnectedAt time.Time `json:"connectedAt,omitempty"`
    LastSeen    time.Time `json:"lastSeen,omitempty"`
}

// QRCodeInfo contém informações do QR Code
type QRCodeInfo struct {
    Code      string    `json:"code"`
    Base64    string    `json:"base64"`
    ExpiresAt time.Time `json:"expiresAt"`
}

// MessageResult representa o resultado do envio de mensagem
type MessageResult struct {
    MessageID string    `json:"messageId"`
    Status    string    `json:"status"`
    SentAt    time.Time `json:"sentAt"`
}

// MediaData representa dados de mídia
type MediaData struct {
    MimeType string `json:"mimeType"`
    Data     []byte `json:"data"`
    FileName string `json:"fileName,omitempty"`
    Caption  string `json:"caption,omitempty"`
}

// Location representa uma localização geográfica
type Location struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Name      string  `json:"name,omitempty"`
    Address   string  `json:"address,omitempty"`
}

// ContactInfo representa informações de contato
type ContactInfo struct {
    Name        string `json:"name"`
    PhoneNumber string `json:"phoneNumber"`
}
```

**Implementação (Adapter):**
```go
// internal/adapters/waclient/whatsapp_adapter.go
package waclient

import "zpwoot/internal/core/ports/output"

type WhatsAppAdapter struct {
    client *WAClient
}

// Implementa output.WhatsAppClient
func (a *WhatsAppAdapter) CreateSession(ctx context.Context, sessionID string) error {
    // Implementação usando whatsmeow
}
```

---

### **Exemplo: Notification Service Port**

```go
// internal/core/ports/output/notification.go
package output

import (
    "context"
    "time"
)

// NotificationService define o contrato para envio de notificações
type NotificationService interface {
    // Webhook genérico
    SendWebhook(ctx context.Context, event *WebhookEvent) error
    
    // Notificações específicas
    NotifySessionConnected(ctx context.Context, sessionID string, deviceJID string) error
    NotifySessionDisconnected(ctx context.Context, sessionID string) error
    NotifyQRCodeGenerated(ctx context.Context, sessionID string, qrCode string, expiresAt time.Time) error
    NotifyMessageReceived(ctx context.Context, sessionID string, message *MessageEvent) error
    NotifyMessageSent(ctx context.Context, sessionID string, messageID string) error
}

// WebhookEvent representa um evento genérico de webhook
type WebhookEvent struct {
    Type      string      `json:"type"`
    SessionID string      `json:"sessionId"`
    Event     interface{} `json:"event"`
    Timestamp time.Time   `json:"timestamp"`
}

// MessageEvent representa um evento de mensagem
type MessageEvent struct {
    ID        string    `json:"id"`
    Chat      string    `json:"chat"`
    Sender    string    `json:"sender"`
    PushName  string    `json:"pushName"`
    Timestamp time.Time `json:"timestamp"`
    FromMe    bool      `json:"fromMe"`
    Type      string    `json:"type"`
    IsGroup   bool      `json:"isGroup"`
    Content   string    `json:"content,omitempty"`
}

// Constantes de tipos de eventos
const (
    EventTypeSessionConnected    = "session.connected"
    EventTypeSessionDisconnected = "session.disconnected"
    EventTypeQRCodeGenerated     = "qr.generated"
    EventTypeMessageReceived     = "message.received"
    EventTypeMessageSent         = "message.sent"
)
```

---

### **Exemplo: Logger Port (RECOMENDADO)**

```go
// internal/core/ports/output/logger.go
package output

import "context"

// Logger define o contrato para logging
type Logger interface {
    // Níveis de log
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
    
    // Context-aware logging
    WithContext(ctx context.Context) Logger
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    WithError(err error) Logger
    
    // Domain-specific helpers
    WithComponent(component string) Logger
    WithRequestID(requestID string) Logger
    WithSessionID(sessionID string) Logger
}

// Field representa um campo de log estruturado
type Field struct {
    Key   string
    Value interface{}
}
```

**Implementação (Adapter):**
```go
// internal/adapters/logger/logger_adapter.go
package logger

import "zpwoot/internal/core/ports/output"

type LoggerAdapter struct {
    logger *Logger // implementação zerolog
}

// Implementa output.Logger
func (l *LoggerAdapter) Info(msg string, fields ...output.Field) {
    event := l.logger.Info()
    for _, f := range fields {
        event = event.Interface(f.Key, f.Value)
    }
    event.Msg(msg)
}
```

---

## 🎬 Input Ports (Portas de Entrada) - OPCIONAL

### **Definição**
Interfaces que definem os **Use Cases** (casos de uso).

### **Quando criar Input Ports?**
- ✅ Para melhorar **testabilidade** (mocks de use cases)
- ✅ Para **desacoplar** handlers HTTP dos use cases concretos
- ✅ Para permitir **múltiplas implementações** de use cases

### **Exemplo: Session Use Case Ports**

```go
// internal/core/ports/input/session.go
package input

import (
    "context"
    "zpwoot/internal/core/application/dto"
)

// SessionCreator define o contrato para criar sessões
type SessionCreator interface {
    Execute(ctx context.Context, req *dto.CreateSessionRequest) (*dto.SessionResponse, error)
}

// SessionConnector define o contrato para conectar sessões
type SessionConnector interface {
    Execute(ctx context.Context, sessionID string) (*dto.SessionResponse, error)
}

// SessionDeleter define o contrato para deletar sessões
type SessionDeleter interface {
    Execute(ctx context.Context, sessionID string, force bool) error
}

// SessionGetter define o contrato para obter sessão
type SessionGetter interface {
    Execute(ctx context.Context, sessionID string) (*dto.SessionResponse, error)
}

// SessionLister define o contrato para listar sessões
type SessionLister interface {
    Execute(ctx context.Context, limit, offset int) (*dto.SessionListResponse, error)
}

// QRCodeManager define o contrato para gerenciar QR codes
type QRCodeManager interface {
    GetQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error)
    RefreshQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error)
}
```

**Uso no Handler:**
```go
// internal/adapters/http/handlers/session.go
package handlers

import "zpwoot/internal/core/ports/input"

type SessionHandler struct {
    creator   input.SessionCreator
    connector input.SessionConnector
    deleter   input.SessionDeleter
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
    // Usa a interface, não a implementação concreta
    response, err := h.creator.Execute(r.Context(), req)
}
```

---

## 📋 Regras para Criar Ports

### ✅ **DEVE:**
1. Definir interfaces claras e focadas (Interface Segregation Principle)
2. Usar tipos de domínio ou DTOs nos parâmetros
3. Retornar erros tipados quando possível
4. Documentar cada método com comentários
5. Agrupar métodos relacionados na mesma interface

### ❌ **NÃO DEVE:**
1. Conter implementações (apenas interfaces)
2. Depender de frameworks externos
3. Ter lógica de negócio
4. Ter dependências de infraestrutura

---

## 🔄 Fluxo de Dependências

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Handler                         │
│                   (Adapter - Input)                     │
└─────────────────────────────────────────────────────────┘
                         ↓
                    usa interface
                         ↓
┌─────────────────────────────────────────────────────────┐
│                  Input Port Interface                   │
│              (ports/input/session.go)                   │
└─────────────────────────────────────────────────────────┘
                         ↑
                    implementa
                         ↑
┌─────────────────────────────────────────────────────────┐
│                    Use Case                             │
│           (application/usecase/session/)                │
└─────────────────────────────────────────────────────────┘
                         ↓
                    usa interface
                         ↓
┌─────────────────────────────────────────────────────────┐
│                 Output Port Interface                   │
│            (ports/output/whatsapp.go)                   │
└─────────────────────────────────────────────────────────┘
                         ↑
                    implementa
                         ↑
┌─────────────────────────────────────────────────────────┐
│                  WhatsApp Adapter                       │
│                (Adapter - Output)                       │
└─────────────────────────────────────────────────────────┘
```

---

## 🧪 Testabilidade

### **Vantagem dos Ports:**
Com ports, você pode facilmente criar **mocks** para testes:

```go
// internal/core/ports/output/whatsapp_mock.go
package output

type MockWhatsAppClient struct {
    CreateSessionFunc func(ctx context.Context, sessionID string) error
}

func (m *MockWhatsAppClient) CreateSession(ctx context.Context, sessionID string) error {
    if m.CreateSessionFunc != nil {
        return m.CreateSessionFunc(ctx, sessionID)
    }
    return nil
}
```

**Uso em testes:**
```go
func TestCreateUseCase(t *testing.T) {
    mockWA := &output.MockWhatsAppClient{
        CreateSessionFunc: func(ctx context.Context, sessionID string) error {
            return nil // Simula sucesso
        },
    }
    
    useCase := session.NewCreateUseCase(sessionService, mockWA, mockNotification)
    result, err := useCase.Execute(ctx, req)
    
    assert.NoError(t, err)
}
```

---

## 📚 Referências

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Ports and Adapters Pattern](https://herbertograca.com/2017/09/14/ports-adapters-architecture/)
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)

