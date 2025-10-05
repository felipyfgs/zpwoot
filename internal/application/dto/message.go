package dto

import (
	"time"

	"zpwoot/internal/application/interfaces"
)


type SendMessageRequest struct {
	To       string       `json:"to" validate:"required"`
	Type     string       `json:"type" validate:"required,oneof=text media location contact"`
	Text     string       `json:"text,omitempty"`
	Media    *MediaData   `json:"media,omitempty"`
	Location *Location    `json:"location,omitempty"`
	Contact  *ContactInfo `json:"contact,omitempty"`
}


type MediaData struct {
	URL      string `json:"url,omitempty"`
	Base64   string `json:"base64,omitempty"`
	FileName string `json:"fileName,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	Caption  string `json:"caption,omitempty"`
}


type Location struct {
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}


type ContactInfo struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
	VCard string `json:"vcard,omitempty"`
}


type SendMessageResponse struct {
	MessageID string    `json:"messageId"`
	Status    string    `json:"status"`
	SentAt    time.Time `json:"sentAt"`
}


type MessageInfo struct {
	ID        string    `json:"id"`
	Chat      string    `json:"chat"`
	Sender    string    `json:"sender"`
	PushName  string    `json:"pushName"`
	Timestamp time.Time `json:"timestamp"`
	FromMe    bool      `json:"fromMe"`
	Type      string    `json:"type"`
	IsGroup   bool      `json:"isGroup"`
	Content   string    `json:"content,omitempty"`
}


type ReceiveMessageRequest struct {
	SessionID string      `json:"sessionId"`
	Message   MessageInfo `json:"message"`
}


type MessageHistoryRequest struct {
	ChatJID string `json:"chatJid" validate:"required"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
}


type MessageHistoryResponse struct {
	Messages []MessageInfo `json:"messages"`
	Total    int           `json:"total"`
	HasMore  bool          `json:"hasMore"`
}


func (r *SendMessageRequest) Validate() error {
	if r.To == "" {
		return ErrRecipientRequired
	}

	if r.Type == "" {
		return ErrMessageTypeRequired
	}

	switch r.Type {
	case "text":
		if r.Text == "" {
			return ErrTextRequired
		}
	case "media":
		if r.Media == nil {
			return ErrMediaRequired
		}
		if r.Media.URL == "" && r.Media.Base64 == "" {
			return ErrMediaContentRequired
		}
	case "location":
		if r.Location == nil {
			return ErrLocationRequired
		}
		if err := r.Location.Validate(); err != nil {
			return err
		}
	case "contact":
		if r.Contact == nil {
			return ErrContactRequired
		}
		if err := r.Contact.Validate(); err != nil {
			return err
		}
	default:
		return ErrInvalidMessageType
	}

	return nil
}

func (l *Location) Validate() error {
	if l.Latitude < -90 || l.Latitude > 90 {
		return ErrInvalidLatitude
	}
	if l.Longitude < -180 || l.Longitude > 180 {
		return ErrInvalidLongitude
	}
	return nil
}

func (c *ContactInfo) Validate() error {
	if c.Name == "" {
		return ErrContactNameRequired
	}
	if c.Phone == "" {
		return ErrContactPhoneRequired
	}
	return nil
}


var (
	ErrRecipientRequired    = &ErrorInfo{Code: "RECIPIENT_REQUIRED", Message: "Recipient is required"}
	ErrMessageTypeRequired  = &ErrorInfo{Code: "MESSAGE_TYPE_REQUIRED", Message: "Message type is required"}
	ErrInvalidMessageType   = &ErrorInfo{Code: "INVALID_MESSAGE_TYPE", Message: "Invalid message type"}
	ErrTextRequired         = &ErrorInfo{Code: "TEXT_REQUIRED", Message: "Text is required for text messages"}
	ErrMediaRequired        = &ErrorInfo{Code: "MEDIA_REQUIRED", Message: "Media is required for media messages"}
	ErrMediaContentRequired = &ErrorInfo{Code: "MEDIA_CONTENT_REQUIRED", Message: "Media URL or base64 content is required"}
	ErrLocationRequired     = &ErrorInfo{Code: "LOCATION_REQUIRED", Message: "Location is required for location messages"}
	ErrContactRequired      = &ErrorInfo{Code: "CONTACT_REQUIRED", Message: "Contact is required for contact messages"}
	ErrInvalidLatitude      = &ErrorInfo{Code: "INVALID_LATITUDE", Message: "Latitude must be between -90 and 90"}
	ErrInvalidLongitude     = &ErrorInfo{Code: "INVALID_LONGITUDE", Message: "Longitude must be between -180 and 180"}
	ErrContactNameRequired  = &ErrorInfo{Code: "CONTACT_NAME_REQUIRED", Message: "Contact name is required"}
	ErrContactPhoneRequired = &ErrorInfo{Code: "CONTACT_PHONE_REQUIRED", Message: "Contact phone is required"}
)


func (m *MediaData) ToInterfacesMediaData() *interfaces.MediaData {
	if m == nil {
		return nil
	}


	var data []byte


	return &interfaces.MediaData{
		MimeType: m.MimeType,
		Data:     data,
		FileName: m.FileName,
		Caption:  m.Caption,
	}
}


func (l *Location) ToInterfacesLocation() *interfaces.Location {
	if l == nil {
		return nil
	}
	return &interfaces.Location{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Name:      l.Name,
		Address:   l.Address,
	}
}


func (c *ContactInfo) ToInterfacesContactInfo() *interfaces.ContactInfo {
	if c == nil {
		return nil
	}
	return &interfaces.ContactInfo{
		Name:        c.Name,
		PhoneNumber: c.Phone,
	}
}
