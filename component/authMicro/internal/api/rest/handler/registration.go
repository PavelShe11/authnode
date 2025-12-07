package handler

import (
	"authMicro/utlis/logger"

	"net/http"

	"github.com/labstack/echo/v4"
)

type Register struct {
	logger logger.Logger
}

func NewRegisterHandler(logger logger.Logger) *Register {
	return &Register{
		logger: logger,
	}
}

func (h *Register) SendRegistrationCode(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, req)
}

func (h *Register) RegistrationConfirmEmail(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, req)
}
