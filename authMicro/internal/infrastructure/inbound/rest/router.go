// Package rest provides the HTTP REST API for the auth microservice.
//
// @title           StudBridge Auth API
// @version         1.0
// @description     Сервис аутентификации: регистрация, вход по OTP-коду, управление JWT токенами.
// @BasePath        /auth/v1
//
// @contact.name   Pavel Sheludyakov
// @contact.url    https://github.com/PavelShe11/studBridge
//
// @tag.name registration
// @tag.description Двухшаговая регистрация с OTP-кодом
//
// @tag.name login
// @tag.description Двухшаговый вход с OTP-кодом
//
// @tag.name tokens
// @tag.description Управление JWT токенами
package rest

import (
	"context"
	"os"

	handler2 "github.com/PavelShe11/studbridge/authMicro/internal/infrastructure/inbound/rest/handler"
	httpErrorHandler2 "github.com/PavelShe11/studbridge/authMicro/internal/infrastructure/inbound/rest/httpErrorHandler"
	mymiddleware "github.com/PavelShe11/studbridge/authMicro/internal/infrastructure/inbound/rest/middleware"
	"github.com/PavelShe11/studbridge/common/logger"
	"github.com/PavelShe11/studbridge/common/translator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	e *echo.Echo
}

func NewRouter(
	log logger.Logger,
	translator *translator.Translator,
	regHandler *handler2.Register,
	loginHandler *handler2.Login,
	refreshTokenHandler *handler2.RefreshToken,
) *Router {
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler2.NewHttpErrorHandler(
		httpErrorHandler2.NewBaseErrorHandler(translator, log),
	)

	e.Use(mymiddleware.RequestLogger(log))
	if os.Getenv("LogLevel") == "debug" {
		e.Use(mymiddleware.RequestBodyLogger(log))
	}
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	v1 := e.Group("/auth/v1")
	v1.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	v1.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.URL("/swagger/doc.json")))

	login := v1.Group("/login")
	login.POST("/sendCodeEmail", loginHandler.SendLoginCode)
	login.POST("/confirmEmail", loginHandler.ConfirmEmail)

	registration := v1.Group("/registration")
	registration.POST("", regHandler.SendRegistrationCode)
	registration.POST("/confirmEmail", regHandler.RegistrationConfirmEmail)

	v1.POST("/refreshToken", refreshTokenHandler.RefreshToken)

	return &Router{
		e: e,
	}
}

func (r *Router) Start(address string) error {
	return r.e.Start(address)
}

func (r *Router) Shutdown(ctx context.Context) error {
	return r.e.Shutdown(ctx)
}
