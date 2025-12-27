package handler

import (
	"github.com/PavelShe11/studbridge/auth/internal/api/rest/httpErrorHandler"
	"github.com/PavelShe11/studbridge/auth/internal/service"
	"github.com/PavelShe11/studbridge/common/logger"

	"net/http"

	"github.com/labstack/echo/v4"
)

type Login struct {
	logger       logger.Logger
	loginService *service.LoginService
}

func NewLoginHandler(logger logger.Logger, loginService *service.LoginService) *Login {
	return &Login{
		logger:       logger,
		loginService: loginService,
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
	var answer, err = h.loginService.Login(email)
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

	lang := httpErrorHandler.GetLangFromHeader(c)
	email, ok := req["email"].(string)
	if !ok {
		email = ""
	}
	code, ok := req["code"].(string)
	if !ok {
		code = ""
	}

	var answer, err = h.loginService.ConfirmLoginEmail(email, code, lang)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, answer)
}
