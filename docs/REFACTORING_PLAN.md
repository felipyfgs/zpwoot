# 🔧 Plano de Refatoração de Nomenclatura - zpwoot

## 📊 Análise de Nomes Problemáticos

### 🔴 Prioridade ALTA - Nomes Excessivamente Longos

| Nome Atual | Problema | Nome Proposto | Justificativa |
|------------|----------|---------------|---------------|
| `generateFallbackMessageID()` | Redundante, contexto já é mensagem | `fallbackID()` | Pacote já indica contexto |
| `initializeWhatsAppClient()` | Verbo longo | `initWAClient()` | "init" é idiomático |
| `SessionManagerAdapter` | Redundante "Adapter" | `SessionManager` | Já está em pacote adapter |
| `MessageServiceWrapper` | "Wrapper" não idiomático | `MessageService` | Implementação direta |
| `MessageSenderImpl` | "Impl" não idiomático | `Sender` | Pacote já indica contexto |
| `DefaultEventHandler` | "Default" redundante | `EventHandler` | Única implementação |
| `WAClientAdapter` | Redundante "Adapter" | `Client` | Já está em pacote adapter |
| `NewWAStoreContainer()` | Nome muito longo | `NewStore()` | Contexto claro |
| `GetSessionUseCases()` | Redundante "UseCases" | `Sessions()` | Getter idiomático |
| `GetMessageUseCases()` | Redundante "UseCases" | `Messages()` | Getter idiomático |
| `GetWhatsAppClient()` | Redundante "WhatsApp" | `WAClient()` | Abreviação comum |
| `ExecuteWithAutoReconnect()` | Nome muito longo | `ConnectAuto()` | Mais direto |
| `ExecuteWithValidation()` | Redundante "Execute" | `DeleteForce()` | Mais claro |
| `ExecuteForce()` | Redundante "Execute" | `DisconnectForce()` | Mais claro |
| `ExecuteSimple()` | Redundante "Execute" | `List()` | Mais simples |
| `ExecuteWithFilter()` | Redundante "Execute" | `ListFiltered()` | Mais claro |
| `ProcessIncomingMessage()` | Redundante "Incoming" | `Process()` | Contexto claro |
| `ProcessIncomingMessageBatch()` | Nome muito longo | `ProcessBatch()` | Mais curto |
| `SendContactMessageFromInput()` | Nome muito longo | `SendContact()` | Overload |
| `GetChatInfoAsInput()` | Sufixo redundante | `ChatInfo()` | Tipo de retorno já indica |
| `GetChatsAsInput()` | Sufixo redundante | `Chats()` | Tipo de retorno já indica |
| `GetContactsAsInput()` | Sufixo redundante | `Contacts()` | Tipo de retorno já indica |
| `ToInterfacesContactInfo()` | Nome muito longo | `ToOutput()` | Mais genérico |
| `ToInterfacesLocation()` | Nome muito longo | `ToOutput()` | Mais genérico |
| `ToInterfacesMediaData()` | Nome muito longo | `ToOutput()` | Mais genérico |
| `ToOutputMediaData()` | Redundante "Output" | `ToOutput()` | Padronizado |
| `ToMediaData()` | Ambíguo | `ToOutput()` | Padronizado |
| `buildMigrationObjects()` | Verbo longo | `buildMigrations()` | Mais curto |
| `categorizeMigrationFile()` | Verbo longo | `categorizeFile()` | Contexto claro |
| `extractVersionFromFilename()` | Nome muito longo | `extractVersion()` | Mais curto |
| `isMigrationApplied()` | Redundante "Migration" | `isApplied()` | Contexto claro |
| `processMigrationFiles()` | Redundante "Migration" | `processFiles()` | Contexto claro |
| `readMigrationDirectory()` | Redundante "Migration" | `readDir()` | Mais curto |
| `readMigrationFile()` | Redundante "Migration" | `readFile()` | Mais curto |
| `GetMigrationStatus()` | Redundante "Migration" | `Status()` | Getter idiomático |
| `UpdateSessionStatus()` | Redundante "Session" | `UpdateStatus()` | Contexto claro |
| `UpdateQRCode()` | OK | `UpdateQR()` | Mais curto |
| `SetQRCode()` | OK | `SetQR()` | Mais curto |
| `ClearQRCode()` | OK | `ClearQR()` | Mais curto |
| `GetQRCode()` | OK | `QRCode()` | Getter idiomático |
| `GetQRCodeForSession()` | Redundante "ForSession" | `QRCode()` | Contexto claro |
| `RefreshQRCode()` | OK | `RefreshQR()` | Mais curto |
| `CheckQRCodeStatus()` | Redundante "QRCode" | `CheckQR()` | Mais curto |
| `hasValidQRCode()` | Redundante "QRCode" | `hasValidQR()` | Mais curto |
| `isQRCodeExpired()` | Redundante "QRCode" | `isQRExpired()` | Mais curto |
| `GenerateQRCodeBase64()` | Nome muito longo | `QRBase64()` | Mais curto |
| `CleanupExpiredQRCodes()` | Redundante "QRCodes" | `CleanupQRs()` | Mais curto |
| `StartQRCleanupRoutine()` | Nome muito longo | `StartQRCleanup()` | Mais curto |
| `processAllQRCodesFromEvent()` | Nome muito longo | `processQRs()` | Mais curto |
| `waitForQRCode()` | OK | `waitForQR()` | Mais curto |
| `waitForQRCodeWithTimeout()` | Nome muito longo | `waitForQR()` | Overload |
| `clearQRCode()` | OK | `clearQR()` | Mais curto |

### 🟡 Prioridade MÉDIA - Nomes Redundantes

| Nome Atual | Problema | Nome Proposto |
|------------|----------|---------------|
| `NewMessageUseCases()` | Redundante "UseCases" | `NewMessages()` |
| `NewSessionManagerAdapter()` | Redundante "Adapter" | `NewManager()` |
| `NewWAClientAdapter()` | Redundante "Adapter" | `NewClient()` |
| `NewDefaultEventHandler()` | Redundante "Default" | `NewEventHandler()` |
| `NewMessageServiceWrapper()` | Redundante "Wrapper" | `NewService()` |
| `SessionToCreateResponse()` | Prefixo redundante | `ToCreateResponse()` |
| `SessionToDetailResponse()` | Prefixo redundante | `ToDetailResponse()` |
| `SessionToListInfo()` | Prefixo redundante | `ToListInfo()` |
| `SessionToListResponse()` | Prefixo redundante | `ToListResponse()` |
| `SessionToStatusResponse()` | Prefixo redundante | `ToStatusResponse()` |
| `FromDomainSession()` | Prefixo redundante | `FromDomain()` |
| `NewQRCodeResponse()` | Redundante "QRCode" | `NewQRResponse()` |

### 🟢 Prioridade BAIXA - Melhorias Estéticas

| Nome Atual | Nome Proposto |
|------------|---------------|
| `GetServerAddress()` | `Address()` |
| `IsDevelopment()` | `IsDev()` |
| `IsProduction()` | `IsProd()` |
| `GetDeviceJID()` | `DeviceJID()` |
| `GetZerologLogger()` | `Zerolog()` |
| `GetGlobalLogger()` | `Global()` |

## 📋 Ordem de Execução

### Fase 1: Funções Utilitárias (Baixo Impacto)
1. ✅ `generateFallbackMessageID` → `fallbackID`
2. ✅ `GenerateQRCodeBase64` → `QRBase64`
3. ✅ Funções de QR Code (todas)

### Fase 2: Adapters (Médio Impacto)
4. ✅ `SessionManagerAdapter` → `SessionManager`
5. ✅ `WAClientAdapter` → `Client`
6. ✅ `MessageServiceWrapper` → `MessageService`
7. ✅ `MessageSenderImpl` → `Sender`
8. ✅ `DefaultEventHandler` → `EventHandler`

### Fase 3: Use Cases (Alto Impacto)
9. ✅ Métodos `Execute*` → Nomes específicos
10. ✅ Getters redundantes

### Fase 4: Container e Inicialização
11. ✅ `initializeWhatsAppClient` → `initWAClient`
12. ✅ Getters do Container

### Fase 5: Conversores e Helpers
13. ✅ Métodos `To*` redundantes
14. ✅ Funções de validação

## 🎯 Métricas Esperadas

| Métrica | Antes | Depois | Melhoria |
|---------|-------|--------|----------|
| Média de caracteres/nome | 28 | 18 | -36% |
| Nomes > 30 caracteres | 45 | 5 | -89% |
| Nomes com "Message" redundante | 23 | 0 | -100% |
| Nomes com "Session" redundante | 18 | 0 | -100% |
| Nomes com "QRCode" redundante | 12 | 0 | -100% |

## ✅ Regras de Refatoração

1. **Getters**: Remover prefixo "Get" quando possível
2. **Setters**: Manter prefixo "Set" (idiomático)
3. **Construtores**: Manter "New" (idiomático)
4. **Validadores**: Manter "Validate" (idiomático)
5. **Conversores**: Usar "To" + tipo destino
6. **Checkers**: Usar "Is" ou "Has" (idiomático)
7. **Handlers**: Usar "Handle" + evento
8. **Executores**: Remover "Execute", usar verbo direto

## 🚀 Próximos Passos

1. Começar pela Fase 1 (baixo impacto)
2. Compilar após cada mudança
3. Validar que tudo funciona
4. Prosseguir para próxima fase
5. Documentar mudanças significativas

