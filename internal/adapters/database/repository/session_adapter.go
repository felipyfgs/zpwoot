package repository

import (
	"context"

	"zpwoot/internal/core/domain/session"
)


type SessionRepositoryAdapter struct {
	repo *SessionRepository
}


func NewSessionRepositoryAdapter(repo *SessionRepository) *SessionRepositoryAdapter {
	return &SessionRepositoryAdapter{
		repo: repo,
	}
}


func (a *SessionRepositoryAdapter) GetSession(ctx context.Context, sessionID string) (*session.Session, error) {
	return a.repo.GetByID(ctx, sessionID)
}


func (a *SessionRepositoryAdapter) GetSessionByName(ctx context.Context, name string) (*session.Session, error) {
	return a.repo.GetByName(ctx, name)
}


func (a *SessionRepositoryAdapter) CreateSession(ctx context.Context, sess *session.Session) error {
	return a.repo.Create(ctx, sess)
}


func (a *SessionRepositoryAdapter) UpdateSession(ctx context.Context, sess *session.Session) error {
	return a.repo.Update(ctx, sess)
}


func (a *SessionRepositoryAdapter) DeleteSession(ctx context.Context, sessionID string) error {
	return a.repo.Delete(ctx, sessionID)
}


func (a *SessionRepositoryAdapter) ListSessions(ctx context.Context) ([]*session.Session, error) {

	return a.repo.List(ctx, 0, 0)
}

