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

func getPersonIdFromHeaders(c echo.Context) (uint, error) {
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

	return uint(personId), nil
}

func getBodyFromRequest[T any](c echo.Context, input *T) error {
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

func LoadFeatureFlagsRoutes(e *echo.Echo, handler *EchoHandler) {
	group := e.Group("/api/feature-flags")

	group.POST("/v1/feature-flags", handler.createFeatureFlagHandler)
	group.GET("/v1/feature-flags", handler.getFeatureFlagHandler)
	group.PUT("/v1/feature-flags/:id", handler.updateFeatureFlagByIdHandler)
}

func (e *EchoHandler) createFeatureFlagHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	var input entity.FeatureFlag
	if err := getBodyFromRequest(c, &input); err != nil {
		return response.ErrorHandler(http.StatusBadRequest, err)
	}

	personId, err := getPersonIdFromHeaders(c)
	if err != nil {
		return response.ErrorHandler(http.StatusUnauthorized, err)
	}

	if err := e.FeatureFlagService.CreateFeatureFlag(input, personId); err != nil {
		if err.Error() == "feature flag already exists" {
			return response.ErrorHandler(http.StatusConflict, err)
		}
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	return response.SuccessHandler(http.StatusCreated, handleResponseMessage("Feature Flag Created"))
}

func (e *EchoHandler) getFeatureFlagHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	id, _ := strconv.Atoi(c.QueryParam("id"))
	personId, _ := strconv.Atoi(c.QueryParam("personId"))
	name := c.QueryParam("name")

	// TODO: check it again, it is terrible
	isActiveStr := c.QueryParam("isActive")
	var isActive *bool
	if isActiveStr == "true" {
		trueValue := true
		isActive = &trueValue
	} else if isActiveStr == "false" {
		falseValue := false
		isActive = &falseValue
	} else if isActiveStr != "" {
		return response.ErrorHandler(http.StatusBadRequest, errors.New("invalid isActive value"))
	}

	if page <= 1 {
		page = 1 // Default page
	}
	if limit <= 0 {
		limit = 10 // Default limit
	}

	featureFlag, totalCount, err := e.FeatureFlagService.GetFeatureFlag(page, limit, name, isActive, uint(id), uint(personId))
	if err != nil {
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	// TODO: check it again, it is terrible
	// Convert featureFlag to []interface{}
	interfaceSlice := make([]interface{}, len(featureFlag))
	for i, v := range featureFlag {
		interfaceSlice[i] = v
	}

	return response.PaginationHandler(interfaceSlice, totalCount)
}

func (e *EchoHandler) updateFeatureFlagByIdHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	var input entity.UpdateFeatureFlag
	if err := getBodyFromRequest(c, &input); err != nil {
		return response.ErrorHandler(http.StatusBadRequest, err)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorHandler(http.StatusBadRequest, errors.New("feature flag id is not a number"))
	}

	// TODO: get personId to audit
	_, err = getPersonIdFromHeaders(c)
	if err != nil {
		return response.ErrorHandler(http.StatusUnauthorized, err)
	}

	if err := e.FeatureFlagService.UpdateFeatureFlagById(uint(id), input); err != nil {
		if err.Error() == "no feature flag updated" {
			return response.ErrorHandler(http.StatusNotFound, err)
		}
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	return response.SuccessHandler(http.StatusOK, handleResponseMessage("Feature Flag Updated"))
}
