package dto


type ListGroupsResponse struct {
	Groups []WhatsAppGroupInfo `json:"groups"`
} //@name ListGroupsResponse


type WhatsAppGroupInfo struct {
	JID          string   `json:"jid" example:"123456789@g.us"`
	Name         string   `json:"name" example:"Meu Grupo"`
	Topic        string   `json:"topic,omitempty" example:"Descrição do grupo"`
	Participants []string `json:"participants,omitempty"`
	IsAnnounce   bool     `json:"isAnnounce" example:"false"`
	IsLocked     bool     `json:"isLocked" example:"false"`
	CreatedAt    int64    `json:"createdAt,omitempty" example:"1696570882"`
} //@name WhatsAppGroupInfo


type GetGroupInfoRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
} //@name GetGroupInfoRequest


type GetGroupInviteLinkRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Reset    bool   `json:"reset,omitempty" example:"false"`
} //@name GetGroupInviteLinkRequest


type GetInviteLinkResponse struct {
	InviteLink string `json:"inviteLink" example:"https://chat.whatsapp.com/ABC123"`
} //@name GetInviteLinkResponse


type JoinGroupRequest struct {
	Code string `json:"code" validate:"required" example:"ABC123DEF456"`
} //@name JoinGroupRequest


type CreateGroupRequest struct {
	Name         string   `json:"name" validate:"required" example:"Meu Grupo"`
	Participants []string `json:"participants" validate:"required,min=1" example:"5511999999999,5511888888888"`
} //@name CreateGroupRequest


type SetGroupLockedRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Locked   bool   `json:"locked" validate:"required" example:"true"`
} //@name SetGroupLockedRequest


type SetDisappearingRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Duration string `json:"duration" validate:"required,oneof=24h 7d 90d off" example:"7d"`
} //@name SetDisappearingRequest


type RemovePhotoRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
} //@name RemovePhotoRequest


type UpdateParticipantsRequest struct {
	GroupJID     string   `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Participants []string `json:"participants" validate:"required,min=1" example:"5511999999999"`
	Action       string   `json:"action" validate:"required,oneof=add remove promote demote" example:"add"`
} //@name UpdateParticipantsRequest


type GetInviteInfoRequest struct {
	Code string `json:"code" validate:"required" example:"ABC123DEF456"`
} //@name GetInviteInfoRequest


type SetGroupPhotoRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Image    string `json:"image" validate:"required" example:"data:image/jpeg;base64,..."`
} //@name SetGroupPhotoRequest


type SetGroupPhotoResponse struct {
	PictureID string `json:"pictureId" example:"abc123"`
} //@name SetGroupPhotoResponse


type SetGroupNameRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Name     string `json:"name" validate:"required" example:"Novo Nome"`
} //@name SetGroupNameRequest


type SetGroupTopicRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Topic    string `json:"topic" validate:"required" example:"Nova descrição"`
} //@name SetGroupTopicRequest


type LeaveGroupRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
} //@name LeaveGroupRequest


type SetGroupAnnounceRequest struct {
	GroupJID string `json:"groupJid" validate:"required" example:"123456789@g.us"`
	Announce bool   `json:"announce" validate:"required" example:"true"`
} //@name SetGroupAnnounceRequest


type GroupActionResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation completed successfully"`
} //@name GroupActionResponse
