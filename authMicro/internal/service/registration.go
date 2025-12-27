package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PavelShe11/studbridge/auth/internal/config"
	"github.com/PavelShe11/studbridge/auth/internal/domain"
	"github.com/PavelShe11/studbridge/auth/internal/repository"
	"github.com/PavelShe11/studbridge/auth/utlis/converter"
	"github.com/PavelShe11/studbridge/auth/utlis/generator"
	"github.com/PavelShe11/studbridge/auth/utlis/hash"
	"github.com/PavelShe11/studbridge/authMicro/grpcApi"
	commondomain "github.com/PavelShe11/studbridge/common/domain"
	"github.com/PavelShe11/studbridge/common/logger"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type RegisterAnswer struct {
	CodeExpires int64  `json:"codeExpires"`
	CodePattern string `json:"codePattern"`
}

type RegistrationService struct {
	registrationSessionRepository repository.RegistrationSessionRepository
	accountServiceClient          grpcApi.AccountServiceClient
	logger                        logger.Logger
	CodeGenConfig                 *config.CodeGenConfig
}

func NewRegistrationService(registrationSessionRepository repository.RegistrationSessionRepository, accountServiceClient grpcApi.AccountServiceClient, logger logger.Logger, codeGenConfig *config.CodeGenConfig) RegistrationService {
	return RegistrationService{
		registrationSessionRepository: registrationSessionRepository,
		accountServiceClient:          accountServiceClient,
		logger:                        logger,
		CodeGenConfig:                 codeGenConfig,
	}
}

func (r *RegistrationService) validateRegistrationData(userData map[string]any, lang string) (map[string]*structpb.Value, error) {
	grpcMap, err := converter.ConvertToGrpcMap(userData)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Add metadata
	md := metadata.Pairs("lang", lang)
	ctx = metadata.NewOutgoingContext(ctx, md)

	validationResponse, err := r.accountServiceClient.ValidateAccountData(
		ctx,
		&grpcApi.ValidateAccountRequest{UserData: grpcMap},
	)
	if err != nil {
		st, _ := status.FromError(err)
		r.logger.Error(fmt.Errorf("ValidateAccountData error: %v, grpc status: %v", err, st))
		return nil, commondomain.NewInternalError()
	}
	if validationResponse.Error != nil {
		return nil, domain.GrpcErrorMapToError(validationResponse.Error)
	}

	return grpcMap, nil
}

func (r *RegistrationService) getAccountByEmail(email string) (*grpcApi.GetAccountResponse, error) {
	ctxG, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accountGrpc, err := r.accountServiceClient.GetAccountByEmail(
		ctxG,
		&grpcApi.GetAccountByEmailRequest{Email: email},
	)

	if err != nil {
		st, _ := status.FromError(err)
		r.logger.Error(fmt.Errorf("GetAccountByEmail error: %v, grpc status: %v", err, st))
		return nil, commondomain.NewInternalError()
	}
	return accountGrpc, nil
}

func (r *RegistrationService) cleanupExpiredSessions() {
	cleanupCtx, cancelCleanup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCleanup()
	if err := r.registrationSessionRepository.CleanExpired(cleanupCtx); err != nil {
		r.logger.Error(fmt.Errorf("error cleaning expired registration sessions: %w", err))
	}
}

func (r *RegistrationService) Register(userData map[string]any, lang string) (*RegisterAnswer, error) {
	r.cleanupExpiredSessions()

	if _, err := r.validateRegistrationData(userData, lang); err != nil {
		return nil, err
	}

	email, ok := userData["email"].(string)
	if !ok {
		r.logger.Error(errors.New("email not found in response"))
		return nil, commondomain.NewInternalError()
	}

	accountGrpc, err := r.getAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	var session *domain.RegistrationSession

	if account, ok := accountGrpc.Result.(*grpcApi.GetAccountResponse_Account); ok && account != nil {
		session, err = r.createOrUpdateSession(email, "")
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}
	} else {
		var plaintextCode string
		plaintextCode, err = generator.Reggen(r.CodeGenConfig.CodePattern, r.CodeGenConfig.CodeMaxLength)
		if err != nil {
			r.logger.Error(err)
			return nil, commondomain.NewInternalError()
		}

		session, err = r.createOrUpdateSession(email, plaintextCode)
		if err != nil {
			r.logger.Error(fmt.Errorf("failed to create or update session: %w", err))
			return nil, commondomain.NewInternalError()
		}
	}

	return &RegisterAnswer{
		CodeExpires: session.CodeExpires.Unix(),
		CodePattern: r.CodeGenConfig.CodePattern,
	}, nil
}

func (r *RegistrationService) createOrUpdateSession(email string, code string) (*domain.RegistrationSession, error) {
	session, err := r.registrationSessionRepository.FindByEmail(email)

	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	originalCode := code
	if code != "" {
		code, err = hash.HashCode(code)
		if err != nil {
			r.logger.Error(fmt.Errorf("failed to hash verification code: %w", err))
			return nil, commondomain.NewInternalError()
		}
	}

	if session == nil {
		session = &domain.RegistrationSession{
			Code:        code,
			Email:       email,
			CodeExpires: time.Now().Add(r.CodeGenConfig.CodeTTL),
			CreateAt:    time.Now(),
		}
	} else if session.CodeExpires.Before(time.Now()) {
		session.Code = code
		session.CodeExpires = time.Now().Add(r.CodeGenConfig.CodeTTL)
	} else {
		return session, nil
	}

	if err := r.registrationSessionRepository.Save(session); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	debugSession := *session
	debugSession.Code = originalCode
	r.logger.Debug(debugSession)

	return session, nil
}

func (r *RegistrationService) validateConfirmationCode(email string, userData map[string]any) error {
	session, err := r.registrationSessionRepository.FindByEmail(email)
	if err != nil {
		r.logger.Error(err)
		return commondomain.NewInternalError()
	}
	if session == nil {
		return domain.NewInvalidCodeError()
	}
	if session.CodeExpires.Before(time.Now()) {
		return domain.NewCodeExpiredError()
	}

	submittedCode, ok := userData["code"].(string)
	if !ok || submittedCode == "" || !hash.VerifyCode(session.Code, submittedCode) {
		return domain.NewInvalidCodeError()
	}

	return nil
}

func (r *RegistrationService) ConfirmRegistration(userData map[string]any, lang string) error {
	grpcMap, err := r.validateRegistrationData(userData, lang)
	if err != nil {
		return err
	}

	email, ok := userData["email"].(string)
	if !ok {
		r.logger.Error(errors.New("email not found in userData"))
		return commondomain.NewInternalError()
	}

	if err := r.validateConfirmationCode(email, userData); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	md := metadata.Pairs("lang", lang)
	ctx = metadata.NewOutgoingContext(ctx, md)

	createAccountResponse, err := r.accountServiceClient.CreateAccount(
		ctx,
		&grpcApi.CreateAccountRequest{UserData: grpcMap},
	)

	if err != nil {
		st, _ := status.FromError(err)
		r.logger.Error(fmt.Errorf("CreateAccount error: %v, grpc status: %v", err, st))
		return commondomain.NewInternalError()
	}

	if createAccountResponse.Error != nil {
		return domain.GrpcErrorMapToError(createAccountResponse.Error)
	}

	if err := r.registrationSessionRepository.DeleteByEmail(email); err != nil {
		r.logger.Error(err)
		return commondomain.NewInternalError()
	}

	r.logger.Info("Account successfully created for email=" + email)
	return nil
}
