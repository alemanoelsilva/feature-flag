package utils

import (
	"errors"
	auth "ff/internal/auth"

	"github.com/labstack/echo/v4"
)

func GetAuthenticatedPerson(c echo.Context, personId *int) error {
	authInfo := c.Get("auth_info").(auth.AuthUserResponse)

	if authInfo.PersonID == 0 {
		return errors.New("you are not logged in")
	}

	*personId = authInfo.PersonID

	return nil
}
