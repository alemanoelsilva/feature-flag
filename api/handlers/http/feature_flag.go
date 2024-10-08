package http

import (
	"errors"
	"ff/api/middlewares"
	featureflag "ff/internal/feature_flag"
	"ff/internal/feature_flag/entity"
	"ff/pkg/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FeatureFlagEchoHandler struct {
	FeatureFlagService featureflag.FeatureFlagService
}

func NewFeatureFlagEchoHandler(featureflag featureflag.FeatureFlagService) *echo.Echo {
	handler := &FeatureFlagEchoHandler{
		FeatureFlagService: featureflag,
	}

	e := echo.New()

	e.Use(middlewares.LoggerMiddleware())

	LoadFeatureFlagsRoutes(e, handler)

	return e
}

func LoadFeatureFlagsRoutes(e *echo.Echo, handler *FeatureFlagEchoHandler) {
	group := e.Group("/api/feature-flags")

	group.POST("/v1/feature-flags", handler.createFeatureFlagHandler)
	group.GET("/v1/feature-flags", handler.getFeatureFlagHandler)
	group.PUT("/v1/feature-flags/:id", handler.updateFeatureFlagByIdHandler)
}

func (e *FeatureFlagEchoHandler) createFeatureFlagHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	var input entity.FeatureFlag
	if err := utils.GetBodyFromRequest(c, &input); err != nil {
		return response.ErrorHandler(http.StatusBadRequest, err)
	}

	var personId int
	if err := utils.GetAuthenticatedPerson(c, &personId); err != nil {
		return response.ErrorHandler(http.StatusUnauthorized, err)
	}

	if err := e.FeatureFlagService.CreateFeatureFlag(input, uint(personId)); err != nil {
		if err.Error() == "feature flag already exists" {
			return response.ErrorHandler(http.StatusConflict, err)
		}
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	return response.SuccessHandlerMessage(http.StatusCreated, "Feature Flag Created")
}

func (e *FeatureFlagEchoHandler) getFeatureFlagHandler(c echo.Context) error {
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

func (e *FeatureFlagEchoHandler) updateFeatureFlagByIdHandler(c echo.Context) error {
	response := ResponseJSON{c: c}

	var input entity.UpdateFeatureFlag
	if err := utils.GetBodyFromRequest(c, &input); err != nil {
		return response.ErrorHandler(http.StatusBadRequest, err)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.ErrorHandler(http.StatusBadRequest, errors.New("feature flag id is not a number"))
	}

	// TODO: get personId to audit
	var personId int
	if err = utils.GetAuthenticatedPerson(c, &personId); err != nil {
		return response.ErrorHandler(http.StatusUnauthorized, err)
	}

	if err := e.FeatureFlagService.UpdateFeatureFlagById(uint(id), input); err != nil {
		if err.Error() == "no feature flag updated" {
			return response.SuccessHandlerMessage(http.StatusOK, "no feature flag updated")
		}
		if err.Error() == "feature flag not found" {
			return response.ErrorHandler(http.StatusNotFound, err)
		}
		return response.ErrorHandler(http.StatusInternalServerError, err)
	}

	return response.SuccessHandlerMessage(http.StatusOK, "Feature Flag Updated")
}
