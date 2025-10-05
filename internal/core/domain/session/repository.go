package session

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, session *Session) error

	GetByID(ctx context.Context, id string) (*Session, error)

	GetByName(ctx context.Context, name string) (*Session, error)

	GetByJID(ctx context.Context, jid string) (*Session, error)

	Update(ctx context.Context, session *Session) error

	Delete(ctx context.Context, id string) error

	List(ctx context.Context, limit, offset int) ([]*Session, error)

	UpdateStatus(ctx context.Context, id string, status Status) error

	UpdateQRCode(ctx context.Context, id string, qrCode string) error
}
