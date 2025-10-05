package shared

import "errors"

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionAlreadyExists = errors.New("session already exists")
	ErrSessionNotConnected  = errors.New("session not connected")
	ErrInvalidSessionStatus = errors.New("invalid session status")
	ErrInvalidStatus        = errors.New("invalid session status")

	ErrMessageNotFound     = errors.New("message not found")
	ErrInvalidMessageType  = errors.New("invalid message type")
	ErrEmptyMessageContent = errors.New("message content cannot be empty")
	ErrInvalidRecipient    = errors.New("invalid recipient")

	ErrContactNotFound = errors.New("contact not found")
	ErrInvalidJID      = errors.New("invalid JID format")

	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternalError = errors.New("internal server error")
)

type DomainError struct {
	Code    string
	Message string
	Cause   error
}

func (e DomainError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e DomainError) Unwrap() error {
	return e.Cause
}

func NewDomainError(code, message string, cause error) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}
