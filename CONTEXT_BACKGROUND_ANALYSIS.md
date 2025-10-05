# 🔍 Análise de context.Background() - Sprint 3

**Total de Ocorrências:** 17  
**Arquivos Afetados:** 2

---

## 📊 Ocorrências por Arquivo

### **internal/adapters/waclient/client.go** (16 ocorrências)

| Linha | Função | Contexto | Prioridade |
|-------|--------|----------|------------|
| 55 | `loadSessionsFromDatabase()` | Carregamento inicial | 🟡 MÉDIA |
| 93 | `autoReconnect()` | Auto-reconnect | 🔴 ALTA |
| 257 | `handleQRCode()` | GetQRChannel | 🔴 ALTA |
| 260 | `handleQRCode()` | updateSessionStatus | 🔴 ALTA |
| 266 | `handleQRCode()` | updateSessionStatus | 🔴 ALTA |
| 274 | `handleQRCode()` | updateSessionStatus | 🔴 ALTA |
| 426 | `handleConnected()` | updateSessionStatus | 🔴 ALTA |
| 435 | `handleDisconnected()` | updateSessionStatus | 🔴 ALTA |
| 444 | `handleLoggedOut()` | updateSessionStatus | 🔴 ALTA |
| 454 | `handleStreamReplaced()` | updateSessionStatus | 🔴 ALTA |
| 485 | `sendEvent()` | SendWebhook | 🟡 MÉDIA |
| 511 | `handleQRCode()` | updateSessionStatus | 🔴 ALTA |
| 526 | `handleConnected()` | updateSessionStatus | 🔴 ALTA |
| 541 | `handleDisconnected()` | updateSessionStatus | 🔴 ALTA |
| 548 | `handleLoggedOut()` | updateSessionStatus | 🔴 ALTA |
| 555 | `handleStreamReplaced()` | updateSessionStatus | 🔴 ALTA |
| 610 | `NewWAStoreContainer()` | Inicialização | 🟢 BAIXA |

### **internal/adapters/waclient/event_handler.go** (1 ocorrência)

| Linha | Função | Contexto | Prioridade |
|-------|--------|----------|------------|
| 233 | `handleMessage()` | WithTimeout | 🟡 MÉDIA |

---

## 🎯 Estratégia de Correção

### **Grupo 1: Event Handlers (12 ocorrências - ALTA prioridade)**
Funções que são chamadas por eventos do WhatsApp:
- `handleQRCode()` - 5 ocorrências
- `handleConnected()` - 2 ocorrências
- `handleDisconnected()` - 2 ocorrências
- `handleLoggedOut()` - 2 ocorrências
- `handleStreamReplaced()` - 2 ocorrências

**Solução:**
- Adicionar campo `ctx context.Context` no struct `Client`
- Propagar contexto do client para event handlers
- Usar `client.ctx` ao invés de `context.Background()`

### **Grupo 2: Goroutines de Background (2 ocorrências - MÉDIA prioridade)**
- `loadSessionsFromDatabase()` - linha 55
- `sendEvent()` - linha 485

**Solução:**
- Criar contexto de aplicação no `WAClient`
- Usar contexto de aplicação para operações de background
- Implementar graceful shutdown

### **Grupo 3: Inicialização (1 ocorrência - BAIXA prioridade)**
- `NewWAStoreContainer()` - linha 610

**Solução:**
- Manter `context.Background()` (aceitável em inicialização)
- Ou receber contexto como parâmetro

### **Grupo 4: Timeout Wrapper (1 ocorrência - MÉDIA prioridade)**
- `handleMessage()` - linha 233

**Solução:**
- Usar `context.WithTimeout(client.ctx, 30*time.Second)`

---

## 📝 Plano de Ação

### **Passo 1: Adicionar contexto ao Client**
```go
type Client struct {
    // ... campos existentes
    ctx    context.Context
    cancel context.CancelFunc
}
```

### **Passo 2: Propagar contexto em event handlers**
```go
// Antes
wac.updateSessionStatus(context.Background(), client)

// Depois
wac.updateSessionStatus(client.ctx, client)
```

### **Passo 3: Adicionar contexto de aplicação ao WAClient**
```go
type WAClient struct {
    // ... campos existentes
    appCtx context.Context
}
```

### **Passo 4: Implementar graceful shutdown**
```go
func (wac *WAClient) Shutdown(ctx context.Context) error {
    // Cancelar todos os clients
    // Esperar goroutines terminarem
}
```

---

## ✅ Critérios de Sucesso

- [ ] Nenhum `context.Background()` em event handlers
- [ ] Contexto propagado corretamente em todas as chamadas
- [ ] Cancelamento funciona quando requisição é cancelada
- [ ] Graceful shutdown implementado
- [ ] Código compila sem erros

