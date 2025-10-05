package waclient

import (
	"context"

	"zpwoot/internal/application/interfaces"
)


type WhatsAppAdapter struct {
	client *WAClient
}


func NewWhatsAppAdapter(client *WAClient) *WhatsAppAdapter {
	return &WhatsAppAdapter{
		client: client,
	}
}


func (a *WhatsAppAdapter) CreateSession(ctx context.Context, sessionID string) error {
	config := &SessionConfig{
		SessionID: sessionID,
	}
	_, err := a.client.CreateSession(ctx, config)
	return err
}


func (a *WhatsAppAdapter) GetSessionStatus(ctx context.Context, sessionID string) (*interfaces.SessionStatus, error) {
	client, err := a.client.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	deviceJID := ""
	if client.WAClient != nil && client.WAClient.Store.ID != nil {
		deviceJID = client.WAClient.Store.ID.String()
	}

	return &interfaces.SessionStatus{
		SessionID:   client.SessionID,
		Connected:   client.Status == "connected",
		LoggedIn:    deviceJID != "",
		DeviceJID:   deviceJID,
		ConnectedAt: client.ConnectedAt,
		LastSeen:    client.LastSeen,
	}, nil
}


func (a *WhatsAppAdapter) DeleteSession(ctx context.Context, sessionID string) error {
	return a.client.DeleteSession(ctx, sessionID)
}


func (a *WhatsAppAdapter) ConnectSession(ctx context.Context, sessionID string) error {
	return a.client.ConnectSession(ctx, sessionID)
}


func (a *WhatsAppAdapter) DisconnectSession(ctx context.Context, sessionID string) error {
	return a.client.DisconnectSession(ctx, sessionID)
}


func (a *WhatsAppAdapter) LogoutSession(ctx context.Context, sessionID string) error {
	return a.client.LogoutSession(ctx, sessionID)
}


func (a *WhatsAppAdapter) IsConnected(ctx context.Context, sessionID string) bool {
	client, err := a.client.GetSession(ctx, sessionID)
	if err != nil {
		return false
	}
	return client.Status == "connected"
}


func (a *WhatsAppAdapter) IsLoggedIn(ctx context.Context, sessionID string) bool {
	client, err := a.client.GetSession(ctx, sessionID)
	if err != nil {
		return false
	}
	return client.WAClient != nil && client.WAClient.Store.ID != nil
}


func (a *WhatsAppAdapter) GetQRCode(ctx context.Context, sessionID string) (*interfaces.QRCodeInfo, error) {
	qrEvent, err := a.client.GetQRCodeForSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &interfaces.QRCodeInfo{
		Code:      qrEvent.Code,
		Base64:    qrEvent.Base64,
		ExpiresAt: qrEvent.ExpiresAt,
	}, nil
}


func (a *WhatsAppAdapter) SendTextMessage(ctx context.Context, sessionID, to, text string) (*interfaces.MessageResult, error) {

	return nil, &interfaces.WhatsAppError{
		Code:    "NOT_IMPLEMENTED",
		Message: "SendTextMessage not implemented yet",
	}
}


func (a *WhatsAppAdapter) SendMediaMessage(ctx context.Context, sessionID, to string, media *interfaces.MediaData) (*interfaces.MessageResult, error) {

	return nil, &interfaces.WhatsAppError{
		Code:    "NOT_IMPLEMENTED",
		Message: "SendMediaMessage not implemented yet",
	}
}


func (a *WhatsAppAdapter) SendLocationMessage(ctx context.Context, sessionID, to string, location *interfaces.Location) (*interfaces.MessageResult, error) {

	return nil, &interfaces.WhatsAppError{
		Code:    "NOT_IMPLEMENTED",
		Message: "SendLocationMessage not implemented yet",
	}
}


func (a *WhatsAppAdapter) SendContactMessage(ctx context.Context, sessionID, to string, contact *interfaces.ContactInfo) (*interfaces.MessageResult, error) {

	return nil, &interfaces.WhatsAppError{
		Code:    "NOT_IMPLEMENTED",
		Message: "SendContactMessage not implemented yet",
	}
}

