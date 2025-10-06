# Roadmap de Implementação - zpwoot API

## 📊 Status Atual

```
█████████████████████████████░░░ 76% (35/46 endpoints)
```

**Implementados:** 35 endpoints
**Faltando:** 11 endpoints
**Total Aplicável:** 46 endpoints
**Removidos (Não Aplicáveis):** 16 endpoints

---

## ✅ Implementado (35 endpoints)

### Sessões (8/9) - 89%
```
✅ Create, List, Get, Delete
✅ Connect, Disconnect, Logout
✅ GetQR
❌ PairPhone
```

### Mensagens (14/18) - 78%
```
✅ Text, Image, Audio, Video, Document, Sticker
✅ Location, Contact, ContactsArray
✅ Buttons, List, Poll, Template
✅ Reaction
❌ Delete, Edit, MarkRead, HistorySync
```

### Grupos (13/13) - 100% ✅
```
✅ List, Create, Leave, Join
✅ GetInfo, GetInviteLink, GetInviteInfo
✅ UpdateParticipants
✅ SetName, SetTopic, SetPhoto, RemovePhoto
✅ SetLocked, SetAnnounce, SetDisappearing
```

### Contatos (0/6) - 0%
```
❌ CheckUser, GetUser, GetAvatar
❌ GetContacts, SendPresence, ChatPresence
```

---

## ❌ Faltando (11 endpoints)

### 🟡 ALTA PRIORIDADE (5 endpoints)

#### Autenticação (1)
```
❌ PairPhone - Pareamento por telefone
```

#### Mensagens Avançadas (4)
```
❌ DeleteMessage   - Deletar mensagem
❌ EditMessage     - Editar mensagem
❌ MarkRead        - Marcar como lida
❌ HistorySync     - Sincronizar histórico
```

---

### 🟡 MÉDIA PRIORIDADE (6 endpoints)

#### Contatos e Presença (6)
```
❌ CheckUser       - Verificar número
❌ GetUser         - Info do usuário
❌ GetAvatar       - Foto de perfil
❌ GetContacts     - Listar contatos
❌ SendPresence    - Online/Typing
❌ ChatPresence    - Presença em chat
```

---

## ⚪ Removidos - Não Aplicáveis (16 endpoints)

### Webhooks (4) - Sistema Diferente
```
⚪ GetWebhook, SetWebhook, UpdateWebhook, DeleteWebhook
```
**Motivo:** zpwoot implementará sistema de eventos próprio (WebSocket/SSE)

### Download de Mídia (4) - Não Necessário
```
⚪ DownloadImage, DownloadDocument, DownloadVideo, DownloadAudio
```
**Motivo:** Mídias acessíveis diretamente via whatsmeow

### Admin Multi-Tenant (5) - Arquitetura Diferente
```
⚪ ListUsers, AddUser, EditUser, DeleteUser, DeleteUserComplete
```
**Motivo:** zpwoot usa sessões, não multi-tenant

### Configurações S3/Proxy (7) - Não Aplicável
```
⚪ SetHistory, GetHistory, SetProxy
⚪ ConfigureS3, GetS3Config, TestS3Connection, DeleteS3Config
```
**Motivo:** Não necessário na arquitetura do zpwoot

---

## 🎯 Plano de Implementação

### Sprint 1: Mensagens Avançadas (2-3 dias) 🟡
**Prioridade:** ALTA

```
1. DeleteMessage    - Deletar mensagem
2. EditMessage      - Editar mensagem
3. MarkRead         - Marcar como lida
4. HistorySync      - Sincronizar histórico
5. PairPhone        - Pareamento por telefone
```

**Entregável:**
- Operações avançadas de mensagens
- Autenticação alternativa

---

### Sprint 2: Contatos e Presença (2-3 dias) 🟡
**Prioridade:** MÉDIA

```
1. CheckUser        - Verificar número
2. GetUser          - Info do usuário
3. GetAvatar        - Foto de perfil
4. GetContacts      - Listar contatos
5. SendPresence     - Online/Typing
6. ChatPresence     - Presença em chat
```

**Entregável:**
- Gerenciamento de contatos
- Sistema de presença

---

## 📈 Progresso por Sprint

```
Sprint 0 (Atual):  █████████████████████████████░░░ 76% (35/46)
Sprint 1 (+5):     ████████████████████████████████ 87% (40/46)
Sprint 2 (+6):     ████████████████████████████████ 100% (46/46) ✅
```

---

## 🎯 Metas

### Curto Prazo (1-2 semanas)
- ✅ Implementar Mensagens Avançadas
- 🎯 **Meta:** 87% de cobertura

### Médio Prazo (2-3 semanas)
- ✅ Implementar Contatos e Presença
- 🎯 **Meta:** 100% de cobertura ✅

### Longo Prazo (1 mês)
- ✅ Sistema de eventos (WebSocket/SSE)
- ✅ Otimizações e melhorias
- 🎯 **Meta:** API completa e otimizada

---

## 📊 Comparação com Referência

| Categoria | zpwoot | Aplicável | Gap |
|-----------|--------|-----------|-----|
| Sessões | 8 | 9 | -1 |
| Mensagens | 14 | 18 | -4 |
| Grupos | 13 | 13 | ✅ |
| Contatos | 0 | 6 | -6 |
| **Total Aplicável** | **35** | **46** | **-11** |

### Removidos (Não Aplicáveis)
| Categoria | Quantidade |
|-----------|------------|
| Webhooks | 4 |
| Download | 4 |
| Admin | 5 |
| Config S3/Proxy | 7 |
| **Total Removido** | **20** |

---

## 🚀 Próximos Passos

1. ✅ Grupos implementados e testados (100%)
2. ⏳ Implementar Sprint 1: Mensagens Avançadas
3. ⏳ Implementar Sprint 2: Contatos e Presença
4. ⏳ Avaliar sistema de eventos (WebSocket/SSE)

---

## 📝 Notas

- **Grupos:** ✅ 100% implementado e testado
- **Cobertura Atual:** 76% (35/46 endpoints aplicáveis)
- **Removidos:** 20 endpoints não aplicáveis ao zpwoot
- **Webhooks:** Sistema de eventos próprio será implementado
- **Multi-tenant:** Não aplicável (zpwoot usa sessões)
- **Download/S3/Proxy:** Não necessário na arquitetura atual

---

**Última Atualização:** 2025-10-06  
**Versão:** 1.0  
**Status:** Em Desenvolvimento

