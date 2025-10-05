package repository

import (
	"context"

	"zpwoot/internal/core/domain/session"
)

type SessionRepo struct {
	repo *SessionRepository
}

func NewSessionRepo(repo *SessionRepository) *SessionRepo {
	return &SessionRepo{
		repo: repo,
	}
}

func (r *SessionRepo) GetByID(ctx context.Context, sessionID string) (*session.Session, error) {
	return r.repo.GetByID(ctx, sessionID)
}

func (r *SessionRepo) GetByName(ctx context.Context, name string) (*session.Session, error) {
	return r.repo.GetByName(ctx, name)
}

func (r *SessionRepo) Create(ctx context.Context, sess *session.Session) error {
	return r.repo.Create(ctx, sess)
}

func (r *SessionRepo) Update(ctx context.Context, sess *session.Session) error {
	return r.repo.Update(ctx, sess)
}

func (r *SessionRepo) Delete(ctx context.Context, sessionID string) error {
	return r.repo.Delete(ctx, sessionID)
}

func (r *SessionRepo) List(ctx context.Context, limit, offset int) ([]*session.Session, error) {
	return r.repo.List(ctx, limit, offset)
}
