package input

import (
	"context"
	"zpwoot/internal/core/application/dto"
)

type SessionCreator interface {
	Execute(ctx context.Context, req *dto.CreateSessionRequest) (*dto.CreateSessionResponse, error)
}

type SessionConnector interface {
	Execute(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error)
}

type SessionDisconnector interface {
	Execute(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error)
}

type SessionLogout interface {
	Execute(ctx context.Context, sessionID string) error
}

type SessionGetter interface {
	Execute(ctx context.Context, sessionID string) (*dto.SessionDetailResponse, error)
}

type SessionLister interface {
	Execute(ctx context.Context, req *dto.PaginationRequest) (*dto.SessionListResponse, error)
}

type SessionDeleter interface {
	Execute(ctx context.Context, sessionID string) error
}

type QRCodeManager interface {
	GetQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error)
	RefreshQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error)
}

type SessionUseCases interface {
	CreateSession(ctx context.Context, req *dto.CreateSessionRequest) (*dto.CreateSessionResponse, error)
	ConnectSession(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error)
	DisconnectSession(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error)
	LogoutSession(ctx context.Context, sessionID string) error
	DeleteSession(ctx context.Context, sessionID string) error

	GetSession(ctx context.Context, sessionID string) (*dto.SessionDetailResponse, error)
	ListSessions(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginationResponse, error)

	GetQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error)
	RefreshQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error)
}
