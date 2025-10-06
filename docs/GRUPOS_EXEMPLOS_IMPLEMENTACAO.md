# Exemplos de Implementação - Grupos

## 📋 Baseado em

- `docs/referencia-handlers.bak` (linhas 3222-4100)
- Padrão zpwoot (Clean Architecture)

---

## 1️⃣ Exemplo Completo: ListGroups

### DTO (`dto/group.go`)
```go
type ListGroupsResponse struct {
    Groups []GroupInfo `json:"groups"`
} //@name ListGroupsResponse

type GroupInfo struct {
    JID          string   `json:"jid" example:"123456789@g.us"`
    Name         string   `json:"name" example:"Meu Grupo"`
    Topic        string   `json:"topic,omitempty" example:"Descrição do grupo"`
    Participants []string `json:"participants,omitempty"`
    IsAnnounce   bool     `json:"isAnnounce" example:"false"`
    IsLocked     bool     `json:"isLocked" example:"false"`
    CreatedAt    int64    `json:"createdAt,omitempty" example:"1696570882"`
} //@name GroupInfo
```

### Interface (`ports/input/group.go`)
```go
type GroupService interface {
    ListGroups(ctx context.Context, sessionID string) (*dto.ListGroupsResponse, error)
}
```

### Implementação (`waclient/groups.go`)
```go
// Baseado em referencia-handlers.bak linha 3222
func (gm *GroupManager) ListGroups(ctx context.Context, sessionID string) (*dto.ListGroupsResponse, error) {
    client, err := gm.clientManager.GetClient(sessionID)
    if err != nil {
        return nil, fmt.Errorf("session not found: %w", err)
    }
    
    // Referência: client.GetJoinedGroups(r.Context())
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
```

### Handler (`handlers/group.go`)
```go
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
        h.logger.Error().
            Err(err).
            Str("session_id", sessionID).
            Msg("Failed to list groups")
        h.handleGroupError(w, err)
        return
    }
    
    h.logger.Info().
        Str("session_id", sessionID).
        Int("count", len(groups.Groups)).
        Msg("Groups listed successfully")
    
    h.writeJSON(w, groups)
}
```

### Rota (`router/routes.go`)
```go
r.Get("/", h.Group.ListGroups)
```

### Exemplo de Uso
```bash
curl -X GET http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/groups \
  -H "Authorization: YOUR_API_KEY"
```

---

## 2️⃣ Exemplo Completo: CreateGroup

### DTO
```go
type CreateGroupRequest struct {
    Name         string   `json:"name" validate:"required" example:"Meu Grupo"`
    Participants []string `json:"participants" validate:"required,min=1" example:"5511999999999,5511888888888"`
} //@name CreateGroupRequest
```

### Implementação (`waclient/groups.go`)
```go
// Baseado em referencia-handlers.bak linha 3428
func (gm *GroupManager) CreateGroup(ctx context.Context, sessionID string, name string, participants []string) (*dto.GroupInfo, error) {
    client, err := gm.clientManager.GetClient(sessionID)
    if err != nil {
        return nil, fmt.Errorf("session not found: %w", err)
    }
    
    if name == "" {
        return nil, errors.New("group name is required")
    }
    
    if len(participants) < 1 {
        return nil, errors.New("at least one participant is required")
    }
    
    // Parse participant JIDs
    participantJIDs := make([]types.JID, len(participants))
    for i, phone := range participants {
        jid, err := parseJID(phone)
        if err != nil {
            return nil, fmt.Errorf("invalid participant phone %s: %w", phone, err)
        }
        participantJIDs[i] = jid
    }
    
    // Referência: whatsmeow.ReqCreateGroup
    req := whatsmeow.ReqCreateGroup{
        Name:         name,
        Participants: participantJIDs,
    }
    
    groupInfo, err := client.WAClient.CreateGroup(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create group: %w", err)
    }
    
    // Converter para DTO
    participantStrings := make([]string, len(groupInfo.Participants))
    for i, p := range groupInfo.Participants {
        participantStrings[i] = p.JID.String()
    }
    
    return &dto.GroupInfo{
        JID:          groupInfo.JID.String(),
        Name:         groupInfo.Name,
        Topic:        groupInfo.Topic,
        Participants: participantStrings,
        IsAnnounce:   groupInfo.IsAnnounce,
        IsLocked:     groupInfo.IsLocked,
        CreatedAt:    groupInfo.GroupCreated.Unix(),
    }, nil
}
```

### Handler
```go
// @Summary      Create group
// @Description  Create a new WhatsApp group
// @Tags         Groups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        sessionId   path      string                    true  "Session ID"
// @Param        group       body      dto.CreateGroupRequest    true  "Group data"
// @Success      200  {object}  dto.GroupInfo
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /sessions/{sessionId}/groups/create [post]
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
    sessionID := chi.URLParam(r, "sessionId")
    if sessionID == "" {
        h.writeError(w, http.StatusBadRequest, "validation_error", "sessionId is required")
        return
    }
    
    var req dto.CreateGroupRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
        return
    }
    
    // Validação
    if req.Name == "" {
        h.writeError(w, http.StatusBadRequest, "validation_error", "name is required")
        return
    }
    
    if len(req.Participants) < 1 {
        h.writeError(w, http.StatusBadRequest, "validation_error", "at least one participant is required")
        return
    }
    
    group, err := h.groupService.CreateGroup(r.Context(), sessionID, req.Name, req.Participants)
    if err != nil {
        h.logger.Error().
            Err(err).
            Str("session_id", sessionID).
            Str("group_name", req.Name).
            Msg("Failed to create group")
        h.handleGroupError(w, err)
        return
    }
    
    h.logger.Info().
        Str("session_id", sessionID).
        Str("group_jid", group.JID).
        Str("group_name", group.Name).
        Msg("Group created successfully")
    
    h.writeJSON(w, group)
}
```

---

## 3️⃣ Exemplo: UpdateGroupParticipants

### Implementação (`waclient/groups.go`)
```go
// Baseado em referencia-handlers.bak linha 3678
func (gm *GroupManager) UpdateGroupParticipants(ctx context.Context, sessionID string, groupJID string, participants []string, action string) error {
    client, err := gm.clientManager.GetClient(sessionID)
    if err != nil {
        return fmt.Errorf("session not found: %w", err)
    }
    
    jid, err := parseJID(groupJID)
    if err != nil {
        return fmt.Errorf("invalid group JID: %w", err)
    }
    
    if len(participants) < 1 {
        return errors.New("at least one participant is required")
    }
    
    // Parse participant JIDs
    participantJIDs := make([]types.JID, len(participants))
    for i, phone := range participants {
        pjid, err := parseJID(phone)
        if err != nil {
            return fmt.Errorf("invalid participant phone %s: %w", phone, err)
        }
        participantJIDs[i] = pjid
    }
    
    // Parse action
    var participantChange whatsmeow.ParticipantChange
    switch action {
    case "add":
        participantChange = whatsmeow.ParticipantChangeAdd
    case "remove":
        participantChange = whatsmeow.ParticipantChangeRemove
    case "promote":
        participantChange = whatsmeow.ParticipantChangePromote
    case "demote":
        participantChange = whatsmeow.ParticipantChangeDemote
    default:
        return fmt.Errorf("invalid action: %s (must be add, remove, promote, or demote)", action)
    }
    
    // Referência: client.UpdateGroupParticipants
    _, err = client.WAClient.UpdateGroupParticipants(jid, participantJIDs, participantChange)
    if err != nil {
        return fmt.Errorf("failed to update group participants: %w", err)
    }
    
    return nil
}
```

---

## 4️⃣ Exemplo: SetDisappearingTimer

### Implementação (`waclient/groups.go`)
```go
// Baseado em referencia-handlers.bak linha 3553
func (gm *GroupManager) SetDisappearingTimer(ctx context.Context, sessionID string, groupJID string, duration string) error {
    client, err := gm.clientManager.GetClient(sessionID)
    if err != nil {
        return fmt.Errorf("session not found: %w", err)
    }
    
    jid, err := parseJID(groupJID)
    if err != nil {
        return fmt.Errorf("invalid group JID: %w", err)
    }
    
    // Parse duration
    var timer time.Duration
    switch duration {
    case "24h":
        timer = 24 * time.Hour
    case "7d":
        timer = 7 * 24 * time.Hour
    case "90d":
        timer = 90 * 24 * time.Hour
    case "off":
        timer = 0
    default:
        return fmt.Errorf("invalid duration: %s (must be 24h, 7d, 90d, or off)", duration)
    }
    
    // Referência: client.SetDisappearingTimer
    err = client.WAClient.SetDisappearingTimer(jid, timer)
    if err != nil {
        return fmt.Errorf("failed to set disappearing timer: %w", err)
    }
    
    return nil
}
```

---

## 5️⃣ Exemplo: SetGroupPhoto

### Implementação (`waclient/groups.go`)
```go
// Baseado em referencia-handlers.bak linha 3819
func (gm *GroupManager) SetGroupPhoto(ctx context.Context, sessionID string, groupJID string, imageData []byte) (string, error) {
    client, err := gm.clientManager.GetClient(sessionID)
    if err != nil {
        return "", fmt.Errorf("session not found: %w", err)
    }
    
    jid, err := parseJID(groupJID)
    if err != nil {
        return "", fmt.Errorf("invalid group JID: %w", err)
    }
    
    if len(imageData) == 0 {
        return "", errors.New("image data is required")
    }
    
    // Validar formato JPEG (WhatsApp requer JPEG)
    if len(imageData) < 3 || imageData[0] != 0xFF || imageData[1] != 0xD8 || imageData[2] != 0xFF {
        return "", errors.New("image must be in JPEG format")
    }
    
    // Referência: client.SetGroupPhoto
    pictureID, err := client.WAClient.SetGroupPhoto(jid, imageData)
    if err != nil {
        return "", fmt.Errorf("failed to set group photo: %w", err)
    }
    
    return pictureID, nil
}
```

---

## 🔧 Funções Auxiliares

### parseJID
```go
func parseJID(phone string) (types.JID, error) {
    // Remove caracteres não numéricos
    cleaned := strings.Map(func(r rune) rune {
        if r >= '0' && r <= '9' {
            return r
        }
        return -1
    }, phone)
    
    // Se já tem @, parsear diretamente
    if strings.Contains(phone, "@") {
        jid, err := types.ParseJID(phone)
        if err != nil {
            return types.JID{}, fmt.Errorf("invalid JID: %w", err)
        }
        return jid, nil
    }
    
    // Adicionar sufixo apropriado
    var suffix string
    if strings.HasSuffix(phone, "@g.us") {
        suffix = ""
    } else if len(cleaned) > 0 {
        suffix = "@s.whatsapp.net"
    } else {
        return types.JID{}, errors.New("invalid phone number")
    }
    
    jid, err := types.ParseJID(cleaned + suffix)
    if err != nil {
        return types.JID{}, fmt.Errorf("invalid JID: %w", err)
    }
    
    return jid, nil
}
```

### handleGroupError
```go
func (h *GroupHandler) handleGroupError(w http.ResponseWriter, err error) {
    if err == nil {
        return
    }
    
    errMsg := err.Error()
    
    switch {
    case strings.Contains(errMsg, "session not found"):
        h.writeError(w, http.StatusNotFound, "session_not_found", "Session not found")
    case strings.Contains(errMsg, "not connected"):
        h.writeError(w, http.StatusPreconditionFailed, "not_connected", "Session not connected")
    case strings.Contains(errMsg, "invalid"):
        h.writeError(w, http.StatusBadRequest, "invalid_request", errMsg)
    case strings.Contains(errMsg, "required"):
        h.writeError(w, http.StatusBadRequest, "validation_error", errMsg)
    default:
        h.writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
    }
}
```

---

## 📝 Notas de Implementação

1. **Sempre validar sessionID** antes de qualquer operação
2. **Parsear JIDs corretamente** (grupos usam `@g.us`, usuários `@s.whatsapp.net`)
3. **Validar formatos** (JPEG para fotos, durações válidas, etc.)
4. **Logar operações** para auditoria
5. **Tratar erros específicos** do whatsmeow
6. **Seguir padrão zpwoot** (Clean Architecture, DTOs, Ports)
7. **Documentar Swagger** para cada endpoint
8. **Adicionar testes unitários** para cada método

---

## 🎯 Próximos Passos

1. Implementar todos os 15 métodos seguindo estes exemplos
2. Adicionar validações robustas
3. Implementar testes unitários
4. Atualizar documentação Swagger
5. Atualizar docs/API.md

