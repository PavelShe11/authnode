package handler

import (
	"github.com/PavelShe11/studbridge/authMicro/internal/service"
	"github.com/PavelShe11/studbridge/common/logger"

	"net/http"

	"github.com/labstack/echo/v4"
)

type Login struct {
	logger       logger.Logger
	loginService *service.LoginService
	tokenService *service.TokenService
}

func NewLoginHandler(logger logger.Logger, loginService *service.LoginService, tokenService *service.TokenService) *Login {
	return &Login{
		logger:       logger,
		loginService: loginService,
		tokenService: tokenService,
	}
}

func (h *Login) SendLoginCode(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		h.logger.Error(err)
		return err
	}

	email, ok := req["email"].(string)
	if !ok {
		email = ""
	}
	var answer, err = h.loginService.Login(c.Request().Context(), email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, answer)
}

func (h *Login) ConfirmEmail(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		h.logger.Error(err)
		return err
	}

	email, ok := req["email"].(string)
	if !ok {
		email = ""
	}
	code, ok := req["code"].(string)
	if !ok {
		code = ""
	}

	accountId, err := h.loginService.ConfirmLogin(c.Request().Context(), email, code)
	if err != nil {
		return err
	}

	tokens, err := h.tokenService.CreateTokens(c.Request().Context(), accountId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tokens)
}
