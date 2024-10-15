package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetNumberParamFromRequest(keyName string, c echo.Context) (int, error) {
	valueString := c.Param(keyName)
	value, err := strconv.Atoi(valueString)
	if err != nil {
		return 0, err
	}
	return value, nil
}
