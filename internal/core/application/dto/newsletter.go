package dto

// CreateNewsletterRequest - Requisição para criar newsletter
type CreateNewsletterRequest struct {
	Name        string `json:"name" binding:"required" example:"Meu Newsletter"`
	Description string `json:"description,omitempty" example:"Descrição do newsletter"`
} //@name CreateNewsletterRequest

// NewsletterInfo - Informações do newsletter
type NewsletterInfo struct {
	JID             string `json:"jid" example:"123456789@newsletter"`
	Name            string `json:"name" example:"Meu Newsletter"`
	Description     string `json:"description,omitempty" example:"Descrição do newsletter"`
	SubscriberCount int    `json:"subscriber_count" example:"150"`
	IsOwner         bool   `json:"is_owner" example:"true"`
	IsFollowing     bool   `json:"is_following" example:"true"`
	IsMuted         bool   `json:"is_muted" example:"false"`
	CreatedAt       int64  `json:"created_at,omitempty" example:"1696570882"`
} //@name NewsletterInfo

// ListNewslettersResponse - Lista de newsletters
type ListNewslettersResponse struct {
	Newsletters []NewsletterInfo `json:"newsletters"`
} //@name ListNewslettersResponse

// FollowNewsletterRequest - Requisição para seguir newsletter
type FollowNewsletterRequest struct {
	NewsletterJID string `json:"newsletter_jid,omitempty" example:"123456789@newsletter"`
	InviteCode    string `json:"invite_code,omitempty" example:"abc123def456"`
} //@name FollowNewsletterRequest

// NewsletterMessage - Mensagem do newsletter
type NewsletterMessage struct {
	ID        string `json:"id" example:"msg123"`
	ServerID  string `json:"server_id" example:"srv456"`
	Content   string `json:"content" example:"Conteúdo da mensagem"`
	Type      string `json:"type" example:"text"` // text, image, video, etc.
	Timestamp int64  `json:"timestamp" example:"1696570882"`
	ViewCount int    `json:"view_count,omitempty" example:"25"`
} //@name NewsletterMessage

// ListNewsletterMessagesResponse - Lista de mensagens do newsletter
type ListNewsletterMessagesResponse struct {
	Messages []NewsletterMessage `json:"messages"`
	HasMore  bool                `json:"has_more" example:"true"`
	Cursor   string              `json:"cursor,omitempty" example:"cursor123"`
} //@name ListNewsletterMessagesResponse

// GetNewsletterMessagesRequest - Requisição para obter mensagens do newsletter
type GetNewsletterMessagesRequest struct {
	Count  int    `json:"count,omitempty" example:"50"`
	Before string `json:"before,omitempty" example:"cursor123"`
} //@name GetNewsletterMessagesRequest

// NewsletterReactionRequest - Requisição para reagir a mensagem do newsletter
type NewsletterReactionRequest struct {
	MessageID string `json:"message_id" binding:"required" example:"msg123"`
	ServerID  string `json:"server_id" binding:"required" example:"srv456"`
	Reaction  string `json:"reaction" binding:"required" example:"👍"`
} //@name NewsletterReactionRequest

// NewsletterMuteRequest - Requisição para silenciar/dessilenciar newsletter
type NewsletterMuteRequest struct {
	Mute bool `json:"mute" example:"true"`
} //@name NewsletterMuteRequest

// NewsletterMarkViewedRequest - Requisição para marcar mensagens como visualizadas
type NewsletterMarkViewedRequest struct {
	ServerIDs []string `json:"server_ids" binding:"required" example:"srv456,srv789"`
} //@name NewsletterMarkViewedRequest

// NewsletterInfoWithInviteRequest - Requisição para obter info do newsletter via convite
type NewsletterInfoWithInviteRequest struct {
	InviteKey string `json:"invite_key" binding:"required" example:"abc123def456"`
} //@name NewsletterInfoWithInviteRequest
