package dto

import (
	"fmt"
	"time"
)

// APIResponse represents a generic API response wrapper
type APIResponse struct {
	Success   bool        `json:"success" example:"true" description:"Whether the request was successful"`
	Data      interface{} `json:"data,omitempty" description:"Response data (present on success)"`
	Error     *ErrorInfo  `json:"error,omitempty" description:"Error information (present on failure)"`
	Timestamp time.Time   `json:"timestamp" example:"2025-01-15T10:30:00Z" description:"Response timestamp"`
}

// ErrorResponse represents an error response (for Swagger compatibility)
type ErrorResponse struct {
	Error   string `json:"error" example:"validation_error" description:"Error code"`
	Message string `json:"message" example:"name is required" description:"Human readable error message"`
}

// ErrorInfo represents detailed error information
type ErrorInfo struct {
	Code    string                 `json:"code" example:"validation_error" description:"Error code"`
	Message string                 `json:"message" example:"Validation failed" description:"Human readable error message"`
	Details map[string]interface{} `json:"details,omitempty" description:"Additional error details"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field" example:"name" description:"Field that failed validation"`
	Message string `json:"message" example:"name is required" description:"Validation error message"`
}

// Error implements the error interface
func (e *ErrorInfo) Error() string {
	return e.Message
}

// Error implements the error interface for ValidationError
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Limit  int `json:"limit" form:"limit" example:"20" validate:"min=1,max=100" description:"Number of items per page (1-100)"`
	Offset int `json:"offset" form:"offset" example:"0" validate:"min=0" description:"Number of items to skip"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Items   interface{} `json:"items" description:"Paginated items"`
	Total   int         `json:"total" example:"100" description:"Total number of items"`
	Limit   int         `json:"limit" example:"20" description:"Number of items per page"`
	Offset  int         `json:"offset" example:"0" description:"Number of items skipped"`
	HasMore bool        `json:"hasMore" example:"true" description:"Whether there are more items"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status" example:"healthy" description:"Overall health status"`
	Version   string            `json:"version" example:"1.0.0" description:"Application version"`
	Timestamp time.Time         `json:"timestamp" example:"2025-01-15T10:30:00Z" description:"Health check timestamp"`
	Services  map[string]string `json:"services" description:"Status of individual services"`
}

// Helper functions

// NewSuccessResponse creates a successful API response
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewErrorResponse creates an error API response
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

// NewErrorResponseWithDetails creates an error API response with details
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

// NewValidationError creates a validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// NewValidationErrorResponse creates a validation error response
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

// Common error codes
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

// Default pagination values
const (
	DefaultLimit  = 20
	MaxLimit      = 100
	DefaultOffset = 0
)

// ApplyDefaults applies default values to pagination request
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

// Validate validates pagination parameters
func (p *PaginationRequest) Validate() error {
	if p.Limit < 1 || p.Limit > MaxLimit {
		return NewValidationError("limit", fmt.Sprintf("limit must be between 1 and %d", MaxLimit))
	}
	if p.Offset < 0 {
		return NewValidationError("offset", "offset must be non-negative")
	}
	return nil
}
