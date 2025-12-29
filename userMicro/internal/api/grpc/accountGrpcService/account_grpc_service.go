package accountGrpcService

import (
	"context"
	"errors"

	"github.com/PavelShe11/studbridge/authMicro/grpcApi"
	commondomain "github.com/PavelShe11/studbridge/common/domain"
	"github.com/PavelShe11/studbridge/common/translator"
	"github.com/PavelShe11/studbridge/user/internal/domain"
	"github.com/PavelShe11/studbridge/user/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

type accountGrpcService struct {
	grpcApi.UnimplementedAccountServiceServer
	accountService service.AccountService
	translator     *translator.Translator
}

func Register(server *grpc.Server, accountService service.AccountService, trans *translator.Translator) {
	grpcApi.RegisterAccountServiceServer(server, &accountGrpcService{
		accountService: accountService,
		translator:     trans,
	})
}

func valueToString(m map[string]*structpb.Value, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	return v.GetStringValue()
}

func (a accountGrpcService) CreateAccount(ctx context.Context, request *grpcApi.CreateAccountRequest) (*grpcApi.CreateAccountResponse, error) {
	lang := getLangFromContext(ctx)

	err := a.accountService.CreateAccount(domain.Account{
		FirstName: valueToString(request.UserData, "firstName"),
		LastName:  valueToString(request.UserData, "lastName"),
		Email:     valueToString(request.UserData, "email"),
	})

	if err != nil {
		a.translator.TranslateError(err, lang)
	}

	return &grpcApi.CreateAccountResponse{
		Error: mapToGrpcError(err),
	}, nil
}

func (a accountGrpcService) GetAccountByEmail(_ context.Context, request *grpcApi.GetAccountByEmailRequest) (*grpcApi.GetAccountResponse, error) {
	return a.accountMapToGetAccountResponse(
		a.accountService.GetAccountByEmail(
			request.GetEmail(),
		),
	)
}

func (a accountGrpcService) accountMapToGetAccountResponse(account *domain.Account, err error) (*grpcApi.GetAccountResponse, error) {
	if err != nil {
		return &grpcApi.GetAccountResponse{
			Result: &grpcApi.GetAccountResponse_Error{
				Error: mapToGrpcError(err),
			},
		}, nil
	}

	if account == nil {
		return &grpcApi.GetAccountResponse{
			Result: &grpcApi.GetAccountResponse_Error{
				Error: &grpcApi.Error{Code: grpcApi.ErrorCode_INTERNAL},
			},
		}, nil
	}

	return &grpcApi.GetAccountResponse{
		Result: &grpcApi.GetAccountResponse_Account_{
			Account: &grpcApi.GetAccountResponse_Account{
				AccountId: account.Id,
				Email:     account.Email,
			},
		},
	}, nil
}

func (a accountGrpcService) ValidateAccountData(ctx context.Context, request *grpcApi.ValidateAccountRequest) (*grpcApi.ValidateAccountResponse, error) {
	lang := getLangFromContext(ctx)

	err := a.accountService.ValidateAccountData(domain.Account{
		FirstName: valueToString(request.UserData, "firstName"),
		LastName:  valueToString(request.UserData, "lastName"),
		Email:     valueToString(request.UserData, "email"),
	})

	if err != nil {
		a.translator.TranslateError(err, lang)
	}

	return &grpcApi.ValidateAccountResponse{
		Error: mapToGrpcError(err),
	}, nil
}

func mapToGrpcError(e error) *grpcApi.Error {
	if e == nil {
		return nil
	}

	errs := make([]*grpcApi.Error_FieldError, 0)

	var validErr *commondomain.BaseValidationError
	if errors.As(e, &validErr) {
		for _, err := range validErr.FieldErrors {
			errs = append(errs, &grpcApi.Error_FieldError{
				Name:    err.NameField,
				Message: err.Message,
			})
		}
		return &grpcApi.Error{
			Code:           grpcApi.ErrorCode_VALIDATION,
			DetailedErrors: errs,
		}
	}

	var baseError *commondomain.BaseError
	if errors.As(e, &baseError) {
		return &grpcApi.Error{
			Code:           grpcApi.ErrorCode_INTERNAL,
			DetailedErrors: errs,
		}
	}

	return &grpcApi.Error{
		Code:           grpcApi.ErrorCode_INTERNAL,
		DetailedErrors: errs,
	}
}

func (a accountGrpcService) GetAccessTokenPayload(
	_ context.Context,
	request *grpcApi.GetAccessTokenPayloadRequest,
) (*grpcApi.GetAccessTokenPayloadResponse, error) {
	// Получить account из БД по ID
	account, err := a.accountService.GetAccountById(request.GetAccountId())

	if err != nil {
		return &grpcApi.GetAccessTokenPayloadResponse{
			Result: &grpcApi.GetAccessTokenPayloadResponse_Error{
				Error: mapToGrpcError(err),
			},
		}, nil
	}

	if account == nil {
		return nil, nil
	}

	// Создать map с данными для токена
	values := make(map[string]*structpb.Value)
	values["sub"] = structpb.NewStringValue(account.Id)

	return &grpcApi.GetAccessTokenPayloadResponse{
		Result: &grpcApi.GetAccessTokenPayloadResponse_Claims{
			Claims: &grpcApi.AccessTokenClaims{
				Values: values,
			},
		},
	}, nil
}

func getLangFromContext(ctx context.Context) string {
	lang := "en"
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if langs := md.Get("lang"); len(langs) > 0 {
			lang = langs[0]
		}
	}
	return lang
}
