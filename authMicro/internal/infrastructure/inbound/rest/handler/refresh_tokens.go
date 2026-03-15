package handler

import (
	"errors"

	"github.com/PavelShe11/authnode/authMicro/internal/infrastructure/inbound/rest/models"
	"github.com/PavelShe11/authnode/authMicro/internal/service"
	"github.com/PavelShe11/authnode/common/logger"

	"net/http"

	"github.com/labstack/echo/v4"
)

type RefreshToken struct {
	logger       logger.Logger
	tokenService *service.TokenService
}

func NewRefreshTokenHandler(logger logger.Logger, tokenService *service.TokenService) *RefreshToken {
	return &RefreshToken{
		logger:       logger,
		tokenService: tokenService,
	}
}

// RefreshToken godoc
// @Summary      Обновить токены
// @Description  Проверяет refresh токен и возвращает новую пару JWT токенов (access + refresh). Старый refresh токен инвалидируется.
// @Description  Обязательное поле: refreshToken (string).
// @Tags         tokens
// @Accept       json
// @Produce      json
// @Param        request  body      object                true  "Refresh токен"
// @Success      200      {object}  models.TokensResponse
// @Failure      400      {object}  entity.BaseError
// @Failure      401      "Токен недействителен или истёк"
// @Failure      500      {object}  entity.BaseError
// @Router       /refreshToken [post]
func (h *RefreshToken) RefreshToken(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	refreshToken, ok := req["refreshToken"].(string)
	if !ok || refreshToken == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	tokens, err := h.tokenService.RefreshTokens(c.Request().Context(), refreshToken)
	if err != nil {
		if errors.Is(err, service.InvalidRefreshTokenError) ||
			errors.Is(err, service.RefreshTokenExpiredError) ||
			errors.Is(err, service.UnauthorizedRefreshTokenError) {

			return c.NoContent(http.StatusUnauthorized)
		}

		return err
	}

	return c.JSON(http.StatusOK, models.NewTokensResponse(tokens))
}
