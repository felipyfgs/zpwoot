package input

import (
	"context"
	"zpwoot/internal/core/application/dto"
)

// GroupService define operações de gerenciamento de grupos WhatsApp
type GroupService interface {
	// Informações
	ListGroups(ctx context.Context, sessionID string) (*dto.ListGroupsResponse, error)
	GetGroupInfo(ctx context.Context, sessionID string, groupJID string) (*dto.WhatsAppGroupInfo, error)
	GetGroupInviteInfo(ctx context.Context, sessionID string, code string) (*dto.WhatsAppGroupInfo, error)

	// Convites
	GetGroupInviteLink(ctx context.Context, sessionID string, groupJID string, reset bool) (string, error)
	JoinGroup(ctx context.Context, sessionID string, code string) error

	// Gerenciamento
	CreateGroup(ctx context.Context, sessionID string, name string, participants []string) (*dto.WhatsAppGroupInfo, error)
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
