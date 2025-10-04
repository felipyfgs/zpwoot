package contact

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, contact *Contact) error
	GetByID(ctx context.Context, id uuid.UUID) (*Contact, error)
	Update(ctx context.Context, contact *Contact) error
	Delete(ctx context.Context, id uuid.UUID) error

	GetByZpJID(ctx context.Context, sessionID uuid.UUID, zpJID string) (*Contact, error)
	ExistsByZpJID(ctx context.Context, sessionID uuid.UUID, zpJID string) (bool, error)

	GetByCwContactID(ctx context.Context, cwContactID int) (*Contact, error)
	GetByCwConversationID(ctx context.Context, cwConversationID int) (*Contact, error)

	List(ctx context.Context, req *ListContactsRequest) ([]*Contact, int64, error)
	ListBySession(ctx context.Context, sessionID uuid.UUID, limit, offset int) ([]*Contact, error)
	ListBySyncStatus(ctx context.Context, status SyncStatus, limit, offset int) ([]*Contact, error)

	UpdateSyncStatus(ctx context.Context, id uuid.UUID, status SyncStatus, cwContactID, cwConversationID *int) error
	GetPendingSyncContacts(ctx context.Context, sessionID uuid.UUID, limit int) ([]*Contact, error)
	GetFailedSyncContacts(ctx context.Context, sessionID uuid.UUID, limit int) ([]*Contact, error)
	MarkAsSynced(ctx context.Context, id uuid.UUID, cwContactID, cwConversationID int) error
	MarkAsFailed(ctx context.Context, id uuid.UUID, errorReason string) error

	Count(ctx context.Context) (int64, error)
	CountBySession(ctx context.Context, sessionID uuid.UUID) (int64, error)
	CountBySyncStatus(ctx context.Context, status SyncStatus) (int64, error)
	CountByType(ctx context.Context, contactType ContactType) (int64, error)

	GetStats(ctx context.Context) (*ContactStats, error)
	GetStatsBySession(ctx context.Context, sessionID uuid.UUID) (*ContactStats, error)

	DeleteOldContacts(ctx context.Context, olderThanDays int) (int64, error)
	DeleteBySession(ctx context.Context, sessionID uuid.UUID) (int64, error)
	CleanupFailedContacts(ctx context.Context, olderThanHours int) (int64, error)
}

type ContactService interface {
	CreateContact(ctx context.Context, req *CreateContactRequest) (*Contact, error)
	GetContact(ctx context.Context, id uuid.UUID) (*Contact, error)
	UpdateContact(ctx context.Context, req *UpdateContactRequest) (*Contact, error)
	DeleteContact(ctx context.Context, id uuid.UUID) error

	SyncContact(ctx context.Context, id uuid.UUID, cwContactID, cwConversationID int) error
	SyncPendingContacts(ctx context.Context, sessionID uuid.UUID, limit int) ([]*Contact, error)
	RetryFailedContacts(ctx context.Context, sessionID uuid.UUID, limit int) ([]*Contact, error)

	ListContacts(ctx context.Context, req *ListContactsRequest) ([]*Contact, int64, error)
	GetContactsBySession(ctx context.Context, sessionID uuid.UUID, limit, offset int) ([]*Contact, error)
	SearchContacts(ctx context.Context, sessionID uuid.UUID, query string, limit, offset int) ([]*Contact, error)

	GetContactStats(ctx context.Context, sessionID *uuid.UUID) (*ContactStats, error)
	GetDashboardStats(ctx context.Context) (*ContactStats, error)

	CleanupOldContacts(ctx context.Context, sessionID uuid.UUID, olderThanDays int) (int64, error)
	ReprocessFailedContacts(ctx context.Context, sessionID uuid.UUID) (int64, error)
}

type WhatsAppGateway interface {
	GetContactInfo(ctx context.Context, sessionID uuid.UUID, jid string) (*WhatsAppContact, error)
	GetContactList(ctx context.Context, sessionID uuid.UUID) ([]*WhatsAppContact, error)
	GetContactPresence(ctx context.Context, sessionID uuid.UUID, jid string) (*ContactPresence, error)

	UpdateContactName(ctx context.Context, sessionID uuid.UUID, jid, name string) error
	BlockContact(ctx context.Context, sessionID uuid.UUID, jid string) error
	UnblockContact(ctx context.Context, sessionID uuid.UUID, jid string) error

	SyncToChatwoot(ctx context.Context, contact *Contact) error
	SyncFromChatwoot(ctx context.Context, cwContactID int) (*Contact, error)
}

type EventPublisher interface {
	PublishContactCreated(ctx context.Context, contact *Contact) error
	PublishContactUpdated(ctx context.Context, contact *Contact) error
	PublishContactSynced(ctx context.Context, contact *Contact) error
	PublishContactSyncFailed(ctx context.Context, contact *Contact, reason string) error
	PublishContactDeleted(ctx context.Context, contactID uuid.UUID) error

	PublishSyncStarted(ctx context.Context, sessionID uuid.UUID, contactCount int) error
	PublishSyncCompleted(ctx context.Context, sessionID uuid.UUID, syncedCount, failedCount int) error
	PublishSyncFailed(ctx context.Context, sessionID uuid.UUID, reason string) error
}

type ContactValidator interface {
	ValidateCreateRequest(req *CreateContactRequest) error
	ValidateUpdateRequest(req *UpdateContactRequest) error
	ValidateContact(contact *Contact) error
	ValidatePhoneNumber(phone string) error
	ValidateEmail(email string) error
	ValidateZpJID(jid string) error
	ValidateSyncStatus(status string) error
}

type WhatsAppContact struct {
	JID        string `json:"jid"`
	Name       string `json:"name"`
	PushName   string `json:"push_name"`
	ShortName  string `json:"short_name"`
	Avatar     string `json:"avatar"`
	IsGroup    bool   `json:"is_group"`
	IsBusiness bool   `json:"is_business"`
	IsBlocked  bool   `json:"is_blocked"`
}

type ContactPresence struct {
	JID        string `json:"jid"`
	IsOnline   bool   `json:"is_online"`
	LastSeen   string `json:"last_seen"`
	LastStatus string `json:"last_status"`
}
