package dto

// CheckUserRequest representa a requisição de verificação de usuários
type CheckUserRequest struct {
	Phones []string `json:"phones" validate:"required,min=1" example:"5511999999999,5511888888888"`
} //@name CheckUserRequest

// WhatsAppUserInfo representa informações de um usuário
type WhatsAppUserInfo struct {
	Query        string `json:"query" example:"5511999999999"`
	IsInWhatsApp bool   `json:"isInWhatsApp" example:"true"`
	JID          string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	VerifiedName string `json:"verifiedName,omitempty" example:"John Doe Business"`
} //@name WhatsAppUserInfo

// CheckUserResponse representa a resposta de verificação de usuários
type CheckUserResponse struct {
	Users []UserInfo `json:"users"`
} //@name CheckUserResponse

// GetUserRequest representa a requisição de informações do usuário
type GetUserRequest struct {
	Phone string `json:"phone" validate:"required" example:"5511999999999"`
} //@name GetUserRequest

// GetUserResponse representa a resposta de informações do usuário
type GetUserResponse struct {
	JID          string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	VerifiedName string `json:"verifiedName,omitempty" example:"John Doe Business"`
	Status       string `json:"status,omitempty" example:"Hey there! I am using WhatsApp."`
	PictureID    string `json:"pictureId,omitempty" example:"1234567890"`
} //@name GetUserResponse

// GetAvatarRequest representa a requisição de foto de perfil
type GetAvatarRequest struct {
	Phone   string `json:"phone" validate:"required" example:"5511999999999"`
	Preview bool   `json:"preview,omitempty" example:"false" description:"Get preview (low quality) instead of full image"`
} //@name GetAvatarRequest

// GetAvatarResponse representa a resposta de foto de perfil
type GetAvatarResponse struct {
	URL       string `json:"url,omitempty" example:"https://pps.whatsapp.net/v/..."`
	ID        string `json:"id,omitempty" example:"1234567890"`
	Type      string `json:"type,omitempty" example:"image"`
	DirectURL string `json:"directUrl,omitempty" example:"https://..."`
} //@name GetAvatarResponse

// GetContactsResponse representa a resposta de listagem de contatos
type GetContactsResponse struct {
	Contacts []ContactDetail `json:"contacts"`
} //@name GetContactsResponse

// ContactDetail representa detalhes de um contato
type ContactDetail struct {
	JID          string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	Name         string `json:"name,omitempty" example:"John Doe"`
	Notify       string `json:"notify,omitempty" example:"John"`
	VerifiedName string `json:"verifiedName,omitempty" example:"John Doe Business"`
	BusinessName string `json:"businessName,omitempty" example:"John's Store"`
} //@name ContactDetail

// SendPresenceRequest representa a requisição de envio de presença
type SendPresenceRequest struct {
	Presence string `json:"presence" validate:"required,oneof=available unavailable composing recording paused" example:"available" description:"Presence type: available, unavailable, composing, recording, paused"`
} //@name SendPresenceRequest

// SendPresenceResponse representa a resposta de envio de presença
type SendPresenceResponse struct {
	Success bool `json:"success" example:"true"`
} //@name SendPresenceResponse

// ChatPresenceRequest representa a requisição de presença em chat
type ChatPresenceRequest struct {
	Phone    string `json:"phone" validate:"required" example:"5511999999999"`
	Presence string `json:"presence" validate:"required,oneof=composing paused" example:"composing" description:"Chat presence: composing (typing), paused (stopped typing)"`
	Media    string `json:"media,omitempty" example:"text" description:"Media type: text, audio"`
} //@name ChatPresenceRequest

// ChatPresenceResponse representa a resposta de presença em chat
type ChatPresenceResponse struct {
	Success bool `json:"success" example:"true"`
} //@name ChatPresenceResponse

