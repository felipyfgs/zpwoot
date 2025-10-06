# Rotas de Grupos Disponíveis (Referência)

## 📋 Resumo

Identificadas **14 rotas de grupos** no arquivo `docs/referencia-handlers.bak`.

---

## 🔍 Rotas Identificadas

### 1. **ListGroups** - Listar Grupos
**Linha:** 3222

**Método:** GET  
**Endpoint:** `/groups/list` (presumido)

**Funcionalidade:**
- Lista todos os grupos que o usuário participa
- Retorna array de `GroupInfo`

**Request:** Nenhum body necessário

**Response:**
```json
{
  "Groups": [
    {
      "JID": "123456789@g.us",
      "Name": "Meu Grupo",
      "Participants": [...],
      ...
    }
  ]
}
```

**Implementação:**
```go
resp, err := client.GetJoinedGroups(ctx)
```

---

### 2. **GetGroupInfo** - Informações do Grupo
**Linha:** 3263

**Método:** GET  
**Endpoint:** `/groups/info` (presumido)

**Funcionalidade:**
- Obtém informações detalhadas de um grupo específico
- Nome, participantes, admins, descrição, etc.

**Request:** Query parameter
```
?groupJID=123456789@g.us
```

**Response:**
```json
{
  "JID": "123456789@g.us",
  "Name": "Meu Grupo",
  "Topic": "Descrição do grupo",
  "Participants": [...],
  "GroupCreated": 1234567890,
  ...
}
```

**Implementação:**
```go
resp, err := client.GetGroupInfo(groupJID)
```

---

### 3. **GetGroupInviteLink** - Link de Convite
**Linha:** 3313

**Método:** GET  
**Endpoint:** `/groups/invite-link` (presumido)

**Funcionalidade:**
- Obtém link de convite do grupo
- Pode resetar o link (gerar novo)

**Request:** Query parameters
```
?groupJID=123456789@g.us&reset=false
```

**Response:**
```json
{
  "InviteLink": "https://chat.whatsapp.com/ABC123DEF456"
}
```

**Implementação:**
```go
resp, err := client.GetGroupInviteLink(groupJID, reset)
```

---

### 4. **GroupJoin** - Entrar no Grupo via Link
**Linha:** 3377

**Método:** POST  
**Endpoint:** `/groups/join` (presumido)

**Funcionalidade:**
- Entra em um grupo usando código de convite

**Request:**
```json
{
  "code": "ABC123DEF456"
}
```

**Response:**
```json
{
  "Details": "Group joined successfully"
}
```

**Implementação:**
```go
_, err := client.JoinGroupWithLink(code)
```

---

### 5. **CreateGroup** - Criar Grupo
**Linha:** 3428

**Método:** POST  
**Endpoint:** `/groups/create` (presumido)

**Funcionalidade:**
- Cria um novo grupo
- Adiciona participantes iniciais

**Request:**
```json
{
  "name": "Meu Novo Grupo",
  "participants": [
    "5511999999999",
    "5511888888888"
  ]
}
```

**Response:**
```json
{
  "JID": "123456789@g.us",
  "Name": "Meu Novo Grupo",
  "Participants": [...],
  ...
}
```

**Implementação:**
```go
req := whatsmeow.ReqCreateGroup{
    Name: name,
    Participants: participantJIDs,
}
groupInfo, err := client.CreateGroup(ctx, req)
```

---

### 6. **SetGroupLocked** - Bloquear Configurações
**Linha:** 3500

**Método:** POST  
**Endpoint:** `/groups/settings/locked` (presumido)

**Funcionalidade:**
- Bloqueia/desbloqueia configurações do grupo
- Apenas admins podem editar quando bloqueado

**Request:**
```json
{
  "groupjid": "123456789@g.us",
  "locked": true
}
```

**Response:**
```json
{
  "Details": "Group Locked setting updated successfully"
}
```

**Implementação:**
```go
err := client.SetGroupLocked(groupJID, locked)
```

---

### 7. **SetDisappearingTimer** - Mensagens Temporárias
**Linha:** 3553

**Método:** POST  
**Endpoint:** `/groups/settings/disappearing` (presumido)

**Funcionalidade:**
- Configura timer de mensagens temporárias (ephemeral)
- Opções: 24h, 7d, 90d, off

**Request:**
```json
{
  "groupjid": "123456789@g.us",
  "duration": "7d"
}
```

**Durations válidas:**
- `"24h"` - 24 horas
- `"7d"` - 7 dias
- `"90d"` - 90 dias
- `"off"` - Desativar

**Response:**
```json
{
  "Details": "Disappearing timer set successfully"
}
```

**Implementação:**
```go
err := client.SetDisappearingTimer(groupJID, duration, time.Now())
```

---

### 8. **RemoveGroupPhoto** - Remover Foto do Grupo
**Linha:** 3626

**Método:** POST  
**Endpoint:** `/groups/photo/remove` (presumido)

**Funcionalidade:**
- Remove a foto do grupo

**Request:**
```json
{
  "groupjid": "123456789@g.us"
}
```

**Response:**
```json
{
  "Details": "Group Photo removed successfully"
}
```

**Implementação:**
```go
_, err := client.SetGroupPhoto(groupJID, nil)
```

---

### 9. **UpdateGroupParticipants** - Gerenciar Participantes
**Linha:** 3678

**Método:** POST  
**Endpoint:** `/groups/participants` (presumido)

**Funcionalidade:**
- Adicionar participantes
- Remover participantes
- Promover a admin
- Rebaixar de admin

**Request:**
```json
{
  "GroupJID": "123456789@g.us",
  "Phone": ["5511999999999", "5511888888888"],
  "Action": "add"
}
```

**Actions válidas:**
- `"add"` - Adicionar participantes
- `"remove"` - Remover participantes
- `"promote"` - Promover a admin
- `"demote"` - Rebaixar de admin

**Response:**
```json
{
  "Details": "Group Participants updated successfully"
}
```

**Implementação:**
```go
_, err := client.UpdateGroupParticipants(groupJID, participantJIDs, action)
```

---

### 10. **GetGroupInviteInfo** - Info do Convite
**Linha:** 3769

**Método:** POST  
**Endpoint:** `/groups/invite-info` (presumido)

**Funcionalidade:**
- Obtém informações de um grupo via código de convite
- Sem entrar no grupo

**Request:**
```json
{
  "code": "ABC123DEF456"
}
```

**Response:**
```json
{
  "JID": "123456789@g.us",
  "Name": "Nome do Grupo",
  "Size": 50,
  "Description": "Descrição do grupo",
  ...
}
```

**Implementação:**
```go
groupInfo, err := client.GetGroupInfoFromLink(code)
```

---

### 11. **SetGroupPhoto** - Definir Foto do Grupo
**Linha:** 3819

**Método:** POST  
**Endpoint:** `/groups/photo` (presumido)

**Funcionalidade:**
- Define foto do grupo
- Aceita apenas JPEG em Base64

**Request:**
```json
{
  "GroupJID": "123456789@g.us",
  "Image": "data:image/jpeg;base64,/9j/4AAQSkZJRg..."
}
```

**Validações:**
- Deve ser formato JPEG
- Deve começar com `data:image/`
- WhatsApp aceita apenas JPEG para fotos de grupo

**Response:**
```json
{
  "Details": "Group Photo set successfully",
  "PictureID": "abc123"
}
```

**Implementação:**
```go
pictureID, err := client.SetGroupPhoto(groupJID, imageData)
```

---

### 12. **SetGroupName** - Alterar Nome do Grupo
**Linha:** 3905

**Método:** POST  
**Endpoint:** `/groups/name` (presumido)

**Funcionalidade:**
- Altera o nome do grupo

**Request:**
```json
{
  "GroupJID": "123456789@g.us",
  "Name": "Novo Nome do Grupo"
}
```

**Response:**
```json
{
  "Details": "Group Name set successfully"
}
```

**Implementação:**
```go
err := client.SetGroupName(groupJID, name)
```

---

### 13. **SetGroupTopic** - Alterar Descrição
**Linha:** 3963

**Método:** POST  
**Endpoint:** `/groups/topic` (presumido)

**Funcionalidade:**
- Altera a descrição (topic) do grupo

**Request:**
```json
{
  "GroupJID": "123456789@g.us",
  "Topic": "Nova descrição do grupo"
}
```

**Response:**
```json
{
  "Details": "Group Topic set successfully"
}
```

**Implementação:**
```go
err := client.SetGroupTopic(groupJID, "", "", topic)
```

---

### 14. **GroupLeave** - Sair do Grupo
**Linha:** 4021

**Método:** POST  
**Endpoint:** `/groups/leave` (presumido)

**Funcionalidade:**
- Sai de um grupo

**Request:**
```json
{
  "GroupJID": "123456789@g.us"
}
```

**Response:**
```json
{
  "Details": "Group left successfully"
}
```

**Implementação:**
```go
err := client.LeaveGroup(groupJID)
```

---

### 15. **SetGroupAnnounce** - Modo Anúncio
**Linha:** 4073

**Método:** POST  
**Endpoint:** `/groups/settings/announce` (presumido)

**Funcionalidade:**
- Ativa/desativa modo anúncio
- Quando ativo, apenas admins podem enviar mensagens

**Request:**
```json
{
  "GroupJID": "123456789@g.us",
  "Announce": true
}
```

**Response:**
```json
{
  "Details": "Group Announce setting updated successfully"
}
```

**Implementação:**
```go
err := client.SetGroupAnnounce(groupJID, announce)
```

---

## 📊 Resumo por Categoria

### Informações (3 rotas)
- ✅ ListGroups - Listar grupos
- ✅ GetGroupInfo - Info do grupo
- ✅ GetGroupInviteInfo - Info via convite

### Convites (2 rotas)
- ✅ GetGroupInviteLink - Obter link
- ✅ GroupJoin - Entrar via link

### Gerenciamento Básico (3 rotas)
- ✅ CreateGroup - Criar grupo
- ✅ GroupLeave - Sair do grupo
- ✅ UpdateGroupParticipants - Gerenciar participantes

### Configurações (4 rotas)
- ✅ SetGroupName - Nome
- ✅ SetGroupTopic - Descrição
- ✅ SetGroupLocked - Bloquear configs
- ✅ SetGroupAnnounce - Modo anúncio

### Mídia (2 rotas)
- ✅ SetGroupPhoto - Definir foto
- ✅ RemoveGroupPhoto - Remover foto

### Avançado (1 rota)
- ✅ SetDisappearingTimer - Mensagens temporárias

---

## ✅ Status no zpwoot

**Atualmente o zpwoot NÃO possui NENHUMA rota de grupos implementada.**

Todas as 15 rotas identificadas precisam ser implementadas seguindo a arquitetura Clean do zpwoot.

