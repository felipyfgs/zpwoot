# 📊 Melhorias de Logging - zpwoot 2025

## 🎯 Objetivo

Modernizar o sistema de logging do zpwoot seguindo as melhores práticas de 2025, mantendo a formalidade mas facilitando a leitura, análise e observabilidade profunda.

---

## 📋 Análise Atual

### ✅ Pontos Fortes
- ✓ Uso do zerolog (alta performance)
- ✓ Logs estruturados com campos contextuais
- ✓ Timestamps padronizados (RFC3339)
- ✓ Separação por níveis (DEBUG, INFO, WARN, ERROR)
- ✓ Package tracking automático

### ⚠️ Pontos de Melhoria

#### 1. **Payloads JSON Excessivamente Grandes**
**Problema Atual:**
```
2025-10-08T14:55:04Z INF events.go:101 > Message received chat=559988527530@s.whatsapp.net from=559981769536:83@s.whatsapp.net from_me=true is_group=false payload={"message":{"conversation":"Bom dia","messageContextInfo":{...MUITO GRANDE...}}}
```

**Impacto:**
- Dificulta leitura rápida
- Aumenta custos de armazenamento
- Complica análise visual

#### 2. **Falta de Correlation IDs**
**Problema:** Impossível rastrear uma mensagem através de múltiplos componentes

#### 3. **Dados Sensíveis Expostos**
**Problema:** Números de telefone completos, nomes, conteúdo de mensagens nos logs

#### 4. **Mensagens Genéricas**
**Problema:** "Message received", "Receipt received" - pouco contexto de negócio

#### 5. **Falta de Métricas de Performance**
**Problema:** Sem duração, latência, ou métricas de processamento

---

## 🚀 Melhorias Recomendadas

### 1. **Structured Logging Aprimorado**

#### Antes:
```
2025-10-08T14:55:04Z INF Message received chat=559988527530@s.whatsapp.net payload={...}
```

#### Depois:
```json
{
  "timestamp": "2025-10-08T14:55:04Z",
  "level": "info",
  "service": "zpwoot",
  "component": "waclient",
  "event": "message.received",
  "correlation_id": "msg_abc123def456",
  "session_id": "a66e10ce-91ef-41d8-ba30-436e693ca10b",
  "message": {
    "id": "3EB0136EE801BBAFA08822",
    "type": "text",
    "direction": "outgoing",
    "chat_type": "private"
  },
  "sender": {
    "id_hash": "sha256:abc123...",
    "is_me": true
  },
  "metrics": {
    "processing_time_ms": 45,
    "queue_depth": 12
  },
  "context": {
    "trace_id": "trace_xyz789",
    "span_id": "span_123"
  }
}
```

**Benefícios:**
- ✅ Dados sensíveis mascarados
- ✅ Métricas de performance incluídas
- ✅ Correlation ID para rastreamento
- ✅ Contexto de negócio claro
- ✅ Payload resumido (detalhes em DEBUG)

---

### 2. **Níveis de Log Otimizados**

#### **INFO** - Eventos de Negócio Importantes
```json
{
  "level": "info",
  "event": "message.sent",
  "message_id": "msg_123",
  "chat_type": "group",
  "delivery_status": "sent",
  "duration_ms": 234
}
```

#### **DEBUG** - Detalhes Técnicos (com payload completo)
```json
{
  "level": "debug",
  "event": "message.payload",
  "message_id": "msg_123",
  "payload": { /* payload completo aqui */ }
}
```

#### **WARN** - Situações Anormais mas Recuperáveis
```json
{
  "level": "warn",
  "event": "qr.expired",
  "session_id": "sess_123",
  "retry_count": 2,
  "next_retry_in_seconds": 30
}
```

#### **ERROR** - Erros que Requerem Atenção
```json
{
  "level": "error",
  "event": "webhook.send.failed",
  "error": {
    "code": "CONNECTION_TIMEOUT",
    "message": "Failed to reach webhook endpoint",
    "retry_count": 3
  },
  "webhook_url_hash": "sha256:xyz...",
  "duration_ms": 5000
}
```

---

### 3. **Mascaramento de Dados Sensíveis**

#### Implementação:

```go
// internal/adapters/logger/sanitizer.go
package logger

import (
    "crypto/sha256"
    "encoding/hex"
    "regexp"
)

var (
    phoneRegex = regexp.MustCompile(`\d{10,15}`)
    emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
)

// HashSensitiveID cria hash de IDs sensíveis
func HashSensitiveID(id string) string {
    hash := sha256.Sum256([]byte(id))
    return "sha256:" + hex.EncodeToString(hash[:])[:16]
}

// MaskPhone mascara número de telefone
func MaskPhone(phone string) string {
    if len(phone) < 4 {
        return "***"
    }
    return "***" + phone[len(phone)-4:]
}

// SanitizeMessage remove conteúdo sensível de mensagens
func SanitizeMessage(msg string) string {
    msg = phoneRegex.ReplaceAllString(msg, "[PHONE]")
    msg = emailRegex.ReplaceAllString(msg, "[EMAIL]")
    return msg
}
```

---

### 4. **Correlation IDs e Trace Context**

#### Implementação:

```go
// internal/adapters/logger/context.go
package logger

import (
    "context"
    "github.com/google/uuid"
)

type contextKey string

const (
    correlationIDKey contextKey = "correlation_id"
    traceIDKey       contextKey = "trace_id"
    spanIDKey        contextKey = "span_id"
)

// WithCorrelationID adiciona correlation ID ao contexto
func WithCorrelationID(ctx context.Context) context.Context {
    correlationID := uuid.New().String()
    return context.WithValue(ctx, correlationIDKey, correlationID)
}

// GetCorrelationID extrai correlation ID do contexto
func GetCorrelationID(ctx context.Context) string {
    if id, ok := ctx.Value(correlationIDKey).(string); ok {
        return id
    }
    return ""
}

// LoggerWithContext cria logger com contexto completo
func LoggerWithContext(ctx context.Context) *Logger {
    logger := New()
    
    if correlationID := GetCorrelationID(ctx); correlationID != "" {
        logger = logger.WithField("correlation_id", correlationID)
    }
    
    if traceID := GetTraceID(ctx); traceID != "" {
        logger = logger.WithField("trace_id", traceID)
    }
    
    return logger
}
```

---

### 5. **Métricas de Performance**

#### Implementação:

```go
// internal/adapters/logger/metrics.go
package logger

import (
    "time"
)

// PerformanceMetrics armazena métricas de performance
type PerformanceMetrics struct {
    StartTime     time.Time
    Duration      time.Duration
    MemoryUsageMB float64
    QueueDepth    int
}

// StartMetrics inicia medição de performance
func StartMetrics() *PerformanceMetrics {
    return &PerformanceMetrics{
        StartTime: time.Now(),
    }
}

// End finaliza medição e retorna duração
func (pm *PerformanceMetrics) End() time.Duration {
    pm.Duration = time.Since(pm.StartTime)
    return pm.Duration
}

// ToFields converte métricas para campos de log
func (pm *PerformanceMetrics) ToFields() map[string]interface{} {
    return map[string]interface{}{
        "duration_ms":     pm.Duration.Milliseconds(),
        "memory_usage_mb": pm.MemoryUsageMB,
        "queue_depth":     pm.QueueDepth,
    }
}
```

---

## 📝 Exemplos Práticos de Uso

### Exemplo 1: Log de Mensagem Recebida (Melhorado)

```go
// internal/adapters/waclient/events.go

func (eh *DefaultEventHandler) handleMessage(client *Client, evt *events.Message) error {
    metrics := logger.StartMetrics()
    ctx := logger.WithCorrelationID(context.Background())
    log := logger.LoggerWithContext(ctx)
    
    // Processar mensagem...
    
    metrics.End()
    
    // Log INFO com resumo
    log.Info().
        Str("event", "message.received").
        Str("message_id", evt.Info.ID).
        Str("message_type", getMessageType(evt.Message)).
        Str("chat_type", getChatType(evt.Info.IsGroup)).
        Str("direction", getDirection(evt.Info.IsFromMe)).
        Str("sender_hash", logger.HashSensitiveID(evt.Info.Sender.String())).
        Str("session_id", client.SessionID).
        Int64("duration_ms", metrics.Duration.Milliseconds()).
        Msg("WhatsApp message processed")
    
    // Log DEBUG com payload completo (apenas em desenvolvimento)
    log.Debug().
        Str("event", "message.payload").
        Str("message_id", evt.Info.ID).
        Interface("payload", evt.Message).
        Msg("Message payload details")
    
    return eh.sendWebhookIfEnabled(client, EventMessage, webhookData)
}
```

### Exemplo 2: Log de Erro com Contexto

```go
func (wac *WAClient) sendWebhook(client *Client, eventType EventType, event interface{}) {
    metrics := logger.StartMetrics()
    
    err := wac.webhookSender.SendWebhook(client.ctx, webhookEvent)
    metrics.End()
    
    if err != nil {
        wac.logger.Error().
            Err(err).
            Str("event", "webhook.send.failed").
            Str("event_type", string(eventType)).
            Str("webhook_url_hash", logger.HashSensitiveID(client.WebhookURL)).
            Str("session_id", client.SessionID).
            Int64("duration_ms", metrics.Duration.Milliseconds()).
            Int("retry_count", 3).
            Msg("Failed to deliver webhook")
        return
    }
    
    wac.logger.Info().
        Str("event", "webhook.sent").
        Str("event_type", string(eventType)).
        Str("session_id", client.SessionID).
        Int64("duration_ms", metrics.Duration.Milliseconds()).
        Msg("Webhook delivered successfully")
}
```

---

## 🔧 Configuração Recomendada

### Variáveis de Ambiente

```bash
# Produção
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
LOG_SAMPLING_RATE=0.1  # 10% de logs INFO, 100% de ERROR
LOG_MASK_SENSITIVE=true

# Desenvolvimento
LOG_LEVEL=debug
LOG_FORMAT=console
LOG_OUTPUT=stderr
LOG_SAMPLING_RATE=1.0  # 100% dos logs
LOG_MASK_SENSITIVE=false
```

---

## 📊 Benefícios Esperados

### 1. **Redução de Custos**
- ⬇️ 60-80% menos armazenamento (com sampling e compressão)
- ⬇️ 40% menos tráfego de rede

### 2. **Melhor Observabilidade**
- ✅ Rastreamento end-to-end com correlation IDs
- ✅ Métricas de performance em tempo real
- ✅ Alertas baseados em padrões de log

### 3. **Segurança e Compliance**
- 🔒 Dados sensíveis mascarados automaticamente
- 🔒 Conformidade com LGPD/GDPR
- 🔒 Auditoria completa de eventos

### 4. **Produtividade do Time**
- ⚡ 60-70% menos tempo de debugging
- ⚡ Análise de problemas mais rápida
- ⚡ Melhor compreensão do comportamento do sistema

---

## 🎯 Próximos Passos

1. ✅ **Fase 1:** Implementar sanitização de dados sensíveis
2. ✅ **Fase 2:** Adicionar correlation IDs e trace context
3. ✅ **Fase 3:** Implementar métricas de performance
4. ✅ **Fase 4:** Otimizar níveis de log e sampling
5. ✅ **Fase 5:** Integrar com OpenTelemetry (opcional)

---

## 📚 Referências

- [Uptrace - Structured Logging Best Practices 2025](https://uptrace.dev/glossary/structured-logging)
- [OpenTelemetry Logging Specification](https://opentelemetry.io/docs/specs/otel/logs/)
- [Zerolog Documentation](https://github.com/rs/zerolog)
- [Go Logging Best Practices](https://www.reddit.com/r/golang/comments/180jnpd/best_practice_for_logging/)

