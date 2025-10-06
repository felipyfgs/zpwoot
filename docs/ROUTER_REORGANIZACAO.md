# Análise e Reorganização do Router

## 📊 Estrutura Atual

### Arquivo: `internal/adapters/http/router/routes.go` (117 linhas)

```
NewRouter()
├── setupPublicRoutes()
│   ├── GET  /
│   ├── GET  /health
│   └── GET  /swagger/*
│
└── setupAPIRoutes()
    └── setupSessionRoutes()
        ├── /sessions
        │   ├── POST   /create
        │   ├── GET    /list
        │   ├── GET    /{sessionId}/info
        │   ├── DELETE /{sessionId}/delete
        │   ├── POST   /{sessionId}/connect
        │   ├── POST   /{sessionId}/disconnect
        │   ├── POST   /{sessionId}/logout
        │   └── GET    /{sessionId}/qr
        │
        ├── setupMessageRoutes()
        │   └── /sessions/{sessionId}/send/message
        │       ├── POST /text
        │       ├── POST /image
        │       ├── POST /audio
        │       ├── POST /video
        │       ├── POST /document
        │       ├── POST /sticker
        │       ├── POST /location
        │       ├── POST /contact
        │       ├── POST /contacts
        │       ├── POST /reaction
        │       ├── POST /template
        │       ├── POST /buttons
        │       ├── POST /list
        │       └── POST /poll
        │
        └── setupGroupRoutes()
            └── /sessions/{sessionId}/groups
                ├── GET    /
                ├── GET    /info
                ├── POST   /invite-info
                ├── GET    /invite-link
                ├── POST   /join
                ├── POST   /create
                ├── POST   /leave
                ├── POST   /participants
                ├── POST   /name
                ├── POST   /topic
                ├── POST   /settings/locked
                ├── POST   /settings/announce
                ├── POST   /settings/disappearing
                ├── POST   /photo
                └── DELETE /photo
```

---

## 🎯 Problemas Identificados

1. **Hierarquia confusa**: `setupSessionRoutes` chama `setupMessageRoutes` e `setupGroupRoutes`
2. **Nomenclatura inconsistente**: Sessões, mensagens e grupos estão em níveis diferentes
3. **Falta de organização lógica**: Rotas de sessão misturadas com rotas de recursos
4. **Difícil manutenção**: Adicionar novos recursos requer modificar múltiplas funções

---

## ✅ Proposta de Reorganização

### Estrutura Proposta (Mais Compacta e Organizada)

```go
NewRouter()
├── setupPublicRoutes()      // Rotas públicas (sem auth)
└── setupProtectedRoutes()   // Rotas protegidas (com auth)
    ├── setupSessionRoutes()
    ├── setupMessageRoutes()
    └── setupGroupRoutes()
```

### Vantagens:
- ✅ Hierarquia clara e plana
- ✅ Cada recurso em sua própria função
- ✅ Fácil adicionar novos recursos
- ✅ Separação clara entre público e protegido
- ✅ Menos linhas de código
- ✅ Mais legível

---

## 📝 Implementação Proposta

```go
package router

import (
	"net/http"

	"zpwoot/internal/adapters/http/handlers"
	"zpwoot/internal/adapters/http/middleware"
	"zpwoot/internal/container"

	_ "zpwoot/docs/swagger"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter cria e configura o router principal
func NewRouter(c *container.Container) http.Handler {
	r := chi.NewRouter()

	// Middleware global
	middleware.SetupMiddleware(r)

	// Criar handlers
	h := handlers.NewHandlers(
		c.GetDatabase(),
		c.GetLogger(),
		c.GetConfig(),
		c.GetSessionUseCases(),
		c.GetMessageUseCases(),
		c.GetWhatsAppClient(),
	)

	// Configurar rotas
	setupPublicRoutes(r, h)
	setupProtectedRoutes(r, c, h)

	return r
}

// setupPublicRoutes configura rotas públicas (sem autenticação)
func setupPublicRoutes(r *chi.Mux, h *handlers.Handlers) {
	r.Get("/", h.Health.Info)
	r.Get("/health", h.Health.Health)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}

// setupProtectedRoutes configura rotas protegidas (com autenticação)
func setupProtectedRoutes(r *chi.Mux, c *container.Container, h *handlers.Handlers) {
	r.Group(func(r chi.Router) {
		middleware.SetupAuthMiddleware(r, c.GetConfig())

		setupSessionRoutes(r, h)
		setupMessageRoutes(r, h)
		setupGroupRoutes(r, h)
	})
}

// setupSessionRoutes configura rotas de gerenciamento de sessões
func setupSessionRoutes(r chi.Router, h *handlers.Handlers) {
	r.Route("/sessions", func(r chi.Router) {
		// CRUD de sessões
		r.Post("/create", h.Session.Create)
		r.Get("/list", h.Session.List)
		r.Get("/{sessionId}/info", h.Session.Get)
		r.Delete("/{sessionId}/delete", h.Session.Delete)

		// Controle de conexão
		r.Post("/{sessionId}/connect", h.Session.Connect)
		r.Post("/{sessionId}/disconnect", h.Session.Disconnect)
		r.Post("/{sessionId}/logout", h.Session.Logout)
		r.Get("/{sessionId}/qr", h.Session.QRCode)
	})
}

// setupMessageRoutes configura rotas de envio de mensagens
func setupMessageRoutes(r chi.Router, h *handlers.Handlers) {
	r.Route("/sessions/{sessionId}/send/message", func(r chi.Router) {
		// Mensagens de texto
		r.Post("/text", h.Message.SendText)

		// Mensagens de mídia
		r.Post("/image", h.Message.SendImage)
		r.Post("/audio", h.Message.SendAudio)
		r.Post("/video", h.Message.SendVideo)
		r.Post("/document", h.Message.SendDocument)
		r.Post("/sticker", h.Message.SendSticker)

		// Mensagens especiais
		r.Post("/location", h.Message.SendLocation)
		r.Post("/contact", h.Message.SendContact)
		r.Post("/contacts", h.Message.SendContactsArray)
		r.Post("/reaction", h.Message.SendReaction)

		// Mensagens interativas
		r.Post("/buttons", h.Message.SendButtons)
		r.Post("/list", h.Message.SendList)
		r.Post("/poll", h.Message.SendPoll)

		// Templates
		r.Post("/template", h.Message.SendTemplate)
	})
}

// setupGroupRoutes configura rotas de gerenciamento de grupos
func setupGroupRoutes(r chi.Router, h *handlers.Handlers) {
	r.Route("/sessions/{sessionId}/groups", func(r chi.Router) {
		// Informações
		r.Get("/", h.Group.ListGroups)
		r.Get("/info", h.Group.GetGroupInfo)
		r.Post("/invite-info", h.Group.GetGroupInviteInfo)

		// Convites
		r.Get("/invite-link", h.Group.GetGroupInviteLink)
		r.Post("/join", h.Group.JoinGroup)

		// Gerenciamento básico
		r.Post("/create", h.Group.CreateGroup)
		r.Post("/leave", h.Group.LeaveGroup)
		r.Post("/participants", h.Group.UpdateGroupParticipants)

		// Configurações do grupo
		r.Post("/name", h.Group.SetGroupName)
		r.Post("/topic", h.Group.SetGroupTopic)

		// Configurações avançadas
		r.Route("/settings", func(r chi.Router) {
			r.Post("/locked", h.Group.SetGroupLocked)
			r.Post("/announce", h.Group.SetGroupAnnounce)
			r.Post("/disappearing", h.Group.SetDisappearingTimer)
		})

		// Mídia
		r.Post("/photo", h.Group.SetGroupPhoto)
		r.Delete("/photo", h.Group.RemoveGroupPhoto)
	})
}
```

---

## 📊 Comparação

### Antes (117 linhas)
```
setupPublicRoutes()
setupAPIRoutes()
  └── setupSessionRoutes()
      ├── setupMessageRoutes()
      └── setupGroupRoutes()
```

### Depois (120 linhas, mas mais organizado)
```
setupPublicRoutes()
setupProtectedRoutes()
  ├── setupSessionRoutes()
  ├── setupMessageRoutes()
  └── setupGroupRoutes()
```

---

## 🎯 Benefícios

1. **Hierarquia Clara**: Cada recurso tem sua própria função independente
2. **Fácil Manutenção**: Adicionar novos recursos não afeta os existentes
3. **Melhor Legibilidade**: Comentários organizados por categoria
4. **Escalável**: Fácil adicionar novos grupos de rotas (contacts, status, etc.)
5. **Consistente**: Todas as rotas seguem o mesmo padrão

---

## 🔄 Próximos Passos

1. ✅ Aplicar reorganização no `routes.go`
2. ✅ Testar compilação
3. ✅ Verificar se todas as rotas funcionam
4. ✅ Atualizar documentação

---

## 📝 Notas Adicionais

### Padrão de Nomenclatura
- **Recursos**: Plural (`/sessions`, `/groups`)
- **Ações**: Verbos claros (`/create`, `/join`, `/leave`)
- **Sub-recursos**: Hierarquia lógica (`/settings/locked`)

### Organização de Handlers
- **Informações**: GET requests
- **Ações**: POST requests
- **Remoção**: DELETE requests
- **Atualização**: POST/PUT requests

### Agrupamento Lógico
- **CRUD**: Create, Read, Update, Delete
- **Controle**: Connect, Disconnect, Logout
- **Configurações**: Settings agrupadas em sub-rota
- **Mídia**: Photo, Video, etc.

