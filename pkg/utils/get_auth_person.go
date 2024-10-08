package utils

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAuthenticatedPerson(c echo.Context, personId *int) error {
	personIdStr := c.Request().Header.Get("Personid")
	if personIdStr == "" {
		return errors.New("missing Personid header")
	}

	id, err := strconv.Atoi(personIdStr)
	if err != nil {
		return errors.New("invalid Personid format")
	}

	if id == 0 {
		return errors.New("you are not logged in")
	}

	*personId = id

	return nil
}
