package domain

import (
	commondomain "github.com/PavelShe11/studbridge/common/domain"
)

// NewValidationError creates a new instance of ValidationError
func NewValidationError() *commondomain.BaseValidationError {
	return &commondomain.BaseValidationError{
		BaseError:   commondomain.BaseError{Code: "validationError"},
		FieldErrors: make([]commondomain.FieldError, 0),
	}
}
