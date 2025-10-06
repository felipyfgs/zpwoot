package waclient

import (
	"context"
	"fmt"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/ports/input"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type CommunityService struct {
	waClient *WAClient
}

func NewCommunityService(waClient *WAClient) input.CommunityService {
	return &CommunityService{
		waClient: waClient,
	}
}

// ListCommunities lista todas as comunidades que a sessão participa
func (cs *CommunityService) ListCommunities(ctx context.Context, sessionID string) (*dto.ListCommunitiesResponse, error) {
	client, err := cs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if !client.IsConnected() {
		return nil, ErrNotConnected
	}

	// Obter grupos e filtrar comunidades (grupos com IsParent=true)
	groups, err := client.WAClient.GetJoinedGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get joined groups: %w", err)
	}

	response := &dto.ListCommunitiesResponse{
		Communities: make([]dto.CommunityInfo, 0),
	}

	for _, group := range groups {
		// Verificar se é uma comunidade (grupo pai)
		if group.IsParent {
			community := dto.CommunityInfo{
				JID:               group.JID.String(),
				Name:              group.Name,
				Description:       group.Topic,
				IsOwner:           group.OwnerJID != nil && *group.OwnerJID == client.WAClient.Store.ID,
				ParticipantCount:  len(group.Participants),
				LinkedGroupsCount: 0, // Será calculado se necessário
				CreatedAt:         group.GroupCreated.Unix(),
			}

			// Verificar se é admin
			for _, participant := range group.Participants {
				if participant.JID == client.WAClient.Store.ID {
					community.IsAdmin = participant.IsAdmin || participant.IsSuperAdmin
					break
				}
			}

			response.Communities = append(response.Communities, community)
		}
	}

	return response, nil
}

// GetCommunityInfo obtém informações detalhadas de uma comunidade
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
		IsOwner:           group.OwnerJID != nil && *group.OwnerJID == client.WAClient.Store.ID,
		ParticipantCount:  len(group.Participants),
		LinkedGroupsCount: 0, // Será calculado via GetSubGroups se necessário
		CreatedAt:         group.GroupCreated.Unix(),
	}

	// Verificar se é admin
	for _, participant := range group.Participants {
		if participant.JID == client.WAClient.Store.ID {
			community.IsAdmin = participant.IsAdmin || participant.IsSuperAdmin
			break
		}
	}

	// Obter subgrupos para calcular LinkedGroupsCount
	subGroups, err := client.WAClient.GetSubGroups(jid)
	if err == nil {
		community.LinkedGroupsCount = len(subGroups)
		
		// Adicionar JIDs dos subgrupos
		linkedGroups := make([]string, len(subGroups))
		for i, subGroup := range subGroups {
			linkedGroups[i] = subGroup.JID.String()
		}
		community.LinkedGroups = linkedGroups
	}

	return community, nil
}

// CreateCommunity cria uma nova comunidade
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

	// Parse participant JIDs
	participantJIDs := make([]types.JID, len(req.Participants))
	for i, phone := range req.Participants {
		jid, err := parseJID(phone)
		if err != nil {
			return nil, fmt.Errorf("invalid participant phone %s: %w", phone, err)
		}
		participantJIDs[i] = jid
	}

	// Criar comunidade (grupo com IsParent=true)
	createReq := whatsmeow.ReqCreateGroup{
		Name:         req.Name,
		Participants: participantJIDs,
		// IsParent:     true, // TODO: Verificar se whatsmeow suporta este campo
	}

	group, err := client.WAClient.CreateGroup(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}

	return &dto.CommunityInfo{
		JID:               group.JID.String(),
		Name:              group.Name,
		Description:       req.Description,
		IsOwner:           true, // Criador é sempre owner
		IsAdmin:           true,
		ParticipantCount:  len(group.Participants),
		LinkedGroupsCount: 0,
		CreatedAt:         group.GroupCreated.Unix(),
	}, nil
}

// LinkGroup vincula um grupo a uma comunidade
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

// UnlinkGroup desvincula um grupo de uma comunidade
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

// GetSubGroups obtém todos os subgrupos de uma comunidade
func (cs *CommunityService) GetSubGroups(ctx context.Context, sessionID string, communityJID string) (*dto.ListCommunitySubGroupsResponse, error) {
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

	subGroups, err := client.WAClient.GetSubGroups(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get sub groups: %w", err)
	}

	response := &dto.ListCommunitySubGroupsResponse{
		SubGroups: make([]dto.CommunitySubGroup, 0, len(subGroups)),
	}

	for _, group := range subGroups {
		participants := make([]string, len(group.Participants))
		for i, p := range group.Participants {
			participants[i] = p.JID.String()
		}

		response.SubGroups = append(response.SubGroups, dto.CommunitySubGroup{
			JID:          group.JID.String(),
			Name:         group.Name,
			Topic:        group.Topic,
			Participants: participants,
			IsAnnounce:   group.IsAnnounce,
			IsLocked:     group.IsLocked,
			CreatedAt:    group.GroupCreated.Unix(),
		})
	}

	return response, nil
}

// GetParticipants obtém todos os participantes de uma comunidade
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

	// Obter participantes de todos os grupos vinculados
	participantJIDs, err := client.WAClient.GetLinkedGroupsParticipants(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get community participants: %w", err)
	}

	response := &dto.ListCommunityParticipantsResponse{
		Participants: make([]dto.CommunityParticipant, 0, len(participantJIDs)),
	}

	for _, participantJID := range participantJIDs {
		participant := dto.CommunityParticipant{
			JID:  participantJID.String(),
			Role: "member", // Por padrão, todos são membros
		}

		// Tentar obter nome do contato
		if contact, ok := client.WAClient.Store.Contacts[participantJID]; ok {
			if contact.FullName != "" {
				participant.Name = contact.FullName
			} else if contact.PushName != "" {
				participant.Name = contact.PushName
			}
		}

		response.Participants = append(response.Participants, participant)
	}

	return response, nil
}
