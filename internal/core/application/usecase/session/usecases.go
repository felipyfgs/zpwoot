package session

import (
	"zpwoot/internal/core/application/interfaces"
	"zpwoot/internal/core/domain/session"
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
	whatsappClient interfaces.WhatsAppClient,
	notificationSvc interfaces.NotificationService,
) *UseCases {
	return &UseCases{
		Create:     NewCreateUseCase(sessionService, whatsappClient, notificationSvc),
		Connect:    NewConnectUseCase(sessionService, whatsappClient, notificationSvc),
		Disconnect: NewDisconnectUseCase(sessionService, whatsappClient, notificationSvc),
		Logout:     NewLogoutUseCase(sessionService, whatsappClient, notificationSvc),
		Get:        NewGetUseCase(sessionService, whatsappClient),
		List:       NewListUseCase(sessionService, whatsappClient),
		Delete:     NewDeleteUseCase(sessionService, whatsappClient, notificationSvc),
		QR:         NewQRUseCase(sessionService, whatsappClient, notificationSvc),
	}
}

