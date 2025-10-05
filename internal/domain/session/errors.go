package session

import (
	"errors"
	"fmt"
)

var (
	ErrSessionNotFound         = errors.New("session not found")
	ErrSessionAlreadyExists    = errors.New("session already exists")
	ErrInvalidSessionID        = errors.New("invalid session ID")
	ErrInvalidSessionName      = errors.New("invalid session name")
	ErrSessionIDTooLong        = errors.New("session ID too long")
	ErrSessionAlreadyConnected = errors.New("session already connected")
	ErrSessionNotConnected     = errors.New("session not connected")
	ErrQRCodeNotAvailable      = errors.New("QR code not available")
	ErrQRCodeExpired           = errors.New("QR code expired")
	ErrInvalidProxyConfig      = errors.New("invalid proxy configuration")
	ErrSessionConnectionFailed = errors.New("session connection failed")
	ErrSessionDisconnectFailed = errors.New("session disconnect failed")
	ErrInvalidSessionState     = errors.New("invalid session state")
	ErrSessionValidationFailed = errors.New("session validation failed")
	ErrProxyConnectionFailed   = errors.New("proxy connection failed")
	ErrQRCodeGenerationFailed  = errors.New("QR code generation failed")
	ErrSessionRestoreFailed    = errors.New("session restore failed")
	ErrInvalidDeviceInfo       = errors.New("invalid device info")
	ErrSessionTimeout          = errors.New("session timeout")
	ErrInvalidDeviceJID        = errors.New("invalid device JID format")
	ErrConnectionFailed        = errors.New("failed to connect to WhatsApp")
	ErrPairingFailed           = errors.New("device pairing failed")
	ErrLogoutFailed            = errors.New("failed to logout session")
	ErrSessionBusy             = errors.New("session is busy with another operation")
	ErrInvalidOperation        = errors.New("invalid operation for current session state")
	ErrOperationTimeout        = errors.New("operation timed out")
)

type SessionError struct {
	Code    string
	Message string
	Cause   error
	Field   string
}

func (e *SessionError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("session error [%s]: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	if e.Field != "" {
		return fmt.Sprintf("session error [%s]: %s (field: %s)", e.Code, e.Message, e.Field)
	}
	return fmt.Sprintf("session error [%s]: %s", e.Code, e.Message)
}

func (e *SessionError) Unwrap() error {
	return e.Cause
}

func NewSessionError(code, message string, cause error) *SessionError {
	return &SessionError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func NewSessionValidationError(field, message string) *SessionError {
	return &SessionError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Field:   field,
	}
}

const (
	ErrCodeSessionNotFound         = "SESSION_NOT_FOUND"
	ErrCodeSessionAlreadyExists    = "SESSION_ALREADY_EXISTS"
	ErrCodeInvalidSessionID        = "INVALID_SESSION_ID"
	ErrCodeSessionIDTooLong        = "SESSION_ID_TOO_LONG"
	ErrCodeSessionAlreadyConnected = "SESSION_ALREADY_CONNECTED"
	ErrCodeSessionNotConnected     = "SESSION_NOT_CONNECTED"
	ErrCodeQRCodeNotAvailable      = "QR_CODE_NOT_AVAILABLE"
	ErrCodeQRCodeExpired           = "QR_CODE_EXPIRED"
	ErrCodeInvalidProxyConfig      = "INVALID_PROXY_CONFIG"
	ErrCodeSessionConnectionFailed = "SESSION_CONNECTION_FAILED"
	ErrCodeSessionDisconnectFailed = "SESSION_DISCONNECT_FAILED"
	ErrCodeInvalidSessionState     = "INVALID_SESSION_STATE"
	ErrCodeSessionValidationFailed = "SESSION_VALIDATION_FAILED"
	ErrCodeProxyConnectionFailed   = "PROXY_CONNECTION_FAILED"
	ErrCodeQRCodeGenerationFailed  = "QR_CODE_GENERATION_FAILED"
	ErrCodeSessionRestoreFailed    = "SESSION_RESTORE_FAILED"
	ErrCodeInvalidDeviceInfo       = "INVALID_DEVICE_INFO"
	ErrCodeSessionTimeout          = "SESSION_TIMEOUT"
)

func ErrSessionNotFoundWithID(id string) *SessionError {
	return NewSessionError(ErrCodeSessionNotFound, fmt.Sprintf("session with ID %s not found", id), nil)
}

func ErrSessionAlreadyExistsWithName(name string) *SessionError {
	return NewSessionError(ErrCodeSessionAlreadyExists, fmt.Sprintf("session with name %s already exists", name), nil)
}

func ErrInvalidSessionIDFormat(id string) *SessionError {
	return NewSessionError(ErrCodeInvalidSessionID, fmt.Sprintf("invalid session ID format: %s", id), nil)
}

func ErrSessionConnectionFailedWithReason(reason string) *SessionError {
	return NewSessionError(ErrCodeSessionConnectionFailed, fmt.Sprintf("session connection failed: %s", reason), nil)
}

func ErrQRCodeGenerationFailedWithReason(reason string) *SessionError {
	return NewSessionError(ErrCodeQRCodeGenerationFailed, fmt.Sprintf("QR code generation failed: %s", reason), nil)
}

func ErrProxyConnectionFailedWithReason(reason string) *SessionError {
	return NewSessionError(ErrCodeProxyConnectionFailed, fmt.Sprintf("proxy connection failed: %s", reason), nil)
}
