package http

import (
	featureflag "ff/internal/feature-flag"

	"github.com/labstack/echo/v4"
)

type EchoHandler struct {
	FeatureFlagService featureflag.FeatureFlagService
}

func NewEchoHandler(featureflag featureflag.FeatureFlagService) *echo.Echo {
	handler := &EchoHandler{
		FeatureFlagService: featureflag,
	}

	router := echo.New()

	LoadFeatureFlagsRoutes(router, handler)

	return router
}

func handleResponseMessage(msg string) interface{} {
	return map[string]interface{}{"message": msg}
}

type ResponseJSON struct {
	c echo.Context
}

func (s ResponseJSON) SuccessHandler(code int, data interface{}) error {
	return s.c.JSON(code, data)
}

func (s ResponseJSON) ErrorHandler(code int, err error) error {
	return s.c.JSON(code, map[string]interface{}{"error": err.Error()})
}
