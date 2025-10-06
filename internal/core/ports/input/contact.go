package input

import (
	"context"
)

// ContactService define operações de contatos e presença
type ContactService interface {
	CheckUser(ctx context.Context, sessionID string, phones []string) ([]UserCheckResult, error)
	GetUser(ctx context.Context, sessionID string, phone string) (*UserDetail, error)
	GetAvatar(ctx context.Context, sessionID string, phone string, preview bool) (*AvatarInfo, error)
	GetContacts(ctx context.Context, sessionID string) ([]Contact, error)
	SendPresence(ctx context.Context, sessionID string, presence string) error
	ChatPresence(ctx context.Context, sessionID string, phone string, presence string, media string) error
}

// UserCheckResult representa o resultado da verificação de usuário
type UserCheckResult struct {
	Query        string
	IsInWhatsApp bool
	JID          string
	VerifiedName string
}

// UserDetail representa detalhes de um usuário
type UserDetail struct {
	JID          string
	VerifiedName string
	Status       string
	PictureID    string
}

// AvatarInfo representa informações de avatar
type AvatarInfo struct {
	URL       string
	ID        string
	Type      string
	DirectURL string
}

// Contact representa um contato
type Contact struct {
	JID          string
	Name         string
	Notify       string
	VerifiedName string
	BusinessName string
}
