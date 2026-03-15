package service

import (
	"context"

	commonEntity "github.com/PavelShe11/authnode/common/entity"
	"github.com/PavelShe11/authnode/common/logger"
	"github.com/PavelShe11/authnode/common/validation"
	"github.com/PavelShe11/authnode/user/internal/entity"
	"github.com/PavelShe11/authnode/user/internal/port"
)

type AccountService struct {
	accountRepository port.AccountRepository
	logger            logger.Logger
	validator         *validation.Validator
}

func NewAccountService(
	accountRepository port.AccountRepository,
	l logger.Logger,
	validator *validation.Validator,
) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
		logger:            l,
		validator:         validator,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, account entity.Account) error {
	errs := s.ValidateAccountData(account)
	if errs != nil {
		return errs
	}
	err := s.accountRepository.CreateAccount(ctx, account)
	if err != nil {
		s.logger.Error(err)
		return commonEntity.NewInternalError()
	}
	return nil
}

func (s *AccountService) GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error) {
	account, err := s.accountRepository.GetAccountByEmail(ctx, email)
	if account == nil && err == nil {
		return nil, nil
	}
	if err != nil {
		s.logger.Error(err)
		return nil, commonEntity.NewInternalError()
	}
	return account, nil
}

func (s *AccountService) GetAccountById(ctx context.Context, id string) (*entity.Account, error) {
	account, err := s.accountRepository.GetAccountById(ctx, id)
	if err != nil {
		s.logger.Error(err)
		return nil, commonEntity.NewInternalError()
	}
	return account, nil
}

func (s *AccountService) ValidateAccountData(account entity.Account) error {
	errs := commonEntity.NewValidationError()
	errs.FieldErrors = s.validator.Struct(&account)
	if len(errs.FieldErrors) > 0 {
		return errs
	}
	return nil
}
