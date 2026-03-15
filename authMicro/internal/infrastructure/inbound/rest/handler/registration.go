package handler

import (
	"net/http"

	"github.com/PavelShe11/authnode/authMicro/internal/infrastructure/inbound/rest/httpErrorHandler"
	"github.com/PavelShe11/authnode/authMicro/internal/infrastructure/inbound/rest/models"
	"github.com/PavelShe11/authnode/authMicro/internal/service"
	"github.com/PavelShe11/authnode/common/logger"
	"github.com/PavelShe11/authnode/common/translator" // Added translator import

	"github.com/labstack/echo/v4"
)

type Register struct {
	logger              logger.Logger
	registrationService *service.RegistrationService
	translator          *translator.Translator // Added translator field
}

func NewRegisterHandler(
	logger logger.Logger,
	registrationService *service.RegistrationService,
	translator *translator.Translator, // Added translator parameter
) *Register {
	return &Register{
		logger:              logger,
		registrationService: registrationService,
		translator:          translator, // Assign translator
	}
}

// SendRegistrationCode godoc
// @Summary      Отправить код регистрации
// @Description  Создаёт или продлевает сессию регистрации и генерирует OTP-код.
// @Description  Набор полей в теле запроса определяется правилами валидации user-service (ValidateAccountData).
// @Description  Auth service передаёт данные напрямую — без фиксированной схемы.
// @Tags         registration
// @Accept       json
// @Produce      json
// @Param        request  body      object                    true  "Данные пользователя (схема определяется user-service)"
// @Success      200      {object}  models.RegistrationResponse
// @Failure      400      {object}  entity.BaseValidationError
// @Failure      500      {object}  entity.BaseError
// @Router       /registration [post]
func (h *Register) SendRegistrationCode(c echo.Context) error {
	var req map[string]any
	if err := c.Bind(&req); err != nil {
		h.logger.Error(err)
		return err
	}

	lang := httpErrorHandler.GetLangFromHeader(c)
	answer, err := h.registrationService.Register(c.Request().Context(), req, lang)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.NewRegistrationResponse(answer))
}

// RegistrationConfirmEmail godoc
// @Summary      Подтвердить email при регистрации
// @Description  Проверяет OTP-код из сессии регистрации, затем передаёт данные в user-service для создания аккаунта.
// @Description  Набор полей определяется user-service. Обязательные поля: email (string), code (string).
// @Tags         registration
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Email, OTP-код и данные пользователя (схема определяется user-service)"
// @Success      200      "Регистрация успешно завершена"
// @Failure      400      {object}  entity.BaseValidationError
// @Failure      500      {object}  entity.BaseError
// @Router       /registration/confirmEmail [post]
func (h *Register) RegistrationConfirmEmail(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		h.logger.Error(err)
		return err
	}

	lang := httpErrorHandler.GetLangFromHeader(c)
	err := h.registrationService.ConfirmRegistration(c.Request().Context(), req, lang)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
