package message

import (
	"errors"
	"fmt"
)

var (
	ErrMessageNotFound         = errors.New("message not found")
	ErrMessageAlreadyExists    = errors.New("message already exists")
	ErrInvalidMessageData      = errors.New("invalid message data")
	ErrMessageSyncFailed       = errors.New("message sync failed")
	ErrInvalidMessageType      = errors.New("invalid message type")
	ErrInvalidMessageContent   = errors.New("invalid message content")
	ErrMessageSendFailed       = errors.New("message send failed")
	ErrInvalidSyncStatus       = errors.New("invalid sync status")
	ErrChatwootSyncFailed      = errors.New("chatwoot sync failed")
	ErrMessageUpdateFailed     = errors.New("message update failed")
	ErrInvalidMessageID        = errors.New("invalid message ID")
	ErrMessageValidationFailed = errors.New("message validation failed")
	ErrMediaProcessingFailed   = errors.New("media processing failed")
	ErrMessageTooLarge         = errors.New("message too large")
	ErrInvalidRecipient        = errors.New("invalid recipient")
)

type MessageError struct {
	Code    string
	Message string
	Cause   error
	Field   string
}

func (e *MessageError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("messaging error [%s]: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	if e.Field != "" {
		return fmt.Sprintf("messaging error [%s]: %s (field: %s)", e.Code, e.Message, e.Field)
	}
	return fmt.Sprintf("messaging error [%s]: %s", e.Code, e.Message)
}

func (e *MessageError) Unwrap() error {
	return e.Cause
}

func NewMessageError(code, message string, cause error) *MessageError {
	return &MessageError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func NewMessageValidationError(field, message string) *MessageError {
	return &MessageError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Field:   field,
	}
}

const (
	ErrCodeMessageNotFound         = "MESSAGE_NOT_FOUND"
	ErrCodeMessageAlreadyExists    = "MESSAGE_ALREADY_EXISTS"
	ErrCodeInvalidMessageData      = "INVALID_MESSAGE_DATA"
	ErrCodeMessageSyncFailed       = "MESSAGE_SYNC_FAILED"
	ErrCodeInvalidMessageType      = "INVALID_MESSAGE_TYPE"
	ErrCodeInvalidMessageContent   = "INVALID_MESSAGE_CONTENT"
	ErrCodeMessageSendFailed       = "MESSAGE_SEND_FAILED"
	ErrCodeInvalidSyncStatus       = "INVALID_SYNC_STATUS"
	ErrCodeChatwootSyncFailed      = "CHATWOOT_SYNC_FAILED"
	ErrCodeMessageUpdateFailed     = "MESSAGE_UPDATE_FAILED"
	ErrCodeInvalidMessageID        = "INVALID_MESSAGE_ID"
	ErrCodeMessageValidationFailed = "MESSAGE_VALIDATION_FAILED"
	ErrCodeMediaProcessingFailed   = "MEDIA_PROCESSING_FAILED"
	ErrCodeMessageTooLarge         = "MESSAGE_TOO_LARGE"
	ErrCodeInvalidRecipient        = "INVALID_RECIPIENT"
)

func ErrMessageNotFoundWithID(id string) *MessageError {
	return NewMessageError(ErrCodeMessageNotFound, fmt.Sprintf("message with ID %s not found", id), nil)
}

func ErrMessageAlreadyExistsWithZpID(zpMessageID string) *MessageError {
	return NewMessageError(ErrCodeMessageAlreadyExists, fmt.Sprintf("message with ZP ID %s already exists", zpMessageID), nil)
}

func ErrInvalidMessageTypeValue(msgType string) *MessageError {
	return NewMessageError(ErrCodeInvalidMessageType, fmt.Sprintf("invalid message type: %s", msgType), nil)
}

func ErrInvalidSyncStatusValue(status string) *MessageError {
	return NewMessageError(ErrCodeInvalidSyncStatus, fmt.Sprintf("invalid sync status: %s", status), nil)
}

func ErrMessageSendFailedWithReason(reason string) *MessageError {
	return NewMessageError(ErrCodeMessageSendFailed, fmt.Sprintf("message send failed: %s", reason), nil)
}

func ErrChatwootSyncFailedWithReason(reason string) *MessageError {
	return NewMessageError(ErrCodeChatwootSyncFailed, fmt.Sprintf("chatwoot sync failed: %s", reason), nil)
}

func ErrMediaProcessingFailedWithReason(reason string) *MessageError {
	return NewMessageError(ErrCodeMediaProcessingFailed, fmt.Sprintf("media processing failed: %s", reason), nil)
}

func ErrMessageTooLargeWithSize(size int64, maxSize int64) *MessageError {
	return NewMessageError(ErrCodeMessageTooLarge, fmt.Sprintf("message size %d bytes exceeds maximum %d bytes", size, maxSize), nil)
}

func ErrInvalidRecipientJID(jid string) *MessageError {
	return NewMessageError(ErrCodeInvalidRecipient, fmt.Sprintf("invalid recipient JID: %s", jid), nil)
}
