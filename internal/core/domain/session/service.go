package session

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"zpwoot/internal/core/domain/shared"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, name string) (*Session, error) {
	if name == "" {
		return nil, errors.New("session name cannot be empty")
	}

	existingSession, err := s.repo.GetByName(ctx, name)
	if err != nil && !errors.Is(err, shared.ErrSessionNotFound) {
		return nil, fmt.Errorf("failed to check existing session: %w", err)
	}

	if existingSession != nil {
		return nil, shared.ErrSessionAlreadyExists
	}

	session := NewSession(name)

	if err := s.repo.Create(ctx, session); err != nil {
		if isUniqueConstraintError(err) {
			return nil, shared.ErrSessionAlreadyExists
		}

		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Session, error) {
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (s *Service) Update(ctx context.Context, session *Session) error {
	if err := s.repo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status Status) error {
	if !status.IsValid() {
		return shared.ErrInvalidStatus
	}

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	return nil
}

func (s *Service) UpdateQR(ctx context.Context, id string, qrCode string) error {
	if err := s.repo.UpdateQRCode(ctx, id, qrCode); err != nil {
		return fmt.Errorf("failed to update QR code: %w", err)
	}

	return nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*Session, error) {
	sessions, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	return sessions, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	return strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "violates unique")
}
