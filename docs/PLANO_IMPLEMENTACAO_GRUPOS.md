# Plano de Implementação - Rotas de Grupos

## 📋 Visão Geral

Implementação completa de 15 rotas de gerenciamento de grupos WhatsApp seguindo a Clean Architecture do zpwoot.

**Referências:**
- `docs/referencia-handlers.bak` - Handlers de referência
- `docs/referencia-main.bak` - Lógica de eventos
- `docs/GRUPOS_DISPONIVEIS.md` - Análise das rotas

---

## 🏗️ Arquitetura

### Estrutura de Pastas

```
internal/
├── adapters/
│   ├── http/
│   │   ├── handlers/
│   │   │   ├── common.go          # Adicionar GroupHandler
│   │   │   ├── group.go           # NOVO - Handler de grupos
│   │   │   ├── message.go
│   │   │   └── session.go
│   │   └── router/
│   │       └── routes.go          # Adicionar setupGroupRoutes
│   └── waclient/
│       ├── groups.go              # NOVO - Operações de grupos
│       └── service.go             # Adicionar GroupService
├── core/
│   ├── application/
│   │   └── dto/
│   │       └── group.go           # NOVO - DTOs de grupos
│   └── ports/
│       ├── input/
│       │   └── group.go           # NOVO - Interface GroupService
│       └── output/
│           └── whatsapp.go        # Adicionar métodos de grupos
```

---

## 📝 Fase 1: DTOs (Data Transfer Objects)

### Arquivo: `internal/core/application/dto/group.go`

```go
package dto

// ListGroupsResponse - Lista de grupos
type ListGroupsResponse struct {
    Groups []GroupInfo `json:"groups"`
} //@name ListGroupsResponse

type GroupInfo struct {
    JID          string   `json:"jid"`
    Name         string   `json:"name"`
    Topic        string   `json:"topic,omitempty"`
    Participants []string `json:"participants,omitempty"`
    IsAnnounce   bool     `json:"isAnnounce"`
    IsLocked     bool     `json:"isLocked"`
    CreatedAt    int64    `json:"createdAt,omitempty"`
} //@name GroupInfo

// GetGroupInfoRequest - Obter info do grupo
type GetGroupInfoRequest struct {
    GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
} //@name GetGroupInfoRequest

// GetGroupInviteLinkRequest - Obter link de convite
type GetGroupInviteLinkRequest struct {
    GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
    Reset    bool   `json:"reset,omitempty" example:"false"`
} //@name GetGroupInviteLinkRequest

type GetGroupInviteLinkResponse struct {
    InviteLink string `json:"inviteLink" example:"https://chat.whatsapp.com/ABC123"`
} //@name GetGroupInviteLinkResponse

// JoinGroupRequest - Entrar no grupo
type JoinGroupRequest struct {
    Code string `json:"code" validate:"required" example:"ABC123DEF456"`
} //@name JoinGroupRequest

// CreateGroupRequest - Criar grupo
type CreateGroupRequest struct {
    Name         string   `json:"name" validate:"required" example:"Meu Grupo"`
    Participants []string `json:"participants" validate:"required,min=1" example:"5511999999999,5511888888888"`
} //@name CreateGroupRequest

// SetGroupLockedRequest - Bloquear configurações
type SetGroupLockedRequest struct {
    GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
    Locked   bool   `json:"locked" validate:"required" example:"true"`
} //@name SetGroupLockedRequest

// SetDisappearingTimerRequest - Mensagens temporárias
type SetDisappearingTimerRequest struct {
    GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
    Duration string `json:"duration" validate:"required,oneof=24h 7d 90d off" example:"7d"`
} //@name SetDisappearingTimerRequest

// RemoveGroupPhotoRequest - Remover foto
type RemoveGroupPhotoRequest struct {
    GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
} //@name RemoveGroupPhotoRequest

// UpdateGroupParticipantsRequest - Gerenciar participantes
type UpdateGroupParticipantsRequest struct {
    GroupJID     string   `json:"groupJid" validate:"required" example:"123456789@g.us"`
    Participants []string `json:"participants" validate:"required,min=1" example:"5511999999999"`
    Action       string   `json:"action" validate:"required,oneof=add remove promote demote" example:"add"`
} //@name UpdateGroupParticipantsRequest

// GetGroupInviteInfoRequest - Info do convite
type GetGroupInviteInfoRequest struct {
    Code string `json:"code" validate:"required" example:"ABC123DEF456"`
} //@name GetGroupInviteInfoRequest

// SetGroupPhotoRequest - Definir foto
type SetGroupPhotoRequest struct {
    GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
    Image    string `json:"image" validate:"required" example:"data:image/jpeg;base64,..."`
} //@name SetGroupPhotoRequest

type SetGroupPhotoResponse struct {
    PictureID string `json:"pictureId" example:"abc123"`
} //@name SetGroupPhotoResponse

// SetGroupNameRequest - Alterar nome
type SetGroupNameRequest struct {
    GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
    Name     string `json:"name" validate:"required" example:"Novo Nome"`
} //@name SetGroupNameRequest

// SetGroupTopicRequest - Alterar descrição
type SetGroupTopicRequest struct {
    GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
    Topic    string `json:"topic" validate:"required" example:"Nova descrição"`
} //@name SetGroupTopicRequest

// LeaveGroupRequest - Sair do grupo
type LeaveGroupRequest struct {
    GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
} //@name LeaveGroupRequest

// SetGroupAnnounceRequest - Modo anúncio
type SetGroupAnnounceRequest struct {
    GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
    Announce bool   `json:"announce" validate:"required" example:"true"`
} //@name SetGroupAnnounceRequest

// Respostas genéricas
type GroupActionResponse struct {
    Success bool   `json:"success" example:"true"`
    Message string `json:"message" example:"Operation completed successfully"`
} //@name GroupActionResponse
```

---

## 📝 Fase 2: Ports (Interfaces)

### Arquivo: `internal/core/ports/input/group.go`

```go
package input

import (
    "context"
    "zpwoot/internal/core/application/dto"
)

type GroupService interface {
    // Informações
    ListGroups(ctx context.Context, sessionID string) (*dto.ListGroupsResponse, error)
    GetGroupInfo(ctx context.Context, sessionID string, groupJID string) (*dto.GroupInfo, error)
    GetGroupInviteInfo(ctx context.Context, sessionID string, code string) (*dto.GroupInfo, error)
    
    // Convites
    GetGroupInviteLink(ctx context.Context, sessionID string, groupJID string, reset bool) (string, error)
    JoinGroup(ctx context.Context, sessionID string, code string) error
    
    // Gerenciamento
    CreateGroup(ctx context.Context, sessionID string, name string, participants []string) (*dto.GroupInfo, error)
    LeaveGroup(ctx context.Context, sessionID string, groupJID string) error
    UpdateGroupParticipants(ctx context.Context, sessionID string, groupJID string, participants []string, action string) error
    
    // Configurações
    SetGroupName(ctx context.Context, sessionID string, groupJID string, name string) error
    SetGroupTopic(ctx context.Context, sessionID string, groupJID string, topic string) error
    SetGroupLocked(ctx context.Context, sessionID string, groupJID string, locked bool) error
    SetGroupAnnounce(ctx context.Context, sessionID string, groupJID string, announce bool) error
    SetDisappearingTimer(ctx context.Context, sessionID string, groupJID string, duration string) error
    
    // Mídia
    SetGroupPhoto(ctx context.Context, sessionID string, groupJID string, imageData []byte) (string, error)
    RemoveGroupPhoto(ctx context.Context, sessionID string, groupJID string) error
}
```

---

## 📝 Fase 3: Implementação waclient

### Arquivo: `internal/adapters/waclient/groups.go`

```go
package waclient

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types"
    "zpwoot/internal/core/application/dto"
)

type GroupManager struct {
    clientManager *Manager
}

func NewGroupManager(manager *Manager) *GroupManager {
    return &GroupManager{
        clientManager: manager,
    }
}

// ListGroups - Lista grupos participantes
func (gm *GroupManager) ListGroups(ctx context.Context, sessionID string) (*dto.ListGroupsResponse, error) {
    client, err := gm.clientManager.GetClient(sessionID)
    if err != nil {
        return nil, err
    }
    
    groups, err := client.WAClient.GetJoinedGroups()
    if err != nil {
        return nil, fmt.Errorf("failed to get joined groups: %w", err)
    }
    
    response := &dto.ListGroupsResponse{
        Groups: make([]dto.GroupInfo, 0, len(groups)),
    }
    
    for _, group := range groups {
        response.Groups = append(response.Groups, dto.GroupInfo{
            JID:        group.JID.String(),
            Name:       group.Name,
            Topic:      group.Topic,
            IsAnnounce: group.IsAnnounce,
            IsLocked:   group.IsLocked,
        })
    }
    
    return response, nil
}

// GetGroupInfo - Informações do grupo
func (gm *GroupManager) GetGroupInfo(ctx context.Context, sessionID string, groupJID string) (*dto.GroupInfo, error) {
    client, err := gm.clientManager.GetClient(sessionID)
    if err != nil {
        return nil, err
    }
    
    jid, err := parseJID(groupJID)
    if err != nil {
        return nil, err
    }
    
    group, err := client.WAClient.GetGroupInfo(jid)
    if err != nil {
        return nil, fmt.Errorf("failed to get group info: %w", err)
    }
    
    participants := make([]string, 0, len(group.Participants))
    for _, p := range group.Participants {
        participants = append(participants, p.JID.String())
    }
    
    return &dto.GroupInfo{
        JID:          group.JID.String(),
        Name:         group.Name,
        Topic:        group.Topic,
        Participants: participants,
        IsAnnounce:   group.IsAnnounce,
        IsLocked:     group.IsLocked,
        CreatedAt:    group.GroupCreated.Unix(),
    }, nil
}

// Continua com outros métodos...
```

---

## 📝 Fase 4: Handler HTTP

### Arquivo: `internal/adapters/http/handlers/group.go`

Estrutura baseada em `message.go` e `session.go`:

```go
package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "zpwoot/internal/adapters/logger"
    "zpwoot/internal/core/application/dto"
    "zpwoot/internal/core/ports/input"
)

type GroupHandler struct {
    groupService input.GroupService
    logger       *logger.Logger
}

func NewGroupHandler(groupService input.GroupService, logger *logger.Logger) *GroupHandler {
    return &GroupHandler{
        groupService: groupService,
        logger:       logger,
    }
}

// ListGroups godoc
// @Summary      List groups
// @Description  List all groups the session is part of
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string  true  "Session ID"
// @Success      200  {object}  dto.ListGroupsResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups [get]
func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
    sessionID := chi.URLParam(r, "sessionId")
    if sessionID == "" {
        h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
        return
    }
    
    groups, err := h.groupService.ListGroups(r.Context(), sessionID)
    if err != nil {
        h.handleGroupError(w, err)
        return
    }
    
    h.writeJSON(w, groups)
}

// Continua com outros handlers...
```

---

## 📝 Fase 5: Rotas

### Arquivo: `internal/adapters/http/router/routes.go`

```go
func setupGroupRoutes(r chi.Router, h *handlers.Handlers) {
    r.Route("/sessions/{sessionId}/groups", func(r chi.Router) {
        // Informações
        r.Get("/", h.Group.ListGroups)
        r.Get("/info", h.Group.GetGroupInfo)
        r.Post("/invite-info", h.Group.GetGroupInviteInfo)
        
        // Convites
        r.Get("/invite-link", h.Group.GetGroupInviteLink)
        r.Post("/join", h.Group.JoinGroup)
        
        // Gerenciamento
        r.Post("/create", h.Group.CreateGroup)
        r.Post("/leave", h.Group.LeaveGroup)
        r.Post("/participants", h.Group.UpdateGroupParticipants)
        
        // Configurações
        r.Post("/name", h.Group.SetGroupName)
        r.Post("/topic", h.Group.SetGroupTopic)
        r.Post("/settings/locked", h.Group.SetGroupLocked)
        r.Post("/settings/announce", h.Group.SetGroupAnnounce)
        r.Post("/settings/disappearing", h.Group.SetDisappearingTimer)
        
        // Mídia
        r.Post("/photo", h.Group.SetGroupPhoto)
        r.Delete("/photo", h.Group.RemoveGroupPhoto)
    })
}
```

Adicionar em `setupSessionRoutes`:
```go
func setupSessionRoutes(r chi.Router, h *handlers.Handlers) {
    // ... código existente ...
    
    setupMessageRoutes(r, h)
    setupGroupRoutes(r, h)  // NOVO
}
```

---

## 🔄 Ordem de Implementação

### Sprint 1: Fundação (2-3 dias)
1. ✅ Criar DTOs (`dto/group.go`)
2. ✅ Criar interfaces (`ports/input/group.go`)
3. ✅ Criar estrutura básica do GroupManager (`waclient/groups.go`)
4. ✅ Criar GroupHandler básico (`handlers/group.go`)
5. ✅ Adicionar rotas (`router/routes.go`)
6. ✅ Atualizar `common.go` para incluir GroupHandler

### Sprint 2: Informações (1-2 dias)
7. ✅ Implementar ListGroups
8. ✅ Implementar GetGroupInfo
9. ✅ Implementar GetGroupInviteInfo
10. ✅ Testes unitários

### Sprint 3: Convites (1 dia)
11. ✅ Implementar GetGroupInviteLink
12. ✅ Implementar JoinGroup
13. ✅ Testes unitários

### Sprint 4: Gerenciamento (2 dias)
14. ✅ Implementar CreateGroup
15. ✅ Implementar LeaveGroup
16. ✅ Implementar UpdateGroupParticipants
17. ✅ Testes unitários

### Sprint 5: Configurações (2 dias)
18. ✅ Implementar SetGroupName
19. ✅ Implementar SetGroupTopic
20. ✅ Implementar SetGroupLocked
21. ✅ Implementar SetGroupAnnounce
22. ✅ Implementar SetDisappearingTimer
23. ✅ Testes unitários

### Sprint 6: Mídia (1 dia)
24. ✅ Implementar SetGroupPhoto
25. ✅ Implementar RemoveGroupPhoto
26. ✅ Testes unitários

### Sprint 7: Documentação e Testes (1 dia)
27. ✅ Atualizar Swagger
28. ✅ Atualizar docs/API.md
29. ✅ Testes de integração
30. ✅ Validação final

**Total estimado: 10-12 dias**

---

## 📋 Checklist de Implementação

- [ ] Criar `internal/core/application/dto/group.go`
- [ ] Criar `internal/core/ports/input/group.go`
- [ ] Criar `internal/adapters/waclient/groups.go`
- [ ] Criar `internal/adapters/http/handlers/group.go`
- [ ] Atualizar `internal/adapters/http/handlers/common.go`
- [ ] Atualizar `internal/adapters/http/router/routes.go`
- [ ] Implementar todos os 15 métodos
- [ ] Adicionar validações
- [ ] Adicionar tratamento de erros
- [ ] Adicionar logs
- [ ] Documentar Swagger
- [ ] Atualizar docs/API.md
- [ ] Testes unitários
- [ ] Testes de integração
- [ ] Code review

---

## 🎯 Próximos Passos

1. Revisar e aprovar este plano
2. Criar branch `feature/groups`
3. Iniciar Sprint 1
4. Implementar incrementalmente
5. Testar continuamente
6. Documentar progressivamente

