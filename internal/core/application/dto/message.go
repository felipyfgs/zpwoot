package dto

import (
	"encoding/base64"
	"errors"
	"strings"
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
	URL      string `json:"url,omitempty" example:"https://example.com/image.jpg"`
	Base64   string `json:"base64,omitempty" example:"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="`
	FileName string `json:"fileName,omitempty" example:"image.jpg"`
	MimeType string `json:"mimeType,omitempty" example:"image/jpeg"`
	Caption  string `json:"caption,omitempty" example:"Check out this image!"`
} //@name MediaData

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90" example:"-23.550520"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180" example:"-46.633308"`
	Name      string  `json:"name,omitempty" example:"São Paulo"`
	Address   string  `json:"address,omitempty" example:"Av. Paulista, 1578 - Bela Vista, São Paulo - SP"`
} //@name Location

type ContactInfo struct {
	Name  string `json:"name" validate:"required" example:"John Doe"`
	Phone string `json:"phone" validate:"required" example:"5511999999999"`
	VCard string `json:"vcard,omitempty" example:"BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+5511999999999\nEND:VCARD"`
} //@name ContactInfo

type SendMessageResponse struct {
	Success     bool         `json:"success" example:"true"`
	ID          string       `json:"id" example:"3EB0A9253FA64269E11C9D"`
	To          string       `json:"to" example:"5511999999999@s.whatsapp.net"`
	Type        string       `json:"type" example:"text"`
	Content     string       `json:"content,omitempty" example:"Hello! This is a test message."`
	Timestamp   int64        `json:"timestamp" example:"1696570882"`
	Status      string       `json:"status" example:"sent"`
	ContextInfo *ContextInfo `json:"contextInfo,omitempty"`
} //@name SendMessageResponse

type ContextInfo struct {
	StanzaID    string `json:"stanzaId,omitempty" example:"3EB0A9253FA64269E11C9D"`
	Participant string `json:"participant,omitempty" example:"5511888888888@s.whatsapp.net"`
	QuotedID    string `json:"quotedId,omitempty" example:"3EB0A9253FA64269E11C9D"`
} //@name ContextInfo

type SendTextMessageRequest struct {
	Phone       string              `json:"phone" validate:"required" example:"5511999999999"`
	Text        string              `json:"text" validate:"required" example:"Hello! This is a test message from zpwoot API."`
	ContextInfo *ContextInfoRequest `json:"contextInfo,omitempty"`
} //@name SendTextMessageRequest

type ContextInfoRequest struct {
	StanzaID    string `json:"stanzaId" example:"3EB0A9253FA64269E11C9D"`
	Participant string `json:"participant,omitempty" example:"5511888888888@s.whatsapp.net"`
} //@name ContextInfoRequest

type SendImageMessageRequest struct {
	Phone       string              `json:"phone" validate:"required" example:"5511999999999"`
	File        string              `json:"file" validate:"required" example:"https://example.com/image.jpg"`
	MimeType    string              `json:"mimeType,omitempty" example:"image/jpeg"`
	FileName    string              `json:"fileName,omitempty" example:"image.jpg"`
	Caption     string              `json:"caption,omitempty" example:"Check out this beautiful image!"`
	ViewOnce    bool                `json:"viewOnce,omitempty" example:"false"`
	ContextInfo *ContextInfoRequest `json:"contextInfo,omitempty"`
} //@name SendImageMessageRequest

type SendAudioMessageRequest struct {
	Phone       string              `json:"phone" validate:"required" example:"5511999999999"`
	File        string              `json:"file" validate:"required" example:"https://example.com/audio.mp3"`
	MimeType    string              `json:"mimeType,omitempty" example:"audio/mpeg"`
	FileName    string              `json:"fileName,omitempty" example:"audio.mp3"`
	ViewOnce    bool                `json:"viewOnce,omitempty" example:"false"`
	ContextInfo *ContextInfoRequest `json:"contextInfo,omitempty"`
} //@name SendAudioMessageRequest

type SendVideoMessageRequest struct {
	Phone       string              `json:"phone" validate:"required" example:"5511999999999"`
	File        string              `json:"file" validate:"required" example:"https://example.com/video.mp4"`
	MimeType    string              `json:"mimeType,omitempty" example:"video/mp4"`
	FileName    string              `json:"fileName,omitempty" example:"video.mp4"`
	Caption     string              `json:"caption,omitempty" example:"Watch this amazing video!"`
	ViewOnce    bool                `json:"viewOnce,omitempty" example:"false"`
	ContextInfo *ContextInfoRequest `json:"contextInfo,omitempty"`
} //@name SendVideoMessageRequest

type SendDocumentMessageRequest struct {
	Phone       string              `json:"phone" validate:"required" example:"5511999999999"`
	File        string              `json:"file" validate:"required" example:"https://example.com/document.pdf"`
	MimeType    string              `json:"mimeType,omitempty" example:"application/pdf"`
	FileName    string              `json:"fileName,omitempty" example:"document.pdf"`
	Caption     string              `json:"caption,omitempty" example:"Important document attached"`
	ContextInfo *ContextInfoRequest `json:"contextInfo,omitempty"`
} //@name SendDocumentMessageRequest

type SendStickerMessageRequest struct {
	Phone       string              `json:"phone" validate:"required" example:"5511999999999"`
	File        string              `json:"file" validate:"required" example:"https://example.com/sticker.webp"`
	MimeType    string              `json:"mimeType,omitempty" example:"image/webp"`
	FileName    string              `json:"fileName,omitempty" example:"sticker.webp"`
	ContextInfo *ContextInfoRequest `json:"contextInfo,omitempty"`
} //@name SendStickerMessageRequest

type SendLocationMessageRequest struct {
	Phone       string              `json:"phone" validate:"required" example:"5511999999999"`
	Latitude    float64             `json:"latitude" validate:"required,min=-90,max=90" example:"-23.550520"`
	Longitude   float64             `json:"longitude" validate:"required,min=-180,max=180" example:"-46.633308"`
	Name        string              `json:"name,omitempty" example:"São Paulo"`
	Address     string              `json:"address,omitempty" example:"Av. Paulista, 1578 - Bela Vista, São Paulo - SP"`
	ContextInfo *ContextInfoRequest `json:"contextInfo,omitempty"`
} //@name SendLocationMessageRequest

type SendContactMessageRequest struct {
	Phone       string              `json:"phone" validate:"required" example:"5511999999999"`
	Contact     *ContactInfo        `json:"contact" validate:"required"`
	ContextInfo *ContextInfoRequest `json:"contextInfo,omitempty"`
} //@name SendContactMessageRequest

type SendReactionMessageRequest struct {
	Phone     string `json:"phone" validate:"required" example:"5511999999999"`
	MessageID string `json:"messageId" validate:"required" example:"3EB0C767D0D1A6F4FD29"`
	Reaction  string `json:"reaction" validate:"required" example:"👍"`

	FromMe *bool `json:"fromMe,omitempty" example:"false"`
} //@name SendReactionMessageRequest

type SendPollMessageRequest struct {
	Phone                  string   `json:"phone" validate:"required" example:"5511999999999"`
	Name                   string   `json:"name" validate:"required" example:"What's your favorite color?"`
	Options                []string `json:"options" validate:"required,min=2,max=12" example:"Red,Blue,Green,Yellow"`
	SelectableOptionsCount int      `json:"selectableOptionsCount" validate:"required,min=1" example:"1"`
} //@name SendPollMessageRequest

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

	return r.validateByType()
}

func (r *SendMessageRequest) validateByType() error {
	switch r.Type {
	case "text":
		return validators.ValidateMessageText(r.Text)
	case "media":
		return r.validateMedia()
	case "location":
		return r.validateLocation()
	case "contact":
		return r.validateContact()
	default:
		return ErrInvalidMessageType
	}
}

func (r *SendMessageRequest) validateMedia() error {
	if r.Media == nil {
		return ErrMediaRequired
	}

	if r.Media.URL == "" && r.Media.Base64 == "" {
		return ErrMediaContentRequired
	}

	if err := validators.ValidateCaption(r.Media.Caption); err != nil {
		return err
	}

	return validators.ValidateFileName(r.Media.FileName)
}

func (r *SendMessageRequest) validateLocation() error {
	if r.Location == nil {
		return ErrLocationRequired
	}

	return r.Location.Validate()
}

func (r *SendMessageRequest) validateContact() error {
	if r.Contact == nil {
		return ErrContactRequired
	}

	return r.Contact.Validate()
}

func (l *Location) Validate() error {
	if err := validators.ValidateLatitude(l.Latitude); err != nil {
		return err
	}

	if err := validators.ValidateLongitude(l.Longitude); err != nil {
		return err
	}

	if err := validators.ValidateLocationName(l.Name); err != nil {
		return err
	}

	if err := validators.ValidateLocationAddress(l.Address); err != nil {
		return err
	}

	return nil
}

func (c *ContactInfo) Validate() error {
	if err := validators.ValidateContactName(c.Name); err != nil {
		return err
	}

	if err := validators.ValidatePhoneNumber(c.Phone); err != nil {
		return err
	}

	if c.Phone == "" {
		return ErrContactPhoneRequired
	}

	return nil
}

var (
	ErrRecipientRequired      = errors.New("recipient is required")
	ErrSessionIDRequired      = errors.New("sessionId is required")
	ErrInvalidJSON            = errors.New("invalid JSON body")
	ErrMessageTypeRequired    = errors.New("message type is required")
	ErrInvalidMessageType     = errors.New("invalid message type")
	ErrTextRequired           = errors.New("text is required for text messages")
	ErrMediaRequired          = errors.New("media is required for media messages")
	ErrMediaContentRequired   = errors.New("media URL or base64 content is required")
	ErrLocationRequired       = errors.New("location is required for location messages")
	ErrContactRequired        = errors.New("contact is required for contact messages")
	ErrInvalidLatitude        = errors.New("latitude must be between -90 and 90")
	ErrInvalidLongitude       = errors.New("longitude must be between -180 and 180")
	ErrContactNameRequired    = errors.New("contact name is required")
	ErrContactPhoneRequired   = errors.New("contact phone is required")
	ErrUnsupportedMessageType = errors.New("unsupported message type")
)

func (m *MediaData) ToInterfacesMediaData() *output.MediaData {
	if m == nil {
		return nil
	}

	var data []byte

	var err error

	if m.Base64 != "" {
		data, err = m.decodeBase64Data()
		if err != nil {
			data = []byte{}
		}
	} else if m.URL != "" {
		data = []byte{}
	}

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

type SendMultipleContactsRequest struct {
	Phone    string         `json:"phone" validate:"required" example:"5511999999999"`
	Contacts []*ContactInfo `json:"contacts" validate:"required,min=1"`
} //@name SendMultipleContactsRequest

type Button struct {
	ID   string `json:"id" validate:"required" example:"btn_1"`
	Text string `json:"text" validate:"required" example:"Click Me"`
} //@name Button

type SendButtonsMessageRequest struct {
	Phone   string    `json:"phone" validate:"required" example:"5511999999999"`
	Text    string    `json:"text" validate:"required" example:"Please choose an option:"`
	Buttons []*Button `json:"buttons" validate:"required,min=1,max=3"`
} //@name SendButtonsMessageRequest

type ListRow struct {
	ID          string `json:"id" validate:"required" example:"row_1"`
	Title       string `json:"title" validate:"required" example:"Option 1"`
	Description string `json:"description,omitempty" example:"Description for option 1"`
} //@name ListRow

type ListSection struct {
	Title string     `json:"title" validate:"required" example:"Section 1"`
	Rows  []*ListRow `json:"rows" validate:"required,min=1"`
} //@name ListSection

type SendListMessageRequest struct {
	Phone    string         `json:"phone" validate:"required" example:"5511999999999"`
	Text     string         `json:"text" validate:"required" example:"Please select an option from the list"`
	Title    string         `json:"title" validate:"required" example:"Menu Options"`
	Sections []*ListSection `json:"sections" validate:"required,min=1"`
} //@name SendListMessageRequest

type TemplateParameter struct {
	Type string `json:"type" validate:"required" example:"text"`
	Text string `json:"text,omitempty" example:"John"`
} //@name TemplateParameter

type TemplateComponent struct {
	Type       string               `json:"type" validate:"required" example:"body"`
	Parameters []*TemplateParameter `json:"parameters,omitempty"`
} //@name TemplateComponent

type TemplateMessage struct {
	Name       string               `json:"name" validate:"required" example:"hello_world"`
	Language   string               `json:"language" validate:"required" example:"en_US"`
	Components []*TemplateComponent `json:"components,omitempty"`
} //@name TemplateMessage

type SendTemplateMessageRequest struct {
	Phone    string           `json:"phone" validate:"required" example:"5511999999999"`
	Template *TemplateMessage `json:"template" validate:"required"`
} //@name SendTemplateMessageRequest
type DeleteMessageRequest struct {
	Phone     string `json:"phone" validate:"required" example:"5511999999999"`
	MessageID string `json:"messageId" validate:"required" example:"3EB0C767D0D1A2B5F8"`
} //@name DeleteMessageRequest
type DeleteMessageResponse struct {
	Success   bool   `json:"success" example:"true"`
	MessageID string `json:"messageId" example:"3EB0C767D0D1A2B5F8"`
	Timestamp int64  `json:"timestamp" example:"1696570882"`
} //@name DeleteMessageResponse
type EditMessageRequest struct {
	Phone     string `json:"phone" validate:"required" example:"5511999999999"`
	MessageID string `json:"messageId" validate:"required" example:"3EB0C767D0D1A2B5F8"`
	Text      string `json:"text" validate:"required" example:"Edited message text"`
} //@name EditMessageRequest
type EditMessageResponse struct {
	Success   bool   `json:"success" example:"true"`
	MessageID string `json:"messageId" example:"3EB0C767D0D1A2B5F8"`
	Timestamp int64  `json:"timestamp" example:"1696570882"`
} //@name EditMessageResponse
type MarkReadRequest struct {
	Phone      string   `json:"phone" validate:"required" example:"5511999999999"`
	MessageIDs []string `json:"messageIds" validate:"required,min=1" example:"3EB0C767D0D1A2B5F8,3EB0C767D0D1A2B5F9"`
} //@name MarkReadRequest
type MarkReadResponse struct {
	Success bool `json:"success" example:"true"`
} //@name MarkReadResponse
type HistorySyncRequest struct {
	Count int `json:"count,omitempty" example:"50" description:"Number of messages to sync (default: 50)"`
} //@name HistorySyncRequest
type HistorySyncResponse struct {
	Success   bool  `json:"success" example:"true"`
	Timestamp int64 `json:"timestamp" example:"1696570882"`
} //@name HistorySyncResponse

// decodeBase64Data extracts and decodes base64 data, handling data URI format
func (m *MediaData) decodeBase64Data() ([]byte, error) {
	base64Data := m.Base64
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	return base64.StdEncoding.DecodeString(base64Data)
}
