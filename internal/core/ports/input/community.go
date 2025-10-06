package input

import (
	"context"
	"zpwoot/internal/core/application/dto"
)

// CommunityService define operações de gerenciamento de comunidades WhatsApp
type CommunityService interface {
	// Informações
	ListCommunities(ctx context.Context, sessionID string) (*dto.ListCommunitiesResponse, error)
	GetCommunityInfo(ctx context.Context, sessionID string, communityJID string) (*dto.CommunityInfo, error)

	// Criação
	CreateCommunity(ctx context.Context, sessionID string, req *dto.CreateCommunityRequest) (*dto.CommunityInfo, error)

	// Gerenciamento de grupos
	LinkGroup(ctx context.Context, sessionID string, communityJID string, req *dto.LinkGroupRequest) error
	UnlinkGroup(ctx context.Context, sessionID string, communityJID string, req *dto.UnlinkGroupRequest) error
	GetSubGroups(ctx context.Context, sessionID string, communityJID string) (*dto.ListCommunitySubGroupsResponse, error)

	// Participantes
	GetParticipants(ctx context.Context, sessionID string, communityJID string) (*dto.ListCommunityParticipantsResponse, error)
}
