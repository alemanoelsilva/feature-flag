package utils

import (
	"encoding/json"
	"io"

	"github.com/labstack/echo/v4"
)

func GetBodyFromRequest[T any](c echo.Context, input *T) error {
	// manually decoding the json body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	return nil
}
