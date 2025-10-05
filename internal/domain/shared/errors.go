package shared

import "errors"

// Domain errors
var (
	// Session errors
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionAlreadyExists = errors.New("session already exists")
	ErrSessionNotConnected  = errors.New("session not connected")
	ErrInvalidSessionStatus = errors.New("invalid session status")
	ErrInvalidStatus        = errors.New("invalid session status")
	
	// Message errors
	ErrMessageNotFound     = errors.New("message not found")
	ErrInvalidMessageType  = errors.New("invalid message type")
	ErrEmptyMessageContent = errors.New("message content cannot be empty")
	ErrInvalidRecipient    = errors.New("invalid recipient")
	
	// Contact errors
	ErrContactNotFound = errors.New("contact not found")
	ErrInvalidJID      = errors.New("invalid JID format")
	
	// General errors
	ErrInvalidInput    = errors.New("invalid input")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrInternalError   = errors.New("internal server error")
)

// DomainError represents a domain-specific error
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

// NewDomainError creates a new domain error
func NewDomainError(code, message string, cause error) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}
