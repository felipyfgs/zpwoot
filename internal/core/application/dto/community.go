package dto

// CreateCommunityRequest - Requisição para criar comunidade
type CreateCommunityRequest struct {
	Name         string   `json:"name" binding:"required" example:"Minha Comunidade"`
	Description  string   `json:"description,omitempty" example:"Descrição da comunidade"`
	Participants []string `json:"participants,omitempty" example:"5511999999999,5511888888888"`
} //@name CreateCommunityRequest

// CommunityInfo - Informações da comunidade
type CommunityInfo struct {
	JID                  string   `json:"jid" example:"123456789@g.us"`
	Name                 string   `json:"name" example:"Minha Comunidade"`
	Description          string   `json:"description,omitempty" example:"Descrição da comunidade"`
	AnnouncementGroupJID string   `json:"announcement_group_jid,omitempty" example:"123456789@g.us"`
	IsOwner              bool     `json:"is_owner" example:"true"`
	IsAdmin              bool     `json:"is_admin" example:"false"`
	ParticipantCount     int      `json:"participant_count" example:"25"`
	LinkedGroupsCount    int      `json:"linked_groups_count" example:"3"`
	LinkedGroups         []string `json:"linked_groups,omitempty"`
	CreatedAt            int64    `json:"created_at,omitempty" example:"1696570882"`
} //@name CommunityInfo

// ListCommunitiesResponse - Lista de comunidades
type ListCommunitiesResponse struct {
	Communities []CommunityInfo `json:"communities"`
} //@name ListCommunitiesResponse

// LinkGroupRequest - Requisição para vincular grupo à comunidade
type LinkGroupRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"123456789@g.us"`
} //@name LinkGroupRequest

// UnlinkGroupRequest - Requisição para desvincular grupo da comunidade
type UnlinkGroupRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"123456789@g.us"`
} //@name UnlinkGroupRequest

// CommunitySubGroup - Informações de subgrupo da comunidade
type CommunitySubGroup struct {
	JID          string   `json:"jid" example:"123456789@g.us"`
	Name         string   `json:"name" example:"Subgrupo 1"`
	Topic        string   `json:"topic,omitempty" example:"Descrição do subgrupo"`
	Participants []string `json:"participants,omitempty"`
	IsAnnounce   bool     `json:"is_announce" example:"false"`
	IsLocked     bool     `json:"is_locked" example:"false"`
	CreatedAt    int64    `json:"created_at,omitempty" example:"1696570882"`
} //@name CommunitySubGroup

// ListCommunitySubGroupsResponse - Lista de subgrupos da comunidade
type ListCommunitySubGroupsResponse struct {
	SubGroups []CommunitySubGroup `json:"sub_groups"`
} //@name ListCommunitySubGroupsResponse

// CommunityParticipant - Participante da comunidade
type CommunityParticipant struct {
	JID  string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	Name string `json:"name,omitempty" example:"João Silva"`
	Role string `json:"role" example:"member"` // owner, admin, member
} //@name CommunityParticipant

// ListCommunityParticipantsResponse - Lista de participantes da comunidade
type ListCommunityParticipantsResponse struct {
	Participants []CommunityParticipant `json:"participants"`
} //@name ListCommunityParticipantsResponse
