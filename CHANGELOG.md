# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

#### ViewOnce Refactoring - Breaking Change ⚠️

**ViewOnce agora é um parâmetro, não uma rota separada!**

- **ANTES**: Endpoint separado `/sessions/{sessionId}/send/message/viewonce`
- **AGORA**: Parâmetro `viewOnce: true` nos endpoints de mídia existentes

**Motivação:**
- ViewOnce é uma propriedade da mensagem, não um tipo diferente de mensagem
- Elimina duplicação de código
- Permite combinar ViewOnce com outros recursos (como contextInfo para respostas)
- Suporte a áudio ViewOnce agora disponível

**Migração:**

```diff
- POST /sessions/{sessionId}/send/message/viewonce
+ POST /sessions/{sessionId}/send/message/image

{
  "phone": "5511999999999",
  "file": "https://example.com/image.jpg",
- "caption": "ViewOnce message"
+ "caption": "ViewOnce message",
+ "viewOnce": true
}
```

**Compatibilidade:**
- O endpoint `/viewonce` ainda funciona mas está **DEPRECATED**
- Será removido em versões futuras (v2.0.0)
- Recomendamos migrar para a nova abordagem

**Documentação:**
- Ver [docs/VIEWONCE_MIGRATION.md](docs/VIEWONCE_MIGRATION.md) para guia completo
- Ver [examples/viewonce_examples.sh](examples/viewonce_examples.sh) para exemplos práticos

### Added

- ✨ Campo `viewOnce` adicionado aos DTOs:
  - `SendImageMessageRequest`
  - `SendVideoMessageRequest`
  - `SendAudioMessageRequest`
- ✨ Campo `ViewOnce` adicionado à estrutura `output.MediaData`
- ✨ Suporte a áudio ViewOnce (antes não era possível)
- ✨ Possibilidade de combinar ViewOnce com contextInfo (respostas/citações)
- 📚 Documentação de migração ViewOnce
- 📚 Exemplos de uso do ViewOnce

### Deprecated

- ⚠️ Endpoint `/sessions/{sessionId}/send/message/viewonce` marcado como DEPRECATED
  - Ainda funciona para compatibilidade retroativa
  - Será removido na versão 2.0.0
  - Use os endpoints de mídia com `viewOnce: true` em vez disso

### Technical Changes

- 🔧 Refatoração do `SendMediaMessage` no waclient para suportar ViewOnce
- 🔧 Atualização dos handlers de mensagem para processar o parâmetro ViewOnce
- 🔧 Remoção da rota `/viewonce` do router (endpoint ainda disponível via handler)
- 📝 Atualização da documentação Swagger

## [1.0.0] - 2025-10-06

### Added

- 🎉 Versão inicial do zpwoot
- ✅ Gerenciamento de sessões WhatsApp
- ✅ Envio de mensagens (texto, imagem, vídeo, áudio, documento, etc.)
- ✅ Integração com Chatwoot
- ✅ Suporte a PostgreSQL
- ✅ Migrações automáticas de banco de dados
- ✅ Documentação Swagger/OpenAPI
- ✅ Docker e Docker Compose
- ✅ Hot reload com Air
- ✅ Logs estruturados com Zerolog
- ✅ Autenticação via API Key
- ✅ Health checks

### Features

#### Session Management
- Criar sessões
- Listar sessões
- Obter informações da sessão
- Deletar sessões
- Conectar/Desconectar
- Logout
- QR Code para autenticação

#### Message Sending
- Mensagens de texto
- Imagens (com caption)
- Vídeos (com caption)
- Áudio/Voice notes
- Documentos
- Stickers
- Localização
- Contatos
- Reações
- Enquetes (Polls)
- Botões
- Listas
- Templates

#### Media Support
- Base64
- URLs
- Caminhos de arquivo
- Auto-detecção de MIME type

#### Advanced Features
- ContextInfo (respostas/citações)
- Webhooks globais
- Integração Chatwoot
- Processamento de mídia

---

## Migration Guides

### ViewOnce Migration (v1.x → v2.x)

Para migrar do endpoint deprecated `/viewonce` para a nova abordagem:

1. **Identifique o tipo de mídia** que você está enviando (imagem, vídeo ou áudio)
2. **Use o endpoint correspondente** (`/image`, `/video`, ou `/audio`)
3. **Adicione o parâmetro** `viewOnce: true` ao payload
4. **Teste** a nova implementação
5. **Remova** as chamadas ao endpoint `/viewonce`

**Exemplo completo:**

```bash
# ANTES (deprecated)
curl -X POST "http://localhost:8080/sessions/my-session/send/message/viewonce" \
  -H "Authorization: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "file": "https://example.com/image.jpg",
    "caption": "ViewOnce message"
  }'

# DEPOIS (recomendado)
curl -X POST "http://localhost:8080/sessions/my-session/send/message/image" \
  -H "Authorization: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "file": "https://example.com/image.jpg",
    "caption": "ViewOnce message",
    "viewOnce": true
  }'
```

---

## Support

Para questões, bugs ou sugestões:
- Abra uma issue no GitHub
- Consulte a documentação em `/docs`
- Veja os exemplos em `/examples`

---

[Unreleased]: https://github.com/your-org/zpwoot/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/your-org/zpwoot/releases/tag/v1.0.0

