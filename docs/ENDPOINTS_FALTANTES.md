# Análise de Endpoints Faltantes - zpwoot vs Referência

## 📊 Resumo Executivo

**Total de Endpoints Aplicáveis:** 46
**Implementados no zpwoot:** 35
**Faltando:** 11
**Cobertura:** 76%

**Removidos (Não Aplicáveis ao zpwoot):** 16
- Webhooks (4) - Sistema de eventos diferente
- Download de Mídia (4) - Não necessário
- Admin Multi-tenant (5) - zpwoot usa sessões
- Configurações S3/Proxy (7) - Não aplicável

---

## ✅ Endpoints Implementados (35)

### Sessões (8)
- [x] Connect - Conectar sessão
- [x] Disconnect - Desconectar sessão
- [x] Logout - Logout e desvinculação
- [x] GetQR - Obter QR Code
- [x] GetStatus - Status da sessão (via Get)
- [x] Create - Criar sessão
- [x] List - Listar sessões
- [x] Delete - Deletar sessão

### Mensagens (14)
- [x] SendMessage - Enviar mensagem de texto
- [x] SendDocument - Enviar documento
- [x] SendAudio - Enviar áudio
- [x] SendImage - Enviar imagem
- [x] SendSticker - Enviar sticker
- [x] SendVideo - Enviar vídeo
- [x] SendContact - Enviar contato
- [x] SendLocation - Enviar localização
- [x] SendButtons - Enviar botões
- [x] SendList - Enviar lista
- [x] SendPoll - Enviar enquete
- [x] SendTemplate - Enviar template
- [x] SendContactsArray - Enviar múltiplos contatos
- [x] React - Reagir a mensagem (SendReaction)

### Grupos (13)
- [x] ListGroups - Listar grupos
- [x] GetGroupInfo - Info do grupo
- [x] GetGroupInviteLink - Link de convite
- [x] GroupJoin - Entrar no grupo
- [x] CreateGroup - Criar grupo
- [x] SetGroupLocked - Bloquear configurações
- [x] SetDisappearingTimer - Mensagens temporárias
- [x] RemoveGroupPhoto - Remover foto
- [x] UpdateGroupParticipants - Gerenciar participantes
- [x] GetGroupInviteInfo - Info via convite
- [x] SetGroupPhoto - Definir foto
- [x] SetGroupName - Alterar nome
- [x] SetGroupTopic - Alterar descrição
- [x] GroupLeave - Sair do grupo
- [x] SetGroupAnnounce - Modo anúncio

---

## ❌ Endpoints Faltantes (11)

### ⚠️ NOTA: Webhooks Removidos
O zpwoot implementará sistema de eventos próprio (WebSocket/SSE) ao invés de webhooks HTTP.
Webhooks da referência não serão implementados.

---

### 1. PairPhone (1 endpoint) 🟡 MÉDIA PRIORIDADE
```
❌ PairPhone - Pareamento por telefone (sem QR Code)
```

**Importância:** MÉDIA
**Complexidade:** BAIXA
**Descrição:** Permite autenticação usando número de telefone + código

**Referência:** Linha 572 em `referencia-handlers.bak`

**Exemplo:**
```go
linkingCode, err := client.PairPhone(
    context.Background(),
    phone,
    true,
    whatsmeow.PairClientChrome,
    "Chrome (Linux)"
)
```

---

### 2. Mensagens - Operações Avançadas (4 endpoints) 🟡 MÉDIA PRIORIDADE
```
❌ DeleteMessage      - Deletar mensagem enviada
❌ SendEditMessage    - Editar mensagem enviada
❌ MarkRead           - Marcar como lida
❌ RequestHistorySync - Sincronizar histórico
```

**Importância:** MÉDIA
**Complexidade:** BAIXA-MÉDIA
**Descrição:** Operações avançadas de mensagens

**Referência:**
- DeleteMessage: Linha 2041
- SendEditMessage: Linha 2106
- MarkRead: Linha 3169
- RequestHistorySync: Linha 2199

---

### 3. Contatos e Presença (6 endpoints) 🟡 MÉDIA PRIORIDADE
```
❌ CheckUser      - Verificar se número está no WhatsApp
❌ GetUser        - Obter informações do usuário
❌ SendPresence   - Enviar presença (online/offline/typing)
❌ ChatPresence   - Presença em chat específico
❌ GetAvatar      - Obter foto de perfil
❌ GetContacts    - Listar contatos
```

**Importância:** MÉDIA
**Complexidade:** BAIXA
**Descrição:** Gerenciamento de contatos e presença

**Referência:**
- CheckUser: Linha 2413
- GetUser: Linha 2479
- SendPresence: Linha 2546
- ChatPresence: Linha 2699
- GetAvatar: Linha 2602
- GetContacts: Linha 2669

---

---

## ⚪ Endpoints Removidos - Não Aplicáveis ao zpwoot (16)

### Webhooks (4 endpoints) - Sistema Diferente
```
⚪ GetWebhook, SetWebhook, UpdateWebhook, DeleteWebhook
```
**Motivo:** zpwoot implementará sistema de eventos próprio (WebSocket/SSE)

### Download de Mídia (4 endpoints) - Não Necessário
```
⚪ DownloadImage, DownloadDocument, DownloadVideo, DownloadAudio
```
**Motivo:** Mídias podem ser acessadas diretamente via whatsmeow

### Admin Multi-Tenant (5 endpoints) - Arquitetura Diferente
```
⚪ ListUsers, AddUser, EditUser, DeleteUser, DeleteUserComplete
```
**Motivo:** zpwoot usa sessões, não sistema multi-tenant

### Configurações S3/Proxy (7 endpoints) - Não Aplicável
```
⚪ SetHistory, GetHistory, SetProxy
⚪ ConfigureS3, GetS3Config, TestS3Connection, DeleteS3Config
```
**Motivo:** Configurações não necessárias na arquitetura do zpwoot

---

## 🎯 Priorização de Implementação

### Sprint 1: Mensagens Avançadas 🟡 (IMPORTANTE)
**Tempo Estimado:** 2-3 dias

1. DeleteMessage
2. SendEditMessage
3. MarkRead
4. PairPhone

**Justificativa:** Funcionalidades muito solicitadas pelos usuários

---

### Sprint 2: Contatos e Presença 🟡 (IMPORTANTE)
**Tempo Estimado:** 2-3 dias

1. CheckUser
2. GetUser
3. SendPresence
4. ChatPresence
5. GetAvatar
6. GetContacts

**Justificativa:** Melhorar experiência do usuário

---

## 📊 Estatísticas por Categoria

| Categoria | Implementados | Faltando | Total | % |
|-----------|---------------|----------|-------|---|
| **Sessões** | 8 | 1 | 9 | 89% |
| **Mensagens** | 14 | 4 | 18 | 78% |
| **Grupos** | 13 | 0 | 13 | 100% |
| **Contatos** | 0 | 6 | 6 | 0% |
| **TOTAL** | **35** | **11** | **46** | **76%** |

### Removidos (Não Aplicáveis)
| Categoria | Quantidade | Motivo |
|-----------|------------|--------|
| **Webhooks** | 4 | Sistema de eventos próprio |
| **Download** | 4 | Não necessário |
| **Admin** | 5 | Arquitetura diferente |
| **Config S3/Proxy** | 7 | Não aplicável |
| **TOTAL REMOVIDO** | **20** | - |

---

## 🎯 Recomendações

### Implementar Imediatamente (Sprint 1):
1. ✅ **PairPhone** - Autenticação sem QR Code
2. ✅ **DeleteMessage** - Deletar mensagens
3. ✅ **EditMessage** - Editar mensagens
4. ✅ **MarkRead** - Marcar como lida

### Implementar em Breve (Sprint 2):
5. ✅ **CheckUser** - Verificar número
6. ✅ **GetAvatar** - Foto de perfil
7. ✅ **SendPresence** - Presença online/typing
8. ✅ **GetContacts** - Listar contatos
9. ✅ **GetUser** - Info do usuário
10. ✅ **ChatPresence** - Presença em chat

### Não Implementar (Removidos):
- ❌ **Webhooks** - Sistema de eventos próprio
- ❌ **Download de Mídia** - Não necessário
- ❌ **Admin Multi-tenant** - Arquitetura diferente
- ❌ **S3/Proxy Config** - Não aplicável

---

## 📝 Próximos Passos

1. ✅ Grupos implementados e testados (100%)
2. ⏳ Implementar Mensagens Avançadas (Sprint 1)
3. ⏳ Implementar Contatos e Presença (Sprint 2)
4. ⏳ Avaliar sistema de eventos (WebSocket/SSE)

---

## 🔗 Referências

- `docs/referencia-handlers.bak` - Handlers de referência
- `docs/referencia-main.bak` - Lógica de eventos
- `docs/GRUPOS_DISPONIVEIS.md` - Análise de grupos
- `docs/TESTES_GRUPOS_RESULTADOS.md` - Testes realizados

