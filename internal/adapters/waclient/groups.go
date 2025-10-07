package waclient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"zpwoot/internal/core/application/dto"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type GroupService struct {
	waClient *WAClient
}

func NewGroupService(waClient *WAClient) *GroupService {
	return &GroupService{
		waClient: waClient,
	}
}
func (gs *GroupService) ListGroups(ctx context.Context, sessionID string) (*dto.ListGroupsResponse, error) {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	groups, err := client.WAClient.GetJoinedGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get joined groups: %w", err)
	}

	response := &dto.ListGroupsResponse{
		Groups: make([]dto.WhatsAppGroupInfo, 0, len(groups)),
	}

	for _, group := range groups {
		response.Groups = append(response.Groups, dto.WhatsAppGroupInfo{
			JID:        group.JID.String(),
			Name:       group.Name,
			Topic:      group.Topic,
			IsAnnounce: group.IsAnnounce,
			IsLocked:   group.IsLocked,
		})
	}

	return response, nil
}
func (gs *GroupService) GetGroupInfo(ctx context.Context, sessionID string, groupJID string) (*dto.WhatsAppGroupInfo, error) {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID: %w", err)
	}

	group, err := client.WAClient.GetGroupInfo(jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %w", err)
	}

	participants := make([]string, 0, len(group.Participants))
	for _, p := range group.Participants {
		participants = append(participants, p.JID.String())
	}

	return &dto.WhatsAppGroupInfo{
		JID:          group.JID.String(),
		Name:         group.Name,
		Topic:        group.Topic,
		Participants: participants,
		IsAnnounce:   group.IsAnnounce,
		IsLocked:     group.IsLocked,
		CreatedAt:    group.GroupCreated.Unix(),
	}, nil
}
func (gs *GroupService) GetGroupInviteInfo(ctx context.Context, sessionID string, code string) (*dto.WhatsAppGroupInfo, error) {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if code == "" {
		return nil, errors.New("invite code is required")
	}

	group, err := client.WAClient.GetGroupInfoFromLink(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info from link: %w", err)
	}

	return &dto.WhatsAppGroupInfo{
		JID:        group.JID.String(),
		Name:       group.Name,
		Topic:      group.Topic,
		IsAnnounce: group.IsAnnounce,
		IsLocked:   group.IsLocked,
		CreatedAt:  group.GroupCreated.Unix(),
	}, nil
}
func (gs *GroupService) GetGroupInviteLink(ctx context.Context, sessionID string, groupJID string, reset bool) (string, error) {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return "", fmt.Errorf("invalid group JID: %w", err)
	}

	link, err := client.WAClient.GetGroupInviteLink(jid, reset)
	if err != nil {
		return "", fmt.Errorf("failed to get group invite link: %w", err)
	}

	return link, nil
}
func (gs *GroupService) JoinGroup(ctx context.Context, sessionID string, code string) error {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if code == "" {
		return errors.New("invite code is required")
	}

	_, err = client.WAClient.JoinGroupWithLink(code)
	if err != nil {
		return fmt.Errorf("failed to join group: %w", err)
	}

	return nil
}
func (gs *GroupService) CreateGroup(ctx context.Context, sessionID string, name string, participants []string) (*dto.WhatsAppGroupInfo, error) {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if name == "" {
		return nil, errors.New("group name is required")
	}

	if len(participants) < 1 {
		return nil, errors.New("at least one participant is required")
	}

	participantJIDs := make([]types.JID, len(participants))

	for i, phone := range participants {
		jid, err := parseJID(phone)
		if err != nil {
			return nil, fmt.Errorf("invalid participant phone %s: %w", phone, err)
		}

		participantJIDs[i] = jid
	}

	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: participantJIDs,
	}

	group, err := client.WAClient.CreateGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	participantStrings := make([]string, len(group.Participants))
	for i, p := range group.Participants {
		participantStrings[i] = p.JID.String()
	}

	return &dto.WhatsAppGroupInfo{
		JID:          group.JID.String(),
		Name:         group.Name,
		Topic:        group.Topic,
		Participants: participantStrings,
		IsAnnounce:   group.IsAnnounce,
		IsLocked:     group.IsLocked,
		CreatedAt:    group.GroupCreated.Unix(),
	}, nil
}
func (gs *GroupService) LeaveGroup(ctx context.Context, sessionID string, groupJID string) error {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	err = client.WAClient.LeaveGroup(jid)
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}

	return nil
}
func (gs *GroupService) UpdateGroupParticipants(ctx context.Context, sessionID string, groupJID string, participants []string, action string) error {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	if len(participants) < 1 {
		return errors.New("at least one participant is required")
	}

	participantJIDs := make([]types.JID, len(participants))

	for i, phone := range participants {
		pjid, err := parseJID(phone)
		if err != nil {
			return fmt.Errorf("invalid participant phone %s: %w", phone, err)
		}

		participantJIDs[i] = pjid
	}

	var participantChange whatsmeow.ParticipantChange
	switch action {
	case "add":
		participantChange = whatsmeow.ParticipantChangeAdd
	case "remove":
		participantChange = whatsmeow.ParticipantChangeRemove
	case "promote":
		participantChange = whatsmeow.ParticipantChangePromote
	case "demote":
		participantChange = whatsmeow.ParticipantChangeDemote
	default:
		return fmt.Errorf("invalid action: %s (must be add, remove, promote, or demote)", action)
	}

	_, err = client.WAClient.UpdateGroupParticipants(jid, participantJIDs, participantChange)
	if err != nil {
		return fmt.Errorf("failed to update group participants: %w", err)
	}

	return nil
}
func (gs *GroupService) SetGroupName(ctx context.Context, sessionID string, groupJID string, name string) error {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	if name == "" {
		return errors.New("group name is required")
	}

	err = client.WAClient.SetGroupName(jid, name)
	if err != nil {
		return fmt.Errorf("failed to set group name: %w", err)
	}

	return nil
}
func (gs *GroupService) SetGroupTopic(ctx context.Context, sessionID string, groupJID string, topic string) error {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	err = client.WAClient.SetGroupTopic(jid, "", "", topic)
	if err != nil {
		return fmt.Errorf("failed to set group topic: %w", err)
	}

	return nil
}
func (gs *GroupService) SetGroupLocked(ctx context.Context, sessionID string, groupJID string, locked bool) error {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	err = client.WAClient.SetGroupLocked(jid, locked)
	if err != nil {
		return fmt.Errorf("failed to set group locked: %w", err)
	}

	return nil
}
func (gs *GroupService) SetGroupAnnounce(ctx context.Context, sessionID string, groupJID string, announce bool) error {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	err = client.WAClient.SetGroupAnnounce(jid, announce)
	if err != nil {
		return fmt.Errorf("failed to set group announce: %w", err)
	}

	return nil
}
func (gs *GroupService) SetDisappearingTimer(ctx context.Context, sessionID string, groupJID string, duration string) error {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	var timer time.Duration
	switch duration {
	case "24h":
		timer = 24 * time.Hour
	case "7d":
		timer = 7 * 24 * time.Hour
	case "90d":
		timer = 90 * 24 * time.Hour
	case "off":
		timer = 0
	default:
		return fmt.Errorf("invalid duration: %s (must be 24h, 7d, 90d, or off)", duration)
	}

	err = client.WAClient.SetDisappearingTimer(jid, timer, time.Now())
	if err != nil {
		return fmt.Errorf("failed to set disappearing timer: %w", err)
	}

	return nil
}
func (gs *GroupService) SetGroupPhoto(ctx context.Context, sessionID string, groupJID string, imageData []byte) (string, error) {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return "", fmt.Errorf("invalid group JID: %w", err)
	}

	if len(imageData) == 0 {
		return "", errors.New("image data is required")
	}

	if len(imageData) < 3 || imageData[0] != 0xFF || imageData[1] != 0xD8 || imageData[2] != 0xFF {
		return "", errors.New("image must be in JPEG format")
	}

	pictureID, err := client.WAClient.SetGroupPhoto(jid, imageData)
	if err != nil {
		return "", fmt.Errorf("failed to set group photo: %w", err)
	}

	return pictureID, nil
}
func (gs *GroupService) RemoveGroupPhoto(ctx context.Context, sessionID string, groupJID string) error {
	client, err := gs.waClient.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	jid, err := parseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	_, err = client.WAClient.SetGroupPhoto(jid, nil)
	if err != nil {
		return fmt.Errorf("failed to remove group photo: %w", err)
	}

	return nil
}

