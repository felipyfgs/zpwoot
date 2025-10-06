package dto

import (
	"time"

	"zpwoot/internal/core/application/validators"
	"zpwoot/internal/core/ports/output"
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

type SendTextMessageRequest struct {
	To   string `json:"to" validate:"required"`
	Text string `json:"text" validate:"required"`
}

type SendImageMessageRequest struct {
	To      string     `json:"to" validate:"required"`
	Image   *MediaData `json:"image" validate:"required"`
	Caption string     `json:"caption,omitempty"`
}

type SendAudioMessageRequest struct {
	To    string     `json:"to" validate:"required"`
	Audio *MediaData `json:"audio" validate:"required"`
}

type SendVideoMessageRequest struct {
	To      string     `json:"to" validate:"required"`
	Video   *MediaData `json:"video" validate:"required"`
	Caption string     `json:"caption,omitempty"`
}

type SendDocumentMessageRequest struct {
	To       string     `json:"to" validate:"required"`
	Document *MediaData `json:"document" validate:"required"`
	Caption  string     `json:"caption,omitempty"`
}

type SendStickerMessageRequest struct {
	To      string     `json:"to" validate:"required"`
	Sticker *MediaData `json:"sticker" validate:"required"`
}

type SendLocationMessageRequest struct {
	To        string  `json:"to" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

type SendContactMessageRequest struct {
	To      string       `json:"to" validate:"required"`
	Contact *ContactInfo `json:"contact" validate:"required"`
}

type SendReactionMessageRequest struct {
	To        string `json:"to" validate:"required"`
	MessageID string `json:"messageId" validate:"required"`
	Reaction  string `json:"reaction" validate:"required"`
}

type SendPollMessageRequest struct {
	To                     string   `json:"to" validate:"required"`
	Name                   string   `json:"name" validate:"required"`
	Options                []string `json:"options" validate:"required,min=2,max=12"`
	SelectableOptionsCount int      `json:"selectableOptionsCount" validate:"required,min=1"`
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
		if err := validators.ValidateMessageText(r.Text); err != nil {
			return &ErrorInfo{Code: "TEXT_INVALID", Message: err.Error()}
		}
	case "media":
		if r.Media == nil {
			return ErrMediaRequired
		}
		if r.Media.URL == "" && r.Media.Base64 == "" {
			return ErrMediaContentRequired
		}
		if err := validators.ValidateCaption(r.Media.Caption); err != nil {
			return &ErrorInfo{Code: "CAPTION_INVALID", Message: err.Error()}
		}
		if err := validators.ValidateFileName(r.Media.FileName); err != nil {
			return &ErrorInfo{Code: "FILENAME_INVALID", Message: err.Error()}
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
	if err := validators.ValidateLatitude(l.Latitude); err != nil {
		return &ErrorInfo{Code: "LATITUDE_INVALID", Message: err.Error()}
	}
	if err := validators.ValidateLongitude(l.Longitude); err != nil {
		return &ErrorInfo{Code: "LONGITUDE_INVALID", Message: err.Error()}
	}
	if err := validators.ValidateLocationName(l.Name); err != nil {
		return &ErrorInfo{Code: "LOCATION_NAME_INVALID", Message: err.Error()}
	}
	if err := validators.ValidateLocationAddress(l.Address); err != nil {
		return &ErrorInfo{Code: "LOCATION_ADDRESS_INVALID", Message: err.Error()}
	}
	return nil
}

func (c *ContactInfo) Validate() error {
	if err := validators.ValidateContactName(c.Name); err != nil {
		return &ErrorInfo{Code: "CONTACT_NAME_INVALID", Message: err.Error()}
	}
	if err := validators.ValidatePhoneNumber(c.Phone); err != nil {
		return &ErrorInfo{Code: "CONTACT_PHONE_INVALID", Message: err.Error()}
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

func (m *MediaData) ToInterfacesMediaData() *output.MediaData {
	if m == nil {
		return nil
	}

	var data []byte

	return &output.MediaData{
		MimeType: m.MimeType,
		Data:     data,
		FileName: m.FileName,
		Caption:  m.Caption,
	}
}

func (m *MediaData) ToMediaData() *output.MediaData {
	return m.ToInterfacesMediaData()
}

func (m *MediaData) ToOutputMediaData() *output.MediaData {
	return m.ToInterfacesMediaData()
}

func (l *Location) ToInterfacesLocation() *output.Location {
	if l == nil {
		return nil
	}
	return &output.Location{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Name:      l.Name,
		Address:   l.Address,
	}
}

func (c *ContactInfo) ToInterfacesContactInfo() *output.ContactInfo {
	if c == nil {
		return nil
	}
	return &output.ContactInfo{
		Name:        c.Name,
		PhoneNumber: c.Phone,
	}
}

type SendContactsArrayMessageRequest struct {
	To       string         `json:"to" validate:"required"`
	Contacts []*ContactInfo `json:"contacts" validate:"required,min=1"`
}

type Button struct {
	ID   string `json:"id" validate:"required"`
	Text string `json:"text" validate:"required"`
}

type SendButtonsMessageRequest struct {
	To      string    `json:"to" validate:"required"`
	Text    string    `json:"text" validate:"required"`
	Buttons []*Button `json:"buttons" validate:"required,min=1,max=3"`
}

type ListRow struct {
	ID          string `json:"id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description,omitempty"`
}

type ListSection struct {
	Title string     `json:"title" validate:"required"`
	Rows  []*ListRow `json:"rows" validate:"required,min=1"`
}

type SendListMessageRequest struct {
	To       string         `json:"to" validate:"required"`
	Text     string         `json:"text" validate:"required"`
	Title    string         `json:"title" validate:"required"`
	Sections []*ListSection `json:"sections" validate:"required,min=1"`
}

type TemplateParameter struct {
	Type string `json:"type" validate:"required"`
	Text string `json:"text,omitempty"`
}

type TemplateComponent struct {
	Type       string               `json:"type" validate:"required"`
	Parameters []*TemplateParameter `json:"parameters,omitempty"`
}

type TemplateMessage struct {
	Name       string               `json:"name" validate:"required"`
	Language   string               `json:"language" validate:"required"`
	Components []*TemplateComponent `json:"components,omitempty"`
}

type SendTemplateMessageRequest struct {
	To       string           `json:"to" validate:"required"`
	Template *TemplateMessage `json:"template" validate:"required"`
}

type SendViewOnceMessageRequest struct {
	To      string     `json:"to" validate:"required"`
	Media   *MediaData `json:"media" validate:"required"`
	Caption string     `json:"caption,omitempty"`
}
