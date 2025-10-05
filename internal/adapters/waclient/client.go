package waclient

import (
	"context"

	"zpwoot/internal/core/ports/output"
)

type WAClientAdapter struct {
	client *WAClient
}

func NewWAClientAdapter(client *WAClient) *WAClientAdapter {
	return &WAClientAdapter{
		client: client,
	}
}

func (w *WAClientAdapter) CreateSession(ctx context.Context, sessionID string) error {
	config := &SessionConfig{
		SessionID: sessionID,
	}
	_, err := w.client.CreateSession(ctx, config)
	return err
}

func (w *WAClientAdapter) GetSessionStatus(ctx context.Context, sessionID string) (*output.SessionStatus, error) {
	client, err := w.client.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	deviceJID := ""
	if client.WAClient != nil && client.WAClient.Store.ID != nil {
		deviceJID = client.WAClient.Store.ID.String()
	}

	return &output.SessionStatus{
		SessionID:   client.SessionID,
		Connected:   client.Status == "connected",
		LoggedIn:    deviceJID != "",
		DeviceJID:   deviceJID,
		ConnectedAt: client.ConnectedAt,
		LastSeen:    client.LastSeen,
	}, nil
}

func (w *WAClientAdapter) DeleteSession(ctx context.Context, sessionID string) error {
	return w.client.DeleteSession(ctx, sessionID)
}

func (w *WAClientAdapter) ConnectSession(ctx context.Context, sessionID string) error {
	return w.client.ConnectSession(ctx, sessionID)
}

func (w *WAClientAdapter) DisconnectSession(ctx context.Context, sessionID string) error {
	return w.client.DisconnectSession(ctx, sessionID)
}

func (w *WAClientAdapter) LogoutSession(ctx context.Context, sessionID string) error {
	return w.client.LogoutSession(ctx, sessionID)
}

func (w *WAClientAdapter) IsConnected(ctx context.Context, sessionID string) bool {
	client, err := w.client.GetSession(ctx, sessionID)
	if err != nil {
		return false
	}
	return client.Status == "connected"
}

func (w *WAClientAdapter) IsLoggedIn(ctx context.Context, sessionID string) bool {
	client, err := w.client.GetSession(ctx, sessionID)
	if err != nil {
		return false
	}
	return client.WAClient != nil && client.WAClient.Store.ID != nil
}

func (w *WAClientAdapter) GetQRCode(ctx context.Context, sessionID string) (*output.QRCodeInfo, error) {
	qrEvent, err := w.client.GetQRCodeForSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &output.QRCodeInfo{
		Code:      qrEvent.Code,
		Base64:    qrEvent.Base64,
		ExpiresAt: qrEvent.ExpiresAt,
	}, nil
}

func (w *WAClientAdapter) SendTextMessage(ctx context.Context, sessionID, to, text string) (*output.MessageResult, error) {

	return nil, &output.WhatsAppError{
		Code:    "NOT_IMPLEMENTED",
		Message: "SendTextMessage not implemented yet",
	}
}

func (w *WAClientAdapter) SendMediaMessage(ctx context.Context, sessionID, to string, media *output.MediaData) (*output.MessageResult, error) {

	return nil, &output.WhatsAppError{
		Code:    "NOT_IMPLEMENTED",
		Message: "SendMediaMessage not implemented yet",
	}
}

func (w *WAClientAdapter) SendLocationMessage(ctx context.Context, sessionID, to string, location *output.Location) (*output.MessageResult, error) {

	return nil, &output.WhatsAppError{
		Code:    "NOT_IMPLEMENTED",
		Message: "SendLocationMessage not implemented yet",
	}
}

func (w *WAClientAdapter) SendContactMessage(ctx context.Context, sessionID, to string, contact *output.ContactInfo) (*output.MessageResult, error) {

	return nil, &output.WhatsAppError{
		Code:    "NOT_IMPLEMENTED",
		Message: "SendContactMessage not implemented yet",
	}
}
