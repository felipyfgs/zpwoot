# 📋 Checklist de Migração para Clean Architecture

Este documento guia a reorganização completa do projeto para seguir os princípios de **Clean Architecture** e **Hexagonal Architecture (Ports & Adapters)**.

---

## 🎯 Objetivo

Garantir que o projeto tenha:
- ✅ Separação clara de responsabilidades
- ✅ Core independente de frameworks e infraestrutura
- ✅ Interfaces bem definidas (Ports)
- ✅ Adapters implementando as interfaces
- ✅ Fluxo de dependências correto (sempre apontando para o Core)

---

## 📊 Fase 0: Preparação e Análise

### 0.1 Mapeamento Inicial

- [ ] **Listar todos os arquivos `.go` do projeto**
  ```bash
  find internal -name "*.go" -type f | sort > files_inventory.txt
  ```

- [ ] **Identificar dependências externas usadas**
  - [ ] HTTP (chi, handlers, middleware)
  - [ ] Database (sqlx, postgres)
  - [ ] WhatsApp (whatsmeow)
  - [ ] Logger (zerolog)
  - [ ] Config (godotenv)
  - [ ] Outras: _______________

- [ ] **Criar backup do projeto**
  ```bash
  git checkout -b migration-clean-architecture
  git add .
  git commit -m "Backup antes da migração para Clean Architecture"
  ```

- [ ] **Ler documentação criada**
  - [ ] `internal/core/README.md`
  - [ ] `internal/core/ports/README.md`
  - [ ] `internal/core/CURRENT_STATE.md`

---

## 📂 Fase 1: Análise de Arquivos por Camada

### 1.1 Analisar `internal/core/domain/`

Para cada arquivo em `domain/`:

**`domain/session/entity.go`**
- [ ] Verificar se contém apenas:
  - [ ] Struct `Session` com campos de negócio
  - [ ] Métodos de comportamento da entidade
  - [ ] Value Objects (ex: `Status`)
  - [ ] Sem dependências externas (apenas stdlib)
- [ ] **Decisão:** ✅ Manter / ❌ Mover para: _______________

**`domain/session/repository.go`**
- [ ] Verificar se é apenas interface (não implementação)
- [ ] Verificar se métodos usam apenas tipos do domínio
- [ ] **Decisão:** ✅ Manter / ❌ Mover para: _______________

**`domain/session/service.go`**
- [ ] Verificar se contém apenas lógica de negócio pura
- [ ] Verificar se depende apenas de `Repository` interface
- [ ] Verificar se não tem dependências de adapters
- [ ] **Decisão:** ✅ Manter / ❌ Mover para: _______________

**`domain/shared/errors.go`**
- [ ] Verificar se contém apenas erros de domínio
- [ ] **Decisão:** ✅ Manter / ❌ Mover para: _______________

### 1.2 Analisar `internal/core/application/`

**`application/dto/`**
- [ ] Listar todos os DTOs existentes:
  - [ ] `common.go` - Response, ErrorResponse, Pagination
  - [ ] `session.go` - CreateSessionRequest, SessionResponse, etc.
  - [ ] `message.go` - SendMessageRequest, MessageResponse, etc.
  - [ ] Outros: _______________

- [ ] Para cada DTO, verificar:
  - [ ] É usado para comunicação externa (API)?
  - [ ] Tem conversão de/para entidades de domínio?
  - [ ] **Decisão:** ✅ Manter em `application/dto/` / ❌ Mover para: _______________

**`application/interfaces/`** ⚠️ **ATENÇÃO: Deve ser movido para `ports/output/`**

- [ ] Listar todas as interfaces:
  - [ ] `whatsapp.go` - WhatsAppClient
  - [ ] `notification.go` - NotificationService
  - [ ] Outras: _______________

- [ ] Para cada interface, decidir:
  - [ ] É uma dependência externa que o Core precisa?
  - [ ] Deve estar em `ports/output/`?
  - [ ] **Ação:** 🔄 Mover para `ports/output/`

**`application/usecase/`**

- [ ] Listar todos os use cases:
  - [ ] `session/create.go`
  - [ ] `session/connect.go`
  - [ ] `session/disconnect.go`
  - [ ] `session/logout.go`
  - [ ] `session/get.go`
  - [ ] `session/list.go`
  - [ ] `session/delete.go`
  - [ ] `session/qr.go`
  - [ ] `message/send.go`
  - [ ] `message/receive.go`
  - [ ] Outros: _______________

- [ ] Para cada use case, verificar:
  - [ ] Tem dependências de Domain Services?
  - [ ] Tem dependências de Ports (interfaces)?
  - [ ] Não tem dependências diretas de Adapters?
  - [ ] Retorna DTOs (não entidades de domínio)?
  - [ ] **Decisão:** ✅ Manter / ❌ Refatorar

**`application/validators/`**

- [ ] Listar validadores:
  - [ ] `session.go`
  - [ ] `message.go`
  - [ ] Outros: _______________

- [ ] Verificar se validam apenas entrada de API (não regras de negócio)
- [ ] **Decisão:** ✅ Manter / ❌ Mover para: _______________

### 1.3 Analisar `internal/core/ports/`

**Estado atual:**
- [ ] Diretório existe mas está vazio

**Ações necessárias:**
- [ ] Criar `internal/core/ports/output/`
- [ ] Criar `internal/core/ports/input/` (opcional)

---

## 🔧 Fase 2: Reorganização de Ports

### 2.1 Criar estrutura de Ports

```bash
mkdir -p internal/core/ports/output
mkdir -p internal/core/ports/input  # opcional
```

- [ ] Estrutura criada

### 2.2 Mover interfaces para `ports/output/`

**WhatsApp Client:**
- [ ] Mover `application/interfaces/whatsapp.go` → `ports/output/whatsapp.go`
  ```bash
  mv internal/core/application/interfaces/whatsapp.go internal/core/ports/output/whatsapp.go
  ```
- [ ] Atualizar package de `package interfaces` para `package output`
- [ ] Verificar se tipos relacionados foram movidos juntos:
  - [ ] `SessionStatus`
  - [ ] `QRCodeInfo`
  - [ ] `MessageResult`
  - [ ] `MediaData`
  - [ ] `Location`
  - [ ] `ContactInfo`
  - [ ] `WhatsAppError`

**Notification Service:**
- [ ] Mover `application/interfaces/notification.go` → `ports/output/notification.go`
  ```bash
  mv internal/core/application/interfaces/notification.go internal/core/ports/output/notification.go
  ```
- [ ] Atualizar package de `package interfaces` para `package output`
- [ ] Verificar se tipos relacionados foram movidos juntos:
  - [ ] `WebhookEvent`
  - [ ] `MessageEvent`
  - [ ] `SessionEvent`
  - [ ] `QRCodeEvent`
  - [ ] Constantes de eventos

### 2.3 Criar Logger Port (NOVO)

- [ ] Criar arquivo `internal/core/ports/output/logger.go`
- [ ] Definir interface `Logger` com métodos:
  - [ ] `Debug(msg string, fields ...Field)`
  - [ ] `Info(msg string, fields ...Field)`
  - [ ] `Warn(msg string, fields ...Field)`
  - [ ] `Error(msg string, fields ...Field)`
  - [ ] `Fatal(msg string, fields ...Field)`
  - [ ] `WithContext(ctx context.Context) Logger`
  - [ ] `WithField(key string, value interface{}) Logger`
  - [ ] `WithFields(fields map[string]interface{}) Logger`
  - [ ] `WithError(err error) Logger`
  - [ ] `WithComponent(component string) Logger`
  - [ ] `WithRequestID(requestID string) Logger`
  - [ ] `WithSessionID(sessionID string) Logger`
- [ ] Definir struct `Field` para campos estruturados

### 2.4 Criar Input Ports (OPCIONAL)

Se decidir criar interfaces para use cases:

- [ ] Criar `internal/core/ports/input/session.go`
  - [ ] Interface `SessionCreator`
  - [ ] Interface `SessionConnector`
  - [ ] Interface `SessionDisconnector`
  - [ ] Interface `SessionDeleter`
  - [ ] Interface `SessionGetter`
  - [ ] Interface `SessionLister`
  - [ ] Interface `QRCodeManager`

- [ ] Criar `internal/core/ports/input/message.go`
  - [ ] Interface `MessageSender`
  - [ ] Interface `MessageReceiver`

### 2.5 Remover diretório vazio

- [ ] Verificar se `application/interfaces/` está vazio
- [ ] Remover diretório:
  ```bash
  rmdir internal/core/application/interfaces
  ```

---

## 🔄 Fase 3: Atualizar Imports

### 3.1 Atualizar imports em Use Cases

Para cada use case em `application/usecase/`:

**`session/create.go`**
- [ ] Trocar `"zpwoot/internal/core/application/interfaces"` por `"zpwoot/internal/core/ports/output"`
- [ ] Trocar `interfaces.WhatsAppClient` por `output.WhatsAppClient`
- [ ] Trocar `interfaces.NotificationService` por `output.NotificationService`
- [ ] Verificar compilação: `go build ./internal/core/application/usecase/session/`

**`session/connect.go`**
- [ ] Atualizar imports
- [ ] Verificar compilação

**`session/disconnect.go`**
- [ ] Atualizar imports
- [ ] Verificar compilação

**`session/logout.go`**
- [ ] Atualizar imports
- [ ] Verificar compilação

**`session/get.go`**
- [ ] Atualizar imports
- [ ] Verificar compilação

**`session/list.go`**
- [ ] Atualizar imports
- [ ] Verificar compilação

**`session/delete.go`**
- [ ] Atualizar imports
- [ ] Verificar compilação

**`session/qr.go`**
- [ ] Atualizar imports
- [ ] Verificar compilação

**`message/send.go`**
- [ ] Atualizar imports
- [ ] Verificar compilação

**`message/receive.go`**
- [ ] Atualizar imports
- [ ] Verificar compilação

### 3.2 Atualizar imports em Container

**`internal/container/container.go`**
- [ ] Trocar `"zpwoot/internal/core/application/interfaces"` por `"zpwoot/internal/core/ports/output"`
- [ ] Trocar `interfaces.WhatsAppClient` por `output.WhatsAppClient`
- [ ] Trocar `interfaces.NotificationService` por `output.NotificationService`
- [ ] Verificar compilação: `go build ./internal/container/`

### 3.3 Atualizar imports em Adapters

**`adapters/waclient/whatsapp_adapter.go`**
- [ ] Trocar `"zpwoot/internal/core/application/interfaces"` por `"zpwoot/internal/core/ports/output"`
- [ ] Trocar `interfaces.SessionStatus` por `output.SessionStatus`
- [ ] Trocar `interfaces.QRCodeInfo` por `output.QRCodeInfo`
- [ ] Trocar `interfaces.MessageResult` por `output.MessageResult`
- [ ] Verificar compilação: `go build ./internal/adapters/waclient/`

**`adapters/http/router/router.go`**
- [ ] Verificar se usa interfaces
- [ ] Atualizar imports se necessário
- [ ] Verificar compilação: `go build ./internal/adapters/http/`

### 3.4 Atualizar imports automaticamente (alternativa)

- [ ] Executar substituição em massa:
  ```bash
  find internal -name "*.go" -type f -exec sed -i 's|zpwoot/internal/core/application/interfaces|zpwoot/internal/core/ports/output|g' {} \;
  ```
- [ ] Executar `gofmt` para formatar:
  ```bash
  gofmt -w internal/
  ```

---

## 🏗️ Fase 4: Implementar Logger Adapter

### 4.1 Criar Logger Adapter

- [ ] Criar arquivo `internal/adapters/logger/logger_adapter.go`
- [ ] Implementar struct `LoggerAdapter` que implementa `output.Logger`
- [ ] Wrapper para `*logger.Logger` existente
- [ ] Implementar todos os métodos da interface

### 4.2 Atualizar Container para usar Logger Port

**`internal/container/container.go`**
- [ ] Trocar `logger *logger.Logger` por `logger output.Logger`
- [ ] Atualizar inicialização para usar `LoggerAdapter`
- [ ] Verificar compilação

---

## ✅ Fase 5: Verificação e Testes

### 5.1 Verificar Compilação

- [ ] Compilar todo o projeto:
  ```bash
  go build ./...
  ```
- [ ] Resolver erros de compilação (se houver)

### 5.2 Executar go mod tidy

- [ ] Limpar dependências:
  ```bash
  go mod tidy
  ```

### 5.3 Executar Testes

- [ ] Executar todos os testes:
  ```bash
  go test ./...
  ```
- [ ] Corrigir testes quebrados (se houver)

### 5.4 Verificar Estrutura Final

- [ ] Executar tree para verificar estrutura:
  ```bash
  tree internal/core -L 3
  ```

- [ ] Estrutura esperada:
  ```
  internal/core/
  ├── domain/
  │   ├── session/
  │   │   ├── entity.go
  │   │   ├── repository.go
  │   │   └── service.go
  │   └── shared/
  │       └── errors.go
  ├── application/
  │   ├── dto/
  │   ├── usecase/
  │   └── validators/
  └── ports/
      ├── output/
      │   ├── whatsapp.go
      │   ├── notification.go
      │   └── logger.go
      └── input/  (opcional)
  ```

---

## 📝 Fase 6: Documentação

### 6.1 Atualizar ARCHITECTURE.md

- [ ] Atualizar diagrama de camadas
- [ ] Atualizar seção de Ports
- [ ] Adicionar exemplos de uso

### 6.2 Atualizar README.md do projeto

- [ ] Atualizar seção de arquitetura
- [ ] Adicionar link para `internal/core/README.md`

### 6.3 Criar README por camada (se não existir)

- [ ] `internal/adapters/README.md`
- [ ] `internal/config/README.md`
- [ ] `internal/container/README.md`

---

## 🎯 Fase 7: Validação Final

### 7.1 Checklist de Qualidade

- [ ] **Domain não depende de nada** (exceto stdlib)
- [ ] **Application depende apenas de Domain e Ports**
- [ ] **Ports define apenas interfaces** (sem implementações)
- [ ] **Adapters implementam Ports**
- [ ] **Fluxo de dependências está correto** (sempre para dentro)
- [ ] **Não há imports com alias desnecessários**
- [ ] **Cada arquivo tem responsabilidade única**
- [ ] **Nomes de arquivos são consistentes**

### 7.2 Code Review

- [ ] Revisar cada arquivo movido
- [ ] Verificar se lógica de negócio não vazou para adapters
- [ ] Verificar se adapters não têm lógica de negócio
- [ ] Verificar se DTOs são usados corretamente

### 7.3 Executar Aplicação

- [ ] Iniciar aplicação:
  ```bash
  make run
  ```
- [ ] Testar endpoints principais:
  - [ ] POST /sessions (criar sessão)
  - [ ] GET /sessions (listar sessões)
  - [ ] POST /sessions/:id/connect (conectar)
  - [ ] GET /sessions/:id/qr (obter QR)
  - [ ] POST /messages/send (enviar mensagem)

---

## 🚀 Fase 8: Commit e Deploy

### 8.1 Commit das Mudanças

- [ ] Adicionar arquivos:
  ```bash
  git add .
  ```
- [ ] Commit:
  ```bash
  git commit -m "refactor: migrate to Clean Architecture with Ports & Adapters

  - Move interfaces to internal/core/ports/output/
  - Create Logger port interface
  - Update all imports to use ports
  - Remove application/interfaces/ directory
  - Update documentation"
  ```

### 8.2 Merge para main

- [ ] Criar Pull Request
- [ ] Code Review
- [ ] Merge para main

---

## 📊 Métricas de Sucesso

Após a migração, o projeto deve ter:

- ✅ **0 dependências** do Domain para Application/Ports/Adapters
- ✅ **0 dependências** do Application para Adapters
- ✅ **100% das interfaces** em Ports
- ✅ **100% dos adapters** implementando interfaces de Ports
- ✅ **0 imports com alias** desnecessários
- ✅ **Compilação sem erros**
- ✅ **Testes passando**
- ✅ **Aplicação funcionando**

---

## 🆘 Troubleshooting

### Erro: "undefined: interfaces.WhatsAppClient"

**Solução:** Atualizar import de `interfaces` para `output`

### Erro: "cannot use ... as type output.Logger"

**Solução:** Criar LoggerAdapter que implementa a interface

### Erro: "import cycle not allowed"

**Solução:** Verificar se Domain não está importando Application

---

## ✅ Conclusão

Ao completar esta checklist, o projeto estará:
- ✅ Seguindo Clean Architecture
- ✅ Com Ports & Adapters bem definidos
- ✅ Testável e manutenível
- ✅ Independente de frameworks
- ✅ Pronto para crescer

**Data de conclusão:** _______________  
**Responsável:** _______________

