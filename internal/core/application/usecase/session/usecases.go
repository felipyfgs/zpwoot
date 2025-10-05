package session

import (
	"zpwoot/internal/core/domain/session"
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
}

func NewUseCases(
	sessionService *session.Service,
	whatsappClient output.WhatsAppClient,
) *UseCases {
	return &UseCases{
		Create:     NewCreateUseCase(sessionService, whatsappClient),
		Connect:    NewConnectUseCase(sessionService, whatsappClient),
		Disconnect: NewDisconnectUseCase(sessionService, whatsappClient),
		Logout:     NewLogoutUseCase(sessionService, whatsappClient),
		Get:        NewGetUseCase(sessionService, whatsappClient),
		List:       NewListUseCase(sessionService, whatsappClient),
		Delete:     NewDeleteUseCase(sessionService, whatsappClient),
		QR:         NewQRUseCase(sessionService, whatsappClient),
	}
}
