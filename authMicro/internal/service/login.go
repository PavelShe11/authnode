package service

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/PavelShe11/studbridge/auth/internal/config"
	"github.com/PavelShe11/studbridge/auth/internal/domain"
	"github.com/PavelShe11/studbridge/auth/internal/repository"
	"github.com/PavelShe11/studbridge/auth/utlis/generator"
	"github.com/PavelShe11/studbridge/auth/utlis/hash"
	"github.com/PavelShe11/studbridge/authMicro/grpcApi"
	commondomain "github.com/PavelShe11/studbridge/common/domain"
	"github.com/PavelShe11/studbridge/common/logger"
	"google.golang.org/grpc/status"
)

/**
TODO: Вход когда аккаунта ещё нет (фейковая сессия)
TODO: Вход когда аккаунт существует (реальная сессия)
TODO: Смена аккаунта
- Пользователь создаёт сессию входа
- Не завершая её удаляет аккаунт (через другую сессию)
- Создаёт новый аккаунт с тем же email
*/

type LoginAnswer struct {
	CodeExpires int64  `json:"code_expires"`
	CodePattern string `json:"code_pattern"`
}

type ConfirmLoginEmailAnswer struct {
	accessToken  string
	accessTTL    int
	refreshToken string
	refreshTTL   int
}

type LoginService struct {
	loginSessionRepository repository.LoginSessionRepository
	accountService         grpcApi.AccountServiceClient
	logger                 logger.Logger
	CodeGenConfig          *config.CodeGenConfig
}

func NewLoginService(
	loginSessionRepository repository.LoginSessionRepository,
	accountService grpcApi.AccountServiceClient,
	logger logger.Logger,
	codeGenConfig *config.CodeGenConfig,
) *LoginService {
	return &LoginService{
		loginSessionRepository: loginSessionRepository,
		accountService:         accountService,
		logger:                 logger,
		CodeGenConfig:          codeGenConfig,
	}
}

func (l *LoginService) cleanupExpiredSessions() {
	cleanupCtx, cancelCleanup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCleanup()
	if err := l.loginSessionRepository.CleanExpired(cleanupCtx); err != nil {
		l.logger.Error(fmt.Errorf("error cleaning expired login sessions: %w", err))
	}
}

func (l *LoginService) getAccountByEmail(email string) (*grpcApi.GetAccountResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accountGrpc, err := l.accountService.GetAccountByEmail(
		ctx,
		&grpcApi.GetAccountByEmailRequest{Email: email},
	)

	if err != nil {
		st, _ := status.FromError(err)
		l.logger.Error(fmt.Errorf("GetAccountByEmail error: %v, grpc status: %v", err, st))
		return nil, commondomain.NewInternalError()
	}
	return accountGrpc, nil
}

func (l *LoginService) validateEmail(email string) error {
	if email == "" {
		validationError := domain.NewValidationError()
		validationError.FieldErrors = append(validationError.FieldErrors, commondomain.FieldError{
			NameField: "email",
			Message:   "required",
			Params:    nil,
		})
		return validationError
	}

	if _, err := mail.ParseAddress(email); err != nil {
		validationError := domain.NewValidationError()
		validationError.FieldErrors = append(validationError.FieldErrors, commondomain.FieldError{
			NameField: "email",
			Message:   "email",
			Params:    nil,
		})
		return validationError
	}

	return nil
}

func (l *LoginService) createOrUpdateSession(email string, accountId *string, code string) (*domain.LoginSession, error) {
	session, err := l.loginSessionRepository.FindByEmail(email)
	if err != nil {
		l.logger.Error(err)
		return nil, err
	}

	originalCode := code
	if code != "" {
		code, err = hash.HashCode(code)
		if err != nil {
			l.logger.Error(fmt.Errorf("failed to hash verification code: %w", err))
			return nil, commondomain.NewInternalError()
		}
	}

	if session == nil {
		session = &domain.LoginSession{
			AccountId:   accountId,
			Email:       email,
			Code:        code,
			CodeExpires: time.Now().Add(l.CodeGenConfig.CodeTTL),
			CreateAt:    time.Now(),
		}
	} else {
		accountIdChanged := (session.AccountId == nil && accountId != nil) ||
			(session.AccountId != nil && accountId == nil) ||
			(session.AccountId != nil && accountId != nil && *session.AccountId != *accountId)

		if session.CodeExpires.Before(time.Now()) || accountIdChanged {
			session.AccountId = accountId
			session.Code = code
			session.CodeExpires = time.Now().Add(l.CodeGenConfig.CodeTTL)
		} else {
			return session, nil
		}
	}

	if err := l.loginSessionRepository.Save(session); err != nil {
		l.logger.Error(err)
		return nil, err
	}

	debugSession := *session
	debugSession.Code = originalCode
	l.logger.Debug(debugSession)

	return session, nil
}

func (l *LoginService) Login(email string) (*LoginAnswer, error) {
	l.cleanupExpiredSessions()

	if err := l.validateEmail(email); err != nil {
		return nil, err
	}

	accountGrpc, err := l.getAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	var accountId *string
	if account, ok := accountGrpc.Result.(*grpcApi.GetAccountResponse_Account); ok && account != nil && account.Account != nil {
		if account.Account.AccountId != "" {
			accountId = &account.Account.AccountId
		}
	}

	var session *domain.LoginSession

	if accountId != nil {
		plaintextCode, err := generator.Reggen(l.CodeGenConfig.CodePattern, l.CodeGenConfig.CodeMaxLength)
		if err != nil {
			l.logger.Error(err)
			return nil, commondomain.NewInternalError()
		}

		session, err = l.createOrUpdateSession(email, accountId, plaintextCode)
		if err != nil {
			l.logger.Error(fmt.Errorf("failed to create or update login session: %w", err))
			return nil, commondomain.NewInternalError()
		}
	} else {
		session, err = l.createOrUpdateSession(email, nil, "")
		if err != nil {
			l.logger.Error(err)
			return nil, err
		}
	}

	return &LoginAnswer{
		CodeExpires: session.CodeExpires.Unix(),
		CodePattern: l.CodeGenConfig.CodePattern,
	}, nil
}

func (l *LoginService) ConfirmLoginEmail(email string, code string, lang string) (*ConfirmLoginEmailAnswer, error) {
	return nil, nil
}
