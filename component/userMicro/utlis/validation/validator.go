package validation

import (
	"errors"
	"userMicro/internal/domain"

	"github.com/go-playground/validator/v10"
)

func Var(nameField string, field interface{}, tag string, error *domain.Error) {
	err := validator.New().Var(field, tag)
	if err == nil {
		return
	}
	errorField := domain.FieldError{
		Name: nameField,
	}
	var validErr validator.ValidationErrors
	errors.As(err, &validErr)
	for i, err := range validErr {
		errorField.Message += err.Tag()
		if len(validErr) < i+1 {
			errorField.Message += ","
		}
	}
	error.FieldErrors = append(error.FieldErrors, errorField)
}

func Struct(s interface{}) []domain.FieldError {
	result := make([]domain.FieldError, 0)
	err := validator.New().Struct(s)
	if err == nil {
		return nil
	}
	var validErr validator.ValidationErrors
	errors.As(err, &validErr)
	for _, err := range validErr {
		result = append(result, domain.FieldError{
			Name:    err.Field(),
			Message: err.Tag(),
		})
	}
	return result
}
