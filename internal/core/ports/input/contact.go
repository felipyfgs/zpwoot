package input

import (
	"context"
)

type ContactService interface {
	CheckUser(ctx context.Context, sessionID string, phones []string) ([]UserCheckResult, error)
	GetUser(ctx context.Context, sessionID string, phone string) (*UserDetail, error)
	GetAvatar(ctx context.Context, sessionID string, phone string, preview bool) (*AvatarInfo, error)
	GetContacts(ctx context.Context, sessionID string) ([]Contact, error)
	SendPresence(ctx context.Context, sessionID string, presence string) error
	ChatPresence(ctx context.Context, sessionID string, phone string, presence string, media string) error
}
type UserCheckResult struct {
	Query        string
	IsInWhatsApp bool
	JID          string
	VerifiedName string
}
type UserDetail struct {
	JID          string
	VerifiedName string
	Status       string
	PictureID    string
}
type AvatarInfo struct {
	URL       string
	ID        string
	Type      string
	DirectURL string
}
type Contact struct {
	JID          string
	Name         string
	Notify       string
	VerifiedName string
	BusinessName string
}
