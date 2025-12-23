package domain

import (
	"authMicro/internal/api/grpcService"
	"fmt"
	"strings"
)

type Error struct {
	Name        string       `json:"name"`
	FieldErrors []FieldError `json:"fieldErrors"`
}

type FieldError struct {
	Name    string            `json:"name"`
	Message string            `json:"message"`
	Params  map[string]string `json:"-"` // Parameters for validation errors (not serialized to JSON)
}

// Error implements the error interface
func (e *Error) Error() string {
	if len(e.FieldErrors) == 0 {
		return e.Name
	}

	var fieldMessages []string
	for _, fe := range e.FieldErrors {
		fieldMessages = append(fieldMessages, fmt.Sprintf("%s: %s", fe.Name, fe.Message))
	}

	return fmt.Sprintf("%s [%s]", e.Name, strings.Join(fieldMessages, ", "))
}

func GrpcErrorMapToError(errs *grpcService.Error) *Error {
	result := Error{
		Name:        errs.Name,
		FieldErrors: make([]FieldError, 0),
	}
	for _, err := range errs.DetailedErrors {
		result.FieldErrors = append(result.FieldErrors, FieldError{
			Name:    err.Name,
			Message: err.Message,
		})
	}
	return &result
}
