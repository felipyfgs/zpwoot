package dto


type CreateCommunityRequest struct {
	Name         string   `json:"name" binding:"required" example:"Minha Comunidade"`
	Description  string   `json:"description,omitempty" example:"Descrição da comunidade"`
	Participants []string `json:"participants,omitempty" example:"5511999999999,5511888888888"`
} //@name CreateCommunityRequest


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


type ListCommunitiesResponse struct {
	Communities []CommunityInfo `json:"communities"`
} //@name ListCommunitiesResponse


type LinkGroupRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"123456789@g.us"`
} //@name LinkGroupRequest


type UnlinkGroupRequest struct {
	GroupJID string `json:"group_jid" binding:"required" example:"123456789@g.us"`
} //@name UnlinkGroupRequest


type CommunitySubGroup struct {
	JID          string   `json:"jid" example:"123456789@g.us"`
	Name         string   `json:"name" example:"Subgrupo 1"`
	Topic        string   `json:"topic,omitempty" example:"Descrição do subgrupo"`
	Participants []string `json:"participants,omitempty"`
	IsAnnounce   bool     `json:"is_announce" example:"false"`
	IsLocked     bool     `json:"is_locked" example:"false"`
	CreatedAt    int64    `json:"created_at,omitempty" example:"1696570882"`
} //@name CommunitySubGroup


type ListCommunitySubGroupsResponse struct {
	SubGroups []CommunitySubGroup `json:"sub_groups"`
} //@name ListCommunitySubGroupsResponse


type CommunityParticipant struct {
	JID  string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	Name string `json:"name,omitempty" example:"João Silva"`
	Role string `json:"role" example:"member"`
} //@name CommunityParticipant


type ListCommunityParticipantsResponse struct {
	Participants []CommunityParticipant `json:"participants"`
} //@name ListCommunityParticipantsResponse
