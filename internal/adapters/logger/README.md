# Logger Centralizado - zpwoot

Este é o sistema de logging centralizado do zpwoot, baseado no [zerolog](https://github.com/rs/zerolog) e seguindo o padrão de uso direto da biblioteca.

## Características

- **Logging Estruturado**: Suporte completo a JSON e console formatado
- **Configuração por Ambiente**: Controlado via variáveis de ambiente
- **Padrão Zerolog**: Segue exatamente o padrão de uso do zerolog
- **Context Aware**: Suporte a campos contextuais (request_id, session_id, component, etc.)
- **Performance**: Baseado no zerolog, uma das bibliotecas de logging mais rápidas do Go

## Configuração

O logger é configurado através das seguintes variáveis de ambiente:

```bash
LOG_LEVEL=info          # trace, debug, info, warn, error, fatal, panic, disabled
LOG_FORMAT=console      # console, json
LOG_OUTPUT=stderr       # stdout, stderr
```

## Inicialização

### Inicialização Básica
```go
import "zpwoot/internal/adapters/logger"

// Inicialização simples
logger.Init("info")
```

### Inicialização com Configuração Completa
```go
import (
    "zpwoot/internal/adapters/config"
    "zpwoot/internal/adapters/logger"
)

cfg := config.Load()
logger.InitWithConfig(cfg)
```

## Uso Básico

### Logging Direto (Padrão Zerolog)
```go
import (
    "zpwoot/internal/adapters/logger"
    "github.com/rs/zerolog"
)

// Usando funções globais
logger.Info().Msg("info message")
logger.Debug().Msg("debug message")
logger.Warn().Msg("warn message")
logger.Error().Msg("error message")
logger.WithLevel(zerolog.FatalLevel).Msg("fatal message")

// Com campos estruturados
logger.Info().
    Str("user_id", "12345").
    Str("action", "login").
    Int("attempts", 3).
    Msg("User logged in")

// Com erro
err := errors.New("something went wrong")
logger.Error().
    Err(err).
    Str("component", "database").
    Msg("Database operation failed")
```

### Usando Instância do Logger
```go
log := logger.New()

log.Info().Msg("info message")
log.Debug().Str("key", "value").Msg("debug with field")
log.Error().Err(err).Msg("error with details")
```

## Logging Contextual

### Criando Loggers com Contexto
```go
// Logger com componente
authLogger := logger.WithComponent("auth")
authLogger.Info().Msg("Authentication started")

// Logger com múltiplos campos
contextLogger := logger.WithFields(map[string]interface{}{
    "request_id": "req-123",
    "user_id":    "user-456",
    "session_id": "sess-789",
})
contextLogger.Info().Msg("Processing request")

// Encadeamento de contexto
apiLogger := logger.WithComponent("api").
    WithRequestID("req-123").
    WithSessionID("sess-456")
apiLogger.Info().Str("endpoint", "/users").Msg("API call")
```

### Métodos de Contexto Disponíveis
```go
logger.WithComponent("component_name")
logger.WithRequestID("request_id")
logger.WithSessionID("session_id")
logger.WithError(err)
logger.WithFields(map[string]interface{}{...})
```

## Formatos de Saída

### Console (Desenvolvimento)
```bash
LOG_FORMAT=console
```
Saída colorida e formatada para desenvolvimento:
```
2025-10-05T12:17:28Z INF main.go:25 > Starting zpwoot application component=main pkg=main
```

### JSON (Produção)
```bash
LOG_FORMAT=json
```
Saída estruturada em JSON para produção:
```json
{"level":"info","time":"2025-10-05T12:17:28Z","caller":"main.go:25","pkg":"main","service":"zpwoot","version":"1.0.0","message":"Starting zpwoot application"}
```

## Níveis de Log

- **trace**: Informações muito detalhadas para debugging
- **debug**: Informações de debugging
- **info**: Informações gerais
- **warn**: Avisos que não impedem o funcionamento
- **error**: Erros que não param a aplicação
- **fatal**: Erros críticos que param a aplicação
- **panic**: Erros que causam panic
- **disabled**: Desabilita todos os logs

## Campos Automáticos

O logger adiciona automaticamente os seguintes campos:

- **timestamp**: Timestamp do log
- **level**: Nível do log
- **caller**: Arquivo e linha onde o log foi chamado (formato limpo: `arquivo.go:linha`)
- **pkg**: Package/módulo de origem do log (detectado automaticamente)
- **service**: Nome do serviço (zpwoot)
- **version**: Versão da aplicação (apenas em formato JSON)

## Métodos de Conveniência

Para casos simples onde você só precisa logar uma mensagem:

```go
log := logger.New()

log.InfoMsg("Simple info message")
log.DebugMsg("Simple debug message")
log.ErrorMsg("Simple error message")
```

## Acesso Direto ao Zerolog

Se precisar de funcionalidades avançadas do zerolog:

```go
zerologLogger := logger.GetZerologLogger()
zerologLogger.Info().Dict("user", zerolog.Dict().
    Str("name", "John").
    Int("age", 30)).
    Msg("Complex structured log")
```

## Exemplo Completo

```go
package main

import (
    "errors"
    "zpwoot/internal/adapters/config"
    "zpwoot/internal/adapters/logger"
    "github.com/rs/zerolog"
)

func main() {
    // Inicializar logger
    cfg := config.Load()
    logger.InitWithConfig(cfg)

    // Logging básico
    logger.Info().Msg("Application started")
    
    // Logging com campos
    logger.Info().
        Str("version", "1.0.0").
        Int("port", 8080).
        Msg("Server configuration")
    
    // Logging contextual
    apiLogger := logger.WithComponent("api")
    apiLogger.Info().
        Str("method", "GET").
        Str("path", "/health").
        Int("status", 200).
        Msg("API request processed")
    
    // Logging de erro
    err := errors.New("database connection failed")
    logger.Error().
        Err(err).
        Str("component", "database").
        Str("host", "localhost").
        Int("port", 5432).
        Msg("Failed to connect to database")
}
```

## Integração com o Projeto

O logger está integrado com o sistema de configuração do zpwoot e pode ser usado em qualquer parte da aplicação após a inicialização no `main.go`.
