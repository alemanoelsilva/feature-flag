package middlewares

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	logger := zerolog.New(os.Stdout)
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Status >= 400 {
				logger.Error().
					Err(v.Error).
					Str("URI", v.URI).
					Str("method", c.Request().Method).
					Int("status", v.Status).
					Msg("request error")
			} else {
				logger.Info().
					Str("URI", v.URI).
					Str("method", c.Request().Method).
					Int("status", v.Status).
					Msg("ok")
			}

			return nil
		},
	})
}
