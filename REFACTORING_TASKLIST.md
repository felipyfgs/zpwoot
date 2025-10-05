# 📋 Tasklist - Refatoração Crítica de Clean Architecture

**Projeto:** zpwoot  
**Objetivo:** Corrigir violações críticas de Clean Architecture e eliminar duplicações  
**Total de Tarefas:** 47 tarefas organizadas em 5 sprints

---

## 📊 Visão Geral

```
┌─────────────────────────────────────────────────────────────┐
│  REFATORAÇÃO CRÍTICA - CLEAN ARCHITECTURE                   │
│  ════════════════════════════════════════════════════════   │
│  Sprint 1: Eliminar Duplicação (7 tarefas)                 │
│  Sprint 2: Implementar Use Cases (10 tarefas)              │
│  Sprint 3: Corrigir Context/Concorrência (10 tarefas)      │
│  Sprint 4: Implementar Transações (10 tarefas)             │
│  Sprint 5: Validação e Documentação (10 tarefas)           │
│  ════════════════════════════════════════════════════════   │
│  Total: 47 tarefas | Estimativa: 38-48 horas               │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎯 Sprint 1: Eliminar Duplicação de Entidades

**Objetivo:** Remover duplicação de Session/SessionInfo e Status entre domain e waclient  
**Prioridade:** 🔴 CRÍTICA  
**Estimativa:** 10-12 horas

### Tarefas:

- [ ] **1.1** Remover SessionInfo de waclient/types.go
  - Deletar `type SessionInfo` e todas as referências
  - Usar `domain.Session` em todo o código

- [ ] **1.2** Remover SessionStatus de waclient/types.go
  - Deletar `type SessionStatus`
  - Usar `domain.Status` em todo waclient

- [ ] **1.3** Remover DBSessionManager de waclient
  - Deletar `session_manager.go`
  - Usar apenas `database/repository/session.go`

- [ ] **1.4** Atualizar waclient para usar domain.Session
  - Modificar `WAClient` para trabalhar com `domain.Session`
  - Atualizar todos os métodos

- [ ] **1.5** Atualizar waclient para usar domain.Status
  - Substituir todas as referências de `SessionStatus`
  - Importar e usar `domain.Status`

- [ ] **1.6** Conectar domain.Repository ao domain.Service
  - Injetar `database/repository` no `domain/session/service.go`
  - Implementar Dependency Injection

- [ ] **1.7** Testar e validar eliminação de duplicação
  - Compilar código
  - Executar testes
  - Garantir que não há regressões

**Arquivos Afetados:**
- `internal/adapters/waclient/types.go`
- `internal/adapters/waclient/session_manager.go` (deletar)
- `internal/adapters/waclient/client.go`
- `internal/domain/session/service.go`

---

## 🎯 Sprint 2: Implementar Use Cases

**Objetivo:** Criar e implementar Use Cases para desacoplar handlers de adapters  
**Prioridade:** 🔴 CRÍTICA  
**Estimativa:** 12-16 horas

### Tarefas:

- [ ] **2.1** Criar SessionUseCase interface
  - Definir interface em `application/usecase/session/`
  - Documentar métodos

- [ ] **2.2** Implementar CreateSessionUseCase
  - Criar `create.go` com lógica de criação
  - Orquestrar domain service + waclient

- [ ] **2.3** Implementar ConnectSessionUseCase
  - Criar `connect.go` com lógica de conexão
  - Gerenciar QR code e estados

- [ ] **2.4** Implementar DisconnectSessionUseCase
  - Criar `disconnect.go` com lógica de desconexão
  - Manter credenciais para reconexão

- [ ] **2.5** Implementar LogoutSessionUseCase
  - Criar `logout.go` com lógica de logout
  - Limpar credenciais e device

- [ ] **2.6** Implementar DeleteSessionUseCase
  - Criar `delete.go` com lógica de deleção
  - Remover de domain + waclient

- [ ] **2.7** Implementar GetQRCodeUseCase
  - Criar `qr.go` com lógica de QR code
  - Gerenciar geração e refresh

- [ ] **2.8** Atualizar handlers para usar Use Cases
  - Modificar `handlers/session.go`
  - Injetar Use Cases via DI

- [ ] **2.9** Remover dependência direta de waclient nos handlers
  - Deletar `import "zpwoot/internal/adapters/waclient"`
  - Usar apenas interfaces de Use Cases

- [ ] **2.10** Testar Use Cases isoladamente
  - Criar testes unitários para cada Use Case
  - Mockar dependências

**Arquivos Criados:**
- `internal/application/usecase/session/interface.go`
- `internal/application/usecase/session/create.go`
- `internal/application/usecase/session/connect.go`
- `internal/application/usecase/session/disconnect.go`
- `internal/application/usecase/session/logout.go`
- `internal/application/usecase/session/delete.go`
- `internal/application/usecase/session/qr.go`

**Arquivos Modificados:**
- `internal/adapters/http/handlers/session.go`

---

## 🎯 Sprint 3: Corrigir Context e Concorrência

**Objetivo:** Substituir context.Background() e eliminar time.Sleep()  
**Prioridade:** 🟡 ALTA  
**Estimativa:** 8-10 horas

### Tarefas:

- [ ] **3.1** Identificar todos os context.Background()
  - Buscar e listar 9 ocorrências em waclient
  - Documentar cada uso

- [ ] **3.2** Propagar contexto em handleQRCode
  - Usar contexto do client
  - Remover `context.Background()`

- [ ] **3.3** Propagar contexto em updateSessionStatus
  - Receber contexto da requisição
  - Propagar em todas as chamadas

- [ ] **3.4** Propagar contexto em event handlers
  - Passar contexto para `handleConnected`, `handleDisconnected`, etc
  - Garantir cancelamento funciona

- [ ] **3.5** Eliminar time.Sleep em handlers
  - Substituir `time.Sleep(500ms)` por polling
  - Usar `context.WithTimeout`

- [ ] **3.6** Eliminar time.Sleep em qr_manager
  - Substituir por channels ou polling
  - Implementar timeout adequado

- [ ] **3.7** Implementar polling com timeout para QR
  - Criar função `waitForQRCode`
  - Usar `context.WithTimeout`

- [ ] **3.8** Adicionar WaitGroup para goroutines
  - Controlar lifecycle em `loadSessionsFromDatabase`
  - Controlar lifecycle em `autoReconnect`

- [ ] **3.9** Implementar graceful shutdown
  - Adicionar método `Shutdown()`
  - Esperar goroutines terminarem

- [ ] **3.10** Testar cancelamento de contexto
  - Validar que cancelamento propaga
  - Testar timeout funciona

**Arquivos Afetados:**
- `internal/adapters/waclient/client.go`
- `internal/adapters/waclient/qr_manager.go`
- `internal/adapters/http/handlers/session.go`

---

## 🎯 Sprint 4: Implementar Transações

**Objetivo:** Adicionar transações em operações críticas e melhorar validações  
**Prioridade:** 🟡 ALTA  
**Estimativa:** 8-10 horas

### Tarefas:

- [ ] **4.1** Criar Unit of Work pattern
  - Implementar `UnitOfWork` em `database/`
  - Gerenciar transações

- [ ] **4.2** Adicionar transação em CreateSession
  - Envolver criação (domain + waclient) em transação
  - Garantir atomicidade

- [ ] **4.3** Adicionar transação em DeleteSession
  - Envolver deleção (domain + waclient) em transação
  - Garantir atomicidade

- [ ] **4.4** Implementar rollback automático
  - Garantir rollback em caso de erro
  - Usar `defer tx.Rollback()`

- [ ] **4.5** Criar validators em application/validators
  - Implementar `SessionValidator`
  - Validações robustas

- [ ] **4.6** Validar tamanho de strings
  - Adicionar validação min/max length
  - Validar `name`, `deviceJID`, etc

- [ ] **4.7** Validar formato de IDs
  - Validar que IDs são UUIDs válidos
  - Rejeitar IDs inválidos

- [ ] **4.8** Validar caracteres permitidos
  - Sanitizar entrada
  - Validar caracteres em campos de texto

- [ ] **4.9** Adicionar validação de unicidade
  - Validar que nome de sessão é único
  - Verificar antes de criar

- [ ] **4.10** Testar transações e validações
  - Criar testes para atomicidade
  - Testar validações funcionando

**Arquivos Criados:**
- `internal/adapters/database/uow.go`
- `internal/application/validators/session_validator.go`

**Arquivos Modificados:**
- `internal/application/usecase/session/create.go`
- `internal/application/usecase/session/delete.go`

---

## 🎯 Sprint 5: Validação Final e Documentação

**Objetivo:** Validar todas as refatorações, criar testes e atualizar documentação  
**Prioridade:** 🟢 MÉDIA  
**Estimativa:** 8-10 horas

### Tarefas:

- [ ] **5.1** Executar todos os testes unitários
  - Rodar suite completa
  - Garantir 100% de sucesso

- [ ] **5.2** Executar testes de integração
  - Testar integração entre camadas
  - Validar fluxo completo

- [ ] **5.3** Validar cobertura de testes > 80%
  - Executar `go test -cover`
  - Adicionar testes se necessário

- [ ] **5.4** Executar linter sem warnings
  - Rodar `golangci-lint`
  - Corrigir todos os warnings

- [ ] **5.5** Validar compilação sem erros
  - `go build` deve compilar
  - Sem erros ou warnings

- [ ] **5.6** Testar em ambiente de staging
  - Deploy em staging
  - Executar testes end-to-end

- [ ] **5.7** Atualizar ARCHITECTURE.md
  - Documentar nova arquitetura
  - Adicionar diagramas atualizados

- [ ] **5.8** Atualizar README.md
  - Atualizar documentação de setup
  - Atualizar exemplos de uso

- [ ] **5.9** Criar migration guide
  - Documentar mudanças breaking
  - Explicar como migrar

- [ ] **5.10** Code review final
  - Revisar todo o código refatorado
  - Aprovar com equipe

---

## 📈 Estimativas e Cronograma

| Sprint | Tarefas | Estimativa | Prioridade |
|--------|---------|------------|------------|
| Sprint 1 | 7 | 10-12h | 🔴 Crítica |
| Sprint 2 | 10 | 12-16h | 🔴 Crítica |
| Sprint 3 | 10 | 8-10h | 🟡 Alta |
| Sprint 4 | 10 | 8-10h | 🟡 Alta |
| Sprint 5 | 10 | 8-10h | 🟢 Média |
| **Total** | **47** | **38-48h** | - |

**Cronograma Sugerido (1 dev full-time):**
- Semana 1-2: Sprint 1 + Sprint 2
- Semana 3-4: Sprint 3 + Sprint 4
- Semana 5: Sprint 5

**Cronograma Alternativo (2 devs):**
- Semana 1: Sprint 1 + Sprint 2 (paralelo)
- Semana 2: Sprint 3 + Sprint 4 (paralelo)
- Semana 3: Sprint 5

---

## ✅ Critérios de Aceitação

### Sprint 1:
- [ ] Sem duplicação de `Session`/`SessionInfo`
- [ ] Sem duplicação de `Status`/`SessionStatus`
- [ ] `DBSessionManager` removido
- [ ] Código compila sem erros

### Sprint 2:
- [ ] Use Cases implementados e testados
- [ ] Handlers dependem de Use Cases
- [ ] Sem import de `waclient` em handlers
- [ ] Testes unitários passando

### Sprint 3:
- [ ] Sem `context.Background()` em produção
- [ ] Sem `time.Sleep()` em produção
- [ ] Graceful shutdown implementado
- [ ] Cancelamento de contexto funciona

### Sprint 4:
- [ ] Transações implementadas
- [ ] Rollback automático funciona
- [ ] Validações robustas implementadas
- [ ] Testes de transação passando

### Sprint 5:
- [ ] Cobertura de testes > 80%
- [ ] Linter sem warnings
- [ ] Documentação atualizada
- [ ] Code review aprovado

---

## 🚀 Próximos Passos

1. ✅ Revisar tasklist com equipe técnica
2. ✅ Alocar recursos (1-2 devs)
3. ✅ Criar branch de refatoração
4. ✅ Iniciar Sprint 1
5. ✅ Daily standups para acompanhamento

---

**Documentos Relacionados:**
- `CODE_ANALYSIS_REPORT.md` - Análise completa
- `ARCHITECTURE_ANALYSIS.md` - Análise arquitetural
- `CODE_ANALYSIS_SUMMARY.md` - Resumo executivo

