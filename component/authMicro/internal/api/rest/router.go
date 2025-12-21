package rest

import (
	"authMicro/internal/api/rest/handler"
	"authMicro/utlis/logger"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	e *echo.Echo
}

func NewRouter(
	log logger.Logger,
	regHandler *handler.Register,
	loginHandler *handler.Login,
	refreshTokenHandler *handler.RefreshToken,
) *Router {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogHeaders:  []string{"Content-Type", "User-Agent"},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.Infof("%s %s %d %s %s", v.Method, v.URI, v.Status, v.Latency, v.RemoteIP)
			} else {
				log.Errorf("%s %s %d %s %s error=%v", v.Method, v.URI, v.Status, v.Latency, v.RemoteIP, v.Error)
			}
			return nil
		},
	}))
	if os.Getenv("LogLevel") == "debug" {
		e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			maxSize := 1024
			if len(reqBody) > 0 {
				reqStr := formatBodyForLog(reqBody, maxSize)
				log.Debugf("Request Body:\n%s", reqStr)
			}
			if len(resBody) > 0 {
				resStr := formatBodyForLog(resBody, maxSize)
				log.Debugf("Response Body:\n%s", resStr)
			}
		}))
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

// formatBodyForLog форматирует body для логирования.
// Если это JSON - форматирует с отступами, если нет - возвращает как есть
func formatBodyForLog(body []byte, maxSize int) string {
	if len(body) == 0 {
		return ""
	}
	var js interface{}
	if err := json.Unmarshal(body, &js); err != nil {
		// Не JSON - вернуть как есть (убрать trailing newlines)
		str := strings.TrimRight(string(body), "\n\r")
		if len(str) > maxSize {
			return str[:maxSize] + "..."
		}
		return str
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, body, "", "  "); err != nil {
		str := strings.TrimRight(string(body), "\n\r")
		if len(str) > maxSize {
			return str[:maxSize] + "..."
		}
		return str
	}
	formatted := strings.TrimRight(buf.String(), "\n\r")
	if len(formatted) > maxSize {
		return formatted[:maxSize] + "..."
	}
	return formatted
}
