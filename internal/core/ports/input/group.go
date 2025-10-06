package input

import (
	"context"
	"zpwoot/internal/core/application/dto"
)

type GroupService interface {
	ListGroups(ctx context.Context, sessionID string) (*dto.ListGroupsResponse, error)
	GetGroupInfo(ctx context.Context, sessionID string, groupJID string) (*dto.WhatsAppGroupInfo, error)
	GetGroupInviteInfo(ctx context.Context, sessionID string, code string) (*dto.WhatsAppGroupInfo, error)
	GetGroupInviteLink(ctx context.Context, sessionID string, groupJID string, reset bool) (string, error)
	JoinGroup(ctx context.Context, sessionID string, code string) error
	CreateGroup(ctx context.Context, sessionID string, name string, participants []string) (*dto.WhatsAppGroupInfo, error)
	LeaveGroup(ctx context.Context, sessionID string, groupJID string) error
	UpdateGroupParticipants(ctx context.Context, sessionID string, groupJID string, participants []string, action string) error
	SetGroupName(ctx context.Context, sessionID string, groupJID string, name string) error
	SetGroupTopic(ctx context.Context, sessionID string, groupJID string, topic string) error
	SetGroupLocked(ctx context.Context, sessionID string, groupJID string, locked bool) error
	SetGroupAnnounce(ctx context.Context, sessionID string, groupJID string, announce bool) error
	SetDisappearingTimer(ctx context.Context, sessionID string, groupJID string, duration string) error
	SetGroupPhoto(ctx context.Context, sessionID string, groupJID string, imageData []byte) (string, error)
	RemoveGroupPhoto(ctx context.Context, sessionID string, groupJID string) error
}
