package service

import (
	"github.com/PavelShe11/studbridge/authMicro/grpcApi"
)

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
	accountService grpcApi.AccountServiceClient
}

func NewLoginService(accountService grpcApi.AccountServiceClient) LoginService {
	return LoginService{
		accountService: accountService,
	}
}

func (l *LoginService) Login(email string, lang string) (*LoginAnswer, error) {
	return nil, nil
}

func (l *LoginService) ConfirmLoginEmail(email string, code string, lang string) (*ConfirmLoginEmailAnswer, error) {
	return nil, nil
}
