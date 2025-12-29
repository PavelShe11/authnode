package domain

import (
	"github.com/PavelShe11/studbridge/authMicro/grpcApi"
	commondomain "github.com/PavelShe11/studbridge/common/domain"
)

// NewInvalidCodeError creates a new instance of InvalidCode error
func NewInvalidCodeError() *commondomain.BaseValidationError {
	return &commondomain.BaseValidationError{
		BaseError: commondomain.BaseError{Code: "invalidCode"},
		FieldErrors: []commondomain.FieldError{{
			NameField: "code",
			Message:   "invalidCode",
			Params:    nil,
		}},
	}
}

// NewCodeExpiredError creates a new instance of CodeExpired error
func NewCodeExpiredError() *commondomain.BaseValidationError {
	return &commondomain.BaseValidationError{
		BaseError: commondomain.BaseError{Code: "codeExpired"},
		FieldErrors: []commondomain.FieldError{{
			NameField: "code",
			Message:   "codeExpired",
			Params:    nil,
		}},
	}
}

// NewInvalidRefreshTokenError creates a new instance of InvalidRefreshToken error
func NewInvalidRefreshTokenError() *commondomain.BaseValidationError {
	return &commondomain.BaseValidationError{
		BaseError: commondomain.BaseError{Code: "invalidRefreshToken"},
		FieldErrors: []commondomain.FieldError{{
			NameField: "refreshToken",
			Message:   "invalidRefreshToken",
			Params:    nil,
		}},
	}
}

// NewRefreshTokenExpiredError creates a new instance of RefreshTokenExpired error
func NewRefreshTokenExpiredError() *commondomain.BaseValidationError {
	return &commondomain.BaseValidationError{
		BaseError: commondomain.BaseError{Code: "refreshTokenExpired"},
		FieldErrors: []commondomain.FieldError{{
			NameField: "refreshToken",
			Message:   "refreshTokenExpired",
			Params:    nil,
		}},
	}
}

func GrpcErrorMapToError(grpcErr *grpcApi.Error) error {
	if grpcErr == nil {
		return nil
	}

	fieldErrors := make([]commondomain.FieldError, 0, len(grpcErr.DetailedErrors))
	for _, err := range grpcErr.DetailedErrors {
		fieldErrors = append(fieldErrors, commondomain.FieldError{
			NameField: err.Name,
			Message:   err.Message,
		})
	}

	switch grpcErr.Code {
	case grpcApi.ErrorCode_VALIDATION:
		validationError := commondomain.NewValidationError()
		validationError.FieldErrors = fieldErrors
		return validationError
	default:
		return commondomain.NewInternalError()
	}
}
