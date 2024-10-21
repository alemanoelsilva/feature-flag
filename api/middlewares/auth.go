package middlewares

import (
	auth "ff/internal/auth"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ValidateCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the cookie from the request
		cookie := c.Request().Header.Get("Cookie")

		if cookie == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing auth cookie")
		}

		authInfo, err := auth.GetAuthInfo(cookie)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid auth cookie")
		}
		if !authInfo.IsAdmin {
			return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to perform this action")
		}

		fmt.Printf("AUTH %v", authInfo)

		c.Set("auth_info", authInfo)

		return next(c)
	}
}
