package dto

// ListGroupsResponse - Lista de grupos
type ListGroupsResponse struct {
	Groups []WhatsAppGroupInfo `json:"groups"`
} //@name ListGroupsResponse

// WhatsAppGroupInfo - Informações do grupo
type WhatsAppGroupInfo struct {
	JID          string   `json:"jid" example:"123456789@g.us"`
	Name         string   `json:"name" example:"Meu Grupo"`
	Topic        string   `json:"topic,omitempty" example:"Descrição do grupo"`
	Participants []string `json:"participants,omitempty"`
	IsAnnounce   bool     `json:"isAnnounce" example:"false"`
	IsLocked     bool     `json:"isLocked" example:"false"`
	CreatedAt    int64    `json:"createdAt,omitempty" example:"1696570882"`
} //@name WhatsAppGroupInfo

// GetGroupInfoRequest - Obter informações do grupo
type GetGroupInfoRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
} //@name GetGroupInfoRequest

// GetGroupInviteLinkRequest - Obter link de convite
type GetGroupInviteLinkRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Reset    bool   `json:"reset,omitempty" example:"false"`
} //@name GetGroupInviteLinkRequest

// GetInviteLinkResponse - Resposta com link de convite
type GetInviteLinkResponse struct {
	InviteLink string `json:"inviteLink" example:"https://chat.whatsapp.com/ABC123"`
} //@name GetInviteLinkResponse

// JoinGroupRequest - Entrar no grupo via link
type JoinGroupRequest struct {
	Code string `json:"code" validate:"required" example:"ABC123DEF456"`
} //@name JoinGroupRequest

// CreateGroupRequest - Criar grupo
type CreateGroupRequest struct {
	Name         string   `json:"name" validate:"required" example:"Meu Grupo"`
	Participants []string `json:"participants" validate:"required,min=1" example:"5511999999999,5511888888888"`
} //@name CreateGroupRequest

// SetGroupLockedRequest - Bloquear configurações do grupo
type SetGroupLockedRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Locked   bool   `json:"locked" validate:"required" example:"true"`
} //@name SetGroupLockedRequest

// SetDisappearingRequest - Configurar mensagens temporárias
type SetDisappearingRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Duration string `json:"duration" validate:"required,oneof=24h 7d 90d off" example:"7d"`
} //@name SetDisappearingRequest

// RemovePhotoRequest - Remover foto do grupo
type RemovePhotoRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
} //@name RemovePhotoRequest

// UpdateParticipantsRequest - Gerenciar participantes
type UpdateParticipantsRequest struct {
	GroupJID     string   `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Participants []string `json:"participants" validate:"required,min=1" example:"5511999999999"`
	Action       string   `json:"action" validate:"required,oneof=add remove promote demote" example:"add"`
} //@name UpdateParticipantsRequest

// GetInviteInfoRequest - Obter informações do convite
type GetInviteInfoRequest struct {
	Code string `json:"code" validate:"required" example:"ABC123DEF456"`
} //@name GetInviteInfoRequest

// SetGroupPhotoRequest - Definir foto do grupo
type SetGroupPhotoRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Image    string `json:"image" validate:"required" example:"data:image/jpeg;base64,..."`
} //@name SetGroupPhotoRequest

// SetGroupPhotoResponse - Resposta ao definir foto
type SetGroupPhotoResponse struct {
	PictureID string `json:"pictureId" example:"abc123"`
} //@name SetGroupPhotoResponse

// SetGroupNameRequest - Alterar nome do grupo
type SetGroupNameRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Name     string `json:"name" validate:"required" example:"Novo Nome"`
} //@name SetGroupNameRequest

// SetGroupTopicRequest - Alterar descrição do grupo
type SetGroupTopicRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Topic    string `json:"topic" validate:"required" example:"Nova descrição"`
} //@name SetGroupTopicRequest

// LeaveGroupRequest - Sair do grupo
type LeaveGroupRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
} //@name LeaveGroupRequest

// SetGroupAnnounceRequest - Configurar modo anúncio
type SetGroupAnnounceRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Announce bool   `json:"announce" validate:"required" example:"true"`
} //@name SetGroupAnnounceRequest

// GroupActionResponse - Resposta genérica de ações de grupo
type GroupActionResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation completed successfully"`
} //@name GroupActionResponse
