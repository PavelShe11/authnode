package httpErrorHandler

import (
	"errors"
	"net/http"

	commonEntity "github.com/PavelShe11/authnode/common/entity"
	"github.com/PavelShe11/authnode/common/logger"
	"github.com/PavelShe11/authnode/common/translator"

	"github.com/labstack/echo/v4"
)

type baseErrorhandler struct {
	translator *translator.Translator
	log        logger.Logger
}

func NewBaseErrorHandler(translator *translator.Translator, log logger.Logger) DomainErrorHandler {
	return &baseErrorhandler{
		translator: translator,
		log:        log,
	}
}

func (h *baseErrorhandler) handle(err error, c echo.Context) bool {
	var domainErr commonEntity.AbstractError
	ok := errors.As(err, &domainErr)
	if !ok {
		return false
	}

	statusCode, err := getStatusCodeForBaseError(domainErr.GetCode())
	if err != nil {
		statusCode = http.StatusInternalServerError
		domainErr = commonEntity.NewInternalError()
	}

	lang := GetLangFromHeader(c)

	h.translator.TranslateError(domainErr, lang)

	if err := c.JSON(statusCode, domainErr); err != nil {
		h.log.Error("Failed to send error response", "error", err)
	}

	return true
}

func getStatusCodeForBaseError(base string) (int, error) {
	switch base {
	case "internalError":
		return http.StatusInternalServerError, nil
	case "invalidCode", "codeExpired", "validationError":
		return http.StatusBadRequest, nil
	default:
		return 0, errors.New("no mapping was added to the http code error for the error")
	}
}
