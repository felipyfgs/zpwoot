package session

import (
	"context"

	"zpwoot/internal/core/application/dto"
	"zpwoot/internal/core/domain/session"
	"zpwoot/internal/core/ports/input"
	"zpwoot/internal/core/ports/output"
)

type UseCases struct {
	Create     *CreateUseCase
	Connect    *ConnectUseCase
	Disconnect *DisconnectUseCase
	Logout     *LogoutUseCase
	Get        *GetUseCase
	List       *ListUseCase
	Delete     *DeleteUseCase
	QR         *QRUseCase
	Pair       *PairUseCase
}

func NewUseCases(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
	logger output.Logger,
) *UseCases {
	return &UseCases{
		Create:     NewCreateUseCase(sessionService, whatsappClient, logger),
		Connect:    NewConnectUseCase(sessionService, whatsappClient, logger),
		Disconnect: NewDisconnectUseCase(sessionService, whatsappClient, logger),
		Logout:     NewLogoutUseCase(sessionService, whatsappClient, logger),
		Get:        NewGetUseCase(sessionService, whatsappClient, logger),
		List:       NewListUseCase(sessionService, whatsappClient, logger),
		Delete:     NewDeleteUseCase(sessionService, whatsappClient, logger),
		QR:         NewQRUseCase(sessionService, whatsappClient, logger),
		Pair:       NewPairUseCase(whatsappClient, logger),
	}
}

func (uc *UseCases) CreateSession(ctx context.Context, req *dto.CreateRequest) (*dto.CreateSessionResponse, error) {
	return uc.Create.Execute(ctx, req)
}

func (uc *UseCases) ConnectSession(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
	return uc.Connect.Execute(ctx, sessionID)
}

func (uc *UseCases) DisconnectSession(ctx context.Context, sessionID string) (*dto.SessionStatusResponse, error) {
	return uc.Disconnect.Execute(ctx, sessionID)
}

func (uc *UseCases) LogoutSession(ctx context.Context, sessionID string) error {
	return uc.Logout.Execute(ctx, sessionID)
}

func (uc *UseCases) DeleteSession(ctx context.Context, sessionID string) error {
	return uc.Delete.Execute(ctx, sessionID)
}

func (uc *UseCases) GetSession(ctx context.Context, sessionID string) (*dto.SessionDetailResponse, error) {
	return uc.Get.Execute(ctx, sessionID)
}

func (uc *UseCases) ListSessions(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	return uc.List.Execute(ctx, req)
}

func (uc *UseCases) GetQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {
	return uc.QR.GetQRCode(ctx, sessionID)
}

func (uc *UseCases) RefreshQRCode(ctx context.Context, sessionID string) (*dto.QRCodeResponse, error) {
	return uc.QR.RefreshQRCode(ctx, sessionID)
}

func (uc *UseCases) PairPhone(ctx context.Context, sessionID string, phone string) (*dto.PairPhoneResponse, error) {
	return uc.Pair.Execute(ctx, sessionID, phone)
}

var _ input.SessionUseCases = (*UseCases)(nil)
