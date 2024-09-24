package http

import (
	"encoding/json"
	"errors"
	"ff/internal/feature_flag/entity"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func getPersonIdFromHeaders(c echo.Context) (int, error) {
	personIdStr := c.Request().Header.Get("Personid")
	if personIdStr == "" {
		return 0, errors.New("missing Personid header")
	}

	personId, err := strconv.Atoi(personIdStr)
	if err != nil {
		return 0, errors.New("invalid Personid format")
	}

	if personId == 0 {
		return 0, errors.New("you are not logged in")
	}

	return personId, nil
}

func LoadFeatureFlagsRoutes(router *echo.Echo, handler *EchoHandler) {
	router.POST("/api/feature-flags/v1/feature-flags", handler.createFeatureFlagHandler)
}

func (e *EchoHandler) createFeatureFlagHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	// bind the json body to the struct, but it's not working when the body comes from the test case
	// var input entity.FeatureFlag
	// if err := c.Bind(&input); err != nil {
	// 	return response.ErrorHandler(http.StatusBadRequest, err)
	// }

	// manually decode the json body
	var input entity.FeatureFlag
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.ErrorHandler(http.StatusBadRequest, err)
	}
	if err := json.Unmarshal(body, &input); err != nil {
		return response.ErrorHandler(http.StatusBadRequest, err)
	}

	if err := input.Validate(); err != nil {
		return response.ErrorHandler(http.StatusBadRequest, err)
	}

	personId, err := getPersonIdFromHeaders(c)
	if err != nil {
		return response.ErrorHandler(http.StatusUnauthorized, err)
	}

	if err := e.FeatureFlagService.CreateFeatureFlag(input, personId); err != nil {
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	return response.SuccessHandler(http.StatusCreated, handleResponseMessage("Feature Flag Created"))
}
