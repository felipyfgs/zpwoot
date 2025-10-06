package input

import (
	"context"
	"zpwoot/internal/core/application/dto"
)

type CommunityService interface {
	ListCommunities(ctx context.Context, sessionID string) (*dto.ListCommunitiesResponse, error)
	GetCommunityInfo(ctx context.Context, sessionID string, communityJID string) (*dto.CommunityInfo, error)
	CreateCommunity(ctx context.Context, sessionID string, req *dto.CreateCommunityRequest) (*dto.CommunityInfo, error)
	LinkGroup(ctx context.Context, sessionID string, communityJID string, req *dto.LinkGroupRequest) error
	UnlinkGroup(ctx context.Context, sessionID string, communityJID string, req *dto.UnlinkGroupRequest) error
	GetSubGroups(ctx context.Context, sessionID string, communityJID string) (*dto.ListCommunitySubGroupsResponse, error)
	GetParticipants(ctx context.Context, sessionID string, communityJID string) (*dto.ListCommunityParticipantsResponse, error)
}
