package waclient

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

const (
	ownerRole = "owner"
)

type CommunityService struct {
	waClient *WAClient
}

func NewCommunityService(waClient *WAClient) input.CommunityService {
	return &CommunityService{waClient: waClient}
}

func (cs *CommunityService) ListCommunities(ctx context.Context, sessionID string) (*dto.ListCommunitiesResponse, error) {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	groups, err := client.WAClient.GetJoinedGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get joined groups: %w", err)
	}

	response := &dto.ListCommunitiesResponse{
		Communities: make([]dto.CommunityInfo, 0),
	}

	for _, group := range groups {
		if group.IsParent {
			community := dto.CommunityInfo{
				JID:               group.JID.String(),
				Name:              group.Name,
				Description:       group.Topic,
				IsOwner:           !group.OwnerJID.IsEmpty() && group.OwnerJID.String() == client.WAClient.Store.ID.String(),
				ParticipantCount:  len(group.Participants),
				LinkedGroupsCount: 0,
				CreatedAt:         group.GroupCreated.Unix(),
			}

			if client.WAClient.Store.ID != nil {
				for _, participant := range group.Participants {
					if participant.JID.String() == client.WAClient.Store.ID.String() {
						community.IsAdmin = participant.IsAdmin || participant.IsSuperAdmin
						break
					}
				}
			}

			response.Communities = append(response.Communities, community)
		}
	}

	return response, nil
}
func (cs *CommunityService) GetCommunityInfo(ctx context.Context, sessionID string, communityJID string) (*dto.CommunityInfo, error) {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	jid, err := parseJID(communityJID)
	if err != nil {
		return nil, fmt.Errorf("invalid community JID: %w", err)
	}

	group, err := client.WAClient.GetGroupInfo(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get community info: %w", err)
	}

	if !group.IsParent {
		return nil, fmt.Errorf("JID is not a community")
	}

	community := &dto.CommunityInfo{
		JID:               group.JID.String(),
		Name:              group.Name,
		Description:       group.Topic,
		IsOwner:           !group.OwnerJID.IsEmpty() && group.OwnerJID.String() == client.WAClient.Store.ID.String(),
		ParticipantCount:  len(group.Participants),
		LinkedGroupsCount: 0,
		CreatedAt:         group.GroupCreated.Unix(),
	}

	if client.WAClient.Store.ID != nil {
		for _, participant := range group.Participants {
			if participant.JID.String() == client.WAClient.Store.ID.String() {
				community.IsAdmin = participant.IsAdmin || participant.IsSuperAdmin
				break
			}
		}
	}

	community.LinkedGroupsCount = 0

	return community, nil
}
func (cs *CommunityService) CreateCommunity(ctx context.Context, sessionID string, req *dto.CreateCommunityRequest) (*dto.CommunityInfo, error) {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	if req.Name == "" {
		return nil, fmt.Errorf("community name is required")
	}

	participantJIDs := make([]types.JID, len(req.Participants))

	for i, phone := range req.Participants {
		jid, err := parseJID(phone)
		if err != nil {
			return nil, fmt.Errorf("invalid participant phone %s: %w", phone, err)
		}

		participantJIDs[i] = jid
	}

	createReq := whatsmeow.ReqCreateGroup{
		Name:         req.Name,
		Participants: participantJIDs,
	}

	group, err := client.WAClient.CreateGroup(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}

	return &dto.CommunityInfo{
		JID:               group.JID.String(),
		Name:              group.Name,
		Description:       req.Description,
		IsOwner:           true,
		IsAdmin:           true,
		ParticipantCount:  len(group.Participants),
		LinkedGroupsCount: 0,
		CreatedAt:         group.GroupCreated.Unix(),
	}, nil
}
func (cs *CommunityService) LinkGroup(ctx context.Context, sessionID string, communityJID string, req *dto.LinkGroupRequest) error {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	parentJID, err := parseJID(communityJID)
	if err != nil {
		return fmt.Errorf("invalid community JID: %w", err)
	}

	childJID, err := parseJID(req.GroupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	err = client.WAClient.LinkGroup(parentJID, childJID)
	if err != nil {
		return fmt.Errorf("failed to link group to community: %w", err)
	}

	return nil
}
func (cs *CommunityService) UnlinkGroup(ctx context.Context, sessionID string, communityJID string, req *dto.UnlinkGroupRequest) error {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return ErrNotConnected
	}

	parentJID, err := parseJID(communityJID)
	if err != nil {
		return fmt.Errorf("invalid community JID: %w", err)
	}

	childJID, err := parseJID(req.GroupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	err = client.WAClient.UnlinkGroup(parentJID, childJID)
	if err != nil {
		return fmt.Errorf("failed to unlink group from community: %w", err)
	}

	return nil
}
func (cs *CommunityService) GetSubGroups(ctx context.Context, sessionID string, communityJID string) (*dto.ListCommunitySubGroupsResponse, error) {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	response := &dto.ListCommunitySubGroupsResponse{
		SubGroups: make([]dto.CommunitySubGroup, 0),
	}

	return response, nil
}
func (cs *CommunityService) GetParticipants(ctx context.Context, sessionID string, communityJID string) (*dto.ListCommunityParticipantsResponse, error) {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	jid, err := parseJID(communityJID)
	if err != nil {
		return nil, fmt.Errorf("invalid community JID: %w", err)
	}

	group, err := client.WAClient.GetGroupInfo(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get community info: %w", err)
	}

	response := &dto.ListCommunityParticipantsResponse{
		Participants: make([]dto.CommunityParticipant, 0, len(group.Participants)),
	}

	for _, participant := range group.Participants {
		communityParticipant := dto.CommunityParticipant{
			JID:  participant.JID.String(),
			Role: "member",
		}
		if participant.IsSuperAdmin {
			communityParticipant.Role = ownerRole
		} else if participant.IsAdmin {
			communityParticipant.Role = "admin"
		}

		contact, err := client.WAClient.Store.Contacts.GetContact(ctx, participant.JID)
		if err == nil {
			if contact.FullName != "" {
				communityParticipant.Name = contact.FullName
			} else if contact.PushName != "" {
				communityParticipant.Name = contact.PushName
			}
		}

		response.Participants = append(response.Participants, communityParticipant)
	}

	return response, nil
}
