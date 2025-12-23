package handler

import (
	"authMicro/internal/service"
	"authMicro/utlis/logger"
	"authMicro/utlis/translator" // Added translator import
	"net/http"
	"strings" // Added strings import

	"github.com/labstack/echo/v4"
)

type Register struct {
	logger              logger.Logger
	registrationService service.RegistrationService
	translator          *translator.Translator // Added translator field
}

func NewRegisterHandler(
	logger logger.Logger,
	registrationService service.RegistrationService,
	translator *translator.Translator, // Added translator parameter
) *Register {
	return &Register{
		logger:              logger,
		registrationService: registrationService,
		translator:          translator, // Assign translator
	}
}

func (h *Register) SendRegistrationCode(c echo.Context) error {
	var req map[string]any
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	lang := getLangFromHeader(c)                             // Get language from header
	answer, err := h.registrationService.Register(req, lang) // Pass language to service
	if err != nil {
		// Translate error before returning
		h.translator.TranslateError(err, lang)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, answer)
}

func (h *Register) RegistrationConfirmEmail(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	lang := getLangFromHeader(c)                                // Get language from header
	err := h.registrationService.ConfirmRegistration(req, lang) // Pass language to service
	if err != nil {
		// Translate error before returning
		h.translator.TranslateError(err, lang)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.NoContent(http.StatusOK)
}

// getLangFromHeader extracts language preferences from the Accept-Language header.
// It returns the preferred language or "en" as default.
func getLangFromHeader(c echo.Context) string {
	acceptLanguage := c.Request().Header.Get("Accept-Language")
	if acceptLanguage == "" {
		return "en" // Default language
	}

	// Basic parsing: split by comma and take the first one.
	// In a real app, you'd use a more robust Accept-Language parser.
	langs := strings.Split(acceptLanguage, ",")
	if len(langs) > 0 {
		return strings.TrimSpace(strings.Split(langs[0], ";")[0])
	}

	return "en" // Fallback
}
