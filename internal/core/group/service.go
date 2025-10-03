package group

import (
	"fmt"
	"regexp"
	"strings"
)

type service struct {
	validator Validator
}

func NewService(validator Validator) Service {
	return &service{
		validator: validator,
	}
}

func (s *service) ValidateGroupCreation(req *CreateGroupRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if err := s.ValidateGroupName(req.Name); err != nil {
		return err
	}

	if err := s.ValidateGroupDescription(req.Description); err != nil {
		return err
	}

	if err := s.ValidateParticipants(req.Participants); err != nil {
		return err
	}

	return nil
}

func (s *service) ValidateGroupName(name string) error {
	if name == "" {
		return fmt.Errorf("group name cannot be empty")
	}

	if len(name) > 25 {
		return fmt.Errorf("group name cannot exceed 25 characters")
	}

	if strings.Contains(name, "\n") || strings.Contains(name, "\r") {
		return fmt.Errorf("group name cannot contain line breaks")
	}

	return nil
}

func (s *service) ValidateGroupDescription(description string) error {
	if len(description) > 512 {
		return fmt.Errorf("group description cannot exceed 512 characters")
	}

	return nil
}

func (s *service) ValidateParticipants(participants []string) error {
	if len(participants) == 0 {
		return fmt.Errorf("at least one participant is required")
	}

	if len(participants) > 256 {
		return fmt.Errorf("too many participants (max 256)")
	}

	seen := make(map[string]bool)
	for _, participant := range participants {
		if err := s.ValidateJID(participant); err != nil {
			return fmt.Errorf("invalid participant %s: %w", participant, err)
		}

		normalized := s.NormalizeJID(participant)
		if seen[normalized] {
			return fmt.Errorf("duplicate participant: %s", participant)
		}
		seen[normalized] = true
	}

	return nil
}

func (s *service) ValidateInviteLink(inviteLink string) error {
	if inviteLink == "" {
		return fmt.Errorf("invite link cannot be empty")
	}

	whatsappLinkPattern := `^https://chat\.whatsapp\.com/[A-Za-z0-9]+$`
	matched, err := regexp.MatchString(whatsappLinkPattern, inviteLink)
	if err != nil {
		return fmt.Errorf("error validating invite link: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid WhatsApp invite link format")
	}

	return nil
}

func (s *service) ValidateJID(jid string) error {
	if jid == "" {
		return fmt.Errorf("JID cannot be empty")
	}

	jidPattern := `^[0-9]+@(s\.whatsapp\.net|g\.us)$`
	matched, err := regexp.MatchString(jidPattern, jid)
	if err != nil {
		return fmt.Errorf("error validating JID: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid JID format")
	}

	return nil
}

func (s *service) CanPerformAction(userJID, groupJID string, action GroupAction, groupInfo *GroupInfo) error {
	if groupInfo == nil {
		return fmt.Errorf("group not found")
	}

	if !groupInfo.HasParticipant(userJID) {
		return fmt.Errorf("user is not a participant of this group")
	}

	switch action {
	case GroupActionAddParticipant, GroupActionPromoteParticipant, GroupActionDemoteParticipant,
		GroupActionSetName, GroupActionSetDescription, GroupActionSetPhoto,
		GroupActionSetSettings, GroupActionGetInviteLink, GroupActionRevokeInviteLink:
		if !groupInfo.IsParticipantAdmin(userJID) {
			return fmt.Errorf("only group admins can perform this action")
		}

	case GroupActionRemoveParticipant:
		if !groupInfo.IsParticipantAdmin(userJID) {
			return fmt.Errorf("only group admins can perform this action")
		}

		for _, participant := range groupInfo.Participants {
			if participant.JID == userJID && userJID == groupInfo.Owner {
				return fmt.Errorf("cannot remove group owner")
			}
		}

	case GroupActionLeave:
		if userJID == groupInfo.Owner {
			return fmt.Errorf("group owner cannot leave the group")
		}

	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	return nil
}

func (s *service) IsGroupAdmin(userJID, groupJID string, groupInfo *GroupInfo) bool {
	if groupInfo == nil {
		return false
	}
	return groupInfo.IsParticipantAdmin(userJID)
}

func (s *service) IsGroupOwner(userJID, groupJID string, groupInfo *GroupInfo) bool {
	if groupInfo == nil {
		return false
	}
	return groupInfo.Owner == userJID
}

func (s *service) ProcessParticipantChanges(req *UpdateParticipantsRequest, currentGroup *GroupInfo) error {
	if req == nil || currentGroup == nil {
		return fmt.Errorf("invalid request or group info")
	}

	switch req.Action {
	case "add":
		for _, participant := range req.Participants {
			if currentGroup.HasParticipant(participant) {
				return fmt.Errorf("participant %s is already in the group", participant)
			}
			if err := s.ValidateJID(participant); err != nil {
				return fmt.Errorf("invalid participant %s: %w", participant, err)
			}
		}

	case "remove":
		for _, participant := range req.Participants {
			if !currentGroup.HasParticipant(participant) {
				return fmt.Errorf("participant %s is not in the group", participant)
			}
			if participant == currentGroup.Owner {
				return fmt.Errorf("cannot remove group owner")
			}
		}

	case "promote":
		for _, participant := range req.Participants {
			if !currentGroup.HasParticipant(participant) {
				return fmt.Errorf("participant %s is not in the group", participant)
			}
			if currentGroup.IsParticipantAdmin(participant) {
				return fmt.Errorf("participant %s is already an admin", participant)
			}
		}

	case "demote":
		for _, participant := range req.Participants {
			if !currentGroup.HasParticipant(participant) {
				return fmt.Errorf("participant %s is not in the group", participant)
			}
			if !currentGroup.IsParticipantAdmin(participant) {
				return fmt.Errorf("participant %s is not an admin", participant)
			}
			if participant == currentGroup.Owner {
				return fmt.Errorf("cannot demote group owner")
			}
		}

	default:
		return fmt.Errorf("invalid action: %s", req.Action)
	}

	return nil
}

func (s *service) ProcessSettingsChanges(req *UpdateGroupSettingsRequest, currentGroup *GroupInfo) error {
	if req == nil || currentGroup == nil {
		return fmt.Errorf("invalid request or group info")
	}

	if req.JoinApprovalMode != "" {
		if req.JoinApprovalMode != "auto" && req.JoinApprovalMode != "admin_approval" {
			return fmt.Errorf("invalid join approval mode: %s", req.JoinApprovalMode)
		}
	}

	if req.MemberAddMode != "" {
		if req.MemberAddMode != "all_members" && req.MemberAddMode != "only_admins" {
			return fmt.Errorf("invalid member add mode: %s", req.MemberAddMode)
		}
	}

	return nil
}

func (s *service) NormalizeJID(jid string) string {

	normalized := strings.ToLower(strings.TrimSpace(jid))

	if !strings.Contains(normalized, "@") {
		normalized += "@s.whatsapp.net"
	}

	return normalized
}

func (s *service) ExtractPhoneNumber(jid string) string {
	parts := strings.Split(jid, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return jid
}

func (s *service) FormatGroupJID(groupID string) string {
	if strings.Contains(groupID, "@") {
		return groupID
	}
	return groupID + "@g.us"
}

type defaultValidator struct{}

func NewDefaultValidator() Validator {
	return &defaultValidator{}
}

func (v *defaultValidator) ValidateGroupName(name string) error {
	if name == "" {
		return fmt.Errorf("group name cannot be empty")
	}
	if len(name) > 25 {
		return fmt.Errorf("group name cannot exceed 25 characters")
	}
	return nil
}

func (v *defaultValidator) ValidateGroupDescription(description string) error {
	if len(description) > 512 {
		return fmt.Errorf("group description cannot exceed 512 characters")
	}
	return nil
}

func (v *defaultValidator) ValidateParticipantJID(jid string) error {
	if jid == "" {
		return fmt.Errorf("JID cannot be empty")
	}
	jidPattern := `^[0-9]+@(s\.whatsapp\.net|g\.us)$`
	matched, _ := regexp.MatchString(jidPattern, jid)
	if !matched {
		return fmt.Errorf("invalid JID format")
	}
	return nil
}

func (v *defaultValidator) ValidateInviteLink(link string) error {
	if link == "" {
		return fmt.Errorf("invite link cannot be empty")
	}
	whatsappLinkPattern := `^https://chat\.whatsapp\.com/[A-Za-z0-9]+$`
	matched, _ := regexp.MatchString(whatsappLinkPattern, link)
	if !matched {
		return fmt.Errorf("invalid WhatsApp invite link format")
	}
	return nil
}

func (v *defaultValidator) ValidateGroupSettings(settings *GroupSettings) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	if settings.JoinApprovalMode != "" &&
		settings.JoinApprovalMode != "auto" &&
		settings.JoinApprovalMode != "admin_approval" {
		return fmt.Errorf("invalid join approval mode")
	}

	if settings.MemberAddMode != "" &&
		settings.MemberAddMode != "all_members" &&
		settings.MemberAddMode != "only_admins" {
		return fmt.Errorf("invalid member add mode")
	}

	return nil
}
