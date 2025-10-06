package dto

import (
	"fmt"
	"time"
)

type APIResponse struct {
	Success   bool        `json:"success" example:"true" description:"Whether the request was successful"`
	Data      interface{} `json:"data,omitempty" description:"Response data (present on success)"`
	Error     *ErrorInfo  `json:"error,omitempty" description:"Error information (present on failure)"`
	Timestamp time.Time   `json:"timestamp" example:"2025-01-15T10:30:00Z" description:"Response timestamp"`
} //@name APIResponse

type ErrorResponse struct {
	Error   string `json:"error" example:"validation_error" description:"Error code"`
	Message string `json:"message" example:"name is required" description:"Human readable error message"`
} //@name ErrorResponse

type APIErrorInfo struct {
	Code    string                 `json:"code" example:"validation_error" description:"Error code"`
	Message string                 `json:"message" example:"Validation failed" description:"Human readable error message"`
	Details map[string]interface{} `json:"details,omitempty" description:"Additional error details"`
} //@name APIErrorInfo

type ValidationError struct {
	Field   string `json:"field" example:"name" description:"Field that failed validation"`
	Message string `json:"message" example:"name is required" description:"Validation error message"`
}

func (e *ErrorInfo) Error() string {
	return e.Message
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type PaginationRequest struct {
	Limit  int `json:"limit" form:"limit" example:"20" validate:"min=1,max=100" description:"Number of items per page (1-100)"`
	Offset int `json:"offset" form:"offset" example:"0" validate:"min=0" description:"Number of items to skip"`
}

type PaginationResponse struct {
	Items   interface{} `json:"items" description:"Paginated items"`
	Total   int         `json:"total" example:"100" description:"Total number of items"`
	Limit   int         `json:"limit" example:"20" description:"Number of items per page"`
	Offset  int         `json:"offset" example:"0" description:"Number of items skipped"`
	HasMore bool        `json:"hasMore" example:"true" description:"Whether there are more items"`
}

type HealthResponse struct {
	Status    string            `json:"status" example:"healthy" description:"Overall health status"`
	Version   string            `json:"version" example:"1.0.0" description:"Application version"`
	Timestamp time.Time         `json:"timestamp" example:"2025-01-15T10:30:00Z" description:"Health check timestamp"`
	Services  map[string]string `json:"services" description:"Status of individual services"`
}

func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
	}
}

func NewErrorResponse(code, message string) *APIResponse {
	return &APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
		Timestamp: time.Now(),
	}
}

func NewErrorResponseWithDetails(code, message string, details map[string]interface{}) *APIResponse {
	return &APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func NewValidationErrorResponse(field, message string) *APIResponse {
	return NewErrorResponseWithDetails(
		"validation_error",
		"Validation failed",
		map[string]interface{}{
			"field":   field,
			"message": message,
		},
	)
}

const (
	ErrorCodeValidation    = "validation_error"
	ErrorCodeNotFound      = "not_found"
	ErrorCodeUnauthorized  = "unauthorized"
	ErrorCodeForbidden     = "forbidden"
	ErrorCodeConflict      = "conflict"
	ErrorCodeInternalError = "internal_error"
	ErrorCodeBadRequest    = "bad_request"
	ErrorCodeServiceError  = "service_error"
	ErrorCodeTimeout       = "timeout"
	ErrorCodeRateLimit     = "rate_limit"
)

const (
	DefaultLimit  = 20
	MaxLimit      = 100
	DefaultOffset = 0
)

func (p *PaginationRequest) ApplyDefaults() {
	if p.Limit <= 0 {
		p.Limit = DefaultLimit
	}
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
	if p.Offset < 0 {
		p.Offset = DefaultOffset
	}
}

func (p *PaginationRequest) Validate() error {
	if p.Limit < 1 || p.Limit > MaxLimit {
		return NewValidationError("limit", fmt.Sprintf("limit must be between 1 and %d", MaxLimit))
	}
	if p.Offset < 0 {
		return NewValidationError("offset", "offset must be non-negative")
	}
	return nil
}
