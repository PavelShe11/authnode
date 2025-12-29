package entity

import (
	"github.com/PavelShe11/studbridge/authMicro/grpcApi"
	commonEntity "github.com/PavelShe11/studbridge/common/entity"
)

// NewInvalidCodeError creates a new instance of InvalidCode error
func NewInvalidCodeError() *commonEntity.BaseValidationError {
	return &commonEntity.BaseValidationError{
		BaseError: commonEntity.BaseError{Code: "invalidCode"},
		FieldErrors: []commonEntity.FieldError{{
			NameField: "code",
			Message:   "invalidCode",
			Params:    nil,
		}},
	}
}

// NewCodeExpiredError creates a new instance of CodeExpired error
func NewCodeExpiredError() *commonEntity.BaseValidationError {
	return &commonEntity.BaseValidationError{
		BaseError: commonEntity.BaseError{Code: "codeExpired"},
		FieldErrors: []commonEntity.FieldError{{
			NameField: "code",
			Message:   "codeExpired",
			Params:    nil,
		}},
	}
}

// NewInvalidRefreshTokenError creates a new instance of InvalidRefreshToken error
func NewInvalidRefreshTokenError() *commonEntity.BaseValidationError {
	return &commonEntity.BaseValidationError{
		BaseError: commonEntity.BaseError{Code: "invalidRefreshToken"},
		FieldErrors: []commonEntity.FieldError{{
			NameField: "refreshToken",
			Message:   "invalidRefreshToken",
			Params:    nil,
		}},
	}
}

// NewRefreshTokenExpiredError creates a new instance of RefreshTokenExpired error
func NewRefreshTokenExpiredError() *commonEntity.BaseValidationError {
	return &commonEntity.BaseValidationError{
		BaseError: commonEntity.BaseError{Code: "refreshTokenExpired"},
		FieldErrors: []commonEntity.FieldError{{
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

	fieldErrors := make([]commonEntity.FieldError, 0, len(grpcErr.DetailedErrors))
	for _, err := range grpcErr.DetailedErrors {
		fieldErrors = append(fieldErrors, commonEntity.FieldError{
			NameField: err.Name,
			Message:   err.Message,
		})
	}

	switch grpcErr.Code {
	case grpcApi.ErrorCode_VALIDATION:
		validationError := commonEntity.NewValidationError()
		validationError.FieldErrors = fieldErrors
		return validationError
	default:
		return commonEntity.NewInternalError()
	}
}
