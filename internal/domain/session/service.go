package session

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrSessionAlreadyExists = errors.New("session already exists")
	ErrInvalidStatus       = errors.New("invalid session status")
)

// Service defines the business logic for session management
type Service struct {
	repo Repository
}

// NewService creates a new session service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateSession creates a new session
func (s *Service) CreateSession(ctx context.Context, name string) (*Session, error) {
	if name == "" {
		return nil, errors.New("session name cannot be empty")
	}

	session := NewSession(name)
	
	if err := s.repo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// GetSession retrieves a session by ID
func (s *Service) GetSession(ctx context.Context, id string) (*Session, error) {
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return session, nil
}

// UpdateSessionStatus updates the session status
func (s *Service) UpdateSessionStatus(ctx context.Context, id string, status Status) error {
	if !status.IsValid() {
		return ErrInvalidStatus
	}

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	return nil
}

// UpdateQRCode updates the session QR code
func (s *Service) UpdateQRCode(ctx context.Context, id string, qrCode string) error {
	if err := s.repo.UpdateQRCode(ctx, id, qrCode); err != nil {
		return fmt.Errorf("failed to update QR code: %w", err)
	}

	return nil
}

// ListSessions retrieves all sessions with pagination
func (s *Service) ListSessions(ctx context.Context, limit, offset int) ([]*Session, error) {
	sessions, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	return sessions, nil
}

// DeleteSession deletes a session
func (s *Service) DeleteSession(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
