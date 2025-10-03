package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"zpwoot/internal/core/session"
)

type SessionResolverService struct {
	repository session.Repository
}

func NewSessionResolver(repository session.Repository) session.SessionResolver {
	return &SessionResolverService{
		repository: repository,
	}
}

func (r *SessionResolverService) ResolveToID(ctx context.Context, sessionName string) (uuid.UUID, error) {

	if id, err := uuid.Parse(sessionName); err == nil {

		_, err := r.repository.GetByID(ctx, id)
		if err != nil {
			return uuid.Nil, fmt.Errorf("session with ID %s not found: %w", sessionName, err)
		}
		return id, nil
	}

	sess, err := r.repository.GetByName(ctx, sessionName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("session with name '%s' not found: %w", sessionName, err)
	}

	return sess.ID, nil
}

func (r *SessionResolverService) ResolveToName(ctx context.Context, sessionID uuid.UUID) (string, error) {
	sess, err := r.repository.GetByID(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("session with ID %s not found: %w", sessionID.String(), err)
	}

	return sess.Name, nil
}

func (r *SessionResolverService) Resolve(ctx context.Context, sessionName string) (*session.ResolveResult, error) {

	if id, err := uuid.Parse(sessionName); err == nil {

		sess, err := r.repository.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("session with ID %s not found: %w", sessionName, err)
		}

		return &session.ResolveResult{
			ID:      sess.ID,
			Name:    sess.Name,
			Session: sess,
		}, nil
	}

	sess, err := r.repository.GetByName(ctx, sessionName)
	if err != nil {
		return nil, fmt.Errorf("session with name '%s' not found: %w", sessionName, err)
	}

	return &session.ResolveResult{
		ID:      sess.ID,
		Name:    sess.Name,
		Session: sess,
	}, nil
}
