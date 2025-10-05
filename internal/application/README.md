# Application Layer - zpwoot

Esta camada implementa os **Use Cases** da aplicação seguindo os princípios da **Clean Architecture**.

## 📁 Estrutura

```
internal/application/
├── dto/                    # Data Transfer Objects
│   ├── common.go          # DTOs comuns (Response, Pagination, etc.)
│   ├── session.go         # DTOs para operações de sessão
│   └── message.go         # DTOs para operações de mensagem
├── interfaces/            # Interfaces para adapters externos
│   ├── whatsapp.go       # Interface do cliente WhatsApp
│   └── notification.go   # Interface do serviço de notificação
├── usecase/              # Use Cases organizados por domínio
│   ├── container.go      # Container de Use Cases
│   ├── session/          # Use Cases de sessão
│   │   ├── create.go     # Criar sessão
│   │   ├── get.go        # Obter detalhes da sessão
│   │   ├── list.go       # Listar sessões
│   │   ├── connect.go    # Conectar sessão
│   │   ├── disconnect.go # Desconectar sessão
│   │   ├── delete.go     # Deletar sessão
│   │   └── qr.go         # Operações de QR code
│   └── message/          # Use Cases de mensagem
│       ├── send.go       # Enviar mensagem
│       └── receive.go    # Processar mensagem recebida
└── README.md             # Este arquivo
```

## 🎯 Responsabilidades

### **DTOs (Data Transfer Objects)**
- **Propósito**: Transferir dados entre camadas sem expor entidades de domínio
- **Validação**: Contêm validação de entrada e conversão para/de domínio
- **Serialização**: Preparados para JSON serialization/deserialization

### **Interfaces**
- **WhatsAppClient**: Define operações do cliente WhatsApp (implementado na camada de infraestrutura)
- **NotificationService**: Define operações de notificação (webhooks, eventos)

### **Use Cases**
Implementam a lógica de aplicação orquestrando:
- **Domain Services** (lógica de negócio)
- **External Adapters** (WhatsApp client, notifications)
- **Repository** (persistência via domain services)

## 🔄 Fluxo de Dados

```
HTTP Handler → Use Case → Domain Service + External Adapters → Repository
     ↓              ↓              ↓                    ↓
   Request DTO → Validation → Domain Entity → Database/External API
     ↑              ↑              ↑                    ↑
HTTP Response ← Response DTO ← Domain Entity ← Database/External API
```

## 📋 Use Cases Implementados

### **Session Use Cases**

#### **CreateUseCase**
- Cria nova sessão no domínio
- Inicializa cliente WhatsApp
- Envia notificações de criação

#### **GetUseCase**
- Obtém detalhes da sessão
- Sincroniza status com WhatsApp client
- Retorna informações completas

#### **ListUseCase**
- Lista sessões com paginação
- Sincroniza status de cada sessão
- Suporte a filtros (futuro)

#### **ConnectUseCase**
- Conecta sessão ao WhatsApp
- Gerencia QR codes
- Notifica eventos de conexão

#### **DisconnectUseCase**
- Desconecta sessão graciosamente
- Atualiza status no domínio
- Notifica desconexão

#### **DeleteUseCase**
- Remove sessão completamente
- Cleanup de recursos
- Suporte a força (force delete)

#### **QRUseCase**
- Gerencia códigos QR
- Refresh de QR expirado
- Validação de status

### **Message Use Cases**

#### **SendUseCase**
- Envia mensagens de diferentes tipos
- Validação de entrada
- Métodos de conveniência por tipo

#### **ReceiveUseCase**
- Processa mensagens recebidas
- Atualiza histórico
- Envia notificações

## 🛠️ Padrões Utilizados

### **Dependency Injection**
Use cases recebem dependências via construtor:
```go
func NewCreateUseCase(
    sessionService *session.Service,
    whatsappClient interfaces.WhatsAppClient,
    notificationSvc interfaces.NotificationService,
) *CreateUseCase
```

### **Error Handling**
- Erros de domínio são convertidos para DTOs
- Rollback automático em caso de falha
- Logging estruturado

### **Async Operations**
- Notificações são enviadas de forma assíncrona
- Fire-and-forget para operações não críticas
- Context para timeout/cancelamento

### **Validation**
- Validação de entrada nos DTOs
- Validação de negócio no domínio
- Sanitização de dados

## 🔧 Container de Use Cases

O `Container` centraliza a criação e gerenciamento dos use cases:

```go
container := usecase.NewContainer(
    sessionService,
    whatsappClient,
    notificationSvc,
)

// Acesso agrupado
sessionUseCases := container.SessionUseCases()
messageUseCases := container.MessageUseCases()
```

## 🚀 Próximos Passos

1. **Testes Unitários**: Implementar testes para todos os use cases
2. **Métricas**: Adicionar instrumentação e métricas
3. **Cache**: Implementar cache para operações frequentes
4. **Rate Limiting**: Controle de taxa para operações externas
5. **Retry Logic**: Lógica de retry para operações que podem falhar

## 📚 Referências

- [Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [WhatsApp Business API](https://developers.facebook.com/docs/whatsapp)
- [whatsmeow Documentation](https://pkg.go.dev/go.mau.fi/whatsmeow)
