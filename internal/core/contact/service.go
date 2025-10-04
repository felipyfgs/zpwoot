package contact

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repository Repository
	validator  ContactValidator
}

func NewService(repo Repository, validator ContactValidator) *Service {
	return &Service{
		repository: repo,
		validator:  validator,
	}
}

func (s *Service) CreateContact(ctx context.Context, req *CreateContactRequest) (*Contact, error) {
	if err := s.validator.ValidateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid create request: %w", err)
	}

	exists, err := s.repository.ExistsByZpJID(ctx, req.SessionID, req.ZpJID)
	if err != nil {
		return nil, fmt.Errorf("failed to check contact existence: %w", err)
	}
	if exists {
		return nil, ErrContactAlreadyExistsWithJID(req.ZpJID)
	}

	now := time.Now()
	contact := &Contact{
		ID:          uuid.New(),
		SessionID:   req.SessionID,
		ZpJID:       req.ZpJID,
		ZpName:      req.ZpName,
		ZpPushName:  req.ZpPushName,
		ZpShortName: req.ZpShortName,
		ZpAvatar:    req.ZpAvatar,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		IsGroup:     req.IsGroup,
		IsBusiness:  req.IsBusiness,
		SyncStatus:  string(SyncStatusPending),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repository.Create(ctx, contact); err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}

	return contact, nil
}

func (s *Service) GetContact(ctx context.Context, id uuid.UUID) (*Contact, error) {
	contact, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	return contact, nil
}

func (s *Service) UpdateContact(ctx context.Context, req *UpdateContactRequest) (*Contact, error) {
	if err := s.validator.ValidateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid update request: %w", err)
	}

	contact, err := s.repository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	contact.ZpName = req.ZpName
	contact.ZpPushName = req.ZpPushName
	contact.ZpShortName = req.ZpShortName
	contact.ZpAvatar = req.ZpAvatar
	contact.Email = req.Email
	contact.IsBlocked = req.IsBlocked
	contact.LastSeen = req.LastSeen
	contact.IsOnline = req.IsOnline
	contact.LastStatus = req.LastStatus
	contact.UpdatedAt = time.Now()

	if err := s.repository.Update(ctx, contact); err != nil {
		return nil, fmt.Errorf("failed to update contact: %w", err)
	}

	return contact, nil
}

func (s *Service) DeleteContact(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}

	return nil
}

func (s *Service) UpdateSyncStatus(ctx context.Context, id uuid.UUID, status SyncStatus, cwContactID, cwConversationID *int) error {
	if !IsValidSyncStatus(string(status)) {
		return ErrInvalidSyncStatusValue(string(status))
	}

	if err := s.repository.UpdateSyncStatus(ctx, id, status, cwContactID, cwConversationID); err != nil {
		return fmt.Errorf("failed to update sync status: %w", err)
	}

	return nil
}

func (s *Service) MarkAsSynced(ctx context.Context, id uuid.UUID, cwContactID, cwConversationID int) error {
	if err := s.repository.MarkAsSynced(ctx, id, cwContactID, cwConversationID); err != nil {
		return fmt.Errorf("failed to mark contact as synced: %w", err)
	}

	return nil
}

func (s *Service) MarkAsFailed(ctx context.Context, id uuid.UUID, errorReason string) error {
	if err := s.repository.MarkAsFailed(ctx, id, errorReason); err != nil {
		return fmt.Errorf("failed to mark contact as failed: %w", err)
	}

	return nil
}

func (s *Service) ListContacts(ctx context.Context, req *ListContactsRequest) ([]*Contact, int64, error) {
	if err := s.validateListRequest(req); err != nil {
		return nil, 0, fmt.Errorf("invalid list request: %w", err)
	}

	contacts, total, err := s.repository.List(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list contacts: %w", err)
	}

	return contacts, total, nil
}

func (s *Service) GetPendingSyncContacts(ctx context.Context, sessionID uuid.UUID, limit int) ([]*Contact, error) {
	contacts, err := s.repository.GetPendingSyncContacts(ctx, sessionID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending sync contacts: %w", err)
	}

	return contacts, nil
}

func (s *Service) GetStats(ctx context.Context) (*ContactStats, error) {
	stats, err := s.repository.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact stats: %w", err)
	}

	return stats, nil
}

func (s *Service) GetStatsBySession(ctx context.Context, sessionID uuid.UUID) (*ContactStats, error) {
	stats, err := s.repository.GetStatsBySession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact stats by session: %w", err)
	}

	return stats, nil
}

func (s *Service) validateListRequest(req *ListContactsRequest) error {
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	if req.SessionID != "" {
		if _, err := uuid.Parse(req.SessionID); err != nil {
			return NewContactValidationError("session_id", "invalid session ID format")
		}
	}

	return nil
}

type DefaultValidator struct{}

func NewDefaultValidator() ContactValidator {
	return &DefaultValidator{}
}

func (v *DefaultValidator) ValidateCreateRequest(req *CreateContactRequest) error {
	if req == nil {
		return NewContactError(ErrCodeInvalidContactData, "request cannot be nil", nil)
	}

	if req.SessionID == uuid.Nil {
		return NewContactValidationError("session_id", "session ID is required")
	}

	if err := v.ValidateZpJID(req.ZpJID); err != nil {
		return err
	}

	if req.Email != "" {
		if err := v.ValidateEmail(req.Email); err != nil {
			return err
		}
	}

	if req.PhoneNumber != "" {
		if err := v.ValidatePhoneNumber(req.PhoneNumber); err != nil {
			return err
		}
	}

	return nil
}

func (v *DefaultValidator) ValidateUpdateRequest(req *UpdateContactRequest) error {
	if req == nil {
		return NewContactError(ErrCodeInvalidContactData, "request cannot be nil", nil)
	}

	if req.ID == uuid.Nil {
		return NewContactValidationError("id", "contact ID is required")
	}

	if req.Email != "" {
		if err := v.ValidateEmail(req.Email); err != nil {
			return err
		}
	}

	return nil
}

func (v *DefaultValidator) ValidateContact(contact *Contact) error {
	if contact == nil {
		return NewContactError(ErrCodeInvalidContactData, "contact cannot be nil", nil)
	}

	if contact.SessionID == uuid.Nil {
		return NewContactValidationError("session_id", "session ID is required")
	}

	if err := v.ValidateZpJID(contact.ZpJID); err != nil {
		return err
	}

	if contact.Email != "" {
		if err := v.ValidateEmail(contact.Email); err != nil {
			return err
		}
	}

	if !IsValidSyncStatus(contact.SyncStatus) {
		return ErrInvalidSyncStatusValue(contact.SyncStatus)
	}

	return nil
}

func (v *DefaultValidator) ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return nil
	}

	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	if len(cleaned) < 10 || len(cleaned) > 15 {
		return ErrInvalidPhoneNumberFormat(phone)
	}

	return nil
}

func (v *DefaultValidator) ValidateEmail(email string) error {
	if email == "" {
		return nil
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmailFormat(email)
	}

	return nil
}

func (v *DefaultValidator) ValidateZpJID(jid string) error {
	if jid == "" {
		return NewContactValidationError("zp_jid", "WhatsApp JID is required")
	}

	jidPattern := `^[0-9]+@(s\.whatsapp\.net|g\.us)$`
	matched, err := regexp.MatchString(jidPattern, jid)
	if err != nil {
		return NewContactError(ErrCodeContactValidationFailed, "error validating JID", err)
	}

	if !matched {
		return NewContactValidationError("zp_jid", "invalid WhatsApp JID format")
	}

	return nil
}

func (v *DefaultValidator) ValidateSyncStatus(status string) error {
	if !IsValidSyncStatus(status) {
		return ErrInvalidSyncStatusValue(status)
	}
	return nil
}
