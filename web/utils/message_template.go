package utils

import (
	"ff/web/components"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorMessage(c echo.Context, msg string) error {
	c.Response().Header().Add("HX-Retarget", "#message")
	c.Response().Header().Add("HX-Reswap", "outerHTML")
	return Render(c, http.StatusConflict, components.Message(true, msg, true))
}
