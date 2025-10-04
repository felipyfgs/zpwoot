package contact

import (
	"errors"
	"fmt"
)

var (
	ErrContactNotFound         = errors.New("contact not found")
	ErrContactAlreadyExists    = errors.New("contact already exists")
	ErrInvalidContactData      = errors.New("invalid contact data")
	ErrContactSyncFailed       = errors.New("contact sync failed")
	ErrInvalidPhoneNumber      = errors.New("invalid phone number")
	ErrInvalidEmail            = errors.New("invalid email address")
	ErrContactBlocked          = errors.New("contact is blocked")
	ErrInvalidSyncStatus       = errors.New("invalid sync status")
	ErrChatwootSyncFailed      = errors.New("chatwoot sync failed")
	ErrContactUpdateFailed     = errors.New("contact update failed")
	ErrInvalidContactType      = errors.New("invalid contact type")
	ErrContactValidationFailed = errors.New("contact validation failed")
)

type ContactError struct {
	Code    string
	Message string
	Cause   error
	Field   string
}

func (e *ContactError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("contact error [%s]: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	if e.Field != "" {
		return fmt.Sprintf("contact error [%s]: %s (field: %s)", e.Code, e.Message, e.Field)
	}
	return fmt.Sprintf("contact error [%s]: %s", e.Code, e.Message)
}

func (e *ContactError) Unwrap() error {
	return e.Cause
}

func NewContactError(code, message string, cause error) *ContactError {
	return &ContactError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func NewContactValidationError(field, message string) *ContactError {
	return &ContactError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Field:   field,
	}
}

const (
	ErrCodeContactNotFound         = "CONTACT_NOT_FOUND"
	ErrCodeContactAlreadyExists    = "CONTACT_ALREADY_EXISTS"
	ErrCodeInvalidContactData      = "INVALID_CONTACT_DATA"
	ErrCodeContactSyncFailed       = "CONTACT_SYNC_FAILED"
	ErrCodeInvalidPhoneNumber      = "INVALID_PHONE_NUMBER"
	ErrCodeInvalidEmail            = "INVALID_EMAIL"
	ErrCodeContactBlocked          = "CONTACT_BLOCKED"
	ErrCodeInvalidSyncStatus       = "INVALID_SYNC_STATUS"
	ErrCodeChatwootSyncFailed      = "CHATWOOT_SYNC_FAILED"
	ErrCodeContactUpdateFailed     = "CONTACT_UPDATE_FAILED"
	ErrCodeInvalidContactType      = "INVALID_CONTACT_TYPE"
	ErrCodeContactValidationFailed = "CONTACT_VALIDATION_FAILED"
)

func ErrContactNotFoundWithID(id string) *ContactError {
	return NewContactError(ErrCodeContactNotFound, fmt.Sprintf("contact with ID %s not found", id), nil)
}

func ErrContactAlreadyExistsWithJID(jid string) *ContactError {
	return NewContactError(ErrCodeContactAlreadyExists, fmt.Sprintf("contact with JID %s already exists", jid), nil)
}

func ErrInvalidPhoneNumberFormat(phone string) *ContactError {
	return NewContactError(ErrCodeInvalidPhoneNumber, fmt.Sprintf("invalid phone number format: %s", phone), nil)
}

func ErrInvalidEmailFormat(email string) *ContactError {
	return NewContactError(ErrCodeInvalidEmail, fmt.Sprintf("invalid email format: %s", email), nil)
}

func ErrInvalidSyncStatusValue(status string) *ContactError {
	return NewContactError(ErrCodeInvalidSyncStatus, fmt.Sprintf("invalid sync status: %s", status), nil)
}

func ErrChatwootSyncFailedWithReason(reason string) *ContactError {
	return NewContactError(ErrCodeChatwootSyncFailed, fmt.Sprintf("chatwoot sync failed: %s", reason), nil)
}
