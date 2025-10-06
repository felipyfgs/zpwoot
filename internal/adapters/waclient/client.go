package waclient

import (
	"context"
	"crypto/rand"
	"time"

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

func (w *WAClientAdapter) GetWAClient() *WAClient {
	return w.client
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
		return nil, w.convertError(err)
	}

	return &output.SessionStatus{
		SessionID:   client.SessionID,
		Connected:   client.IsConnected(),
		LoggedIn:    client.IsLoggedIn(),
		DeviceJID:   client.GetDeviceJID(),
		ConnectedAt: client.ConnectedAt,
		LastSeen:    client.LastSeen,
	}, nil
}

func (w *WAClientAdapter) DeleteSession(ctx context.Context, sessionID string) error {
	return w.convertError(w.client.DeleteSession(ctx, sessionID))
}

func (w *WAClientAdapter) ConnectSession(ctx context.Context, sessionID string) error {
	return w.convertError(w.client.ConnectSession(ctx, sessionID))
}

func (w *WAClientAdapter) DisconnectSession(ctx context.Context, sessionID string) error {
	return w.convertError(w.client.DisconnectSession(ctx, sessionID))
}

func (w *WAClientAdapter) LogoutSession(ctx context.Context, sessionID string) error {
	return w.convertError(w.client.LogoutSession(ctx, sessionID))
}

func (w *WAClientAdapter) IsConnected(ctx context.Context, sessionID string) bool {
	client, err := w.client.GetSession(ctx, sessionID)
	if err != nil {
		return false
	}
	return client.IsConnected()
}

func (w *WAClientAdapter) IsLoggedIn(ctx context.Context, sessionID string) bool {
	client, err := w.client.GetSession(ctx, sessionID)
	if err != nil {
		return false
	}
	return client.IsLoggedIn()
}

func (w *WAClientAdapter) GetQRCode(ctx context.Context, sessionID string) (*output.QRCodeInfo, error) {
	qrEvent, err := w.client.GetQRCodeForSession(ctx, sessionID)
	if err != nil {
		return nil, w.convertError(err)
	}

	return &output.QRCodeInfo{
		Code:      qrEvent.Code,
		Base64:    qrEvent.Base64,
		ExpiresAt: qrEvent.ExpiresAt,
	}, nil
}

func (w *WAClientAdapter) SendTextMessage(ctx context.Context, sessionID, to, text string) (*output.MessageResult, error) {
	messageSender := NewMessageSender(w.client)
	resp, err := messageSender.SendTextMessage(ctx, sessionID, to, text, nil)
	if err != nil {
		return nil, w.convertError(err)
	}

	return &output.MessageResult{
		MessageID: string(resp.ID),
		Status:    "sent",
		SentAt:    resp.Timestamp,
	}, nil
}

func (w *WAClientAdapter) SendMediaMessage(ctx context.Context, sessionID, to string, media *output.MediaData) (*output.MessageResult, error) {
	messageSender := NewMessageSender(w.client)
	resp, err := messageSender.SendMediaMessage(ctx, sessionID, to, media)
	if err != nil {
		return nil, w.convertError(err)
	}

	return &output.MessageResult{
		MessageID: string(resp.ID),
		Status:    "sent",
		SentAt:    resp.Timestamp,
	}, nil
}

func (w *WAClientAdapter) SendLocationMessage(ctx context.Context, sessionID, to string, location *output.Location) (*output.MessageResult, error) {
	messageSender := NewMessageSender(w.client)
	err := messageSender.SendLocationMessage(ctx, sessionID, to, location.Latitude, location.Longitude, location.Name)
	if err != nil {
		return nil, w.convertError(err)
	}

	return &output.MessageResult{
		MessageID: generateMessageID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}, nil
}

func (w *WAClientAdapter) SendContactMessage(ctx context.Context, sessionID, to string, contact *output.ContactInfo) (*output.MessageResult, error) {
	messageSender := NewMessageSender(w.client)
	contactInfo := &ContactInfo{
		Name:  contact.Name,
		Phone: contact.PhoneNumber,
	}
	err := messageSender.SendContactMessage(ctx, sessionID, to, contactInfo)
	if err != nil {
		return nil, w.convertError(err)
	}

	return &output.MessageResult{
		MessageID: generateMessageID(),
		Status:    "sent",
		SentAt:    time.Now(),
	}, nil
}

func (w *WAClientAdapter) convertError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*output.WhatsAppError); ok {
		return err
	}

	switch err {
	case ErrSessionNotFound:
		return output.ErrSessionNotFound
	case ErrNotConnected:
		return output.ErrSessionNotConnected
	case ErrInvalidJID:
		return output.ErrInvalidJID
	case ErrConnectionFailed:
		return output.ErrConnectionFailed
	default:

		return &output.WhatsAppError{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}
	}
}

func generateMessageID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {

		for i := range b {
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		}
		return string(b)
	}

	for i := range b {
		b[i] = charset[randomBytes[i]%byte(len(charset))]
	}
	return string(b)
}

type ContactInfo struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	VCard string `json:"vcard,omitempty"`
}

type SessionManagerAdapter struct {
	client *WAClient
}

func NewSessionManagerAdapter(client *WAClient) *SessionManagerAdapter {
	return &SessionManagerAdapter{
		client: client,
	}
}

func (s *SessionManagerAdapter) CreateSession(ctx context.Context, sessionID string) error {
	config := &SessionConfig{
		SessionID: sessionID,
	}
	_, err := s.client.CreateSession(ctx, config)
	return err
}

func (s *SessionManagerAdapter) GetSessionStatus(ctx context.Context, sessionID string) (*output.SessionStatus, error) {
	client, err := s.client.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &output.SessionStatus{
		SessionID:   client.SessionID,
		Connected:   client.IsConnected(),
		LoggedIn:    client.IsLoggedIn(),
		DeviceJID:   client.GetDeviceJID(),
		ConnectedAt: client.ConnectedAt,
		LastSeen:    client.LastSeen,
	}, nil
}

func (s *SessionManagerAdapter) DeleteSession(ctx context.Context, sessionID string) error {
	return s.client.DeleteSession(ctx, sessionID)
}

func (s *SessionManagerAdapter) ConnectSession(ctx context.Context, sessionID string) error {
	return s.client.ConnectSession(ctx, sessionID)
}

func (s *SessionManagerAdapter) DisconnectSession(ctx context.Context, sessionID string) error {
	return s.client.DisconnectSession(ctx, sessionID)
}

func (s *SessionManagerAdapter) LogoutSession(ctx context.Context, sessionID string) error {
	return s.client.LogoutSession(ctx, sessionID)
}

func (s *SessionManagerAdapter) IsConnected(ctx context.Context, sessionID string) bool {
	client, err := s.client.GetSession(ctx, sessionID)
	if err != nil {
		return false
	}
	return client.IsConnected()
}

func (s *SessionManagerAdapter) IsLoggedIn(ctx context.Context, sessionID string) bool {
	client, err := s.client.GetSession(ctx, sessionID)
	if err != nil {
		return false
	}
	return client.IsLoggedIn()
}

func (s *SessionManagerAdapter) GetQRCode(ctx context.Context, sessionID string) (*output.QRCodeInfo, error) {
	qrEvent, err := s.client.GetQRCodeForSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &output.QRCodeInfo{
		Code:      qrEvent.Code,
		Base64:    qrEvent.Base64,
		ExpiresAt: qrEvent.ExpiresAt,
	}, nil
}
