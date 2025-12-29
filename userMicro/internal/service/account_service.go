package service

import (
	commonEntity "github.com/PavelShe11/studbridge/common/entity"
	"github.com/PavelShe11/studbridge/common/logger"
	"github.com/PavelShe11/studbridge/common/validation"
	"github.com/PavelShe11/studbridge/user/internal/entity"
	"github.com/PavelShe11/studbridge/user/internal/repository"
)

type AccountService struct {
	accountRepository *repository.AccountRepository
	logger            logger.Logger
	validator         *validation.Validator
}

func NewAccountService(
	accountRepository *repository.AccountRepository,
	l logger.Logger,
	validator *validation.Validator,
) *AccountService {
	return &AccountService{
		accountRepository: accountRepository,
		logger:            l,
		validator:         validator,
	}
}

func (s *AccountService) CreateAccount(account entity.Account) error {
	errs := s.ValidateAccountData(account)
	if errs != nil {
		return errs
	}
	err := s.accountRepository.CreateAccount(account)
	if err != nil {
		s.logger.Error(err)
		return commonEntity.NewInternalError()
	}
	return nil
}

func (s *AccountService) GetAccountByEmail(email string) (*entity.Account, error) {
	account, err := s.accountRepository.GetAccountByEmail(email)
	if account == nil && err == nil {
		return nil, nil
	}
	if err != nil {
		s.logger.Error(err)
		return nil, commonEntity.NewInternalError()
	}
	return account, nil
}

func (s *AccountService) GetAccountById(id string) (*entity.Account, error) {
	account, err := s.accountRepository.GetAccountById(id)
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
